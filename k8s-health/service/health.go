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
	"flakybit.net/psl/k8s-health/client"
	"flakybit.net/psl/k8s-health/config"
)

type HealthChecker interface {
	IsHealthy() bool
}

type HealthCheckService struct {
	conf        config.Config
	client      *client.K8sClient
	dsChecker   *DaemonSetChecker
	loadChecker *NodeLoadChecker
}

func NewHealthCheckService(conf config.Config, k8sClient *client.K8sClient) *HealthCheckService {
	nodeInfo := k8sClient.GetNodeInfo(conf.NodeName)
	hcSvc := HealthCheckService{
		conf,
		k8sClient,
		NewDaemonSetChecker(conf, k8sClient, nodeInfo),
		NewNodeLoadChecker(conf, k8sClient, nodeInfo),
	}
	if conf.DaemonSetHC.Enabled {
		go hcSvc.dsChecker.Run()
	}
	if conf.NodeLoadHC.Enabled {
		go hcSvc.loadChecker.Run()
	}
	return &hcSvc
}

func (hcs *HealthCheckService) IsHealthy() bool {
	return hcs.dsChecker.IsHealthy() && hcs.loadChecker.IsHealthy()
}
