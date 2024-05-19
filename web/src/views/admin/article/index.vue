<template>
  <div class="container layout-padding">
    <el-card shadow="hover" class="layout-padding-auto">
      <div class="mb15">
        <el-input v-model="state.queryParams.field_params_list[0].condition_value" size="default"
                  style="max-width: 180px"></el-input>
        <el-button size="default" type="primary" class="ml10" @click="getArticleList(state.queryParams)">
          <el-icon>
            <ele-Search/>
          </el-icon>
          {{$t('message.common.query')}}
        </el-button>
        <el-button size="default" type="success" class="ml10" @click="onOpenDialog('add')">
          <el-icon>
            <ele-FolderAdd/>
          </el-icon>
          {{$t('message.adminArticle.addArticle')}}
        </el-button>
      </div>
      <el-table :data="articleStoreData.articleList.value.data" stripe height="100%">
        <el-table-column fixed type="index" :label="$t('message.adminArticle.Article.index')" width="60"/>
        <el-table-column prop="title" :label="$t('message.adminArticle.Article.title')" show-overflow-tooltip width="200"></el-table-column>
        <el-table-column prop="id" :label="$t('message.adminArticle.Article.id')" show-overflow-tooltip width="60"></el-table-column>
        <el-table-column prop="status" :label="$t('message.adminArticle.Article.status')" show-overflow-tooltip width="80">
          <template #default="{ row }">
            <el-button v-if="row.status" type="primary">{{$t('message.common.enable')}}</el-button>
            <el-button v-else type="info">{{$t('message.common.disable')}}</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="type" :label="$t('message.adminArticle.Article.type')" show-overflow-tooltip width="100">
          <template #default="{ row }">
            <el-button v-if="row.type==='home'" type="primary">{{$t('message.constant.ARTICLE_TYPE_HOME')}}</el-button>
            <el-button v-else-if="row.type==='dialog'" type="primary">{{$t('message.constant.ARTICLE_TYPE_DIALOG')}}</el-button>
            <el-button v-else type="primary">{{$t('message.constant.ARTICLE_TYPE_NOTICE')}}</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="introduction" :label="$t('message.adminArticle.Article.introduction')" show-overflow-tooltip></el-table-column>
        <el-table-column :label="$t('message.common.operate')" width="100">
          <template #default="scope">
            <el-button size="small" text type="primary"
                       @click="onOpenDialog('edit', scope.row)">{{$t('message.common.modify')}}
            </el-button>
            <el-button :disabled="scope.row.id === 1 || scope.row.id === 2" size="small" text type="danger"
                       @click="deleteArticle(scope.row)">{{$t('message.common.delete')}}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
          background
          :page-sizes="[10, 30, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="onHandleSizeChange" @current-change="onHandleCurrentChange"
          v-model:current-page="state.queryParams.pagination.page_num"
          v-model:page-size="state.queryParams.pagination.page_size"
          :total="articleStoreData.articleList.value.total"
      />
    </el-card>
    <ArticleDialog ref="articleDialogRef" @refresh="getArticleList(state.queryParams)"></ArticleDialog>
  </div>
</template>

<script lang="ts" setup>

import {ElMessage, ElMessageBox} from "element-plus";
import {defineAsyncComponent, onMounted, reactive, ref} from "vue";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
import {storeToRefs} from "pinia";
import { useAdminArticleStore } from "/@/stores/admin_logic/articleStore";
import { useI18n } from "vue-i18n";


const ArticleDialog = defineAsyncComponent(() => import('/@/views/admin/article/dialog.vue'))
const articleDialogRef = ref()
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
const articleStore = useAdminArticleStore()
const articleStoreData = storeToRefs(articleStore)
const {t} = useI18n()

//定义变量
const state = reactive({
  queryParams:{
    table_name: 'article',
    field_params_list: [
      {field: 'title', field_chinese_name: '', field_type: '', condition: 'like', condition_value: "", operator: ''},
    ] as FieldParams[],
    pagination: { page_num: 1, page_size: 30, order_by: '', } as Pagination,//分页参数
  },
})

function deleteArticle(row: any) {
  ElMessageBox.confirm(t('message.common.message_confirm_delete'), t('message.common.tip'), {
    confirmButtonText: t('message.common.button_confirm'),
    cancelButtonText: t('message.common.button_cancel'),
    type: 'warning',
  })
      .then(() => {
        request(apiStoreData.adminApi.value.deleteArticle, row).then(()=>{
          getArticleList(state.queryParams)
        })
      })
      .catch(() => {
      });
}

//打开弹窗
const onOpenDialog = (type: string, row?: any) => {
  articleDialogRef.value.openDialog(type, row)
};
//获取article列表
const getArticleList = (params: object) => {
 articleStore.getArticleList(state.queryParams)
}

// 分页改变
const onHandleSizeChange = (val: number) => {
  state.queryParams.pagination.page_size = val;
  getArticleList(state.queryParams)
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
  state.queryParams.pagination.page_num = val;
  getArticleList(state.queryParams)
};

//加载时
onMounted(() => {
  getArticleList(state.queryParams)
});


</script>


<style scoped lang="scss">
.container {
  :deep(.el-card__body) {
    display: flex;
    flex-direction: column;
    flex: 1;
    overflow: auto;
    .el-table {
      flex: 1;
    }
  }
}
</style>