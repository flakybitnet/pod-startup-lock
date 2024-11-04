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
	"flakybit.net/psl/common/util"
	"flakybit.net/psl/k8s-health/client"
	"flakybit.net/psl/k8s-health/config"
	"fmt"
	"log"
	"time"

	AppsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

type DaemonSetChecker struct {
	conf       config.Config
	client     *client.K8sClient
	nodeLabels map[string]string
	healthy    bool
}

func NewDaemonSetChecker(conf config.Config, client *client.K8sClient) *DaemonSetChecker {
	nodeLabels := client.GetNodeLabels(conf.NodeName)
	return &DaemonSetChecker{conf, client, nodeLabels, false}
}

func (dsc *DaemonSetChecker) IsHealthy() bool {
	return dsc.healthy
}

func (dsc *DaemonSetChecker) Run() {
	for {
		if dsc.check() {
			log.Print("HealthCheck passed")
			dsc.healthy = true
			time.Sleep(dsc.conf.DaemonSetHC.PeriodOnPass)
		} else {
			log.Print("HealthCheck failed")
			dsc.healthy = false
			time.Sleep(dsc.conf.DaemonSetHC.PeriodOnPass)
		}
	}
}

func (dsc *DaemonSetChecker) check() bool {
	log.Print("---")
	log.Print("HealthCheck:")
	daemonSets := dsc.client.GetDaemonSets(dsc.conf.DaemonSetHC.Namespace)
	if dsc.checkAllDaemonSetsReady(daemonSets) {
		return true
	}
	nodePods := dsc.client.GetNodePods(dsc.conf.NodeName)
	return dsc.checkAllDaemonSetsPodsAvailableOnNode(daemonSets, nodePods)
}

func (dsc *DaemonSetChecker) checkAllDaemonSetsReady(daemonSets []AppsV1.DaemonSet) bool {
	for _, ds := range daemonSets {
		if required, reason := dsc.checkRequired(&ds); !required {
			log.Print(reason)
			continue
		}
		status := ds.Status
		if status.DesiredNumberScheduled != status.NumberReady {
			log.Printf("'%v' daemonSet not ready: Desired: '%v', Ready: '%v'",
				ds.Name, status.DesiredNumberScheduled, status.NumberReady)
			return false
		}
		log.Printf("'%v': ok", ds.Name)
	}
	log.Print("All DaemonSets ok")
	return true
}

func (dsc *DaemonSetChecker) checkAllDaemonSetsPodsAvailableOnNode(daemonSets []AppsV1.DaemonSet, pods []v1.Pod) bool {
	for _, ds := range daemonSets {
		if required, reason := dsc.checkRequired(&ds); !required {
			log.Print(reason)
			continue
		}
		log.Printf("'%v' daemonSet: Looking for Pods on node", ds.Name)
		pod, found := findDaemonSetPod(&ds, pods)
		if !found {
			log.Printf("'%v' daemonSet: No Pods found", ds.Name)
			return false
		}
		log.Printf("'%v' daemonSet: Found Pod: '%v'", ds.Name, pod.Name)
		if !isPodReady(pod) {
			return false
		}
	}
	log.Print("All DaemonSets Pods available on node")
	return true
}

func (dsc *DaemonSetChecker) checkRequired(ds *AppsV1.DaemonSet) (bool, string) {
	reason := fmt.Sprintf("'%v' daemonSet Excluded from healthcheck: ", ds.Name)
	if len(dsc.conf.DaemonSetHC.Exclude) > 0 && util.MapContainsAny(ds.Labels, dsc.conf.DaemonSetHC.Exclude) {
		return false, reason + "matches exclude labels"
	}
	if len(dsc.conf.DaemonSetHC.Include) > 0 && !util.MapContainsAll(ds.Labels, dsc.conf.DaemonSetHC.Include) {
		return false, reason + "not matches include labels"
	}
	if dsc.conf.DaemonSetHC.HostNetwork && !ds.Spec.Template.Spec.HostNetwork {
		return false, reason + "not on host network"
	}
	nodeSelector := ds.Spec.Template.Spec.NodeSelector
	if !util.MapContainsAll(dsc.nodeLabels, nodeSelector) {
		return false, reason + "not eligible for scheduling on node"
	}
	return true, fmt.Sprintf("'%v' daemonSet healthcheck required", ds.Name)
}

func findDaemonSetPod(ds *AppsV1.DaemonSet, pods []v1.Pod) (*v1.Pod, bool) {
	for _, pod := range pods {
		if isPodOwnedByDs(&pod, ds) {
			return &pod, true
		}
	}
	return nil, false
}

func isPodReady(pod *v1.Pod) bool {
	if pod.Status.Phase != "Running" {
		log.Printf("'%v' Pod: Not running: Phase: '%v'", pod.Name, pod.Status.Phase)
		return false
	}
	for _, cond := range pod.Status.Conditions {
		if cond.Type == "Ready" && cond.Status == "True" {
			log.Printf("'%v' Pod: Ready", pod.Name)
			return true
		}
	}
	log.Printf("'%v' Pod: Not Ready: '%v'", pod.Name, pod.Status.Conditions)
	return false
}

func isPodOwnedByDs(pod *v1.Pod, ds *AppsV1.DaemonSet) bool {
	for _, ref := range pod.ObjectMeta.OwnerReferences {
		if ds.ObjectMeta.UID == ref.UID {
			return true
		}
	}
	return false
}
