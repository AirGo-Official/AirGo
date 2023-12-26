<template>
  <el-dialog v-model="state.isShowDialog" :title="state.title" width="80%" destroy-on-close align-center>
    <el-form v-model="orderStoreData.orderManageData.value.currentOrder" label-width="80" label-position="top">
      <el-form-item label="订单ID">
        <el-input v-model="orderStoreData.orderManageData.value.currentOrder.out_trade_no"></el-input>
      </el-form-item>
      <el-form-item label="ID">
        <el-input v-model.number="orderStoreData.orderManageData.value.currentOrder.id"></el-input>
      </el-form-item>
      <el-form-item label="订单金额">
        <el-input v-model="orderStoreData.orderManageData.value.currentOrder.total_amount"></el-input>
      </el-form-item>
      <el-form-item label="实收金额">
        <el-input v-model="orderStoreData.orderManageData.value.currentOrder.receipt_amount"></el-input>
      </el-form-item>
      <el-form-item label="付款金额">
        <el-input v-model="orderStoreData.orderManageData.value.currentOrder.buyer_logon_id"></el-input>
      </el-form-item>


      <el-form-item label="用户">
        <el-input v-model="orderStoreData.orderManageData.value.currentOrder.user_name"></el-input>
      </el-form-item>
      <el-form-item label="用户ID">
        <el-input v-model.number="orderStoreData.orderManageData.value.currentOrder.user_id"></el-input>
      </el-form-item>

      <el-form-item label="商品ID">
        <el-input v-model="orderStoreData.orderManageData.value.currentOrder.goods_id"></el-input>
      </el-form-item>
      <el-form-item label="商品标题">
        <el-input v-model="orderStoreData.orderManageData.value.currentOrder.subject"></el-input>
      </el-form-item>
      <el-form-item label="商品类型">
        <el-input v-model="orderStoreData.orderManageData.value.currentOrder.goods_type"></el-input>
      </el-form-item>
      <el-form-item label="商品价格">
        <el-input v-model="orderStoreData.orderManageData.value.currentOrder.price"></el-input>
      </el-form-item>
      <el-form-item label="发货类型">
        <el-input v-model="orderStoreData.orderManageData.value.currentOrder.deliver_type"></el-input>
      </el-form-item>
      <el-form-item label="发货内容">
        <v-md-editor v-model="orderStoreData.orderManageData.value.currentOrder.deliver_text"></v-md-editor>
      </el-form-item>

      <el-form-item label="支付ID">
        <el-input v-model="orderStoreData.orderManageData.value.currentOrder.pay_id"></el-input>
      </el-form-item>
      <el-form-item label="支付类型">
        <el-input v-model="orderStoreData.orderManageData.value.currentOrder.pay_type"></el-input>
      </el-form-item>

      <el-form-item label="折扣码ID">
        <el-input v-model="orderStoreData.orderManageData.value.currentOrder.coupon_id"></el-input>
      </el-form-item>
      <el-form-item label="折扣码">
        <el-input v-model="orderStoreData.orderManageData.value.currentOrder.coupon_name"></el-input>
      </el-form-item>
      <el-form-item label="折扣金额">
        <el-input v-model="orderStoreData.orderManageData.value.currentOrder.coupon_amount"></el-input>
      </el-form-item>

      <el-form-item label="旧套餐折扣金额">
        <el-input v-model="orderStoreData.orderManageData.value.currentOrder.deduction_amount"></el-input>
      </el-form-item>
      <el-form-item label="余额折扣金额">
        <el-input v-model="orderStoreData.orderManageData.value.currentOrder.remain_amount"></el-input>
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
import {useOrderStore} from "/@/stores/orderStore";
import {storeToRefs} from "pinia";

const orderStore = useOrderStore()
const orderStoreData = storeToRefs(orderStore)

const state = reactive({
  isShowDialog: false,
  title: "订单编辑",
})


// 打开弹窗
const openDialog = (row: Order) => {
  state.isShowDialog = true
  orderStoreData.orderManageData.value.currentOrder = row
}
// 关闭弹窗
const closeDialog = () => {
  state.isShowDialog = false
};

//确认提交
const onSubmit=() =>{
  orderStore.updateUserOrder(orderStoreData.orderManageData.value.currentOrder).then((res)=>{
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