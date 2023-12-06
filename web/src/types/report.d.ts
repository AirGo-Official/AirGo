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

declare interface FieldParams{
    table_name: string
    field_params_list: FieldParams[],
    pagination: Pagination,
}


//搜索条件
declare interface FieldParams {
    field: string;
    field_chinese_name: string;
    field_type: string;
    condition: string;
    condition_value: string;
    operator: string;
}


declare interface FieldTableNew {
    field: string;
    condition: string;
    condition_value: string;
}



//分页条件
declare interface Pagination {
    page_num: number
    page_size: number
    order_by: string,
}


declare interface PaginationParams {
    page_num: number;
    page_size: number;
    search: string;
    date: [];
}


