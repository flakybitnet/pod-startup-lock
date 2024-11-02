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

package config

import (
	"context"
	"github.com/sethvargo/go-envconfig"
	log "log/slog"
	"os"
	"time"
)

type Config struct {
	LockHost     string        `env:"PSL_LOCK_HOST, required"`            // Lock service host
	LockPort     int           `env:"PSL_LOCK_PORT, default=8080"`        // Lock service port
	LockDuration time.Duration `env:"PSL_LOCK_DURATION"`                  // Custom lock duration to request
	Period       time.Duration `env:"PSL_LOCK_CHECK_PERIOD, default=60s"` // Pause between lock acquisition attempts
	Timeout      time.Duration `env:"PSL_LOCK_CHECK_TIMEOUT, default=1s"` // Timeout of lock request
}

func NewConfig(ctx context.Context) Config {
	var conf Config
	err := envconfig.Process(ctx, &conf)
	if err != nil {
		log.Error("cannot not process configuration", err)
		os.Exit(1)
	}
	valid := conf.validate()
	if !valid {
		os.Exit(1)
	}
	log.Info("application configured", log.Any("config", conf))
	return conf
}

func (c *Config) validate() bool {
	valid := true
	if c.Period < 0 {
		log.Error("lock check period is lesser than 0")
		valid = false
	}
	if c.Timeout < 0 {
		log.Error("check timeout is lesser than 0")
		valid = false
	}
	return valid
}
