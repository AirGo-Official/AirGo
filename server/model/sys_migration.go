package model

type Migration struct {
	PanelType  string `json:"panel_type"  binding:"required"` //v2board, sspanel
	DBAddress  string `json:"db_address"  binding:"required"`
	DBPort     int64  `json:"db_port"     binding:"required"`
	DBUsername string `json:"db_username" binding:"required"`
	DBPassword string `json:"db_password" binding:"required"`
	DBName     string `json:"db_name"     binding:"required"`
}
