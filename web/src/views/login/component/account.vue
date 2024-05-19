<!-- eslint-disable no-console -->
<template>
  <el-form size="large" class="login-content-form" ref="ruleFormRef" :rules="loginRules" :model="loginData">
    <el-form-item prop="user_name">
      <el-input text :placeholder="$t('message.adminUser.SysUser.user_name')" v-model="loginData.user_name" clearable
                autocomplete="off">
        <template #prefix>
          <el-icon>
            <ele-User/>
          </el-icon>
        </template>
      </el-input>
    </el-form-item>
    <el-form-item prop="password">
      <el-input :type="state.isShowPassword ? 'text' : 'password'" :placeholder="$t('message.adminUser.SysUser.password')"
                v-model="loginData.password" autocomplete="off">
        <template #prefix>
          <el-icon>
            <ele-Unlock/>
          </el-icon>
        </template>
        <template #suffix>
          <i class="iconfont  login-content-password"
             :class="state.isShowPassword ? 'icon-yincangmima' : 'icon-xianshimima'"
             @click="state.isShowPassword = !state.isShowPassword">
          </i>
        </template>
      </el-input>
    </el-form-item>

    <el-form-item v-if="state.enableResetPassword">
      <el-col :span="13">
        <el-input text maxlength="4" :placeholder="$t('message.login.placeholder1')" v-model="loginData.email_code" clearable
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
        <el-button class="login-content-code"
                   type="primary"
                   :disabled="state.isCountDown || loginData.user_name === ''"
                   @click="onGetEmailCode">
          {{ state.isCountDown ? `$t('message.login.retry'):${state.countDownTime}s` : $t('message.login.codeText') }}
        </el-button>
      </el-col>
    </el-form-item>

    <el-form-item>
      <el-col :span="11">
        <el-button v-if="!state.enableResetPassword"
                   type="primary"
                   class="login-content-submit"
                   @click="submitForm(ruleFormRef)">
          <span>{{$t('message.login.signIn')}}</span>
        </el-button>
      </el-col>
      <el-col :span="2">
      </el-col>
      <el-col :span="11">
        <el-button v-if="!state.enableResetPassword"
                   class="login-content-resetPassword"
                   @click="onResetPassword">
          <span>{{$t('message.login.resetPassword')}}</span>
        </el-button>
      </el-col>
      <el-col :span="11">
        <el-button v-if="state.enableResetPassword"
                   @click="onSubmitResetPassword"
                   class="login-content-resetPassword"
                   type="danger">
          <span>{{$t('message.login.resetPassword')}}</span>
        </el-button>
      </el-col>
      <el-col :span="2">
      </el-col>
      <el-col :span="11">
        <el-button v-if="state.enableResetPassword"
                   @click="state.enableResetPassword=false"
                   class="login-content-resetPassword"
                   type="primary">
          <span>{{$t('message.login.signIn')}}</span>
        </el-button>
      </el-col>
    </el-form-item>
    <el-form-item>
      <!--    语言切换-->
      <el-dropdown :show-timeout="70" :hide-timeout="50" trigger="click" @command="onLanguageChange">
        <div>
            <i style="font-size: 25px; color:white;          text-shadow: #000000 0px 0 3px;" class="ri-global-line"></i>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="zh-cn" :disabled="state.disabledI18n === 'zh-cn'">简体中文</el-dropdown-item>
            <el-dropdown-item command="zh-tw" :disabled="state.disabledI18n === 'zh-tw'">繁體中文</el-dropdown-item>
            <el-dropdown-item command="en" :disabled="state.disabledI18n === 'en'">English</el-dropdown-item>

          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts" name="loginAccount">
import {computed, reactive, ref} from 'vue';
import {ElMessage, FormInstance, FormRules} from 'element-plus';
import {Session} from '/@/utils/storage';
import {formatAxis} from '/@/utils/formatTime';
import {NextLoading} from '/@/utils/loading';
import {initBackEndControlRoutes} from '/@/router/backEnd';
import {useRoute, useRouter} from 'vue-router';
import {storeToRefs} from 'pinia';
import {useUserStore} from "/@/stores/user_logic/userStore";
import {useThemeConfig} from '/@/stores/themeConfig';
import {useApiStore} from "/@/stores/apiStore";
import { useConstantStore } from "/@/stores/constantStore";
import { usePublicStore } from "/@/stores/publicStore";
import other from "/@/utils/other";
import { useI18n } from "vue-i18n";

