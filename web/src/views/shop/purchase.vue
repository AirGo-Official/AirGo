<template>
  <div >
    <el-dialog v-loading.fullscreen.lock="state.isShowLoading" v-model="state.isShowDialog" title="详情" width="80%" destroy-on-close>
      <el-steps :active="state.active" process-status="wait" finish-status="success">
        <el-step title="选择套餐">
          <template #icon>
            <SvgIcon name="iconfont icon-1_round_solid" :size="40" />
          </template>
        </el-step>
        <el-step title="订单详情">
          <template #icon>
            <SvgIcon name="iconfont icon-2_round_solid" :size="40" />
          </template>
        </el-step>
        <el-step title="支付">
          <template #icon>
            <SvgIcon name="iconfont icon-3_round_solid" :size="40" />
          </template>
        </el-step>
      </el-steps>
      <div >
        <el-row :gutter="50">
          <el-col :xs="24" :sm="24" :md="6" :lg="6" :xl="6">
            <div
              style="margin-top: 10px;height: 100px;border-radius:10px;background: rgba(224,224,224,0.29);padding: 10px">
              <el-image :src="shopStoreData.currentGoods.value.cover_image" style="height: 100%">
                <template #error>
                  <div class="image-slot">
                    <el-icon>
                      <icon-picture />
                    </el-icon>
                  </div>
                </template>
              </el-image>
            </div>
          </el-col>
          <el-col :xs="24" :sm="24" :md="18" :lg="18" :xl="18">
            <div style="border-radius:10px;background: rgba(224,224,224,0.29);padding: 10px;margin-top: 10px">
              <div style="margin-top: 10px;">
                {{ shopStoreData.currentGoods.value.subject }}
              </div>
              <div style="margin-top: 10px;display: flex;">
                <el-tag class="ml-2"
                        v-if="shopStoreData.currentGoods.value.goods_type === constantStore.GOODS_TYPE_SUBSCRIBE">
                  {{ $t("message.adminShop.Goods.goods_type") }}: {{ $t("message.constant.GOODS_TYPE_SUBSCRIBE") }}
                </el-tag>
                <el-tag class="ml-2"
                        v-if="shopStoreData.currentGoods.value.goods_type === constantStore.GOODS_TYPE_RECHARGE">
                  {{ $t("message.adminShop.Goods.goods_type") }}: {{ $t("message.constant.GOODS_TYPE_RECHARGE") }}
                </el-tag>
                <el-tag class="ml-2"
                        v-if="shopStoreData.currentGoods.value.goods_type === constantStore.GOODS_TYPE_GENERAL">
                  {{ $t("message.adminShop.Goods.goods_type") }}: {{ $t("message.constant.GOODS_TYPE_GENERAL") }}
                </el-tag>

                <el-tag class="ml-2"
                        v-if="shopStoreData.currentGoods.value.deliver_type === constantStore.DELIVER_TYPE_NONE"
                        type="success">{{ $t("message.adminShop.Goods.deliver_type") }}:
                  {{ $t("message.constant.DELIVER_TYPE_NONE") }}
                </el-tag>
                <el-tag class="ml-2"
                        v-if="shopStoreData.currentGoods.value.deliver_type === constantStore.DELIVER_TYPE_AUTO"
                        type="success">{{ $t("message.adminShop.Goods.deliver_type") }}:
                  {{ $t("message.constant.DELIVER_TYPE_AUTO") }}
                </el-tag>
                <el-tag class="ml-2"
                        v-if="shopStoreData.currentGoods.value.deliver_type === constantStore.DELIVER_TYPE_MANUAL">
                  type="success"{{ $t("message.adminShop.Goods.deliver_type") }}:
                  {{ $t("message.constant.DELIVER_TYPE_MANUAL") }}
                </el-tag>

                <el-tag type="warning">
                  {{ $t("message.adminShop.Goods.quota") }}：{{ shopStoreData.currentGoods.value.quota }} /
                  {{ $t("message.adminShop.Goods.stock") }}：{{ shopStoreData.currentGoods.value.stock }}
                </el-tag>
              </div>
              <div style="margin-top: 10px">
                <span>
                  <span style="color: red;">￥</span>
                  <span style="color: red;font-size: 30px;">{{ shopStoreData.currentGoods.value.price }}</span>
                  <span
                    v-if="shopStoreData.currentGoods.value.goods_type === constantStore.GOODS_TYPE_SUBSCRIBE"> / {{ $t("message.common.month")
                    }}</span>
                </span>
              </div>
            </div>

          </el-col>
        </el-row>
        <div v-if="state.active === 1"
             style="border-radius:10px;background: rgba(224,224,224,0.29);padding: 10px;margin-top: 10px">
          <div class="card-text">
            <el-input v-model="shopStoreData.currentOrder.value.coupon_name" placeholder="输入折扣码" size="default">
              <template #prepend>
                <el-icon>
                  <Ticket />
                </el-icon>
              </template>
              <template #append>
                <el-button class="card-text-right" color="blue" size="small" @click="getOrderInfo">验证</el-button>
              </template>
            </el-input>
          </div>
          <div class="card-text">
            <el-button class="card-text-left" type="info">订购时长</el-button>
            <el-input-number class="card-text-right"
                             v-model.number="shopStoreData.currentOrder.value.duration"></el-input-number>
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
            <el-text class="card-text-right" style="font-size: 25px;">{{ shopStoreData.currentOrder.value.total_amount
              }}
            </el-text>
          </div>
          <div style="text-align: right;margin-top: 20px">
            <el-button color="blue" @click="nextSubmitOrder">提交订单</el-button>
          </div>
        </div>
      </div>
      <div v-if="state.active === 2"
           style="border-radius:10px;background: rgba(224,224,224,0.29);padding: 10px;margin-top: 10px">
        <div style="margin-bottom: 20px">
          <el-button type="primary">订单信息</el-button>
          <el-descriptions
            :column="1"
            border
            size="small"
            direction="horizontal"
          >
              <el-descriptions-item label="订单号">{{ shopStoreData.currentOrder.value.out_trade_no }}</el-descriptions-item>
              <el-descriptions-item label="创建时间">{{ DateStrToTime(shopStoreData.currentOrder.value.created_at) }}</el-descriptions-item>
              <el-descriptions-item label="订购周期">{{ shopStoreData.currentOrder.value.duration}}</el-descriptions-item>
              <el-descriptions-item label="类型">
                <span v-if="shopStoreData.currentOrder.value.order_type === constantStore.ORDER_TYPE_NEW">{{$t('message.constant.ORDER_TYPE_NEW')}}</span>
                <span v-else-if="shopStoreData.currentOrder.value.order_type === constantStore.ORDER_TYPE_RENEW">{{$t('message.constant.ORDER_TYPE_RENEW')}}</span>
                <span v-else>{{$t('message.constant.ORDER_TYPE_DESTROYED')}}</span>
              </el-descriptions-item>
            <el-descriptions-item label="原始价格">{{ shopStoreData.currentOrder.value.original_amount}}</el-descriptions-item>
            <el-descriptions-item label="折扣">{{ shopStoreData.currentOrder.value.coupon_amount}}</el-descriptions-item>
            <el-descriptions-item label="订单价格">
              <span style="color: red;font-size: 30px">{{ shopStoreData.currentOrder.value.total_amount}}</span>
            </el-descriptions-item>
          </el-descriptions>
        </div>
        <div v-if="state.isShowPayment">
          <el-button type="primary">支付方式</el-button>
          <div>
            <el-radio-group v-model="shopStoreData.currentOrder.value.pay_id"
                            v-for="(v,k) in shopStoreData.payList.value" :key="k">
              <el-radio :label="v.id" style="width: 800px;margin-top: 10px">
                <div style="display: flex;align-items: center">
                        <span style="margin-right: 5px">
                          <el-image :src="v.pay_logo_url" style="height: 15px;"></el-image>
                        </span>
                  <span>{{ v.name }}</span>
                </div>
              </el-radio>
            </el-radio-group>
          </div>
        </div>
        <div style="text-align: right;margin-top: 20px">
          <el-button color="blue" @click="closeDialog">取消</el-button>
          <el-button color="blue" @click="nextPurchase" :disabled="!shopStoreData.currentOrder.value.pay_id">确认购买
          </el-button>
        </div>
      </div>
      <div v-if="state.active === 3"
           style="border-radius:10px;background: rgba(224,224,224,0.29);padding: 10px;margin-top: 10px">
        <div>
          <el-result icon="success" title="支付成功" subTitle="请返回首页查看服务状态" v-if="state.result === 1"></el-result>
          <el-result icon="warning" title="注意" subTitle="请尽快完成支付" v-else-if="state.result === 2">
            <template #extra>
              <div v-if="state.showPayInfo === 1">
                <div>
                  <div class="qrcode-img-warp">
                    <div class="mb30 mt30 qrcode-img">
                      <!-- 二维码弹窗 -->
                      <div id="qrcode" ref="qrcodeRef"></div>
                    </div>
                  </div>
                </div>
                <el-button :href="shopStoreData.currentOrder.value.pay_info">前往支付宝</el-button>
                <el-text>{{ shopStoreData.currentOrder.value.pay_info }}</el-text>
              </div>
              <div v-else-if="state.showPayInfo === 2">
                <el-button :href="shopStoreData.currentOrder.value.pay_info">前往易支付</el-button>
                <el-text>{{ shopStoreData.currentOrder.value.pay_info }}</el-text>
              </div>
            </template>
          </el-result>
          <el-result icon="error" title="错误" subTitle="支付遇到问题" v-else-if="state.result === 3"></el-result>
        </div>
        <div style="text-align: right;margin-top: 20px">
          <el-button color="blue" @click="closeDialog">关闭</el-button>
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
// import { isMobile } from "/@/utils/other";
import qs from "qs";
import QRCode from "qrcodejs2-fixes";
import { useConstantStore } from "/@/stores/constantStore";
import { DateStrToTime } from "/@/utils/formatTime";


