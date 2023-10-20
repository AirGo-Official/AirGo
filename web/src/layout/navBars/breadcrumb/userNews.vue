<template>
  <div class="layout-navbars-breadcrumb-user-news">
    <div class="head-box">
      <div class="head-box-title">通知</div>
      <!--      <div class="head-box-btn"  @click="onAllReadClick">全部已读</div>-->
    </div>
    <div class="content-box">
      <div class="content-box-item" @click="toDetails(v)" v-for="(v,k) in state.articleDate.article_list" :key="k">
        <div>{{ v.title }}</div>
        <div class="content-box-msg">
          {{ v.introduction }}
        </div>
        <div class="head-box">
          <span class="content-box-time">{{ DateStrtoTime(v.created_at) }}</span>
          <el-text type="primary">详情>>></el-text>
        </div>

      </div>
    </div>
    <!--    <div class="foot-box" @click="onGoToMore">更多>>></div>-->
    <el-dialog v-model="state.isShowDialog" width="768px" destroy-on-close align-center>
      <div style="font-size: 40px">
        {{ state.currentArticle.title }}
      </div>
      <div>
        {{ state.currentArticle.introduction }}
      </div>
      <el-divider content-position="left"></el-divider>
      <div>
        <v-md-preview :text="state.currentArticle.content"></v-md-preview>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts" name="layoutBreadcrumbUserNews">
import {onMounted, reactive} from 'vue';
import {DateStrtoTime} from "../../../utils/formatTime";
import {useRoute, useRouter} from 'vue-router';
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
import {storeToRefs} from "pinia";
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)


const router = useRouter();
// 定义变量内容
const state = reactive({
  isShowDialog: false,
  params: {
    search: '',
    page_num: 1,
    page_size: 3,
  },
  articleDate: {
    total: 0,
    article_list: [] as Article[],
  },
  currentArticle: {} as Article,
});
//获取article列表
const getArticleList = (params: object) => {
  request(apiStoreData.api.value.article_getArticle,params).then((res) => {
      state.articleDate = res.data
      // console.log(state.articleDate)
  })
}
//打开详情弹窗
const onOpenDialog = () => {
  state.isShowDialog = true

}
// 全部已读点击
const onAllReadClick = () => {

};
// 前往更多
const onGoToMore = () => {
  // window.open('');
};
// 前往详情
const toDetails = (params: Article) => {
  // router.push({path:'notices',query:{id:1}})
  console.log("k:", params)
  // JSON.parse(JSON.stringify(row))
  state.currentArticle = JSON.parse(JSON.stringify(params))
  state.isShowDialog = true
}
//
onMounted(() => {
  getArticleList(state.params)
});
</script>

<style scoped lang="scss">
.layout-navbars-breadcrumb-user-news {
  .head-box {
    display: flex;
    border-bottom: 1px solid var(--el-border-color-lighter);
    box-sizing: border-box;
    color: var(--el-text-color-primary);
    justify-content: space-between;
    height: 35px;
    align-items: center;

    .head-box-btn {
      color: var(--el-color-primary);
      font-size: 13px;
      cursor: pointer;
      opacity: 0.8;

      &:hover {
        opacity: 1;
      }
    }
  }

  .content-box {
    font-size: 13px;

    .content-box-item {
      border-top: 1px solid var(--el-border-color-lighter);
      padding-top: 12px;

      &:last-of-type {
        padding-bottom: 12px;
      }

      .content-box-msg {
        color: var(--el-text-color-secondary);
        margin-top: 5px;
        margin-bottom: 5px;
      }

      .content-box-time {
        color: var(--el-text-color-secondary);
      }
    }
  }

  .foot-box {
    height: 35px;
    color: var(--el-color-primary);
    font-size: 13px;
    cursor: pointer;
    opacity: 0.8;
    display: flex;
    align-items: center;
    justify-content: center;
    border-top: 3px solid var(--el-border-color-lighter);

    &:hover {
      opacity: 1;
    }
  }

  :deep(.el-empty__description p) {
    font-size: 13px;
  }
}
</style>
