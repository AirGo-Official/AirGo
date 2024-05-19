<template>
<div>
  <el-card style="margin: 1vh;border-radius: 1vh;text-align: center;">
  <div style="font-size: 4vh">{{articleStoreData.currentArticle.value.title}}</div>
  <div style="font-size: 1.8vh">{{articleStoreData.currentArticle.value.introduction}}</div>
  <el-button round style="margin-top: 1.5vh;font-size: 1.9vh"@click='$router.back()'><i class="ri-arrow-left-line" ></i>{{$t('message.user.back')}}</el-button>

</el-card>
  <el-scrollbar height="80vh">
  <v-md-preview style="color:var(--el-text-color-primary);" :text="articleStoreData.currentArticle.value.content"></v-md-preview>
</el-scrollbar>
</div>
</template>

<script setup lang="ts">
import {useArticleStore} from "/@/stores/user_logic/articleStore";
import {storeToRefs} from "pinia";

const articleStore = useArticleStore()
const articleStoreData = storeToRefs(articleStore)


</script>

<style scoped>

</style>