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
	"crypto/tls"
	"net"
	"net/http"
	"errors"
	"time"

	"github.com/gorilla/websocket"
	"gopkg.in/tomb.v2"
)

var ErrStopped = errors.New("server: already stopped")

// The Server manages multiple listeners and sends new streams to the registered
// channel. The serve requires an already created channel to be supplied, as in
// most cases that channel would never close so that the backend can restart
// independently from the broker logic.
type Server struct {
	in        chan<- Conn
	listeners []net.Listener
	upgrader  *websocket.Upgrader

	tomb      tomb.Tomb
}

// NewServer returns a new Server.
func NewServer(channel chan<- Conn) *Server {
	s := &Server{
		in: channel,
		listeners: make([]net.Listener, 0),
		upgrader: &websocket.Upgrader{
			HandshakeTimeout: 60 * time.Second,
			Subprotocols: []string{ "mqtt" },
		},
	}

	// start cleanup function
	go func(){
		select {
		case <-s.tomb.Dying():
			for _, l := range s.listeners {
				l.Close()
			}
		}
	}()

	return s
}

// Error returns the last occurred error. The return value can be consulted
// when the server gets stopped unexpectedly because of an potential error.
func (s *Server) Error() error {
	err := s.tomb.Err()
	if err == tomb.ErrStillAlive {
		return nil
	}

	return err
}

// LaunchTCP will launch a TCP server.
func (s *Server) LaunchTCP(address string) error {
	select {
	case <- s.tomb.Dying():
		return ErrStopped
	default:
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	s.listeners = append(s.listeners, listener)
	s.AcceptConnections(listener)

	return nil
}

// LaunchTLS will launch a TLS server.
func (s *Server) LaunchTLS(address string, config *tls.Config) error {
	select {
	case <- s.tomb.Dying():
		return ErrStopped
	default:
	}

	listener, err := tls.Listen("tcp", address, config)
	if err != nil {
		return err
	}

	s.listeners = append(s.listeners, listener)
	s.AcceptConnections(listener)

	return nil
}

// AcceptConnections accepts and sends new connections as streams to the
// Accept channel.
//
// Note: If the server has been stopped due to Stop() or an error the
// internally started goroutine will return.
func (s *Server) AcceptConnections(listener net.Listener) {
	s.tomb.Go(func() error {
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-s.tomb.Dying():
					return nil
				default:
					return err
				}
			}

			stream := NewNetConn(conn)

			select {
			case s.in <- stream:
			case <-s.tomb.Dying():
				stream.Close()

				return tomb.ErrDying
			}
		}
	})
}

// LaunchWS will launch a WS server.
func (s *Server) LaunchWS(address string) error {
	select {
	case <- s.tomb.Dying():
		return ErrStopped
	default:
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	s.listeners = append(s.listeners, listener)
	s.launchHTTP(listener)

	return nil
}

// LaunchWSS will launch a WSS server.
func (s *Server) LaunchWSS(address string, config *tls.Config) error {
	select {
	case <- s.tomb.Dying():
		return ErrStopped
	default:
	}

	listener, err := tls.Listen("tcp", address, config)
	if err != nil {
		return err
	}

	s.listeners = append(s.listeners, listener)
	s.launchHTTP(listener)

	return nil
}

// a helper to create a http mux and serve it from a listener
func (s *Server) launchHTTP(listener net.Listener) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.RequestHandler())

	h := &http.Server{
		Handler: mux,
	}

	s.tomb.Go(func() error {
		err := h.Serve(listener)

		select {
		case <-s.tomb.Dying():
			return tomb.ErrDying
		default:
			return err
		}
	})
}

// RequestHandler returns the WebsSocket handler function.
// You can directly mount this function in your existing HTTP multiplexer.
//
// Note: If the server has been stopped due to Stop() or an error the
// returned handler will return a HTTP 500 error.
func (s *Server) RequestHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		case <- s.tomb.Dying():
			http.Error(w, "Internal Server Error", 500)
			return
		default:
		}

		conn, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			// upgrader already responded to request
			return
		}

		stream := NewWebSocketConn(conn)

		select {
		case s.in <- stream:
		case <-s.tomb.Dying():
			stream.Close()
		}
	}
}

// Stop will stop the server by closing all listeners.
func (s *Server) Stop() {
	s.tomb.Kill(nil)
	s.tomb.Wait()
}

// Stopped will return a boolean indicating if the server has been
// already stopped by Stop() or an error.
func (s *Server) Stopped() bool {
	return !s.tomb.Alive()
}