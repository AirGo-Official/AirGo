<template>
  <el-dialog v-model="state.isShowDialog" :title="state.title" width="100%" destroy-on-close align-center>
    <el-form label-position="top">
      <el-form-item label="文章主题">
        <el-input v-model="articleStoreData.currentArticle.value.title"/>
      </el-form-item>
      <el-form-item label="文章简介">
        <el-input v-model="articleStoreData.currentArticle.value.introduction"/>
      </el-form-item>
      <el-form-item label="是否显示">
        <el-switch v-model="articleStoreData.currentArticle.value.status" inline-prompt active-text="显示"
                   inactive-text="隐藏"
                   style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
      </el-form-item>
      <el-form-item label="类型"
                    v-if="articleStoreData.currentArticle.value.id !== 1 && articleStoreData.currentArticle.value.id !== 2">
        <el-radio-group v-model="articleStoreData.currentArticle.value.type">
          <el-radio label="notice">公告</el-radio>
          <el-radio label="knowledge">知识库</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="文章内容---使用markdown语法">
        <v-md-editor v-model="articleStoreData.currentArticle.value.content" height="400px"></v-md-editor>
      </el-form-item>
    </el-form>
    <template #footer>
            <span class="dialog-footer">
                <el-button @click="closeDialog">取消</el-button>
                <el-button type="primary" @click="onSubmit">
                    确认
                </el-button>
            </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
// 定义子组件向父组件传值/事件
import {useArticleStore} from "/@/stores/articleStore";
import {reactive} from "vue";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
import {storeToRefs} from "pinia";

const emit = defineEmits(['refresh']);
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
const articleStore = useArticleStore()
const articleStoreData = storeToRefs(articleStore)

//定义参数
const state = reactive({
  type: "",
  title: "",
  isShowDialog: false,
})

// 打开弹窗
const openDialog = (type: string, row?: any) => {
  // console.log("打开弹窗:", type)
  if (type == 'add') {
    state.type = type
    state.title = "新建文章"
    state.isShowDialog = true
    articleStoreData.currentArticle.value.id = 0
  } else {
    articleStoreData.currentArticle.value = row
    state.type = type
    state.title = "修改文章"
    state.isShowDialog = true
  }
}
// 关闭弹窗
const closeDialog = () => {
  state.isShowDialog = false
};

//确认提交
function onSubmit() {
  if (state.type === 'add') {
    request(apiStoreData.api.value.article_newArticle, articleStoreData.currentArticle.value)
    setTimeout(() => {
      emit('refresh');
    }, 500);
  } else {
    request(apiStoreData.api.value.article_updateArticle, articleStoreData.currentArticle.value)
    setTimeout(() => {
      emit('refresh');
    }, 500);
  }
  closeDialog()
}

// 暴露变量
defineExpose({
  openDialog,   // 打开弹窗
});

</script>

<style scoped>

</style>