/*
 * Copyright 2018, Oath Inc.
 * Licensed under the terms of the MIT license. See LICENSE file in the project root for terms.
 */

package service

import (
	"fmt"
	"log"
	"net/http"
)

type Service struct {
	host       string
	port       int
	healthFunc func() bool
}

func NewService(host string, port int, healthFunc func() bool) Service {
	return Service{host, port, healthFunc}
}

func (s *Service) Run() {
	log.Print("Starting Http Service...")
	addr := fmt.Sprintf("%s:%v", s.host, s.port)
	http.HandleFunc("/", s.handler)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Panic("Http Service failed to start: ", err)
	}
}

func (s *Service) handler(w http.ResponseWriter, r *http.Request) {
	if s.healthFunc() {
		respondOk(w, r)
	} else {
		respondLocked(w, r)
	}
}

func respondOk(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	log.Printf("Responding to '%v': %v", r.RemoteAddr, status)
	w.WriteHeader(status)
	w.Write([]byte("HealthCheck OK"))
}

func respondLocked(w http.ResponseWriter, r *http.Request) {
	status := http.StatusPreconditionFailed
	log.Printf("Responding to '%v': %v", r.RemoteAddr, status)
	w.WriteHeader(status)
	w.Write([]byte("HealthCheck Failed"))
}
