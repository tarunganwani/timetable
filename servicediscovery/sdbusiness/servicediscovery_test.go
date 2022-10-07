package sdbusiness_test

import (
	"math/rand"
	"servicediscovery/sdbusiness"
	"servicediscovery/sdmodel"
	"strings"
	"sync"
	"testing"
)

func TestRegisterAndDergister(t *testing.T) {

	svcname := "timetable"
	addr := "localhost"
	port := "8080"
	port2 := "8081"

	err := sdbusiness.RegisterService(svcname, addr, port, nil)
	//uninitialized model should lead to an error
	if err == nil {
		t.Error("should return invalid model error")
	}
	model := sdmodel.NewSDModel()
	err = sdbusiness.RegisterService(svcname, addr, port, model)
	_ = sdbusiness.RegisterService(svcname, addr, port2, model)
	if err != nil {
		t.Error("failed tor register", err)
	}
	svcdata, err := sdbusiness.FetchServiceAddress(svcname, model)
	if err != nil {
		t.Fatal("failed to fetch address", err)
	}
	if svcdata.Address != addr {
		t.Errorf("Expected %#v Got %#v", addr, svcdata.Address)
	}
	if svcdata.Port != port {
		t.Errorf("Expected %#v Got %#v", port, svcdata.Port)
	}
	err = sdbusiness.DeregisterService(svcname, addr, port, model)
	if err != nil {
		t.Fatal("failed to deregister", err)
	}
	svcdata, err = sdbusiness.FetchServiceAddress(svcname, model)
	if err != nil {
		t.Fatal("failed to fetch address", err)
	}
	if svcdata.Address != addr {
		t.Errorf("Expected %#v Got %#v", addr, svcdata.Address)
	}
	if svcdata.Port != port2 {
		t.Errorf("Expected %#v Got %#v", port2, svcdata.Port)
	}
	_ = sdbusiness.DeregisterService(svcname, addr, port2, model)
	_, err = sdbusiness.FetchServiceAddress(svcname, model)
	if err == nil {
		t.Error("no data present, should get an error")
	}
}

func TestNegativeFlows(t *testing.T) {

	svcname := "timetable"
	addr := "localhost"
	port := "8080"
	err := sdbusiness.DeregisterService(svcname, addr, port, nil)
	if err == nil {
		t.Error("should get unitialized model error")
	}

	model := sdmodel.NewSDModel()
	err = sdbusiness.DeregisterService(svcname, addr, port, model)
	if err == nil {
		t.Fatal("Should get not registered error")
	}
	if !strings.Contains(err.Error(), "not registered") {
		t.Fatal("Should get not registered error")
	}

	_ = sdbusiness.RegisterService(svcname, addr, port, model)
	someUnregisteredPort := "9000"
	err = sdbusiness.DeregisterService(svcname, addr, someUnregisteredPort, model)
	if err == nil {
		t.Fatal("Should get not registered error")
	}
	if !strings.Contains(err.Error(), "not registered") {
		t.Fatal("Should get not registered error")
	}
	_, err = sdbusiness.FetchServiceAddress(svcname, nil)
	if err == nil {
		t.Fatal("Should get not uninitialized model error")
	}
	someOtherSvc := "uniform_catalog"
	_, err = sdbusiness.FetchServiceAddress(someOtherSvc, model)
	if err == nil {
		t.Fatal("Should get not service not found error")
	}
	if !strings.Contains(err.Error(), "not registered") {
		t.Fatal("Should get not registered error")
	}
}

func TestConcurrency(t *testing.T) {

	sdSet1 := sdmodel.ServiceData{Name: "timetable", Address: "localhost", Port: "8000"}
	sdSet2 := sdmodel.ServiceData{Name: "timetable", Address: "localhost", Port: "8001"}
	sdSet3 := sdmodel.ServiceData{Name: "uniform_catalog", Address: "localhost", Port: "9000"}
	sdSet4 := sdmodel.ServiceData{Name: "uniform_catalog", Address: "localhost", Port: "9001"}

	const (
		REGISTER   = 0
		DEREGISTER = 1
		FETCH      = 2
		NUMBER_OPS = 3
	)
	model := sdmodel.NewSDModel()
	DoRandomOp := func(sdSet sdmodel.ServiceData, model *sdmodel.Model, t *testing.T) {

		defer func() {
			if err := recover(); err != nil {
				t.Fatal("Program panicked", err)
			}
		}()
		op := rand.Int() % NUMBER_OPS
		switch op {
		case REGISTER:
			err := sdbusiness.RegisterService(sdSet.Name, sdSet.Address, sdSet.Port, model)
			if err != nil {
				t.Error("Register error:", err)
			}
		case DEREGISTER:
			err := sdbusiness.DeregisterService(sdSet.Name, sdSet.Address, sdSet.Port, model)
			if err != nil && strings.Contains(err.Error(), "not registered") == false {
				t.Error("Deregister error:", err)
			}
		case FETCH:
			_, err := sdbusiness.FetchServiceAddress(sdSet.Name, model)
			if err != nil && strings.Contains(err.Error(), "not registered") == false {
				t.Error("Fetch error:", err)
			}

		}
	}

	loopOverRandomOp := func(model *sdmodel.Model, t *testing.T, sdSet sdmodel.ServiceData, N int, wg *sync.WaitGroup) {
		for i := 0; i < N; i++ {
			DoRandomOp(sdSet, model, t)
		}
		wg.Done()
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go loopOverRandomOp(model, t, sdSet1, 100000, wg)
	go loopOverRandomOp(model, t, sdSet2, 100000, wg)
	go loopOverRandomOp(model, t, sdSet3, 100000, wg)
	go loopOverRandomOp(model, t, sdSet4, 100000, wg)
	wg.Wait()
}
