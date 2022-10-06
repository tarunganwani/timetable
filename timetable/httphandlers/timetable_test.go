package httphandlers_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"timetable/httphandlers"

	"github.com/gorilla/mux"
)

func setup() {
	log.Println("Setting up tests.")
	os.Setenv("DATA_DIR", "/home/tg/WS/timetable_git/timetable/resources")
}

// func teardown(){

// }

type TestGetObject struct {
	inputRoute string
	schoolcode string
	grade      string
	division   string
	outStatus  int
}

func TestTimetableRouteHandler(t *testing.T) {
	setup()

	testGetObjArr := []TestGetObject{
		{
			inputRoute: "/timetable",
			schoolcode: "vvhs",
			grade:      "1",
			division:   "D",
			outStatus:  http.StatusOK,
		},
		{
			inputRoute: "/timetable",
			schoolcode: "vvhs",
			grade:      "",
			division:   "",
			outStatus:  http.StatusNotImplemented,
		},
		{
			inputRoute: "/timetable",
			schoolcode: "vvhs",
			grade:      "2",
			division:   "A",
			outStatus:  http.StatusNotFound,
		},
	}
	for _, testGetObj := range testGetObjArr {

		absIpRoute := testGetObj.inputRoute + "/" + testGetObj.schoolcode +
			"/" + testGetObj.grade + "/" + testGetObj.division
		log.Println("Testing route ", absIpRoute)
		request, err := http.NewRequest("GET", absIpRoute, nil)
		if err != nil {
			t.Fatal(err)
		}

		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"schoolcode": testGetObj.schoolcode,
			"grade":      testGetObj.grade,
			"division":   testGetObj.division,
		}

		// set mux url vars
		request = mux.SetURLVars(request, vars)

		//create response recorder
		respRecorder := httptest.NewRecorder()
		handler := http.HandlerFunc(httphandlers.FetchTimetable)
		handler.ServeHTTP(respRecorder, request)

		if status := respRecorder.Code; status != testGetObj.outStatus {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, testGetObj.outStatus)
		}
	}

}

func TestRegisterHandler(t *testing.T) {
	var router *mux.Router = nil
	if err := httphandlers.RegisterTimetableHandlers(router); err == nil {
		t.Fail()
	}
	router, err := httphandlers.InitializeRouter(httphandlers.RouterConfig{Host: "xyz", Port: "8080"})
	if err != nil {
		t.Fail()
	}
	if err = httphandlers.RegisterTimetableHandlers(router); err != nil {
		t.Fail()
	}
}
