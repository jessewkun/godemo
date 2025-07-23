package model

import (
	"time"

	"github.com/jessewkun/gocommon/db/mysql"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	mysql.BaseModel
	Username  string         `gorm:"size:32" json:"username"`         // 用户名
	Password  string         `gorm:"size:128" json:"-"`               // 密码
	Email     string         `gorm:"size:128" json:"email"`           // 邮箱
	DeletedAt mysql.DateTime `gorm:"type:datetime" json:"deleted_at"` // 删除时间
}

func (m *User) BeforeCreate(tx *gorm.DB) (err error) {
	now := mysql.DateTime(time.Now())
	m.CreatedAt = now
	m.ModifiedAt = now
	return err
}

func (m *User) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = mysql.DateTime(time.Now())
	return nil
}
