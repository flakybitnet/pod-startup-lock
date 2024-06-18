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
	endpointChecker := healthcheck.NewHealthChecker(conf, k8sClient)
	srv := service.NewService(conf.Host, conf.Port, endpointChecker.HealthFunction())

	go srv.Run()
	go endpointChecker.Run()

	select {} // Wait forever and let child goroutines run
}
