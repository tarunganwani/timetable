package httphandlers

import (
	"log"

	"github.com/gorilla/mux"
)

type RouterConfig struct {
	Host        string
	Port        string
	Certificate string
}

func InitializeRouter(cfg RouterConfig) (*mux.Router, error) {

	log.Println("Initialing http routes..")
	router := mux.NewRouter()
	err := RegisterTimetableHandlers(router)
	return router, err
}
