// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fscontext_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/fleetspeak/fleetspeak/src/common/fscontext"
)

func TestFromDoneChanTODO(t *testing.T) {
	done := make(chan struct{})

	ctx, cancel := fscontext.FromDoneChanTODO(done)
	defer cancel()

	if err := ctx.Err(); err != nil {
		t.Errorf("done channel still open: ctx.Err() = %v, want nil", err)
	}

	close(done)

	select {
	case <-time.After(time.Second):
		t.Errorf("timeout waiting for context cancelation")
	case <-ctx.Done():
		if err := ctx.Err(); err != context.Canceled {
			t.Errorf("done channel closed: ctx.Err() = %v, want canceled", err)
		}
	}
}