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
*/

package client

import (
	"context"
	. "flakybit.net/psl/init/config"
	"fmt"
	"io"
	log "log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const maxIdleConnections = 1

type LockClient struct {
	conf    Config
	client  *http.Client
	lockUrl string
}

func NewLockClient(conf Config) *LockClient {
	httpClient := &http.Client{
		Transport: &http.Transport{MaxIdleConnsPerHost: maxIdleConnections},
		Timeout:   conf.Timeout,
	}
	client := &LockClient{
		conf,
		httpClient,
		fmt.Sprintf("http://%s:%d", conf.LockHost, conf.LockPort),
	}
	log.Info("configured Lock client", log.String("lock-url", client.lockUrl))
	return client
}

func (c *LockClient) AcquireLock(ctx context.Context) (bool, error) {
	values := url.Values{}
	if c.conf.LockDuration > 0 {
		values.Add("duration", strconv.FormatFloat(c.conf.LockDuration.Seconds(), 'f', 0, 64))
	}
	request, err := http.NewRequestWithContext(ctx, "GET", c.lockUrl, strings.NewReader(values.Encode()))
	if err != nil {
		return false, err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	log.Info("acquiring lock", log.String("url", request.URL.String()))
	response, err := c.client.Do(request)
	if err != nil {
		return false, err
	}
	_, err = io.ReadAll(response.Body)
	err = response.Body.Close()
	if err != nil {
		return false, err
	}

	return response.StatusCode == 200, nil
}