const shopStore = useShopStore();
const shopStoreData = storeToRefs(shopStore);
const apiStore = useApiStore();
const apiStoreData = storeToRefs(apiStore);

const constantStore = useConstantStore();
const qrcodeRef = ref();

const state = reactive({
  active: 1,
  isShowDialog: false,
  isShowQR: false,
  isShowPayment: false,
  isShowLoading: false,
  showPayInfo: 0,
  result: 0,
  QRcode: null,
  epayUrl: ""

});
const openDialog = () => {
  state.active = 1;
  state.isShowDialog = true;
};
const closeDialog = () => {
  state.active = 1;
  state.isShowDialog = false;
  state.isShowQR = false;
  state.isShowPayment = false;
  state.isShowLoading = false;
  state.showPayInfo = 0;
  state.result = 0;
  state.QRcode = null;
  state.epayUrl = "";
};
const getOrderInfo = () => {
  state.isShowDialog = true;
  shopStore.getOrderInfo(shopStoreData.currentOrder.value).then((res) => {
    state.isShowLoading = false;
    shopStoreData.currentOrder.value = {} as Order;
    shopStoreData.currentOrder.value = res.data;
  }).catch(() => {
    state.isShowLoading = false;
  });
};
const next = () => {
  if (state.active === 3) {

  } else {
    state.active++;
  }
};

