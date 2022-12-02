package devices

import (
	"context"
	"errors"
	"fmt"
	devicesstorage "mediator/internal/adapters/db/devicestorage"
	"mediator/internal/config"
	"mediator/pkg/database/pg"

	"github.com/jackc/pgx/v4"
)

type dbs struct {
	conn *pgx.Conn
}

func NewDevStorage(dbConn *pgx.Conn) devicesstorage.DevStorage {
	return &dbs{
		conn: dbConn,
	}
}

// returns a dictionary to get the device type NAME by type ID
func deviceTypeById() map[int]string {
	devType := make(map[int]string)
	devType[1] = "raw"
	devType[2] = "vlan"
	devType[3] = "lacp"
	return devType
}

// returns a dictionary to get the device type ID by type NAME
func deviceTypeByName() map[string]int {
	devType := make(map[string]int)
	devType["raw"] = 1
	devType["vlan"] = 2
	devType["lacp"] = 3
	return devType
}

// Get device from db by device name
func (d *dbs) GetDevByName(devName string) *devicesstorage.DtoDevice {
	var (
		devTypeid int
	)
	ctx := context.Background()
	dbDev := new(devicesstorage.DtoDevice)
	devMapID := deviceTypeById()
	query := "SELECT id, name, enabled, device_typeid FROM device d WHERE name=$1"
	result, _ := d.conn.Query(ctx, query, devName)

	if result.Next() {
		result.Scan(&dbDev.Id, &dbDev.Name, &dbDev.Enabled, &devTypeid)
	}
	dbDev.Type = devMapID[devTypeid]
	result.Close()
	return dbDev
}

// Get all device from db
func (d *dbs) GetAllDev() []*devicesstorage.PostDevice {
	var (
		devTypeid, vlan_id  int
		dbDevices           []*devicesstorage.PostDevice
		name, slave         string
		enabled, forwarding bool
		dst                 *string
	)
	dst = nil
	dbDev := new(devicesstorage.PostDevice)
	devMapID := deviceTypeById()
	query := "SELECT d.name, r.enabled,d.device_typeid,0 as vlan_id,r.forwarding,'' as slave, '' as dst " +
		" FROM device d  inner join RawDevice r on r.device_id =d.id " +
		" UNION " +
		" SELECT d.name, v.enabled,d.device_typeid,v.vlan_id, v.forwarding,  " +
		" (select distinct name from device where id =v.slave)," +
		" COALESCE((select distinct name from device where id =v.dst),'')" +
		" FROM device d inner join VlanDevice v on v.vlan_device_id =d.id "

	ctx := context.Background()
	rows, err := d.conn.Query(ctx, query)
	if err != nil {
		config.Logging.Error(err.Error())
	}

	for rows.Next() {
		err = rows.Scan(&name, &enabled, &devTypeid, &vlan_id, &forwarding, &slave, &dst)
		if err != nil {
			config.Logging.Error(err.Error())
		}
		dbDev.Name = name
		dbDev.Enabled = enabled
		dbDev.Type = devMapID[devTypeid]
		dbDev.Vlan_id = vlan_id
		dbDev.Forwarding = forwarding
		dbDev.Slave = slave
		dbDev.Dst = dst
		dbDevices = append(dbDevices, dbDev)
		dbDev = new(devicesstorage.PostDevice)
	}
	rows.Close()
	return dbDevices

}

// Check exists device in database
func (d *dbs) getDevIdByName(devName string) (int, int) {
	var dev_id, device_typeid int
	query := "SELECT id, device_typeid FROM device WHERE name=$1"
	ctx := context.Background()
	result, _ := d.conn.Query(ctx, query, devName)
	if result.Next() {
		result.Scan(&dev_id, &device_typeid)
	}
	result.Close()
	return dev_id, device_typeid
}

// Check exists device in database
func (d *dbs) CheckExistDevice(devName string) bool {

	var err error
	if !pg.IsConnected(d.conn) {
		d.conn, err = pg.ConnectDB()
	}
	if err != nil {
		config.Logging.Error("Fail check exist device. No connection to database")
		return false
	}

	countRows := 0
	query := "SELECT count(*) cnt FROM device"
	if devName != "" {
		query = fmt.Sprintf("SELECT count(*) cnt FROM device WHERE name='%s' ", devName)
	}
	ctx := context.Background()
	result, err := d.conn.Query(ctx, query)
	if err != nil {
		return false
	}
	if result.Next() {
		result.Scan(&countRows)
	}
	if countRows > 0 {
		result.Close()
		return true
	}
	result.Close()
	return false
}

