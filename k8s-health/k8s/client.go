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

package k8s

import (
	"context"
	"log"

	. "flakybit.net/psl/common/util"
	. "flakybit.net/psl/k8s-health/config"
	AppsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Client struct {
	k8s kubernetes.Clientset
}

func NewClient(appConfig Config) *Client {
	k8sConfig := getK8sConfig(appConfig)
	k8sClient := *kubernetes.NewForConfigOrDie(k8sConfig)
	return &Client{k8sClient}
}

func (c *Client) GetNodeLabels(nodeName string) map[string]string {
	node := (*RetryOrPanicDefault(func() (interface{}, error) {
		return c.k8s.CoreV1().Nodes().Get(context.TODO(), nodeName, meta.GetOptions{})
	})).(*v1.Node)
	return node.Labels
}

func (c *Client) GetDaemonSets(namespace string) []AppsV1.DaemonSet {
	daemonSetList := (*RetryOrPanicDefault(func() (interface{}, error) {
		return c.k8s.AppsV1().DaemonSets(namespace).List(context.TODO(), meta.ListOptions{})
	})).(*AppsV1.DaemonSetList)
	return daemonSetList.Items
}

func (c *Client) GetNodePods(nodeName string) []v1.Pod {
	opt := meta.ListOptions{}
	opt.FieldSelector = "spec.nodeName=" + nodeName

	podList := (*RetryOrPanicDefault(func() (interface{}, error) {
		return c.k8s.CoreV1().Pods("").List(context.TODO(), opt)
	})).(*v1.PodList)
	return podList.Items
}

func getK8sConfig(appConfig Config) *rest.Config {
	if appConfig.K8sApiBaseUrl != "" {
		log.Printf("K8s baseUrl overrided! Using out-of-cluster k8s client config")
		config := rest.Config{}
		config.Host = appConfig.K8sApiBaseUrl
		config.Insecure = true
		return &config
	} else {
		log.Printf("Using in-cluster k8s client config")
		config, err := rest.InClusterConfig()
		PanicOnError(err)
		return config
	}
}
