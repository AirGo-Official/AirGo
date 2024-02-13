<template>
<div v-loading="state.isShowLoading">
  <el-dialog v-model="state.isShowDialog" title="详情" width="80%" destroy-on-close>
    <el-steps :active="state.active" process-status="wait" finish-status="success">
      <el-step title="选择套餐"></el-step>
      <el-step title="订单详情"></el-step>
      <el-step title="支付"></el-step>
    </el-steps>

    <div v-if="state.active === 1" v-loading="state.isLoading">
      <div class="home-card-item" v-if="shopStoreData.currentOrder.value.order_type === constantStore.ORDER_TYPE_NEW">
        <div class="card-text">
          <el-button class="card-text-left" color="blue">订购套餐</el-button>
          <el-text class="card-header-left">{{ shopStoreData.currentGoods.value.subject }}</el-text>
        </div>
        <div class="card-text">
          <el-button class="card-text-left" color="blue">套餐流量</el-button>
          <span class="card-text-right">{{ shopStoreData.currentGoods.value.total_bandwidth }}GB</span>
        </div>
        <div class="card-text">
          <el-button class="card-text-left" color="blue">详情</el-button>
        </div>
        <div v-html="shopStoreData.currentGoods.value.des"></div>
        <div class="card-text">
          <el-button class="card-text-left" color="blue">订购时长</el-button>
        </div>
        <div class="slider-demo-block">
          <el-slider v-model.number="shopStoreData.currentOrder.value.duration" @change="getOrderInfo" :max="24" :min="1" show-input />
        </div>
        <el-divider></el-divider>
        <div class="card-text">
          <el-input v-model="shopStoreData.currentOrder.value.coupon_name" placeholder="输入折扣码" size="default">
            <template #prepend>
              <el-icon><Ticket /></el-icon>
            </template>
            <template #append>
              <el-button class="card-text-right" color="blue" size="small" @click="getOrderInfo">验证</el-button>
            </template>
          </el-input>
        </div>
        <div class="card-text">
          <el-button class="card-text-left" type="info">金额</el-button>
          <el-text class="card-text-right">{{ shopStoreData.currentOrder.value.price }}</el-text>
        </div>
        <div class="card-text">
          <el-button class="card-text-left" type="info">优惠码折扣</el-button>
          <el-text class="card-text-right">-{{ shopStoreData.currentOrder.value.coupon_amount }}</el-text>
        </div>
        <div class="card-text">
          <el-button class="card-text-left" type="info">应付</el-button>
          <el-text class="card-text-right" style="font-size: 25px;">{{ shopStoreData.currentOrder.value.total_amount }}</el-text>
        </div>

        <div style="text-align: right;margin-top: 20px" >
          <el-button color="blue" @click="nextSubmitOrder">提交订单</el-button>
        </div>
      </div>
      <div class="home-card-item" v-else>
        <div class="card-text">
          <el-button class="card-text-left" color="blue">续费套餐</el-button>
          <el-text class="card-header-left">{{ shopStoreData.currentOrder.value.subject }}</el-text>
        </div>
        <div class="card-text">
          <el-button class="card-text-left" type="info">应付</el-button>
          <el-text class="card-text-right" style="font-size: 25px;">{{ shopStoreData.currentOrder.value.total_amount }}</el-text>
        </div>
        <div style="text-align: right;margin-top: 20px" >
          <el-button color="blue" @click="nextSubmitOrder">提交订单</el-button>
        </div>
      </div>
    </div>
    <div v-if="state.active === 2">
      <div class="home-card-item">
        <div class="card-text">
          <el-button type="primary">套餐</el-button>
          <el-text class="card-text-right">{{ shopStoreData.currentOrder.value.subject }}</el-text>
        </div>
        <div class="card-text">
          <el-button type="primary">金额</el-button>
          <el-text class="card-text-right">{{ shopStoreData.currentOrder.value.total_amount }}元</el-text>
        </div>
        <el-divider></el-divider>
        <div v-if="state.isShowPayment">
          <el-button type="primary">支付方式</el-button>
          <div>
            <el-radio-group style="height: 60px;" v-model="shopStoreData.currentOrder.value.pay_id" class="ml-4">
              <el-radio :label="v.id" v-for="(v,k) in shopStoreData.payList.value" :key="k">
                <div style="display: flex;align-items: center">
                  <el-image :src="v.pay_logo_url" style="height: 15px;"></el-image>
                  <span>{{ v.name }}</span>
                </div>
              </el-radio>
            </el-radio-group>
          </div>
        </div>
      </div>
      <div style="text-align: right;margin-top: 20px" >
        <el-button color="blue" @click="nextPurchase">确认购买</el-button>
      </div>
    </div>
    <div v-if="state.active === 3">
      <div v-if="state.isShowQR">
        <el-card shadow="hover">
          <div class="qrcode-img-warp">
            <div class="mb30 mt30 qrcode-img">
              <!-- 二维码弹窗 -->
              <div id="qrcode" class="qrcode" ref="qrcodeRef"></div>
            </div>
          </div>
        </el-card>
      </div>
    </div>

  </el-dialog>

