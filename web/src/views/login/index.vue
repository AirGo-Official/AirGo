<template>
  <div class="login-container flex">
    <div class="login-left">
      <div class="login-left-logo" style="">
          <img :src="themeConfig.logo_link"/>
          <div class="login-left-logo-text">
            <span>{{ getThemeConfig.globalViceTitle }}</span>
            <span class="login-left-logo-text-msg">{{ getThemeConfig.globalViceTitleMsg }}</span>
          </div>
      </div>
    </div>
    <div class="login-right flex">
      <div class="login-right-warp flex-margin">
        <div class="login-right-warp-mian">
          <div class="login-right-warp-main-title">
            <div>
              <div>{{ getThemeConfig.globalTitle }}</div>
            </div>
          </div>
          <div class="login-right-warp-main-form">
            <div v-if="!state.isScan">
              <el-tabs v-model="state.tabsActiveName">
                <el-tab-pane :label="$t('message.login.signIn')" name="login">
                  <Account/>
                </el-tab-pane>
                <el-tab-pane :label="$t('message.login.register')" name="register">
                  <Register ref="RegisterRef" @toLogin="toLogin"/>
                </el-tab-pane>
              </el-tabs>
            </div>
          </div>
        </div>
      </div>
    </div>
    <LayoutFooter class="login-footer" v-if="isFooter"/>
  </div>
</template>

<script setup lang="ts" name="loginIndex">
import {computed, defineAsyncComponent, onMounted, reactive,ref} from 'vue';
import {storeToRefs} from 'pinia';


import {useThemeConfig} from '/@/stores/themeConfig';
import {NextLoading} from '/@/utils/loading';
import {useRoute} from "vue-router";
import { Local, Session } from "/@/utils/storage";

// 引入组件
const Account = defineAsyncComponent(() => import('/@/views/login/component/account.vue'));
const Register = defineAsyncComponent(() => import('/@/views/login/component/register.vue'));
const RegisterRef = ref()
const LayoutFooter = defineAsyncComponent(() => import('/@/layout/footer/index.vue'));


// 定义变量内容
const storesThemeConfig = useThemeConfig();
const {themeConfig} = storeToRefs(storesThemeConfig);
const state = reactive({
  tabsActiveName: 'login',
  isScan: false,
});

// 获取布局配置信息
const getThemeConfig = computed(() => {
  return themeConfig.value;
});
// 设置 footer 显示/隐藏
const route = useRoute();
const isFooter = computed(() => {
  return themeConfig.value.isFooter && !route.meta.isIframe;
});
// 去登录
const toLogin = () => {
  state.tabsActiveName = "login"
}

// 页面加载时
onMounted(() => {
  let i = ''
  if (route.query.i !== undefined){
    i = route.query.i as string
  }
  Session.set('invitation', i)
  NextLoading.done();
});
</script>

<style scoped lang="scss">
.login-container {
  height: 100%;
  background: var(--el-color-white);
  background-image: url("../../assets/bgc/login-body.svg");
  /* 背景图垂直、水平均居中 */
  //background-position: center center;
  //background-position: top 300px right 0px;
  /* 背景图不平铺 */
  background-repeat: no-repeat;
  /* 当内容高度大于图片高度时，背景图像的位置相对于viewport固定 */
  background-attachment: fixed;
  /* 让背景图基于容器大小伸缩 */
  background-size: cover;
  /* 设置背景颜色，背景图加载过程中会显示背景色 */
  //background-color: var(--el-color-primary-light-5);

  .login-left {
    flex: 1;
    position: relative;
    //background-color: var(--el-color-primary-light-5);
    margin-right: 200px;

    .login-left-logo {
      display: flex;
      align-items: center;
      position: absolute;
      top: 30px;
      left: 30px;
      z-index: 1;
      animation: logoAnimation 0.3s ease;

      img {
        width: 52px;
        height: 52px;
      }

      .login-left-logo-text {
        display: flex;
        flex-direction: column;
        width: 400px;

        span {
          margin-left: 10px;
          font-size: 28px;
          color: #000000; //标题颜色
        }

        .login-left-logo-text-msg {
          font-size: 12px;
          color: rgba(134, 109, 109, 0.99); //副标题颜色
        }
      }
    }
  }

  .login-right {
    width: 500px;
    margin-right: 35%;

    .login-right-warp {
      border: 1px solid var(--el-color-primary-light-3); //表单边框
      border-radius: 3px;
      width: 100%;
      //height: 500px; //表单高度
      position: relative;
      overflow: hidden;
      background-color: var(--el-color-white);

      .login-right-warp-mian {
        display: flex;
        flex-direction: column;
        height: 100%;

        .login-right-warp-main-title {
          height: 80px;
          line-height: 80px;
          margin-top: 0px;
          font-size: 27px;
          text-align: center;
          letter-spacing: 3px;
          animation: logoAnimation 0.3s ease;
          animation-delay: 0.3s;
          color: var(--el-text-color-primary);
        }

        .login-right-warp-main-form {
          flex: 1;
          padding: 0 20px 20px;
        }
      }
    }
  }
}

.login-footer {
  position: absolute;
  left: 50%;
  top: 98%;
  transform: translate(-50%, -50%); /* 50%为自身尺寸的一半 */
}
.html{
  overflow:hidden;
}
</style>
