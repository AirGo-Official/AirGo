<template>
  <div>
    <div v-html="articleStoreData.articleID1.value.content"></div>
  </div>
</template>

<script setup lang="ts">
import { useArticleStore } from "/@/stores/user_logic/articleStore";
import { storeToRefs } from "pinia";
import { onMounted } from "vue";
import { NextLoading } from "/@/utils/loading";

const articleStore = useArticleStore()
const articleStoreData = storeToRefs(articleStore)
onMounted(()=>{
  NextLoading.done();
  articleStore.getDefaultArticles()
});
</script>

<style scoped>

</style>