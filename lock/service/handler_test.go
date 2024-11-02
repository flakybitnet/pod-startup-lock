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
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var timeout = time.Duration(10) * time.Second

func TestAcquireIfFirst(t *testing.T) {
	// GIVEN
	permitFunction := func() bool {
		return true
	}
	lock := state.NewLock(1)
	handler := NewLockHandler(&lock, timeout, permitFunction)
	req, _ := http.NewRequest("GET", "/", nil)

	// WHEN
	rr := prepareResponseRecorder(req, handler)

	// THEN
	assertResponseStatusCode(http.StatusOK, rr.Code, t)
}

func TestAcquireIfSecond(t *testing.T) {
	// GIVEN
	permitFunction := func() bool {
		return true
	}

	lock := state.NewLock(1)
	handler := NewLockHandler(&lock, timeout, permitFunction)
	req, _ := http.NewRequest("GET", "/", nil)
	prepareResponseRecorder(req, handler)

	// WHEN
	rr := prepareResponseRecorder(req, handler)

	// THEN
	assertResponseStatusCode(http.StatusLocked, rr.Code, t)
}

func TestAcquireIfWrongTimeoutRequested(t *testing.T) {
	// GIVEN
	permitFunction := func() bool {
		return true
	}

	lock := state.NewLock(1)
	handler := NewLockHandler(&lock, timeout, permitFunction)
	req, _ := http.NewRequest("GET", "/", nil)
	q := req.URL.Query()
	q.Add("duration", "a")
	req.URL.RawQuery = q.Encode()

	prepareResponseRecorder(req, handler)

	// WHEN
	rr := prepareResponseRecorder(req, handler)

	// THEN
	assertResponseStatusCode(http.StatusLocked, rr.Code, t)
}

func TestAcquireIfZeroTimeoutRequested(t *testing.T) {
	// GIVEN
	permitFunction := func() bool {
		return true
	}

	lock := state.NewLock(1)
	handler := NewLockHandler(&lock, timeout, permitFunction)
	req, _ := http.NewRequest("GET", "/", nil)
	q := req.URL.Query()
	q.Add("duration", "0")
	req.URL.RawQuery = q.Encode()

	prepareResponseRecorder(req, handler)

	// WHEN
	rr := prepareResponseRecorder(req, handler)

	// THEN
	assertResponseStatusCode(http.StatusOK, rr.Code, t)
}

func TestAcquireIfDisabled(t *testing.T) {
	// GIVEN
	permitFunction := func() bool {
		return false
	}

	lock := state.NewLock(1)
	handler := NewLockHandler(&lock, timeout, permitFunction)
	req, _ := http.NewRequest("GET", "/", nil)

	// WHEN
	rr := prepareResponseRecorder(req, handler)

	// THEN
	assertResponseStatusCode(http.StatusLocked, rr.Code, t)
}

func assertResponseStatusCode(expected int, actual int, t *testing.T) {
	if actual != expected {
		t.Errorf("handler returned wrong status code: expected %v got %v", expected, actual)
	}
}

func prepareResponseRecorder(req *http.Request, handler http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}
