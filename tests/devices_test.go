package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

var (
	dev_status bool
)

type device struct {
	Name         string   `json:"name"`
	Devtype      string   `json:"type"`
	Enabled      bool     `json:"enabled"`
	Raw_devices  []string `json:"raw_devices"`
	Running      bool     `json:"running"`
	Flow_control bool     `json:"flow_control"`
	Mtu          int      `json:"mtu"`
}

type simpleDevice struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

func TestGetDevices(t *testing.T) {

	devices_uri := fmt.Sprintf("%s/devices", uri)
	resp, err := http.Get(devices_uri)
	if err != nil {
		t.Error(err)
		return

	}
	if resp.StatusCode != 200 {
		t.Errorf("Get Devices bad status code: %d", resp.StatusCode)
		return
	}
	devices := make(map[string][]device)
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &devices)
	if err != nil {
		t.Errorf("parsing devices json response: %v", err)
		return
	}

}

func TestGetDevice(t *testing.T) {

	devices_uri := fmt.Sprintf("%s/devices/eth0", uri)
	resp, err := http.Get(devices_uri)
	if err != nil {
		t.Error(err)
		return
	}
	if resp.StatusCode != 200 {
		t.Errorf("Get Device bad status code: %d", resp.StatusCode)
		return
	}
	devices := new(device)
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &devices)
	if err != nil {
		t.Errorf("parsing device json response: %v", err)
		return
	}
	dev_status = devices.Enabled

}

func TestDeviceStatistics(t *testing.T) {
	devices_uri := fmt.Sprintf("%s/devices/eth0/stat", uri)
	resp, err := http.Get(devices_uri)
	if err != nil {
		t.Error(err)
		return
	}

	if resp.StatusCode != 200 {
		t.Errorf("Device Statistics bad status code: %d", resp.StatusCode)
		return
	}

}

func TestDevicePut(t *testing.T) {
	dev := simpleDevice{
		Name:    "eth0",
		Enabled: dev_status,
	}
	client := &http.Client{}
	devices_uri := fmt.Sprintf("%s/devices/eth0", uri)
	json, err := json.Marshal(dev)
	if err != nil {
		t.Error(err)
		return
	}
	req, err := http.NewRequest(http.MethodPut, devices_uri, bytes.NewBuffer(json))
	if err != nil {
		t.Error(err)
		return
	}
	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
		return
	}

	if resp.StatusCode != 200 {
		t.Errorf("Device PUT request bad status code: %d", resp.StatusCode)
		return
	}
}

func TestDevicePostAndDelete(t *testing.T) {

	// Start create device
	devices := devicesContainer()
	response, err := postRequest(devices["device1"], "devices")
	if err != nil {
		t.Error(err)
	}
	if response.StatusCode != 200 {
		t.Errorf("Device (device1) POST request bad status code: %d", response.StatusCode)
	}

	response, _ = postRequest(devices["device2"], "devices")
	if response.StatusCode == 200 || response.StatusCode == 201 {
		t.Errorf("Device (device2) POST request shouldn't work: %d", response.StatusCode)
	}

	//Start delete device
	url := fmt.Sprintf("%s/devices/vlan777_t", uri)
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Error(err)
	}
	response, err = client.Do(req)
	if err != nil {
		t.Error(err)
	}
	if response.StatusCode != 200 {
		t.Errorf("Delete device vlan777_t bad status code: %d", response.StatusCode)
	}

	// Delete group devices
	url = fmt.Sprintf("%s/devices", uri)
	client = &http.Client{}
	req, err = http.NewRequest("DELETE", url, bytes.NewBuffer([]byte(devices["delete"])))
	if err != nil {
		t.Error(err)
	}
	response, err = client.Do(req)
	if err != nil {
		t.Error(err)
	}
	if response.StatusCode != 200 {
		t.Errorf("Delete group devices bad status code: %d", response.StatusCode)
	}

}
