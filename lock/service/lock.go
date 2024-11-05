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
	. "flakybit.net/psl/lock/config"
	log "log/slog"
	"sync"
	"time"
)

type LockService struct {
	conf  Config
	mutex sync.Mutex
	locks []time.Time
}

func NewLockService(conf Config) *LockService {
	service := &LockService{conf: conf}
	log.Info("configured lock service")
	return service
}

func (ls *LockService) Acquire(duration time.Duration) bool {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()

	ls.removeExpired()
	if len(ls.locks) < ls.conf.ParallelLocks {
		ls.addNew(duration)
		log.Info("lock acquired",
			log.Int("duration", int(duration.Seconds())),
			log.Int("locks", len(ls.locks)))
		return true
	}
	return false
}

func (ls *LockService) addNew(duration time.Duration) {
	expireTime := time.Now().Add(duration)
	ls.locks = append(ls.locks, expireTime)
}

func (ls *LockService) removeExpired() {
	var live []time.Time
	for i := 0; i < len(ls.locks); i++ {
		if !isExpired(ls.locks[i]) {
			live = append(live, ls.locks[i])
		}
	}
	ls.locks = live
}

func isExpired(t time.Time) bool {
	return time.Now().After(t)
}
