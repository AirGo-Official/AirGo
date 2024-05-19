<template>
  <div class="container layout-padding">
    <el-drawer v-model="state.isShowDrawer" size="90%" :title="$t('message.adminNode.access')" destroy-on-close>
      <el-card shadow="hover" class="layout-padding-auto">
        <div class="mb15">
          <el-row :gutter="10" style="width: 768px">
            <el-col :span="4">
              <el-button size="default" type="success" class="ml10" @click="openDialog('add')">
                <el-icon>
                  <ele-FolderAdd />
                </el-icon>
                {{ $t("message.adminNode.addAccess") }}
              </el-button>
            </el-col>
          </el-row>
        </div>
        <el-table :data="nodeStoreData.accessList.value.data" stripe height="100%" style="width: 100%">
          <el-table-column fixed type="index" :label="$t('message.adminNode.Access.index')" width="60" />
          <el-table-column prop="name" :label="$t('message.adminNode.Access.name')" show-overflow-tooltip
                           width="200"></el-table-column>
          <el-table-column prop="route" :label="$t('message.adminNode.Access.route')" show-overflow-tooltip
                           width="800"></el-table-column>
          <el-table-column :label="$t('message.common.operate')">
            <template #default="scope">
              <el-button size="small" text type="primary"
                         @click="openDialog('edit', scope.row)">{{ $t("message.common.modify") }}
              </el-button>
              <el-button size="small" text type="danger"
                         @click="deleteAccess(scope.row)">{{ $t("message.common.delete") }}
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
          :total="nodeStoreData.accessList.value.total"
        />
      </el-card>
      <el-dialog v-model="state.isShowDialog" :title="state.dialogTitle" width="80%" destroy-on-close align-center>
        <el-form label-position="top">
          <el-form-item :label="$t('message.adminNode.Access.name')">
            <el-input v-model="nodeStoreData.currentAccess.value.name" />
          </el-form-item>
          <el-form-item :label="$t('message.adminNode.Access.route')">
            <el-input v-model="nodeStoreData.currentAccess.value.route" type="textarea" autosize />
          </el-form-item>
        </el-form>
        <template #footer>
            <span class="dialog-footer">
                <el-button @click="closeDialog">{{ $t("message.common.button_cancel") }}</el-button>
                <el-button type="primary" @click="onSubmit">
                    {{ $t("message.common.button_confirm") }}
                </el-button>
            </span>
        </template>
      </el-dialog>
    </el-drawer>
  </div>
</template>

<script lang="ts" setup>

//定义参数
import { defineAsyncComponent, onMounted, reactive, ref } from "vue";
import { ElMessageBox } from "element-plus";
import { storeToRefs } from "pinia";
import { useAdminNodeStore } from "/@/stores/admin_logic/nodeStore";
import { useI18n } from "vue-i18n";

const nodeStore = useAdminNodeStore();
const nodeStoreData = storeToRefs(nodeStore);
const { t } = useI18n();
const state = reactive({
  loading: true,
  queryParams: {
    table_name: "access",
    field_params_list: [
      { field: "name", condition: "like", condition_value: "" }
    ] as FieldTableNew[],
    pagination: { page_num: 1, page_size: 30 } as Pagination//分页参数
  } as QueryParams,
  dialogType: "",
  dialogTitle: "",
  isShowDialog: false,
  isShowDrawer: false
});
const openDrawer = () => {
  state.isShowDrawer = true;
  getAccess();
};

//获取数据
const getAccess = () => {
  nodeStore.getAccessList(state.queryParams);
};
//删除
const deleteAccess = (row: any) => {
  ElMessageBox.confirm(t("message.common.message_confirm_delete"), t("message.common.tip"), {
    confirmButtonText: t("message.common.button_confirm"),
    cancelButtonText: t("message.common.button_cancel"),
    type: "warning"
  })
    .then(() => {
      nodeStore.deleteAccess(row);
      setTimeout(() => {
        getAccess();
      }, 500);

    })
    .catch(() => {
    });
};

// 分页改变
const onHandleSizeChange = (val: number) => {
  state.queryParams.pagination.page_size = val;
  getAccess();

};
// 分页改变
const onHandleCurrentChange = (val: number) => {
  state.queryParams.pagination.page_num = val;
  getAccess();
};

// 打开弹窗
const openDialog = (type: string, row?: any) => {
  state.isShowDialog = true;
  state.dialogType = type;
  switch (type) {
    case "add":
      state.dialogTitle = t("message.adminNode.addAccess");
      break;
    case "edit":
      state.dialogTitle = t("message.adminNode.modifyAccess");
      nodeStoreData.currentAccess.value = row;
      break;
    default:
      break;
  }
};
// 关闭弹窗
const closeDialog = () => {
  state.isShowDialog = false;
};

//确认提交
function onSubmit() {
  switch (state.dialogType) {
    case "add":
      nodeStore.newAccess(nodeStoreData.currentAccess.value).then(() => {
        getAccess();
      });
      break;
    case "edit":
      nodeStore.updateAccess(nodeStoreData.currentAccess.value).then(() => {
        getAccess();
      });
      break;
    default:
      break;
  }
  closeDialog();
}

// 暴露变量
defineExpose({
  openDrawer
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