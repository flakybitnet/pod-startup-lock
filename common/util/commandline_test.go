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
	"github.com/stretchr/testify/require"
	"testing"
)

func TestArrayValStringWhenEmpty(t *testing.T) {
	// GIVEN
	arr := ArrayVal{}
	expected := "[]"

	// WHEN
	actual := arr.String()

	// THEN
	require.Equal(t, expected, actual)
}

func TestArrayValStringWhenNotEmpty(t *testing.T) {
	// GIVEN
	arr := ArrayVal{"val1", "val2"}
	expected := "[val1 val2]"

	// WHEN
	actual := arr.String()

	// THEN
	require.Equal(t, expected, actual)
}

func TestArrayValSet(t *testing.T) {
	// GIVEN
	arr := ArrayVal{}
	expected := ArrayVal{"val1", "val2"}

	// WHEN
	arr.Set("val1")
	arr.Set("val2")

	// THEN
	require.Equal(t, expected, arr)
}

func TestNewPairArrayVal(t *testing.T) {
	// GIVEN
	// WHEN
	arrayVal := NewPairArrayVal("-")

	// THEN
	require.Equal(t, "-", arrayVal.sep, "Wrong separator")
	require.Len(t, arrayVal.Get(), 0)
}

func TestPairArrayValString(t *testing.T) {
	// GIVEN
	arrayVal := NewPairArrayVal(":")
	arrayVal.Set("a:1")
	arrayVal.Set("b:2")
	expected := "[{a 1} {b 2}]"

	// WHEN
	actual := arrayVal.String()

	// THEN
	require.Equal(t, expected, actual)
}

func TestPairArrayValSetWhenEmpty(t *testing.T) {
	// GIVEN
	arrayVal := NewPairArrayVal(":")

	// WHEN
	panicFunc := func() { arrayVal.Set("") }

	// THEN
	require.PanicsWithValue(t, "Failed to parse value: ''", panicFunc)
}

func TestPairArrayValSetWhenNoValues(t *testing.T) {
	// GIVEN
	arrayVal := NewPairArrayVal(":")

	// WHEN
	panicFunc := func() { arrayVal.Set(":") }

	// THEN
	require.PanicsWithValue(t, "Failed to parse value: ':'", panicFunc)
}

func TestPairArrayValSetWhenNoKey(t *testing.T) {
	// GIVEN
	arrayVal := NewPairArrayVal(":")

	// WHEN
	panicFunc := func() { arrayVal.Set(":val") }

	// THEN
	require.PanicsWithValue(t, "Failed to parse value: ':val'", panicFunc)
}

func TestPairArrayValSetWhenNoValue(t *testing.T) {
	// GIVEN
	arrayVal := NewPairArrayVal(":")

	// WHEN
	panicFunc := func() { arrayVal.Set("key:") }

	// THEN
	require.PanicsWithValue(t, "Failed to parse value: 'key:'", panicFunc)
}

func TestPairArrayValSetSingle(t *testing.T) {
	// GIVEN
	arrayVal := NewPairArrayVal(":")
	expected := []Pair{{"a", "1"}}

	// WHEN
	arrayVal.Set("a:1")

	// THEN
	require.Equal(t, expected, arrayVal.Get())
}

func TestPairArrayValSetMultiple(t *testing.T) {
	// GIVEN
	arrayVal := NewPairArrayVal(":")
	expected := []Pair{{"a", "1"}, {"b", "2"}}

	// WHEN
	arrayVal.Set("a:1")
	arrayVal.Set("b:2")

	// THEN
	require.Equal(t, expected, arrayVal.Get())
}

func TestPairArrayValSetMultipleInvalid(t *testing.T) {
	// GIVEN
	arrayVal := NewPairArrayVal(":")

	// WHEN
	panicFunc := func() { arrayVal.Set("key::val") }

	// THEN
	require.PanicsWithValue(t, "Failed to parse value: 'key::val'", panicFunc)
}
