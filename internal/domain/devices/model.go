package devices

import "net/http"

type DevService interface {
	UpdateDev(w http.ResponseWriter, r *http.Request) (*http.Response, error)
	DeleteDev(w http.ResponseWriter, r *http.Request) (*http.Response, error)
	CreateNewDev(w http.ResponseWriter, r *http.Request) (*http.Response, error)
	InitDevices() error
	GetAll(r *http.Request) (*http.Response, error)
	GetOne(r *http.Request) (*http.Response, error)
	GetDeviceStat(r *http.Request) (*http.Response, error)
}

type Devices struct {
	Dev_id      int
	Raw_id      int
	Name        string
	Dev_type    string
	Enabled     bool
	Raw_devices []string
	Send_Status bool
}
