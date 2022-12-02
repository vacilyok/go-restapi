package devicesstorage

import "database/sql"

type DevStorage interface {
	UpdateDbDev(map[string]interface{}) error
	DeleteDevByName(devName string) error
	CheckExistDevice(string) bool
	GetDevByName(string) *DtoDevice
	CreateDev(*DtoCreateDevice) (*sql.Tx, error)
	InsertVlan(*PostDevice) error
	GetAllDev() []*PostDevice
}
