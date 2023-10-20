// role
//used 角色信息
declare interface RowRoleType<T = any> {
    created_at: string;
    updated_at: string;
    id: number;    		//角色ID
    role_name: string;    	// 角色名称
    description: string;  	//角色描述
    status: boolean;      	//角色状态
    menus: [];				//角色菜单
    nodes: [];               //角色节点数组
    //sort: number;
    //createTime: string;
    [key: string]: T;
}

//casbin
declare interface CasbinInfo {
    roleID: number;
    casbinItems: [];
}

declare interface CasbinItem {
    method: string;
    path: string;
}

//used
interface SysRoleTableType extends TableType {
    data: RowRoleType[];
}

//used
declare interface SysRoleState {
    tableData: SysRoleTableType;
}
