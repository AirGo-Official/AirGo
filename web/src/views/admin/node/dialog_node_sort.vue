<template>
  <el-dialog
      v-model="state.isShowDialog"
      :title="state.title"
      width="550px" destroy-on-close
      align-center
  >
    <el-table class="nodeSort" :data="state.node_list" row-key="id" height="100%" style="width: 100%;flex: 1;">
      <el-table-column type="index" label="序号" show-overflow-tooltip width="60" fixed></el-table-column>
      <el-table-column prop="remarks" label="节点名称" show-overflow-tooltip width="300" fixed></el-table-column>
      <el-table-column prop="id" label="节点ID" show-overflow-tooltip width="60" fixed></el-table-column>
      <el-table-column label="操作" show-overflow-tooltip fixed>
        <el-icon class="move">
          <Rank/>
        </el-icon>
      </el-table-column>
    </el-table>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="state.isShowDialog = false">取消</el-button>
        <el-button type="primary" @click="onSubmit">
          确认
        </el-button>
      </span>
    </template>
  </el-dialog>

</template>

<script lang="ts" setup>
import {nextTick, reactive} from "vue";
import Sortable from "sortablejs";
import {useNodeStore} from "/@/stores/nodeStore";
import {ElMessage} from "element-plus";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
import {storeToRefs} from "pinia";

const nodeStore = useNodeStore()

const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)

//定义参数
const state = reactive({
  type: "",
  title: "节点排序",
  isShowDialog: false,
  node_list: [] as NodeInfo[],
})
// 打开弹窗
const openDialog = () => {
  state.isShowDialog = true
  request(apiStoreData.api.value.node_getAllNode).then((res) => {
    ElMessage.success(res.msg)
    state.node_list = res.data
  })
  nextTick(() => {
    initSortable("nodeSort")
  })
}
//处理排序后的节点
const nodeSortHandler = (data: Array<any>) => {
  let arr: any = []
  data.forEach((item: any, index: number) => {
    let it = {id: 0, node_order: 0}
    it.id = item.id
    it.node_order = index + 1
    arr.push(it)
  })
  return arr
}
//确认提交
const onSubmit = () => {
  state.isShowDialog = false
  request(apiStoreData.api.value.node_nodeSort, nodeSortHandler(state.node_list)).then((res) => {
    ElMessage.success(res.msg)
    nodeStore.getNodeWithTraffic({search: '', page_num: 1, page_size: 30, date: [],})
  })
}


// 创建sortable实例
function initSortable(className: string) {
  // 获取表格row的父节点
  const table = document.querySelector('.' + className + ' .el-table__body-wrapper tbody');
  // 创建拖拽实例
  let dragTable = Sortable.create(table, {
    animation: 150, //动画
    disabled: false, // 拖拽不可用? false 启用（刚刚渲染表格的时候起作用，后面不起作用）
    handle: ".move", //指定拖拽目标，点击此目标才可拖拽元素(此例中设置操作按钮拖拽)
    // filter: ".disabled", //指定不可拖动的类名（el-table中可通过row-class-name设置行的class）
    dragClass: "dragClass", //设置拖拽样式类名
    ghostClass: "ghostClass", //设置拖拽停靠样式类名
    chosenClass: "chosenClass", //设置选中样式类名
    // 开始拖动事件
    onStart: () => {
      // console.log("开始拖动");
    },
    // 结束拖动事件
    onEnd: (evt: any) => {
      // console.log("结束拖动", `拖动前索引${evt.oldIndex}---拖动后索引${evt.newIndex}`);
      const currRow = state.node_list.splice(evt.oldIndex, 1)[0];
      state.node_list.splice(evt.newIndex, 0, currRow);
      // console.log("结束拖动", state.node_list);
    },
  });
};

// 设置表格row的class
// function tableRowClassName(row: any) {
//   if (row.disabled) {
//     return "disabled";
//   }
//   return "";
// };

// 暴露变量
defineExpose({
  openDialog,   // 打开弹窗
});


</script>

<style scoped lang="scss">

// 拖拽
.dragClass {
  background: rgba($color: #41c21a, $alpha: 0.5) !important;
}

// 停靠
.ghostClass {
  background: rgba($color: #6cacf5, $alpha: 0.5) !important;
}

// 选择
.chosenClass:hover > td {
  background: rgba($color: #f56c6c, $alpha: 0.5) !important;
}
</style>