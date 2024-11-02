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
	Copyright (c) 2022, The PSL (Pod Startup Lock) Authors

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
	"flakybit.net/psl/lock/state"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type lockHandler struct {
	lock            *state.Lock
	defaultTimeout  time.Duration
	permitAcquiring func() bool
}

func NewLockHandler(lock *state.Lock, defaultTimeout time.Duration, permitOperationChecker func() bool) http.Handler {
	return &lockHandler{lock, defaultTimeout, permitOperationChecker}
}

func (h *lockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.permitAcquiring() {
		respondLocked(w, r)
		return
	}
	duration, ok := getRequestedDuration(r.URL.Query())
	if !ok {
		duration = h.defaultTimeout
	}

	if h.lock.Acquire(duration) {
		respondOk(w, r)
	} else {
		respondLocked(w, r)
	}
}

func getRequestedDuration(values url.Values) (time.Duration, bool) {
	durationStr := values.Get("duration")
	if durationStr == "" {
		return 0, false
	}
	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		log.Printf("Invalid duration requested: '%v'", durationStr)
		return 0, false
	}
	return time.Duration(duration) * time.Second, true
}

func respondOk(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	log.Printf("Responding to '%v': %v", r.RemoteAddr, status)
	w.WriteHeader(status)
	w.Write([]byte("Lock acquired"))
}

func respondLocked(w http.ResponseWriter, r *http.Request) {
	status := http.StatusLocked
	log.Printf("Responding to '%v': %v", r.RemoteAddr, status)
	w.WriteHeader(status)
	w.Write([]byte("Locked"))
}
