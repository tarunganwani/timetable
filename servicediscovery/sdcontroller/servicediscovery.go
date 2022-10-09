package sdcontroller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"servicediscovery/sdbusiness"
	"servicediscovery/sdmodel"

	"github.com/gorilla/mux"
	"github.com/tarunganwani/timetable/utility"
)

type ServiceData struct {
	Svcname string `json:"svcname"`
	Host    string `json:"host"`
	Port    string `json:"port"`
}

type RouterConfig struct {
	Host        string
	Port        string
	Certificate string
}

type Controller struct {
	sdmodel *sdmodel.Model
	router  *mux.Router
	cfg     RouterConfig
}

func (controller *Controller) initializeRouter(cfg RouterConfig) error {

	controller.cfg = cfg
	log.Println("Initialing http routes..")
	controller.router = mux.NewRouter()
	err := controller.registerServiceDiscoveryHandlers()
	return err
}

func (controller *Controller) Serve() {
	err := utility.FireHttpServer(controller.cfg.Host, controller.cfg.Port, "", "", controller.router)
	if err != nil {
		log.Fatalf("server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}

func NewSDController(sdmodel *sdmodel.Model, cfg RouterConfig) (*Controller, error) {
	controller := new(Controller)
	if sdmodel == nil {
		return nil, errors.New("model not initialized")
	}
	var err error
	controller.sdmodel = sdmodel
	err = controller.initializeRouter(cfg)
	return controller, err
}

func (controller *Controller) registerServiceDiscoveryHandlers() error {

	log.Println("Registering service discovery handlers")
	router := controller.router
	if router == nil {
		msg := "failed to register service discovery handlers, router not instantiated"
		return errors.New(msg)
	}
	router.HandleFunc("/heartbeat", controller.HeartbeatServiceHandler).Methods("POST")
	router.HandleFunc("/deregister", controller.DeregisterServiceHandler).Methods("POST")
	router.HandleFunc("/fetchaddress/{servicename}", controller.FetchServiceAddressHandler).Methods("GET")
	log.Println("Handlers registered")
	return nil
}

func getRequestData(r *http.Request) (ServiceData, error) {
	decoder := json.NewDecoder(r.Body)
	var data ServiceData
	err := decoder.Decode(&data)
	return data, err
}

func (controller *Controller) HeartbeatServiceHandler(w http.ResponseWriter, r *http.Request) {

	defer func() {
		err := recover()
		if err != nil {
			log.Println("error encountered in proxy handler", err)
			Bake500Response(w, "internal error")
		}
	}()
	log.Println("Got heartbeat")
	data, err := getRequestData(r)
	msg := ""
	if err != nil || data.Svcname == "" || data.Host == "" || data.Port == "" {
		msg = "Invalid input request"
		log.Println(msg)
		if err != nil {
			log.Println(err.Error())
		}
		Bake400Response(w, msg)
		return
	}

	log.Printf("Service details [svc %s host %s port %s] ", data.Svcname, data.Host, data.Port)
	err = sdbusiness.RegisterService(data.Svcname, data.Host, data.Port, controller.sdmodel)
	if err != nil {
		msg = "Error registering service"
		log.Println(msg, err.Error())
		Bake500Response(w, msg)
		return
	}

	msg = "ok"
	log.Println(msg)
	Bake200Response(w, []byte(msg))
}

func (controller *Controller) DeregisterServiceHandler(w http.ResponseWriter, r *http.Request) {

	defer func() {
		err := recover()
		if err != nil {
			log.Println("error encountered in proxy handler", err)
			Bake500Response(w, "internal error")
		}
	}()
	log.Println("Handle deregister request")
	data, err := getRequestData(r)
	msg := ""
	if err != nil || data.Svcname == "" || data.Host == "" || data.Port == "" {
		msg = "Invalid input request"
		log.Println(msg)
		if err != nil {
			log.Println(err.Error())
		}
		Bake400Response(w, msg)
		return
	}

	err = sdbusiness.DeregisterService(data.Svcname, data.Host, data.Port, controller.sdmodel)
	if err != nil {
		msg = "Error deregistering service - probably service is not registered"
		log.Println(msg, err.Error())
		Bake500Response(w, msg)
		return
	}

	msg = "ok"
	log.Println(msg)
	Bake200Response(w, []byte(msg))
}

func (controller *Controller) FetchServiceAddressHandler(w http.ResponseWriter, r *http.Request) {

	defer func() {
		err := recover()
		if err != nil {
			log.Println("error encountered in proxy handler", err)
			Bake500Response(w, "internal error")
		}
	}()

	log.Println("Handle fetch-service request")
	reqvars := mux.Vars(r)
	svcModelData, err := sdbusiness.FetchServiceAddress(reqvars["servicename"], controller.sdmodel)
	svcdata := ServiceData{
		Svcname: svcModelData.Name,
		Host:    svcModelData.Address,
		Port:    svcModelData.Port,
	}
	msg := ""
	if err != nil {
		msg = "Service not found"
		log.Println(msg, err.Error())
		Bake404Response(w, msg)
		return
	}

	svcdataBytes, err := json.Marshal(svcdata)
	if err != nil {
		msg = "Error processing request"
		log.Println(msg, err.Error())
		Bake500Response(w, msg)
		return
	}

	log.Println(string(svcdataBytes))
	Bake200Response(w, []byte(svcdataBytes))
}
