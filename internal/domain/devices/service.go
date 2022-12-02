package devices

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	devicesstorage "mediator/internal/adapters/db/devicestorage"
	"mediator/internal/config"
	"net/http"
	"strconv"

	mux "github.com/gorilla"
)

const devicesURL = "/devices"

type service struct {
	storage devicesstorage.DevStorage
}

// Create new service for device
func NewService(storage devicesstorage.DevStorage) DevService {
	return &service{
		storage: storage,
	}
}

// ***************************************************************************************************
func (s *service) createResponse(msg string, code int) *http.Response {
	response := &http.Response{
		Status:        strconv.Itoa(code),
		StatusCode:    code,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(bytes.NewBufferString(msg)),
		ContentLength: int64(len(msg)),
		Header:        make(http.Header, 0),
	}
	return response
}

// ***************************************************************************************************
// Service create new device
func (s *service) CreateNewDev(w http.ResponseWriter, r *http.Request) (*http.Response, error) {
	var (
		err error
	)
	respBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.New("fail read json body at create new device")
	}
	config.Logging.LogRest(string(respBody), "/devices", "POST")
	postDevs := make(map[string][]devicesstorage.PostDevice)
	err = json.Unmarshal(respBody, &postDevs)
	if err != nil {
		return nil, err
	}

	if devices, ok := postDevs["devices"]; ok {
		for _, dev := range devices {
			if dev.Type == "raw" {
				continue
			}
			if dev_exists := s.storage.CheckExistDevice(dev.Name); dev_exists {
				err_message := "Device with name " + dev.Name + " alredy exists in db"
				config.Logging.Warning(err_message)
				return nil, errors.New(err_message)
			}
			if slave_exists := s.storage.CheckExistDevice(dev.Slave); !slave_exists {
				err_message := "Slave device with name " + dev.Slave + " not found in db"
				config.Logging.Error(err_message)
				return nil, errors.New(err_message)
			}

			if dev.Type == "vlan" {
				virtDev := make(map[string]interface{})
				virtDevs := make(map[string][]map[string]interface{})
				virtDev["name"] = dev.Name
				virtDev["enabled"] = dev.Enabled
				virtDev["type"] = dev.Type
				virtDev["slave"] = dev.Slave
				virtDev["vlan_id"] = dev.Vlan_id
				virtDevs["devices"] = append(virtDevs["devices"], virtDev)
				response, _ := s.postRequestToRPC(virtDevs)
				if response.StatusCode > 201 {
					return response, nil
				}
				err = s.storage.InsertVlan(&dev)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return s.createResponse("Devices created successfully", 200), nil
}

// ***************************************************************************************************
// Send POST request to RPC
func (s *service) postRequestToRPC(obj map[string][]map[string]interface{}) (*http.Response, error) {
	jsonString, _ := json.Marshal(obj)
	// config.Logging.LogRest(string(jsonString), "/devices", "POST")

	url := "http://" + config.Params.RPCHost + ":" + strconv.Itoa(config.Params.RPCPort) + "/devices"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonString))
	if err != nil {
		config.Logging.Error("POST  request to rpc server is fail")
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		config.Logging.Error("Fail response for POST new device " + err.Error())
		return nil, err
	}
	return response, nil
}

// Service for update device status
func (s *service) UpdateDev(w http.ResponseWriter, r *http.Request) (*http.Response, error) {
	var dev_name string
	vars := mux.Vars(r)

	if _, ok := vars["dev_name"]; ok {
		dev_name = vars["dev_name"]
		err := s.UpdateDevByName(dev_name, r)
		if err != nil {
			config.Logging.Error(err.Error())
			return nil, err

		}
		return nil, errors.New(" this route not work")
	}
	resp, err := s.UpdateGroupDev(r)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ***************************************************************************************************
// put request to RPC server
func (s *service) putGroupDevRequest(devItem map[string]interface{}) (*http.Response, error) {

	groupDev := make(map[string][]map[string]interface{})
	groupDev["devices"] = append(groupDev["devices"], devItem)
	jsonString, _ := json.Marshal(groupDev)

	config.Logging.LogRest(string(jsonString), "/devices", "PUT")

	url := "http://" + config.Params.RPCHost + ":" + strconv.Itoa(config.Params.RPCPort) + "/devices"
	request, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonString))
	if err != nil {
		config.Logging.Error("put request to rpc server is fail")
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return response, err
}