// Update device state (enabled: true/false) by device name
func (d *dbs) UpdateDbDev(dev map[string]interface{}) error {
	var (
		dev_id, dev_type int
		query            string
		err              error
	)
	if !pg.IsConnected(d.conn) {
		d.conn, err = pg.ConnectDB()
	}
	if err != nil {
		return errors.New("fail update device, no connection to database")
	}

	if dev_name, ok := dev["name"]; ok {
		dev_id, dev_type = d.getDevIdByName(dev_name.(string))
	}
	switch dev_type {
	case 1:
		if _, ok := dev["enabled"]; ok {
			query = fmt.Sprintf("UPDATE RawDevice set enabled=%t WHERE device_id=%d;", dev["enabled"].(bool), dev_id)
		}
		if _, ok := dev["mtu"]; ok {
			query += fmt.Sprintf("UPDATE RawDevice set mtu=%d WHERE device_id=%d;", dev["mtu"].(int), dev_id)
		}
		_, err := d.conn.Exec(context.Background(), query)
		if err != nil {
			config.Logging.Error(err.Error())
			return err
		}

	case 2:
		if _, ok := dev["enabled"]; ok {
			query = fmt.Sprintf("UPDATE VlanDevice set enabled=%t WHERE vlan_device_id=%d;", dev["enabled"].(bool), dev_id)
		}
		if _, ok := dev["forwarding"]; ok {
			query += fmt.Sprintf("UPDATE VlanDevice set forwarding=%t WHERE vlan_device_id=%d;", dev["forwarding"].(bool), dev_id)
		}

		if dst, ok := dev["dst"]; ok {
			if dst != nil {
				dst_id, _ := d.getDevIdByName(dev["dst"].(string))
				query += fmt.Sprintf("UPDATE VlanDevice set dst=%d WHERE vlan_device_id=%d;", dst_id, dev_id)
			} else {
				query += fmt.Sprintf("UPDATE VlanDevice set dst=NULL WHERE vlan_device_id=%d;", dev_id)
			}
		}
		_, err := d.conn.Exec(context.Background(), query)
		if err != nil {
			config.Logging.Error(err.Error())
			return err
		}
	}
	return nil
}

func (d *dbs) InsertVlan(device *devicesstorage.PostDevice) error {
	var (
		err       error
		query     string
		lastDevId int
	)
	if !pg.IsConnected(d.conn) {
		d.conn, err = pg.ConnectDB()
	}

	if err != nil {
		err := errors.New("Fail create new device " + device.Name + " . No connection to database")
		config.Logging.Error(err.Error())
		return err
	}
	devType := deviceTypeByName()
	query = "INSERT INTO device (name, device_typeid) VALUES ($1, $2) RETURNING id"
	err = d.conn.QueryRow(context.Background(), query, device.Name, devType[device.Type]).Scan(&lastDevId)
	if err != nil {
		config.Logging.Error("Fail create vlan " + device.Name)
		return err
	}
	slave_id, _ := d.getDevIdByName(device.Slave)

	if device.Dst != nil {
		dst_id, _ := d.getDevIdByName(*device.Dst)
		query = "INSERT INTO VlanDevice (vlan_device_id, slave, vlan_id,enabled,forwarding, dst) VALUES ($1, $2,$3,$4,$5,$6)"
		_, err = d.conn.Exec(context.Background(), query, lastDevId, slave_id, device.Vlan_id, device.Enabled, device.Forwarding, dst_id)
	} else {
		query = "INSERT INTO VlanDevice (vlan_device_id, slave, vlan_id,enabled,forwarding) VALUES ($1, $2,$3,$4,$5)"
		_, err = d.conn.Exec(context.Background(), query, lastDevId, slave_id, device.Vlan_id, device.Enabled, device.Forwarding)
	}

	if err != nil {
		config.Logging.Error("Fail create vlan " + device.Name)
		return err
	}
	return nil
}

// create new device in database
func (d *dbs) CreateDev(newDev *devicesstorage.DtoCreateDevice) (pgx.Tx, error) {
	var (
		err       error
		tx        pgx.Tx
		query     string
		lastDevId int
	)

	if !pg.IsConnected(d.conn) {
		d.conn, err = pg.ConnectDB()
	}

	if err != nil {
		err := errors.New(" Fail create new device. No connection to database")
		return nil, err
	}

	devType := deviceTypeByName()
	ctx := context.Background()

	numTypeNewDev := devType[newDev.Type]
	query = "INSERT INTO device (name, device_typeid) VALUES ($1, $2) RETURNING id"
	tx, err = d.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		config.Logging.Error("Create device BeginTx error")
		return nil, err
	}
	err = tx.QueryRow(ctx, query, newDev.Name, devType[newDev.Type]).Scan(&lastDevId)
	if err != nil {
		config.Logging.Error("Error insert device with name " + newDev.Name)
		return nil, err
	}
	if numTypeNewDev == 1 {
		query = "INSERT INTO RawDevice (device_id, running, enabled, routing, forwarding, flow_control, mtu)  VALUES ($1,$2,$3,$4,$5,$6,$7)"
		_, err = tx.Exec(ctx, query, lastDevId, newDev.Running, newDev.Enabled, newDev.Routing, newDev.Forwarding, newDev.Flow_control, newDev.Mtu)
		if err != nil {
			tx.Rollback(ctx)
			config.Logging.Error("fail insert new raw device " + newDev.Name)
			return nil, err
		}
	}
	return tx, err
}

// delete device by name from db
func (d *dbs) DeleteDevByName(devName string) error {
	var err error
	if !pg.IsConnected(d.conn) {
		d.conn, err = pg.ConnectDB()
	}

	if err != nil {
		return errors.New(" Query delete device from db is fail. Not connect to db")
	}
	ctx := context.Background()
	tx, err := d.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, "Delete from  device WHERE name = $1", devName)
	if err != nil {
		return err
	}
	tx.Commit(ctx)
	return nil
}
