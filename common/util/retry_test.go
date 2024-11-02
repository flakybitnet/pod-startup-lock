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

package util

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRetrySuccess(t *testing.T) {
	// GIVEN
	expected := "success"
	successFunc := func() (interface{}, error) {
		return expected, nil
	}

	// WHEN
	actual := (*RetryOrPanic(1, 1, successFunc)).(string)

	// THEN
	require.Equal(t, actual, expected)
}

func TestRetryDefaultSuccess(t *testing.T) {
	// GIVEN
	expected := "success"
	successFunc := func() (interface{}, error) {
		return expected, nil
	}

	// WHEN
	actual := (*RetryOrPanicDefault(successFunc)).(string)

	// THEN
	require.Equal(t, actual, expected)
}

func TestRetryFail(t *testing.T) {
	// GIVEN
	errorFunc := func() (interface{}, error) {
		return nil, fmt.Errorf("error")
	}
	expected := "Failed after 1 attempts, last error: error"

	// WHEN
	panicFunc := func() { RetryOrPanic(1, 1, errorFunc) }

	// THEN
	require.PanicsWithValue(t, expected, panicFunc)
}

func TestMultipleRetryFail(t *testing.T) {
	// GIVEN
	errorFunc := func() (interface{}, error) {
		return nil, fmt.Errorf("error")
	}
	expected := "Failed after 2 attempts, last error: error"

	// WHEN
	panicFunc := func() { RetryOrPanic(2, 0, errorFunc) }

	// THEN
	require.PanicsWithValue(t, expected, panicFunc)
}
