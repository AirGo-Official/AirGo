<template>
  <el-dialog v-model="state.isShowDialog" title="修改角色api权限" height="500px" destroy-on-close>
    <el-transfer v-model="dialogEditApi.casbinInfo.casbinItems"
                 :props="{
      key: 'path',
      label: 'path',
    }"
                 :data="dialogEditApi.allCasbinInfo.casbinItems"/>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="state.isShowDialog = false">取消</el-button>
        <el-button type="primary" @click="onSubmit">
          提交
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import {storeToRefs} from 'pinia';
import {useRoleStore} from "/@/stores/roleStore";
import {reactive} from "vue";

const roleStore = useRoleStore()
const {dialogEditApi} = storeToRefs(roleStore)

// 定义变量内容
const state = reactive({
  isShowDialog: false,
  type: '',
  title: '',
  submitTxt: '',
})

// 打开弹窗
const openDialog = (row: RowRoleType) => {
  state.isShowDialog = true;
  //获取当前roleID
  dialogEditApi.value.casbinInfo.roleID = row.id
  //获取全部api
  roleStore.getAllPolicy()
  //获取当前角色 api (选中)
  roleStore.getPolicyByRoleIds({roleID: row.id})
};
// 关闭弹窗
const closeDialog = () => {
  state.isShowDialog = false;
};
//提交
const onSubmit = () => {
  roleStore.updateCasbinPolicy(dialogEditApi.value.casbinInfo)
  closeDialog()
}

// 暴露变量
defineExpose({
  openDialog,
});


</script>
<style>
/* 定义两边的el-transfer-panel大小的方法,直接设置是没有用的,需要去掉scoped即可。才能成功覆盖原生的样式 */
.el-transfer-panel {
  width: 280px;
  height: 600px;
}

.el-transfer-panel__body {
  height: 600px;
}

.el-transfer-panel__list {
  height: 550px;
}

/*穿梭框内部展示列表的高宽度*/
/*:deep(.el-transfer-panel__list.is-filterable){*/
/*  width:280px;*/
/*  height:500px;*/
/*}*/
</style>
  