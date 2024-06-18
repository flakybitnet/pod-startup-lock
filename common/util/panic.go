/*
 * Copyright 2018, Oath Inc.
 * Copyright 2024, The PSL (Pod Startup Lock) Authors
 * Licensed under the terms of the MIT license. See LICENSE file in the project root for terms.
 */

package util

func PanicOnError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
