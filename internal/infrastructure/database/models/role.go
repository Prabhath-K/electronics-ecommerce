package models

type Role struct {
	ID   uint64 `gorm:"column:role_id;primaryKey;autoIncrement"`
	Name string `gorm:"column:role_name"`
}

func (Role) TableName() string {
	return "role"
}
