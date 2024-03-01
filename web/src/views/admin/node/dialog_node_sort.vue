<template>
  <el-dialog
    v-model="state.isShowDialog"
    :title="$t('message.adminNode.sortNode')"
    width="80%" destroy-on-close
    align-center
  >
    <el-table class="nodeSort" :data="nodeStoreData.nodeList.value.data" row-key="id" height="100%"
              style="width: 100%;flex: 1;">
      <el-table-column type="index" :label="$t('message.adminNode.NodeInfo.index')" show-overflow-tooltip width="60"
                       fixed></el-table-column>
      <el-table-column prop="remarks" :label="$t('message.adminNode.NodeInfo.remarks')" show-overflow-tooltip
                       width="300" fixed></el-table-column>
      <el-table-column prop="id" :label="$t('message.adminNode.NodeInfo.id')" show-overflow-tooltip width="60"
                       fixed></el-table-column>
      <el-table-column :label="$t('message.common.operate')" show-overflow-tooltip fixed>
        <el-icon class="move">
          <Rank />
        </el-icon>
      </el-table-column>
    </el-table>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="state.isShowDialog = false">{{ $t("message.common.button_cancel") }}</el-button>
        <el-button type="primary" @click="onSubmit">
          {{ $t("message.common.button_confirm") }}
        </el-button>
      </span>
    </template>
  </el-dialog>

</template>

<script lang="ts" setup>
import { nextTick, reactive } from "vue";
import Sortable from "sortablejs";
import { useAdminNodeStore } from "/@/stores/admin_logic/nodeStore";
import { request } from "/@/utils/request";
import { useApiStore } from "/@/stores/apiStore";
import { storeToRefs } from "pinia";
import { useI18n } from "vue-i18n";

const nodeStore = useAdminNodeStore();
const nodeStoreData = storeToRefs(nodeStore);
const apiStore = useApiStore();
const apiStoreData = storeToRefs(apiStore);
const { t } = useI18n();

const emit = defineEmits(["refresh", "onGetNode"]);

//定义参数
const state = reactive({
  isShowDialog: false
});
// 打开弹窗
const openDialog = () => {
  state.isShowDialog = true;
  nextTick(() => {
    initSortable("nodeSort");
  });
};
//处理排序后的节点
const nodeSortHandler = (data: Array<any>) => {
  let arr: any = [];
  data.forEach((item: any, index: number) => {
    let it = { id: 0, node_order: 0 };
    it.id = item.id;
    it.node_order = index + 1;
    arr.push(it);
  });
  return arr;
};
//确认提交
const onSubmit = () => {
  state.isShowDialog = false;
  request(apiStoreData.adminApi.value.nodeSort, nodeSortHandler(nodeStoreData.nodeList.value.data)).then((res) => {
    emit("refresh");
  });
};

// 创建sortable实例
function initSortable(className: string) {
  // 获取表格row的父节点
  const table = document.querySelector("." + className + " .el-table__body-wrapper tbody");
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
      const currRow = nodeStoreData.nodeList.value.data.splice(evt.oldIndex, 1)[0];
      nodeStoreData.nodeList.value.data.splice(evt.newIndex, 0, currRow);
    }
  });
};

// 暴露变量
defineExpose({
  openDialog   // 打开弹窗
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