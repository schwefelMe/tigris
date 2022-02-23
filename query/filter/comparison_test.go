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

package filter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestNewMatcher(t *testing.T) {
	matcher, err := NewMatcher(EQ, structpb.NewNumberValue(1))
	require.NoError(t, err)

	_, ok := matcher.(*EqualityMatcher)
	require.True(t, ok)

	matcher, err = NewMatcher("foo", structpb.NewNumberValue(1))
	require.Equal(t, status.Errorf(codes.InvalidArgument, "unsupported operand 'foo'"), err)
	require.Nil(t, matcher)
}
