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
	"context"
	"flakybit.net/psl/init/config"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const maxIdleConnections = 1

func main() {
	conf := config.NewConfig(context.TODO())

	lockUrl := fmt.Sprintf("http://%s:%d", conf.LockHost, conf.LockPort)
	if conf.LockDuration > 0 {
		values := url.Values{}
		values.Add("duration", strconv.FormatFloat(conf.LockDuration.Seconds(), 'f', 0, 64))
		lockUrl = fmt.Sprintf("%s?%s", lockUrl, values.Encode())
	}
	log.Printf("Will try to acquire lock at '%s' each '%d' sec", lockUrl, conf.Period)

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: maxIdleConnections,
		},
		Timeout: conf.Timeout,
	}
	for {
		if acquireLock(client, lockUrl) {
			return
		}
		time.Sleep(conf.Period)
	}
}

func acquireLock(client *http.Client, url string) bool {
	//values := url.Values{}
	//values.Add("duration", strconv.FormatFloat(conf.LockDuration.Seconds(), 'f', 0, 64))
	//req, err := http.NewRequest("GET", url, strings.NewReader(q.Encode()))
	//if err != nil {
	//	return err
	//}
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//req.Close = true
	//resp, err := http.DefaultClient.Do(req)
	//if err != nil {
	//	return err
	//}

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
