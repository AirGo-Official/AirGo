<template>
  <el-dialog v-model="state.isShowDialog" :title="state.title" width="80%" destroy-on-close align-center>
    <el-descriptions
      :column="2"
      border
      title="订单信息"
    >
      <el-descriptions-item label="ID">{{ shopStoreData.currentOrder.value.id }}</el-descriptions-item>
      <el-descriptions-item label="订单ID">{{ shopStoreData.currentOrder.value.out_trade_no }}</el-descriptions-item>
      <el-descriptions-item label="创建时间">{{ DateStrToTime(shopStoreData.currentOrder.value.created_at) }}</el-descriptions-item>
      <el-descriptions-item label="用户ID">{{ shopStoreData.currentOrder.value.user_id}}</el-descriptions-item>
      <el-descriptions-item label="用户">{{shopStoreData.currentOrder.value.user_name }}</el-descriptions-item>
      <el-descriptions-item label="商品ID">{{ shopStoreData.currentOrder.value.goods_id }}</el-descriptions-item>
      <el-descriptions-item label="商品标题">{{ shopStoreData.currentOrder.value.subject }}</el-descriptions-item>
      <el-descriptions-item label="商品类型">{{ shopStoreData.currentOrder.value.goods_type }}</el-descriptions-item>
      <el-descriptions-item label="商品价格">{{ shopStoreData.currentOrder.value.price }}</el-descriptions-item>
    </el-descriptions>
    <el-form v-model="shopStoreData.currentOrder.value" label-width="80" label-position="top">
      <el-form-item label="订单状态">
        <el-radio-group v-model="shopStoreData.currentOrder.value.trade_status" class="ml-4">
          <el-radio :label="constantStore.ORDER_STATUS_CREATED">已创建</el-radio>
          <el-radio :label="constantStore.ORDER_STATUS_WAIT_BUYER_PAY">等待买家付款</el-radio>
          <el-radio :label="constantStore.ORDER_STATUS_TRADE_SUCCESS">支付成功</el-radio>
          <el-radio :label="constantStore.ORDER_STATUS_TRADE_FINISHED">交易结束</el-radio>
          <el-radio :label="constantStore.ORDER_STATUS_TRADE_CLOSED">交易关闭</el-radio>
          <el-radio :label="constantStore.ORDER_STATUS_COMPLETED">已完成</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="发货类型">
        <el-radio-group v-model="shopStoreData.currentOrder.value.deliver_type" class="ml-4">
          <el-radio :label="constantStore.DELIVER_TYPE_NONE">无需发货</el-radio>
          <el-radio :label="constantStore.DELIVER_TYPE_AUTO">自动发货</el-radio>
          <el-radio :label="constantStore.DELIVER_TYPE_MANUAL">手动发货</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="发货内容">
        <v-md-editor v-model="shopStoreData.currentOrder.value.deliver_text"></v-md-editor>
      </el-form-item>
      <el-form-item label="支付ID">
        <el-input v-model="shopStoreData.currentOrder.value.pay_id"></el-input>
      </el-form-item>
      <el-form-item label="支付类型">
        <el-input v-model="shopStoreData.currentOrder.value.pay_type"></el-input>
      </el-form-item>
      <el-form-item label="折扣码ID">
        <el-input v-model="shopStoreData.currentOrder.value.coupon_id"></el-input>
      </el-form-item>
      <el-form-item label="折扣码">
        <el-input v-model="shopStoreData.currentOrder.value.coupon_name"></el-input>
      </el-form-item>
      <el-form-item label="折扣金额">
        <el-input v-model="shopStoreData.currentOrder.value.coupon_amount"></el-input>
      </el-form-item>
      <el-form-item label="余额折扣金额">
        <el-input v-model="shopStoreData.currentOrder.value.balance_amount"></el-input>
      </el-form-item>
      <el-form-item label="订单金额">
        <el-input v-model="shopStoreData.currentOrder.value.total_amount"></el-input>
      </el-form-item>
      <el-form-item label="付款金额">
        <el-input v-model="shopStoreData.currentOrder.value.buyer_logon_id"></el-input>
      </el-form-item>
    </el-form>
    <template #footer>
            <span class="dialog-footer">
                <el-button @click="closeDialog">取消</el-button>
                <el-button type="primary" @click="onSubmit">
                    确认
                </el-button>
            </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">

import {reactive} from "vue";
import {storeToRefs} from "pinia";
import { useAdminShopStore } from "/@/stores/admin_logic/shopStore";
import { useConstantStore } from "/@/stores/constantStore";
import {DateStrToTime} from "/@/utils/formatTime";

const shopStore = useAdminShopStore()
const shopStoreData = storeToRefs(shopStore)
const constantStore = useConstantStore()

const state = reactive({
  isShowDialog: false,
  title: "订单编辑",
})


// 打开弹窗
const openDialog = (row: Order) => {
  state.isShowDialog = true
  shopStoreData.currentOrder.value = row
}
// 关闭弹窗
const closeDialog = () => {
  state.isShowDialog = false
};

//确认提交
const onSubmit=() =>{
  shopStore.updateUserOrder(shopStoreData.currentOrder.value).then((res)=>{
    closeDialog()
  })
}

// 暴露变量
defineExpose({
  openDialog,   // 打开弹窗
});

</script>

<style scoped lang="scss">

</style>