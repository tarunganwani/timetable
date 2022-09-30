package httphandlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type RouterConfig struct {
	Host        string
	Port        string
	Certificate string
}

func InitializeRouter(cfg RouterConfig) error {

	log.Println("Initialing http routes..")
	router := mux.NewRouter()
	err := register_timetable_handlers(router)
	if err != nil {
		return err
	}
	addr := cfg.Host + ":" + cfg.Port
	http.ListenAndServe(addr, router)
	return nil
}
