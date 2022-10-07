package sdbusiness

import (
	"errors"
	"fmt"
	"servicediscovery/sdmodel"
	"time"
)

//TODO MAJOR - create a go routine for maintainig the service registry for alive connections
//TODO MAJOR - check if the service data slice access needs to be protected by mutex

// Utility function - Remove from slice without caring about order
func Remove(s []sdmodel.ServiceData, i int) []sdmodel.ServiceData {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// func finddata
func findSliceIndex(srvData []sdmodel.ServiceData, address, port string) int {
	for idx, val := range srvData {
		if val.Address == address && val.Port == port {
			return idx
		}
	}
	return -1
}

func createAndGetNewServiceData(svcname, address, port string) sdmodel.ServiceData {

	newdata := sdmodel.ServiceData{
		Name:                  svcname,
		Address:               address,
		Port:                  port,
		LastHeartbeatReceived: time.Now(),
	}
	return newdata
}

func RegisterService(svcname, address, port string, modelPtr *sdmodel.Model) error {

	//TODO consider validating host and port.
	//In case these are not valid, maybe throw an error

	if modelPtr == nil {
		return errors.New("invalid model")
	}

	if svcname == "" || port == "" || address == "" {
		return errors.New("invalid request")
	}
	srvdata, found := modelPtr.ServiceDiscoveryMap[svcname]
	if !found {
		newdata := createAndGetNewServiceData(svcname, address, port)
		modelPtr.ServiceDiscoveryMap[svcname] = append(modelPtr.ServiceDiscoveryMap[svcname], newdata)
	} else {
		//check if address and port do not already exist in the list
		idx := findSliceIndex(srvdata, address, port)
		//in case not, add it to the slice
		if idx == -1 {
			newdata := createAndGetNewServiceData(svcname, address, port)
			modelPtr.ServiceDiscoveryMap[svcname] = append(modelPtr.ServiceDiscoveryMap[svcname], newdata)
		} else {
			//otherwise just update the timestamp (treat it as a heart beat)
			modelPtr.ServiceDiscoveryMap[svcname][idx].LastHeartbeatReceived = time.Now()
		}
	}
	return nil
}

func DeregisterService(svcname, address, port string, modelPtr *sdmodel.Model) error {

	if modelPtr == nil {
		return errors.New("invalid model")
	}

	if svcname == "" || port == "" || address == "" {
		return errors.New("invalid request")
	}
	srvdata, found := modelPtr.ServiceDiscoveryMap[svcname]
	if !found {

		msg := fmt.Sprintf("%s service not registered ", svcname)
		return errors.New(msg)
	}
	idx := findSliceIndex(srvdata, address, port)
	if idx == -1 {
		msg := fmt.Sprintf("%s service not registered [host %s port %s] ", svcname, address, port)
		return errors.New(msg)
	}
	//remove the service and host from slice
	newSrvdata := Remove(srvdata, idx)
	modelPtr.ServiceDiscoveryMap[svcname] = newSrvdata
	return nil
}

func FetchServiceAddress(svcname string, modelPtr *sdmodel.Model) (sdmodel.ServiceData, error) {

	if modelPtr == nil {
		return sdmodel.ServiceData{}, errors.New("invalid model")
	}
	srvdata, found := modelPtr.ServiceDiscoveryMap[svcname]
	if !found || srvdata == nil || len(srvdata) == 0 {
		msg := fmt.Sprintf("%s service not registered ", svcname)
		return sdmodel.ServiceData{}, errors.New(msg)
	}

	//TODO load balancer logic could go here in the future
	//For now just return the first host found in the slice
	return srvdata[0], nil
}
