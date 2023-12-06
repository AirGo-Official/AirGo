<template>
  <div class="container layout-padding">
    <el-card shadow="hover" class="layout-padding-auto" v-loading="state.isLoadingTable">
      <el-row>
        <el-col :xs="24" :sm="24" :md="24" :lg="8" :xl="8">
          <el-date-picker
              size="default"
              v-model="state.this_month"
              type="datetimerange"
              :shortcuts="shortcuts"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-col>
        <el-col :xs="24" :sm="24" :md="24" :lg="12" :xl="12">
            <el-button @click="onGetNode()" size="default" type="primary" class="ml10">
              <el-icon>
                <ele-Search/>
              </el-icon>
              查询
            </el-button>
            <el-button size="default" type="success" class="ml10" @click="onOpenEditNode('新建节点','vless')">
              <el-icon>
                <ele-FolderAdd/>
              </el-icon>
              新增节点
            </el-button>

            <el-button size="default" type="warning" class="ml10" @click="onOpenNodeSortDialog">
              <el-icon>
                <DCaret/>
              </el-icon>
              排序
            </el-button>

            <el-button size="default" type="primary" class="ml10" @click="onOpenNodeSharedDialog">
              <el-icon>
                <Share/>
              </el-icon>
              共享节点管理
            </el-button>
        </el-col>
      </el-row>


      <el-table :data="nodeManageData.nodes.node_list" height="100%" stripe style="width: 100%;flex: 1;" @sort-change="sortChange">
        <el-table-column fixed type="index" label="序号" width="60"/>
        <el-table-column prop="remarks" label="节点名称" show-overflow-tooltip width="200" sortable="custom"></el-table-column>
        <el-table-column prop="id" label="节点ID" show-overflow-tooltip width="60" sortable="custom"></el-table-column>
        <el-table-column prop="node_order" label="排序" show-overflow-tooltip width="60" sortable="custom"></el-table-column>
        <el-table-column prop="address" label="节点地址" show-overflow-tooltip width="150" sortable="custom"></el-table-column>
        <el-table-column prop="port" label="节点端口" width="80" show-overflow-tooltip sortable="custom"></el-table-column>
        <el-table-column prop="sort" label="协议类型" width="120" show-overflow-tooltip sortable="custom">
          <template #default="scope">
            <el-button type="success" v-if="scope.row.node_type ==='vmess'">vmess</el-button>
            <el-button type="warning" v-if="scope.row.node_type ==='vless'">vless</el-button>
            <el-button type="info" v-if="scope.row.node_type ==='trojan'">trojan</el-button>
            <el-button type="danger" v-if="scope.row.node_type ==='shadowsocks'">shadowsocks</el-button>
            <el-button type="primary" v-if="scope.row.node_type ==='hysteria'">hysteria</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="total_up" label="上行流量(GB)" show-overflow-tooltip width="200" sortable="custom">
          <template #default="scope">
            <el-tag type="warning">{{ scope.row.total_up / 1024 / 1024 / 1024 }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="total_down" label="下行流量(GB)" show-overflow-tooltip width="200" sortable="custom">
          <template #default="scope">
            <el-tag type="warning">{{ scope.row.total_down / 1024 / 1024 / 1024 }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="enable_transfer" label="节点类型" show-overflow-tooltip sortable="custom">
          <template #default="scope">
            <el-tag type="warning" v-if="scope.row.enable_transfer">中转</el-tag>
            <el-tag type="success" v-else>直连</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="enabled" label="节点状态" show-overflow-tooltip sortable="custom">
          <template #default="scope">
            <el-tag type="success" v-if="scope.row.enabled">启用</el-tag>
            <el-tag type="danger" v-else>禁用</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="node_speedlimit" label="限速" show-overflow-tooltip sortable="custom"></el-table-column>
        <el-table-column prop="traffic_rate" label="倍率" show-overflow-tooltip sortable="custom"></el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="scope">
            <el-button size="small" text type="primary"
                       @click="onOpenEditNode('编辑节点', scope.row)">编辑
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
          @size-change="onHandleSizeChange" @current-change="onHandleCurrentChange"
          v-model:current-page="reportStoreData.reportParams.value.pagination.page_num"
          v-model:page-size="reportStoreData.reportParams.value.pagination.page_size"
          :total="nodeManageData.nodes.total"
      />
    </el-card>
    <NodeDialog ref="nodeDialogRef" @refresh="onGetNode()"/>
    <NodeSortDialog ref="nodeSortDialogRef" @refresh="onGetNode()"></NodeSortDialog>
    <NodeSharedDialog ref="nodeSharedDialogRef"></NodeSharedDialog>

  </div>
</template>

<script setup lang="ts">

import {defineAsyncComponent, onBeforeMount, onMounted, reactive, ref} from "vue";
import {storeToRefs} from "pinia";
import {useNodeStore} from "/@/stores/nodeStore";
import {useUserStore} from "/@/stores/userStore";
import {ElMessageBox} from "element-plus";
import {useAccessStore} from "/@/stores/accessStore";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
import {useReportStore} from "/@/stores/reportStore";
import {formatDate} from "/@/utils/formatTime";

const NodeDialog = defineAsyncComponent(() => import('/@/views/admin/node/dialog_edit.vue'))
const NodeSortDialog = defineAsyncComponent(() => import('/@/views/admin/node/dialog_node_sort.vue'))
const NodeSharedDialog = defineAsyncComponent(() => import('/@/views/admin/node/dialog_node_shared.vue'))
const nodeDialogRef = ref()
const nodeSortDialogRef = ref()
const nodeSharedDialogRef = ref()
const accessStore = useAccessStore()
const accessStoreData = storeToRefs(accessStore)
const nodeStore = useNodeStore()
const {nodeManageData} = storeToRefs(nodeStore)

const userStore = useUserStore()
const {userInfos} = storeToRefs(userStore)
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
const reportStore = useReportStore()
const reportStoreData = storeToRefs(reportStore)
//时间范围
const shortcuts = [
  {
    text: '上周',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
      return [start, end]
    },
  },
  {
    text: '上月',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
      return [start, end]
    },
  },
  {
    text: '最近3个月',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 90)
      return [start, end]
    },
  },
]
//定义参数
const state = reactive({
  isLoadingTable: false,
  this_month: [''],
})

