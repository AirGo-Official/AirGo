<template>
  <div>
    <el-dialog :title="state.title" v-model="state.isShowDialog" width="80%">
      <el-descriptions
        style="margin-bottom: 20px"
        :column="1"
        border
      >
        <el-descriptions-item :label="$t('message.adminUser.SysUser.id') ">{{ userStoreData.currentUser.value.id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('message.adminUser.SysUser.created_at') ">{{ DateStrToTime(userStoreData.currentUser.value.created_at) }}</el-descriptions-item>
      </el-descriptions>
      <el-form ref="userDialogFormRef" :model="userStoreData.currentUser.value" size="default" label-position="top">
            <el-form-item :label="$t('message.adminUser.SysUser.user_name') ">
              <el-input v-model="userStoreData.currentUser.value.user_name"
                        clearable></el-input>
            </el-form-item>
            <el-form-item :label="$t('message.adminUser.SysUser.password') ">
              <el-input v-model="userStoreData.currentUser.value.password"
                        clearable></el-input>
            </el-form-item>
            <el-form-item :label="$t('message.adminUser.SysUser.role_group') ">
              <el-checkbox-group v-model="userStoreData.check_list.value">
                <el-checkbox :label="v.role_name" v-for="(v,index) in roleStoreData.roleList.value.data"
                             :key="index"></el-checkbox>
              </el-checkbox-group>
            </el-form-item>
            <el-form-item :label="$t('message.adminUser.SysUser.enable') ">
              <el-switch v-model="userStoreData.currentUser.value.enable" inline-prompt :active-text="$t('message.common.enable') "
                         inactive-text="$t('message.common.disable') "></el-switch>
            </el-form-item>
            <el-form-item :label="$t('message.adminUser.SysUser.balance') ">
              <el-input-number v-model.number="userStoreData.currentUser.value.balance" type="number"></el-input-number>
            </el-form-item>
        <el-form-item :label="$t('message.adminUser.SysUser.tg_id') ">
          <el-input v-model.number="userStoreData.currentUser.value.tg_id" type="number"></el-input>
        </el-form-item>
            <el-form-item :label="$t('message.adminUser.SysUser.referrer_code') ">
              <el-input v-model="userStoreData.currentUser.value.referrer_code"></el-input>
            </el-form-item>
            <el-form-item :label="$t('message.adminUser.SysUser.invitation_code')">
              <el-input v-model="userStoreData.currentUser.value.invitation_code">
                <template #append>
                  <el-button @click="resetInvitationCode">{{ $t('message.common.reset') }}</el-button>
                </template>
              </el-input>
            </el-form-item>
      </el-form>
      <template #footer>
				<span class="dialog-footer">
					<el-button @click="closeDialog" size="default">{{ $t('message.common.button_cancel') }}</el-button>
					<el-button type="primary" @click="onSubmit" size="default">{{ $t('message.common.button_confirm') }}</el-button>
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
import { useI18n } from "vue-i18n";

const userStore = useAdminUserStore()
const userStoreData = storeToRefs(userStore)
const roleStore = useAdminRoleStore()
const roleStoreData = storeToRefs(roleStore)
const shopStore = useAdminShopStore()
const emit = defineEmits(['refresh']);
const userDialogFormRef = ref();
const {t} = useI18n()

const state = reactive({
  isShowDialog: false,
  type: '',
  title: '',
  submitTxt: '',
})

// 打开弹窗
const openDialog = (type: string, row: SysUser) => {
  state.type = type
  state.isShowDialog = true;
  //打开时加载全部套餐
  shopStore.getGoodsList()
  //打开时加载全部角色，用来设置用户角色
  roleStore.getRoleList({page_num: 1, page_size: 10000})
  if (type === 'edit') {
    state.title = t('message.adminUser.modify_user');
    userStoreData.currentUser.value = row
    userStoreData.currentUser.value.password = ''
    //计算用户的角色
    let currentUserRoleIds: string[] = []
    userStoreData.currentUser.value.role_group.forEach((item: any) => {
      currentUserRoleIds.push(item.role_name)
    })
    userStoreData.check_list.value = currentUserRoleIds
  } else {
    state.title = t('message.adminUser.add_user');
    userStore.resetData();
  }
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
  if (state.type === 'add') {
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
