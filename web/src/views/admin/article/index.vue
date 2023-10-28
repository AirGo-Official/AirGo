<template>
  <div class="container layout-padding">
    <el-card shadow="hover" class="layout-padding-auto">
      <div class="mb15">
        <el-input v-model="state.params.search" size="default" placeholder="请输入名称"
                  style="max-width: 180px"></el-input>
        <el-button size="default" type="primary" class="ml10" @click="getArticleList(state.params)">
          <el-icon>
            <ele-Search/>
          </el-icon>
          查询
        </el-button>
        <el-button size="default" type="success" class="ml10" @click="onOpenDialog('add')">
          <el-icon>
            <ele-FolderAdd/>
          </el-icon>
          新建文章
        </el-button>
      </div>
      <el-table :data="articleStoreData.articleDate.value.article_list" stripe height="100%" style="width: 100%;flex: 1;">
        <el-table-column fixed type="index" label="序号" width="60"/>
        <el-table-column prop="title" label="标题" show-overflow-tooltip width="200"></el-table-column>
        <el-table-column prop="id" label="ID" show-overflow-tooltip width="60"></el-table-column>
        <el-table-column prop="status" label="是否显示" show-overflow-tooltip width="80">
          <template #default="{ row }">
            <el-button v-if="row.status" type="primary">显示</el-button>
            <el-button v-else type="info">隐藏</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="类型" show-overflow-tooltip width="80">
          <template #default="{ row }">
            <el-button v-if="row.type==='notice'" type="primary">公告</el-button>
            <el-button v-else-if="row.type==='knowledge'" type="primary">知识库</el-button>
            <el-button v-else type="info">系统保留</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="introduction" label="简介" show-overflow-tooltip width="200"></el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="scope">
            <el-button size="small" text type="primary"
                       @click="onOpenDialog('edit', scope.row)">修改
            </el-button>
            <el-button :disabled="scope.row.id == 1 || scope.row.id == 2" size="small" text type="primary"
                       @click="deleteArticle(scope.row)">删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
          background
          :page-sizes="[10, 30, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="onHandleSizeChange" @current-change="onHandleCurrentChange"
          v-model:current-page="state.params.page_num"
          v-model:page-size="state.params.page_size"
          :total="articleStoreData.articleDate.value.total"
      />
    </el-card>
    <ArticleDialog ref="articleDialogRef" @refresh="getArticleList(state.params)"></ArticleDialog>
  </div>
</template>

<script lang="ts" setup>

import {ElMessage, ElMessageBox} from "element-plus";
import {defineAsyncComponent, onMounted, reactive, ref} from "vue";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
import {storeToRefs} from "pinia";
import {useArticleStore} from "/@/stores/articleStore";

const ArticleDialog = defineAsyncComponent(() => import('/@/views/admin/article/dialog.vue'))
const articleDialogRef = ref()

const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
const articleStore = useArticleStore()
const articleStoreData = storeToRefs(articleStore)

//定义变量
const state = reactive({
  params: {
    search: '',
    page_num: 1,
    page_size: 30,
  },
})

function deleteArticle(row: any) {
  ElMessageBox.confirm(`此操作将永久删除文章：${row.title}, 是否继续?`, '提示', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        request(apiStoreData.api.value.article_deleteArticle, row)
        setTimeout(() => {
          getArticleList(state.params)
          ElMessage.success('成功');
        }, 500);
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
  request(apiStoreData.api.value.article_getArticle, params).then((res) => {
    articleStoreData.articleDate.value = res.data
  })
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