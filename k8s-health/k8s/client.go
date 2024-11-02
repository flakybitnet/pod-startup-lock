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
	. "flakybit.net/psl/k8s-health/config"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
	"log"
	"time"
)

var defaultRetry = wait.Backoff{
	Duration: 1 * time.Second,
	Factor:   2.0,
	Jitter:   0.1,
	Steps:    5,
}

var defaultRetriable = func(error) bool {
	return true
}

type Client struct {
	k8s kubernetes.Clientset
}

func NewClient(appConfig Config) *Client {
	k8sConfig := getK8sConfig(appConfig)
	k8sClient := *kubernetes.NewForConfigOrDie(k8sConfig)
	return &Client{k8sClient}
}

func (c *Client) GetNodeLabels(nodeName string) map[string]string {
	var node *core.Node
	retryOnError(func() error {
		var err error
		node, err = c.k8s.CoreV1().Nodes().Get(context.TODO(), nodeName, meta.GetOptions{})
		return err
	})
	return node.Labels
}

func (c *Client) GetDaemonSets(namespace string) []apps.DaemonSet {
	var daemonSets *apps.DaemonSetList
	retryOnError(func() error {
		var err error
		daemonSets, err = c.k8s.AppsV1().DaemonSets(namespace).List(context.TODO(), meta.ListOptions{})
		return err
	})
	return daemonSets.Items
}

func (c *Client) GetNodePods(nodeName string) []core.Pod {
	opt := meta.ListOptions{}
	opt.FieldSelector = "spec.nodeName=" + nodeName

	var pods *core.PodList
	retryOnError(func() error {
		var err error
		pods, err = c.k8s.CoreV1().Pods("").List(context.TODO(), opt)
		return err
	})
	return pods.Items
}

func getK8sConfig(appConfig Config) *rest.Config {
	if appConfig.K8sApiBaseUrl != "" {
		log.Printf("K8s baseUrl overrided! Using out-of-cluster k8s client config")
		config := rest.Config{}
		config.Host = appConfig.K8sApiBaseUrl
		config.Insecure = true
		return &config
	}

	log.Printf("Using in-cluster k8s client config")
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	return config
}

func retryOnError(fn func() error) {
	err := retry.OnError(defaultRetry, defaultRetriable, fn)
	if err != nil {
		panic(err)
	}
}
