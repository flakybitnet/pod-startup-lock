/*
This file is part of PSL (Pod Startup LockService).
Copyright (c) 2024, The PSL (Pod Startup LockService) Authors

PSL (Pod Startup LockService) is free software:
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
	"fmt"
	"log"
	"regexp"
)

var endpointPattern = regexp.MustCompile(`^(\S+?)://(.*)$`)
var addressPattern = regexp.MustCompile(`^(\S+):(\d+)$`)

type Endpoint interface {
	Protocol() string
	String() string
	IsHttp() bool
}

type RawEndpoint interface {
	Endpoint
	Address() string
}

type HttpEndpoint interface {
	Endpoint
	Url() string
}

type EndpointData struct {
	protocol string
}

type RawEndpointData struct {
	EndpointData
	address string
}

type HttpEndpointData struct {
	EndpointData
	url string
}

func (e *RawEndpointData) String() string {
	return fmt.Sprintf("%s", e.address)
}

func (e *HttpEndpointData) String() string {
	return fmt.Sprintf("%s", e.url)
}

func (e *EndpointData) Protocol() string {
	return e.protocol
}

func (e *EndpointData) IsHttp() bool {
	return isHttp(e.Protocol())
}

func isHttp(protocol string) bool {
	return protocol == "http" || protocol == "https"
}

func (e *RawEndpointData) Address() string {
	return e.address
}

func (e *HttpEndpointData) Url() string {
	return e.url
}

func ParseEndpoint(str string) Endpoint {
	match := endpointPattern.FindStringSubmatch(str)
	if match == nil || len(match) != 3 {
		log.Panicf("Endpoint malformed: '%s'", str)
	}
	protocol := match[1]
	address := match[2]
	if isHttp(protocol) {
		return CreateHttp(protocol, str)
	} else {
		return CreateRaw(protocol, address)
	}
}

func CreateRaw(protocol string, address string) RawEndpoint {
	match := addressPattern.FindStringSubmatch(address)
	if match == nil || len(match) != 3 {
		log.Panicf("Address malformed: '%s'", address)
	}
	return &RawEndpointData{EndpointData{protocol}, address}
}

func CreateHttp(protocol string, url string) HttpEndpoint {
	return &HttpEndpointData{EndpointData{protocol}, url}
}
