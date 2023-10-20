<template>
  <div>
    <div class="home-card-item" v-for="(v,k) in articleStoreData.articleDate.value.article_list" :key="k">
      <el-card class="box-card">
        <template #header>
          <div class="card-header">
            <span>{{v.title}}</span>
            <span @click="showArticle(v)" style="color: #ff0000">详情>>></span>
          </div>
        </template>
        <div>
          <div class="head-box">
            <div>{{v.introduction}}</div>
          </div>
        </div>
      </el-card>
    </div>
    <el-pagination background
                   class="mt15"
                   layout="total, sizes, prev, pager, next, jumper"
                   :page-sizes="[10, 30, 50]"
                   v-model:current-page="state.params.page_num"
                   v-model:page-size="state.params.page_size"
                   :total="articleStoreData.articleDate.value.total"
                   @size-change="onHandleSizeChange"
                   @current-change="onHandleCurrentChange">
    </el-pagination>
  </div>
</template>

<script setup lang="ts">
import {request} from "/@/utils/request";
import {onMounted, reactive} from "vue";
import {useApiStore} from "/@/stores/apiStore";
import {storeToRefs} from "pinia";
import {useArticleStore} from "/@/stores/articleStore";
import {DateStrtoTime} from "../../utils/formatTime";
import {useRoute, useRouter} from "vue-router";
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
const articleStore = useArticleStore()
const articleStoreData = storeToRefs(articleStore)
// const route = useRoute();
const router = useRouter()



//定义变量
const state = reactive({
  params: {
    search: `type='knowledge' AND status=true`,
    page_num: 1,
    page_size: 30,
  },
})

//获取article列表
const getArticleList = (params: object) => {
  request(apiStoreData.api.value.article_getArticle, params).then((res) => {
    articleStoreData.articleDate.value = res.data
  })
}

//
const showArticle=(v:Article)=>{
  // console.log(v)
  articleStoreData.currentArticle.value=v
  router.push({path:"/static/showArticle"})
}

// 分页改变
const onHandleSizeChange = (val: number) => {
  state.params.page_size = val;
  getArticleList(state.params)
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
  state.params.page_num = val;
  getArticleList(state.params)
};

//加载时
onMounted(() => {
  getArticleList(state.params)
});


</script>

<style scoped>

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
  /*background-image: url("../../assets/bgc/bg-1.svg");*/
  background-repeat: no-repeat;
  background-position: 100%, 100%;
  /*background: rgba(0,0,0,0.3);*/
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