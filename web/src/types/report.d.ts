declare interface MysqlTable {
    table_name: string;
}

declare interface SqliteTable {
    name: string;
}

declare interface MysqlColumn {
    data_type: string;
    column_name: string;
    data_type_long: string;
    column_comment: string;
}

declare interface SqliteColumn {
    name: string;
    type: string;
}


//搜索条件
declare interface FieldTable {
    field: string;
    field_chinese_name: string;
    field_type: string;
    condition: string;
    condition_value: string;
}
declare interface FieldTableNew {
    field: string;
    condition: string;
    condition_value: string;
}
//分页条件
declare interface Pagination {
    search: string
    page_num: number
    page_size: number
    date: [],
}


