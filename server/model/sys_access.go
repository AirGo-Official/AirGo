package model

import "time"

type Access struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ID        int64     `json:"id"      gorm:"primary_key"`
	Name      string    `json:"name"`
	Route     string    `json:"route" gorm:"size:3000;type:text"`
}
