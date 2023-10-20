<template>
  <div class="login-container flex">
    <div class="login-left">
      <div class="login-left-logo">
        <img :src="themeConfig.logo_link"/>
        <div class="login-left-logo-text">
          <span>{{ getThemeConfig.globalViceTitle }}</span>
          <span class="login-left-logo-text-msg">{{ getThemeConfig.globalViceTitleMsg }}</span>
        </div>
      </div>
      <!--      左侧svg-->
      <div class="login-left-img">
        <img :src="loginMain"/>
      </div>
      <img :src="loginBg" class="login-left-waves"/>
    </div>
    <div class="login-right flex">
      <div class="login-right-warp flex-margin">
        <span class="login-right-warp-one"></span>
        <span class="login-right-warp-two"></span>
        <div class="login-right-warp-mian">
          <div class="login-right-warp-main-title">{{ getThemeConfig.globalTitle }}</div>
          <div class="login-right-warp-main-form">
            <div v-if="!state.isScan">
              <el-tabs v-model="state.tabsActiveName">
                <el-tab-pane label="登录" name="account">
                  <Account/>
                </el-tab-pane>
                <el-tab-pane label="注册" name="mobile">
                  <Register/>
                </el-tab-pane>
              </el-tabs>
            </div>
            <Scan v-if="state.isScan"/>
            <div class="login-content-main-sacn" @click="state.isScan = !state.isScan">
              <i class="iconfont" :class="state.isScan ? 'icon-diannao1' : 'icon-barcode-qr'"></i>
              <div class="login-content-main-sacn-delta"></div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <LayoutFooter class="login-footer" v-if="isFooter"/>
  </div>
  ->
</template>

<script setup lang="ts" name="loginIndex">
import {computed, defineAsyncComponent, onMounted, reactive} from 'vue';
import {storeToRefs} from 'pinia';

import {useThemeConfig} from '/@/stores/themeConfig';
import {NextLoading} from '/@/utils/loading';
import loginMain from '/@/assets/bgc/login-main.svg';
import loginBg from '/@/assets/bgc/login-bg.svg';
import {useRoute} from "vue-router";
import {Local} from "/@/utils/storage";

// 引入组件
const Account = defineAsyncComponent(() => import('/@/views/login/component/account.vue'));
const Register = defineAsyncComponent(() => import('/@/views/login/component/register.vue'));
const Scan = defineAsyncComponent(() => import('/@/views/login/component/scan.vue'));
const LayoutFooter = defineAsyncComponent(() => import('/@/layout/footer/index.vue'));

// 定义变量内容
const storesThemeConfig = useThemeConfig();
const {themeConfig} = storeToRefs(storesThemeConfig);
const state = reactive({
  tabsActiveName: 'account',
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

// 页面加载时
onMounted(() => {
  // console.log("route.query",  route.query.i)
  Local.set('invitation', route.query.i)
  NextLoading.done();
});
</script>

<style scoped lang="scss">
.login-container {
  height: 100%;
  background: var(--el-color-white);

  .login-left {
    flex: 1;
    position: relative;
    background-color: rgba(255, 165, 0, 1);
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

        span {
          margin-left: 10px;
          font-size: 28px;
          color: #ffffff; //标题颜色
        }

        .login-left-logo-text-msg {
          font-size: 12px;
          color: #ffffff; //副标题颜色
        }
      }
    }

    .login-left-img {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      width: 100%;
      height: 52%;

      img {
        width: 100%;
        height: 100%;
        animation: error-num 0.6s ease;
      }
    }

    .login-left-waves {
      position: absolute;
      top: 0;
      right: -100px;
    }
  }

  .login-right {
    width: 500px;
    margin-right: 100px;

    .login-right-warp {
      border: 1px solid var(--el-color-primary-light-3); //表单边框
      border-radius: 3px;
      width: 100%;
      //height: 500px; //表单高度
      position: relative;
      overflow: hidden;
      background-color: var(--el-color-white);

      .login-right-warp-one,
      .login-right-warp-two {
        position: absolute;
        display: block;
        width: inherit;
        height: inherit;

        &::before,
        &::after {
          content: '';
          position: absolute;
          z-index: 1;
        }
      }

      .login-right-warp-one {
        &::before {
          filter: hue-rotate(0deg);
          top: 0px;
          left: 0;
          width: 100%;
          height: 5px;
          background: linear-gradient(90deg, transparent, var(--el-color-primary));
          animation: loginLeft 3s linear infinite;
        }

        &::after {
          filter: hue-rotate(60deg);
          top: -100%;
          right: 2px;
          width: 5px;
          height: 100%;
          background: linear-gradient(180deg, transparent, var(--el-color-primary));
          animation: loginTop 3s linear infinite;
          animation-delay: 0.7s;
        }
      }

      .login-right-warp-two {
        &::before {
          filter: hue-rotate(120deg);
          bottom: 2px;
          right: -100%;
          width: 100%;
          height: 5px;
          background: linear-gradient(270deg, transparent, var(--el-color-primary));
          animation: loginRight 3s linear infinite;
          animation-delay: 1.4s;
        }

        &::after {
          filter: hue-rotate(300deg);
          bottom: -100%;
          left: 0px;
          width: 5px;
          height: 100%;
          background: linear-gradient(360deg, transparent, var(--el-color-primary));
          animation: loginBottom 3s linear infinite;
          animation-delay: 2.1s;
        }
      }

      .login-right-warp-mian {
        display: flex;
        flex-direction: column;
        height: 100%;

        .login-right-warp-main-title {
          height: 100px;
          line-height: 100px;
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

          .login-content-main-sacn {
            position: absolute;
            top: 0;
            right: 0;
            width: 50px;
            height: 50px;
            overflow: hidden;
            cursor: pointer;
            transition: all ease 0.3s;
            color: var(--el-color-primary);

            &-delta {
              position: absolute;
              width: 35px;
              height: 70px;
              z-index: 2;
              top: 2px;
              right: 21px;
              background: var(--el-color-white);
              transform: rotate(-45deg);
            }

            &:hover {
              opacity: 1;
              transition: all ease 0.3s;
              color: var(--el-color-primary) !important;
            }

            i {
              width: 47px;
              height: 50px;
              display: inline-block;
              font-size: 48px;
              position: absolute;
              right: 1px;
              top: 0px;
            }
          }
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

//.muyu {
//  position: absolute;
//  left: 50%;
//  top: 90%;
//  transform: translate(-50%, -50%); /* 50%为自身尺寸的一半 */
//}
</style>
