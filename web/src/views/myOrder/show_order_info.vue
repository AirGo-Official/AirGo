<template>
  <el-dialog v-model="state.isShowDialog" :title="state.title" width="80%" destroy-on-close align-center>
   <div>
     <el-descriptions
         column="1"
         border
     >
       <el-descriptions-item label="订单编号">{{ shopStoreData.currentOrder.value.out_trade_no }}</el-descriptions-item>
       <el-descriptions-item label="创建时间">{{ DateStrToTime(shopStoreData.currentOrder.value.created_at)
         }}</el-descriptions-item>
       <el-descriptions-item label="用户">{{ shopStoreData.currentOrder.value.user_name }}</el-descriptions-item>
       <el-descriptions-item label="商品ID">{{ shopStoreData.currentOrder.value.goods_id }}</el-descriptions-item>
       <el-descriptions-item label="商品类型">
         <el-tag class="ml-2" v-if="shopStoreData.currentOrder.value.goods_type === 'subscribe'">订阅</el-tag>
         <el-tag class="ml-2" v-if="shopStoreData.currentOrder.value.goods_type === 'recharge'">充值</el-tag>
         <el-tag class="ml-2" v-if="shopStoreData.currentOrder.value.goods_type === 'general'">普通商品</el-tag>
        </el-descriptions-item>
       <el-descriptions-item label="商品标题">{{ shopStoreData.currentOrder.value.subject }}</el-descriptions-item>
       <el-descriptions-item label="商品价格">{{ shopStoreData.currentOrder.value.price }}</el-descriptions-item>
       <el-descriptions-item label="订单金额">{{ shopStoreData.currentOrder.value.total_amount }}</el-descriptions-item>
       <el-descriptions-item v-if="shopStoreData.currentOrder.value.deliver_type !== 'none'" label="发货内容">
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

const shopStore = useShopStore()
const shopStoreData = storeToRefs(shopStore)

const state = reactive({
  isShowDialog: false,
  title: "订单详情",
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