package business_test

import (
	"os"
	"testing"
	"timetable/business"
)

func TestFetchResource(t *testing.T) {
	//setup
	os.Setenv("DATA_DIR", "/home/tg/WS/timetable_git/timetable/resources")

	_, found, err := business.FetchResource("vvhs", "1", "D")
	if err != nil {
		t.Fail()
	}
	if found == false {
		t.Fail()
	}

	_, found, err = business.FetchResource("euro", "1", "D")
	if err == nil {
		t.Fail()
	}
	if found == true {
		t.Fail()
	}

}
