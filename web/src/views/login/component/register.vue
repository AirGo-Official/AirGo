<template>
  <el-form size="large" class="login-content-form" ref="ruleFormRef" :rules="registerRules" :model="registerData">

    <el-form-item prop="user_name">
      <el-input text :placeholder="$t('message.login.placeholder4')" v-model="registerData.user_name" clearable autocomplete="off">
        <template #prefix>
          <el-icon>
            <ele-User/>
          </el-icon>
        </template>
        <template #append>
          <el-select v-model="registerData.email_suffix" style="width: 130px;">
            <el-option
                v-for="item in publicStoreData.acceptable_email_suffixes_arr.value"
                :key="item"
                :label="item"
                :value="item"
            />
          </el-select>
        </template>
      </el-input>
    </el-form-item>

    <el-form-item prop="password">
      <el-input text :placeholder="$t('message.login.placeholder2')" v-model="registerData.password" clearable autocomplete="off">
        <template #prefix>
          <el-icon>
            <ele-Unlock/>
          </el-icon>
        </template>
      </el-input>
    </el-form-item>
    <el-form-item prop="re_password">
      <el-input text :placeholder="$t('message.login.placeholder3')" v-model="registerData.re_password" clearable autocomplete="off">
        <template #prefix>
          <el-icon>
            <ele-Unlock/>
          </el-icon>
        </template>
      </el-input>
    </el-form-item>
    <!--    base64验证码-->
    <el-form-item prop="b64s" v-if="publicStoreData.publicSetting.value.enable_base64_captcha">
      <el-col :span="10" style="display: flex;align-items: center">
        <el-input text :placeholder="$t('message.login.placeholder1')" v-model="registerData.base64_captcha.b64s" clearable autocomplete="off">
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
      <el-form-item v-if="publicStoreData.publicSetting.value.enable_email_code">
      <el-col :span="13">
        <el-input text maxlength="4" :placeholder="$t('message.login.placeholder1')" v-model="registerData.email_code" clearable
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
          {{ state.isCountDown ? `${t('message.login.retry')}:${state.countDownTime}s` : $t('message.login.codeText') }}
        </el-button>
      </el-col>
    </el-form-item>

    <el-form-item>
      <el-button @click="submitForm(ruleFormRef)" round type="primary" v-waves class="login-content-submit">
        <span>{{$t('message.login.register')}}</span>
      </el-button>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts" name="loginMobile">
import {onMounted, reactive, ref} from 'vue';
import type {FormInstance, FormRules} from 'element-plus'
import {ElMessage} from 'element-plus';
import {useUserStore} from "/@/stores/user_logic/userStore";
import {storeToRefs} from 'pinia';
import {useRouter} from 'vue-router';
import {useThemeConfig} from '/@/stores/themeConfig';
import {useAdminServerStore} from "/@/stores/admin_logic/serverStore";
import {usePublicStore} from "/@/stores/publicStore";
import {Local} from "/@/utils/storage";
import { useConstantStore } from "/@/stores/constantStore";
import { useI18n } from "vue-i18n";
const constantStore = useConstantStore()
const userStore = useUserStore()
const {registerData, loginData} = storeToRefs(userStore)
const router = useRouter();
const storesThemeConfig = useThemeConfig();
const {themeConfig} = storeToRefs(storesThemeConfig);
const serverStore = useAdminServerStore()
const publicStore = usePublicStore()
const publicStoreData = storeToRefs(publicStore)
const {t} = useI18n()
const emit = defineEmits(['toLogin'])

//定义参数
const state = reactive({
  isShowPassword: false,
  isCountDown: false,
  countDownTime: 60,
  loading: {
    signIn: false,
  },
  emailParams: {
    email_type: constantStore.EMAIL_TYPE_USER_REGISTER,
    target_email: "",
  } as EmailRequest,
});

//注册
const onRegister = () => {
  userStore.register().then((res) => {
    ElMessage.success(t('message.login.registerText'))
    Local.clear() //删除缓存（包含邀请码数据）
    userStore.$reset()
    emit('toLogin')
    refreshCaptcha() //刷新base64验证码，方便再次注册
  }).catch(() => {
    refreshCaptcha() //注册失败，刷新base64验证码，停留在当前页面
  })
}
//获取邮件验证码
const onGetEmailCode = () => {
  if (registerData.value.user_name === '') {
    return
  }
  state.emailParams.target_email = registerData.value.user_name + registerData.value.email_suffix
  publicStore.sendEmailCode(state.emailParams).then((res) => {
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
    {required: true, message: t('message.login.loginRules.msg1'), trigger: 'blur'},
    {min: 4, max: 40, message: t('message.login.loginRules.msg2'), trigger: 'blur'},
  ],
  password: [
    {required: true, message: t('message.login.loginRules.msg3'), trigger: 'blur'},
    {min: 4, max: 20, message: t('message.login.loginRules.msg4'), trigger: 'blur'},
  ],
  re_password: [
    {required: true, message: t('message.login.loginRules.msg3'), trigger: 'blur'},
    {min: 4, max: 20, message: t('message.login.loginRules.msg4'), trigger: 'blur'},
  ],
})
// 提交表单，验证表单
const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid, fields) => {
    if (valid) {
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
}

}</style>
