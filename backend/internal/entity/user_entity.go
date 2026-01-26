package entity

import "time"

type User struct {
	ID        uint      `gorm:"column:id;primaryKey;autoIncrement"`
	Username  string    `gorm:"column:username;type:varchar(100);uniqueIndex;not null"`
	Password  string    `gorm:"column:password;type:varchar(255);not null"`
	Role      string    `gorm:"column:role;type:enum('super_admin','admin','user');not null;default:'user'"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime;not null"`
}

func (u *User) TableName() string {
	return "users"
}
