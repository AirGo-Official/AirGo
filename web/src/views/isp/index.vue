<template>
  <div class="container layout-padding">
    <!--    登录卡片-->
    <el-card shadow="hover" class="layout-padding-auto"
             v-if="ispStoreData.isp.value.unicom_config.cookie==='' && ispStoreData.isp.value.telecom_config.telecomToken===''">
      <div>
        <el-tabs type="border-card">
          <el-tab-pane label="联通">
            <el-form>
              <el-form-item>
                <el-input v-model="ispStoreData.isp.value.mobile" placeholder="输入联通手机号"></el-input>
              </el-form-item>
              <el-form-item>
                <el-col :span="18">
                  <el-input v-model="ispStoreData.isp.value.unicom_config.password"
                            placeholder="输入收到的验证码"></el-input>
                </el-col>
                <el-col :span="6">
                  <el-button @click="sendCode(ispStoreData.isp.value,'unicom')">{{
                      ispStoreData.isCountDown.value ? `${ispStoreData.countDownTime.value}s后重新获取` : "获取验证码"
                    }}
                  </el-button>
                </el-col>
              </el-form-item>
              <el-form-item>
                <div style="text-align: center;width: 100%">
                  <el-button @click="ispLogin(ispStoreData.isp.value,'unicom')">登录</el-button>
                </div>
              </el-form-item>
            </el-form>
          </el-tab-pane>
          <el-tab-pane label="电信">
            <el-form>
              <el-form-item>
                <el-input v-model="ispStoreData.isp.value.mobile" placeholder="输入电信手机号"></el-input>
              </el-form-item>
              <el-form-item>
                <el-input v-model="ispStoreData.isp.value.telecom_config.telecomPassword"
                          placeholder="输入登录密码"></el-input>
              </el-form-item>
              <el-form-item>
                <div style="text-align: center;width: 100%">
                  <el-button @click="ispLogin(ispStoreData.isp.value,'telecom')">登录</el-button>
                </div>
              </el-form-item>
            </el-form>
          </el-tab-pane>
        </el-tabs>
      </div>
    </el-card>
    <!--    监控卡片-->
    <el-card shadow="hover" class="layout-padding-auto"
             v-if="ispStoreData.isp.value.unicom_config.cookie !=='' || ispStoreData.isp.value.telecom_config.telecomToken !==''">
      <div class="card-text">
        <el-tag class="card-text-left" type="warning">运营商</el-tag>
        <span class="card-text-right" v-if="ispStoreData.isp.value.isp_type==='unicom'">联通</span>
        <span class="card-text-right" v-if="ispStoreData.isp.value.isp_type==='telecom'">电信</span>
      </div>
      <div class="card-text">
        <el-tag class="card-text-left" type="warning">号码</el-tag>
        <span class="card-text-right">{{ ispStoreData.isp.value.mobile }}</span>
      </div>
      <div class="card-text">
        <el-tag class="card-text-left" type="warning">状态</el-tag>
        <span class="card-text-right">{{ ispStoreData.isp.value.status ? '正在监控ing' : '请重新登录' }}</span>
      </div>
      <div class="card-text">
        <el-button class="card-text-left" type="warning" @click="loginAgain()">重新登录</el-button>
        <el-button class="card-text-left" type="warning" @click="copyUrl()">复制url</el-button>
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {useISPStore} from "/@/stores/ispStore";
import {storeToRefs} from "pinia";
import {onMounted} from "vue";
import {TelecomMobileHandler, TelecomRSAEncrypt} from "/@/utils/encrypt"
import commonFunction from '/@/utils/commonFunction';
import {Local} from "/@/utils/storage";
import {getFormatDate, randomString} from "/@/utils/formatTime";

const {copyText} = commonFunction();
const ispStore = useISPStore()
const ispStoreData = storeToRefs(ispStore)
import {useApiStore} from "/@/stores/apiStore";
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)

//获取验证码
const sendCode = (params: Isp, isp_type: string) => {
  params.isp_type = isp_type

  switch (params.isp_type) {
    case "unicom":
    case "telecom":
      const num = TelecomMobileHandler(params.mobile)
      params.telecom_config.phoneNum = num

  }
  // console.log("params:", params)
  ispStore.sendCode(params)
  ispStoreData.isCountDown.value = true
  handleTimeChange()

}

// 电信loginAuthCipherAsymmertric字段解密感谢 huangqikang511@github 技术指导
const ispLogin = (params: Isp, isp_type: string) => {
  params.isp_type = isp_type
  switch (params.isp_type) {
    case "unicom":
    case "telecom":
      // 判断是否有deviceUid
      if (params.telecom_config.deviceUid === '') {
        // console.log("deviceUid为空")
        params.telecom_config.deviceUid = randomString(16)
      }
      //处理手机号
      params.telecom_config.phoneNum = TelecomMobileHandler(params.mobile)
      //处理loginAuthCipherAsymmertric
      const tm = getFormatDate("yyyyMMddHHmm00")
      const au = `iPhone 14 15.4.${params.telecom_config.deviceUid.slice(0, 12)}${params.mobile}${tm}${params.telecom_config.telecomPassword}0$$$0.`
      // console.log("au:",au)
      params.telecom_config.loginAuthCipherAsymmertric = TelecomRSAEncrypt(au)
      //时间戳
      params.telecom_config.timestamp = tm

  }
  // console.log("登录:", params)
  ispStore.ispLogin(params)
}
//复制url
const copyUrl = () => {
  const url = import.meta.env.VITE_API_URL + apiStoreData.staticApi.value.public_queryPackage.path + "?id=" + Local.get("token")
  // console.log("url:", url)
  copyText(url)
}
//重新登录
const loginAgain = () => {
  ispStore.ispLogin({isp_type: "loginAgain"})
}

//倒计时
const handleTimeChange = () => {
  if (ispStoreData.countDownTime.value <= 0) {
    ispStoreData.isCountDown.value = false;
    ispStoreData.countDownTime.value = 60;
  } else {
    setTimeout(() => {
      ispStoreData.countDownTime.value--;
      handleTimeChange();
    }, 1000);
  }
};


onMounted(() => {
  ispStore.getMonitorByUserID()
});


</script>

<style scoped>

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.home-card-item {
  width: 100%;
  height: 100%;
  border-radius: 4px;
  transition: all ease 0.3s;
  padding: 20px;
  overflow: hidden;
  background: var(--el-color-white);
  color: var(--el-text-color-primary);
  border: 1px solid var(--next-border-color-light);
}

.el-card {
  background-image: url("../../assets/bgc/bg-1.svg");
  background-repeat: no-repeat;
  background-position: 100%, 100%;
  /*background: rgba(0, 0, 0, 0.3);*/
}

.card-text {
  display: flex;
  justify-content: space-between;
  height: 60px
}

.card-text-left {
  margin-top: auto;
  margin-bottom: auto;
}

.card-text-right {
  margin-top: auto;
  margin-bottom: auto;
  font-size: 30px;
}

.card-header-left {
  font-size: 15px;
  color: #AC96F1;
}
</style>