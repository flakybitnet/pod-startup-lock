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
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewEndpointIfHttp(t *testing.T) {
	// GIVEN
	expectedProtocol := "http"

	// WHEN
	actual := ParseEndpoint("http://localhost:1234")

	// THEN
	require.Equal(t, expectedProtocol, actual.Protocol())
	require.True(t, actual.IsHttp())
}

func TestNewEndpointIfHttps(t *testing.T) {
	// GIVEN
	expectedProtocol := "https"

	// WHEN
	actual := ParseEndpoint("https://localhost:1234")

	// THEN
	require.Equal(t, expectedProtocol, actual.Protocol())
	require.True(t, actual.IsHttp())
}

func TestNewEndpointIfTcp(t *testing.T) {
	// GIVEN
	expectedProtocol := "tcp"

	// WHEN
	actual := ParseEndpoint("tcp://localhost:1234")

	// THEN
	require.Equal(t, expectedProtocol, actual.Protocol())
	require.False(t, actual.IsHttp())
}

func TestNewEndpointIfInvalidString(t *testing.T) {
	// GIVEN
	// WHEN
	panicFunc := func() { ParseEndpoint("localhost_1234") }

	// THEN
	require.PanicsWithValue(t, "Endpoint malformed: 'localhost_1234'", panicFunc)
}

func TestNewEndpointIfInvalidPort(t *testing.T) {
	// GIVEN
	// WHEN
	panicFunc := func() { ParseEndpoint("localhost:abcd") }

	// THEN
	require.PanicsWithValue(t, "Endpoint malformed: 'localhost:abcd'", panicFunc)
}

func TestNewEndpointIfInvalidProtocol(t *testing.T) {
	// GIVEN
	// WHEN
	panicFunc := func() { ParseEndpoint("localhost:1234") }

	// THEN
	require.PanicsWithValue(t, "Endpoint malformed: 'localhost:1234'", panicFunc)
}

func TestRawEndpointAddress(t *testing.T) {
	// GIVEN
	expectedAddress := "localhost:1234"

	// WHEN
	actual := ParseEndpoint("tcp://localhost:1234").(RawEndpoint)

	// THEN
	require.Equal(t, expectedAddress, actual.Address())
}

func TestHttpEndpointUrl(t *testing.T) {
	// GIVEN
	expectedUrl := "http://localhost:1234"

	// WHEN
	actual := ParseEndpoint("http://localhost:1234").(HttpEndpoint)

	// THEN
	require.Equal(t, expectedUrl, actual.Url())
}

func TestRawEndpointAddressNoPort(t *testing.T) {
	// GIVEN
	// WHEN
	panicFunc := func() { ParseEndpoint("tcp://localhost") }

	// THEN
	require.PanicsWithValue(t, "Address malformed: 'localhost'", panicFunc)
}

func TestHttpEndpointUrlNoPort(t *testing.T) {
	// GIVEN
	expectedUrl := "http://localhost"

	// WHEN
	actual := ParseEndpoint("http://localhost").(HttpEndpoint)

	// THEN
	require.Equal(t, expectedUrl, actual.Url())
}

func TestRawEndpointString(t *testing.T) {
	// GIVEN
	expected := "localhost:1234"

	// WHEN
	actual := ParseEndpoint("tcp://localhost:1234")

	// THEN
	require.Equal(t, expected, actual.String())
}

func TestHttpEndpointString(t *testing.T) {
	// GIVEN
	expected := "http://localhost:1234"

	// WHEN
	actual := ParseEndpoint("http://localhost:1234")

	// THEN
	require.Equal(t, expected, actual.String())
}
