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
	"log"
	"strings"
)

type ArrayVal []string

func (v *ArrayVal) String() string {
	return fmt.Sprintf("%v", []string(*v))
}

func (v *ArrayVal) Set(value string) error {
	*v = append(*v, value)
	return nil
}

type Pair struct {
	A string
	B string
}

type PairArrayVal struct {
	sep    string
	values []Pair
}

func NewPairArrayVal(sep string) PairArrayVal {
	return PairArrayVal{sep, make([]Pair, 0)}
}

func (v *PairArrayVal) String() string {
	return fmt.Sprintf("%v", v.values)
}

func (v *PairArrayVal) Set(value string) error {
	chunks := strings.Split(value, v.sep)
	if len(chunks) != 2 || chunks[0] == "" || chunks[1] == "" {
		log.Panicf("Failed to parse value: '%s'", value)
	}
	pair := Pair{chunks[0], chunks[1]}
	v.values = append(v.values, pair)
	return nil
}

func (v *PairArrayVal) Get() []Pair {
	return v.values
}
