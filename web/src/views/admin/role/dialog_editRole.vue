<template>
  <div class="system-role-dialog-container">
    <el-dialog :title="state.title" v-model="state.isShowDialog" width="769px" destroy-on-close>
      <el-form ref="roleDialogFormRef" :model="roleStoreData.currentRole.value" size="default" label-width="90px">
            <el-form-item label="角色名称">
              <el-input v-model="roleStoreData.currentRole.value.role_name" placeholder="请输入角色名称" clearable></el-input>
            </el-form-item>
            <el-form-item label="角色状态">
              <el-switch v-model="roleStoreData.currentRole.value.status" inline-prompt active-text="启"
                         inactive-text="禁"></el-switch>
            </el-form-item>
            <el-form-item label="角色描述">
              <el-input v-model="roleStoreData.currentRole.value.description" type="textarea" placeholder="请输入角色描述"
                        maxlength="150"></el-input>
            </el-form-item>
            <el-form-item label="菜单权限">
              <el-tree ref="menu_tree_ref" node-key="id"
                       :data="menuStoreData.allMenuTree.value"
                       :props="{children:'children',label:'title'}"
                       :default-checked-keys="roleStoreData.checkedMenuIDs.value"
                       show-checkbox class="menu-data-tree"/>
            </el-form-item>
        <el-form-item label="API权限">
          <el-tree ref="api_tree_ref" node-key="path"
                   :data="roleStoreData.allCasbinInfo.value.casbinItems"
                   :props="{label:'path'}"
                   :default-checked-keys="roleStoreData.checkedCasbinPath.value"
                   show-checkbox class="menu-data-tree"/>
        </el-form-item>
      </el-form>
      <template #footer>
				<span class="dialog-footer">
					<el-button @click="onCancel" size="default">取 消</el-button>
					<el-button type="primary" @click="onSubmit" size="default">{{ state.submitTxt }}</el-button>
				</span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import {reactive, ref} from 'vue';

import {useAdminRoleStore} from "/@/stores/admin_logic/roleStore";
import {storeToRefs} from 'pinia';
import {ElMessage} from 'element-plus';
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
import { useAdminMenuStore } from "/@/stores/admin_logic/menuStore";
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
// 定义子组件向父组件传值/事件
const emit = defineEmits(['refresh']);
const menuStore = useAdminMenuStore()
const menuStoreData = storeToRefs(menuStore)
const roleStore = useAdminRoleStore()
const roleStoreData = storeToRefs(roleStore)
const menu_tree_ref = ref()
const api_tree_ref = ref()

// 定义变量内容
const roleDialogFormRef = ref();
const state = reactive({
  isShowDialog: false,
  type: '',
  title: '',
  submitTxt: '',
})

// 打开弹窗
const openDialog = (type: string, row: RoleInfo) => {
  state.isShowDialog = true;
  state.type = type
  if (type === 'edit') {
    state.title = '修改角色';
    state.submitTxt = '修 改';

    roleStoreData.currentRole.value = row;
    roleStore.muneIDsHandler()

    roleStoreData.currentCasbin.value.roleID = row.id
    //获取当前角色 api
    roleStore.getPolicyByID()

  } else {
    state.title = '新增角色';
    state.submitTxt = '新 增';
  }
  //全部菜单tree
  menuStore.getAllMenuTree()
  //获取全部api
  roleStore.getAllPolicy()
};
// 关闭弹窗
const closeDialog = () => {
  state.isShowDialog = false;
};
// 取消
const onCancel = () => {
  closeDialog();
};
// 提交
const onSubmit = () => {
  if (state.type === 'edit') {

    const checkedMenuIDs=[...menu_tree_ref.value.getCheckedKeys(), ...menu_tree_ref.value.getHalfCheckedKeys()];
    const checkedCasbinPath=[...api_tree_ref.value.getCheckedKeys()];

    roleStore.updateRole(checkedMenuIDs,checkedCasbinPath).then((res) => {
      ElMessage.success('修改成功');
      //父组件重新加载
      emit('refresh');
    })
    //处理api权限

    //关闭编辑弹窗
    closeDialog();
  } else {

    const checkedMenuIDs=[...menu_tree_ref.value.getCheckedKeys(), ...menu_tree_ref.value.getHalfCheckedKeys()];
    const checkedCasbinPath=[...api_tree_ref.value.getCheckedKeys()];

    roleStore.newRole(checkedMenuIDs,checkedCasbinPath).then((res) => {
      ElMessage.success('新建角色成功');
      //父组件重新加载
      emit('refresh');
    })
    //处理api权限

    //关闭编辑弹窗
    closeDialog();
  }
};


// 暴露变量
defineExpose({
  openDialog,
});
</script>

<style scoped lang="scss">
.system-role-dialog-container {
  .menu-data-tree {
    width: 100%;
    border: 1px solid var(--el-border-color);
    border-radius: var(--el-input-border-radius, var(--el-border-radius-base));
    padding: 5px;
  }
}
</style>
