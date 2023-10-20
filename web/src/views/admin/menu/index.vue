<template>
  <div class="system-menu-container layout-pd">
    <el-card shadow="hover">
      <div class="system-menu-search mb15">
        <el-input v-model="state.search.meta.title" size="default" placeholder="请输入菜单名称"
                  style="max-width: 180px">
        </el-input>
        <el-button @click="onSearch" size="default" type="primary" class="ml10">
          <el-icon>
            <ele-Search/>
          </el-icon>
          查询
        </el-button>
        <el-button size="default" type="success" class="ml10" @click="onOpenAddMenu('add')">
          <el-icon>
            <ele-FolderAdd/>
          </el-icon>
          新增菜单
        </el-button>
      </div>
      <el-table :data="allRoutesList" v-loading="state.tableData.loading" stripe style="width: 100%" row-key="path"
                :tree-props="{ children: 'children', hasChildren: 'hasChildren' }">
        <el-table-column label="菜单名称" show-overflow-tooltip width="200px">
          <template #default="scope">
            <SvgIcon :name="scope.row.meta.icon"/>
            <span class="ml10">{{ scope.row.meta.title }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="path" label="路由路径" show-overflow-tooltip width="150px"></el-table-column>
        <el-table-column label="组件路径" show-overflow-tooltip width="200px">
          <template #default="scope">
            <span>{{ scope.row.component }}</span>
          </template>
        </el-table-column>
        <el-table-column label="菜单ID" show-overflow-tooltip width="80">
          <template #default="scope">
            {{ scope.row.id }}
          </template>
        </el-table-column>
        <el-table-column label="操作" show-overflow-tooltip>
          <template #default="scope">
            <el-button size="small" text type="primary"
                       @click="onOpenEditMenu('edit', scope.row)">修改
            </el-button>
            <el-button size="small" text type="primary" @click="onTableRowDel(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <MenuDialog ref="menuDialogRef" @refresh="getTableData()"/>
  </div>
</template>

<script setup lang="ts" name="systemMenu">
import {defineAsyncComponent, onMounted, reactive, ref} from 'vue';
import {RouteRecordRaw} from 'vue-router';
import {ElMessage, ElMessageBox} from 'element-plus';
import {storeToRefs} from 'pinia';
import {useRoutesStore} from '/@/stores/routesStore';

const MenuDialog = defineAsyncComponent(() => import('/@/views/admin/menu/dialog.vue'));
const stores = useRoutesStore();
const {allRoutesList} = storeToRefs(stores);
const menuDialogRef = ref();
const state = reactive({
  tableData: {
    loading: true,
  },
  //查询参数
  search: {
    meta: {
      title: ""
    }
  } as Route
});

// 获取路由数据
const getTableData = () => {
  state.tableData.loading = true;
  stores.setAllRoutesList()
  setTimeout(() => {
    state.tableData.loading = false;
  }, 500);

};
//查询动态路由
const onSearch = () => {
  if (state.search.meta.title === '') {
    stores.setAllRoutesList()
  } else {
    stores.findRoutesListByTitle(state.search)
  }
}
// 打开新增菜单弹窗
const onOpenAddMenu = (type: string) => {
  menuDialogRef.value.openDialog(type);
};
// 打开编辑菜单弹窗
const onOpenEditMenu = (type: string, row: RouteRecordRaw) => {
  menuDialogRef.value.openDialog(type, row);
};
// 删除当前行
const onTableRowDel = (row: RouteRecordRaw) => {
  ElMessageBox.confirm(`此操作将永久删除路由：${row.path}, 是否继续?`, '提示', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        //console.log("删除当前行 row:",row)
        stores.delDynamicRoute(row)
        setTimeout(() => {
          getTableData();
          ElMessage.success('删除成功');
        }, 1000);
      })
      .catch(() => {
      });
};
// 页面加载时
onMounted(() => {
  getTableData();
});
</script>
