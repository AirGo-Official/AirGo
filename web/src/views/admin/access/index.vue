<template>
  <div class="container layout-padding">
    <el-card shadow="hover" class="layout-padding-auto">
      <div class="mb15">
      <el-row :gutter="10" style="width: 768px">
        <el-col :span="4">
          <el-input v-model="search.condition_value" size="default" placeholder="请输入名称"
                    style="max-width: 180px"></el-input>
        </el-col>
        <el-col :span="3">
          <el-button  size="default" type="primary" class="ml10" @click="getAccess()">
            <el-icon>
              <ele-Search/>
            </el-icon>
            查询
          </el-button>
        </el-col>
        <el-col :span="4">
          <el-button size="default" type="success" class="ml10" @click="openDialog('新建')">
            <el-icon>
              <ele-FolderAdd/>
            </el-icon>
            新增
          </el-button>
        </el-col>
      </el-row>
      </div>
      <el-table :data="accessStoreData.routes_list.value.data"  stripe height="100%" style="width: 100%">
        <el-table-column fixed type="index" label="序号" width="60"/>
        <el-table-column prop="name" label="名称" show-overflow-tooltip width="200"></el-table-column>
        <el-table-column prop="route" label="禁止路由" show-overflow-tooltip width="800"></el-table-column>
        <el-table-column label="操作" >
          <template #default="scope">
            <el-button  size="small" text type="primary"
                       @click="openDialog('编辑', scope.row)">编辑
            </el-button>
            <el-button  size="small" text type="primary"
                       @click="onRowDel(scope.row)">删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
          background
          :page-sizes="[10, 30, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="onHandleSizeChange" @current-change="onHandleCurrentChange"
          v-model:current-page="accessStoreData.params.value.pagination.page_num"
          v-model:page-size="accessStoreData.params.value.pagination.page_size"
          :total="accessStoreData.routes_list.value.total"
      />
    </el-card>
    <Dialog ref="DialogRef"></Dialog>
  </div>
</template>

<script lang="ts" setup>

//定义参数
import {defineAsyncComponent, onMounted, reactive, ref} from "vue";
import {ElMessageBox} from "element-plus";
import {useAccessStore} from "/@/stores/accessStore";
import {storeToRefs} from "pinia";

const Dialog = defineAsyncComponent(() => import('/@/views/admin/access/dialog.vue'))
const DialogRef = ref()

const accessStore = useAccessStore()
const accessStoreData = storeToRefs(accessStore)
const search:FieldTableNew = accessStoreData.params.value.field_params_list[0]

const state = reactive({
  loading: true,
})

//打开弹窗
function openDialog(title: string, row?:any) {
  DialogRef.value.openDialog(title, row)
}
//获取数据
const getAccess=()=>{
  accessStore.getRoutesList(accessStoreData.params.value)
}
//删除
function onRowDel(row: any) {
  ElMessageBox.confirm(`此操作将永久删除：${row.name}，是否继续?`, '提示', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        accessStore.deleteRoutes(row)
        setTimeout(()=>{
          getAccess()
        },500)

      })
      .catch(() => {
      });
}

// 分页改变
const onHandleSizeChange = (val: number) => {
  accessStoreData.params.value.pagination.page_size = val;
  getAccess()

};
// 分页改变
const onHandleCurrentChange = (val: number) => {
  accessStoreData.params.value.pagination.page_num = val;
  getAccess()
};
onMounted(() => {
  getAccess()
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