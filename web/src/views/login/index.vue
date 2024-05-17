<template>
  <div class="login-container flex">
    <div class="login-left">
      <div class="login-left-logo" style="text-shadow: #000000 0px 0 3px;">
          <img src="https://ice.frostsky.com/2024/03/23/f9c5cb33538da41e490d98e6051ee556.png"/>
          <div class="login-left-logo-text">
            <span class="">風嶼Link</span>
            <span class="login-left-logo-text-msg">「正在不断进化中...」</span>
          </div>
      </div>
    </div>
    <div class="login-right flex" >
      <div class="login-right-warp flex-margin" >
        <div class="login-right-warp-mian">
          <div class="login-right-warp-main-title" style="margin-top: clamp(1.3rem, 2vw, 1.8rem);">
            <div>
              <h3 style="font-size: clamp(1.5rem, 2vw, 1.8rem);" ><i class="ri-chat-1-fill"></i> 有朋自远方来，不亦悦乎。<br>
              (=・ω・=)</h3>
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
  //获取aff  http://localhost:8080/#/login?aff=12345
  let aff = ''
  if (route.query.aff){
    aff = route.query.aff as string
    // console.log("邀请码：",aff)
    Session.set('invitation', aff)
  }
  NextLoading.done();
});
</script>

<style scoped lang="scss">
.login-container {
  height: 100%;
  background-image: url("../../assets/bgc/login-body.svg");
  background-repeat: no-repeat;
  background-size: cover;
  background-attachment: fixed;
  
  .login-left {
    flex: 1;
    position: relative;
    margin-right: 200px;

    .login-left-logo {
      display: flex;
      align-items: center;
      position: absolute;
      top: 30px;
      left: 30px;
      z-index: 1;
     // animation: logoAnimation 0.3s ease;

      img {
        width: 51px;
        height: 51px;
      }

      .login-left-logo-text { //logo
        display: flex;
        flex-direction: column;
        min-width:200px;
        max-width: 300px;

        span {
          margin-left: 10px;
          font-size: 28px;
          color: #ffffff; //标题颜色
          text-shadow: #000000 0px 0 3px;
          font-weight: 900;
        }

        .login-left-logo-text-msg {
          font-size: 12px;
          color: #ffffff; //副标题颜色
          text-shadow: #000000 0px 0 3px;
          font-weight: 700  ;
      } }
    }
  }


  .login-right {
    width: 500px;
    margin-right: 35%;
    

    .login-right-warp {
      //border: 1px solid var(--el-color-primary-light-3); //表单边框
      border-radius: 9px;
      width: 100%;
      
      //height: 500px; //表单高度
      position: relative;
      overflow: hidden;
      background: rgba(219, 227, 255, 0.245);
	    -webkit-backdrop-filter: blur(10px);
	    backdrop-filter: blur(10px);
      box-shadow: 0px 0px 45px rgba(0, 0, 0, 0.486);

      
      .login-right-warp-mian {
        display: flex;
        flex-direction: column;
        height: 100%;
        
        
        .login-right-warp-main-title {
          height: 80px;
          line-height: 35px;
          margin-top: 10px;
          font-size: 27px;
          text-align: center;
          letter-spacing: 2px;
          // animation: logoAnimation 0.3s ease;
          //animation-delay: 0.3s;
          color: white;
          text-shadow: #000000 0px 0 4px;
          font-weight: 700  ;
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
  position: fixed;
  left: 50%;
  top: 96%;
  transform: translate(-50%, -50%); /* 50%为自身尺寸的一半 */
}
.html{
  overflow:hidden;
}

</style>
<style>
.el-tabs__item.is-active{
    color:rgb(255, 255, 255);
    font-weight: 700  ;

  }

.el-tabs__active-bar{
    background-color:rgb(255, 255, 255);
    }
.el-tabs__nav-wrap::after{
  height:0px !important
}
.el-tabs__item{
  color: rgba(255, 255, 255, 0.777);
}
.el-link__inner{
  color: rgba(255, 255, 255, 0.400);
}
.el-input__wrapper{
  background-color: transparent;
}
.el-input__inner{
  color: white;
}

</style>