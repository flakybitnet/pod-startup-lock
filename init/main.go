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

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const defaultHost = "localhost"
const defaultPort = 8888
const defaultPause = 1
const defaultTimeout = 0

const maxIdleConnections = 1
const requestTimeout = 1 * time.Second

func main() {
	host := flag.String("host", defaultHost, "Lock service host")
	port := flag.Int("port", defaultPort, "Lock service port")
	duration := flag.Int("duration", defaultTimeout, "Custom lock duration to request, sec")
	pauseSec := flag.Int("pause", defaultPause, "Pause between lock attempts, sec")
	flag.Parse()

	pause := time.Duration(*pauseSec) * time.Second
	url := fmt.Sprintf("http://%s:%v", *host, *port)
	if *duration > 0 {
		url = fmt.Sprintf("%s?duration=%v", url, *duration)
	}
	log.Printf("Will try to acquire lock at '%s' each '%v' sec", url, *pauseSec)

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: maxIdleConnections,
		},
		Timeout: requestTimeout,
	}
	for {
		if acquireLock(client, url) {
			return
		}
		time.Sleep(pause)
	}
}

func acquireLock(client *http.Client, url string) bool {
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Error occurred: '%v'", err)
		return false
	}
	io.Copy(io.Discard, resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Lock not acquired, waiting (status: %v)", resp.StatusCode)
		return false
	}
	log.Print("Lock acquired, exiting")
	return true
}
