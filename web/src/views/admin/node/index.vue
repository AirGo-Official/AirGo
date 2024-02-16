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
              :range-separator="$t('message.common.to')"
              :start-placeholder="$t('message.common.startDate')"
              :end-placeholder="$t('message.common.endDate')"
              value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-col>
        <el-col :xs="24" :sm="24" :md="24" :lg="12" :xl="12">
            <el-button @click="onGetNode()" size="default" type="primary" class="ml10">
              <el-icon>
                <ele-Search/>
              </el-icon>
              {{$t('message.common.query')}}
            </el-button>
            <el-button size="default" type="success" class="ml10" @click="onOpenEditNode('add',{})">
              <el-icon>
                <ele-FolderAdd/>
              </el-icon>
              {{$t('message.common.add')}}
            </el-button>

            <el-button size="default" type="warning" class="ml10" @click="onOpenNodeSortDialog">
              <el-icon>
                <DCaret/>
              </el-icon>
              {{$t('message.common.sort')}}
            </el-button>

            <el-button size="default" type="primary" class="ml10" @click="onOpenNodeSharedDialog">
              <el-icon>
                <Share/>
              </el-icon>
              {{$t('message.adminNode.sharedNode')}}
            </el-button>
        </el-col>
      </el-row>
      <el-table :data="nodeStoreData.nodeList.value.data" height="100%" stripe style="width: 100%;flex: 1;" @sort-change="sortChange">
        <el-table-column fixed type="index" :label="$t('message.adminNode.NodeInfo.index')" width="60"/>
        <el-table-column prop="id" :label="$t('message.adminNode.NodeInfo.id')" show-overflow-tooltip width="80" sortable="custom" fixed></el-table-column>
        <el-table-column prop="remarks" :label="$t('message.adminNode.NodeInfo.remarks')" show-overflow-tooltip width="200" sortable="custom"></el-table-column>
        <el-table-column prop="enabled" :label="$t('message.adminNode.NodeInfo.enabled')" width="120" show-overflow-tooltip sortable="custom">
          <template #default="scope">
            <el-tag type="success" v-if="scope.row.enabled">{{$t('message.common.enable')}}</el-tag>
            <el-tag type="danger" v-else>{{$t('message.common.disable')}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="enable_transfer" :label="$t('message.adminNode.NodeInfo.enable_transfer')" width="120" show-overflow-tooltip sortable="custom">
          <template #default="scope">
            <el-tag type="warning" v-if="scope.row.enable_transfer">{{$t('message.adminNode.NodeInfo.node_type_transfer')}}</el-tag>
            <el-tag type="success" v-else>{{$t('message.adminNode.NodeInfo.node_type_direct')}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="address" :label="$t('message.adminNode.NodeInfo.address')" show-overflow-tooltip width="150" sortable="custom">
          <template #default="{row}">
            <span v-if="row.enable_transfer">{{row.transfer_address}}</span>
            <span v-else>{{row.address}}</span>
          </template>
        </el-table-column>
        <el-table-column prop="port" :label="$t('message.adminNode.NodeInfo.port')" width="100" show-overflow-tooltip sortable="custom"></el-table-column>
        <el-table-column prop="node_type" :label="$t('message.adminNode.NodeInfo.node_type')" width="120" show-overflow-tooltip sortable="custom">
          <template #default="scope">
            <el-button type="success" v-if="scope.row.node_type ===constantStore.NODE_TYPE_VMESS">vmess</el-button>
            <el-button type="warning" v-if="scope.row.node_type ===constantStore.NODE_TYPE_VLESS">vless</el-button>
            <el-button type="info" v-if="scope.row.node_type ===constantStore.NODE_TYPE_TROJAN">trojan</el-button>
            <el-button type="danger" v-if="scope.row.node_type ===constantStore.NODE_TYPE_SHADOWSOCKS">shadowsocks</el-button>
            <el-button type="primary" v-if="scope.row.node_type ===constantStore.NODE_TYPE_HYSTERIA">hysteria</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="total_up" :label="$t('message.adminNode.NodeInfo.total_up')" show-overflow-tooltip width="200">
          <template #default="scope">
            <el-tag type="warning">{{ scope.row.total_up / 1024 / 1024 / 1024 }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="total_down" :label="$t('message.adminNode.NodeInfo.total_down')" show-overflow-tooltip width="200">
          <template #default="scope">
            <el-tag type="warning">{{ scope.row.total_down / 1024 / 1024 / 1024 }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="node_speed_limit" :label="$t('message.adminNode.NodeInfo.node_speed_limit')" show-overflow-tooltip sortable="custom"></el-table-column>
        <el-table-column prop="traffic_rate" :label="$t('message.adminNode.NodeInfo.traffic_rate')" show-overflow-tooltip sortable="custom"></el-table-column>
        <el-table-column :label="$t('message.common.operate')" width="100">
          <template #default="scope">
            <el-button size="small" text type="primary"
                       @click="onOpenEditNode('edit', scope.row)">{{$t('message.common.modify')}}
            </el-button>
            <el-button size="small" text type="primary"
                       @click="onRowDel(scope.row)">{{$t('message.common.delete')}}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
          background
          :page-sizes="[10, 30, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="onHandleSizeChange" @current-change="onHandleCurrentChange"
          v-model:current-page="state.queryParams.pagination.page_num"
          v-model:page-size="state.queryParams.pagination.page_size"
          :total="nodeStoreData.nodeList.value.total"
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
import {useAdminNodeStore} from "/@/stores/admin_logic/nodeStore";
import {ElMessageBox} from "element-plus";
import {formatDate} from "/@/utils/formatTime";
import { useConstantStore } from "/@/stores/constantStore";
import { useI18n } from "vue-i18n";

const NodeDialog = defineAsyncComponent(() => import('/@/views/admin/node/dialog_edit.vue'))
const NodeSortDialog = defineAsyncComponent(() => import('/@/views/admin/node/dialog_node_sort.vue'))
const NodeSharedDialog = defineAsyncComponent(() => import('/@/views/admin/node/dialog_node_shared.vue'))
const nodeDialogRef = ref()
const nodeSortDialogRef = ref()
const nodeSharedDialogRef = ref()
const nodeStore = useAdminNodeStore()
const nodeStoreData = storeToRefs(nodeStore)
const constantStore = useConstantStore()
const {t} = useI18n()
//时间范围
const shortcuts = [
  {
    text: t('message.common.lastWeek'),
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
      return [start, end]
    },
  },
  {
    text: t('message.common.lastMonth'),
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
      return [start, end]
    },
  },
  {
    text: t('message.common.lastThreeMonths'),
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
  queryParams:{
    table_name: 'node',
    field_params_list: [
      {field: 'created_at', field_chinese_name: '', field_type: '', condition: '>', condition_value: "", operator: ''},
      {field: 'created_at', field_chinese_name: '', field_type: '', condition: '<', condition_value: "", operator: 'AND',}
    ] as FieldParams[],
    pagination: {page_num: 1, page_size: 30, order_by: 'node_order',} as Pagination,//分页参数
  } as QueryParams,
})

//打开新建节点，修改节点弹窗
function onOpenEditNode(title: string, row?: NodeInfo) {
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

//查询节点
function onGetNode() {
  state.queryParams.field_params_list[0].condition_value = state.this_month[0]
  state.queryParams.field_params_list[1].condition_value = state.this_month[1]
  nodeStore.getNodeWithTraffic(state.queryParams)
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
  ElMessageBox.confirm(t('message.common.message_confirm_delete'), t('message.common.tip'), {
    confirmButtonText: t('message.common.button_confirm'),
    cancelButtonText: t('message.common.button_cancel'),
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
  state.queryParams.pagination.page_size = val;
  onGetNode()
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
  state.queryParams.pagination.page_num = val;
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
      state.queryParams.pagination.order_by=p+" ASC"
      break
    default:
      state.queryParams.pagination.order_by=p+" DESC"
      break
  }
  onGetNode()

}
onBeforeMount(()=>{
  initDate()
});

onMounted(() => {
  onGetNode()
  nodeStore.getAccessList({
    table_name: 'access',
      field_params_list: [
      {field: 'id', condition: '<>', condition_value: '',},
    ] as FieldTableNew[],
      pagination: { page_num: 1, page_size: 9999,} as Pagination,
  } as QueryParams)
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