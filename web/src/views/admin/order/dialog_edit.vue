<template>
  <el-dialog v-model="state.isShowDialog" :title="$t('message.adminOrder.modify_order')" width="80%" destroy-on-close align-center>
    <el-descriptions
      :column="2"
      border
    >
      <el-descriptions-item :label="$t('message.adminOrder.Order.id')">{{ shopStoreData.currentOrder.value.id }}</el-descriptions-item>
      <el-descriptions-item :label="$t('message.adminOrder.Order.out_trade_no')">{{ shopStoreData.currentOrder.value.out_trade_no }}</el-descriptions-item>
      <el-descriptions-item :label="$t('message.adminOrder.Order.created_at')">{{ DateStrToTime(shopStoreData.currentOrder.value.created_at) }}</el-descriptions-item>
      <el-descriptions-item :label="$t('message.adminOrder.Order.user_id')">{{ shopStoreData.currentOrder.value.user_id}}</el-descriptions-item>
      <el-descriptions-item :label="$t('message.adminOrder.Order.user_name')">{{shopStoreData.currentOrder.value.user_name }}</el-descriptions-item>
      <el-descriptions-item :label="$t('message.adminOrder.Order.goods_id')">{{ shopStoreData.currentOrder.value.goods_id }}</el-descriptions-item>
      <el-descriptions-item :label="$t('message.adminOrder.Order.subject')">{{ shopStoreData.currentOrder.value.subject }}</el-descriptions-item>
      <el-descriptions-item :label="$t('message.adminOrder.Order.goods_type')">{{ shopStoreData.currentOrder.value.goods_type }}</el-descriptions-item>
      <el-descriptions-item :label="$t('message.adminOrder.Order.price')">{{ shopStoreData.currentOrder.value.price }}</el-descriptions-item>
    </el-descriptions>
    <el-form v-model="shopStoreData.currentOrder.value" label-position="top">
      <el-form-item :label="$t('message.adminOrder.Order.trade_status')">
        <el-radio-group v-model="shopStoreData.currentOrder.value.trade_status" class="ml-4">
          <el-radio :label="constantStore.ORDER_STATUS_CREATED">{{$t('message.constant.ORDER_STATUS_CREATED')}}</el-radio>
          <el-radio :label="constantStore.ORDER_STATUS_WAIT_BUYER_PAY">{{$t('message.constant.ORDER_STATUS_WAIT_BUYER_PAY')}}</el-radio>
          <el-radio :label="constantStore.ORDER_STATUS_TRADE_SUCCESS">{{$t('message.constant.ORDER_STATUS_TRADE_SUCCESS')}}</el-radio>
          <el-radio :label="constantStore.ORDER_STATUS_TRADE_FINISHED">{{$t('message.constant.ORDER_STATUS_TRADE_FINISHED')}}</el-radio>
          <el-radio :label="constantStore.ORDER_STATUS_TRADE_CLOSED">{{$t('message.constant.ORDER_STATUS_TRADE_CLOSED')}}</el-radio>
          <el-radio :label="constantStore.ORDER_STATUS_COMPLETED">{{$t('message.constant.ORDER_STATUS_COMPLETED')}}</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item :label="$t('message.adminOrder.Order.deliver_type')">
        <el-radio-group v-model="shopStoreData.currentOrder.value.deliver_type" class="ml-4">
          <el-radio :label="constantStore.DELIVER_TYPE_NONE">{{$t('message.constant.DELIVER_TYPE_NONE')}}</el-radio>
          <el-radio :label="constantStore.DELIVER_TYPE_AUTO">{{$t('message.constant.DELIVER_TYPE_AUTO')}}</el-radio>
          <el-radio :label="constantStore.DELIVER_TYPE_MANUAL">{{$t('message.constant.DELIVER_TYPE_MANUAL')}}</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item :label="$t('message.adminOrder.Order.deliver_text')" v-if="shopStoreData.currentOrder.value.deliver_type !== constantStore.DELIVER_TYPE_NONE">
        <v-md-editor v-model="shopStoreData.currentOrder.value.deliver_text"></v-md-editor>
      </el-form-item>
      <el-form-item :label="$t('message.adminOrder.Order.pay_id')">
        <el-input v-model="shopStoreData.currentOrder.value.pay_id"></el-input>
      </el-form-item>
      <el-form-item :label="$t('message.adminOrder.Order.pay_type')">
        <el-input v-model="shopStoreData.currentOrder.value.pay_type"></el-input>
      </el-form-item>
      <el-form-item :label="$t('message.adminOrder.Order.coupon_id')">
        <el-input v-model="shopStoreData.currentOrder.value.coupon_id"></el-input>
      </el-form-item>
      <el-form-item :label="$t('message.adminOrder.Order.coupon_name')">
        <el-input v-model="shopStoreData.currentOrder.value.coupon_name"></el-input>
      </el-form-item>
      <el-form-item :label="$t('message.adminOrder.Order.coupon_amount')">
        <el-input v-model="shopStoreData.currentOrder.value.coupon_amount"></el-input>
      </el-form-item>
      <el-form-item :label="$t('message.adminOrder.Order.balance_amount')">
        <el-input v-model="shopStoreData.currentOrder.value.balance_amount"></el-input>
      </el-form-item>
      <el-form-item :label="$t('message.adminOrder.Order.total_amount')">
        <el-input v-model="shopStoreData.currentOrder.value.total_amount"></el-input>
      </el-form-item>
      <el-form-item :label="$t('message.adminOrder.Order.buyer_pay_amount')">
        <el-input v-model="shopStoreData.currentOrder.value.buyer_pay_amount"></el-input>
      </el-form-item>
      <el-form-item :label="$t('message.adminOrder.Order.buyer_logon_id')">
        <el-input v-model="shopStoreData.currentOrder.value.buyer_logon_id"></el-input>
      </el-form-item>
    </el-form>
    <template #footer>
            <span class="dialog-footer">
                <el-button @click="closeDialog">{{$t('message.common.button_cancel')}}</el-button>
                <el-button type="primary" @click="onSubmit">
                    {{$t('message.common.button_confirm')}}
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