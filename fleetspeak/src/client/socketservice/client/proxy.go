// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package client is a go client library for socketservice.Service.
package client

import (
	"os"
	"time"

	log "github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/fleetspeak/fleetspeak/src/client/channel"

	fcpb "github.com/google/fleetspeak/fleetspeak/src/client/channel/proto/fleetspeak_channel"
	fspb "github.com/google/fleetspeak/fleetspeak/src/common/proto/fleetspeak"
)

// Back-off constants for retrying channel connections.
// Current values will result in the following sequence of
// delays (seconds, rounded off): [1, 1.5, 2.3, 3.4, 5.0, 7.6, 11.4, 15, 15, 15, ...]
const (
	maxChannelRetryDelay      = 15 * time.Second
	channelRetryBackoffFactor = 1.5
)

func backOffChannelRetryDelay(currentDelay time.Duration) time.Duration {
	newDelay := time.Duration(float64(currentDelay) * channelRetryBackoffFactor)
	if newDelay > maxChannelRetryDelay {
		return maxChannelRetryDelay
	} else {
		return newDelay
	}
}

// OpenChannel creates a channel.RelentlessChannel to a fleetspeak client
// through an agreed upon unix domain socket.
func OpenChannel(socketPath string, version string) *channel.RelentlessChannel {
	return channel.NewRelentlessChannel(
		func() (*channel.Channel, func()) {
			sd, err := ptypes.MarshalAny(&fcpb.StartupData{Pid: int64(os.Getpid()), Version: version})
			if err != nil {
				log.Fatalf("unable to marshal StartupData: %v", err)
			}
			m := &fspb.Message{
				MessageType: "StartupData",
				Destination: &fspb.Address{ServiceName: "system"},
				Data:        sd,
			}
		L:
			for {
				ch, fin := buildChannel(socketPath)
				if ch == nil {
					return ch, fin
				}
				select {
				case e := <-ch.Err:
					log.Errorf("Channel failed with error: %v", e)
					fin()
					continue L
				case ch.Out <- m:
					return ch, fin
				}
			}
		})
}
