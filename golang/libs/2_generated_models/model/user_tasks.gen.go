// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameUserTask = "user_tasks"

// UserTask mapped from table <user_tasks>
type UserTask struct {
	ID          string     `gorm:"column:id;type:text;primaryKey" json:"id"`
	UserID      string     `gorm:"column:user_id;type:text;not null" json:"user_id"`
	TaskID      string     `gorm:"column:task_id;type:text;not null" json:"task_id"`
	CompletedAt *time.Time `gorm:"column:completed_at;type:datetime;default:CURRENT_TIMESTAMP" json:"completed_at"`
}

// TableName UserTask's table name
func (*UserTask) TableName() string {
	return TableNameUserTask
}
