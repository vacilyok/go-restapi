package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGetMetrics(t *testing.T) {
	metrics_uri := fmt.Sprintf("%s/metrics", uri)
	resp, err := http.Get(metrics_uri)
	if err != nil {
		t.Error(err)
		return
	}

	if resp.StatusCode != 200 {
		t.Errorf("bad status code: %d", resp.StatusCode)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	metrics := make(map[string]int)
	err = json.Unmarshal(body, &metrics)
	if err != nil {
		t.Errorf("parsing metrics json response: %v", err)
		return
	}
}
