<template>
<div>
  <div style="font-size: 50px;text-align: center">{{articleStoreData.currentArticle.value.title}}</div>
  <div style="color: #9b9da1;text-align: center">{{articleStoreData.currentArticle.value.introduction}}</div>
  <v-md-preview :text="articleStoreData.currentArticle.value.content"></v-md-preview>
</div>
</template>

<script setup lang="ts">
import {useArticleStore} from "/@/stores/articleStore";
import {storeToRefs} from "pinia";

const articleStore = useArticleStore()
const articleStoreData = storeToRefs(articleStore)

</script>

<style scoped>

</style>