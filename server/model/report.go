package model

type DbInfoReq struct {
	DbType       string   `json:"db_type"`
	DatabaseList []string `json:"database_list"`
}

type DbTableReq struct {
	DbName    string `json:"db_name"`
	TableName string `json:"table_name"`
}

// 查询参数 请求
type QueryParams struct {
	TableName       string            `json:"table_name" binding:"required"`
	FieldParamsList []FieldParamsItem `json:"field_params_list"`
	Pagination      Pagination        `json:"pagination"` //分页参数
}

// 查询字段参数
type FieldParamsItem struct {
	Operator       string `json:"operator"` // AND OR
	Field          string `json:"field"`
	FieldType      string `json:"field_type"`
	Condition      string `json:"condition"` //= > < <> like
	ConditionValue string `json:"condition_value"`
}

// 分页参数
type Pagination struct {
	PageNum  int64  `json:"page_num"`
	PageSize int64  `json:"page_size"`
	OrderBy  string `json:"order_by"`
}

type CommonDataResp struct {
	Total int64 `json:"total"`
	Data  any   `json:"data"`
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
