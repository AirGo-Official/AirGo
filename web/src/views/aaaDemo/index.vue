<template>
  <div class="container layout-padding">
    <el-card shadow="hover" class="layout-padding-auto">
      <el-row :gutter="10" style="width: 768px">
        <el-col :span="4">
          <el-input v-model="state.params.search" size="default" placeholder="请输入名称"
                    style="max-width: 180px"></el-input>
        </el-col>
        <el-col :span="3">
          <el-button  size="default" type="primary" class="ml10">
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

      <el-table :data="state.data" height="100%" stripe style="width: 100%;flex: 1;">
        <el-table-column fixed type="index" label="序号" width="60"/>
        <el-table-column prop="enabled" label="状态" show-overflow-tooltip>
          <template #default="scope">
            <el-tag type="success" v-if="scope.row.id===1">启用</el-tag>
            <el-tag type="danger" v-else>禁用</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="scope">
            <el-button  size="small" text type="primary"
                       @click="openDialog('edit', scope.row)">修改
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
          v-model:current-page="state.params.pagination_params.page_num"
          v-model:page-size="state.params.pagination_params.page_size"
          :total="state.total"
      />
    </el-card>
    <Dialog ref="DialogRef"></Dialog>
  </div>
</template>

<script lang="ts" setup>

//定义参数
import {defineAsyncComponent, onMounted, reactive, ref} from "vue";
import {ElMessageBox} from "element-plus";

const Dialog = defineAsyncComponent(() => import('/@/views/aaaDemo/dialog.vue'))
const DialogRef = ref()

const state = reactive({
  loading: true,
  data:[
    {id:1},
    {id:2},
  ],
  total:30,
  //查询的条件参数
  params: {
    table_name: '',
    field_params_list: [] as FieldParams[],    //搜索条件列表 {field: '', field_chinese_name: '', field_type: '', condition: '=', condition_value: '',}
    pagination_params: {
      page_num: 1,
      page_size: 30,
    } as Pagination,//分页参数
  },

})

//打开弹窗
function openDialog(title: string, row?: any) {
  DialogRef.value.openDialog(title, row)
}



//删除节点
function onRowDel(row: any) {
  ElMessageBox.confirm(`此操作将永久删除：${row.id}，是否继续?`, '提示', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {

      })
      .catch(() => {
      });
}

// 分页改变
const onHandleSizeChange = (val: number) => {
  state.params.pagination_params.page_size = val;

};
// 分页改变
const onHandleCurrentChange = (val: number) => {
  state.params.pagination_params.page_num = val;

};

onMounted(() => {

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