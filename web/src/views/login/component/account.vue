<!-- eslint-disable no-console -->
<template>
  <el-form size="large" class="login-content-form" ref="ruleFormRef" :rules="loginRules" :model="loginData">
    <el-form-item prop="user_name">
      <el-input text placeholder="邮箱" v-model="loginData.user_name" clearable
                autocomplete="off">
        <template #prefix>
          <el-icon>
            <ele-User/>
          </el-icon>
        </template>
      </el-input>
    </el-form-item>
    <el-form-item prop="password">
      <el-input :type="state.isShowPassword ? 'text' : 'password'" placeholder="密码"
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
        <el-input text maxlength="4" placeholder="请输入验证码" v-model="loginData.email_code" clearable
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
        <el-button class="login-content-code" type="primary" :disabled="state.isCountDown || loginData.user_name === ''"
                   @click="onGetEmailCode">
          {{ state.isCountDown ? `${state.countDownTime}s后重新获取` : "获取验证码" }}
        </el-button>
      </el-col>
    </el-form-item>

    <el-form-item>
      <el-col :span="11">
        <el-button v-if="!state.enableResetPassword" type="primary" class="login-content-submit" round
                   @click="submitForm(ruleFormRef)">
          <span>登 录</span>
        </el-button>
      </el-col>
      <el-col :span="2">
      </el-col>
      <el-col :span="11">
        <el-button v-if="!state.enableResetPassword" class="login-content-resetPassword" round @click="onResetPassword">
          <span>重置密码</span>
        </el-button>
      </el-col>
      <el-col :span="11">
        <el-button v-if="state.enableResetPassword" @click="onSubmitResetPassword" class="login-content-resetPassword"
                   round type="danger">
          <span>确认重置密码</span>
        </el-button>
      </el-col>
      <el-col :span="2">
      </el-col>
      <el-col :span="11">
        <el-button v-if="state.enableResetPassword" @click="state.enableResetPassword=false"
                   class="login-content-resetPassword" round type="primary">
          <span>返回登录</span>
        </el-button>
      </el-col>
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
import {useUserStore} from "/@/stores/userStore";
import {useThemeConfig} from '/@/stores/themeConfig';
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";

const route = useRoute();
const router = useRouter();
const userStore = useUserStore()
const {loginData} = storeToRefs(userStore)
const storesThemeConfig = useThemeConfig();
const {themeConfig} = storeToRefs(storesThemeConfig);
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
//定义参数
const state = reactive({
  enableResetPassword: false,
  isShowPassword: false,
  isCountDown: false,
  countDownTime: 60,
  loading: {
    signIn: false,
  },
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
  //state.loading.signIn = true;
  await userStore.userLogin(loginData.value)
  //添加完动态路由，再进行 router 跳转，否则可能报错 No match found for location with path "/"
  const isNoPower = await initBackEndControlRoutes();
  //执行完 initBackEndControlRoutes，再执行 signInSuccess
  signInSuccess(isNoPower);

};
// 登录成功后的跳转
const signInSuccess = (isNoPower: boolean | undefined) => {
  console.log("3")
  if (isNoPower) {
    ElMessage.warning('抱歉，您没有登录权限');
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
      router.push('/');
    }
    // 登录成功提示
    const signInText = '欢迎回来！';
    ElMessage.success(`${currentTimeInfo}，${signInText}`);
    // 添加 loading，防止第一次进入界面时出现短暂空白
    NextLoading.start();
  }
  //state.loading.signIn = false;
};
//获取邮箱验证码
const onGetEmailCode = () => {
  if (loginData.value.user_name === '') {
    return
  }
  state.isCountDown = true
  request(apiStoreData.staticApi.value.public_getEmailCode, loginData.value).then((res) => {
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
    {required: true, message: '请输入用户名', trigger: 'blur'},
    {min: 4, max: 40, message: '长度4～40', trigger: 'blur'},
  ],
  password: [
    {required: true, message: '请输入密码', trigger: 'blur'},
    {min: 4, max: 20, message: '密码长度4～20', trigger: 'blur'},
  ],
})

const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid, fields) => {
    if (valid) {
      console.log('submit!')
      onSignIn()
    } else {
      console.log('error submit!', fields)
    }
  })
}


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
    //	letter-spacing: 5px;
  }

  .login-content-submit {
    width: 100%;
    letter-spacing: 2px;
    font-weight: 300;
  }

  .login-content-resetPassword {
    width: 100%;
    letter-spacing: 2px;
    font-weight: 300;
    margin-left: 0;
  }
}</style>
