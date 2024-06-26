/*
 * Copyright 2018, Oath Inc.
 * Copyright 2024, The PSL (Pod Startup Lock) Authors
 * Licensed under the terms of the MIT license. See LICENSE file in the project root for terms.
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
