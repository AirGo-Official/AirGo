//角色-store
import {defineStore, storeToRefs} from "pinia";


import {ElMessage} from "element-plus";

import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";

const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)

//暴露role store
export const useRoleStore = defineStore("roleStore", {
    state: () => ({
        //角色管理页面参数
        roleManageData: {
            roles: {
                total: 0,
                role_list: [] as RowRoleType[],
            },
        },
        //编辑角色弹窗参数
        dialog: {
            roleForm: {       //当前编辑角色的参数
                id: 0,        //角色ID
                role_name: '',     // 角色名称
                description: '',  //角色描述
                status: true,     //角色状态
                menus: [],        //角色菜单
                nodes: [],
            } as RowRoleType,
        },
        //编辑角色权限弹窗参数
        dialogEditApi: {
            casbinInfo: {    //当前角色的权限
                roleID: 0,
                casbinItems: [''],
            },
            allCasbinInfo: {    //全部权限
                roleID: 0,
                casbinItems: [],
            } as CasbinInfo,
        },

    }),
    actions: {
        //获取角色列表
        async getRoleList(params?: object) {
            const res: any = await request(apiStoreData.api.value.role_getRoleList, params)
            this.roleManageData.roles = res.data
            ElMessage.success(res.msg)
        },

        //获取当前角色的权限
        async getPolicyByRoleIds(params?: object) {
            const res = await request(apiStoreData.api.value.casbin_getPolicyByRoleIds, params)
            var casbinRes: CasbinInfo = res.data
            if (casbinRes.casbinItems !== null) {
                var oldArr: string[] = []
                casbinRes.casbinItems.forEach((item: CasbinItem) => {
                    oldArr.push(item.path)
                });
                this.dialogEditApi.casbinInfo.casbinItems = oldArr
            }
        },
        //获取全部权限
        async getAllPolicy() {
            const res = await request(apiStoreData.api.value.casbin_getAllPolicy)
            this.dialogEditApi.allCasbinInfo = res.data
            ElMessage.success(res.msg)
        },
        //更新角色权限
        async updateCasbinPolicy(params?: object) {
            const res = await request(apiStoreData.api.value.casbin_updateCasbinPolicy, params)
            ElMessage.success(res.msg)
        }
    }
})