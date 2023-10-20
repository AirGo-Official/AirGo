<template>
  <div class="container layout-padding">
    <el-card shadow="hover" class="layout-padding-auto">
      <div class="mb15">
        <el-input v-model="state.params.search" size="default" placeholder="请输入用户名称"
                  style="max-width: 180px"></el-input>
        <el-button @click="userStore.getUserList(state.params)" size="default" type="primary" class="ml10">
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
        <el-button size="default" type="primary" class="ml10" @click="onShowCollapse">
          <el-icon>
            <ele-Search/>
          </el-icon>
          高级查询
        </el-button>
        <div v-if="state.isShowCollapse">
          <!--          report组件-->
          <ReportComponent ref="reportRef" @getReportData="getReportDataHandler"></ReportComponent>
        </div>
      </div>
      <el-table :data="userManageData.users.user_list" stripe style="width: 100%;flex: 1;">
        <el-table-column type="index" label="序号" width="60" fixed/>
        <el-table-column prop="user_name" label="账户名称" show-overflow-tooltip width="150"></el-table-column>
        <el-table-column prop="id" label="账户ID" show-overflow-tooltip width="60"></el-table-column>
        <el-table-column prop="created_at" label="创建日期" show-overflow-tooltip width="150">
          <template #default="{row}">
            <span>{{ DateStrtoTime(row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="enable" label="用户状态" show-overflow-tooltip width="80">
          <template #default="scope">
            <el-tag type="success" v-if="scope.row.enable">启用</el-tag>
            <el-tag type="danger" v-else>禁用</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="subscribe_info.sub_status" label="订阅状态" show-overflow-tooltip width="80">
          <template #default="scope">
            <el-tag type="success" v-if="scope.row.subscribe_info.sub_status">启用</el-tag>
            <el-tag type="danger" v-else>禁用</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="subscribe_info.expired_at" label="订阅到期时间" show-overflow-tooltip width="150">
          <template #default="scope">
            <el-tag type="info">
              {{ DateStrtoTime(scope.row.subscribe_info.expired_at) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="subscribe_info.subscribe_url" label="通用订阅url" show-overflow-tooltip width="400">
          <template #default="scope">
            <el-tag type="info">
              {{ serverStore.publicServerConfig.backend_url }}user/getSub?link={{
                scope.row.subscribe_info.subscribe_url
              }}&type=v2ray
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="subscribe_info.goods_id" label="商品ID" show-overflow-tooltip
                         width="60"></el-table-column>
        <el-table-column prop="subscribe_info.t" label="总流量(GB)" show-overflow-tooltip>
          <template #default="scope">
            <el-tag type="info">{{ scope.row.subscribe_info.t / 1024 / 1024 / 1024 }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="已用流量(GB)" show-overflow-tooltip width="100px">
          <template #default="scope">
            <el-tag type="info">
              {{ (scope.row.subscribe_info.d + scope.row.subscribe_info.u) / 1024 / 1024 / 1024 }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="余额" show-overflow-tooltip>
          <template #default="scope">
            <el-tag type="info">
              {{ scope.row.remain }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="100">
          <template #default="scope">
            <el-button size="small" text type="primary" @click="onOpenEditUser('edit', scope.row)"
            >修改
            </el-button
            >
            <el-button size="small" text type="primary" @click="onRowDel(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
          background
          class="mt15"
          layout="total, sizes, prev, pager, next, jumper"
          :page-sizes="[10, 30, 50]"
          v-model:current-page="state.params.page_num"
          v-model:page-size="state.params.page_size"
          :total="userManageData.users.total"
          @size-change="onHandleSizeChange"
          @current-change="onHandleCurrentChange"
      >
      </el-pagination>
    </el-card>
    <UserDialog ref="userDialogRef" @refresh=""/>
  </div>
</template>

<script setup lang="ts" name="systemUser">
import {defineAsyncComponent, onMounted, reactive, ref} from 'vue';
import {ElMessageBox} from 'element-plus';


//store
import {storeToRefs} from 'pinia';
import {useUserStore} from '/@/stores/userStore'
import {useServerStore} from "/@/stores/serverStore";
import {useReportStore} from "/@/stores/reportStore"

import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
import {DateStrtoTime} from "../../../utils/formatTime";

const serverStore = useServerStore()
const userStore = useUserStore()
const {userManageData} = storeToRefs(userStore)
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
const reportStore = useReportStore()
const {reportTable} = storeToRefs(reportStore)
const UserDialog = defineAsyncComponent(() => import('/@/views/admin/user/dialog.vue'));
const ReportComponent = defineAsyncComponent(() => import('/@/components/report/index.vue'))
const userDialogRef = ref();
const reportRef = ref()

// 定义变量内容

const state = reactive({
  activeCollapseNames: '1', //当前激活的折叠面板
  isShowCollapse: false,
  params: {
    search: '',
    page_num: 1,
    page_size: 30,
  },
});
// 打开新增用户弹窗
const onOpenAddUser = (type: string) => {
  userDialogRef.value.openDialog(type);
};
// 打开修改用户弹窗
const onOpenEditUser = (type: string, row: SysUser) => {
  userDialogRef.value.openDialog(type, row);
};

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
          userStore.getUserList(state.params)
        }, 500)
      })
      .catch(() => {
      });
};
// 分页改变
const onHandleSizeChange = (val: number) => {
  if (state.isShowCollapse) {
    getReportDataHandler()
  } else {
    state.params.page_size = val;
    userStore.getUserList(state.params)
  }
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
  if (state.isShowCollapse) {
    getReportDataHandler()
  } else {
    state.params.page_num = val;
    userStore.getUserList(state.params)
  }
};
// 页面加载时
onMounted(() => {
  userStore.getUserList(state.params)
});

//开启高级查询折叠面板
const onShowCollapse = () => {
  state.isShowCollapse = !state.isShowCollapse
  //防止子组件渲染太慢，导致undefined问题
  setTimeout(() => {
    if (state.isShowCollapse) {
      reportRef.value.openReportComponent("user")  //参数：user库表
    }
  }, 500)
}
//请求数据
const getReportDataHandler = () => {
  //拼接分页参数
  reportTable.value.pagination_params = state.params;
  // console.log("参数：",reportTable.value)
  //请求数据
  request(apiStoreData.api.value.report_reportSubmit, reportTable.value).then((res) => {
    userManageData.value.users.user_list = res.data.table_data
    userManageData.value.users.total = res.data.total
  })
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
      flex: 1;
    }
  }
}
</style>


