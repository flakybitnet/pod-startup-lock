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
	"errors"
	"github.com/sethvargo/go-envconfig"
	log "log/slog"
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
	Timeout      time.Duration `env:"TIMEOUT, default=10s"`     // Timeout of health checks
}

func NewConfig(ctx context.Context) (Config, error) {
	var conf Config
	err := envconfig.Process(ctx, &conf)
	if err != nil {
		return conf, err
	}
	err = conf.validate()
	if err != nil {
		return conf, err
	}
	log.Info("application configured", log.Any("config", conf))
	return conf, err
}

func (c *Config) validate() error {
	var parallelLocksError error
	if c.ParallelLocks < 1 {
		parallelLocksError = errors.New("parallel locks is lesser than 0")
	}
	var lockDurationError error
	if c.LockDuration < 0 {
		lockDurationError = errors.New("lock duration is lesser than 0")
	}
	var hcPeriodPassError error
	if c.HealthCheck.PeriodOnPass < 0 {
		hcPeriodPassError = errors.New("period on pass is lesser than 0")
	}
	var hcPeriodFailError error
	if c.HealthCheck.PeriodOnFail < 0 {
		hcPeriodFailError = errors.New("period on fail is lesser than 0")
	}
	var hcEndpointsError error
	if c.HealthCheck.Enabled && len(c.HealthCheck.Endpoints) == 0 {
		hcEndpointsError = errors.New("endpoints health check is enabled, but endpoint list is empty")
	}
	return errors.Join(parallelLocksError, lockDurationError, hcPeriodPassError, hcPeriodFailError, hcEndpointsError)
}
