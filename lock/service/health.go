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
	"context"
	. "flakybit.net/psl/lock/client"
	. "flakybit.net/psl/lock/config"
	log "log/slog"
	"time"
)

type HealthCheckService struct {
	conf      Config
	client    *HealthClient
	endpoints []Endpoint
	healthy   bool
}

func NewHealthCheckService(conf Config, client *HealthClient) *HealthCheckService {
	var endpoints []Endpoint
	for _, url := range conf.HealthCheck.Endpoints {
		endpoints = append(endpoints, ParseEndpoint(url))
	}
	checker := &HealthCheckService{
		conf,
		client,
		endpoints,
		false,
	}
	log.Info("configured health check service")
	return checker
}

func (hcs *HealthCheckService) IsHealthy() bool {
	if !hcs.conf.HealthCheck.Enabled {
		return true
	}
	return hcs.healthy
}

func (hcs *HealthCheckService) Run(ctx context.Context) {
	ticker := time.NewTicker(hcs.conf.HealthCheck.PeriodOnFail)
	defer ticker.Stop()

	for {
		checkStatus := hcs.checkAll(ctx, hcs.endpoints)
		if checkStatus != hcs.healthy {
			if checkStatus {
				ticker.Reset(hcs.conf.HealthCheck.PeriodOnPass)
			} else {
				ticker.Reset(hcs.conf.HealthCheck.PeriodOnFail)
			}
		}
		log.Debug("performed health checks", log.Bool("healthy", checkStatus))
		hcs.healthy = checkStatus

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return
		}
	}
}

func (hcs *HealthCheckService) checkAll(ctx context.Context, endpoints []Endpoint) bool {
	for _, endpoint := range endpoints {
		if !hcs.check(ctx, endpoint) {
			return false
		}
	}
	return true
}

func (hcs *HealthCheckService) check(ctx context.Context, endpoint Endpoint) bool {
	if endpoint.IsHttp() {
		healthy, err := hcs.client.CheckHttp(ctx, endpoint.(HttpEndpoint).Url())
		if err != nil {
			log.ErrorContext(ctx, "failed to check endpoint", log.Any("error", err))
		}
		return healthy
	} else {
		healthy, err := hcs.client.CheckRaw(ctx, endpoint.Protocol(), endpoint.(RawEndpoint).Address())
		if err != nil {
			log.ErrorContext(ctx, "failed to check endpoint", log.Any("error", err))
		}
		return healthy
	}
}
