<template>
  <el-dialog v-model="state.isShowDialog" :title="state.title" width="80%">
    <el-form label-width="120px">
      <el-form-item label="支付别名">
        <el-input v-model="payStoreData.payInfo.value.name"/>
      </el-form-item>
      <el-form-item label="支付logo">
        <el-select
            v-model="payStoreData.payInfo.value.pay_logo_url"
            filterable
            allow-create
            default-first-option
            :reserve-keyword="false"
            placeholder="输入支付logo url"
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
        <div style="color: #9b9da1">如需自定义请输入图片url</div>
      </el-form-item>
      <el-form-item label="是否启用">
        <el-switch v-model="payStoreData.payInfo.value.status" inline-prompt active-text="开启"
                   inactive-text="关闭"
                   style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
      </el-form-item>
      <el-form-item label="支付类型">
        <el-radio-group v-model="payStoreData.payInfo.value.pay_type">
          <el-radio label="epay">epay</el-radio>
          <el-radio label="alipay">alipay</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item v-if="payStoreData.payInfo.value.pay_type==='alipay'" label="支付宝appID">
        <el-input v-model="payStoreData.payInfo.value.alipay.alipay_app_id"/>
      </el-form-item>
      <el-form-item v-if="payStoreData.payInfo.value.pay_type==='alipay'" label="支付宝应用私钥">
        <el-input v-model="payStoreData.payInfo.value.alipay.alipay_app_private_key" type="textarea"/>
      </el-form-item>
      <el-form-item v-if="payStoreData.payInfo.value.pay_type==='alipay'" label="支付宝公钥">
        <el-input v-model="payStoreData.payInfo.value.alipay.alipay_ali_public_key" type="textarea"/>
      </el-form-item>
      <el-form-item v-if="payStoreData.payInfo.value.pay_type==='alipay'" label="支付宝接口加密密钥">
        <el-input v-model="payStoreData.payInfo.value.alipay.alipay_encrypt_key" placeholder="没有则不填"
                  type="textarea"/>
      </el-form-item>


      <el-form-item v-if="payStoreData.payInfo.value.pay_type==='epay'" label="易支付商户ID">
        <el-input v-model.number="payStoreData.payInfo.value.epay.epay_pid"/>
      </el-form-item>
      <el-form-item v-if="payStoreData.payInfo.value.pay_type==='epay'" label="易支付商户密钥">
        <el-input v-model="payStoreData.payInfo.value.epay.epay_key"/>
      </el-form-item>
      <el-form-item v-if="payStoreData.payInfo.value.pay_type==='epay'" label="易支付api地址">
        <el-input v-model="payStoreData.payInfo.value.epay.epay_api_url"/>
        <div style="color: #9b9da1">*页面跳转支付的地址，例如：http://abc.com/submit.php</div>
      </el-form-item>

      <!--      <el-form-item v-if="payStoreData.payInfo.value.pay_type==='epay'" label="异步通知地址">-->
      <!--        <el-input v-model="payStoreData.payInfo.value.epay.epay_notify_url"/>-->
      <!--      </el-form-item>-->
      <!--      <el-form-item v-if="payStoreData.payInfo.value.pay_type==='epay'" label="页面跳转通知地址">-->
      <!--        <el-input v-model="payStoreData.payInfo.value.epay.epay_return_url"/>-->
      <!--      </el-form-item>-->
    </el-form>
    <template #footer>
            <span class="dialog-footer">
                <el-button @click="closeDialog">取消</el-button>
                <el-button type="primary" @click="onSubmit">
                    确认
                </el-button>
            </span>
    </template>
  </el-dialog>

</template>

<script setup lang="ts">

import {usePayStore} from "/@/stores/payStore";
import {storeToRefs} from "pinia";
import {reactive} from "vue";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
import alipayLogo from '/@/assets/icon/alipay.jpeg';
import epayLogo from '/@/assets/icon/epay.png';

const payStore = usePayStore()
const payStoreData = storeToRefs(payStore)
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)

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
  ],
})

// 打开弹窗
const openDialog = (type: string, row?: any) => {
  if (type == 'add') {
    state.type = type
    state.title = "新建支付"
    state.isShowDialog = true
    //清空store中的payInfo
    payStore.wipePayInfo()

  } else {
    state.type = type
    state.title = "修改支付"
    state.isShowDialog = true
    payStoreData.payInfo.value = row
  }
}
// 关闭弹窗
const closeDialog = () => {
  state.isShowDialog = false
};

//确认提交
function onSubmit() {
  if (state.type === 'add') {
    request(apiStoreData.api.value.pay_newPay, payStoreData.payInfo.value).then((res) => {
      setTimeout(() => {
        emit('refresh');
      }, 500);
    })

  } else {
    request(apiStoreData.api.value.pay_updatePay, payStoreData.payInfo.value).then((res) => {
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