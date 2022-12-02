	  CREATE TABLE IF NOT EXISTS device_type (
		id int(10) unsigned NOT NULL AUTO_INCREMENT,
		type_name varchar(10),
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
	
	CREATE TABLE IF NOT EXISTS device (
		id int(10) unsigned NOT NULL AUTO_INCREMENT,
		name varchar(10),
		device_typeid int(10) unsigned NOT NULL,
		PRIMARY KEY (id),
		UNIQUE (name),
		FOREIGN KEY (device_typeid)  REFERENCES device_type (Id)
	  ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
	

	CREATE TABLE IF NOT EXISTS ip (
		id int(10) unsigned NOT NULL AUTO_INCREMENT,
		address varchar(16) NOT NULL,
		mask int(2) NOT NULL,
		nexthop varchar(16),
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
	 
	 CREATE TABLE IF NOT EXISTS RawDevice (
		id int(10) unsigned NOT NULL AUTO_INCREMENT,
		device_id int(10) unsigned,
		running bool default true,
		enabled bool default true,
		routing bool default false,
		forwarding bool default false,
		flow_control bool default true,
		dst int(10) unsigned,
		mtu int(10) unsigned default 1500,
		ip_id int(10) unsigned,
		PRIMARY KEY (id),
		FOREIGN KEY (device_id)  REFERENCES device (Id) ON DELETE CASCADE,
		FOREIGN KEY (ip_id)  REFERENCES ip (id)
	  ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

		CREATE TABLE IF NOT EXISTS VlanDevice (
		vlan_device_id int(10) unsigned,
		slave int(10) unsigned,
		vlan_id int(10) unsigned,
		enabled bool default true,
		routing bool default false,
		forwarding bool default false,
		dst int(10) unsigned default NULL,
		FOREIGN KEY (vlan_device_id)  REFERENCES device (Id) ON DELETE CASCADE,
		FOREIGN KEY (slave)  REFERENCES device (Id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

	CREATE TABLE IF NOT EXISTS logger (
	 	id int unsigned NOT NULL AUTO_INCREMENT,
	 	date datetime NOT NULL,
	 	source varchar(50) NOT NULL,
	 	message text NOT NULL,
	 	level tinyint NOT NULL,
	 	PRIMARY KEY (id)
	 	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

	CREATE TABLE IF NOT EXISTS rest_logger (
	 		id int unsigned NOT NULL AUTO_INCREMENT,
	 		date datetime NOT NULL,
	 		source varchar(50) NOT NULL,
	 		route varchar(50),
	 		method varchar(6),
	 		message json,
	 		PRIMARY KEY (id)
	 	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;		
	
	-- CREATE TABLE IF NOT EXISTS LacpDevice (
	-- 	lacp_device_id int(10) unsigned,
	-- 	device_id int(10) unsigned,
	-- 	FOREIGN KEY (lacp_device_id)  REFERENCES device (Id),
	-- 	FOREIGN KEY (device_id)  REFERENCES device (Id)
	--   ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
	  