<template>
  <div class="system-menu-container layout-pd">
    <el-card shadow="hover">
      <div class="system-menu-search mb15">
        <el-button size="default" type="success" class="ml10" @click="onOpenAddMenu('add')">
          <el-icon>
            <ele-FolderAdd/>
          </el-icon>
          {{$t('message.adminMenu.add_menu')}}
        </el-button>
      </div>
      <el-table :data="menuStoreData.allMenuList.value" v-loading="state.tableData.loading" stripe style="width: 100%" row-key="path"
                :tree-props="{ children: 'children', hasChildren: 'hasChildren' }">
        <el-table-column :label="$t('message.adminMenu.Route.title')" show-overflow-tooltip width="200px">
          <template #default="scope">
            <SvgIcon :name="scope.row.meta.icon"/>
            <span class="ml10">{{ $t(scope.row.meta.title) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="path" :label="$t('message.adminMenu.Route.path')" show-overflow-tooltip width="150px"></el-table-column>
        <el-table-column :label="$t('message.adminMenu.Route.component')" show-overflow-tooltip width="200px">
          <template #default="scope">
            <span>{{ scope.row.component }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('message.adminMenu.Route.id')" show-overflow-tooltip width="80">
          <template #default="scope">
            {{ scope.row.id }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('message.common.operate')" show-overflow-tooltip>
          <template #default="scope">
            <el-button size="small" text type="primary"
                       @click="onOpenEditMenu('edit', scope.row)">{{$t('message.common.modify')}}
            </el-button>
            <el-button size="small" text type="danger" @click="onTableRowDel(scope.row)">{{$t('message.common.delete')}}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <MenuDialog ref="menuDialogRef" @refresh="getMenuList()"/>
  </div>
</template>

<script setup lang="ts" name="systemMenu">
import {defineAsyncComponent, onMounted, reactive, ref} from 'vue';
import {RouteRecordRaw} from 'vue-router';
import {ElMessageBox} from 'element-plus';
import {storeToRefs} from 'pinia';
import { useAdminMenuStore } from "/@/stores/admin_logic/menuStore";
import { useI18n } from "vue-i18n";


const MenuDialog = defineAsyncComponent(() => import('/@/views/admin/menu/dialog.vue'));
const menuStore = useAdminMenuStore()
const menuStoreData = storeToRefs(menuStore);
const menuDialogRef = ref();
const {t} = useI18n()
const state = reactive({
  tableData: {
    loading: true,
  },
});

// 获取菜单列表
const getMenuList = () => {
  state.tableData.loading = true;
  menuStore.getAllMenuList()
  setTimeout(() => {
    state.tableData.loading = false;
  }, 500);

};
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
  ElMessageBox.confirm(t('message.common.message_confirm_delete'), t('message.common.tip'), {
    confirmButtonText: t('message.common.button_confirm'),
    cancelButtonText: t('message.common.button_cancel'),
    type: 'warning',
  })
      .then(() => {
        menuStore.delMenu(row).then(()=>{
          getMenuList();
        })
        })
};
// 页面加载时
onMounted(() => {
  getMenuList();
});
</script>
