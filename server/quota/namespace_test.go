// Copyright 2022 Tigris Data, Inc.
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

package quota

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tigrisdata/tigris/server/config"
	"github.com/tigrisdata/tigris/server/metadata"
	"go.uber.org/atomic"
)

type testMetricsBackend struct {
	ReadRate  atomic.Int64
	WriteRate atomic.Int64
}

func (d *testMetricsBackend) CurRates(_ context.Context, _ string) (int64, int64, error) {
	log.Debug().Int64("read_units", d.ReadRate.Load()).Int64("write_units", d.WriteRate.Load()).Msg("Update read rates")
	if d.ReadRate.Load() < 0 {
		return 0, 0, fmt.Errorf("simulate error")
	}

	return d.ReadRate.Load(), d.WriteRate.Load(), nil
}

func TestNamespaceQuota(t *testing.T) {
	tenants, ctx, cancel := metadata.NewTestTenantMgr(kvStore)
	defer cancel()

	ns := fmt.Sprintf("ns-test-tenantQuota-1-%x", rand.Uint64()) //nolint:golint,gosec
	id := rand.Uint32()                                          //nolint:golint,gosec

	_, err := tenants.CreateOrGetTenant(ctx, metadata.NewTenantNamespace(ns, id))
	require.NoError(t, err)

	id = rand.Uint32() //nolint:golint,gosec
	_, err = tenants.CreateOrGetTenant(ctx, metadata.NewTenantNamespace(ns+"_other", id))
	require.NoError(t, err)

	tb := &testMetricsBackend{}

	m := initNamespace(tenants, &config.QuotaConfig{
		Namespace: config.NamespaceLimitsConfig{
			Default: config.LimitsConfig{
				ReadUnits:  100,
				WriteUnits: 50,
			},
			Node: config.LimitsConfig{
				ReadUnits:  10,
				WriteUnits: 5,
			},
			RefreshInterval: 1 * time.Millisecond,
			Namespaces: map[string]config.LimitsConfig{
				ns + "_other": {
					ReadUnits:  100,
					WriteUnits: 50,
				},
			},
		},
	}, tb)

	time.Sleep(3 * time.Millisecond)

	is := m.getState(ns)
	assert.Equal(t, int64(10), is.setReadLimit.Load())
	assert.Equal(t, int64(5), is.setWriteLimit.Load())

	require.NoError(t, m.Allow(ctx, ns, 6000, false)) // 2 units: 4096 + something
	require.NoError(t, m.Allow(ctx, ns, 1500, true))  // 2 units: 1024 + something

	i := 0
	for ; err != ErrReadUnitsExceeded; i++ {
		err = m.Allow(ctx, ns, 2048, false)
	}
	assert.Equal(t, 9, i)

	i = 0
	for ; err != ErrWriteUnitsExceeded; i++ {
		err = m.Allow(ctx, ns, 512, true) // < 1024 = 1 unit
	}
	assert.Equal(t, 4, i)

	log.Debug().Msg("simulate rate surge")

	tb.ReadRate.Store(11000)
	tb.WriteRate.Store(50000)

	time.Sleep(30 * time.Millisecond)

	// regulator rates should be minimized
	assert.Equal(t, int64(1), is.setReadLimit.Load())
	assert.Equal(t, int64(1), is.setWriteLimit.Load())

	log.Debug().Msg("simulate rate drop to minimal")

	tb.ReadRate.Store(1)
	tb.WriteRate.Store(1)

	time.Sleep(30 * time.Millisecond)

	// set rates should be maximized
	assert.Equal(t, int64(10), is.setReadLimit.Load())
	assert.Equal(t, int64(5), is.setWriteLimit.Load())

	m.Cleanup()
}

func TestNamespaceQuotaCalcLimits(t *testing.T) {
	cases := []struct {
		name string

		maxNodeLimit int64
		maxNamespace int64
		curNamespace int64
		setNodeLimit int64
		hysteresis   int64
		increment    int64
		exp          int64
	}{
		{
			name:         "traffic_below_limit",
			setNodeLimit: 1,
			maxNodeLimit: 100,
			curNamespace: 100,
			maxNamespace: 1000,
			hysteresis:   10,
			increment:    20,
			exp:          21,
		},
		{
			name:         "traffic_above_limit",
			setNodeLimit: 35,
			maxNodeLimit: 100,
			curNamespace: 1100,
			maxNamespace: 1000,
			hysteresis:   10,
			increment:    20,
			exp:          15,
		},
		{
			name:         "minimum_bound",
			setNodeLimit: 1,
			maxNodeLimit: 100,
			curNamespace: 1100,
			maxNamespace: 1000,
			hysteresis:   10,
			increment:    20,
			exp:          1,
		},
		{
			name:         "maximum_bound",
			setNodeLimit: 90,
			maxNodeLimit: 100,
			curNamespace: 100,
			maxNamespace: 1000,
			hysteresis:   10,
			increment:    20,
			exp:          100,
		},
		{
			name:         "first_time_set",
			setNodeLimit: 0,
			maxNodeLimit: 100,
			curNamespace: 1100,
			maxNamespace: 1000,
			hysteresis:   10,
			increment:    20,
			exp:          80,
		},
		{
			name:         "in_high_threshold",
			setNodeLimit: 17,
			maxNodeLimit: 100,
			curNamespace: 1005,
			maxNamespace: 1000,
			hysteresis:   10,
			increment:    20,
			exp:          17,
		},
		{
			name:         "in_low_threshold",
			setNodeLimit: 18,
			maxNodeLimit: 100,
			curNamespace: 997,
			maxNamespace: 1000,
			hysteresis:   10,
			increment:    20,
			exp:          18,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := calcLimit(
				c.setNodeLimit,
				c.maxNodeLimit,
				c.curNamespace,
				c.maxNamespace,
				c.hysteresis,
				c.increment,
			)
			assert.Equal(t, c.exp, res)
		})
	}
}