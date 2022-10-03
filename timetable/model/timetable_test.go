package model

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/tarunganwani/timetable/utility"
)

// func test_string_equals(t *testing.T, field, expected, actual string) {
// 	if actual != expected {
// 		t.Errorf("Invalid %v. Expected %v Got %v", field, expected, actual)
// 	}
// }

func TestTimetableUnmarshal(t *testing.T) {

	jsonData, err := utility.ReadEntireFile("./test_resources/test_vvhs_1_D.json")
	if err != nil {
		t.Fatalf("Error reading file " + err.Error())
	}
	var tt Timetable
	err = json.Unmarshal(jsonData, &tt)
	if err != nil {
		t.Fatalf("Error unmarshalling json " + err.Error())
	}

	ttExpected := Timetable{
		Schoolcode: "vvhs",
		Grade:      "1",
		Division:   "D",
		Days: []Day{
			{
				Name:        "Monday",
				UniformType: "REGULAR",
				Periods: []Period{
					{
						Type:    "1H",
						Subject: "Math",
						Book:    "Book1",
					},
					{
						Type:    "1H",
						Subject: "English",
						Book:    "Book2",
					},
					{
						Type: "SHORT_BREAK",
						Name: "Recess1",
					},
					{
						Type:    "1H",
						Subject: "Computers",
						Book:    "Book3",
					},
				},
			},
			{
				Name:        "Tuesday",
				UniformType: "PE",
				Periods: []Period{
					{
						Type:    "1H",
						Subject: "Math",
						Book:    "Book1",
					},
					{
						Type:    "1H",
						Subject: "PE",
					},
					{
						Type: "SHORT_BREAK",
						Name: "Recess1",
					},
					{
						Type:    "1H",
						Subject: "Computers",
						Book:    "Book3",
					},
				},
			},
		},
	}

	if reflect.DeepEqual(ttExpected, tt) != true {
		t.Errorf("Time table decoding error. \nExpected %v \nGot %v", ttExpected, tt)
	}
}
