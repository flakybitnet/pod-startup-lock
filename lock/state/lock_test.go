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

package state

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var duration = time.Duration(10) * time.Second

func TestAcquireSingleIfFirst(t *testing.T) {
	// GIVEN
	lock := NewLock(1)

	// WHEN
	success := lock.Acquire(duration)

	// THEN
	require.True(t, success)
}

func TestAcquireSingleIfSecond(t *testing.T) {
	// GIVEN
	lock := NewLock(1)
	lock.Acquire(duration)

	// WHEN
	success := lock.Acquire(duration)

	// THEN
	require.False(t, success)
}

func TestAcquireSingleIfReleased(t *testing.T) {
	// GIVEN
	lock := NewLock(1)
	lock.Acquire(0)

	// WHEN
	success := lock.Acquire(duration)

	// THEN
	require.True(t, success)
}

func TestAcquireMultipleIfFirst(t *testing.T) {
	// GIVEN
	lock := NewLock(2)

	// WHEN
	success := lock.Acquire(duration)

	// THEN
	require.True(t, success)
}

func TestAcquireMultipleIfSecond(t *testing.T) {
	// GIVEN
	lock := NewLock(2)
	lock.Acquire(duration)

	// WHEN
	success := lock.Acquire(duration)

	// THEN
	require.True(t, success)
}

func TestAcquireMultipleIfExceed(t *testing.T) {
	// GIVEN
	lock := NewLock(2)
	lock.Acquire(duration)
	lock.Acquire(duration)

	// WHEN
	success := lock.Acquire(duration)

	// THEN
	require.False(t, success)
}

func TestAcquireMultipleIfReleased(t *testing.T) {
	// GIVEN
	lock := NewLock(2)
	lock.Acquire(0)
	lock.Acquire(duration)

	// WHEN
	success := lock.Acquire(duration)

	// THEN
	require.True(t, success)
}
