// Copyright 2022-2023 Tigris Data, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package workers

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"github.com/tigrisdata/tigris/schema"
	"github.com/tigrisdata/tigris/server/metadata"
	"github.com/tigrisdata/tigris/server/services/v1/database"
	"github.com/tigrisdata/tigris/server/transaction"
	ulog "github.com/tigrisdata/tigris/util/log"
)

const (
	MAX_ERROR_COUNT = 10
	LEASE_TIME      = 1 * time.Minute
	PEAK_JOB_ITEMS  = 5 // number of items to fetch from the queue to see which one to select
)

type WorkerTestTask struct {
	Sleep            time.Duration `json:"sleep"`
	NumErrors        int           `json:"error_count"`
	ShouldStopWorker bool          `json:"should_stop_worker"`
}

type Worker struct {
	id        uint
	queue     *metadata.QueueSubspace
	txMgr     *transaction.Manager
	tenantMgr *metadata.TenantManager
	errCount  uint
	// Indicated to the worker to shutdown
	done chan struct{}
	// How long the worker sleeps between checking queue
	sleepTime    time.Duration
	itemEvent    chan<- Event
	hearbeatChan chan<- uint
}

type Event struct {
	Success  bool
	Item     metadata.QueueItem
	WorkerId uint
}

func newEvent(success bool, item metadata.QueueItem, workerId uint) Event {
	return Event{
		success,
		item,
		workerId,
	}
}

func newWorker(id uint, queue *metadata.QueueSubspace, txMgr *transaction.Manager, tenantMgr *metadata.TenantManager, sleepTime time.Duration, itemEvent chan<- Event, heartbeatChan chan<- uint) *Worker {
	return &Worker{
		id:           id,
		queue:        queue,
		txMgr:        txMgr,
		tenantMgr:    tenantMgr,
		errCount:     0,
		done:         make(chan struct{}, 1),
		sleepTime:    sleepTime,
		itemEvent:    itemEvent,
		hearbeatChan: heartbeatChan,
	}
}

// Add a slightly random jitter for each worker
// This helps avoid a situation where every worker
// starts or peaks at the exact time and all tries to read from the job queue.
func (w *Worker) jitterSleep() {
	//nolint:gosec // it is ok to use a weak random generator here
	jitterSleep := rand.Intn(500)
	time.Sleep(time.Duration(jitterSleep)*time.Millisecond + w.sleepTime)
}

func (w *Worker) Start() {
	w.jitterSleep()
	w.Loop()
}

func (w *Worker) Loop() {
	hbTicker := time.NewTicker(2 * w.sleepTime)
	c := 0
	for {
		select {
		case <-w.done:
			log.Info().Msgf("Worker %d shutting down", w.id)
			hbTicker.Stop()
			return
		case <-hbTicker.C:
			w.hearbeatChan <- w.id
		default:
			c++
			err := w.peekAndProcess()
			if err != nil {
				w.errCount++
				log.Err(err).Msgf("Worker %d error while processing", w.id)
				// if the worker is getting to many errors in a row
				// shut the worker down and restart a new one
				if w.errCount >= MAX_ERROR_COUNT {
					log.Err(err).Msgf("Worker %d exceeded error count and shutting down", w.id)
					return
				}
			} else {
				w.errCount = 0
			}
		}
		w.jitterSleep()
	}
}

func (w *Worker) Stop() {
	w.done <- struct{}{}
}

func (w *Worker) peekAndProcess() error {
	ctx := context.Background()
	tx, err := w.txMgr.StartTx(ctx)
	if err != nil {
		return err
	}
	items, err := w.queue.Peek(ctx, tx, PEAK_JOB_ITEMS)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return tx.Rollback(ctx)
	}
	selectedItem, err := w.queue.ObtainLease(ctx, tx, &items[0], LEASE_TIME)
	log.Debug().Msgf("Worker %d: processing task %s", w.id, selectedItem.Id)
	if err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}

	if err = w.processItem(selectedItem); err != nil {
		log.Err(err).Msgf("Worker %d: failed to process %s", w.id, selectedItem.Id)
		return w.handleFailedProcessing(ctx, selectedItem)
	}

	log.Info().Msgf("Worker %d: Completed %s", w.id, selectedItem.Id)
	w.itemEvent <- newEvent(true, *selectedItem, w.id)
	return nil
}

