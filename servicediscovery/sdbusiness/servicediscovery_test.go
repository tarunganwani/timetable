package sdbusiness_test

import (
	"servicediscovery/sdbusiness"
	"servicediscovery/sdmodel"
	"strings"
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