const route = useRoute();
const router = useRouter();
const constantStore = useConstantStore()
const publicStore = usePublicStore()
const userStore = useUserStore()
const {loginData} = storeToRefs(userStore)
const storesThemeConfig = useThemeConfig();
const {themeConfig} = storeToRefs(storesThemeConfig);
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
const { locale, t } = useI18n();
//定义参数
const state = reactive({
  enableResetPassword: false,
  isShowPassword: false,
  isCountDown: false,
  countDownTime: 60,
  loading: {
    signIn: false,
  },
  disabledI18n: "zh-cn",
});

// 时间获取
const currentTime = computed(() => {
  return formatAxis(new Date());
});
//重置密码
const onResetPassword = () => {
  state.enableResetPassword = true
}
//确认重置密码
const onSubmitResetPassword = () => {
  userStore.submitResetPassword().then((res) => {
    ElMessage.success(res.msg)
  })
}
// 登录
const onSignIn = async () => {
  await userStore.userLogin(loginData.value)
  //添加完动态路由，再进行 router 跳转，否则可能报错 No match found for location with path "/"
  const isNoPower = await initBackEndControlRoutes();
  //执行完 initBackEndControlRoutes，再执行 signInSuccess
  signInSuccess(isNoPower);
};
// 登录成功后的跳转
const signInSuccess = (isNoPower: boolean | undefined) => {
  if (isNoPower) {
    ElMessage.warning(t('message.staticRoutes.noPower'));
    Session.clear();
  } else {
    // 初始化登录成功时间问候语
    let currentTimeInfo = currentTime.value;
    // 登录成功，跳到转首页
    // 如果是复制粘贴的路径，非首页/登录页，那么登录成功后重定向到对应的路径中
    if (route.query?.redirect) {
      router.push({
        path: <string>route.query?.redirect,
        query: Object.keys(<string>route.query?.params).length > 0 ? JSON.parse(<string>route.query?.params) : '',
      });
    } else {
      router.push('/home');
    }
    // 登录成功提示
    const signInText = t('message.login.signInText');
    ElMessage.success(`${currentTimeInfo}，${signInText}`);
    // 添加 loading，防止第一次进入界面时出现短暂空白
    NextLoading.start();
  }
};
//获取邮箱验证码
const onGetEmailCode = () => {
  if (loginData.value.user_name === '') {
    return
  }
  state.isCountDown = true
  publicStore.sendEmailCode({email_type:constantStore.EMAIL_TYPE_USER_RESETPWD,target_email:loginData.value.user_name} as EmailRequest).then((res)=>{
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
    }, 1000);
  }
};

//表单验证
const ruleFormRef = ref<FormInstance>()
const loginRules = reactive<FormRules<RegisterForm>>({
  user_name: [
    {required: true, message: t('message.login.loginRules.msg1'), trigger: 'blur'},
    {min: 4, max: 40, message: t('message.login.loginRules.msg2'), trigger: 'blur'},
  ],
  password: [
    {required: true, message: t('message.login.loginRules.msg3'), trigger: 'blur'},
    {min: 4, max: 20, message: t('message.login.loginRules.msg4'), trigger: 'blur'},
  ],
})

const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid, fields) => {
    if (valid) {
      onSignIn()
    } else {
      console.log('error submit!', fields)
    }
  })
}
// 初始化组件大小/i18n
const initI18nOrSize = (value: string, attr: string) => {
  (<any>state)[attr] = Session.get("themeConfig")[value];
};
// 语言切换
const onLanguageChange = (lang: string) => {
  Session.remove("themeConfig");
  themeConfig.value.globalI18n = lang;
  Session.set("themeConfig", themeConfig.value);
  locale.value = lang;
  other.useTitle();
  initI18nOrSize("globalI18n", "disabledI18n");
};


</script>

<style scoped lang="scss">
.login-content-form {
  margin-top: 10px;

  @for $i from 1 through 4 {
    .login-animation#{$i} {
      opacity: 0;
      animation-name: error-num;
      animation-duration: 0.5s;
      animation-fill-mode: forwards;
      animation-delay: calc($i/10) + s;
    }
  }

  .login-content-password {
    display: inline-block;
    width: 20px;
    cursor: pointer;

    &:hover {
      color: #909399;
    }
  }

  .login-content-code {
    width: 100%;
    padding: 0;
    font-weight: bold;
  }

  .login-content-submit {
    width: 100%;

  }

  .login-content-resetPassword {
    width: 100%;
  }

}
</style>