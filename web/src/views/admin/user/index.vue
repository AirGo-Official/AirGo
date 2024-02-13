<template>
  <div >
    <el-card shadow="hover" class="layout-padding-auto">
      <div class="mb15">
        <el-input v-model="state.queryParams.field_params_list[0].condition_value" size="default"
                  placeholder="请输入用户名称"
                  style="max-width: 180px"></el-input>
        <el-button @click="findUser" size="default" type="primary" class="ml10">
          <el-icon>
            <ele-Search/>
          </el-icon>
          查询
        </el-button>
        <el-button size="default" type="success" class="ml10" @click="onOpenAddUser('add')">
          <el-icon>
            <ele-FolderAdd/>
          </el-icon>
          新增用户
        </el-button>
        <el-button size="default" color="blue" class="ml10" @click="onShowCollapse">
          <el-icon>
            <ele-Search/>
          </el-icon>
          高级查询
        </el-button>
        <div v-if="state.isShowCollapse">
          <!--          report组件-->
          <ReportComponent ref="reportRef" @getReportData="getUserList"></ReportComponent>
        </div>
      </div>
      <el-table :data="userStoreData.userList.value.data" stripe style="width: 100%" @sort-change="sortChange">
        <el-table-column type="index" label="序号" width="60" fixed/>
        <el-table-column prop="user_name" label="账户名称" show-overflow-tooltip width="150"
                         sortable="custom"></el-table-column>
        <el-table-column prop="id" label="账户ID" show-overflow-tooltip width="80" sortable="custom"></el-table-column>
        <el-table-column prop="created_at" label="创建日期" show-overflow-tooltip width="150" sortable="custom">
          <template #default="{row}">
            <span>{{ DateStrToTime(row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="enable" label="用户状态" show-overflow-tooltip width="100" sortable="custom">
          <template #default="scope">
            <el-tag type="success" v-if="scope.row.enable">启用</el-tag>
            <el-tag type="danger" v-else>禁用</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="余额" show-overflow-tooltip sortable="custom">
          <template #default="scope">
            <el-tag type="info">
              {{ scope.row.balance }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200px">
          <template #default="scope">
            <el-button size="small" text type="primary" @click="onOpenCustomerService(scope.row)">客户服务</el-button>
            <el-button size="small" text type="primary" @click="onOpenEditUser('edit', scope.row)">修改</el-button>
            <el-button size="small" text type="primary" @click="onRowDel(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
          background
          class="mt15"
          layout="total, sizes, prev, pager, next, jumper"
          :page-sizes="[10, 30, 50]"
          v-model:current-page="state.queryParams.pagination.page_num"
          v-model:page-size="state.queryParams.pagination.page_size"
          :total="userStoreData.userList.value.total"
          @size-change="onHandleSizeChange"
          @current-change="onHandleCurrentChange"
      >
      </el-pagination>
    </el-card>
    <UserDialog ref="userDialogRef" @refresh="getUserList"/>
    <CustomerServiceDialog ref="customerServiceDialogRef"></CustomerServiceDialog>
  </div>
</template>

<script setup lang="ts" name="systemUser">
import {defineAsyncComponent, onMounted, reactive, ref} from 'vue';
import {ElMessageBox} from 'element-plus';
import {storeToRefs} from 'pinia';

import {DateStrToTime} from "/@/utils/formatTime";
import { useAdminUserStore } from "/@/stores/admin_logic/userStore";
const userStore = useAdminUserStore()
const userStoreData = storeToRefs(userStore)
const UserDialog = defineAsyncComponent(() => import('/@/views/admin/user/dialog_user_edit.vue'));
const ReportComponent = defineAsyncComponent(() => import('/@/components/report/index.vue'))
const CustomerServiceDialog =  defineAsyncComponent(() => import('/@/views/admin/user/dialog_customer_service.vue'));
const userDialogRef = ref();
const reportRef = ref()
const customerServiceDialogRef = ref()


// 定义变量内容
const state = reactive({
  activeCollapseNames: '1', //当前激活的折叠面板
  isShowCollapse: false,
  queryParams:{
    table_name: 'user',
    field_params_list: [
      { field: 'user_name', field_chinese_name: '', field_type: '', condition: '<>', condition_value: '', operator: '', }
    ] as FieldParams[],
    pagination: { page_num: 1, page_size: 30, order_by: 'id',
    } as Pagination,//分页参数
  } as QueryParams,
});

// 打开新增用户弹窗
const onOpenAddUser = (type: string) => {
  userDialogRef.value.openDialog(type);
};
// 打开修改用户弹窗
const onOpenEditUser = (type: string, row: SysUser) => {
  userDialogRef.value.openDialog(type, row);
};

//
const onOpenCustomerService=(row: SysUser)=>{
  customerServiceDialogRef.value.onOpenDrawer(row);
}
//
const findUser = () => {
  state.queryParams.field_params_list[0].condition = 'like'
  getUserList()
}
//请求数据
const getUserList = () => {
  userStore.getUserList(state.queryParams)
}

// 删除用户
const onRowDel = (row: SysUser) => {
  ElMessageBox.confirm(`此操作将永久删除账户名称：“${row.user_name}”，是否继续?`, '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        userStore.deleteUser(row)
        setTimeout(() => {
          getUserList()
        }, 500)
      })
      .catch(() => {
      });
};
// 分页改变
const onHandleSizeChange = (val: number) => {
  if (state.isShowCollapse) {
    getUserList()
  } else {
    state.queryParams.pagination.page_size = val;
    getUserList()
  }
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
  if (state.isShowCollapse) {
    getUserList()
  } else {
    state.queryParams.pagination.page_num = val;
    getUserList()
  }
};
//排序监听
const sortChange = (column: any) => {
  //处理嵌套字段
  let p = (column.prop as string)
  if (p.indexOf('.') !== -1) {
    p = p.slice(p.indexOf('.') + 1)
  }
  switch (column.order) {
    case 'ascending':
      state.queryParams.pagination.order_by = p + " ASC"
      break
    default:
      state.queryParams.pagination.order_by = p + " DESC"
      break
  }
  getUserList()
}
onMounted(() => {
  getUserList()
});

//开启高级查询折叠面板
const onShowCollapse = () => {
  state.isShowCollapse = !state.isShowCollapse
  setTimeout(() => {
    if (state.isShowCollapse) {
      reportRef.value.openReportComponent("user")  //参数：user库表
    }
  }, 500)
}

</script>

<style scoped lang="scss">
.container {
  :deep(.el-card__body) {
    display: flex;
    flex-direction: column;
    flex: 1;
    overflow: auto;

    .el-table {
      overflow: auto;
      flex: 1;
    }
  }
}
</style>


