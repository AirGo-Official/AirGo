declare interface RoleInfo
{
    created_at: string;
    updated_at: string;
    id: number;    		//角色ID
    role_name: string;   // 角色名称
    description: string; //角色描述
    status: boolean //角色状态
    user_group: [] //用户
    menus: Route[];		//角色菜单
    casbins:CasbinItem[]
}

//casbin
declare interface CasbinInfo {
    roleID: number;
    casbinItems:CasbinItem[];
}

declare interface CasbinItem {
    method: string;
    path: string;
}

//used
interface SysRoleTableType extends TableType {
    data: RoleInfo[];
}

//used
declare interface SysRoleState {
    tableData: SysRoleTableType;
}
