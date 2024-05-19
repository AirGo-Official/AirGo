<template>
  <div class="container layout-padding">
    <el-card shadow="hover" class="layout-padding-auto">
      <div class="mb15">
        <el-input v-model="state.queryParams.field_params_list[0].condition_value" size="default"
                  style="max-width: 180px"></el-input>
        <el-button @click="findUser" size="default" type="primary" class="ml10">
          <el-icon>
            <ele-Search />
          </el-icon>
          {{ $t("message.adminUser.query") }}
        </el-button>
        <el-button size="default" type="success" class="ml10" @click="openEditDialog('add')">
          <el-icon>
            <ele-FolderAdd />
          </el-icon>
          {{ $t("message.adminUser.add_user") }}
        </el-button>
        <el-button size="default" color="blue" class="ml10" @click="onShowCollapse">
          <el-icon>
            <ele-Search />
          </el-icon>
          {{ $t("message.adminUser.advanced_query") }}
        </el-button>
        <div v-if="state.isShowCollapse">
          <!--          report组件-->
          <ReportComponent ref="reportRef" @getReportData="emitsParams"></ReportComponent>
        </div>
      </div>
      <el-table :data="userStoreData.userList.value.data" stripe style="width: 100%" @sort-change="sortChange">
        <el-table-column type="index" :label="$t('message.adminUser.SysUser.index')" width="60" fixed />
        <el-table-column prop="avatar" :label="$t('message.adminUser.SysUser.avatar')" show-overflow-tooltip width="150"
                         sortable="custom">
          <template #default="{row}">
            <el-avatar size="small" :src="row.avatar" @error="true">
              <img
                src="https://cube.elemecdn.com/e/fd/0fc7d20532fdaf769a25683617711png.png"
              />
            </el-avatar>
          </template>
        </el-table-column>
        <el-table-column prop="user_name" :label="$t('message.adminUser.SysUser.user_name')" show-overflow-tooltip
                         width="150"
                         sortable="custom"></el-table-column>
        <el-table-column prop="id" :label="$t('message.adminUser.SysUser.id')" show-overflow-tooltip width="80"
                         sortable="custom"></el-table-column>
        <el-table-column prop="created_at" :label="$t('message.adminUser.SysUser.created_at')" show-overflow-tooltip
                         width="150" sortable="custom">
          <template #default="{row}">
            <span>{{ DateStrToTime(row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="enable" :label="$t('message.adminUser.SysUser.enable')" show-overflow-tooltip width="100"
                         sortable="custom">
          <template #default="scope">
            <el-tag type="success" v-if="scope.row.enable">{{ $t("message.common.enable") }}</el-tag>
            <el-tag type="danger" v-else>{{ $t("message.common.disable") }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('message.adminUser.SysUser.balance')" show-overflow-tooltip sortable="custom">
          <template #default="scope">
            <el-tag type="info">
              {{ scope.row.balance }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('message.common.operate')" width="200px">
          <template #default="scope">
            <el-button size="small" text type="primary" @click="openCustomerServiceDrawer(scope.row)">
              {{ $t("message.adminUser.customerService") }}
            </el-button>
            <el-button size="small" text type="primary" @click="openEditDialog('edit', scope.row)">
              {{ $t("message.common.modify") }}
            </el-button>
            <el-button size="small" text type="danger" @click="deleteUser(scope.row)">
              {{ $t("message.common.delete") }}
            </el-button>
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
    <UserDialog ref="userDialogRef" @refresh="getUserList" />
    <CustomerServiceDrawer ref="CustomerServiceDrawerRef"></CustomerServiceDrawer>
  </div>
</template>

<script setup lang="ts" name="systemUser">
import { defineAsyncComponent, onMounted, reactive, ref } from "vue";
import { ElMessageBox } from "element-plus";
import { storeToRefs } from "pinia";

import { DateStrToTime } from "/@/utils/formatTime";
import { useAdminUserStore } from "/@/stores/admin_logic/userStore";
import { useI18n } from "vue-i18n";

const userStore = useAdminUserStore();
const userStoreData = storeToRefs(userStore);
const UserDialog = defineAsyncComponent(() => import("/@/views/admin/user/dialog_user_edit.vue"));
const ReportComponent = defineAsyncComponent(() => import("/@/components/report/index.vue"));
const CustomerServiceDrawer = defineAsyncComponent(() => import("/@/views/admin/user/drawer_customer_service.vue"));
const userDialogRef = ref();
const reportRef = ref();
const CustomerServiceDrawerRef = ref();
const { t } = useI18n();

// 定义变量内容
const state = reactive({
  activeCollapseNames: "1", //当前激活的折叠面板
  isShowCollapse: false,
  queryParams: {
    table_name: "user",
    field_params_list: [
      {
        field: "user_name",
        field_chinese_name: "",
        field_type: "",
        condition: "<>",
        condition_value: "",
        operator: ""
      }
    ] as FieldParams[],
    pagination: {
      page_num: 1, page_size: 30, order_by: "id"
    } as Pagination//分页参数
  } as QueryParams
});

// 打开用户弹窗
const openEditDialog = (type: string, row?: SysUser) => {
  userDialogRef.value.openDialog(type, row);
};

//
const openCustomerServiceDrawer = (row: SysUser) => {
  CustomerServiceDrawerRef.value.openDrawer(row);
};
//
const findUser = () => {
  state.queryParams.field_params_list[0].condition = "like";
  getUserList();
};

//接受子组件传值
const emitsParams=(params:QueryParams)=>{
state.queryParams=params
  getUserList()
}
//请求数据
const getUserList = () => {
  userStore.getUserList(state.queryParams);
};

// 删除用户
const deleteUser = (row: SysUser) => {
  ElMessageBox.confirm(t("message.common.message_confirm_delete"), t("message.common.tip"), {
    confirmButtonText: t("message.common.button_confirm"),
    cancelButtonText: t("message.common.button_cancel"),
    type: "warning"
  })
    .then(() => {
      userStore.deleteUser(row).then(() => {
        getUserList();
      });
    })
    .catch(() => {
    });
};
// 分页改变
const onHandleSizeChange = (val: number) => {
  if (state.isShowCollapse) {
    getUserList();
  } else {
    state.queryParams.pagination.page_size = val;
    getUserList();
  }
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
  if (state.isShowCollapse) {
    getUserList();
  } else {
    state.queryParams.pagination.page_num = val;
    getUserList();
  }
};
//排序监听
const sortChange = (column: any) => {
  //处理嵌套字段
  let p = (column.prop as string);
  if (p.indexOf(".") !== -1) {
    p = p.slice(p.indexOf(".") + 1);
  }
  switch (column.order) {
    case "ascending":
      state.queryParams.pagination.order_by = p + " ASC";
      break;
    default:
      state.queryParams.pagination.order_by = p + " DESC";
      break;
  }
  getUserList();
};
onMounted(() => {
  getUserList();
});

//开启高级查询折叠面板
const onShowCollapse = () => {
  state.isShowCollapse = !state.isShowCollapse;
  setTimeout(() => {
    if (state.isShowCollapse) {
      reportRef.value.openReportComponent("user");  //参数：user库表
    }
  }, 500);
};

</script>

<style scoped lang="scss">
.container {
  :deep(.el-card__body) {
    display: flex;
    flex-direction: column;
    flex: 1;
    overflow: auto;

    .el-table {
      flex: 1;
    }
  }
}
</style>