func (w *Worker) handleFailedProcessing(ctx context.Context, selectedItem *metadata.QueueItem) error {
	selectedItem.ErrorCount++
	tx, err := w.txMgr.StartTx(ctx)
	if err != nil {
		return err
	}
	// The queue item has had to many errors the job will remove it from the queue
	if selectedItem.ErrorCount >= MAX_ERROR_COUNT {
		log.Err(err).Msgf("Worker %d: Max fail count, dropping item %s from the queue", w.id, selectedItem.Id)
		w.itemEvent <- newEvent(false, *selectedItem, w.id)
		if err = w.queue.Dequeue(ctx, tx, selectedItem); err != nil {
			return err
		}
	} else {
		if err = w.queue.Requeue(ctx, tx, selectedItem, w.sleepTime*time.Duration(selectedItem.ErrorCount)); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (w *Worker) processItem(queueItem *metadata.QueueItem) error {
	switch queueItem.TaskType {
	case metadata.BUILD_INDEX_QUEUE_TASK:
		return w.buildIndexTask(queueItem)
	case metadata.TEST_QUEUE_TASK:
		return w.testQueueTask(queueItem)
	}

	return fmt.Errorf("unknown job type")
}

func (w *Worker) testQueueTask(item *metadata.QueueItem) error {
	ctx := context.Background()
	var testTask WorkerTestTask
	if err := jsoniter.Unmarshal(item.Data, &testTask); err != nil {
		return err
	}

	if testTask.ShouldStopWorker {
		testTask.ShouldStopWorker = false
		item.Data, _ = jsoniter.Marshal(testTask)
		w.Stop()
		return fmt.Errorf("forced worker stop %d", w.id)
	}

	if testTask.NumErrors > 0 {
		testTask.NumErrors--
		item.Data, _ = jsoniter.Marshal(testTask)
		return fmt.Errorf("test error generated %d", testTask.NumErrors)
	}

	time.Sleep(testTask.Sleep)

	tx, err := w.txMgr.StartTx(ctx)
	if err != nil {
		return err
	}
	if err = w.queue.Complete(ctx, tx, item); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (w *Worker) buildIndexTask(queueItem *metadata.QueueItem) error {
	var task metadata.IndexBuildTask
	if err := jsoniter.Unmarshal(queueItem.Data, &task); err != nil {
		return err
	}

	ctx := context.Background()
	dbBranch := metadata.NewDatabaseNameWithBranch(task.ProjName, task.Branch)
	tenant, _ := w.tenantMgr.GetTenant(ctx, task.NamespaceId)

	project, err := tenant.GetProject(task.ProjName)
	if err != nil {
		return err
	}

	db, err := project.GetDatabase(dbBranch)
	if err != nil {
		return err
	}

	coll := db.GetCollection(task.CollName)
	indexer := database.NewSecondaryIndexer(coll)

	// Extend the lease of the queue item so that another worker does not
	// try and process it
	progressUpdate := func(ctx context.Context, tx transaction.Tx) error {
		return w.queue.RenewLease(ctx, tx, queueItem, LEASE_TIME)
	}
	if err = indexer.BuildCollection(ctx, w.txMgr, progressUpdate); err != nil {
		return err
	}

	for _, index := range coll.SecondaryIndexes.All {
		index.State = schema.INDEX_ACTIVE
	}

	tx, err := w.txMgr.StartTx(ctx)
	if err != nil {
		return err
	}

	if err = tenant.UpdateCollectionIndexes(ctx, tx, db, coll.Name, coll.SecondaryIndexes.All); ulog.E(err) {
		return err
	}

	if err = w.queue.Complete(ctx, tx, queueItem); ulog.E(err) {
		return err
	}

	return tx.Commit(ctx)
}

type WorkerInfo struct {
	worker       *Worker
	lastHearbeat time.Time
}

type WorkerPool struct {
	sync.Mutex
	maxWorkers      uint
	queue           *metadata.QueueSubspace
	txMgr           *transaction.Manager
	tenantMgr       *metadata.TenantManager
	workers         []*WorkerInfo
	nextWorkerId    uint
	workerSleepTime time.Duration
	poolSleepTime   time.Duration
	heartbeatChan   chan uint
	stopChan        chan struct{}
	eventChan       chan Event
	eventListeners  []chan<- Event
}

type Complete func(item *metadata.QueueItem)

func NewWorkerPool(maxWorkers uint, queue *metadata.QueueSubspace, txMgr *transaction.Manager, tenantMgr *metadata.TenantManager, workerSleepTime time.Duration, poolSleepTime time.Duration) *WorkerPool {
	return &WorkerPool{
		maxWorkers:      1,
		queue:           queue,
		txMgr:           txMgr,
		tenantMgr:       tenantMgr,
		workers:         make([]*WorkerInfo, 0),
		nextWorkerId:    0,
		workerSleepTime: workerSleepTime,
		poolSleepTime:   poolSleepTime,
		heartbeatChan:   make(chan uint, maxWorkers*3),
		stopChan:        make(chan struct{}, 1),
		eventChan:       make(chan Event, maxWorkers*3),
		// This is for testing
		eventListeners: make([]chan<- Event, 0),
	}
}

func (pool *WorkerPool) Start() error {
	for i := uint(0); i < pool.maxWorkers; i++ {
		pool.nextWorkerId = i
		pool.workers = append(pool.workers, pool.newWorker(i))
	}

	go pool.Loop()
	return nil
}

func (pool *WorkerPool) newWorker(id uint) *WorkerInfo {
	worker := newWorker(id, pool.queue, pool.txMgr, pool.tenantMgr, pool.workerSleepTime, pool.eventChan, pool.heartbeatChan)
	go worker.Start()
	return &WorkerInfo{
		worker:       worker,
		lastHearbeat: time.Now(),
	}
}

// Internal testing for worker events.
func (pool *WorkerPool) subscribe(listener chan<- Event) {
	pool.Lock()
	defer pool.Unlock()
	pool.eventListeners = append(pool.eventListeners, listener)
}

func (pool *WorkerPool) notify(event Event) {
	pool.Lock()
	defer pool.Unlock()
	for _, listener := range pool.eventListeners {
		select {
		case listener <- event:
			continue
		case <-time.After(5 * time.Millisecond):
			continue
		}
	}
}

func (pool *WorkerPool) Loop() {
	ticker := time.NewTicker(pool.poolSleepTime)
	for {
		select {
		case <-pool.stopChan:
			log.Info().Msg("shutting down worker pool")
			return
		case id := <-pool.heartbeatChan:
			pool.rxHeartbeats(id)
		case event := <-pool.eventChan:
			pool.notify(event)
		case <-ticker.C:
			pool.checkHeartbeats()
		}
	}
}

func (pool *WorkerPool) rxHeartbeats(workerId uint) {
	pool.Lock()
	defer pool.Unlock()

	for _, info := range pool.workers {
		if info.worker.id == workerId {
			info.lastHearbeat = time.Now()
			return
		}
	}

	log.Error().Msgf("received heartbeat from missing worker id %d", workerId)
}

func (pool *WorkerPool) checkHeartbeats() {
	pool.Lock()
	defer pool.Unlock()

	now := time.Now()
	for i, info := range pool.workers {
		if now.Sub(info.lastHearbeat) > 5*pool.workerSleepTime {
			info.worker.Stop()
			pool.nextWorkerId++
			log.Error().Msgf("No response from worker %d adding new worker", info.worker.id)
			pool.workers[i] = pool.newWorker(pool.nextWorkerId)
		}
	}
}

func (pool *WorkerPool) Stop() {
	pool.Lock()
	defer pool.Unlock()
	for _, info := range pool.workers {
		info.worker.Stop()
	}
	pool.stopChan <- struct{}{}
}