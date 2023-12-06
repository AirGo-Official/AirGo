<template>
  <div>
    <el-dialog v-model="state.isShowDialog" :title="state.title" width="80%" destroy-on-close align-center>
      <el-form v-model="ticketStoreData.newTicketInfo.value" size="default" label-position="top">
        <el-form-item label="标题">
          <el-input v-model="ticketStoreData.newTicketInfo.value.title"></el-input>
        </el-form-item>
        <el-form-item label="详情">
          <el-input v-model="ticketStoreData.newTicketInfo.value.details" type="textarea" autosize></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
				<span class="dialog-footer">
					<el-button @click="closeDialog" size="default">取 消</el-button>
					<el-button type="primary" @click="onSubmit" size="default">提交</el-button>
				</span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">

import {reactive} from "vue";
import {useTicketStore} from "/@/stores/ticketStore";
import {storeToRefs} from "pinia";

const ticketStore = useTicketStore()
const ticketStoreData = storeToRefs(ticketStore)

const state = reactive({
  isShowDialog: false,
  title: "新建工单",
})

// 定义子组件向父组件传值/事件
const emit = defineEmits(['refresh']);

// 打开弹窗
const openDialog = (row?: any) => {
  state.isShowDialog = true
}
// 关闭弹窗
const closeDialog = () => {
  state.isShowDialog = false
};
//确认提交
const onSubmit = () => {
  ticketStore.newTicket()
  setTimeout(()=>{
    emit('refresh')
  },2000)
  closeDialog()
}
// 暴露变量
defineExpose({
  openDialog,   // 打开弹窗
});

</script>

<style scoped lang="scss">

</style>