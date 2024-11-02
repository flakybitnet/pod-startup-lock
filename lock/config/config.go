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
	BindHost      string            `env:"PSL_BIND_HOST"`                  // Host/Ip to bind
	BindPort      int               `env:"PSL_BIND_PORT, default=8080"`    // Port to bind
	ParallelLocks int               `env:"PSL_PARALLEL_LOCKS, default=1"`  // Number of locks allowed to acquire simultaneously
	LockDuration  time.Duration     `env:"PSL_LOCK_DURATION, default=10s"` // Default lock duration
	HealthCheck   HealthCheckConfig `env:", prefix=PSL_HC_"`
}

type HealthCheckConfig struct {
	Enabled      bool          `env:"ENABLED, default=false"`
	Endpoints    []string      `env:"ENDPOINTS"`                // Health check tcp endpoint, host:port
	PeriodOnFail time.Duration `env:"PERIOD_FAIL, default=10s"` // Pause between endpoint health checks if previous failed
	PeriodOnPass time.Duration `env:"PERIOD_PASS, default=60s"` // Pause between endpoint health checks if previous succeeded
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
	if c.ParallelLocks < 1 {
		log.Error("parallel locks is lesser than 0")
		valid = false
	}
	if c.LockDuration < 0 {
		log.Error("lock duration is lesser than 0")
		valid = false
	}
	if c.HealthCheck.PeriodOnPass < 0 {
		log.Error("period on pass is lesser than 0")
		valid = false
	}
	if c.HealthCheck.PeriodOnFail < 0 {
		log.Error("period on fail is lesser than 0")
		valid = false
	}
	if c.HealthCheck.Enabled && len(c.HealthCheck.Endpoints) == 0 {
		log.Error("endpoints health check is enabled, but endpoint list is empty")
		valid = false
	}
	return valid
}
