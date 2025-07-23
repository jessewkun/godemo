package provider

import (
	"context"
	"fmt"

	"github.com/jessewkun/gocommon/db/mysql"
	"gorm.io/gorm"
)

type MainDB struct{ *gorm.DB }

type MainDBName string

var MainDBNameValue MainDBName = "main"

// WithContext 为UserDB添加WithContext方法
func (db MainDB) WithContext(ctx context.Context) *gorm.DB {
	return db.DB.WithContext(ctx)
}

func ProvideMainDB(name MainDBName) MainDB {
	conn, err := mysql.GetConn(string(name))
	if err != nil {
		panic(fmt.Errorf("get db conn error: %s", err))
	}
	return MainDB{conn}
}
