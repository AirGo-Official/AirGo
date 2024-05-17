<template>
  <div class="system-role-container layout-padding">
    <div class="system-role-padding layout-padding-auto layout-padding-view">
      <div class="system-user-search mb15">
        <el-button size="default" type="success" class="ml10" @click="onOpenAddRole('add')">
          <el-icon>
            <ele-FolderAdd/>
          </el-icon>
          {{$t('message.adminRole.add_role')}}
        </el-button>
      </div>
      <el-table :data="roleStoreData.roleList.value.data" stripe v-loading="state.loading"
                style="width: 100%">
        <el-table-column type="index" :label="$t('message.adminRole.RoleInfo.index')" width="60" fixed/>
        <el-table-column prop="role_name" :label="$t('message.adminRole.RoleInfo.role_name')" show-overflow-tooltip fixed></el-table-column>
        <el-table-column prop="id" :label="$t('message.adminRole.RoleInfo.id')" show-overflow-tooltip></el-table-column>
        <el-table-column prop="status" :label="$t('message.adminRole.RoleInfo.status')" show-overflow-tooltip>
          <template #default="scope">
            <el-tag type="success" v-if="scope.row.status">{{$t('message.common.enable')}}</el-tag>
            <el-tag type="info" v-else>{{$t('message.common.disable')}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('message.common.operate')" width="300">
          <template #default="scope">
            <el-button size="small" text type="primary"
                       @click="onOpenEditRole('edit', scope.row)">{{$t('message.common.modify')}}
            </el-button>
            <el-button size="small" text type="danger"
                       @click="deleteRole(scope.row)">{{$t('message.common.delete')}}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <RoleDialog ref="roleDialogRef" @refresh="getRoleList()"/>
  </div>
</template>

<script setup lang="ts" name="systemRole">
import {defineAsyncComponent, onMounted, reactive, ref} from 'vue';
import {ElMessage, ElMessageBox} from 'element-plus';

import {storeToRefs} from 'pinia';
import {useAdminRoleStore} from "/@/stores/admin_logic/roleStore";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
import { useI18n } from "vue-i18n";
const roleStore = useAdminRoleStore()
const roleStoreData = storeToRefs(roleStore)
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
const RoleDialog = defineAsyncComponent(() => import('/@/views/admin/role/dialog_editRole.vue'));
const roleDialogRef = ref();
const {t} = useI18n()


//定义参数
const state = reactive({
  loading: false,
})
// 打开新增角色弹窗
const onOpenAddRole = (type: string) => {
  roleDialogRef.value.openDialog(type);
};
// 打开修改角色弹窗
const onOpenEditRole = (type: string, row: Object) => {
  roleDialogRef.value.openDialog(type, row);
};

// 删除角色
const deleteRole = (row: RoleInfo) => {
  ElMessageBox.confirm(t('message.common.message_confirm_delete'), t('message.common.tip'), {
    confirmButtonText: t('message.common.button_confirm'),
    cancelButtonText: t('message.common.button_cancel'),
    type: 'warning',
  }).then(() => {
    roleStore.deleteRole({id: row.id} as RoleInfo).then(()=>{
      getRoleList()
    })
  })
};
//查询
const getRoleList = () => {
  roleStore.getRoleList()
}

// 页面加载时
onMounted(() => {
  getRoleList()
});
</script>

<style scoped lang="scss">
.system-role-container {
  .system-role-padding {
    padding: 15px;

    .el-table {
      flex: 1;
    }
  }
}
</style>
