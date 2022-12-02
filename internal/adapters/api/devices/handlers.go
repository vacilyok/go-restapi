package devicehandlers

import (
	"errors"
	"io/ioutil"
	"mediator/internal/adapters/api"
	"mediator/internal/domain/devices"

	"mediator/internal/config"
	"net/http"

	mux "github.com/gorilla"
)

const (
	devicesURL    = "/devices"
	deviceURL     = "/devices/{dev_name}"
	deviceStatURL = "/devices/{dev_name}/{action}"
)

type deviceHandler struct {
	service devices.DevService
}

// create new device handler
func NewDevHandlers(service devices.DevService) api.Handlers {
	return &deviceHandler{
		service: service,
	}
}

// register routes
func (dh *deviceHandler) Register(router *mux.Router) {
	router.HandleFunc(devicesURL, dh.GetAll).Methods("GET")
	router.HandleFunc(deviceStatURL, dh.DeviceStat).Methods("GET")
	router.HandleFunc(devicesURL, dh.Create).Methods("POST")
	router.HandleFunc(deviceURL, dh.GetOne).Methods("GET")
	router.HandleFunc(deviceURL, dh.Update).Methods("PUT")
	router.HandleFunc(devicesURL, dh.Update).Methods("PUT")
	router.HandleFunc(deviceURL, dh.Delete).Methods("DELETE")
	router.HandleFunc(devicesURL, dh.Delete).Methods("DELETE")

}

// http error message response
func (dh *deviceHandler) errorResponse(w http.ResponseWriter, err error, code int) {
	config.Logging.Error(err.Error())
	w.WriteHeader(code)
	w.Write([]byte(err.Error()))
}

// view http response
func (dh *deviceHandler) viewResponse(writer http.ResponseWriter, resp *http.Response) {
	contentType := resp.Header.Get("Content-type")
	writer.Header().Set("Content-Type", contentType)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		dh.errorResponse(writer, err, 500)
		return
	}
	if resp.StatusCode > 205 {
		err = errors.New(string(respBody))
		dh.errorResponse(writer, err, resp.StatusCode)
		return

	}
	writer.WriteHeader(resp.StatusCode)
	writer.Write(respBody)
	if len(respBody) > 0 && len(respBody) < 300 {
		config.Logging.Info(string(respBody))
	}

}

// send request rpc server for get all devices state
func (dh *deviceHandler) GetAll(writer http.ResponseWriter, r *http.Request) {
	resp, err := dh.service.GetAll(r)
	if err != nil {
		dh.errorResponse(writer, err, 500)
		return
	}
	dh.viewResponse(writer, resp)
}

// handler get one device
func (dh *deviceHandler) GetOne(writer http.ResponseWriter, r *http.Request) {
	resp, err := dh.service.GetOne(r)
	if err != nil {
		dh.errorResponse(writer, err, http.StatusBadRequest)
		return
	}
	dh.viewResponse(writer, resp)
}

// handler for device stat
func (dh *deviceHandler) DeviceStat(writer http.ResponseWriter, r *http.Request) {
	stat, err := dh.service.GetDeviceStat(r)
	if err != nil {
		dh.errorResponse(writer, err, http.StatusBadRequest)
		return
	}
	dh.viewResponse(writer, stat)
}

// handler for create new device
func (dh *deviceHandler) Create(writer http.ResponseWriter, r *http.Request) {
	resp, err := dh.service.CreateNewDev(writer, r)
	if err != nil {
		dh.errorResponse(writer, err, http.StatusBadRequest)
		return
	}
	dh.viewResponse(writer, resp)
}

// handler for  update device
func (dh *deviceHandler) Update(writer http.ResponseWriter, r *http.Request) {
	resp, err := dh.service.UpdateDev(writer, r)
	if err != nil {
		dh.errorResponse(writer, err, http.StatusBadRequest)
		return
	}
	dh.viewResponse(writer, resp)
}

// handler for  delete device
func (dh *deviceHandler) Delete(writer http.ResponseWriter, r *http.Request) {
	resp, err := dh.service.DeleteDev(writer, r)
	if err != nil {
		dh.errorResponse(writer, err, http.StatusBadRequest)
		return
	}
	dh.viewResponse(writer, resp)
}