// ***************************************************************************************************
// Update one dev by name
func (s *service) UpdateDevByName(dev_name string, r *http.Request) error {
	respBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		config.Logging.Error(err.Error())
		return err
	}
	devItem := make(map[string]interface{})
	err = json.Unmarshal(respBody, &devItem)
	if err != nil {
		return err
	}
	if _, ok := devItem["name"]; !ok {
		err_msg := errors.New(" wrong json format. Missing field name")
		config.Logging.Error(err_msg.Error())
		return err_msg
	}
	s.putGroupDevRequest(devItem)
	return nil
}

// ***************************************************************************************************
// update array devices
func (s *service) UpdateGroupDev(r *http.Request) (*http.Response, error) {
	var resp *http.Response
	respBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		config.Logging.Error(err.Error())
		return nil, err
	}
	groupDev := make(map[string][]map[string]interface{})
	err = json.Unmarshal(respBody, &groupDev)
	if err != nil {
		config.Logging.Error(err.Error())
		return nil, err
	}
	if _, ok := groupDev["devices"]; !ok {
		err_msg := errors.New(" wrong json format. Missing field type")
		config.Logging.Error(err_msg.Error())
		return nil, err_msg
	}
	for _, dev := range groupDev["devices"] {
		devName := dev["name"].(string)
		if dev_exists := s.storage.CheckExistDevice(devName); !dev_exists {
			err_message := "Device with name " + devName + " not exists"

			config.Logging.Error(err_message)
			return nil, errors.New(err_message)
		}
		if devSlave, ok := dev["slave"]; ok {
			if slave_exists := s.storage.CheckExistDevice(devSlave.(string)); !slave_exists {
				err_message := "Slave device with name " + devSlave.(string) + " not exists"
				config.Logging.Error(err_message)
				return nil, errors.New(err_message)
			}
		}
		if dst, ok := dev["dst"]; ok {
			if dst != nil {
				if slave_exists := s.storage.CheckExistDevice(dst.(string)); !slave_exists {
					err_message := "dst device with name " + dst.(string) + " not exists"
					config.Logging.Error(err_message)
					return nil, errors.New(err_message)
				}
			}
		}
		resp, _ = s.putGroupDevRequest(dev)
		if resp.StatusCode == 200 {
			err := s.storage.UpdateDbDev(dev)
			if err != nil {
				return nil, err
			}
		}
	}
	return resp, nil
}

// ***************************************************************************************************
// Get all device from db and send device param to RPC
func (s *service) sendDevState() {
	devices := s.storage.GetAllDev()
	// create devices
	for _, dev := range devices {

		if dev.Type == "vlan" {
			virtDev := make(map[string]interface{})
			virtDevs := make(map[string][]map[string]interface{})
			virtDev["name"] = dev.Name
			virtDev["type"] = "vlan"
			virtDev["enabled"] = dev.Enabled
			virtDev["slave"] = dev.Slave
			virtDev["vlan_id"] = dev.Vlan_id
			virtDevs["devices"] = append(virtDevs["devices"], virtDev)
			s.postRequestToRPC(virtDevs)

		}
		config.Logging.Info("Create device if not exists " + dev.Name)
	}
	// update devices
	for _, dev := range devices {
		devItem := make(map[string]interface{})
		if dev.Type == "raw" {
			devItem["name"] = dev.Name
			devItem["enabled"] = dev.Enabled
		}
		if dev.Type == "vlan" {
			devItem["name"] = dev.Name
			devItem["enabled"] = dev.Enabled
			devItem["dst"] = dev.Dst
			devItem["forwarding"] = dev.Forwarding
		}
		s.putGroupDevRequest(devItem)
		config.Logging.Info("Update device state " + dev.Name)
	}

}

