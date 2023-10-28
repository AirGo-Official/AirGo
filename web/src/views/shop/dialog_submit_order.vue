<template>
  <div>
    <el-dialog v-model="state.isShowSubmitOrderDialog" title="订单详情" width="80%">
      <div class="home-card-item">
        <div class="card-text">
          <el-button class="card-text-left" type="warning">订购套餐</el-button>
          <el-text class="card-header-left">{{ shopData.currentGoods.subject }}</el-text>
        </div>
        <div class="card-text">
          <el-button class="card-text-left" type="warning">套餐流量</el-button>
          <span class="card-text-right">{{ shopData.currentGoods.total_bandwidth }}GB</span>
        </div>
        <div class="card-text">
          <el-button class="card-text-left" type="warning">有效期</el-button>
          <span class="card-text-right">{{ shopData.currentGoods.expiration_date }}天</span>
        </div>
        <div v-html="shopData.currentGoods.des"></div>
        <el-divider></el-divider>
        <div class="card-text">
          <el-input v-model="shopData.currentOrder.coupon_name" placeholder="输入折扣码"></el-input>
          <el-button class="card-text-right" color="blue" size="default" @click="varifyCoupon">验证</el-button>
        </div>
        <div class="card-text">
          <el-button class="card-text-left" type="info">金额</el-button>
          <el-text class="card-text-right">{{ shopData.currentOrder.price }}</el-text>
        </div>
        <div class="card-text">
          <el-button class="card-text-left" type="info">优惠码折扣</el-button>
          <el-text class="card-text-right">-{{ shopData.currentOrder.coupon_amount }}</el-text>
        </div>
        <div class="card-text">
          <el-button class="card-text-left" type="info">旧套餐抵扣</el-button>
          <el-text class="card-text-right">-{{ shopData.currentOrder.deduction_amount }}</el-text>
        </div>
        <div class="card-text">
          <el-button class="card-text-left" type="info">余额抵扣</el-button>
          <el-text class="card-text-right">-{{ shopData.currentOrder.remain_amount }}</el-text>
        </div>
        <div class="card-text">
          <el-text class="card-text-left" style="font-size: 25px;">应付</el-text>
          <el-text class="card-text-right" style="font-size: 25px;">{{ shopData.currentOrder.total_amount }}</el-text>
        </div>
      </div>
      <template #footer>
            <span class="dialog-footer">
                <el-button @click="state.isShowSubmitOrderDialog = false">取消</el-button>
                <el-button type="primary" @click="onSubmitOrder" color="#FC3D08">
                    提交订单
                </el-button>
            </span>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {reactive} from 'vue';
import {ElMessage} from "element-plus";
import {storeToRefs} from 'pinia';
import {useShopStore} from "/@/stores/shopStore";
import {useOrderStore} from "/@/stores/orderStore"
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";

const shopStore = useShopStore()
const {shopData} = storeToRefs(shopStore)
const orderStore = useOrderStore()
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)

//定义变量
const state = reactive({
  isShowSubmitOrderDialog: false,

})
//获取订单详情
const getOrderInfo = (params: object) => {
  orderStore.getOrderInfo(params).then((res) => {
    ElMessage.success(res.msg)
    shopData.value.currentOrder = {} as Order
    shopData.value.currentOrder = res.data
  }).catch()
}
//打开弹窗
const openDialog = () => {
  state.isShowSubmitOrderDialog = true
  //获取订单详情（计算价格等）
  getOrderInfo({goods_id: shopData.value.currentGoods.id})
}
//验证折扣码
const varifyCoupon = () => {
  getOrderInfo({goods_id: shopData.value.currentOrder.goods_id, coupon_name: shopData.value.currentOrder.coupon_name})
}

//提交订单按钮
const onSubmitOrder = () => {
  shopData.value.currentOrder.id = 0
  request(apiStoreData.api.value.shop_preCreatePay, shopData.value.currentOrder).then((res) => {
    //保存订单信息到pinia
    shopData.value.currentOrder = res.data
    //关闭弹窗
    state.isShowSubmitOrderDialog = false
    //调用父组件 支付弹窗
    emits('openPurchaseDialog')
  }).catch(() => {
    state.isShowSubmitOrderDialog = false
  })
}
//子组件调用父组件
const emits = defineEmits(['openPurchaseDialog'])
//暴露变量
defineExpose({
  openDialog,
})
</script>

<style scoped>
.el-card {
  background-image: url("../../assets/bgc/bg-3.svg");
  background-repeat: no-repeat;
  background-position: 100%, 100%;
}

.card-text {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 35px
}

.card-text-left {
  margin-top: auto;
  margin-bottom: auto;
}

.card-text-right {
  margin-top: auto;
  margin-bottom: auto;
  font-size: 20px;
}

.card-header-left {
  font-size: 30px;
  color: #FC3D08;
}
</style>