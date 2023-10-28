<template>
  <div class="container layout-padding">
    <el-card shadow="hover" class="layout-padding-auto">
      <el-button type="info" plain>只显示最近10次订单</el-button>
      <el-table :data="orderPersonal.allOrders.order_list" stripe fit height="100%" style="width: 100%;">
        <el-table-column prop="subject" label="商品标题" show-overflow-tooltip width="150px"/>
        <el-table-column prop="total_amount" label="金额" show-overflow-tooltip width="60px"/>
        <el-table-column prop="trade_status" label="状态" show-overflow-tooltip width="90px">
          <template #default="scope">
            <el-tag type="success" v-if="scope.row.trade_status==='TRADE_SUCCESS'">支付成功</el-tag>
            <el-tag type="warning" v-else-if="scope.row.trade_status==='WAIT_BUYER_PAY'">等待付款</el-tag>
            <el-tag type="danger" v-else-if="scope.row.trade_status==='TRADE_CLOSED'">超时关闭</el-tag>
            <el-tag type="success" v-else-if="scope.row.trade_status==='TRADE_FINISHED'">交易结束</el-tag>
            <el-tag type="info" v-else-if="scope.row.trade_status==='Created'">已创建</el-tag>
            <el-tag type="success" v-else-if="scope.row.trade_status==='Completed'">已完成</el-tag>
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
    </el-card>
    <!-- 引入确认购买弹窗组件 -->
    <PurchaseDialog ref="PurchaseDialogRef" @openQRDialog="openQRDialog"></PurchaseDialog>
    <!-- 引入二维码弹窗 -->
    <QRDialog ref="QRDialogRef"></QRDialog>
  </div>
</template>

<script setup lang="ts">
import {storeToRefs} from "pinia";
import {defineAsyncComponent, onMounted, ref} from "vue";
import {useOrderStore} from "/@/stores/orderStore";
import {useShopStore} from "/@/stores/shopStore";
import {isMobile} from "/@/utils/other";

const orderStore = useOrderStore()
const {orderPersonal} = storeToRefs(orderStore)
const shopStore = useShopStore()
const {shopData} = storeToRefs(shopStore)
const PurchaseDialog = defineAsyncComponent(() => import('/@/views/shop/dialog_purchase.vue'))
const QRDialog = defineAsyncComponent(() => import('/@/views/shop/dialog_QR.vue'))
const PurchaseDialogRef = ref()
const QRDialogRef = ref()

onMounted(() => {
  orderStore.getOrder() //获取用户最近10次订单
})

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