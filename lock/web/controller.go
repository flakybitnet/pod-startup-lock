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

package web

import (
	. "flakybit.net/psl/common"
	. "flakybit.net/psl/lock/config"
	. "flakybit.net/psl/lock/service"
	"fmt"
	log "log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Controller struct {
	conf          Config
	healthChecker HealthChecker
	lockService   *LockService
}

func NewController(conf Config, healthChecker HealthChecker, lockService *LockService) *Controller {
	controller := &Controller{conf, healthChecker, lockService}
	log.Info("configured web controller")
	return controller
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	message := "Lock acquired"

	if c.healthChecker.IsHealthy() {
		duration := c.getRequestedDuration(r.URL.Query())
		if duration == 0 {
			duration = c.conf.LockDuration
		}
		acquired := c.lockService.Acquire(duration)
		if !acquired {
			status = http.StatusLocked
			message = "Locked"
		}
	} else {
		status = http.StatusLocked
		message = "Locked"
	}

	log.Info("responding to health request",
		log.String("client-ip", r.RemoteAddr),
		log.Int("status", status))

	w.WriteHeader(status)
	_, err := fmt.Fprint(w, message)
	if err != nil {
		log.Error("failed to respond to health check request",
			log.String("client-ip", r.RemoteAddr),
			log.Int("status", status),
			log.Any("error", err))
	}
}

func (c *Controller) getRequestedDuration(values url.Values) time.Duration {
	durationStr := values.Get("duration")
	if durationStr == "" {
		return 0
	}
	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		log.Error("invalid requested duration",
			log.String("duration", durationStr),
			log.Any("error", err))
		return 0
	}
	return time.Duration(duration) * time.Second
}
