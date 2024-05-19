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
    <div class="login-right flex" >
      <div class="login-right-warp flex-margin" >
        <div class="login-right-warp-mian">
          <div class="login-right-warp-main-title" style="margin-top: clamp(1.3rem, 2vw, 1.8rem);">
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
        width:clamp(10rem, 80vw, 100rem);


        span {
          margin-left: 10px;
          font-size: 28px;
          color: #ffffff; //标题颜色
          text-shadow: #000000 0px 0 2.5px;
          font-weight: 900;
        }

        .login-left-logo-text-msg {
          font-size: 12px;
          color: #ffffff; //副标题颜色
          text-shadow: #000000 0px 0 2.5px;
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
          height: 2em;
          line-height: 35px;
          margin-top: 10px;
          font-size: clamp(1.5rem, 2vw, 1.8rem);
          text-align: center;
          letter-spacing: 1px;
          color: white;
          text-shadow: #000000 0px 0 2.5px;
          font-weight: 800  ;
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
    font-weight: 1000  ;
    text-shadow: #000000 0px 0 1.5px;
    font-size: 1.1em;
  }

.el-tabs__active-bar{
    background-color:var(--el-color-primary);
    }
.el-tabs__nav-wrap::after{
  height:0px !important
}
.el-tabs__item{
  color:rgb(255, 255, 255);
  font-weight: 500 ;
  font-size: 1.1em;
  text-shadow: #000000 0px 0 1.5px;
}

.el-link__inner{
  color: rgba(255, 255, 255, 0.400);
}

.el-input__wrapper{
  background-color: transparent;
}
.el-input__inner{
  color: white;
  text-shadow: #000000 0px 0 1px;
  font-size: 1.1em;
  font-weight: 700  ;

}

.el-form-item__error{
  font-weight: 700  ;

}

</style>