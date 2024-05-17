import {defineStore, storeToRefs} from 'pinia';
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
const apiStore = useApiStore()

export const useAdminUserStore = defineStore('adminUserStore', {
  state: () => ({
      userList: {
        total: 0,
        data: [] as SysUser[],
      },
        currentUser: {
          user_name: '',
          nick_name: '',
          password: '123456',
          avatar: '',
          enable: true,
          role_group: [] as RoleInfo[],
        } as SysUser,
        check_list: ['普通用户'], //选中的角色
    checkedRoleIDs:[] as number[],
    queryParams:{
      table_name: 'user',
      field_params_list: [
        {field:"id",field_chinese_name:"",field_type:"",condition:"=",condition_value:"",operator:""}
      ] as FieldParams[],
      pagination: { page_num: 1, page_size: 1, order_by: 'id ASC', } as Pagination,//分页参数
    } as QueryParams,
    userSummary:{
        lastMonth:[] as UserSummary[],
        thisMonth:[] as UserSummary[],
    }


  }),
  actions: {
    // 重置数据
    async resetData() {
      this.currentUser = {
        id:0,
        user_name: '',
        nick_name: '',
        password: '123456',
        avatar: '',
        enable: true,
        role_group: [] as RoleInfo[],
      } as SysUser
      this.check_list = ['普通用户']
    },
    //获取一位用户
    async firstUserByID(id:number){
      this.queryParams.field_params_list[0].condition_value = id.toString()
      const res = await  request(apiStore.adminApi.getUserList, this.queryParams)
      const list:SysUser[] = res.data.data
      if (list.length > 0){
        this.currentUser = list[0]
      }
    },

    //获取用户列表
    async getUserList(params?: QueryParams) {
      const res = await request(apiStore.adminApi.getUserList, params)
      this.userList = res.data
    },
    //新建用户
    async newUser(params?: SysUser) {
     return  request(apiStore.adminApi.newUser, params)
    },
    //修改用户
    async updateUser(params?: SysUser) {
      return  request(apiStore.adminApi.updateUser, params)
    },
    //删除用户
    async deleteUser(params?: SysUser) {
     return  request(apiStore.adminApi.deleteUser, params)
    },
    roleIDsHandler(user:SysUser){
      user.role_group.forEach((item)=>{
        this.checkedRoleIDs = []
        this.checkedRoleIDs.push(item.id)
      })
    },
    async getUserSummary(params:QueryParams,m:number){
      let mm = new Date().getMonth()
      const res = await request(apiStore.adminApi.userSummary,params)
      if (res.data === null){
        return
      }
      if (m === (mm+1)){ //本月=4，则mm=3;mm+1=当前月
        this.userSummary.thisMonth = res.data
      } else {
        this.userSummary.lastMonth = res.data
      }
    },
  },
});
