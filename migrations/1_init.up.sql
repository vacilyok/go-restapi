CREATE TABLE IF NOT EXISTS device_type (
		id SERIAL PRIMARY KEY,
		type_name varchar(10)
	  );
CREATE TABLE IF NOT EXISTS device (
		id SERIAL PRIMARY KEY,
		name varchar(10),
		device_typeid integer NOT NULL,
		UNIQUE (name),
		FOREIGN KEY (device_typeid)  REFERENCES device_type (Id)
);
CREATE TABLE IF NOT EXISTS ip (
		id SERIAL PRIMARY KEY,
		address varchar(16) NOT NULL,
		mask integer NOT NULL,
		nexthop varchar(16)
);
CREATE TABLE IF NOT EXISTS RawDevice (
		id SERIAL PRIMARY KEY,
		device_id integer,
		running bool default true,
		enabled bool default true,
		routing bool default false,
		forwarding bool default false,
		flow_control bool default true,
		dst integer,
		mtu integer default 1500,
		ip_id integer,
		FOREIGN KEY (device_id)  REFERENCES device (Id),
		FOREIGN KEY (ip_id)  REFERENCES ip (id)
);
CREATE TABLE IF NOT EXISTS VlanDevice (
		vlan_device_id integer,
		slave integer,
		vlan_id integer,
		enabled bool default true,
		routing bool default false,
		forwarding bool default false,
		dst integer default NULL,
		FOREIGN KEY (vlan_device_id)  REFERENCES device (Id) ON DELETE CASCADE,
		FOREIGN KEY (slave)  REFERENCES device (Id)
) ;
INSERT INTO device_type (id,type_name) VALUES (1,'raw'),(2,'vlan'), (3,'lacp');	  
CREATE TABLE IF NOT EXISTS rules (
   id SERIAL PRIMARY KEY,
   body json DEFAULT NULL,
   active bool DEFAULT true
);