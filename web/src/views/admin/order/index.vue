<template>
  <div class="container layout-padding">
    <el-card shadow="hover" class="layout-padding-auto">
      <div class="mb15">
        <el-input v-model="state.queryParams.field_params_list[0].condition_value" placeholder="请输入订单号"
                  style="max-width: 180px" size="default"></el-input>
        <el-button size="default" type="primary" class="ml10" @click="getOrderList()">
          <el-icon>
            <ele-Search/>
          </el-icon>
          查询
        </el-button>
        <el-button size="default" type="success" class="ml10" @click="onShowCollapse">
          <el-icon>
            <ele-Search/>
          </el-icon>
          高级查询
        </el-button>

        <el-collapse v-if="state.isShowCollapse" v-model="state.activeCollapseNames">
          <el-collapse-item name="1">
            <!--          report组件-->
            <ReportComponent ref="reportRef" @getReportData="getReportData"></ReportComponent>
          </el-collapse-item>
        </el-collapse>
      </div>
      <el-table :data="shopStoreData.orderList.value.data" stripe style="width: 100%;flex: 1;" @sort-change="sortChange">
        <el-table-column type="index" label="序号" fixed width="60px"/>
        <el-table-column prop="out_trade_no" label="订单号" width="200" sortable="custom"/>
        <el-table-column prop="id" label="订单ID" width="100px" sortable="custom"/>
        <el-table-column prop="created_at" label="下单日期" width="150" sortable="custom">
          <template #default="scope">
            <el-tag type="success">{{ DateStrToTime(scope.row.created_at) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="user_name" label="用户" width="180" sortable="custom"/>
        <el-table-column prop="goods_id" label="商品ID" show-overflow-tooltip width="100" sortable="custom"/>
        <el-table-column prop="subject" label="商品标题" show-overflow-tooltip width="200" sortable="custom"/>
        <el-table-column prop="total_amount" label="订单金额" show-overflow-tooltip width="100" sortable="custom"/>
        <el-table-column prop="receipt_amount" label="实收金额" show-overflow-tooltip width="100" sortable="custom"/>
        <el-table-column prop="deliver_type" label="发货状态" show-overflow-tooltip width="100" sortable="custom">
          <template #default="{row}">
            <el-tag type="info" v-if="row.deliver_type===constantStore.DELIVER_TYPE_NONE">无需发货</el-tag>
            <el-tag type="success" v-if="row.deliver_type===constantStore.DELIVER_TYPE_AUTO">自动发货</el-tag>
            <el-tag type="warning" v-if="row.deliver_type===constantStore.DELIVER_TYPE_MANUAL">手动发货</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="trade_status" label="交易状态" show-overflow-tooltip sortable="custom" width="100">
          <template #default="scope">
            <el-tag type="info" v-if="scope.row.trade_status===constantStore.ORDER_STATUS_CREATED">已创建</el-tag>
            <el-tag type="warning" v-else-if="scope.row.trade_status===constantStore.ORDER_STATUS_WAIT_BUYER_PAY">等待买家付款</el-tag>
            <el-tag type="success" v-else-if="scope.row.trade_status===constantStore.ORDER_STATUS_TRADE_SUCCESS">支付成功</el-tag>
            <el-tag type="danger" v-else-if="scope.row.trade_status===constantStore.ORDER_STATUS_TRADE_CLOSED">交易关闭</el-tag>
            <el-tag type="success" v-else-if="scope.row.trade_status===constantStore.ORDER_STATUS_TRADE_FINISHED">交易结束</el-tag>
            <el-tag type="success" v-else-if="scope.row.trade_status===constantStore.ORDER_STATUS_COMPLETED">已完成</el-tag>
            <el-tag type="danger" v-else>未知状态</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="scope">
            <el-button size="small" text type="primary"
                       @click="onEditOrder(scope.row)">编辑
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination background
                     class="mt15"
                     layout="total, sizes, prev, pager, next, jumper"
                     :page-sizes="[10, 30, 50]"
                     v-model:current-page="state.queryParams.pagination.page_num"
                     v-model:page-size="state.queryParams.pagination.page_size"
                     :total="shopStoreData.orderList.value.total"
                     @size-change="onHandleSizeChange"
                     @current-change="onHandleCurrentChange">
      </el-pagination>
    </el-card>
    <DialogEditOrder ref="DialogEditOrderRef"></DialogEditOrder>
  </div>
</template>

<script setup lang="ts">
import {defineAsyncComponent, onBeforeMount, onMounted, reactive, ref} from "vue";
import {DateStrToTime} from "/@/utils/formatTime"

import {storeToRefs} from "pinia";
import {ElMessage, ElMessageBox} from "element-plus";
import { useAdminShopStore } from "/@/stores/admin_logic/shopStore";
import { useConstantStore } from "/@/stores/constantStore";

const shopStore = useAdminShopStore()
const shopStoreData = storeToRefs(shopStore)
const ReportComponent = defineAsyncComponent(() => import('/@/components/report/index.vue'))
const reportRef = ref()
const DialogEditOrder = defineAsyncComponent(() => import('/@/views/admin/order/dialog_edit.vue'))
const DialogEditOrderRef = ref()
const constantStore = useConstantStore()


//定义参数
const state = reactive({
  activeCollapseNames: '1', //当前激活的折叠面板
  isShowCollapse: false,
  queryParams:{
    table_name: 'order',
    field_params_list: [
      {field:"out_trade_no",field_chinese_name:"",field_type:"",condition:"like",condition_value:"",operator:""}
    ] as FieldParams[],
    pagination: { page_num: 1, page_size: 30, order_by: 'id DESC',
    } as Pagination,//分页参数
  } as QueryParams,
})
//
const getOrderList = () => {
  shopStore.getOrderList(state.queryParams)
}
//编辑订单
const onEditOrder=(row: Order)=>{
  DialogEditOrderRef.value.openDialog(row)
}

// 分页改变
const onHandleSizeChange = (val: number) => {
  state.queryParams.pagination.page_size = val;
    getOrderList()
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
  state.queryParams.pagination.page_num = val;
    getOrderList()
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
  getOrderList()
}
//开启高级查询折叠面板
const onShowCollapse = () => {
  state.isShowCollapse = !state.isShowCollapse
  //防止子组件渲染太慢，导致undefined问题
  setTimeout(() => {
    if (state.isShowCollapse) {
      reportRef.value.openReportComponent("order")
    }
  }, 500)
}
//
const getReportData = (data: any) => {
  getOrderList()
}
onMounted(() => {
  getOrderList()
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