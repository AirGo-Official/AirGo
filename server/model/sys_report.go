package model

type DbNameAndTableReq struct {
	Database  string `json:"database" gorm:"column:database"`
	TableName string `json:"table_name" gorm:"column:table_name"`
}

type DbInfoReq struct {
	DbType       string   `json:"db_type"`
	DatabaseList []string `json:"database_list"`
}

// 查询参数 请求
type FieldParamsReq struct {
	TableName        string            `json:"table_name" binding:"required"`
	FieldParamsList  []FieldParamsItem `json:"field_params_list"`
	PaginationParams PaginationParams  `json:"pagination_params"` //分页参数
}
type FieldParamsItem struct {
	Operator       string `json:"operator"` // AND OR
	Field          string `json:"field"`
	FieldType      string `json:"field_type"`
	Condition      string `json:"condition"` //= > < <> like
	ConditionValue string `json:"condition_value"`
}

// mysql 表名
type DbMysqlTable struct {
	TableName string `json:"table_name" gorm:"column:table_name"`
}

// mysql 字段
type DbMysqlColumn struct {
	DataType      string `json:"data_type" gorm:"column:data_type"`
	ColumnName    string `json:"column_name" gorm:"column:column_name"`
	DataTypeLong  string `json:"data_type_long" gorm:"column:data_type_long"`
	ColumnComment string `json:"column_comment" gorm:"column:column_comment"`
}

// sqlite 表名
type DbSqliteTable struct {
	TableName string `json:"name" gorm:"column:name"`
}

// sqlite 字段
type DbSqliteColumn struct {
	Name     string `json:"name" gorm:"column:name"`
	NameType string `json:"type" gorm:"column:type"`
}
