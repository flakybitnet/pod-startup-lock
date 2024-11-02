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

package config

import (
	"flag"
	"flakybit.net/psl/common/util"
	"log"
	"time"
)

const defaultPort = 8888
const defaultParallelLocks = 1
const defaultLockTimeout = 10
const defaultFailTimeout = 10
const defaultPassTimeout = 60

func Parse() Config {
	host := flag.String("host", "", "Host/Ip to bind")
	port := flag.Int("port", defaultPort, "Port to bind")
	parallelLocks := flag.Int("locks", defaultParallelLocks, "Count of locks allowed to acquire in parallel")
	lockTimeout := flag.Int("timeout", defaultLockTimeout, "Default lock timeout, sec")
	failTimeout := flag.Int("failHc", defaultFailTimeout, "Pause between endpoint health checks if previous failed, sec")
	passTimeout := flag.Int("passHc", defaultPassTimeout, "Pause between endpoint health checks if previous succeeded, sec")

	var healthEndpoints util.ArrayVal
	flag.Var(&healthEndpoints, "check", "HealthCheck tcp endpoint, host:port")
	flag.Parse()

	config := Config{
		*host,
		*port,
		*parallelLocks,
		time.Duration(*lockTimeout) * time.Second,
		time.Duration(*failTimeout) * time.Second,
		time.Duration(*passTimeout) * time.Second,
		parseEndpoints(healthEndpoints),
	}
	log.Printf("Application config:\n%+v", config)
	return config
}

type Config struct {
	Host              string
	Port              int
	ParallelLocks     int
	LockTimeout       time.Duration
	HealthFailTimeout time.Duration
	HealthPassTimeout time.Duration
	HealthEndpoints   []Endpoint
}

func parseEndpoints(urls []string) []Endpoint {
	var endpoints []Endpoint
	for _, url := range urls {
		endpoints = append(endpoints, ParseEndpoint(url))
	}
	return endpoints
}
