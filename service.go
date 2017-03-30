package main

import (
	"sync"
	"time"
)

// Service holds the data of a service
type Service struct {
	Name   string
	Status bool
	Who    string
	Until  time.Time
	mutex  sync.Mutex
	FUntil string
}

func NewService(Name string) *Service {
	var svc *Service
	svc = new(Service)
	svc.Name = Name
	svc.Status = false
	return svc
}

// Updates the service status if the Until date has expired
func (svc Service) UpdateUntil() {
	if svc.Until.Before(time.Now()) {
		svc.mutex.Lock()
		svc.Status = false
		svc.mutex.Unlock()
	}
}
