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

This file incorporates work covered by the following copyright and permission notice:
	Copyright (c) 2018, Oath Inc.
	Copyright (c) 2022, The PSL (Pod Startup Lock) Authors

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE.
*/

package service

import (
	"context"
	. "flakybit.net/psl/common/util"
	. "flakybit.net/psl/k8s-health/client"
	. "flakybit.net/psl/k8s-health/config"
	"fmt"
	log "log/slog"
	"time"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
)

type DaemonSetChecker struct {
	conf       Config
	client     *K8sClient
	nodeLabels map[string]string
	healthy    bool
}

func NewDaemonSetChecker(conf Config, client *K8sClient, node *core.Node) *DaemonSetChecker {
	checker := &DaemonSetChecker{conf, client, node.Labels, false}
	log.Info("configured DaemonSet checker")
	return checker
}

func (dsc *DaemonSetChecker) IsHealthy() bool {
	if !dsc.conf.DaemonSetHC.Enabled {
		return true
	}
	return dsc.healthy
}

func (dsc *DaemonSetChecker) Run(ctx context.Context) {
	ticker := time.NewTicker(dsc.conf.DaemonSetHC.PeriodOnFail)
	defer ticker.Stop()

	for {
		checkStatus := dsc.check(ctx)
		if checkStatus != dsc.healthy {
			log.Info("DaemonSet health check status changed",
				log.Bool("old", dsc.healthy),
				log.Bool("new", checkStatus))
			if checkStatus {
				ticker.Reset(dsc.conf.DaemonSetHC.PeriodOnPass)
			} else {
				ticker.Reset(dsc.conf.DaemonSetHC.PeriodOnFail)
			}
		}
		log.Debug("performed DaemonSet health check", log.Bool("healthy", checkStatus))
		dsc.healthy = checkStatus

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return
		}
	}
}

func (dsc *DaemonSetChecker) check(ctx context.Context) bool {
	daemonSets := dsc.client.GetDaemonSets(ctx, dsc.conf.DaemonSetHC.Namespace)
	requiredDaemonSets := dsc.getRequiredDaemonSets(daemonSets)
	if dsc.checkDaemonSetsReady(requiredDaemonSets) {
		return true
	}
	nodePods := dsc.client.GetNodePods(ctx, dsc.conf.NodeName)
	return dsc.checkDaemonSetsPodsAvailableOnNode(daemonSets, nodePods)
}

func (dsc *DaemonSetChecker) getRequiredDaemonSets(daemonSets []apps.DaemonSet) []apps.DaemonSet {
	var requiredDaemonSets []apps.DaemonSet
	for _, ds := range daemonSets {
		required, reason := dsc.checkRequired(&ds)
		if required {
			requiredDaemonSets = append(requiredDaemonSets, ds)
		} else {
			log.Debug("skipping DaemonSet", log.String("reason", reason))
		}
	}
	return requiredDaemonSets
}

func (dsc *DaemonSetChecker) checkRequired(ds *apps.DaemonSet) (bool, string) {
	reason := fmt.Sprintf("'%s/%s' daemonSet Excluded from healthcheck: ", ds.Namespace, ds.Name)
	if len(dsc.conf.DaemonSetHC.Exclude) > 0 && MapContainsAny(ds.Labels, dsc.conf.DaemonSetHC.Exclude) {
		return false, reason + "matches exclude labels"
	}
	if len(dsc.conf.DaemonSetHC.Include) > 0 && !MapContainsAll(ds.Labels, dsc.conf.DaemonSetHC.Include) {
		return false, reason + "not matches include labels"
	}
	if dsc.conf.DaemonSetHC.HostNetwork && !ds.Spec.Template.Spec.HostNetwork {
		return false, reason + "not on host network"
	}
	nodeSelector := ds.Spec.Template.Spec.NodeSelector
	if len(nodeSelector) > 0 && !MapContainsAll(dsc.nodeLabels, nodeSelector) {
		return false, reason + "not eligible for scheduling on node"
	}
	return true, fmt.Sprintf("'%s/%s' daemonSet healthcheck required", ds.Namespace, ds.Name)
}

func (dsc *DaemonSetChecker) checkDaemonSetsReady(daemonSets []apps.DaemonSet) bool {
	for _, ds := range daemonSets {
		status := ds.Status
		if status.DesiredNumberScheduled != status.NumberReady {
			log.Info("DaemonSet is not ready",
				log.String("daemon-set", ds.Name),
				log.Int("desired", int(status.DesiredNumberScheduled)),
				log.Int("ready", int(status.NumberReady)))
			return false
		}
		log.Debug("DaemonSet is ready", log.String("daemon-set", ds.Name))
	}
	log.Debug("all DaemonSets are ready")
	return true
}

func (dsc *DaemonSetChecker) checkDaemonSetsPodsAvailableOnNode(daemonSets []apps.DaemonSet, pods []core.Pod) bool {
	for _, ds := range daemonSets {
		log.Debug("looking for pods on node", log.String("daemon-set", ds.Name))
		pod, found := findDaemonSetPod(&ds, pods)
		if !found {
			log.Info("no pod found", log.String("daemon-set", ds.Name))
			return false
		}
		log.Debug("pod found", log.String("daemon-set", ds.Name), log.String("pod", pod.Name))
		if !isPodReady(pod) {
			log.Info("pod is not ready", log.String("daemon-set", ds.Name), log.String("pod", pod.Name))
			return false
		}
	}
	log.Debug("all DaemonSets pods are available on node")
	return true
}

func findDaemonSetPod(ds *apps.DaemonSet, pods []core.Pod) (*core.Pod, bool) {
	for _, pod := range pods {
		if isPodOwnedByDs(&pod, ds) {
			return &pod, true
		}
	}
	return nil, false
}

func isPodOwnedByDs(pod *core.Pod, ds *apps.DaemonSet) bool {
	for _, ref := range pod.ObjectMeta.OwnerReferences {
		if ds.ObjectMeta.UID == ref.UID {
			return true
		}
	}
	return false
}

func isPodReady(pod *core.Pod) bool {
	if pod.Status.Phase != "Running" {
		log.Debug("pod is not running",
			log.String("pod", pod.Name),
			log.String("phase", string(pod.Status.Phase)))
		return false
	}
	for _, cond := range pod.Status.Conditions {
		if cond.Type == "Ready" && cond.Status == "True" {
			log.Debug("pod is ready", log.String("pod", pod.Name))
			return true
		}
	}
	log.Debug("pod is not ready",
		log.String("pod", pod.Name),
		log.Any("conditions", pod.Status.Conditions))
	return false
}
