<template>
  <div>
    <el-dialog :title="state.title" v-model="state.isShowDialog" width="769px" destroy-on-close>
      <el-form ref="roleDialogFormRef" :model="roleStoreData.currentRole.value" size="default" label-position="top">
            <el-form-item :label="$t('message.adminRole.RoleInfo.role_name')">
              <el-input v-model="roleStoreData.currentRole.value.role_name" clearable></el-input>
            </el-form-item>
            <el-form-item :label="$t('message.adminRole.RoleInfo.status')">
              <el-switch v-model="roleStoreData.currentRole.value.status" inline-prompt :active-text="$t('message.common.enable')"
                         :inactive-text="$t('message.common.disable')"></el-switch>
            </el-form-item>
            <el-form-item :label="$t('message.adminRole.RoleInfo.description')">
              <el-input v-model="roleStoreData.currentRole.value.description" type="textarea" maxlength="150"></el-input>
            </el-form-item>
        <!--                       :data="menuStoreData.allMenuTree.value"-->
            <el-form-item :label="$t('message.adminRole.RoleInfo.menus')">
              <el-tree ref="menu_tree_ref" node-key="id"
                       :data="menuStoreData.allMenuList.value"
                       :props="{children:'children',label:'name'}"
                       :default-checked-keys="roleStoreData.checkedMenuIDs.value"
                       show-checkbox class="data-tree">
                <template #default="{ node, data }">
                  <span class="custom-tree-node">
                    <span>{{ $t(data.meta.title) }}</span>
                  </span>
                </template>
              </el-tree>
            </el-form-item>
        <el-form-item :label="$t('message.adminRole.RoleInfo.casbins')">
          <el-tree ref="api_tree_ref" node-key="path"
                   :data="roleStoreData.allCasbinInfo.value.casbinItems"
                   :props="{label:'path'}"
                   :default-checked-keys="roleStoreData.checkedCasbinPath.value"
                   show-checkbox class="data-tree"/>
        </el-form-item>
      </el-form>
      <template #footer>
				<span class="dialog-footer">
					<el-button @click="onCancel" size="default">{{$t('message.common.button_cancel')}}</el-button>
					<el-button type="primary" @click="onSubmit" size="default">{{$t('message.common.button_confirm')}}</el-button>
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
import { useI18n } from "vue-i18n";
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
const {t} = useI18n()

// 定义变量内容
const roleDialogFormRef = ref();
const state = reactive({
  isShowDialog: false,
  type: '',
  title: '',
})

// 打开弹窗
const openDialog = (type: string, row: RoleInfo) => {
  state.isShowDialog = true;
  state.type = type
  if (type === 'edit') {
    state.title = t('message.adminRole.modify_role')

    roleStoreData.currentRole.value = row;
    roleStore.muneIDsHandler()

    roleStoreData.currentCasbin.value.roleID = row.id
    //获取当前角色 api
    roleStore.getPolicyByID()

  } else {
    state.title = t('message.adminRole.add_role')
  }
  //获取全部菜单
  menuStore.getAllMenuList()

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
      //父组件重新加载
      emit('refresh');
    })
    //关闭编辑弹窗
    closeDialog();
  } else {
    const checkedMenuIDs=[...menu_tree_ref.value.getCheckedKeys(), ...menu_tree_ref.value.getHalfCheckedKeys()];
    const checkedCasbinPath=[...api_tree_ref.value.getCheckedKeys()];
    roleStore.newRole(checkedMenuIDs,checkedCasbinPath).then((res) => {
      //父组件重新加载
      emit('refresh');
    })
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

  .data-tree {
    width: 100%;
    border: 1px solid var(--el-border-color);
    border-radius: var(--el-input-border-radius, var(--el-border-radius-base));
    padding: 5px;
  }
  .custom-tree-node {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 14px;
    padding-right: 8px;
  }

</style>
