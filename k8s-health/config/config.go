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
	BindHost    string                     `env:"PSL_BIND_HOST"`               // Address to bind
	BindPort    int                        `env:"PSL_BIND_PORT, default=8080"` // Port to bind
	NodeName    string                     `env:"PSL_NODE_NAME, required"`     // K8s node name which the current app instance runs on
	K8sApiUrl   string                     `env:"PSL_K8S_API_URL"`             // K8s API URL, for out-of-cluster usage only
	DaemonSetHC DaemonSetHealthCheckConfig `env:", prefix=PSL_HC_DAEMONSET_"`
	NodeLoadHC  NodeLoadHealthCheckConfig  `env:", prefix=PSL_HC_NODELOAD_"`
}

type DaemonSetHealthCheckConfig struct {
	Enabled      bool              `env:"ENABLED, default=true"`
	Namespace    string            `env:"NAMESPACE"`                   // K8s Namespace to check DaemonSets in, blank for all namespaces
	HostNetwork  bool              `env:"HOST_NETWORK, default=false"` // Host network DaemonSets only
	Include      map[string]string `env:"INCLUDE_LABELS"`              // Include DaemonSet labels, "label1:value1,label2:value2"
	Exclude      map[string]string `env:"EXCLUDE_LABELS"`              // Exclude DaemonSet labels, "label1:value1,label2:value2"
	PeriodOnFail time.Duration     `env:"PERIOD_FAIL, default=10s"`    // Period of health checks if previous failed
	PeriodOnPass time.Duration     `env:"PERIOD_PASS, default=60s"`    // Period of health checks if previous succeeded
}

type NodeLoadHealthCheckConfig struct {
	Enabled      bool          `env:"ENABLED, default=false"`
	CpuThreshold int           `env:"CPU_THRESHOLD, default=80"` // Node CPU utilisation in percent above which it is treated as unhealthy
	Period       time.Duration `env:"PERIOD, default=10s"`       // Period of health checks
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
	var dsIncludeExcludeError error
	if len(c.DaemonSetHC.Include) > 0 && len(c.DaemonSetHC.Exclude) > 0 {
		dsIncludeExcludeError = errors.New("cannot specify both Included and Excluded DaemonSet ")
	}
	var dsPeriodPassError error
	if c.DaemonSetHC.PeriodOnPass < 0 {
		dsPeriodPassError = errors.New("period of success DaemonSet check is lesser than 0")
	}
	var dsPeriodFailError error
	if c.DaemonSetHC.PeriodOnFail < 0 {
		dsPeriodFailError = errors.New("period of failed DaemonSet check is lesser than 0")
	}
	var nlThresholdError error
	if c.NodeLoadHC.CpuThreshold < 0 || c.NodeLoadHC.CpuThreshold > 100 {
		nlThresholdError = errors.New("cpu threshold of node load check is out of interval [0, 100]")
	}
	var nlPeriodError error
	if c.NodeLoadHC.Period < 0 {
		nlPeriodError = errors.New("period of node load check is lesser than 0")
	}
	return errors.Join(dsIncludeExcludeError, dsPeriodPassError, dsPeriodFailError, nlThresholdError, nlPeriodError)
}
