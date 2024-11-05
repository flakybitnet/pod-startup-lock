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
*/

package main

import (
	"context"
	. "flakybit.net/psl/init/client"
	. "flakybit.net/psl/init/config"
	. "flakybit.net/psl/init/service"
	slogenv "github.com/cbrewster/slog-env"
	log "log/slog"
	"os"
)

func main() {
	var err error
	ctx := context.Background()

	logHandler := slogenv.NewHandler(
		log.NewTextHandler(os.Stderr, nil),
		slogenv.WithEnvVarName("PSL_LOG"))
	log.SetDefault(log.New(logHandler))

	conf, err := NewConfig(ctx)
	if err != nil {
		log.ErrorContext(ctx, "failed to configure application", err)
		panic(err)
	}

	lockClient := NewLockClient(conf)
	lockService := NewLockService(conf, lockClient)
	lockService.Run(ctx)
}
