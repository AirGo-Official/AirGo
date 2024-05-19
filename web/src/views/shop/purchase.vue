<template>
  <div>
    <el-dialog
      v-model="state.isShowDialog"
      :title="$t('message.common.details')"
      
      @close="closeDialog"
      destroy-on-close>
      <!--      1、现实进度条-->
      <el-steps :active="state.active" process-status="wait" finish-status="success">
        <el-step :title="$t('message.adminShop.selectGoods')">
          <template #icon>
            <i style="font-size: 2.4em;" class="ri-number-1"></i>
          </template>
        </el-step>
        <el-step :title="$t('message.adminOrder.orderDetails')">
          <template #icon>
            <i style="font-size: 2.4em;" class="ri-number-2"></i>
          </template>
        </el-step>
        <el-step :title="$t('message.adminShop.purchase')">
          <template #icon>
            <i style="font-size: 2.4em;" class="ri-number-3"></i>
          </template>
        </el-step>
      </el-steps>
      <!--      2、显示商品信息-->
      <div>
        <div v-if="state.active === 1 || state.active === 2">
          <!--          2-1、续费商品-->
          <el-row v-if="shopStoreData.currentOrder.value.order_type === constantStore.ORDER_TYPE_RENEW">
            <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24">
              <el-card style="margin-top: 10px;border-radius:10px;">
                <el-descriptions
                  :column="1"
                  border
                  size="small"
                  direction="horizontal"
                >
                  <el-descriptions-item :label="$t('message.adminUser.CustomerService.subject')">
                    {{ customerServiceStoreData.customerService.value.subject }}
                  </el-descriptions-item>
                  <el-descriptions-item :label="$t('message.adminUser.CustomerService.service_start_at')">
                    {{ DateStrToTime(customerServiceStoreData.customerService.value.service_start_at) }}
                  </el-descriptions-item>
                  <el-descriptions-item :label="$t('message.adminUser.CustomerService.service_end_at')">
                    {{ DateStrToTime(customerServiceStoreData.customerService.value.service_end_at) }}
                  </el-descriptions-item>
                  <el-descriptions-item :label="$t('message.adminUser.CustomerService.duration')">
                    {{ customerServiceStoreData.customerService.value.duration }}
                  </el-descriptions-item>
                  <el-descriptions-item :label="$t('message.adminUser.CustomerService.renewal_amount')">
                    {{ customerServiceStoreData.customerService.value.renewal_amount }}
                  </el-descriptions-item>

                </el-descriptions>
              </el-card>
            </el-col>
          </el-row>
          <!--          2-2、新购-->
          <el-row :gutter="50" v-else>
            <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24">
              <div style="margin-top: 10px;border-radius: 1vh;text-align: center;" >
                <div style="text-align: center">
              <el-image :src="shopStoreData.currentGoods.value.cover_image"
                        lazy
                        style="height: 20vh; border-radius: 1vh;"
                        fit="cover"
                        :preview-src-list="[shopStoreData.currentGoods.value.cover_image]">
                <template #error>
                  <div class="image-slot">
                    <i class="ri-signal-wifi-error-line"></i>
                  </div>
                </template>
              </el-image>
              <div class="item-title">
                  {{ shopStoreData.currentGoods.value.subject }}
                </div>
            </div>    

              </div>
            </el-col>

          </el-row>
        </div>
        <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24">
          <el-card v-if="state.active === 1
          && shopStoreData.currentOrder.value.order_type === constantStore.ORDER_TYPE_NEW
          && shopStoreData.currentOrder.value.goods_type === constantStore.GOODS_TYPE_SUBSCRIBE"
                   style="border-radius:10px;padding: 10px;margin-top: 0.8em;">
                <div style="margin-bottom: 20px;font-size: 20px" size="large" type="primary">{{ $t("message.adminOrder.Order.duration") }}</div>
                <div>
                    <el-radio-group v-model.number="shopStoreData.currentOrder.value.duration" size="small" @change="getOrderInfo">
                      <el-radio v-if="shopStoreData.currentGoods.value.price" class="mb15" :value="1" border>单月</el-radio>
                      <el-radio v-if="shopStoreData.currentGoods.value.price_3_month" class="mb15" :value="3" border>3个月</el-radio>
                      <el-radio v-if="shopStoreData.currentGoods.value.price_6_month" class="mb15" :value="6"  border>6个月</el-radio>
                      <el-radio v-if="shopStoreData.currentGoods.value.price_12_month" class="mb15" :value="12"  border>12个月</el-radio>
                      <el-radio v-if="shopStoreData.currentGoods.value.price_unlimited_duration" class="mb15" :value="-1"  border>不限时</el-radio>
                    </el-radio-group>
                </div>
          </el-card>
        <el-card v-if="state.active === 1"
             style="border-radius:10px;padding: 10px;margin-top: 0.8em;">
          <div v-if="shopStoreData.currentOrder.value.order_type === constantStore.ORDER_TYPE_NEW">
            <div class="card-text" v-if="shopStoreData.currentOrder.value.goods_type !== constantStore.GOODS_TYPE_RECHARGE" >
              <el-input v-model="shopStoreData.currentOrder.value.coupon_name"
                        :placeholder="$t('message.adminShop.Coupon.name')" size="default"
                        >
                <template #prepend>
                  <el-icon>
                    <Ticket />
                  </el-icon>
                </template>
                <template #append>
                  <el-button size="small" @click="getOrderInfo">
                    {{ $t("message.common.button_confirm") }}
                  </el-button>
                </template>
              </el-input>
            </div>
            <div class="card-text">
              <el-tag class="card-text-left" type="primary">{{ $t("message.adminOrder.Order.price") }}</el-tag>
              <el-text class="card-text-right">{{ shopStoreData.currentOrder.value.price }}</el-text>
              </div>
            <div class="card-text">
              <el-tag class="card-text-left" type="primary">{{ $t("message.adminOrder.Order.coupon_amount") }}
              </el-tag>
              <el-text class="card-text-right">-{{ shopStoreData.currentOrder.value.coupon_amount }}</el-text>
            </div>
            <div class="card-text">
              <el-tag class="card-text-left" type="primary">{{ $t("message.adminOrder.Order.total_amount") }}
              </el-tag>
              <el-text class="card-text-right" style="font-size: 25px;">{{ shopStoreData.currentOrder.value.total_amount
                }}
              </el-text>
            </div>
          </div>
        </el-card>
        </el-col>
      </div>

        <template #footer>
           <div v-if="state.active === 1">
            <el-button type="primary" @click="nextSubmitOrder">{{ $t("message.adminOrder.submitOrder") }}</el-button>
          </div>
          <div v-if="state.active === 2">
            <el-button color="blue" @click="closeDialog">{{ $t("message.common.button_cancel") }}</el-button>
            <el-button color="blue" @click="nextPurchase" :disabled="!shopStoreData.currentOrder.value.pay_id">
              {{ $t("message.adminShop.purchase") }}
            </el-button>
          </div>
          <div v-if="state.active === 3">
          <el-button type="primary" @click="back">{{ $t("message.common.back") }}</el-button>
          <el-button type="primary" @click="closeDialog">{{ $t("message.common.close") }}</el-button>
        </div>
        </template>

      <!--      3、点击购买后-->
      <div v-if="state.active === 2"
           v-loading="state.isShowLoading"
           element-loading-text="Loading..."
           element-loading-background="rgba(122, 122, 122, 0.8)"
           style="border-radius:10px;padding: 10px;margin-top: 10px">
        <div style="margin-bottom: 20px"
             v-if="state.active === 2 || state.active === 3">
        </div>
        <div v-if="state.active === 2">
          <div v-if="state.isShowPayment">
            <el-tag type="primary">{{ $t("message.adminOrder.Order.pay_type") }}</el-tag>
            <div>
              <el-radio-group v-model="shopStoreData.currentOrder.value.pay_id"
                              v-for="(v,k) in shopStoreData.payList.value" :key="k">
                <el-radio :label="v.id" style="width: 800px;margin-top: 10px">
                  <div style="display: flex;align-items: center">
                        <span style="margin-right: 5px">
                          <el-image :src="v.pay_logo_url" style="height: 15px;"></el-image>
                        </span>
                    <span>{{ v.name }}</span>
                    <span style="margin-left: 30px;color: #6cacf5" v-if="v.pay_type === constantStore.PAY_TYPE_BALANCE">[ {{ $t("message.adminUser.SysUser.balance") }}: {{ userInfos.balance }} ]</span>
                  </div>
                </el-radio>
              </el-radio-group>
            </div>

          </div>
        </div>
      </div>
      <!--      4、点击支付后-->
      <div v-if="state.active === 3"
           v-loading="state.isShowLoading"
           element-loading-text="Loading..."
           element-loading-background="rgba(122, 122, 122, 0.8)"
           style="border-radius:10px;padding: 10px;margin-top: 10px">
        <div>
          <el-result icon="success" :title="$t('message.adminShop.resultText1')" v-if="state.result === 1"></el-result>
          <el-result icon="warning" :title="$t('message.adminShop.resultText2')" v-else-if="state.result === 2">
            <template #extra>
              <div v-if="state.showPayInfo === 1">
                <div class="qrcode-img-warp">
                  <div class="mb30 mt30 qrcode-img">
                    <div class="qrcode" ref="qrcodeRef"></div>
                  </div>
                </div>
                <el-button type="warning" :href="state.alipayUrl">{{ $t("message.adminShop.resultText4") }}
                </el-button>
                <el-link type="primary" :href="state.alipayUrl" target="_blank">{{ state.alipayUrl }}</el-link>
              </div>
              <div v-else-if="state.showPayInfo === 2">
                <el-button type="primary" round size="large"  @Click="reject_epay" >{{ $t("message.adminShop.resultText5") }}
                </el-button>
              </div>
            </template>
          </el-result>
          <el-result icon="error" :title="$t('message.adminShop.resultText3')"
                     v-else-if="state.result === 3"></el-result>
        </div>

        <el-tag type="primary" style="margin-bottom: 1em;">{{ $t("message.adminOrder.orderDetails") }}</el-tag>
          <el-descriptions
            :column="1"
            border
            size="small"
            direction="horizontal"
          >
            <el-descriptions-item :label="$t('message.adminOrder.Order.out_trade_no')">
              {{ shopStoreData.currentOrder.value.out_trade_no }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('message.adminOrder.Order.created_at')">
              {{ DateStrToTime(shopStoreData.currentOrder.value.created_at) }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('message.adminOrder.Order.price')">
              {{ shopStoreData.currentOrder.value.price }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('message.adminOrder.Order.duration')">
              {{ shopStoreData.currentOrder.value.duration }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('message.adminOrder.Order.order_type')">
              <span
                v-if="shopStoreData.currentOrder.value.order_type === constantStore.ORDER_TYPE_NEW">{{ $t("message.constant.ORDER_TYPE_NEW")
                }}</span>
              <span
                v-else-if="shopStoreData.currentOrder.value.order_type === constantStore.ORDER_TYPE_RENEW">{{ $t("message.constant.ORDER_TYPE_RENEW")
                }}</span>
              <span v-else>{{ $t("message.constant.ORDER_TYPE_DESTROYED") }}</span>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('message.adminOrder.Order.original_amount')">
              {{ shopStoreData.currentOrder.value.original_amount }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('message.adminOrder.Order.coupon_amount')">
              {{ shopStoreData.currentOrder.value.coupon_amount }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('message.adminOrder.Order.total_amount')">
              <span style="color: red;font-size: 30px">{{ shopStoreData.currentOrder.value.total_amount }}</span>
            </el-descriptions-item>
          </el-descriptions>
      </div>

    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from "vue";
import { useShopStore } from "/@/stores/user_logic/shopStore";
import { storeToRefs } from "pinia";
import { useApiStore } from "/@/stores/apiStore";
import { request } from "/@/utils/request";
import { ElMessage, ElMessageBox } from "element-plus";
import qs from "qs";
import QRCode from "qrcodejs2-fixes";
import { useConstantStore } from "/@/stores/constantStore";
import { DateStrToTime } from "/@/utils/formatTime";
import { useI18n } from "vue-i18n";
import { useCustomerServiceStore } from "/@/stores/user_logic/customerServiceStore";
import { useUserStore } from "/@/stores/user_logic/userStore";
// import { isMobile } from "/@/utils/other";

const shopStore = useShopStore();
const shopStoreData = storeToRefs(shopStore);
const apiStore = useApiStore();
const apiStoreData = storeToRefs(apiStore);
const customerServiceStore = useCustomerServiceStore();
const customerServiceStoreData = storeToRefs(customerServiceStore);
const userStore = useUserStore();
const { userInfos } = storeToRefs(userStore);

const constantStore = useConstantStore();
const qrcodeRef = ref();
const { t } = useI18n();

const state = reactive({
  active: 1,
  isShowDialog: false,
  isShowQR: false,
  isShowPayment: false,
  isShowLoading: false,
  showPayInfo: 0,
  result: 0,
  alipayUrl: "",
  epayUrl: ""

});
const openDialog = (type: string) => {
  state.active = 1;
  state.isShowDialog = true;

  switch (type) {
    case constantStore.ORDER_TYPE_NEW:
      getOrderInfo();
      break;
    case constantStore.ORDER_TYPE_RENEW:
      break;
    case "fromMyOrder": //从我的订单跳转过来
      state.active = 2;
      state.isShowLoading = true;
      loop();
      break;
    default:
      break;
  }

};
const closeDialog = () => {
  state.active = 1;
  state.isShowDialog = false;
  state.isShowQR = false;
  state.isShowPayment = false;
  state.isShowLoading = false;
  state.showPayInfo = 0;
  state.result = 0;
  state.alipayUrl = "";
  state.epayUrl = "";
  shopStoreData.currentOrder.value = {} as Order;
  //如果在响应之前关闭弹窗，也要关闭定时器
  clearInterval(timer);
};
const getOrderInfo = () => {
  state.isShowDialog = true;
  shopStore.getOrderInfo(shopStoreData.currentOrder.value).then((res) => {
    state.isShowLoading = false;
    // shopStoreData.currentOrder.value = {} as Order;
    if (res.code === 10) { //code=10，能正常获取请求数据，但有重要message 需要显式提醒。此处用来提示折扣码是否有效
      shopStoreData.currentOrder.value.coupon_name = "";
      ElMessageBox.confirm(res.msg, t("message.common.tip"), {
        cancelButtonText: t("message.common.button_cancel"),
        type: "warning"
      }).then(() => {
      });
    }
    if (res.data) {
      shopStoreData.currentOrder.value = res.data;
    }
  });
};
const next = () => {
  if (state.active === 3) {

  } else {
    state.active++;
  }
};
const back = () => {
  if (state.active === 1) {

  } else {
    state.active--;
  }
};

//轮询定时器
let timer: NodeJS.Timeout
const loop = () => {
  let i = 0;
  timer = setInterval(() => {
    getOrderInfoWaitPay(timer, i++);
  }, 3000);
};

const getOrderInfoWaitPay = (timer: NodeJS.Timeout, i: number) => {
  state.isShowLoading = true;
  setTimeout(() => {
    // console.log("开始轮循请求,次数：", i);
    //请求
    request(apiStoreData.userApi.value.getOrderInfoWaitPay, shopStoreData.currentOrder.value).then((res) => {
      //保存订单信息到pinia
      shopStoreData.currentOrder.value = res.data;
      if (shopStoreData.currentOrder.value.trade_status === constantStore.ORDER_STATUS_WAIT_BUYER_PAY) {
        // 关闭轮询
        clearInterval(timer);
        // 获取支付列表
        shopStore.getEnabledPayList();
        // 显示支付
        state.isShowPayment = true;
        // 关闭loading
        state.isShowLoading = false;
      }
    });
    if (i >= 100) {
      clearInterval(timer);
      ElMessageBox.alert(t("message.adminShop.resultText6"), t("message.common.tip"), {
        confirmButtonText: t("message.common.close")
      })
        .then(() => {
          closeDialog();
        });

    }
  }, 0);
};

const nextSubmitOrder = () => {
  state.isShowLoading = true;
  shopStoreData.currentOrder.value.id = 0;
  request(apiStoreData.userApi.value.preCreateOrder, shopStoreData.currentOrder.value).then((res) => {
    //保存订单信息到pinia
    shopStoreData.currentOrder.value = res.data;
    //
    next();
    //
    loop();
  });
};
const nextPurchase = () => {
  //el-steps 跳转到第3步
  state.active = 3;
  state.isShowLoading = true;
  request(apiStoreData.userApi.value.purchase, shopStoreData.currentOrder.value).then((res) => {
    state.isShowLoading = false;
    //保存支付信息
    shopStoreData.currentOrder.value = res.data;
    let alipayInfo = shopStoreData.currentOrder.value.pay_info.alipay_info;
    let epayInfo = shopStoreData.currentOrder.value.pay_info.epay_info;

    switch (shopStoreData.currentOrder.value.pay_type) {
      case constantStore.PAY_TYPE_BALANCE:
        ElMessage.success(res.msg);
        state.result = 1;
        break;
      case constantStore.PAY_TYPE_ALIPAY:
        state.result = 2;
        // console.log("alipayInfo.qr_code:", alipayInfo.qr_code);
        if (alipayInfo.qr_code) {       // 1、支付宝
          state.showPayInfo = 1;
          state.alipayUrl = alipayInfo.qr_code;
          showQR();
        }
        break;
      case constantStore.PAY_TYPE_EPAY:
        state.result = 2;
        // console.log("epayInfo.epay_api_url:", epayInfo.epay_api_url);
        if (epayInfo.epay_api_url !== "") { //2、易支付
          state.showPayInfo = 2;
          let params = qs.stringify(epayInfo.epay_pre_create_pay);
          // window.location.href = epayInfo.epay_api_url + "?" + params;
          state.epayUrl = epayInfo.epay_api_url + "?" + params;
        }
        break;
      default:
        state.result = 3;
        break;
    }
  });
};
//
const showQR = () => {
  state.isShowQR = true;
  setTimeout(() => {
    //清除上一次二维码
    qrcodeRef.value.innerHTML = "";
    new QRCode(qrcodeRef.value, {
      text: shopStoreData.currentOrder.value.pay_info.alipay_info.qr_code,
      width: 125,
      height: 125,
      colorDark: "#000000",
      colorLight: "#ffffff"
    });
  }, 500);

};
//暴露变量
defineExpose({
  openDialog
});

const reject_epay = () =>{
  window.location.href = state.epayUrl 

}

</script>

<style scoped lang="scss">

.card-text {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 35px
}

.card-text-left {
  width: 100px;

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

.qrcode-img-warp {
  text-align: center;

  .qrcode-img {
    display: flex;
    width: 100%;
    height: 125px;

    .qrcode {
      margin: auto;
      width: 125px;
      height: 125px;
    }
  }
}

.image {
  padding: 30px 0;
  text-align: center;
  border-right: solid 1px var(--el-border-color);
  display: inline-block;
  width: 20%;
  box-sizing: border-box;
  vertical-align: top;
}

</style>