// Check exists device in database, if count device =0,
// then fill database from get request in rpc server
func (s *service) InitDevices() error {
	config.Logging.Info("Check device in database...")
	if s.storage.CheckExistDevice("") {
		config.Logging.Info("Devices already added to database")
		s.sendDevState()
		return nil
	}
	config.Logging.Info("Adding devices to database")
	url := "http://" + config.Params.RPCHost + ":" + strconv.Itoa(config.Params.RPCPort) + "/devices"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	deveces := new(devicesstorage.DtoDevices)
	err = json.Unmarshal(body, &deveces)
	if err != nil {
		return err
	}
	for _, dev := range deveces.Devices {
		newDev := devicesstorage.DtoCreateDevice{
			Name:         dev.Name,
			Type:         dev.Type,
			Enabled:      dev.Enabled,
			Running:      dev.Running,
			Flow_control: dev.Flow_control,
			Routing:      dev.Routing,
			Forwarding:   dev.Forwarding,
			Dst:          dev.Dst,
			Mtu:          dev.Mtu,
		}
		if dev.Type == "raw" {
			tx, err := s.storage.CreateDev(&newDev)
			if err != nil {
				return err
			}
			tx.Commit(context.Background())
		}
	}
	return nil
}

// get request for all device
func (s *service) GetAll(r *http.Request) (*http.Response, error) {
	config.Logging.LogRest("{}", devicesURL, r.Method)
	url := fmt.Sprintf("http://%s:%d%s", config.Params.RPCHost, config.Params.RPCPort, devicesURL)
	resp, err := http.Get(url)
	if err != nil {
		msg := "error GET request for all devices"
		config.Logging.Error(msg)
		return nil, err
	}
	return resp, nil

}

// get request for one device
func (s *service) GetOne(r *http.Request) (*http.Response, error) {

	var dev_name string
	vars := mux.Vars(r)
	if _, ok := vars["dev_name"]; ok {
		dev_name = vars["dev_name"]
	}
	url := fmt.Sprintf("http://%s:%d%s/%s", config.Params.RPCHost, config.Params.RPCPort, devicesURL, dev_name)
	config.Logging.LogRest("", devicesURL+"/"+dev_name, r.Method)
	resp, err := http.Get(url)
	if err != nil {
		msg := errors.New("rpc server response error")
		config.Logging.Error(msg.Error())
		return nil, msg
	}
	return resp, nil
}

func (s *service) GetDeviceStat(r *http.Request) (*http.Response, error) {
	var (
		dev_name string
		action   string
	)
	vars := mux.Vars(r)
	if _, ok := vars["dev_name"]; ok {
		dev_name = vars["dev_name"]
	}
	if _, ok := vars["action"]; ok {
		action = vars["action"]
	}
	if action != "stat" {
		msg := "wrong request"
		return nil, errors.New(msg)
	}
	url := fmt.Sprintf("http://%s:%d%s/%s/%s", config.Params.RPCHost, config.Params.RPCPort, devicesURL, dev_name, action)
	resp, err := http.Get(url)
	if err != nil {
		msg := "rpc server response error"
		config.Logging.Error(errors.New(msg).Error())
		return nil, err
	}
	return resp, err
}

// Service for delete device
func (s *service) DeleteDev(w http.ResponseWriter, r *http.Request) (*http.Response, error) {
	var dev_name string
	vars := mux.Vars(r)
	if _, ok := vars["dev_name"]; ok {
		dev_name = vars["dev_name"]
		config.Logging.LogRest("", devicesURL+"/"+dev_name, r.Method)
		resp, err := s.deleteDevByName(dev_name, r)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	resp, err := s.deleteGroupDev(r)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// delete one device
func (s *service) deleteDevByName(devName string, r *http.Request) (*http.Response, error) {
	url := fmt.Sprintf("http://%s:%d%s/%s", config.Params.RPCHost, config.Params.RPCPort, devicesURL, devName)
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	err = s.storage.DeleteDevByName(devName)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// delete array devices
func (s *service) deleteGroupDev(r *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	config.Logging.LogRest(string(body), devicesURL, r.Method)
	deveces := make(map[string][]map[string]interface{})
	err = json.Unmarshal(body, &deveces)
	if err != nil {
		return nil, err
	}
	for _, dev := range deveces["devices"] {
		resp, err = s.deleteDevByName(dev["name"].(string), r)
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}
