/*
This file is part of PSL (Pod Startup Lock).
Copyright (c) 2024, The PSL (Pod Startup Lock) Authors

PSL (Pod Startup Lock) is free software:
you can redistribute it and/or modify it under the terms of the GNU General Public License
as published by the Free Software Foundation, version 3 of the License.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY;
without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program.
If not, see <https://www.gnu.org/licenses/>.

This file incorporates work covered by the following copyright and permission notice:
	Copyright (c) 2018, Oath Inc.

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE.
*/

package service

import (
	. "flakybit.net/psl/lock/config"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

const maxIdleConnections = 10
const requestTimeout = 10 * time.Second

var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: maxIdleConnections,
	},
	Timeout: requestTimeout,
}

var dialer = &net.Dialer{
	Timeout: requestTimeout,
}

type EndpointChecker struct {
	waitOnPass time.Duration
	waitOnFail time.Duration
	endpoints  []Endpoint
	isHealthy  bool
}

func NewEndpointChecker(waitOnPass time.Duration, waitOnFail time.Duration, endpoints []Endpoint) EndpointChecker {
	return EndpointChecker{waitOnPass, waitOnFail, endpoints, false}
}

func (c *EndpointChecker) HealthFunction() func() bool {
	return func() bool {
		return c.isHealthy
	}
}

func (c *EndpointChecker) Run() {
	if len(c.endpoints) == 0 {
		log.Print("No Endpoints to check")
		c.isHealthy = true
		return
	}
	for {
		if checkAll(c.endpoints) {
			log.Print("Endpoint Check passed")
			c.isHealthy = true
			time.Sleep(c.waitOnPass)
		} else {
			log.Print("Endpoint Check failed")
			c.isHealthy = false
			time.Sleep(c.waitOnFail)
		}
	}
}

func checkAll(endpoints []Endpoint) bool {
	for _, endpoint := range endpoints {
		if !check(endpoint) {
			return false
		}
	}
	return true
}

func check(endpoint Endpoint) bool {
	if endpoint.IsHttp() {
		return checkHttp(endpoint.(HttpEndpoint))
	} else {
		return checkRaw(endpoint.(RawEndpoint))
	}
}

func checkRaw(endpoint RawEndpoint) bool {
	conn, err := dialer.Dial(endpoint.Protocol(), endpoint.Address())
	if err != nil {
		log.Printf("'%v' endpoint connection failed: '%v'", endpoint, err)
		return false
	}
	conn.Close()
	log.Printf("'%v' endpoint OK", endpoint)
	return true
}

func checkHttp(endpoint HttpEndpoint) bool {
	resp, err := client.Get(endpoint.Url())
	if err != nil {
		log.Printf("'%v' endpoint connection failed: '%v'", endpoint, err)
		return false
	}
	io.Copy(io.Discard, resp.Body)
	defer resp.Body.Close()

	if isSuccessful(resp.StatusCode) {
		log.Printf("'%v' endpoint OK (status: %v)", endpoint, resp.StatusCode)
		return true
	} else {
		log.Printf("'%v' endpoint Fail (status: %v)", endpoint, resp.StatusCode)
		return false
	}
}

func isSuccessful(code int) bool {
	return code >= 200 && code < 300
}
