<template>
  <div class="text">
    <div v-html="state.indexHTML"></div>
    <DefaultDialog ref="DefaultDialogRef"></DefaultDialog>
  </div>
</template>

<script setup lang="ts">
import { useArticleStore } from "/@/stores/user_logic/articleStore";
import { storeToRefs } from "pinia";
import { defineAsyncComponent, nextTick, onMounted, reactive, ref } from "vue";
import { NextLoading } from "/@/utils/loading";
const DefaultDialog = defineAsyncComponent( () => import('/@/views/default/defaultDialog.vue'));
const DefaultDialogRef = ref()


const state = reactive({
  indexHTML:'',
  dialogHTML:'',
  isShowDialog:false,
})

const articleStore = useArticleStore()
const articleStoreData = storeToRefs(articleStore)

onMounted(()=>{
  NextLoading.done();
  articleStore.getDefaultArticles().then(()=>{
    if (articleStoreData.defaultArticles.value.length >=1){
      state.indexHTML = articleStoreData.defaultArticles.value[0].content
    }
    setTimeout(()=>{
      DefaultDialogRef.value.openDialog()
    },1000)
  })
});

</script>

<style scoped>
.text{
  width: 100%;
  overflow-y: scroll;
  height: 100%;
}



</style>