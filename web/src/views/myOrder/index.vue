<template>
  <div class="container layout-padding">
    <el-card shadow="hover" class="layout-padding-auto">
      <el-table :data="shopStoreData.orderList.value.data" stripe fit height="100%" style="width: 100%;" @sort-change="sortChange">
<!--        <el-table-column type="index" :label="$t('message.adminOrder.Order.index')" fixed width="60px"/>-->
        <el-table-column prop="created_at" :label="$t('message.adminOrder.Order.created_at')" width="150" sortable="custom" fixed>
          <template #default="{row}">
            <el-tag type="info">{{ DateStrToTime(row.created_at) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="out_trade_no" :label="$t('message.adminOrder.Order.out_trade_no')" width="150" sortable="custom"/>
        <el-table-column prop="order_type" :label="$t('message.adminOrder.Order.order_type')" width="100" sortable="custom">
          <template #default="{row}">
            <el-button type="primary" v-if="row.order_type === constantStore.ORDER_TYPE_NEW">{{$t('message.constant.ORDER_TYPE_NEW')}}</el-button>
            <el-button type="primary" v-else-if="row.order_type === constantStore.ORDER_TYPE_RENEW">{{$t('message.constant.ORDER_TYPE_RENEW')}}</el-button>
            <el-button type="info" v-else>{{$t('message.constant.ORDER_TYPE_DESTROYED')}}</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="subject" :label="$t('message.adminOrder.Order.subject')" show-overflow-tooltip width="200" sortable="custom"/>
        <el-table-column prop="total_amount" :label="$t('message.adminOrder.Order.total_amount')" show-overflow-tooltip width="150" sortable="custom"/>
        <el-table-column prop="trade_status" :label="$t('message.adminOrder.Order.trade_status')" show-overflow-tooltip sortable="custom" width="150">
          <template #default="{row}">
            <el-tag type="success" v-if="row.trade_status===constantStore.ORDER_STATUS_TRADE_SUCCESS">{{$t('message.constant.ORDER_STATUS_TRADE_SUCCESS')}}</el-tag>
            <el-tag type="warning" v-else-if="row.trade_status===constantStore.ORDER_STATUS_WAIT_BUYER_PAY">{{$t('message.constant.ORDER_STATUS_WAIT_BUYER_PAY')}}</el-tag>
            <el-tag type="info" v-else-if="row.trade_status===constantStore.ORDER_STATUS_TRADE_CLOSED">{{$t('message.constant.ORDER_STATUS_TRADE_CLOSED')}}</el-tag>
            <el-tag type="info" v-else-if="row.trade_status===constantStore.ORDER_STATUS_TRADE_FINISHED">{{$t('message.constant.ORDER_STATUS_TRADE_FINISHED')}}</el-tag>
            <el-tag type="info" v-else-if="row.trade_status===constantStore.ORDER_STATUS_CREATED">{{$t('message.constant.ORDER_STATUS_CREATED')}}</el-tag>
            <el-tag type="success" v-else-if="row.trade_status===constantStore.ORDER_STATUS_COMPLETED">{{$t('message.constant.ORDER_STATUS_COMPLETED')}}</el-tag>
            <el-tag type="danger" v-else>{{$t('message.constant.ORDER_STATUS_UNKNOWN_STATE')}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('message.common.operate')" min-width="400">
          <template #default="{row}">
            <el-button  text type="primary"
                       @click="showOrderInfo(row)">{{$t('message.common.details')}}
            </el-button>
            <el-button v-if="(row.trade_status === constantStore.ORDER_STATUS_WAIT_BUYER_PAY
                       || row.trade_status === constantStore.ORDER_STATUS_CREATED)
                       && (getTimeDifference(row.created_at) < constantStore.CACHE_SUBMIT_ORDER_TIMEOUT * 60*1000)"

                    text type="primary"
                       @click="toPay(row)">
              <span>{{$t('message.adminShop.purchase')}}<el-countdown class="countdown" :value="orderTimeout(row.created_at)" /></span>
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
    <!-- 引入确认购买弹窗组件 -->
    <PurchaseDialog ref="PurchaseDialogRef" @openQRDialog="openQRDialog"></PurchaseDialog>
    <!-- 引入二维码弹窗 -->
    <QRDialog ref="QRDialogRef"></QRDialog>
    <!--    订单详情弹窗-->
    <DialogShowOrderInfo ref="DialogShowOrderInfoRef"></DialogShowOrderInfo>
  </div>
</template>

<script setup lang="ts">
import {storeToRefs} from "pinia";
import { defineAsyncComponent, onBeforeMount, onMounted, reactive, ref } from "vue";
import {useShopStore} from "/@/stores/user_logic/shopStore";
import { DateStrToTime, getTimeDifference } from "/@/utils/formatTime";
import { useConstantStore } from "/@/stores/constantStore";

const shopStore = useShopStore()
const shopStoreData = storeToRefs(shopStore)
const PurchaseDialog = defineAsyncComponent(() => import('/@/views/shop/purchase.vue'))
const QRDialog = defineAsyncComponent(() => import('/@/views/shop/purchase.vue'))
const PurchaseDialogRef = ref()
const QRDialogRef = ref()
const constantStore = useConstantStore()


const DialogShowOrderInfo = defineAsyncComponent(() => import('/@/views/myOrder/show_order_info.vue'))
const DialogShowOrderInfoRef = ref()

const state = reactive({
  queryParams:{
    table_name: 'order',
    field_params_list: [
      { field: 'out_trade_no', field_chinese_name: '', field_type: '', condition: 'like', condition_value: '', operator: '',
      }
    ] as FieldParams[],
    pagination: { page_num: 1, page_size: 30, order_by: 'id DESC',
    } as Pagination,//分页参数
  },
})

//
const getUserOrders = () => {
  shopStore.getOrderList(state.queryParams)
}
// 分页改变
const onHandleSizeChange = (val: number) => {
  state.queryParams.pagination.page_size = val;
  getUserOrders()
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
  state.queryParams.pagination.page_num = val;
  getUserOrders()
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
  getUserOrders()
}

onMounted(() => {
  getUserOrders()
});
//
const showOrderInfo=(row:Order)=>{
  DialogShowOrderInfoRef.value.openDialog(row)

}

//去支付流程
const toPay = (row: Order) => {
  //当前订单存入pinia
  shopStoreData.currentOrder.value = row
  PurchaseDialogRef.value.openDialog('fromMyOrder')
}

//打开二维码弹窗
const openQRDialog = () => {
  //调用子组件打开弹窗
  QRDialogRef.value.openDialog()
}
//订单支付倒计时
const orderTimeout=(strTime:string)=>{
 return  new Date(strTime).getTime()+constantStore.CACHE_SUBMIT_ORDER_TIMEOUT*1000 * 60
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
.countdown {
  :deep(.el-statistic__content ) {
    color: red;
    font-size: 12px;
  }
}

</style>