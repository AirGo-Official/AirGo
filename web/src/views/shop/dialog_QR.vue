<template>
  <el-dialog v-model="state.isShowQRDialog" title="扫描二维码支付" width="80%">
    <div>
      <el-card shadow="hover">
        <div class="qrcode-img-warp">
          <div class="mb30 mt30 qrcode-img">
            <!-- 二维码弹窗 -->
            <div id="qrcode" class="qrcode" ref="qrcodeRef"></div>
          </div>
        </div>
      </el-card>
    </div>
    <template #footer>
            <span class="dialog-footer">
                <el-button @click="state.isShowQRDialog = false">取消</el-button>
                <el-button type="primary" @click="state.isShowQRDialog = false">
                    已支付
                </el-button>
            </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import {nextTick, reactive, ref} from "vue";
import QRCode from 'qrcodejs2-fixes';
import {storeToRefs} from 'pinia';
import {useShopStore} from "/@/stores/shopStore";

const qrcodeRef = ref();
const shopStore = useShopStore()
const {shopData} = storeToRefs(shopStore)
//定义变量
const state = reactive({
  isShowQRDialog: false,
  QRcode: null,
})
//二维码
const onInitQrcode = () => {
  //清除上一次二维码
  const codeHtml = document.getElementById("qrcode");
  codeHtml.innerHTML = "";
  state.QRcode = new QRCode(qrcodeRef.value, {
    text: shopData.value.currentOrder.pay_info.alipay_info.qr_code,
    width: 125,
    height: 125,
    colorDark: '#000000',
    colorLight: '#ffffff',
  });
}
//打开弹窗
const openDialog = () => {
  state.isShowQRDialog = true
  nextTick(() => {
    onInitQrcode()
  })
}

defineExpose({
  openDialog,
})
</script>

<style scoped></style>