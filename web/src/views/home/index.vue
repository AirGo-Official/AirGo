<template>
  <div>
    <div class="home-card-item" v-if="customerServiceStoreData.customerServiceList.value.length === 0">
      <el-card>
        <el-skeleton :rows="5" animated />
        <h2>{{$t('message.home.no_data')}}</h2>
      </el-card>
    </div>

    <el-row :gutter="15" class="home-card-one mb15">
      <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12"
              v-for="(v, k) in customerServiceStoreData.customerServiceList.value" :key="k">
        <div class="home-card-item">
          <el-card>
            <template #header>
              <div class="card-header-left">
                <el-text >{{ v.subject }}</el-text>
              </div>
            </template>
            <el-descriptions
              :column="4"
              border
              size="small"
              direction="vertical"
            >
              <el-descriptions-item :label="$t('message.home.des_start')">{{ DateStrToTime(v.service_start_at) }}</el-descriptions-item>
              <el-descriptions-item :label="$t('message.home.des_end')">{{ DateStrToTime(v.service_end_at) }}</el-descriptions-item>
              <el-descriptions-item :label="$t('message.home.des_SubStatus')">
                <el-icon v-if="v.sub_status" color="green" size="large">
                  <SuccessFilled />
                </el-icon>
                <el-icon v-else color="red" size="large">
                  <CircleCloseFilled />
                </el-icon>
              </el-descriptions-item>
              <el-descriptions-item :label="$t('message.home.des_renewAmount')">{{ v.renewal_amount }}</el-descriptions-item>
              <el-descriptions-item :label="$t('message.home.des_trafficResetDay')">{{ v.traffic_reset_day }}</el-descriptions-item>
              <el-descriptions-item :label="$t('message.home.des_used')">
                <span>{{ ((v.used_up + v.used_down) / 1024 / 1024 / 1024).toFixed(2) }}GB / {{ (v.total_bandwidth / 1024 / 1024 / 1024).toFixed(2)
                  }}GB</span>
              </el-descriptions-item>
              <el-descriptions-item :label="$t('message.home.des_usageRate')">
                <el-progress
                  :color="customColors"
                  striped
                  striped-flow
                  :text-inside="true" :stroke-width="16"
                  :percentage="Number((((v.used_up + v.used_down)/v.total_bandwidth)*100).toFixed(2))">
                </el-progress>
              </el-descriptions-item>
            </el-descriptions>
            <div style="margin-top: 15px">
              <el-button size="small" style="margin-bottom: 10px" type="primary" @click="openDialogCustomerServiceDetails">{{$t('message.home.button_details')}}</el-button>
              <el-button size="small" style="margin-bottom: 10px" type="primary" @click="renew(v)">{{$t('message.home.button_renew')}}</el-button>
              <el-button size="small" style="margin-bottom: 10px" type="primary" @click="openPushDialog">{{$t('message.home.button_push')}}</el-button>
              <el-button size="small" style="margin-bottom: 10px" type="primary" @click="resetSubscribeUUID(v)">{{$t('message.home.button_resetSub')}}</el-button>
              <el-button size="small" style="margin-bottom: 10px" type="primary" @click="openSubDialog(v.sub_uuid)">{{$t('message.home.button_copySub')}}</el-button>
            </div>
          </el-card>
        </div>
      </el-col>
    </el-row>
    <!--    复制订阅弹窗-->
    <el-dialog v-model="state.isShowSubDialog" destroy-on-close>
      <div class="qrcode-img-warp">
        <div class="mb30 mt30 qrcode-img">
          <!-- 二维码弹窗 -->
          <div id="qrcode" class="qrcode" ref="qrcodeRef"></div>
        </div>
      </div>
      <div v-for="(v,k) in state.subType" :key="k">
        <el-row style="margin-top:10px;display: flex; justify-content: space-between; align-items: center;">
          <el-col :span="16">
            <el-button size="large" color="blue" style="width: 100%" @click="copyLink(v)">{{ v }} {{$t('message.home.subscription')}}</el-button>
          </el-col>
          <el-col :span="8">
            <el-button size="large" @click="showQR(v)">
              <el-icon>
                <FullScreen />
              </el-icon>
              QR code
            </el-button>
          </el-col>
        </el-row>
      </div>
    </el-dialog>
    <!--    push弹窗-->
    <el-dialog v-model="state.isShowPushDialog" destroy-on-close align-center>
      <div>
        <el-form :model="customerServiceStoreData.pushCustomerServiceRequest.value" label-position="top">
          <el-form-item :label="$t('message.home.target_username')">
            <el-input v-model="customerServiceStoreData.pushCustomerServiceRequest.value.to_user_name"></el-input>
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closePushDialog">{{$t('message.common.button_cancel')}}</el-button>
          <el-button type="primary" @click="toPush">{{$t('message.common.button_confirm')}}</el-button>
        </div>
      </template>

    </el-dialog>
    <!--    续费弹窗-->
    <Purchase ref="PurchaseRef"></Purchase>
    <!--    详情弹窗-->
    <DialogCustomerServiceDetails ref="DialogCustomerServiceDetailsRef"></DialogCustomerServiceDetails>
<!--    默认弹窗-->
    <DefaultDialog ref="DefaultDialogRef"></DefaultDialog>
  </div>

</template>

