<template>
  <el-form size="large" class="login-content-form" ref="ruleFormRef" :rules="registerRules" :model="registerData">

    <el-form-item prop="user_name">
      <el-input text placeholder="请输入邮箱" v-model="registerData.user_name" clearable autocomplete="off">
        <template #prefix>
          <el-icon>
            <ele-User/>
          </el-icon>
        </template>
        <template #append>
          <el-select v-model="registerData.email_suffix" style="width: 130px;">
            <el-option label="@qq.com" value="@qq.com"/>
            <el-option label="@gmail.com" value="@gmail.com"/>
            <el-option label="@163.com" value="@163.com"/>
            <el-option label="@126.com" value="@126.com"/>
            <el-option label="@outlook.com" value="@outlook.com"/>
            <el-option label="@hotmail.com" value="@hotmail.com"/>
            <el-option label="@foxmail.com" value="@foxmail.com"/>
          </el-select>
        </template>
      </el-input>
    </el-form-item>

    <el-form-item prop="password">
      <el-input text placeholder="请输入密码" v-model="registerData.password" clearable autocomplete="off">
        <template #prefix>
          <el-icon>
            <ele-Unlock/>
          </el-icon>
        </template>
      </el-input>
    </el-form-item>
    <el-form-item prop="re_password">
      <el-input text placeholder="重新输入密码" v-model="registerData.re_password" clearable autocomplete="off">
        <template #prefix>
          <el-icon>
            <ele-Unlock/>
          </el-icon>
        </template>
      </el-input>
    </el-form-item>
    <!--    base64验证码-->
    <el-form-item prop="b64s">
      <el-col :span="10" style="display: flex;align-items: center">
        <el-input text placeholder="输入验证码" v-model="registerData.base64_captcha.b64s" clearable autocomplete="off">
          <template #prefix>
            <el-icon>
              <ele-Position/>
            </el-icon>
          </template>
        </el-input>
      </el-col>
      <el-col :span="4"></el-col>
      <el-col :span="10">
        <img :src="publicStore.base64CaptchaData.b64s" @click="refreshCaptcha"/>
      </el-col>

    </el-form-item>

    <el-form-item v-if="publicServerConfig.enable_email_code">
      <el-col :span="13">
        <el-input text maxlength="4" placeholder="请输入验证码" v-model="registerData.email_code" clearable
                  autocomplete="off">
          <template #prefix>
            <el-icon>
              <ele-Position/>
            </el-icon>
          </template>
        </el-input>
      </el-col>
      <el-col :span="1"></el-col>
      <el-col :span="10">
        <el-button class="login-content-code" type="primary" :disabled="state.isCountDown" @click="onGetEmailCode">
          {{ state.isCountDown ? `${state.countDownTime}s后重新获取` : "获取验证码" }}
        </el-button>
      </el-col>
    </el-form-item>

    <el-form-item>
      <el-button @click="submitForm(ruleFormRef)" round type="primary" v-waves class="login-content-submit">
        <span>注 册</span>
      </el-button>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts" name="loginMobile">
import {onMounted, reactive, ref} from 'vue';
import type {FormInstance, FormRules} from 'element-plus'
import {ElMessage} from 'element-plus';
import {useUserStore} from "/@/stores/userStore";
import {storeToRefs} from 'pinia';

import {useRouter} from 'vue-router';

import {useThemeConfig} from '/@/stores/themeConfig';

import {useServerStore} from "/@/stores/serverStore";

import {usePublicStore} from "/@/stores/publicStore";

import {Local} from "/@/utils/storage";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";

const userStore = useUserStore()
const {registerData, loginData} = storeToRefs(userStore)
const router = useRouter();
const storesThemeConfig = useThemeConfig();
const {themeConfig} = storeToRefs(storesThemeConfig);
const serverStore = useServerStore()
const {publicServerConfig} = storeToRefs(serverStore)
const publicStore = usePublicStore()

const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)


//定义参数
const state = reactive({
  isShowPassword: false,
  isCountDown: false,
  countDownTime: 60,
  loading: {
    signIn: false,
  },
});

//注册
const onRegister = () => {
  userStore.register().then((res) => {
    ElMessage.success('注册成功，前往登录...')
    Local.clear() //删除缓存（包含邀请码数据）
    setTimeout(() => {
      window.location.href = '/'; // 去登录页
      //router.push('/'); // 去登录页
    }, 500)

  }).catch(() => {
    refreshCaptcha()
  })
}
//获取邮件验证码
const onGetEmailCode = () => {
  if (registerData.value.user_name === '') {
    return
  }
  // publicApi.getEmailCodeApi(userStore.registerData).then((res) => {
  request(apiStoreData.staticApi.value.public_getEmailCode, userStore.registerData).then((res) => {
    state.isCountDown = true
    ElMessage.success(res.msg)
    handleTimeChange()
  })
};
//倒计时
const handleTimeChange = () => {
  if (state.countDownTime <= 0) {
    state.isCountDown = false;
    state.countDownTime = 60;
  } else {
    setTimeout(() => {
      state.countDownTime--;
      handleTimeChange();
    }, 500);
  }
};

//表单校验
const ruleFormRef = ref<FormInstance>()
const registerRules = reactive<FormRules<RegisterForm>>({
  user_name: [
    {required: true, message: '请输入用户名', trigger: 'blur'},
    {min: 4, max: 40, message: '长度4～40', trigger: 'blur'},
  ],
  password: [
    {required: true, message: '请输入密码', trigger: 'blur'},
    {min: 4, max: 20, message: '密码长度4～20', trigger: 'blur'},
  ],
  re_password: [
    {required: true, message: '请输入密码', trigger: 'blur'},
    {min: 4, max: 20, message: '密码长度4～20', trigger: 'blur'},
  ],
})
// 提交表单，验证表单
const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid, fields) => {
    if (valid) {
      console.log('submit!')
      onRegister()
    } else {
      console.log('error submit!', fields)
    }
  })
}
// 刷新base64Captcha
const refreshCaptcha = () => {
  publicStore.getBase64Captcha()
}

onMounted(() => {
  publicStore.getBase64Captcha()
});
</script>


<style scoped lang="scss">
.login-content-form {
  margin-top: 20px;
  @for $i from 1 through 4 {
    .login-animation#{$i} {
      opacity: 0;
      animation-name: error-num;
      animation-duration: 0.5s;
      animation-fill-mode: forwards;
      animation-delay: calc($i/10) + s;
    }
  }

  .login-content-code {
    width: 100%;
    padding: 0;
  }

  .login-content-submit {
    width: 100%;
    letter-spacing: 2px;
    font-weight: 300;
  }

}</style>
