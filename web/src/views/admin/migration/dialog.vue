<template>
  <el-dialog v-model="state.isShowDialog" :title="state.title" width="80%" destroy-on-close align-center>

  </el-dialog>
</template>

<script setup lang="ts">

import {reactive} from "vue";

const state = reactive({
  isShowDialog:false,
  title: "",
})


// 打开弹窗
const openDialog = (title: string, row?: any) => {
    state.isShowDialog = true
}
// 关闭弹窗
const closeDialog = () => {
  state.isShowDialog = false
};
//确认提交
function onSubmit() {

}
// 暴露变量
defineExpose({
  openDialog,   // 打开弹窗
});

</script>

<style scoped lang="scss">

</style>