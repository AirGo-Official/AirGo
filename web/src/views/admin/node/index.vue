<template>
  <div class="container layout-padding">
    <el-card shadow="hover" class="layout-padding-auto">
      <el-row :gutter="10" style="width: 768px">
        <el-col :span="4">
          <el-input v-model="state.params.search" size="default" placeholder="请输入名称"
                    style="max-width: 180px"></el-input>
        </el-col>
        <el-col :span="8">
          <el-date-picker
              style="width: 250px"
              size="default"
              v-model="state.params.date"
              type="datetimerange"
              :shortcuts="shortcuts"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-col>
        <el-col :span="3">
          <el-button @click="onGetNode(state.params)" size="default" type="primary" class="ml10">
            <el-icon>
              <ele-Search/>
            </el-icon>
            查询
          </el-button>
        </el-col>
        <el-col :span="4">
          <el-button size="default" type="success" class="ml10" @click="onOpenEditNode('新建节点','vless')">
            <el-icon>
              <ele-FolderAdd/>
            </el-icon>
            新增节点
          </el-button>
        </el-col>
        <el-col :span="3">
          <el-button size="default" type="warning" class="ml10" @click="onOpenNodeSortDialog">
            <el-icon>
              <DCaret/>
            </el-icon>
            排序
          </el-button>
        </el-col>
        <el-col :span="2">
          <el-button size="default" type="primary" class="ml10" @click="onOpenNodeSharedDialog">
            <el-icon>
              <Share/>
            </el-icon>
            共享节点管理
          </el-button>
        </el-col>
      </el-row>

      <el-table :data="nodeManageData.nodes.node_list" height="100%" stripe style="width: 100%;flex: 1;">
        <el-table-column fixed type="index" label="序号" width="60"/>
        <el-table-column prop="remarks" label="节点名称" show-overflow-tooltip width="200"></el-table-column>
        <el-table-column prop="id" label="节点ID" show-overflow-tooltip width="60"></el-table-column>
        <el-table-column prop="address" label="节点地址" show-overflow-tooltip width="150"></el-table-column>
        <el-table-column prop="port" label="节点端口" width="80" show-overflow-tooltip></el-table-column>
        <el-table-column prop="sort" label="协议类型" width="120" show-overflow-tooltip>
          <template #default="scope">
            <el-button type="success" v-if="scope.row.node_type ==='vmess'">vmess</el-button>
            <el-button type="warning" v-if="scope.row.node_type ==='vless'">vless</el-button>
            <el-button type="info" v-if="scope.row.node_type ==='trojan'">trojan</el-button>
            <el-button type="danger" v-if="scope.row.node_type ==='shadowsocks'">shadowsocks</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="total_up" label="上行流量(GB)" show-overflow-tooltip width="200">
          <template #default="scope">
            <el-tag type="warning">{{ scope.row.total_up / 1024 / 1024 / 1024 }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="total_down" label="下行流量(GB)" show-overflow-tooltip width="200">
          <template #default="scope">
            <el-tag type="warning">{{ scope.row.total_down / 1024 / 1024 / 1024 }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="enable_transfer" label="节点类型" show-overflow-tooltip>
          <template #default="scope">
            <el-tag type="warning" v-if="scope.row.enable_transfer">中转</el-tag>
            <el-tag type="success" v-else>直连</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="enabled" label="节点状态" show-overflow-tooltip>
          <template #default="scope">
            <el-tag type="success" v-if="scope.row.enabled">启用</el-tag>
            <el-tag type="danger" v-else>禁用</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="node_speedlimit" label="限速" show-overflow-tooltip></el-table-column>
        <el-table-column prop="traffic_rate" label="倍率" show-overflow-tooltip></el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="scope">
            <el-button :disabled="userInfos.id !== 1" size="small" text type="primary"
                       @click="onOpenEditNode('edit', scope.row)">修改
            </el-button>
            <el-button :disabled="userInfos.id !== 1" size="small" text type="primary"
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
          v-model:current-page="state.params.page_num"
          v-model:page-size="state.params.page_size"
          :total="nodeManageData.nodes.total"
      />
    </el-card>
    <NodeDialog ref="nodeDialogRef" @refresh="onGetNode(state.params)"/>
    <NodeSortDialog ref="nodeSortDialogRef"></NodeSortDialog>
    <NodeSharedDialog ref="nodeSharedDialogRef"></NodeSharedDialog>

  </div>
</template>

<script setup lang="ts" name="NodeManage">

import {defineAsyncComponent, onMounted, reactive, ref} from "vue";
import {storeToRefs} from "pinia";
import {useNodeStore} from "/@/stores/nodeStore";
import {useUserStore} from "/@/stores/userStore";
import {ElMessageBox} from "element-plus";

const NodeDialog = defineAsyncComponent(() => import('/@/views/admin/node/dialog_edit.vue'))
const NodeSortDialog = defineAsyncComponent(() => import('/@/views/admin/node/dialog_node_sort.vue'))
const NodeSharedDialog = defineAsyncComponent(() => import('/@/views/admin/node/dialog_node_shared.vue'))
const nodeDialogRef = ref()
const nodeSortDialogRef = ref()
const nodeSharedDialogRef = ref()

const nodeStore = useNodeStore()
const {nodeManageData} = storeToRefs(nodeStore)

const userStore = useUserStore()
const {userInfos} = storeToRefs(userStore)
//时间范围
const shortcuts = [
  {
    text: '上周',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
      return [start, end]
    },
  },
  {
    text: '上月',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
      return [start, end]
    },
  },
  {
    text: '最近3个月',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 90)
      return [start, end]
    },
  },
]
//定义参数
const state = reactive({
  loading: true,
  params: {
    search: '',
    page_num: 1,
    page_size: 30,
    date: [],
  },
})

//打开新建节点，修改节点弹窗
function onOpenEditNode(title: string, nodeType: string, row?: NodeInfo) {
  nodeDialogRef.value.openDialog(title, nodeType, row)
}

//打开节点排序弹窗
function onOpenNodeSortDialog() {
  nodeSortDialogRef.value.openDialog()
}

//打开共享节点弹窗
function onOpenNodeSharedDialog() {
  nodeSharedDialogRef.value.openDialog()
}

//查询节点
function onGetNode(params?: object) {
  nodeStore.getNodeWithTraffic(params)
}

//删除节点
function onRowDel(row: NodeInfo) {
  ElMessageBox.confirm(`此操作将永久删除节点：${row.remarks}，是否继续?`, '提示', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        nodeStore.deleteNode(row)
        setTimeout(() => {
          state.params.search = ''
          onGetNode(state.params)
        }, 500);
      })
      .catch(() => {
      });
}

// 分页改变
const onHandleSizeChange = (val: number) => {
  state.params.page_size = val;
  onGetNode(state.params)
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
  state.params.page_num = val;
  onGetNode(state.params)
};

onMounted(() => {
  onGetNode(state.params)
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