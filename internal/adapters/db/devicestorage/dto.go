package devicesstorage

type DtoDevice struct {
	Id int `json:"id"`
	// RawId   int    `json:"raw_id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Enabled bool   `json:"enabled"`
}

type DtoUpdateDevice struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type DtoCreateDevice struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Enabled      bool   `json:"enabled"`
	Running      bool   `json:"running"`
	Flow_control bool   `json:"flow_control"`
	Routing      bool   `json:"routing"`
	Forwarding   bool   `json:"forwarding"`
	Mtu          int    `json:"mtu"`
	Dst          string `json:"dst"`
	// RawDevices   []string `json:"raw_devices"`
}

type DtoDevices struct {
	Devices []DtoCreateDevice
}

type PostDevice struct {
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	Enabled    bool    `json:"enabled"`
	Slave      string  `json:"slave"`
	Vlan_id    int     `json:"vlan_id"`
	Routing    bool    `json:"routing"`
	Forwarding bool    `json:"forwarding"`
	Dst        *string `json:"dst"`
}

type PostVlan struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Enabled bool   `json:"enabled"`
	Slave   string `json:"slave"`
	Vlan_id int    `json:"vlan_id"`
}