//打开新建节点，修改节点弹窗
function onOpenEditNode(title: string, row?: NodeInfo) {1
  nodeDialogRef.value.openDialog(title, row)
}

//打开节点排序弹窗
function onOpenNodeSortDialog() {
  nodeSortDialogRef.value.openDialog()
}

//打开共享节点弹窗
function onOpenNodeSharedDialog() {
  nodeSharedDialogRef.value.openDialog()
}
//初始化查询参数
const defaultFieldParams = (start: string, end: string) => {
  reportStoreData.reportParams.value.table_name = 'node'
  reportStoreData.reportParams.value.field_params_list = [
    {field: 'created_at', field_chinese_name: '', field_type: '', condition: '>', condition_value: start, operator: ''},
    {field: 'created_at', field_chinese_name: '', field_type: '', condition: '<', condition_value: end, operator: 'AND',}
  ] as FieldParams[]
  reportStoreData.reportParams.value.pagination = {page_num: 1, page_size: 30, order_by: 'node_order',} as Pagination

}

//查询节点
function onGetNode() {
  nodeStore.getNodeWithTraffic(reportStoreData.reportParams.value)
}
function initDate() {
  // 目标时间范围格式： "2023-05-09 11:56:02"
  let currentDate = new Date();
  let currentY = currentDate.getFullYear();
  let currentM = currentDate.getMonth() + 1;
  // let MonthDayNum = new Date(currentY,currentM,0).getDate();  //计算当月的天数

  //当月
  let startDate = new Date(currentY, currentM - 1, 1);
  let endDate = new Date(currentY, currentM, 0, 23, 59, 59); // new Date(2020,11,0);//表示2020/11/30这天
  let thisMonthStart = formatDate(startDate, "YYYY-mm-dd HH:MM:SS")
  let thisMonthEnd = formatDate(endDate, "YYYY-mm-dd HH:MM:SS")
  //上月
  state.this_month = [thisMonthStart, thisMonthEnd]
}

//删除节点
function onRowDel(row: NodeInfo) {
  ElMessageBox.confirm(`此操作将永久删除节点：${row.remarks}，是否继续?`, '提示', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        nodeStore.deleteNode(row)
        setTimeout(() => {
          onGetNode()
        }, 500);
      })
      .catch(() => {
      });
}

// 分页改变
const onHandleSizeChange = (val: number) => {
  reportStoreData.reportParams.value.pagination.page_size = val;
  onGetNode()
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
  reportStoreData.reportParams.value.pagination.page_num = val;
  onGetNode()
};
//排序监听
const sortChange = (column: any) => {
  //处理嵌套字段
  let p = (column.prop as string)
  if (p.indexOf('.') !== -1) {
    p = p.slice(p.indexOf('.')+1)
  }
  switch (column.order){
    case 'ascending':
      reportStoreData.reportParams.value.pagination.order_by=p+" ASC"
      break
    default:
      reportStoreData.reportParams.value.pagination.order_by=p+" DESC"
      break
  }
  onGetNode()

}
onBeforeMount(()=>{
  initDate()
  defaultFieldParams(state.this_month[0],state.this_month[1])
});

onMounted(() => {
  onGetNode() //获取全部节点
  accessStore.getRoutesList(accessStoreData.params.value)//获取全部access
});

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