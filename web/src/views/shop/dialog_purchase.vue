<template>
  <div>
    <!-- 支付弹窗 -->
    <el-dialog v-model="state.isShowPurchaseDialog" title="等待支付" width="80%">
      <div class="home-card-item">
        <div class="card-text">
          <el-button type="primary">套餐</el-button>
          <el-text class="card-text-right">{{ shopData.currentOrder.subject }}</el-text>
        </div>
        <div class="card-text">
          <el-button type="primary">金额</el-button>
          <el-text class="card-text-right">{{ shopData.currentOrder.total_amount }}元</el-text>
        </div>
        <el-divider></el-divider>
        <div v-if="shopData.currentOrder.total_amount!=='0'">
          <el-button type="primary">支付方式</el-button>
          <div>
            <el-radio-group style="height: 60px;" v-model="shopData.currentOrder.pay_id" class="ml-4">
              <el-radio :label="v.id" v-for="(v,k) in payStoreData.payList.value" :key="k">
                <div style="display: flex;align-items: center">
                  <el-image :src="v.pay_logo_url" style="height: 15px;"></el-image>
                  <span>{{ v.name }}</span>
                </div>
              </el-radio>
            </el-radio-group>
          </div>
        </div>
      </div>
      <template #footer>
            <span class="dialog-footer">
                <el-button @click="state.isShowPurchaseDialog = false">取消</el-button>
                <el-button color="#FC3D08"
                           @click="onPurchase({id:shopData.currentGoods.id})">
                    确认支付
                </el-button>
            </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import {ElMessage} from 'element-plus';
import {reactive} from "vue";

import {storeToRefs} from 'pinia';
import {useShopStore} from "/@/stores/shopStore";
import {isMobile} from "/@/utils/other";

import {usePayStore} from "/@/stores/payStore";

import * as qs from "qs";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
const shopStore = useShopStore()
const {shopData} = storeToRefs(shopStore)
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
const payStore = usePayStore()
const payStoreData = storeToRefs(payStore)

//定义变量
const state = reactive({
  isShowPurchaseDialog: false,
})
//打开弹窗
const openDialog = () => {
  state.isShowPurchaseDialog = true

  //获取支付列表
  payStore.getEnabledPayList()
}
//调用父组件
const emits = defineEmits(['openQRDialog'])

//购买按钮
const onPurchase = (params?: object) => {
  //传out_trade_no，pay_id
  request(apiStoreData.api.value.shop_purchase, {
    out_trade_no: shopData.value.currentOrder.out_trade_no,
    pay_id: shopData.value.currentOrder.pay_id
  }).then((res) => {
    if (res.msg === "购买成功") { //0元购
      ElMessage.success("购买成功")
      state.isShowPurchaseDialog = false
      return
    } else {
      //保存支付信息
      shopData.value.currentOrder.pay_info = res.data
      let pay_info = shopData.value.currentOrder.pay_info
      if (pay_info.alipay_info.qr_code !== "") { //支付宝支付
        const ok = isMobile()
        if (ok) {
          window.location.href = pay_info.alipay_info.qr_code; //手机端跳转支付页面
          return
        } else {
          emits('openQRDialog')   //电脑端打开支付二维码弹窗
        }
      } else if (pay_info.epay_info.epay_api_url !== "") {
        // 对象转url参数
        let params = qs.stringify(pay_info.epay_info.epay_pre_create_pay)
        window.location.href = pay_info.epay_info.epay_api_url + "?" + params
      }
    }
  }).catch(() => {
  })
  //关闭弹窗
  state.isShowPurchaseDialog = false
}

defineExpose({
  openDialog,// 打开弹窗
})
</script>

<style scoped>

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