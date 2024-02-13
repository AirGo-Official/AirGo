package model

import "time"

type Access struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
	ID        int64      `json:"id"      gorm:"primaryKey"`
	Name      string     `json:"name"`
	Route     string     `json:"route" gorm:"type:text"`
}
