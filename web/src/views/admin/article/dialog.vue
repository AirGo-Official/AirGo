<template>
  <el-dialog v-model="state.isShowDialog" :title="state.title" width="100%" destroy-on-close align-center>
    <el-form label-position="top">
      <el-form-item :label="$t('message.adminArticle.Article.title')">
        <el-input v-model="articleStoreData.currentArticle.value.title"/>
      </el-form-item>
      <el-form-item :label="$t('message.adminArticle.Article.introduction')">
        <el-input v-model="articleStoreData.currentArticle.value.introduction"/>
      </el-form-item>
      <el-form-item :label="$t('message.adminArticle.Article.status')" >
        <el-switch v-model="articleStoreData.currentArticle.value.status"
                   inline-prompt
                   :active-text="$t('message.common.enable')"
                   :inactive-text="$t('message.common.disable')"
                   style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
      </el-form-item>
      <el-form-item :label="$t('message.adminArticle.Article.type')">
        <el-radio-group v-model="articleStoreData.currentArticle.value.type">
          <el-radio label="home">{{ $t('message.constant.ARTICLE_TYPE_HOME') }}</el-radio>
          <el-radio label="dialog">{{ $t('message.constant.ARTICLE_TYPE_DIALOG') }}</el-radio>
          <el-radio label="notice">{{ $t('message.constant.ARTICLE_TYPE_NOTICE') }}</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item :label="$t('message.adminArticle.Article.type')">
        <el-input type="textarea" autosize
                  v-if="articleStoreData.currentArticle.value.type === constantStore.ARTICLE_TYPE_HOME"
                  v-model="articleStoreData.currentArticle.value.content"></el-input>
        <v-md-editor v-else v-model="articleStoreData.currentArticle.value.content" height="400px"></v-md-editor>
      </el-form-item>
    </el-form>
    <template #footer>
            <span class="dialog-footer">
                <el-button @click="closeDialog">{{$t('message.common.button_cancel')}}</el-button>
                <el-button type="primary" @click="onSubmit">
                    {{$t('message.common.button_confirm')}}
                </el-button>
            </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">

import {reactive,defineAsyncComponent} from "vue";
import {useApiStore} from "/@/stores/apiStore";
import {storeToRefs} from "pinia";
import { useAdminArticleStore } from "/@/stores/admin_logic/articleStore";
import { useConstantStore } from "/@/stores/constantStore";
import { useI18n } from "vue-i18n";

const emit = defineEmits(['refresh']);
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
const articleStore = useAdminArticleStore()
const articleStoreData = storeToRefs(articleStore)
const constantStore = useConstantStore()
const {t} = useI18n()
// 引入组件
const Editor = defineAsyncComponent(() => import('/@/components/editor/index.vue'));

//定义参数
const state = reactive({
  type: "",
  title: "",
  isShowDialog: false,
  text:"",
})

// 打开弹窗
const openDialog = (type: string, row?: any) => {
  if (type == 'add') {
    state.type = type
    state.title = t('message.adminArticle.addArticle')
    state.isShowDialog = true
    articleStoreData.currentArticle.value.id = 0
  } else {
    articleStoreData.currentArticle.value = row
    state.type = type
    state.title = t('message.adminArticle.modifyArticle')
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
    articleStore.newArticle().then(()=>{
      emit('refresh');
    })
  } else {
    articleStore.updateArticle().then(()=>{
      emit('refresh');
    })
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