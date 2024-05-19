<template>
  <div>
    <el-dialog v-model="state.isShowDialog" width="35em"
               :show-close="true"
               destroy-on-close style="border-radius: 1vh;" >
      <!--      <div v-html="state.dialogHTML"></div>-->
      <v-md-preview :text="state.dialogHTML"></v-md-preview>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive } from "vue";
import { Session } from "/@/utils/storage";
import { useArticleStore } from "/@/stores/user_logic/articleStore";
import { storeToRefs } from "pinia";

const articleStore = useArticleStore();
const articleStoreData = storeToRefs(articleStore);

const state = reactive({
  indexHTML: "",
  dialogHTML: "",
  isShowDialog: false
});
const openDialog = () => {
  if (articleStoreData.defaultArticles.value.length === 2) {
    state.dialogHTML = articleStoreData.defaultArticles.value[1].content;
    //判断是否显示过
    if (!Session.get("defaultDialogDisplayed")) {
      state.isShowDialog = true;
      Session.set("defaultDialogDisplayed", "1");
    }
  }
};

// 暴露变量
defineExpose({
  openDialog
});


</script>

<style scoped>

</style>
