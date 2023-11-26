<template>
  <div>
    <el-row :gutter="10" class="home-card-two mb15">
<!--      套餐详情-->
      <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
        <div class="home-card-item">
          <el-card class="box-card" style="height: 300px">
            <template #header>
              <div class="card-header">
                <el-text class="card-header-left">套餐详情</el-text>
                <el-text size="large" type="primary">{{userInfos.subscribe_info.goods_subject}}</el-text>
              </div>
            </template>
            <div class="card-text">
              <el-tag class="card-text-left" type="info">订阅状态:</el-tag>
              <el-button type="primary" class="card-text-right" style="font-size:20px" v-if="userInfos.subscribe_info.sub_status">有效</el-button>
              <el-button type="info" class="card-text-right" style="font-size:20px"  v-else>已过期</el-button>
            </div>
            <div class="card-text">
              <el-tag class="card-text-left" type="info">剩余流量:</el-tag>
              <span class="card-text-right">{{
                  ((userInfos.subscribe_info.t - userInfos.subscribe_info.d - userInfos.subscribe_info.u) / 1024 / 1024 / 1024).toFixed(2)
                }}GB</span>
            </div>
            <div class="card-text">
              <el-tag class="card-text-left" type="info">到期时间:</el-tag>
              <span class="card-text-right">{{ DateStrtoTime(userInfos.subscribe_info.expired_at).slice(0, 10) }}</span>
            </div>
          </el-card>
        </div>
      </el-col>
      <!--      订阅连接-->
      <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
        <div class="home-card-item">
          <el-card class="box-card" style="height: 300px">
            <template #header>
              <div class="card-header">
                <el-text class="card-header-left">订阅地址</el-text>
                <el-button type="primary" size="large" text plain class="button" @click="onResetSub">重置订阅链接</el-button>
              </div>
            </template>
            <div style="text-align: center">
              <div>
                <el-button size="large" color="blue" style="margin-top: 10px;margin-bottom: 10px;width: 100%" @click="Sub('')">
                  复制通用订阅
                </el-button>
              </div>
              <div>
                <el-button size="large" color="deeppink" style="margin-top: 10px;margin-bottom: 10px;width: 100%" @click="QRSub()">
                  二维码订阅
                </el-button>
              </div>
              <div>
                <el-button size="large" color="brown" style="margin-top: 10px;margin-bottom: 10px;width: 100%" @click="state.isShowSubDialog=true">
                  手动选择订阅
                </el-button>
              </div>
            </div>
          </el-card>
        </div>
      </el-col>

<!--      自定义内容-->
      <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24">
        <div class="home-card-item" >
          <el-card class="box-card">
            <div class="box-card">
              <v-md-preview :text="articleStoreData.articleID1.value.content"></v-md-preview>
            </div>
          </el-card>
        </div>
      </el-col>
    </el-row>
    <el-dialog v-model="state.isShowDialog" :title="state.title" width="80%" destroy-on-close center>
      <v-md-preview :text="articleStoreData.articleID2.value.content"></v-md-preview>
    </el-dialog>
    <el-dialog v-model="state.isShowSubDialog">
      <div>
        <el-button size="large" color="blue" style="margin-top: 10px;margin-bottom: 10px;width: 100%" @click="Sub('NekoBox')">
          NekoBox 订阅
        </el-button>
      </div>
      <div>
        <el-button size="large" color="blue" style="margin-top: 10px;margin-bottom: 10px;width: 100%" @click="Sub('v2rayNG')">
          v2rayNG 订阅
        </el-button>
      </div>
      <div>
        <el-button size="large" color="blue" style="margin-top: 10px;margin-bottom: 10px;width: 100%" @click="Sub('v2rayN')">
          v2rayN 订阅
        </el-button>
      </div>
      <div>
        <el-button size="large" color="blue" style="margin-top: 10px;margin-bottom: 10px;width: 100%" @click="Sub('Clash')">
          Clash 订阅
        </el-button>
      </div>
      <div>
        <el-button size="large" color="blue" style="margin-top: 10px;margin-bottom: 10px;width: 100%" @click="Sub('Shadowrocket')">
          Shadowrocket 订阅
        </el-button>
      </div>
      <div>
        <el-button size="large" color="blue" style="margin-top: 10px;margin-bottom: 10px;width: 100%" @click="Sub('Surge')">
          Surge 订阅
        </el-button>
      </div>
      <div>
        <el-button size="large" color="blue" style="margin-top: 10px;margin-bottom: 10px;width: 100%" @click="Sub('Quantumult')">
          Quantumult 订阅
        </el-button>
      </div>
      <div>
        <el-button size="large" color="blue" style="margin-top: 10px;margin-bottom: 10px;width: 100%" @click="Sub('V2rayU')">
          V2rayU 订阅
        </el-button>
      </div>



    </el-dialog>
    <el-dialog v-model="state.isShowQRSubDialog" width="350px">
      <!-- 二维码弹窗 -->
      <div >
        <div id="qrcode" ref="qrcodeRef"></div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">

