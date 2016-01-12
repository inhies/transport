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

package transport

import (
	"testing"
)

func tcpPreparer(handler func(Conn)) (Conn, chan Conn) {
	channel := make(chan Conn)
	tp := newTestPort()

	server := NewServer(channel)
	server.LaunchTCP(tp.address())

	go func(){
		handler(<-channel)
		close(channel)
		server.Stop()
	}()

	conn, err := testDialer.Dial(tp.url("tcp"))
	if err != nil {
		panic(err)
	}

	return conn, channel
}

func TestNetConnConnection(t *testing.T) {
	abstractConnConnectTest(t, tcpPreparer)
}

func TestNetConnClose(t *testing.T) {
	abstractConnCloseTest(t, tcpPreparer)
}

func TestNetConnEncodeError(t *testing.T) {
	abstractConnEncodeErrorTest(t, tcpPreparer)
}

func TestNetConnDecode1Error(t *testing.T) {
	abstractConnDecodeError1Test(t, tcpPreparer)
}

func TestNetConnDecode2Error(t *testing.T) {
	abstractConnDecodeError2Test(t, tcpPreparer)
}

func TestNetConnDecode3Error(t *testing.T) {
	abstractConnDecodeError3Test(t, tcpPreparer)
}

func TestNetConnSendAfterClose(t *testing.T) {
	abstractConnSendAfterCloseTest(t, tcpPreparer)
}

func TestNetConnCounters(t *testing.T) {
	abstractConnCountersTest(t, tcpPreparer)
}

func TestNetConnReadLimit(t *testing.T) {
	abstractConnReadLimitTest(t, tcpPreparer)
}