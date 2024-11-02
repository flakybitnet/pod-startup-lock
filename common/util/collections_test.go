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

func TestMapContainsAllPairsWhenAllEmpty(t *testing.T) {
	// GIVEN
	haystack := make(map[string]string)
	needle := make(map[string]string)

	// WHEN
	contains := MapContainsAll(haystack, needle)

	// THEN
	require.False(t, contains)
}

func TestMapContainsAllPairsWhenPairsEmpty(t *testing.T) {
	// GIVEN
	haystack := map[string]string{"a": "1", "b": "2"}
	needle := make(map[string]string)

	// WHEN
	contains := MapContainsAll(haystack, needle)

	// THEN
	require.False(t, contains)
}

func TestMapContainsAllPairsWhenNotContains(t *testing.T) {
	// GIVEN
	haystack := map[string]string{"a": "1", "b": "2"}
	needle := map[string]string{"c": "3"}

	// WHEN
	contains := MapContainsAll(haystack, needle)

	// THEN
	require.False(t, contains)
}

func TestMapContainsAllPairsWhenNotContainsValue(t *testing.T) {
	// GIVEN
	haystack := map[string]string{"a": "1", "b": "2"}
	needle := map[string]string{"a": "3"}

	// WHEN
	contains := MapContainsAll(haystack, needle)

	// THEN
	require.False(t, contains)
}

func TestMapContainsAllPairsWhenContainsSingle(t *testing.T) {
	// GIVEN
	haystack := map[string]string{"a": "1", "b": "2"}
	needle := map[string]string{"a": "1"}

	// WHEN
	contains := MapContainsAll(haystack, needle)

	// THEN
	require.True(t, contains)
}

func TestMapContainsAllPairsWhenContainsMultiple(t *testing.T) {
	// GIVEN
	haystack := map[string]string{"a": "1", "b": "2"}
	needle := map[string]string{"a": "1", "b": "2"}

	// WHEN
	contains := MapContainsAll(haystack, needle)

	// THEN
	require.True(t, contains)
}

func TestMapContainsAllPairsWhenContainsOneOfMultiple(t *testing.T) {
	// GIVEN
	haystack := map[string]string{"a": "1", "b": "2"}
	needle := map[string]string{"a": "1", "c": "2"}

	// WHEN
	contains := MapContainsAll(haystack, needle)

	// THEN
	require.False(t, contains)
}

func TestMapContainsAllPairsWhenContainsOneOfMultipleValue(t *testing.T) {
	// GIVEN
	haystack := map[string]string{"a": "1", "b": "2"}
	needle := map[string]string{"a": "1", "b": "3"}

	// WHEN
	contains := MapContainsAll(haystack, needle)

	// THEN
	require.False(t, contains)
}

func TestMapContainsAnyPairWhenAllEmpty(t *testing.T) {
	// GIVEN
	haystack := make(map[string]string)
	needle := make(map[string]string)

	// WHEN
	contains := MapContainsAny(haystack, needle)

	// THEN
	require.False(t, contains)
}

func TestMapContainsAnyPairWhenPairsEmpty(t *testing.T) {
	// GIVEN
	haystack := map[string]string{"a": "1", "b": "2"}
	needle := make(map[string]string)

	// WHEN
	contains := MapContainsAny(haystack, needle)

	// THEN
	require.False(t, contains)
}

func TestMapContainsAnyPairWhenNotContains(t *testing.T) {
	// GIVEN
	haystack := map[string]string{"a": "1", "b": "2"}
	needle := map[string]string{"c": "3"}

	// WHEN
	contains := MapContainsAny(haystack, needle)

	// THEN
	require.False(t, contains)
}

func TestMapContainsAnyPairWhenNotContainsValue(t *testing.T) {
	// GIVEN
	haystack := map[string]string{"a": "1", "b": "2"}
	needle := map[string]string{"a": "3"}

	// WHEN
	contains := MapContainsAny(haystack, needle)

	// THEN
	require.False(t, contains)
}

func TestMapContainsAnyPairWhenContainsSingle(t *testing.T) {
	// GIVEN
	haystack := map[string]string{"a": "1", "b": "2"}
	needle := map[string]string{"a": "1"}

	// WHEN
	contains := MapContainsAny(haystack, needle)

	// THEN
	require.True(t, contains)
}

func TestMapContainsAnyPairWhenContainsMultiple(t *testing.T) {
	// GIVEN
	haystack := map[string]string{"a": "1", "b": "2"}
	needle := map[string]string{"a": "1", "b": "2"}

	// WHEN
	contains := MapContainsAny(haystack, needle)

	// THEN
	require.True(t, contains)
}

func TestMapContainsAnyPairWhenContainsOneOfMultiple(t *testing.T) {
	// GIVEN
	haystack := map[string]string{"a": "1", "b": "2"}
	needle := map[string]string{"a": "1", "c": "3"}

	// WHEN
	contains := MapContainsAny(haystack, needle)

	// THEN
	require.True(t, contains)
}

func TestMapContainsAnyPairWhenContainsOneOfMultipleValue(t *testing.T) {
	// GIVEN
	haystack := map[string]string{"a": "1", "b": "2"}
	needle := map[string]string{"a": "3", "b": "2"}

	// WHEN
	contains := MapContainsAny(haystack, needle)

	// THEN
	require.True(t, contains)
}