import {DateStrtoTime} from "/@/utils/formatTime";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
import {storeToRefs} from "pinia";
import {useUserStore} from "/@/stores/userStore";
import {onMounted, reactive, ref} from 'vue';
import {ElMessage} from 'element-plus';
import {Select} from '@element-plus/icons-vue'
import commonFunction from '/@/utils/commonFunction';
import {useArticleStore} from "/@/stores/articleStore";
import QRCode from "qrcodejs2-fixes";

const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
const userStore = useUserStore()
const {userInfos} = storeToRefs(userStore)
const {copyText} = commonFunction();

const articleStore =useArticleStore()
const articleStoreData = storeToRefs(articleStore)
const qrcodeRef = ref();

//定义参数
const state = reactive({
  host: {
    host: '',
  },
  isShowDialog:false,
  isShowSubDialog:false,
  isShowQRSubDialog:false,
  QRcode: null,
})

//获取首页自定义内容
const getArticleID1=()=>{
  articleStore.getArticleList({search:"id=1 AND status=true",page_num:1,page_size:1}).then((res)=>{
    //保存
    articleStoreData.articleID1.value=res.data.article_list[0]
  })
}
//获取首页弹窗内容
const getArticleID2=()=>{
  articleStore.getArticleList({search:"id=2 AND status=true",page_num:1,page_size:1}).then((res)=>{
    if ((res.data.article_list as Article[]).length === 0){
      return
    }
    //保存
    articleStoreData.articleID2.value=res.data.article_list[0]
    //显示弹窗
    if (articleStoreData.articleID2.value.content !== ''){
      state.isShowDialog = true
    }
  })
}

//修改混淆
const onChangeHost = () => {
  userStore.changeHost(state.host)
  state.host.host = ''
}
//重置订阅
const onResetSub = () => {
  request(apiStoreData.api.value.user_resetSub).then((res) => {
    ElMessage.success(res.msg)
    // window.location.reload()
    userStore.getUserInfo()
  })
}
//复制订阅
const Sub = (type: string) => {
  switch (type) {
    case "NekoBox":
      copyText(userStore.subUrl + "&type=NekoBox")
      break
    case "v2rayNG":
      copyText(userStore.subUrl + "&type=v2rayNG")
      break
    case "v2rayN":
      copyText(userStore.subUrl + "&type=v2rayN")
      break
    case "Clash":
      copyText(userStore.subUrl + "&type=Clash")
      break

    case "Shadowrocket":
      copyText(userStore.subUrl + "&type=Shadowrocket")
      break

    case "Surge":
      copyText(userStore.subUrl + "&type=Surge")
      break

    case "Quantumult":
      copyText(userStore.subUrl + "&type=Quantumult")
      break
    case "V2rayU":
      copyText(userStore.subUrl + "&type=V2rayU")
      break

    default:
      copyText(userStore.subUrl)
      break
  }
}
//二维码订阅
const QRSub=()=>{
  state.isShowQRSubDialog=true
  setTimeout(()=>{
    onInitQrcode()
  },200)


}
//二维码
const onInitQrcode = () => {
  //清除上一次二维码
  const codeHtml = document.getElementById("qrcode");
  codeHtml.innerHTML = "";
  state.QRcode = new QRCode(qrcodeRef.value, {
    text: userStore.subUrl,
    width: 300,
    height: 300,
    colorDark: '#0a55f8',
    colorLight: 'rgb(255,255,255)',
  });
}
// 页面加载时
onMounted(() => {
  getArticleID1();
  getArticleID2();
});

</script>

<style scoped lang="scss">

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.home-card-item {
  padding: 10px;
  overflow: hidden;
  background: var(--el-color-white);
  color: var(--el-text-color-primary);
}

.el-card {
  background-image: url("../../assets/bgc/bg-1.svg");
  background-repeat: no-repeat;
  background-position: 100%, 100%;
  //background: rgba(0,0,0,0.3);
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