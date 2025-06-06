package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`    // 用户ID
	Username  string    `gorm:"size:32" json:"username"` // 用户名
	Password  string    `gorm:"size:128" json:"-"`       // 密码
	Email     string    `gorm:"size:128" json:"email"`   // 邮箱
	CreatedAt time.Time `json:"created_at"`              // 创建时间
	UpdatedAt time.Time `json:"updated_at"`              // 更新时间
	DeletedAt time.Time `gorm:"index" json:"-"`          // 删除时间
}
