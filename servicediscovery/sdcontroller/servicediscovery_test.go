package sdcontroller_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"servicediscovery/sdcontroller"
	"servicediscovery/sdmodel"
	"testing"

	"github.com/gorilla/mux"
)

func setup() *sdcontroller.Controller {
	// Get model
	model := sdmodel.NewSDModel()

	// Initialize controller
	controller, err := sdcontroller.NewSDController(model, sdcontroller.RouterConfig{
		Host:        "",
		Port:        "",
		Certificate: "",
	})
	if err != nil {
		log.Fatalln(err)
	}
	return controller
}

func TestRouteHandlers(t *testing.T) {

	controller := setup()

	fetchaddressRoute := "/fetchaddress"
	log.Println("Testing route ", fetchaddressRoute)
	request, err := http.NewRequest("GET", fetchaddressRoute, nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{
		"servicename": "timetable",
	}
	// set mux url vars
	request = mux.SetURLVars(request, vars)

	//create response recorder
	respRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.FetchServiceAddressHandler)
	handler.ServeHTTP(respRecorder, request)

	if status := respRecorder.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	/////////////////////////////////////////////////////////////////////////////////
	////////////////////// Register timetable service ///////////////////////////////
	/////////////////////////////////////////////////////////////////////////////////

	jsonBody := []byte(`{"svcname":"timetable", "host":"localhost", "port":"8080"}`)
	bodyReader := bytes.NewReader(jsonBody)
	postRequest, err := http.NewRequest("POST", "/heartbeat", bodyReader)
	if err != nil {
		t.Fatal("unexpected error while posting hertbeat")
	}
	handler = http.HandlerFunc(controller.HeartbeatServiceHandler)
	respRecorder1 := httptest.NewRecorder()
	handler.ServeHTTP(respRecorder1, postRequest)

	if status := respRecorder1.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	/////////////////////////////////////////////////////////////////////////////////
	////////////// check fetchaddress for registered service ////////////////////////
	/////////////////////////////////////////////////////////////////////////////////

	request, _ = http.NewRequest("GET", fetchaddressRoute, nil)
	request = mux.SetURLVars(request, vars)
	respRecorder2 := httptest.NewRecorder()
	handler = http.HandlerFunc(controller.FetchServiceAddressHandler)
	handler.ServeHTTP(respRecorder2, request)

	if status := respRecorder2.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	type SDResponse struct {
		Message string
	}
	var sdresponse SDResponse
	var servicedata sdcontroller.ServiceData
	err = json.Unmarshal(respRecorder2.Body.Bytes(), &sdresponse)
	log.Println("response received", respRecorder2.Body.String())
	if err != nil {
		t.Fatal("unexpected error while decoding service discovery response")
	}
	if sdresponse.Message == "" {
		t.Fatal("message can not be nil")
	}
	err = json.Unmarshal([]byte(sdresponse.Message), &servicedata)
	if err != nil {
		t.Fatal("unexpected error while decoding service data")
	}
	if servicedata.Host != "localhost" {
		t.Errorf("Expected host value %#v got %#v", "localhost", servicedata.Host)
	}
	if servicedata.Port != "8080" {
		t.Errorf("Expected port value %#v got %#v", "8080", servicedata.Port)
	}
	if servicedata.Svcname != "timetable" {
		t.Errorf("Expected host value %#v got %#v", "timetable", servicedata.Svcname)
	}

	/////////////////////////////////////////////////////////////////////////////////
	////////////////////// Deegister timetable service //////////////////////////////
	/////////////////////////////////////////////////////////////////////////////////

	jsonBody = []byte(`{"svcname":"timetable", "host":"localhost", "port":"8080"}`)
	bodyReader = bytes.NewReader(jsonBody)
	postRequest, err = http.NewRequest("POST", "/deregister", bodyReader)

	if err != nil {
		t.Fatal("unexpected error while posting hertbeat")
	}
	handler = http.HandlerFunc(controller.DeregisterServiceHandler)
	respRecorder3 := httptest.NewRecorder()
	handler.ServeHTTP(respRecorder3, postRequest)

	if status := respRecorder3.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	/////////////////////////////////////////////////////////////////////////////////
	////////////////////// Deegister negative case //////////////////////////////////
	/////////////////////////////////////////////////////////////////////////////////

	negJsonBody := []byte(`{"svcname":"tiffin_catalog", "host":"localhost", "port":"8080"}`)
	nBodyReader := bytes.NewReader(negJsonBody)
	nPostRequest, err := http.NewRequest("POST", "/deregister", nBodyReader)

	if err != nil {
		t.Fatal("unexpected error while posting deregister request")
	}
	handler = http.HandlerFunc(controller.DeregisterServiceHandler)
	respRecorder4 := httptest.NewRecorder()
	handler.ServeHTTP(respRecorder4, nPostRequest)

	if status := respRecorder4.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	/////////////////////////////////////////////////////////////////////////////////
	////////////////////// Register negative unit test1 /////////////////////////////
	/////////////////////////////////////////////////////////////////////////////////

	//"address" mentioned in request in place of "host"
	n1JsonBody := []byte(`{"svcname":"timetable", "address":"localhost", "port":"8080"}`)
	n1bodyReader := bytes.NewReader(n1JsonBody)
	n1postRequest, err := http.NewRequest("POST", "/heartbeat", n1bodyReader)
	if err != nil {
		t.Fatal("unexpected error while posting hertbeat")
	}
	n1handler := http.HandlerFunc(controller.HeartbeatServiceHandler)
	n1respRecorder1 := httptest.NewRecorder()
	n1handler.ServeHTTP(n1respRecorder1, n1postRequest)

	if status := n1respRecorder1.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	/////////////////////////////////////////////////////////////////////////////////
	////////////////////// Register negative unit test2 /////////////////////////////
	/////////////////////////////////////////////////////////////////////////////////

	//"address" mentioned in request in place of "host"
	n2postRequest, err := http.NewRequest("POST", "/deregister", n1bodyReader)
	if err != nil {
		t.Fatal("unexpected error while posting deregister request")
	}
	n2handler := http.HandlerFunc(controller.DeregisterServiceHandler)
	n2respRecorder := httptest.NewRecorder()
	n2handler.ServeHTTP(n2respRecorder, n2postRequest)

	if status := n2respRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}
