<template>
  <div class="container layout-padding">
    <el-card shadow="hover" class="layout-padding-auto" v-loading="state.isLoadingTable">
      <el-row>
        <el-col :xs="24" :sm="24" :md="24" :lg="12" :xl="12">
          <el-button size="default" type="primary" class="ml10" @click="openNodeDialog('add',{})">
            <el-icon>
              <ele-FolderAdd />
            </el-icon>
            {{ $t("message.common.add") }}
          </el-button>
          <el-button size="default" type="primary" class="ml10" @click="openNodeSortDialog">
            <el-icon>
              <DCaret />
            </el-icon>
            {{ $t("message.common.sort") }}
          </el-button>
          <el-button size="default" type="primary" class="ml10" @click="openAccessDrawer">
            <el-icon>
              <ChromeFilled />
            </el-icon>
            {{ $t("message.adminNode.access") }}
          </el-button>
          <el-button size="default" type="primary" class="ml10" @click="openServerStatusDialog">
              <el-icon>
                <Odometer />
              </el-icon>
            {{ $t("message.adminNode.serverStatus") }}
          </el-button>

        </el-col>
      </el-row>
      <el-table :data="nodeStoreData.nodeList.value.data" height="100%" stripe style="width: 100%;flex: 1;"
                @sort-change="sortChange">
        <!--        <el-table-column fixed type="index" :label="$t('message.adminNode.NodeInfo.index')" width="60" />-->
        <el-table-column prop="id" :label="$t('message.adminNode.NodeInfo.id')" show-overflow-tooltip width="60"
                         sortable="custom" fixed></el-table-column>

        <el-table-column prop="node_type" :label="$t('message.adminNode.NodeInfo.node_type')"
                         width="120"
                         show-overflow-tooltip
                         sortable="custom"
                         fixed>
          <template #default="{row}">
            <el-tag type="success" style="width: 100px" v-if="row.node_type ===constantStore.NODE_TYPE_NORMAL">
              {{ $t("message.constant.NODE_TYPE_NORMAL") }}
            </el-tag>
            <el-tag type="warning" style="width: 100px" v-else-if="row.node_type ===constantStore.NODE_TYPE_TRANSFER">
              {{ $t("message.constant.NODE_TYPE_TRANSFER") }}
            </el-tag>
            <el-tag type="info" style="width: 100px" v-else-if="row.node_type ===constantStore.NODE_TYPE_SHARED">
              {{ $t("message.constant.NODE_TYPE_SHARED") }}
            </el-tag>
          </template>
        </el-table-column>


        <el-table-column prop="protocol" :label="$t('message.adminNode.NodeInfo.protocol')"
                         width="120"
                         show-overflow-tooltip
                         sortable="custom"
        >
          <template #default="{row}">
            <el-button type="success" style="width: 100px" v-if="row.protocol ===constantStore.NODE_PROTOCOL_VMESS">
              {{ $t("message.constant.NODE_PROTOCOL_VMESS") }}
            </el-button>
            <el-button type="warning" style="width: 100px"
                       v-else-if="row.protocol ===constantStore.NODE_PROTOCOL_VLESS">{{ $t("message.constant.NODE_PROTOCOL_VLESS") }}
            </el-button>
            <el-button type="info" style="width: 100px" v-else-if="row.protocol ===constantStore.NODE_PROTOCOL_TROJAN">
              {{ $t("message.constant.NODE_PROTOCOL_TROJAN") }}
            </el-button>
            <el-button type="danger" style="width: 100px"
                       v-else-if="row.protocol ===constantStore.NODE_PROTOCOL_SHADOWSOCKS">{{ $t("message.constant.NODE_PROTOCOL_SHADOWSOCKS") }}
            </el-button>
            <el-button type="primary" style="width: 100px"
                       v-else-if="row.protocol ===constantStore.NODE_PROTOCOL_HYSTERIA">{{ $t("message.constant.NODE_PROTOCOL_HYSTERIA") }}
            </el-button>
          </template>
        </el-table-column>

        <el-table-column prop="remarks" :label="$t('message.adminNode.NodeInfo.remarks')" show-overflow-tooltip
                         width="200" sortable="custom"></el-table-column>
        <el-table-column prop="enabled" :label="$t('message.adminNode.NodeInfo.enabled')" width="120"
                         show-overflow-tooltip sortable="custom">
          <template #default="scope">
            <el-button type="success" v-if="scope.row.enabled">{{ $t("message.common.display") }}</el-button>
            <el-button type="info" v-else>{{ $t("message.common.hide") }}</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="address" :label="$t('message.adminNode.NodeInfo.address')" show-overflow-tooltip
                         width="150" sortable="custom">
          <template #default="{row}">
            <span v-if="row.node_type === constantStore.NODE_TYPE_TRANSFER">{{ row.transfer_address }}</span>
            <span v-else>{{ row.address }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="port" :label="$t('message.adminNode.NodeInfo.port')" width="100" show-overflow-tooltip
                         sortable="custom"></el-table-column>
        <el-table-column prop="node_speed_limit" :label="$t('message.adminNode.NodeInfo.node_speed_limit')"
                         show-overflow-tooltip width="100" sortable="custom"></el-table-column>
        <el-table-column prop="traffic_rate" :label="$t('message.adminNode.NodeInfo.traffic_rate')"
                         show-overflow-tooltip sortable="custom"></el-table-column>
        <el-table-column :label="$t('message.common.operate')">
          <template #default="scope">
            <el-button size="small" text type="primary"
                       @click="openNodeDialog('edit', scope.row)">{{ $t("message.common.modify") }}
            </el-button>
            <el-button size="small" text type="danger"
                       @click="onRowDel(scope.row)">{{ $t("message.common.delete") }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <NodeDialog ref="nodeDialogRef" @refresh="getNodeList()" />
    <NodeSortDialog ref="nodeSortDialogRef" @refresh="getNodeList()"></NodeSortDialog>
    <AccessDrawer ref="AccessDrawerRef"></AccessDrawer>
    <ServerStatusDialog ref="serverStatusDialogRef"></ServerStatusDialog>
  </div>
</template>

<script setup lang="ts">

import { defineAsyncComponent, onBeforeMount, onMounted, reactive, ref } from "vue";
import { storeToRefs } from "pinia";
import { useAdminNodeStore } from "/@/stores/admin_logic/nodeStore";
import { ElMessageBox } from "element-plus";
import { formatDate } from "/@/utils/formatTime";
import { useConstantStore } from "/@/stores/constantStore";
import { useI18n } from "vue-i18n";

const NodeDialog = defineAsyncComponent(() => import("/@/views/admin/node/dialog_edit.vue"));
const NodeSortDialog = defineAsyncComponent(() => import("/@/views/admin/node/dialog_node_sort.vue"));
const AccessDrawer = defineAsyncComponent(() => import("/@/views/admin/node/drawer_access.vue"));
const ServerStatusDialog = defineAsyncComponent(() => import("/@/views/admin/node/dialog_server_status.vue"));
const nodeDialogRef = ref();
const nodeSortDialogRef = ref();
const serverStatusDialogRef = ref()
const AccessDrawerRef = ref();
const nodeStore = useAdminNodeStore();
const nodeStoreData = storeToRefs(nodeStore);
const constantStore = useConstantStore();
const { t } = useI18n();

//定义参数
const state = reactive({
  isLoadingTable: false,
  this_month: [""],
  queryParams: {
    table_name: "node",
    field_params_list: [
      { field: "id", field_chinese_name: "", field_type: "", condition: "<>", condition_value: "", operator: "" }
      // {field: 'created_at', field_chinese_name: '', field_type: '', condition: '<', condition_value: "", operator: 'AND',}
    ] as FieldParams[],
    pagination: { page_num: 1, page_size: 9999, order_by: "node_order" } as Pagination//设为9999，理论能获取全部节点，暂时取消详细的分页设置
  } as QueryParams,
  nodeTypeFilters: [
    { text: constantStore.NODE_PROTOCOL_VLESS, value: constantStore.NODE_PROTOCOL_VLESS },
    { text: constantStore.NODE_PROTOCOL_VMESS, value: constantStore.NODE_PROTOCOL_VMESS },
    { text: constantStore.NODE_PROTOCOL_SHADOWSOCKS, value: constantStore.NODE_PROTOCOL_SHADOWSOCKS },
    { text: constantStore.NODE_PROTOCOL_HYSTERIA, value: constantStore.NODE_PROTOCOL_HYSTERIA },
    { text: constantStore.NODE_TYPE_TRANSFER, value: constantStore.NODE_TYPE_TRANSFER }
  ]
});

//打开新建节点，修改节点弹窗
function openNodeDialog(title: string, row?: NodeInfo) {
  nodeDialogRef.value.openDialog(title, row);
}

//打开节点排序弹窗
function openNodeSortDialog() {
  nodeSortDialogRef.value.openDialog();
}

//打开访问控制抽屉
const openAccessDrawer = () => {
  AccessDrawerRef.value.openDrawer();
};
const openServerStatusDialog=()=>{
  serverStatusDialogRef.value.openDialog();
}

//查询节点
const getNodeList = () => {
  // state.queryParams.field_params_list[0].condition_value = state.this_month[0]
  // state.queryParams.field_params_list[1].condition_value = state.this_month[1]
  nodeStore.getNodeList(state.queryParams);
};

//删除节点
function onRowDel(row: NodeInfo) {
  ElMessageBox.confirm(t("message.common.message_confirm_delete"), t("message.common.tip"), {
    confirmButtonText: t("message.common.button_confirm"),
    cancelButtonText: t("message.common.button_cancel"),
    type: "warning"
  })
    .then(() => {
      nodeStore.deleteNode(row).then(() => {
        getNodeList();
      });
    });
}

// 分页改变
const onHandleSizeChange = (val: number) => {
  state.queryParams.pagination.page_size = val;
  getNodeList();
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
  state.queryParams.pagination.page_num = val;
  getNodeList();
};
//排序监听
const sortChange = (column: any) => {
  //处理嵌套字段
  let p = (column.prop as string);
  if (p.indexOf(".") !== -1) {
    p = p.slice(p.indexOf(".") + 1);
  }
  switch (column.order) {
    case "ascending":
      state.queryParams.pagination.order_by = p + " ASC";
      break;
    default:
      state.queryParams.pagination.order_by = p + " DESC";
      break;
  }
  getNodeList();

};

onMounted(() => {
  getNodeList();
  nodeStore.getAccessList({
    table_name: "access",
    field_params_list: [
      { field: "id", condition: "<>", condition_value: "" }
    ] as FieldTableNew[],
    pagination: { page_num: 1, page_size: 9999 } as Pagination
  } as QueryParams);
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