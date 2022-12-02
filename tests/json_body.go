package test

func rulesContainer() map[string]string {
	rules := make(map[string]string)
	rules["rule1"] = `{
		"rules": []
	}`

	rules["rule2"] = `{
		"rules": [
			{
				"prefix": "10.10.1.0/24",
				"countermeasures": []
			}
		]
	}`

	rules["rule3"] = `{
		"rules": [
		  {
			"prefix": "12.12.12.0/24",
			"countermeasures": [
			  {
				"uuid": "11110000-0000-0000-test-000000000011",
				"matches": [
				  {
					"name": "l3_protocol",
					"options": {
					  "protocol": "tcp"
					}
				  }
				],
				"action": {
				  "name": "limit",
					"options": {
					  "pps": 10000,
					  "burst":1000
					}
				}
			  }
			]
		  }]`

	rules["rule4"] = `{
		"rules": [
			{
			  "prefix": "192.168.13.13/32",
			  "countermeasures": [
				{
				  "uuid": "e7b50000-0000-0000-e9b5-000000000000",
				  "matches": [
					{
					  "name": "l3_protocol",
					  "options": {
						"protocol": "tcp"
					  }
					}
				  ],
				  "action": {
					"name": "limit",
					  "options": {
						"pps": 1000,
						"burst":1000
					  }
				  }
				}
			  ]
			}		
			]
		}`

	ss := `{
		"rules": [
		  {
			"prefix": "8.8.1.0/24",
			"countermeasures": [
			  {
				"uuid": "11110000-0000-0000-e9b5-000000000011",
				"matches": [
				  {
					"name": "l3_protocol",
					"options": {
					  "protocol": "tcp"
					}
				  }
				],
				"action": {
				  "name": "limit",
					"options": {
					  "pps": 10000,
					  "burst":1000
					}
				}
			  }
			]
		  },
		  {
			"prefix": "22.22.22.0/24",
			"countermeasures": [
			  {
				"uuid": "22220000-0000-0000-e9b5-000000002222",
				"matches": [
				  {
					"name": "l3_protocol",
					"options": {
					  "protocol": "udp"
					}
				  }
				],
				"action": {
				  "name": "drop",
					"options": {}
				}
			  },
			  {
				"uuid": "32220000-2222-0000-e9b5-000000002222",
				"matches": [
				  {
					"name": "l3_protocol",
					"options": {
					  "protocol": "tcp"
					}
				  }
				],
				"action": {
				  "name": "accept",
					"options": {}
				}
			  }        
			]
		  }    
		]
	  }`

	rules["rule13"] = ss
	return rules
}

func devicesContainer() map[string]string {
	devices := make(map[string]string)

	devices["device1"] = `{
		"devices":[
		{
		"name":"vlan777_t",
		"enabled":true,
		"forwarding":true,
		"slave":"eth0",
		"type":"vlan",
		"vlan_id":777
		},
	  {
		"name":"vlan888_t",
		"enabled":true,
		"forwarding":true,
		"slave":"eth0",
		"type":"vlan",
		"vlan_id":888
		},
		{
		  "name":"vlan889_t",
		  "enabled":true,
		  "forwarding":true,
		  "slave":"eth0",
		  "type":"vlan",
		  "vlan_id":889
		  }  
		]
	  }`

	devices["device2"] = `{
		"devices":[
		{
		"name":"vlan666_t",
		"enabled":true,
		"forwarding":true,
		"slave":"eth15",
		"type":"vlan",
		"vlan_id":666
		}  
		]
	  }`

	devices["delete"] = `{
		"devices": [
		{
		  "name": "vlan888_t"
		},
		{
		  "name": "vlan889_t"
		}
	  ]
	}`
	return devices
}
