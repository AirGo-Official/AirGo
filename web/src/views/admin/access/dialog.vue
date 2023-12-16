<template>
  <el-dialog v-model="state.isShowDialog" :title="state.title" width="80%" destroy-on-close align-center>
    <el-form label-position="top">
      <el-form-item label="名称">
        <el-input v-model="accessStoreData.current_routes.value.name"/>
      </el-form-item>
      <el-form-item label="禁止路由, 多个数据请换行">
        <el-input v-model="accessStoreData.current_routes.value.route" type="textarea" autosize/>
      </el-form-item>
    </el-form>
    <template #footer>
            <span class="dialog-footer">
                <el-button @click="closeDialog">取消</el-button>
                <el-button type="primary" @click="onSubmit">
                    确认
                </el-button>
            </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">

import {reactive} from "vue";
import {useAccessStore} from "/@/stores/accessStore";
import {storeToRefs} from "pinia";

const accessStore = useAccessStore()
const accessStoreData = storeToRefs(accessStore)


const state = reactive({
  isShowDialog: false,
  title: "",
})

// 打开弹窗
const openDialog = (title: string, row?: any) => {
  state.isShowDialog = true
  state.title = title
  switch (title) {
    case "新建":
      break
    default:
      accessStoreData.current_routes.value = row
      break
  }
}
// 关闭弹窗
const closeDialog = () => {
  state.isShowDialog = false
};

//确认提交
function onSubmit() {
  switch (state.title) {
    case "新建":
      accessStore.newRoutes(accessStoreData.current_routes.value)
      break
    default:
      accessStore.updateRoutes(accessStoreData.current_routes.value)
      break
  }
  setTimeout(()=>{
    accessStore.getRoutesList(accessStoreData.params.value)
  },500)
  closeDialog()
}

// 暴露变量
defineExpose({
  openDialog,   // 打开弹窗
});

</script>

<style scoped lang="scss">

</style>