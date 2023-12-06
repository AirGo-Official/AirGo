<template>
  <div class="container layout-padding">
    <el-card shadow="hover" class="layout-padding-auto">
      <el-table :data="orderPersonal.allOrders.order_list" stripe fit height="100%" style="width: 100%;" @sort-change="sortChange">
        <el-table-column type="index" label="序号" fixed width="60px"/>
        <el-table-column prop="out_trade_no" label="订单号" width="200" sortable="custom"/>
        <el-table-column prop="id" label="订单ID" width="100px" sortable="custom"/>
        <el-table-column prop="created_at" label="下单日期" width="150" sortable="custom">
          <template #default="scope">
            <el-tag type="success">{{ DateStrtoTime(scope.row.created_at) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="user_name" label="用户" width="180" sortable="custom"/>
        <el-table-column prop="goods_id" label="商品ID" show-overflow-tooltip width="100" sortable="custom"/>
        <el-table-column prop="subject" label="商品标题" show-overflow-tooltip width="200" sortable="custom"/>
        <el-table-column prop="total_amount" label="订单金额" show-overflow-tooltip width="100" sortable="custom"/>
        <el-table-column prop="receipt_amount" label="实收金额" show-overflow-tooltip width="100" sortable="custom"/>
        <el-table-column prop="trade_status" label="交易状态" show-overflow-tooltip sortable="custom" width="100">
          <template #default="scope">
            <el-tag type="success" v-if="scope.row.trade_status==='TRADE_SUCCESS'">支付成功</el-tag>
            <el-tag type="warning" v-else-if="scope.row.trade_status==='WAIT_BUYER_PAY'">等待买家付款</el-tag>
            <el-tag type="danger" v-else-if="scope.row.trade_status==='TRADE_CLOSED'">交易超时关闭</el-tag>
            <el-tag type="success" v-else-if="scope.row.trade_status==='TRADE_FINISHED'">交易结束</el-tag>
            <el-tag type="info" v-else-if="scope.row.trade_status==='Created'">订单已创建</el-tag>
            <el-tag type="success" v-else-if="scope.row.trade_status==='Completed'">订单已完成</el-tag>
            <el-tag type="danger" v-else>未知状态</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button v-if="scope.row.trade_status === 'WAIT_BUYER_PAY' || scope.row.trade_status === 'Created'"
                       size="small" text type="primary"
                       @click="toPay(scope.row)">去支付
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination background
                     class="mt15"
                     layout="total, sizes, prev, pager, next, jumper"
                     :page-sizes="[10, 30, 50]"
                     v-model:current-page="reportStoreData.reportParams.value.pagination.page_num"
                     v-model:page-size="reportStoreData.reportParams.value.pagination.page_size"
                     :total="orderPersonal.allOrders.total"
                     @size-change="onHandleSizeChange"
                     @current-change="onHandleCurrentChange">
      </el-pagination>
    </el-card>
    <!-- 引入确认购买弹窗组件 -->
    <PurchaseDialog ref="PurchaseDialogRef" @openQRDialog="openQRDialog"></PurchaseDialog>
    <!-- 引入二维码弹窗 -->
    <QRDialog ref="QRDialogRef"></QRDialog>
  </div>
</template>

<script setup lang="ts">
import {storeToRefs} from "pinia";
import {defineAsyncComponent, onBeforeMount, onMounted, ref} from "vue";
import {useOrderStore} from "/@/stores/orderStore";
import {useShopStore} from "/@/stores/shopStore";
import {isMobile} from "/@/utils/other";
import {DateStrtoTime} from "/@/utils/formatTime"
import {useReportStore} from "/@/stores/reportStore";

const orderStore = useOrderStore()
const {orderPersonal} = storeToRefs(orderStore)
const shopStore = useShopStore()
const {shopData} = storeToRefs(shopStore)
const PurchaseDialog = defineAsyncComponent(() => import('/@/views/shop/dialog_purchase.vue'))
const QRDialog = defineAsyncComponent(() => import('/@/views/shop/dialog_QR.vue'))
const PurchaseDialogRef = ref()
const QRDialogRef = ref()
const reportStore = useReportStore()
const reportStoreData = storeToRefs(reportStore)

//
const getUserOrders=()=>{
  orderStore.getOrder(reportStoreData.reportParams.value)
}
// 分页改变
const onHandleSizeChange = (val: number) => {
  reportStoreData.reportParams.value.pagination.page_size = val;
  getUserOrders()
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
  reportStoreData.reportParams.value.pagination.page_num = val;
  getUserOrders()
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
  getUserOrders()
}
//初始化查询参数
const defaultFieldParams = () => {
  reportStoreData.reportParams.value.table_name = 'orders'
  reportStoreData.reportParams.value.field_params_list = [
    {field: 'out_trade_no', field_chinese_name: '', field_type: '', condition: 'like', condition_value: '', operator: '',} as FieldParams]
  reportStoreData.reportParams.value.pagination = {page_num: 1, page_size: 30, order_by: 'id DESC',} as Pagination
}
//
onBeforeMount(() => {
  defaultFieldParams()
});
onMounted(() => {
  getUserOrders()
});

//去支付流程
const toPay = (row: Order) => {
  //当前订单存入pinia
  shopData.value.currentOrder = row
  PurchaseDialogRef.value.openDialog()
}

//打开二维码弹窗
const openQRDialog = () => {
  //调用子组件打开弹窗
  QRDialogRef.value.openDialog()
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

// 拖拽
.dragClass {
  background: rgba($color: #41c21a, $alpha: 0.5) !important;
}

// 停靠
.ghostClass {
  background: rgba($color: #6cacf5, $alpha: 0.5) !important;
}

// 选择
.chosenClass:hover > td {
  background: rgba($color: #f56c6c, $alpha: 0.5) !important;
}

</style>