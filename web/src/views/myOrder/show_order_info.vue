<template>
  <el-dialog v-model="state.isShowDialog" :title="$t('message.common.details')" width="80%" destroy-on-close align-center>
   <div>
     <el-descriptions
         column="1"
         border
     >
       <el-descriptions-item :label="$t('message.adminOrder.Order.out_trade_no')">{{ shopStoreData.currentOrder.value.out_trade_no }}</el-descriptions-item>
       <el-descriptions-item :label="$t('message.adminOrder.Order.created_at')">{{ DateStrToTime(shopStoreData.currentOrder.value.created_at) }}</el-descriptions-item>
       <el-descriptions-item :label="$t('message.adminOrder.Order.user_name')">{{ shopStoreData.currentOrder.value.user_name }}</el-descriptions-item>
       <el-descriptions-item :label="$t('message.adminOrder.Order.goods_id')">{{ shopStoreData.currentOrder.value.goods_id }}</el-descriptions-item>
       <el-descriptions-item :label="$t('message.adminOrder.Order.goods_type')">
         <el-tag class="ml-2" v-if="shopStoreData.currentOrder.value.goods_type === constantStore.GOODS_TYPE_SUBSCRIBE">{{$t('message.constant.GOODS_TYPE_SUBSCRIBE')}}</el-tag>
         <el-tag class="ml-2" v-if="shopStoreData.currentOrder.value.goods_type === constantStore.GOODS_TYPE_RECHARGE">{{$t('message.constant.GOODS_TYPE_RECHARGE')}}</el-tag>
         <el-tag class="ml-2" v-if="shopStoreData.currentOrder.value.goods_type === constantStore.GOODS_TYPE_GENERAL">{{$t('message.constant.GOODS_TYPE_GENERAL')}}</el-tag>
        </el-descriptions-item>
       <el-descriptions-item :label="$t('message.adminOrder.Order.subject')">{{ shopStoreData.currentOrder.value.subject }}</el-descriptions-item>
       <el-descriptions-item :label="$t('message.adminOrder.Order.price')">{{ shopStoreData.currentOrder.value.price }}</el-descriptions-item>
       <el-descriptions-item :label="$t('message.adminOrder.Order.total_amount')">{{ shopStoreData.currentOrder.value.total_amount }}</el-descriptions-item>
       <el-descriptions-item v-if="shopStoreData.currentOrder.value.deliver_type !== $t('message.constant.DELIVER_TYPE_NONE')"
                             :label="$t('message.adminOrder.Order.deliver_text')">
         <v-md-preview :text="shopStoreData.currentOrder.value.deliver_text"></v-md-preview></el-descriptions-item>
     </el-descriptions>
   </div>
  </el-dialog>
</template>

<script setup lang="ts">

import {reactive} from "vue";
import {storeToRefs} from "pinia";
import {DateStrToTime} from "/@/utils/formatTime"
import { useShopStore } from "/@/stores/user_logic/shopStore";
import { useConstantStore } from "/@/stores/constantStore";

const shopStore = useShopStore()
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

}

// 暴露变量
defineExpose({
  openDialog,   // 打开弹窗
});

</script>

<style scoped lang="scss">

</style>