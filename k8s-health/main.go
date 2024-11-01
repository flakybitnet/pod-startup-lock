/*
 * Copyright 2018, Oath Inc.
 * Copyright 2024, The PSL (Pod Startup Lock) Authors
 * Licensed under the terms of the MIT license. See LICENSE file in the project root for terms.
 */

package main

import (
	"flakybit.net/psl/k8s-health/config"
	"flakybit.net/psl/k8s-health/healthcheck"
	"flakybit.net/psl/k8s-health/k8s"
	"flakybit.net/psl/k8s-health/service"
)

func main() {
	conf := config.Parse()
	conf.Validate()

	k8sClient := k8s.NewClient(conf)
	healthChecker := healthcheck.NewHealthChecker(conf, k8sClient)
	srv := service.NewService(conf.Host, conf.Port, healthChecker.HealthFunction())

	go srv.Run()
	go healthChecker.RunDaemonSetsChecks()
	if conf.NodeCpuLoadThreshold > 0 { // enabled
		go healthChecker.RunLoadChecks()
	}

	select {} // Wait forever and let child goroutines run
}