</div>
</template>

<script setup lang="ts">
import { nextTick, reactive, ref } from "vue";
import { useShopStore } from "/@/stores/user_logic/shopStore";
import { storeToRefs } from "pinia";
import { useApiStore } from "/@/stores/apiStore";
import { request } from "/@/utils/request";

import { ElMessage, ElMessageBox } from "element-plus";
import { isMobile } from "/@/utils/other";
import qs from "qs";
import QRCode from "qrcodejs2-fixes";
import { useConstantStore } from "/@/stores/constantStore";

const shopStore = useShopStore()
const shopStoreData = storeToRefs(shopStore)
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)

const constantStore = useConstantStore()
const qrcodeRef = ref();

const state = reactive({
  active:1,
  isShowDialog: false,
  isShowQR: false,
  isShowPayment: false,
  isLoading:false,
  QRcode: null,
  isShowLoading:false,
})
const openDialog = () => {
  state.active = 1
  state.isShowDialog = true
}
const closeDialog = ()=>{
  state.active = 1
  state.isShowDialog = false
  state.isShowQR = false
  state.isShowPayment = false
  state.isLoading = false
}
const getOrderInfo = () => {
  state.isLoading = true
  shopStore.getOrderInfo(shopStoreData.currentOrder.value).then((res) => {
    state.isShowLoading = false
    shopStoreData.currentOrder.value = {} as Order
    shopStoreData.currentOrder.value = res.data
  }).catch(()=>{
    state.isShowLoading = false
  })
}
const next = () => {
  if (state.active === 3){

  } else {
    state.active++
  }
}

function loop () {
  let i=0;
  let  timer = setInterval(() => {
    getOrderInfoWaitPay(timer,i++)
  }, 3000)
}
function getOrderInfoWaitPay (timer: NodeJS.Timeout, i:number) {
  setTimeout(()=>{
    console.log("开始轮循请求,次数：",i);
    //请求
    request(apiStoreData.userApi.value.getOrderInfoWaitPay, shopStoreData.currentOrder.value).then((res) => {
      //保存订单信息到pinia
      shopStoreData.currentOrder.value = res.data
      if (shopStoreData.currentOrder.value.trade_status === constantStore.ORDER_STATUS_WAIT_BUYER_PAY){
        // 关闭轮询
        clearInterval(timer);
        // 获取支付列表
        shopStore.getEnabledPayList()
        // 显示支付
        state.isShowPayment = true
        // 关闭loading
        state.isLoading = false
      }
    }).catch(() => {
    })
    if(i>=6){
      clearInterval(timer);
      ElMessageBox.alert('请求超时，请稍后检查我的订单，有无订单记录', '提示', {
        confirmButtonText: 'OK',
      })
        .then(() => {
          closeDialog()
        })
        .catch(() => {
        });

    }
  }, 0)
}

const nextSubmitOrder=()=>{
  // 加载loading
  state.isLoading = true
  //
  shopStoreData.currentOrder.value.id = 0
  // 轮训
  request(apiStoreData.userApi.value.preCreatePay, shopStoreData.currentOrder.value).then((res) => {
    //保存订单信息到pinia
    shopStoreData.currentOrder.value = res.data
    //
    next()
    //
    loop()
  }).catch(error => {

  })
}
const nextPurchase=()=>{
  request(apiStoreData.userApi.value.purchase, shopStoreData.currentOrder.value).then((res) => {
    if (res.msg === "Purchase success") { //0元购
      ElMessage.success(res.msg)
      // state.isShowPurchaseDialog = false
      return
    } else {
      //保存支付信息
      shopStoreData.currentOrder.value.pay_info = res.data
      let pay_info = shopStoreData.currentOrder.value.pay_info
      if (pay_info.alipay_info.qr_code !== "") { //支付宝支付
        if (isMobile()) {
          window.location.href = pay_info.alipay_info.qr_code; //手机端跳转支付页面
          return
        } else {
          shoeQR()  //电脑端打开支付二维码弹窗
        }
      } else if (pay_info.epay_info.epay_api_url !== "") { //易支付
        // 对象转url参数
        let params = qs.stringify(pay_info.epay_info.epay_pre_create_pay)
        window.location.href = pay_info.epay_info.epay_api_url + "?" + params
      }
    }
  }).catch(() => {
  })
  //关闭弹窗
  // state.isShowPurchaseDialog = false

}
//
const shoeQR=()=>{
  state.isShowQR = true
  nextTick(() => {
    onInitQrcode()
  })
}
//
const onInitQrcode = () => {
  //清除上一次二维码
  let codeHtml = document.getElementById("qrcode");
  codeHtml.innerHTML = "";
  state.QRcode = new QRCode(qrcodeRef.value, {
    text: shopStoreData.currentOrder.value.pay_info.alipay_info.qr_code,
    width: 125,
    height: 125,
    colorDark: '#000000',
    colorLight: '#ffffff',
  });
}
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
  width: 100px;
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