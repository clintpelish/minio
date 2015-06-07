/*
 * Minimalist Object Storage, (C) 2015 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package quota

import "net/http"

// requestLimitHandler
type connLimit struct {
	handler         http.Handler
	connectionQueue chan bool
}

func (c *connLimit) Add() {
	c.connectionQueue <- true
	return
}

func (c *connLimit) Remove() {
	<-c.connectionQueue
	return
}

// ServeHTTP is an http.Handler ServeHTTP method
func (c *connLimit) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c.Add()
	c.handler.ServeHTTP(w, req)
	c.Remove()
}

// ConnectionLimit limits the number of concurrent connections
func ConnectionLimit(h http.Handler, limit int) http.Handler {
	return &connLimit{
		handler:         h,
		connectionQueue: make(chan bool, limit),
	}
}
