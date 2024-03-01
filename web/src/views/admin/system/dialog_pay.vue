<template>
  <el-dialog v-model="state.isShowDialog" :title="state.title" width="80%">
    <el-form label-width="120px" label-position="top">
      <el-form-item :label="$t('message.adminShop.PayInfo.name')">
        <el-input v-model="adminShopStoreData.payInfo.value.name"/>
      </el-form-item>
      <el-form-item :label="$t('message.adminShop.PayInfo.pay_logo_url')">
        <el-select
            v-model="adminShopStoreData.payInfo.value.pay_logo_url"
            filterable
            allow-create
            default-first-option
            :reserve-keyword="false"
            style="width: 100%"
        >
          <el-option
              v-for="(v,k) in state.defaultPayLogoList"
              :key="k"
              :label="v.name"
              :value="v.url">
            <div style="display: flex;align-items: center">
              <el-text>{{ v.name }}</el-text>
              <el-image style="width: 30px; height: 30px" :src="v.url" fit="fill"/>
            </div>
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item :label="$t('message.adminShop.PayInfo.status')">
        <el-switch v-model="adminShopStoreData.payInfo.value.status" inline-prompt :active-text="$t('message.common.enable')"
                   :inactive-text="$t('message.common.disable')"
                   style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
      </el-form-item>
      <el-form-item :label="$t('message.adminShop.PayInfo.pay_type')">
        <el-radio-group v-model="adminShopStoreData.payInfo.value.pay_type">
          <el-radio label="epay">epay</el-radio>
          <el-radio label="alipay">alipay</el-radio>
          <el-radio label="balance">balance</el-radio>
        </el-radio-group>
      </el-form-item>
      <div v-if="adminShopStoreData.payInfo.value.pay_type==='alipay'">
        <el-form-item :label="$t('message.adminShop.Alipay.alipay_app_id')">
          <el-input v-model="adminShopStoreData.payInfo.value.alipay.alipay_app_id"/>
        </el-form-item>
        <el-form-item  :label="$t('message.adminShop.Alipay.alipay_app_private_key')">
          <el-input v-model="adminShopStoreData.payInfo.value.alipay.alipay_app_private_key" type="textarea"/>
        </el-form-item>
        <el-form-item  :label="$t('message.adminShop.Alipay.alipay_ali_public_key')">
          <el-input v-model="adminShopStoreData.payInfo.value.alipay.alipay_ali_public_key" type="textarea"/>
        </el-form-item>
        <el-form-item  :label="$t('message.adminShop.Alipay.alipay_encrypt_key')">
          <el-input v-model="adminShopStoreData.payInfo.value.alipay.alipay_encrypt_key"
                    type="textarea"/>
        </el-form-item>
      </div>
      <div v-else-if="adminShopStoreData.payInfo.value.pay_type==='epay'">
        <el-form-item :label="$t('message.adminShop.Epay.epay_pid')">
          <el-input v-model.number="adminShopStoreData.payInfo.value.epay.epay_pid"/>
        </el-form-item>
        <el-form-item :label="$t('message.adminShop.Epay.epay_key')">
          <el-input v-model="adminShopStoreData.payInfo.value.epay.epay_key"/>
        </el-form-item>
        <el-form-item  :label="$t('message.adminShop.Epay.epay_api_url')">
          <el-input v-model="adminShopStoreData.payInfo.value.epay.epay_api_url" placeholder="http://xxx.com/submit.php"/>
        </el-form-item>
      </div>

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

import {storeToRefs} from "pinia";
import {reactive} from "vue";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
import alipayLogo from '/@/assets/icon/alipay.jpeg';
import epayLogo from '/@/assets/icon/epay.png';
import balancePayLogo from '/@/assets/icon/balance.jpeg';
import { useAdminShopStore } from "/@/stores/admin_logic/shopStore";
import { useI18n } from "vue-i18n";

const adminShopStore = useAdminShopStore()
const adminShopStoreData = storeToRefs(adminShopStore)
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
const {t} = useI18n()

// 定义子组件向父组件传值/事件
const emit = defineEmits(['refresh']);

//定义参数
const state = reactive({
  isShowDialog: false,
  type: '',
  title: '',
  defaultPayLogoList: [
    {name: "alipay logo", url: alipayLogo},
    {name: "epay logo", url: epayLogo},
    {name: "balance logo", url: balancePayLogo},
  ],
})

// 打开弹窗
const openDialog = (type: string, row?: any) => {
  if (type == 'add') {
    state.type = type
    state.title = t('message.adminServer.addPay')
    state.isShowDialog = true
    //清空store中的payInfo
    adminShopStore.wipePayInfo()

  } else {
    state.type = type
    state.title = t('message.adminServer.modifyPay')
    state.isShowDialog = true
    adminShopStoreData.payInfo.value = row
  }
}
// 关闭弹窗
const closeDialog = () => {
  state.isShowDialog = false
};

//确认提交
function onSubmit() {
  if (state.type === 'add') {
    request(apiStoreData.adminApi.value.newPay, adminShopStoreData.payInfo.value).then((res) => {
      setTimeout(() => {
        emit('refresh');
      }, 500);
    })

  } else {
    request(apiStoreData.adminApi.value.updatePay, adminShopStoreData.payInfo.value).then((res) => {
      setTimeout(() => {
        emit('refresh');
      }, 500);
    })
  }
  closeDialog();
}

// 暴露变量
defineExpose({
  openDialog,
});
</script>

<style scoped>

</style>