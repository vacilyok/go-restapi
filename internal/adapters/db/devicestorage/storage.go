package devicesstorage

import "github.com/jackc/pgx/v4"

type DevStorage interface {
	UpdateDbDev(map[string]interface{}) error
	DeleteDevByName(devName string) error
	CheckExistDevice(string) bool
	GetDevByName(string) *DtoDevice
	CreateDev(*DtoCreateDevice) (pgx.Tx, error)
	InsertVlan(*PostDevice) error
	GetAllDev() []*PostDevice
}
