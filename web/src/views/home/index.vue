<template>
  <div>
    <el-row :gutter="10" class="home-card-two mb15">
      <el-col :xs="24" :sm="12" :md="12" :lg="6" :xl="6">
        <div class="home-card-item">
          <el-card class="box-card">
            <template #header>
              <div class="card-header">
                <el-text class="card-header-left">套餐详情</el-text>
                <el-text size="large" type="primary">{{userInfos.subscribe_info.goods_subject}}</el-text>
              </div>
            </template>
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
      <el-col :xs="24" :sm="12" :md="12" :lg="6" :xl="6">
        <div class="home-card-item">
          <el-card class="box-card" style="width: 100%;flex: 1;">
            <template #header>
              <div class="card-header">
                <el-text class="card-header-left">当前混淆:</el-text>
                <span>{{ userInfos.subscribe_info.host }}</span>
              </div>
            </template>
            <div>
              <el-input v-model="state.host.host" placeholder="输入混淆">
                <template #append>
                  <el-button size="large" @click="onChangeHost" :icon="Select">确认修改</el-button>
                </template>
              </el-input>
            </div>

          </el-card>
        </div>
      </el-col>
      <el-col :xs="24" :sm="14" :md="14" :lg="16" :xl="16">
        <div class="home-card-item">
          <el-card class="box-card">
            <template #header>
              <div class="card-header">
                <el-text class="card-header-left">订阅地址</el-text>
                <el-button type="primary" size="large" text plain class="button" @click="onResetSub">重置订阅链接</el-button>
              </div>
            </template>
            <div>
              <el-button style="margin-top: 10px;margin-bottom: 10px" @click="v2rayNGSub('v2ray')" type="primary" plain>
                复制通用订阅订阅
              </el-button>
              <el-button style="margin-top: 10px;margin-bottom: 10px" @click="v2rayNGSub('clash')" type="success" plain>
                复制Clash Meta订阅
              </el-button>
            </div>
          </el-card>
        </div>
      </el-col>
      <el-col :xs="24" :sm="14" :md="14" :lg="16" :xl="16">
        <div class="home-card-item" >
          <div class="box-card">
            <v-md-preview :text="articleStoreData.articleID1.value.content"></v-md-preview>
          </div>
        </div>
      </el-col>
    </el-row>
    <el-dialog v-model="state.isShowDialog" :title="state.title" width="80%" destroy-on-close center>
      <v-md-preview :text="articleStoreData.articleID2.value.content"></v-md-preview>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">

import {DateStrtoTime} from "/@/utils/formatTime";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
import {storeToRefs} from "pinia";
import {useUserStore} from "/@/stores/userStore";
import {onMounted, reactive} from 'vue';
import {ElMessage} from 'element-plus';
import {Select} from '@element-plus/icons-vue'
import commonFunction from '/@/utils/commonFunction';
import {useArticleStore} from "/@/stores/articleStore";

const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
const userStore = useUserStore()
const {userInfos} = storeToRefs(userStore)
const {copyText} = commonFunction();

const articleStore =useArticleStore()
const articleStoreData = storeToRefs(articleStore)
//定义参数
const state = reactive({
  host: {
    host: '',
  },
  isShowDialog:false,
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
const v2rayNGSub = (type: string) => {
  switch (type) {
    case "v2ray":
      //通用订阅；v2rayNG订阅
      copyText(userStore.subUrl + "&type=v2ray")
      break
    case "clash":
      //Clash订阅
      copyText(userStore.subUrl + "&type=clash")
      break
    default:
      copyText(userStore.subUrl + "&type=v2ray")
      break
  }
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