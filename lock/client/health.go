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

package client

import (
	"context"
	. "flakybit.net/psl/lock/config"
	"fmt"
	"io"
	log "log/slog"
	"net"
	"net/http"
)

const maxIdleConnections = 10

type HealthClient struct {
	conf       Config
	httpClient *http.Client
	rawClient  *net.Dialer
}

func NewHealthClient(conf Config) *HealthClient {
	httpClient := &http.Client{
		Transport: &http.Transport{MaxIdleConnsPerHost: maxIdleConnections},
		Timeout:   conf.HealthCheck.Timeout,
	}
	dialer := &net.Dialer{
		Timeout: conf.HealthCheck.Timeout,
	}

	client := &HealthClient{
		conf,
		httpClient,
		dialer,
	}
	log.Info("configured Lock httpClient")
	return client
}

func (c *HealthClient) CheckRaw(ctx context.Context, protocol, address string) (bool, error) {
	conn, err := c.rawClient.DialContext(ctx, protocol, address)
	log.Debug("checking raw endpoint", log.String("endpoint", fmt.Sprintf("%s://%s", protocol, address)))
	if err != nil {
		return false, err
	}
	err = conn.Close()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *HealthClient) CheckHttp(ctx context.Context, url string) (bool, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false, err
	}

	log.Debug("checking HTTP endpoint", log.String("url", request.URL.String()))
	response, err := c.httpClient.Do(request)
	if err != nil {
		return false, err
	}
	_, err = io.ReadAll(response.Body)
	err = response.Body.Close()
	if err != nil {
		return false, err
	}

	return response.StatusCode >= 200 && response.StatusCode < 300, nil
}
