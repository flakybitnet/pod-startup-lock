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

package service

import (
	"context"
	. "flakybit.net/psl/init/client"
	. "flakybit.net/psl/init/config"
	log "log/slog"
	"time"
)

type LockService struct {
	conf   Config
	client *LockClient
}

func NewLockService(conf Config, client *LockClient) *LockService {
	hcSvc := LockService{
		conf,
		client,
	}
	log.Info("configured lock service")
	return &hcSvc
}

func (ls *LockService) Run(ctx context.Context) {
	ticker := time.NewTicker(ls.conf.Period)
	defer ticker.Stop()

	for {
		success, err := ls.client.AcquireLock(ctx)
		if err != nil {
			log.ErrorContext(ctx, "failed to acquire a lock", log.Any("error", err))
		}
		if success {
			log.Info("lock acquired successfully")
			return
		}

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return
		}
	}
}
