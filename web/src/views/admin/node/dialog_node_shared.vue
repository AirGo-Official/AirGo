<template>
  <el-dialog v-model="state.isShowDialog" :title="state.title" width="769px" destroy-on-close align-center>
    <el-card shadow="hover" class="layout-padding-auto" style="height: 800px">
      <el-row :gutter="10" style="width: 768px">
        <el-col :span="18">
          <el-input v-model="nodeStore.nodeSharedData.newNodeSharedUrl.url" size="default"
          ></el-input>
        </el-col>
        <el-col :span="2">
          <el-button size="default" type="success" class="ml10" @click="newNodeShared()">
            <el-icon>
              <ele-FolderAdd/>
            </el-icon>
            新增节点
          </el-button>
        </el-col>
      </el-row>
      <div style="color: #9b9da1">*支持订阅，单个节点，多个节点，支持base64编码，仅解析vmess,vless,trojan</div>
      <el-table :data="nodeStore.nodeSharedData.nodeList" height="100%" stripe style="width: 100%;flex: 1;">
        <el-table-column prop="id" label="ID" show-overflow-tooltip width="40"></el-table-column>
        <el-table-column prop="node_type" label="类型" show-overflow-tooltip width="80"></el-table-column>
        <el-table-column prop="remarks" label="别名" show-overflow-tooltip width="200"></el-table-column>
        <el-table-column prop="address" label="地址" show-overflow-tooltip width="200"></el-table-column>
        <el-table-column prop="port" label="端口" show-overflow-tooltip width="80"></el-table-column>
        <el-table-column prop="network" label="传输协议" show-overflow-tooltip width="80"></el-table-column>
        <el-table-column prop="host" label="混淆" show-overflow-tooltip width="200"></el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="scope">
            <el-button size="small" text type="primary"
                       @click="deleteNodeShared(scope.row)">删除
            </el-button>
          </template>
        </el-table-column>

      </el-table>
    </el-card>

  </el-dialog>
</template>


<script setup lang="ts">
//定义参数
import {reactive} from "vue";
import {useNodeStore} from "/@/stores/nodeStore";
import {ElMessage} from "element-plus";

const nodeStore = useNodeStore()

const state = reactive({
  type: "",
  title: "共享节点管理",
  isShowDialog: false,
})

// 打开弹窗
const openDialog = () => {
  state.isShowDialog = true
  nodeStore.getNodeSharedList()
}
//新建共享节点
const newNodeShared = () => {
  nodeStore.newNodeShared(nodeStore.nodeSharedData.newNodeSharedUrl).then((res) => {
    ElMessage.success(res.msg)
    getNodeSharedList()
  })

}
//获取共享节点列表
const getNodeSharedList = () => {
  nodeStore.getNodeSharedList()
}
//删除共享节点
const deleteNodeShared = (raw: NodeSharedInfo) => {
  nodeStore.deleteNodeShared(raw).then((res) => {
    ElMessage.success(res.msg)
    getNodeSharedList()
  })
}
// 暴露变量
defineExpose({
  openDialog,   // 打开弹窗
});

</script>
<style scoped lang="scss">

</style>