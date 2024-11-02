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
	. "flakybit.net/psl/common/util"
	"log"
	"os"
	"time"
)

const defaultPort = 9999
const defaultFailTimeout = 10
const defaultPassTimeout = 60

func Parse() Config {
	host := flag.String("host", "", "Host/Ip to bind")
	port := flag.Int("port", defaultPort, "Port to bind")
	baseUrl := flag.String("baseUrl", "", "K8s api base url. For out-of-cluster usage only")
	namespace := flag.String("namespace", "", "K8s Namespace to check DaemonSets in. Blank for all namespaces")
	failTimeout := flag.Int("failHc", defaultFailTimeout, "Pause between DaemonSet health checks if previous failed, sec")
	passTimeout := flag.Int("passHc", defaultPassTimeout, "Pause between DaemonSet health checks if previous succeeded, sec")
	hostNetwork := flag.Bool("hostNet", false, "Host network DaemonSets only")

	nodeName, _ := os.LookupEnv("NODE_NAME")

	includeDs := NewPairArrayVal(":")
	flag.Var(&includeDs, "in", "Include DaemonSet labels, label:value")
	excludeDs := NewPairArrayVal(":")
	flag.Var(&excludeDs, "ex", "Exclude DaemonSet labels, label:value")
	flag.Parse()

	config := Config{
		*host,
		*port,
		*baseUrl,
		*namespace,
		time.Duration(*failTimeout) * time.Second,
		time.Duration(*passTimeout) * time.Second,
		nodeName,
		*hostNetwork,
		includeDs.Get(),
		excludeDs.Get(),
	}
	log.Printf("Application config:\n%+v", config)
	config.Validate()
	return config
}

type Config struct {
	Host              string
	Port              int
	K8sApiBaseUrl     string
	Namespace         string
	HealthFailTimeout time.Duration
	HealthPassTimeout time.Duration
	NodeName          string
	HostNetworkDs     bool
	IncludeDs         []Pair
	ExcludeDs         []Pair
}

func (c *Config) Validate() {
	if c.NodeName == "" {
		log.Panic("NODE_NAME not specified")
	}
	if len(c.IncludeDs) > 0 && len(c.ExcludeDs) > 0 {
		log.Panic("Cannot specify both Included and Excluded DaemonSet labels, choose one")
	}
}
