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

	"github.com/stretchr/testify/require"
)

func TestGlobalDial(t *testing.T) {
	conn, err := Dial("mqtt://localhost:1883")
	require.NoError(t, err)
	conn.Close()
}

func TestDialerURLWithoutPort(t *testing.T) {
	conn, err := Dial("mqtt://localhost")
	require.NoError(t, err)
	conn.Close()
}

func TestDialerBadURL(t *testing.T) {
	_, err := Dial("foo")
	require.Equal(t, DialError, toError(err).Code())
}

func TestDialerUnsupportedProtocol(t *testing.T) {
	_, err := Dial("foo://localhost")
	require.Equal(t, DialError, toError(err).Code())
	require.Equal(t, ErrUnsupportedProtocol, toError(err).Err())
}

func TestDialerTCPError(t *testing.T) {
	_, err := Dial("tcp://localhost:1234567")
	require.Equal(t, DialError, toError(err).Code())
}

func TestDialerTLSError(t *testing.T) {
	_, err := Dial("tls://localhost:1234567")
	require.Equal(t, DialError, toError(err).Code())
}

func TestDialerWSError(t *testing.T) {
	_, err := Dial("ws://localhost:1234567")
	require.Equal(t, DialError, toError(err).Code())
}

func TestDialerWSSError(t *testing.T) {
	_, err := Dial("wss://localhost:1234567")
	require.Equal(t, DialError, toError(err).Code())
}

func TestTCPDefaultPort(t *testing.T) {
	tp := newTestPort()

	server := NewServer(make(chan Conn))
	server.LaunchTCP(tp.address())

	dialer := NewDialer()
	dialer.DefaultTCPPort = tp.port()

	conn, err := dialer.Dial("tcp://localhost")
	require.NoError(t, err)

	conn.Close()
	server.Stop()
}

//TODO: something odd here...
//func TestTLSDefaultPort(t *testing.T) {
//	tp := newTestPort()
//
//	server := NewServer(make(chan Conn))
//	server.LaunchTLS(tp.address(), serverTLSConfig)
//
//	dialer := NewDialer()
//	dialer.TLSClientConfig = clientTLSConfig
//	dialer.DefaultTLSPort = tp.port()
//
//	conn, err := dialer.Dial("tls://localhost")
//	require.NoError(t, err)
//
//	conn.Close()
//	server.Stop()
//}

func TestWSDefaultPort(t *testing.T) {
	tp := newTestPort()

	server := NewServer(make(chan Conn))
	server.LaunchWS(tp.address())

	dialer := NewDialer()
	dialer.DefaultWSPort = tp.port()

	conn, err := dialer.Dial("ws://localhost")
	require.NoError(t, err)

	conn.Close()
	server.Stop()
}

func TestWSSDefaultPort(t *testing.T) {
	tp := newTestPort()

	server := NewServer(make(chan Conn))
	server.LaunchWSS(tp.address(), serverTLSConfig)

	dialer := NewDialer()
	dialer.TLSClientConfig = clientTLSConfig
	dialer.DefaultWSSPort = tp.port()

	conn, err := dialer.Dial("wss://localhost")
	require.NoError(t, err)

	conn.Close()
	server.Stop()
}