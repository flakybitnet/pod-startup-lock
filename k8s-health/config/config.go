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
	BindHost    string                     `env:"PSL_BIND_HOST"`               // Host/Ip to bind
	BindPort    int                        `env:"PSL_BIND_PORT, default=8080"` // Port to bind
	NodeName    string                     `env:"PSL_NODE_NAME, required"`     // K8s node name which the current app instance runs on
	K8sApiUrl   string                     `env:"PSL_K8S_API_URL"`             // K8s API URL, for out-of-cluster usage only
	DaemonSetHC DaemonSetHealthCheckConfig `env:", prefix=PSL_HC_DAEMONSET_"`
}

type DaemonSetHealthCheckConfig struct {
	Enabled      bool              `env:"ENABLED, default=true"`
	Namespace    string            `env:"NAMESPACE"`                   // K8s Namespace to check DaemonSets in, blank for all namespaces
	HostNetwork  bool              `env:"HOST_NETWORK, default=false"` // Host network DaemonSets only
	Include      map[string]string `env:"INCLUDE_LABELS"`              // Include DaemonSet labels, label:value
	Exclude      map[string]string `env:"EXCLUDE_LABELS"`              // Exclude DaemonSet labels, label:value
	PeriodOnFail time.Duration     `env:"PERIOD_FAIL, default=10s"`    // Pause between DaemonSet health checks if previous failed
	PeriodOnPass time.Duration     `env:"PERIOD_PASS, default=60s"`    // Pause between DaemonSet health checks if previous succeeded
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
	if len(c.DaemonSetHC.Include) > 0 && len(c.DaemonSetHC.Exclude) > 0 {
		log.Error("cannot specify both Included and Excluded DaemonSet labels, choose one")
		valid = false
	}
	if c.DaemonSetHC.PeriodOnPass < 0 {
		log.Error("period on pass is lesser than 0")
		valid = false
	}
	if c.DaemonSetHC.PeriodOnFail < 0 {
		log.Error("period on fail is lesser than 0")
		valid = false
	}
	return valid
}
