// Copyright (c) 2014 The gomqtt Authors. All rights reserved.
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

// Package transport implements functionality for handling MQTT 3.1.1
// (http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/) connections.
package transport

import (
	"github.com/gomqtt/packet"
)

// Conn defines the interface for all transport connections.
type Conn interface {
	// Send will write the packet to the underlying connection. It will return
	// an error if there was an error while encoding or writing to the
	// underlying connection.
	Send(packet.Packet) error

	// Receive will read from the underlying connection and return a fully read
	// packet. It will return an error if there was an error while decoding or
	// reading from the underlying connection.
	Receive() (packet.Packet, error)

	// Close will close the underlying connection and cleanup resources.
	Close() error

	// BytesWritten will return the number of bytes successfully written to
	// the underlying connection.
	BytesWritten() int64

	// BytesRead will return the number of bytes successfully read from the
	// underlying connection.
	BytesRead() int64

	// SetReadLimit sets the maximum size for a packet read. If a packet exceeds
	// the limit, the connection gets closed and returns an error.
	SetReadLimit(limit int64)
}