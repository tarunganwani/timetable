package httphandlers

import (
	"errors"
	"log"
	"net/http"
	"timetable/business"

	"github.com/gorilla/mux"
)

func register_timetable_handlers(router *mux.Router) error {
	if router == nil {
		msg := "failed to register timetable handlers, router not instantiated"
		log.Println(msg)
		return errors.New(msg)
	}
	router.HandleFunc("/timetable/{schoolcode}/{grade}/{division}", FetchTimetable).Methods("GET")
	log.Println("Handlers registered")
	return nil
}

func FetchTimetable(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if err := recover(); err != nil {
			msg := "error serving request"
			log.Println(msg)
			Bake500Response(w, msg)
		}
	}()

	vars := mux.Vars(r)

	schoolcode := vars["schoolcode"]
	grade := vars["grade"]
	division := vars["division"]

	log.Printf("Serving request for school %s grade %s division %s\n", schoolcode, grade, division)
	if schoolcode == "" || grade == "" || division == "" {
		msg := "Feature not available"
		log.Println(msg)
		Bake501Response(w, msg)
	}

	jsonData, found, err := business.FetchResource(schoolcode, grade, division)
	if !found {
		msg := "resource not found"
		log.Println(msg)
		Bake404Response(w, msg)
		return
	}
	if err != nil {
		msg := "error serving request"
		log.Println(msg)
		Bake500Response(w, msg)
		return
	}
	err = Bake200Response(w, jsonData)
	if err == nil {
		log.Println("success fetching timetable")
	} else {
		log.Println("error creating response :", err.Error())
	}
}