<script setup lang="ts">
import { onMounted, reactive, ref, defineAsyncComponent } from "vue";
import { useApiStore } from "/@/stores/apiStore";
import { DateStrToTime } from "/@/utils/formatTime";
import commonFunction from "/@/utils/commonFunction";
import { usePublicStore } from "/@/stores/publicStore";
import { storeToRefs } from "pinia";
import QRCode from "qrcodejs2-fixes";
import { useCustomerServiceStore } from "/@/stores/user_logic/customerServiceStore";
import { ElMessage, ElMessageBox } from "element-plus";
import { v4 as uuid } from 'uuid';
import { useShopStore } from "/@/stores/user_logic/shopStore";
import { useConstantStore } from "/@/stores/constantStore";
import { useI18n } from "vue-i18n";
import { useArticleStore } from "/@/stores/user_logic/articleStore";

const constantStore = useConstantStore()
const shopStore = useShopStore()
const shopStoreData = storeToRefs(shopStore)
const apiStore = useApiStore();
const publicStore = usePublicStore();
const publicStoreData = storeToRefs(publicStore);
const customerServiceStore = useCustomerServiceStore();
const customerServiceStoreData = storeToRefs(customerServiceStore);
const { copyText } = commonFunction();
const qrcodeRef = ref();
const Purchase = defineAsyncComponent(() => import("/@/views/shop/purchase.vue"));
const PurchaseRef = ref();
const {t} = useI18n()
const articleStore = useArticleStore()
const DefaultDialog = defineAsyncComponent( () => import('/@/views/default/defaultDialog.vue'));
const DefaultDialogRef = ref();

//组件
const DialogCustomerServiceDetails = defineAsyncComponent(() => import("/@/views/home/dialog_customer_service_details.vue"));
const DialogCustomerServiceDetailsRef = ref();
// const format = (percentage) => (percentage === 100 ? 'Full' : `${percentage}%`)
const state = reactive({
  isShowSubDialog: false,
  isShowPushDialog: false,
  subType: ["NekoBox", "v2rayNG", "v2rayN", "Shadowrocket", "Clash", "Surge", "Quantumult", "V2rayU"],
  currentSubUUID: "",
  QRcode: null,
});
const customColors = [
  { color: "#9af56c", percentage: 20 },
  { color: "#5cb87a", percentage: 40 },
  { color: "#f8c67e", percentage: 60 },
  { color: "#ff785a", percentage: 80 },
  { color: "#fa193b", percentage: 100 }
];
// 获取列表
const getCustomerServiceList = () => {
  customerServiceStore.getCustomerServiceList();
};
const getPublicSetting = () => {
  publicStore.getPublicSetting();
};
const openSubDialog = (subUUID: string) => {
  state.isShowSubDialog = true;
  state.currentSubUUID = subUUID.replace(/-/g, "");
};
const openPushDialog = (cs: CustomerService) => {
  state.isShowPushDialog = true;
  customerServiceStoreData.pushCustomerServiceRequest.value.customer_service_id = cs.id;
};
const closePushDialog = () => {
  state.isShowPushDialog = false;
};
const toPush = () => {
  customerServiceStore.pushCustomerService().then(() => {
    closePushDialog();
  });
};
const renew=(cs:CustomerService)=>{
  //保存用户服务
  customerServiceStoreData.customerService.value = cs
  //构造订单数据。需要3个参数，user_id，order_type，customer_service_id。
  // user_id由后端自动填充，前端传customer_service_id，order_type
  shopStoreData.currentOrder.value = {
    order_type:constantStore.ORDER_TYPE_RENEW,
    customer_service_id:cs.id,
  } as Order

  PurchaseRef.value.openDialog(constantStore.ORDER_TYPE_RENEW);
}
const resetSubscribeUUID=(cs:CustomerService)=>{
  ElMessageBox.alert(t('message.home.message_confirm_reset_sub'),t('message.common.tip'), {
    confirmButtonText: t('message.common.button_confirm'),
  })
    .then(() => {
      customerServiceStore.resetSubscribeUUID({id:cs.id,sub_uuid:uuid()} as CustomerService).then((res)=>{
        ElMessage.success(res.msg)
        getCustomerServiceList()
      })
    })
    .catch(() => {
    });
}
const copyLink = (subType: string) => {
  copyText(publicStoreData.publicSetting.value.backend_url + "/api/public/sub/" + state.currentSubUUID + "?type=" + subType);
};
const showQR = (subType: string) => {
  let link = publicStoreData.publicSetting.value.backend_url + "/api/public/sub/" + state.currentSubUUID + "?type=" + subType;
  //清除上一次二维码
  qrcodeRef.value.innerHTML = "";
  new QRCode(qrcodeRef.value, {
    text: link,
    width: 125,
    height: 125,
    colorDark: "#000000",
    colorLight: "#ffffff"
  });
};
const openDialogCustomerServiceDetails = () => {
  DialogCustomerServiceDetailsRef.value.openDialog();
};
const defaultArticle=()=>{
  articleStore.getDefaultArticles().then(()=>{
    setTimeout(()=>{
      DefaultDialogRef.value.openDialog()
    },2000)
  })
}

onMounted(() => {
  getCustomerServiceList();
  getPublicSetting();
  defaultArticle()
});


</script>

<style scoped>
.home-card-item {
  width: 100%;
  height: 100%;
  border-radius: 4px;
  transition: all ease 0.3s;
  overflow: hidden;
  padding: 10px;
  /*background: var(--el-color-white);*/
  color: var(--el-text-color-primary);
  /*border: 1px solid var(--next-border-color-light);*/
}

.card-header-left {
  display: block;
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
  font-weight: bolder;
}
</style>