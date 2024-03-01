//角色-store
import {defineStore, storeToRefs} from "pinia";


import {ElMessage} from "element-plus";

import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)

export const useAdminRoleStore = defineStore("adminRoleStore", {
    state: () => ({
        // 角色菜单
        checkedMenuIDs:[] as number[],//编辑角色时，显示该角色已经关联的菜单
        allMenuList:[] as Route[],    //编辑角色时，显示全部菜单树
        roleList:{
            total:0,
            data:[] as RoleInfo[],
        },
        currentRole: {
            id: 0,
            role_name: '',
            description: '',
            status: true,
            menus: [] as Route[],
        } as RoleInfo,

        //角色casbin
        currentCasbin:{} as CasbinInfo,
        checkedCasbinPath:[] as string[],
        allCasbinInfo: {} as CasbinInfo,

    }),
    actions: {
        //获取角色列表
        async getRoleList() {
            const res: any = await request(apiStoreData.adminApi.value.getRoleList)
            this.roleList = res.data
        },
        //新建角色
        async newRole(checkedMenuIDs:number[],checkedCasbinPath:string[]){
            this.currentRole.id = 0 //清空上次编辑的id
            this.currentCasbin.roleID = 0 //清空上次编辑的id
            this.checkedMenuIDs = checkedMenuIDs
            this.checkedCasbinPath = checkedCasbinPath
            const role = this.menuIDsHandlerWhenSubmit(this.currentRole,this.checkedMenuIDs)
            const role1 = this.checkedCasbinPathHandlerWhenSubmit(role,this.checkedCasbinPath)
            return  request(apiStoreData.adminApi.value.newRole,role1)
        },
        //更新角色
        async updateRole(checkedMenuIDs:number[],checkedCasbinPath:string[]){
            this.checkedMenuIDs = checkedMenuIDs
            this.checkedCasbinPath = checkedCasbinPath
            const role = this.menuIDsHandlerWhenSubmit(this.currentRole,this.checkedMenuIDs)
            const role1 = this.checkedCasbinPathHandlerWhenSubmit(role,this.checkedCasbinPath)
            return request(apiStoreData.adminApi.value.updateRole,role1)
        },
        //删除角色
        async deleteRole(role:RoleInfo){
            return request(apiStoreData.adminApi.value.delRole,role)
        },
        //获取当前角色的casbin权限
        async getPolicyByID() {
            const res = await request(apiStoreData.adminApi.value.getPolicyByID, {roleID:this.currentCasbin.roleID} as CasbinInfo)
            this.currentCasbin = res.data
            this.checkedCasbinPathHandler()
        },
        //获取全部casbin权限
        async getAllPolicy() {
            const res = await request(apiStoreData.adminApi.value.getAllPolicy)
            this.allCasbinInfo = res.data
        },

        muneIDsHandler(){
            this.checkedMenuIDs = []
            let sonArr: number[] = []
            let sonNewArr: number[] = []
            let parentArr: number[] = []
            this.currentRole.menus.forEach((item: any) => {
                parentArr.push(item.parent_id)
                sonArr.push(item.id)
            });
            //父节点去重
            const newPsrentArr = parentArr.filter((value: any, index: any, array: any) => {
                return array.indexOf(value) === index;
            });
            //子节点去重
            sonArr.forEach((item: any) => {
                if (newPsrentArr.indexOf(item) === -1) {
                    sonNewArr.push(item)
                }
            })
            // console.log("newPsrentArr:", newPsrentArr)
            // console.log("sonArr:", sonArr)
            // console.log("sonNewArr:", sonNewArr)
            this.checkedMenuIDs = sonNewArr
        },
        menuIDsHandlerWhenSubmit(role:RoleInfo,checkedMenuIDs:number[]){
            role.menus = [] as Route[]
            checkedMenuIDs.forEach((value: number, index: number, array: number[])=>{
                role.menus.push({id:value} as Route)
            })
            return role
        },
        checkedCasbinPathHandler(){
            this.checkedCasbinPath=[]
            this.currentCasbin.casbinItems.forEach((value: CasbinItem, index: number, array: CasbinItem[])=>{
                this.checkedCasbinPath.push(value.path)
            })
        },
        checkedCasbinPathHandlerWhenSubmit(role:RoleInfo,checkedCasbinPath:string[]){
            role.casbins = []
            checkedCasbinPath.forEach((value: string, index: number, array: string[])=>{
                role.casbins.push({path:value} as CasbinItem)
            })
            return role
        },

    }
})