function loop() {
  let i = 0;
  let timer = setInterval(() => {
    getOrderInfoWaitPay(timer, i++);
  }, 3000);
}

function getOrderInfoWaitPay(timer: NodeJS.Timeout, i: number) {
  setTimeout(() => {
    console.log("开始轮循请求,次数：", i);
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
    }).catch(() => {
      // 关闭loading
      state.isShowLoading = false;
    });
    if (i >= 6) {
      clearInterval(timer);
      ElMessageBox.alert("请求超时，请稍后检查我的订单，有无订单记录", "提示", {
        confirmButtonText: "OK"
      })
        .then(() => {
          closeDialog();
        })
        .catch(() => {
        });

    }
  }, 0);
}

const nextSubmitOrder = () => {
  // 加载loading
  state.isShowLoading = true;
  //
  shopStoreData.currentOrder.value.id = 0;
  // 轮询
  request(apiStoreData.userApi.value.preCreatePay, shopStoreData.currentOrder.value).then((res) => {
    //保存订单信息到pinia
    shopStoreData.currentOrder.value = res.data;
    //
    next();
    //
    loop();
  }).catch(error => {
    // 关闭loading
    state.isShowLoading = false;
  });
};
const nextPurchase = () => {
  //el-steps 跳转到第3步
  state.active = 3;
  state.isShowLoading = true;
  request(apiStoreData.userApi.value.purchase, shopStoreData.currentOrder.value).then((res) => {
    if (res.msg === "Purchase success") { //交易成功
      ElMessage.success(res.msg);
      state.isShowLoading = false;
      state.result = 1;
    } else {
      state.result = 2;
      //保存支付信息
      shopStoreData.currentOrder.value.pay_info = res.data;
      let alipayInfo = shopStoreData.currentOrder.value.pay_info.alipay_info
      let epayInfo = shopStoreData.currentOrder.value.pay_info.epay_info
      if (alipayInfo.qr_code !== "") {       // 1、支付宝
        state.showPayInfo = 1
        // if (isMobile())
        // window.location.href = alipayInfo.qr_code;
          shoeQR();
      } else if (epayInfo.epay_api_url !== "") { //2、易支付
        state.showPayInfo = 2
        let params = qs.stringify(epayInfo.epay_pre_create_pay);
        // window.location.href = epayInfo.epay_api_url + "?" + params;
        state.epayUrl = epayInfo.epay_api_url + "?" + params;
      }
    }
  }).catch(() => {
    state.result = 3
  });
};
//
const shoeQR = () => {
  state.isShowQR = true;
  nextTick(() => {
    onInitQrcode();
  });
};
//
const onInitQrcode = () => {
  //清除上一次二维码
  let codeHtml = document.getElementById("qrcode");
  codeHtml.innerHTML = "";
  state.QRcode = new QRCode(qrcodeRef.value, {
    text: shopStoreData.currentOrder.value.pay_info.alipay_info.qr_code,
    width: 125,
    height: 125,
    colorDark: "#000000",
    colorLight: "#ffffff"
  });
};
//暴露变量
defineExpose({
  openDialog
});

</script>

<style scoped>

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