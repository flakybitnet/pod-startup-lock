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
	log "log/slog"
)

type HealthCheckService struct {
	conf        Config
	client      *K8sClient
	dsChecker   *DaemonSetChecker
	loadChecker *NodeLoadChecker
}

func NewHealthCheckService(ctx context.Context, conf Config, k8sClient *K8sClient) *HealthCheckService {
	nodeInfo := k8sClient.GetNodeInfo(ctx, conf.NodeName)
	hcSvc := HealthCheckService{
		conf,
		k8sClient,
		NewDaemonSetChecker(conf, k8sClient, nodeInfo),
		NewNodeLoadChecker(conf, k8sClient, nodeInfo),
	}
	log.Info("configured health check service",
		log.Bool("daemon-set-check", conf.DaemonSetHC.Enabled),
		log.Bool("node-load-check", conf.NodeLoadHC.Enabled))
	return &hcSvc
}
func (hcs *HealthCheckService) Run(ctx context.Context) {
	if hcs.conf.DaemonSetHC.Enabled {
		go hcs.dsChecker.Run(ctx)
	}
	if hcs.conf.NodeLoadHC.Enabled {
		go hcs.loadChecker.Run(ctx)
	}
}

func (hcs *HealthCheckService) IsHealthy() bool {
	healthy := hcs.dsChecker.IsHealthy() && hcs.loadChecker.IsHealthy()
	log.Debug("overall health status", log.Bool("status", healthy))
	return healthy
}
