package provider

import (
	"github.com/jessewkun/gocommon/db"
	"gorm.io/gorm"
)

type MainDB struct{ *gorm.DB }

type MainDBName string

var MainDBNameValue MainDBName = "main"

func ProvideMainDB(name MainDBName) MainDB {
	return MainDB{db.GetConn(string(name))}
}
