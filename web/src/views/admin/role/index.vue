<template>
  <div class="system-role-container layout-padding">
    <div class="system-role-padding layout-padding-auto layout-padding-view">
      <div class="system-user-search mb15">
        <el-input v-model="state.queryParams.search" size="default" placeholder="请输入角色名称"
                  style="max-width: 180px"></el-input>
        <el-button size="default" type="success" class="ml10" @click="onOpenAddRole('add')">
          <el-icon>
            <ele-FolderAdd/>
          </el-icon>
          新增角色
        </el-button>
      </div>
      <el-table :data="roleStoreData.roleList.value.data" stripe v-loading="state.loading"
                style="width: 100%">
        <el-table-column type="index" label="序号" width="60" fixed/>
        <el-table-column prop="role_name" label="角色名称" show-overflow-tooltip fixed></el-table-column>
        <el-table-column prop="id" label="角色ID" show-overflow-tooltip></el-table-column>
        <el-table-column prop="status" label="角色状态" show-overflow-tooltip>
          <template #default="scope">
            <el-tag type="success" v-if="scope.row.status">启用</el-tag>
            <el-tag type="info" v-else>禁用</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="300">
          <template #default="scope">
            <el-button size="small" text type="primary"
                       @click="onOpenEditRole('edit', scope.row)">修改
            </el-button>
            <el-button size="small" text type="primary"
                       @click="onRowDel(scope.row)">删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
          background
          :page-sizes="[10, 30, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="onHandleSizeChange" @current-change="onHandleCurrentChange" class="mt15"
          v-model:current-page="state.queryParams.page_num"
          v-model:page-size="state.queryParams.page_size"
          :total="roleStoreData.roleList.value.total">
      </el-pagination>
    </div>
    <RoleDialog ref="roleDialogRef" @refresh="roleStore.getRoleList(state.queryParams)"/>
  </div>
</template>

<script setup lang="ts" name="systemRole">
import {defineAsyncComponent, onMounted, reactive, ref} from 'vue';
import {ElMessage, ElMessageBox} from 'element-plus';

import {storeToRefs} from 'pinia';
import {useAdminRoleStore} from "/@/stores/admin_logic/roleStore";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
const roleStore = useAdminRoleStore()
const roleStoreData = storeToRefs(roleStore)
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
const RoleDialog = defineAsyncComponent(() => import('/@/views/admin/role/dialog_editRole.vue'));
const roleDialogRef = ref();


//定义参数
const state = reactive({
  loading: false,
  queryParams: {
    search: '',
    page_num: 1,
    page_size: 30,
  },
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
const onRowDel = (row: RoleInfo) => {
  ElMessageBox.confirm(`此操作将永久删除角色名称：“${row.role_name}”，是否继续?`, '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(() => {
    request(apiStoreData.adminApi.value.delRole, {id: row.id}).then((res) => {
      ElMessage.success('删除失败');
    }).catch(() => {
      onSearch(state.queryParams)
    })
  })
};
//查询
const onSearch = (params?: object) => {
  roleStore.getRoleList(params)
}
// 分页改变
const onHandleSizeChange = (val: number) => {
  state.queryParams.page_size = val;
  onSearch(state.queryParams)
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
  state.queryParams.page_num = val;
  onSearch(state.queryParams)
};
// 页面加载时
onMounted(() => {
  onSearch(state.queryParams)
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
