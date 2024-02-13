<template>
  <div>
    <el-dialog :title="state.title" v-model="state.isShowDialog" width="80%">
      <el-descriptions
        style="margin-bottom: 20px"
        :column="1"
        border
      >
        <el-descriptions-item label="ID">{{ userStoreData.currentUser.value.id }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ DateStrToTime(userStoreData.currentUser.value.created_at) }}</el-descriptions-item>
      </el-descriptions>
      <el-form ref="userDialogFormRef" :model="userStoreData.currentUser.value" size="default" label-width="90px">
            <el-form-item label="账户邮箱">
              <el-input v-model="userStoreData.currentUser.value.user_name" placeholder="请输入账户邮箱"
                        clearable></el-input>
            </el-form-item>
            <el-form-item label="账户密码">
              <el-input v-model="userStoreData.currentUser.value.password" placeholder="请输入"
                        clearable></el-input>
            </el-form-item>
            <el-form-item label="关联角色">
              <el-checkbox-group v-model="userStoreData.check_list.value">
                <el-checkbox :label="v.role_name" v-for="(v,index) in roleStoreData.roleList.value.data"
                             :key="index"></el-checkbox>
              </el-checkbox-group>
            </el-form-item>
            <el-form-item label="用户状态">
              <el-switch v-model="userStoreData.currentUser.value.enable" inline-prompt active-text="开启"
                         inactive-text="关闭"></el-switch>
            </el-form-item>
            <el-form-item label="余额">
              <el-input-number v-model.number="userStoreData.currentUser.value.balance" type="number"></el-input-number>
            </el-form-item>
        <el-form-item label="TG ID">
          <el-input v-model.number="userStoreData.currentUser.value.tg_id" type="number"></el-input>
        </el-form-item>
            <el-form-item label="推荐人">
              <el-input v-model="userStoreData.currentUser.value.referrer_code"></el-input>
            </el-form-item>
            <el-form-item label="邀请码">
              <el-input v-model="userStoreData.currentUser.value.invitation_code">
                <template #append>
                  <el-button @click="resetInvitationCode">重置</el-button>
                </template>
              </el-input>
            </el-form-item>
      </el-form>
      <template #footer>
				<span class="dialog-footer">
					<el-button @click="closeDialog" size="default">取 消</el-button>
					<el-button type="primary" @click="onSubmit" size="default">{{ state.submitTxt }}</el-button>
				</span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import {reactive, ref} from 'vue';

import {storeToRefs} from 'pinia';
import {useAdminRoleStore} from '/@/stores/admin_logic/roleStore'
import {randomStringNew} from "/@/utils/encrypt"
import { useAdminShopStore } from "/@/stores/admin_logic/shopStore";
import { useAdminUserStore } from "/@/stores/admin_logic/userStore";
import {DateStrToTime} from "/@/utils/formatTime"

const userStore = useAdminUserStore()
const userStoreData = storeToRefs(userStore)
const roleStore = useAdminRoleStore()
const roleStoreData = storeToRefs(roleStore)
const shopStore = useAdminShopStore()
const emit = defineEmits(['refresh']);
const userDialogFormRef = ref();

const state = reactive({
  isShowDialog: false,
  type: '',
  title: '',
  submitTxt: '',
})

// 打开弹窗
const openDialog = (type: string, row: SysUser) => {
  state.isShowDialog = true;
  shopStore.getGoodsList() //获取全部套餐
  if (type === 'edit') {
    state.title = '修改用户';
    state.submitTxt = '修 改';
    // userStoreData.currentUser.value = JSON.parse(JSON.stringify(row)) //深拷贝,防止修改时间报错
    userStoreData.currentUser.value = row
    userStoreData.currentUser.value.password = ''
    //计算用户的角色
    let currentUserRoleIds: string[] = []
    userStoreData.currentUser.value.role_group.forEach((item: any) => {
      currentUserRoleIds.push(item.role_name)
    })
    userStoreData.check_list.value = currentUserRoleIds
  } else {
    state.title = '新增用户';
    state.submitTxt = '新 增';
    userStore.resetData();
  }
//打开时加载全部角色，用来设置用户角色
  roleStore.getRoleList({page_num: 1, page_size: 10000})
};
// 关闭弹窗
const closeDialog = () => {
  state.isShowDialog = false;
};
// 提交
const onSubmit = () => {
  //处理角色
  userStoreData.currentUser.value.role_group = [];
  userStoreData.check_list.value.forEach((value: string, index: number, array: string[]) => {
    userStoreData.currentUser.value.role_group.push({
      role_name: value,
    } as RoleInfo)
  })
  if (state.title === '新增用户') {
    userStore.newUser(userStoreData.currentUser.value)
  } else {
    userStore.updateUser(userStoreData.currentUser.value)
  }
  setTimeout(() => {
    emit('refresh')
  }, 500)
  closeDialog();
};

//重置邀请码
const resetInvitationCode=()=>{
  userStoreData.currentUser.value.invitation_code =  randomStringNew(8)
}
// 暴露变量
defineExpose({
  openDialog,
});
</script>
