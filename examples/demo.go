package examples

import "time"

type Students struct {
	ID            int       `gorm:"column:id;AUTO_INCREMENT;primary_key" json:"id"`
	StudentNumber string    `gorm:"column:student_number;unique;NOT NULL" json:"student_number"`
	Name          string    `gorm:"column:name;NOT NULL" json:"name"`
	Age           int       `gorm:"column:age" json:"age"`
	CreatedAt     time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (s *Students) TableName() string {
	return "students"
}
