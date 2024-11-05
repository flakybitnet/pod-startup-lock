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
	. "flakybit.net/psl/k8s-health/client"
	. "flakybit.net/psl/k8s-health/config"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	log "log/slog"
	"math"
	"time"
)

// https://stackoverflow.com/questions/68497673/kubernetes-rest-api-node-cpu-and-ram-usage-in-percentage
// https://stackoverflow.com/questions/52029656/how-to-retrieve-kubernetes-metrics-via-client-go-and-golang

type NodeLoadChecker struct {
	conf            Config
	client          *K8sClient
	nodeCpuCapacity *resource.Quantity
	healthy         bool
}

func NewNodeLoadChecker(conf Config, client *K8sClient, node *core.Node) *NodeLoadChecker {
	cpuCap := node.Status.Capacity.Cpu()
	checker := &NodeLoadChecker{
		conf,
		client,
		cpuCap,
		false,
	}
	log.Info("configured node load checker",
		log.String("cpu-capacity", cpuCap.String()),
		log.Int("threshold", conf.NodeLoadHC.CpuThreshold))
	return checker
}

func (nlc *NodeLoadChecker) IsHealthy() bool {
	if !nlc.conf.NodeLoadHC.Enabled {
		return true
	}
	return nlc.healthy
}

func (nlc *NodeLoadChecker) Run(ctx context.Context) {
	ticker := time.NewTicker(nlc.conf.NodeLoadHC.Period)
	defer ticker.Stop()

	for {
		checkStatus := nlc.check(ctx)
		log.Debug("performed node load health check", log.Bool("healthy", checkStatus))
		nlc.healthy = checkStatus

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return
		}
	}
}

func (nlc *NodeLoadChecker) check(ctx context.Context) bool {
	metrics := nlc.client.GetNodeMetrics(ctx, nlc.conf.NodeName)
	cpuUsageMilli := metrics.Usage.Cpu().MilliValue()
	cpuUsageShare := float64(cpuUsageMilli) / float64(nlc.nodeCpuCapacity.MilliValue())
	cpuUsagePct := int(math.Round(cpuUsageShare * 100))
	log.Debug("node CPU usage", log.Int64("cpu-milli", cpuUsageMilli), log.Int("cpu-pct", cpuUsagePct))

	if cpuUsagePct > nlc.conf.NodeLoadHC.CpuThreshold {
		log.Info("node overload", log.Int("load", cpuUsagePct), log.Int("threshold", nlc.conf.NodeLoadHC.CpuThreshold))
		return false
	}
	return true
}
