<template>
  <div class="system-menu-dialog-container">
    <el-dialog :title="state.dialog.title" v-model="state.dialog.isShowDialog" width="769px" destroy-on-close>
      <el-form ref="menuDialogFormRef" :model="menuStoreData.currentMenu.value" size="default" value-width="80px" value-position="top">
            <el-form-item :value="$t('message.adminMenu.Route.parent_menu')">
              <el-cascader :options="menuStoreData.allMenuList.value"
                           @change="changeCheckMenu"
                           :props="{ checkStrictly: true, value: 'id', lable: 'name' }"
                           clearable class="w100">
                <template #default="{ node, data }">
                  <span>{{ $t(data.meta.title) }}</span>
                  <span v-if="!node.isLeaf"> ({{ data.children.length }}) </span>
                </template>
              </el-cascader>
            </el-form-item>

            <el-form-item :value="$t('message.adminMenu.Route.parent_id')">
              <el-input v-model="menuStoreData.currentMenu.value.parent_id"></el-input>
            </el-form-item>

            <el-form-item :value="$t('message.adminMenu.Route.remarks')">
              <el-input v-model="menuStoreData.currentMenu.value.remarks"></el-input>
            </el-form-item>

            <el-form-item :value="$t('message.adminMenu.Route.title')">
              <el-input v-model="menuStoreData.currentMenu.value.meta.title"
                        clearable></el-input>
            </el-form-item>
          <template v-if="true">
              <el-form-item :value="$t('message.adminMenu.Route.name')">
                <el-input v-model="menuStoreData.currentMenu.value.name" clearable></el-input>
              </el-form-item>
              <el-form-item :value="$t('message.adminMenu.Route.path')">
                <el-input v-model="menuStoreData.currentMenu.value.path" clearable></el-input>
              </el-form-item>
              <el-form-item :value="$t('message.adminMenu.Route.icon')">
                <IconSelector  v-model="menuStoreData.currentMenu.value.meta.icon"/>
              </el-form-item>
              <el-form-item :value="$t('message.adminMenu.Route.component')">
                <el-input v-model="menuStoreData.currentMenu.value.component"  clearable></el-input>
              </el-form-item>
              <el-form-item :value="$t('message.adminMenu.Route.isLink')">
                <el-input v-model="menuStoreData.currentMenu.value.meta.isLink"
                          clearable>
                </el-input>
              </el-form-item>
          </template>
          <template v-if="true">
              <el-form-item :value="$t('message.adminMenu.Route.isIframe')">
                <el-radio-group v-model="menuStoreData.currentMenu.value.meta.isIframe">
                  <el-radio :value="true">{{$t('message.common.yes')}}</el-radio>
                  <el-radio :value="false">{{$t('message.common.no')}}</el-radio>
                </el-radio-group>
              </el-form-item>

          </template>
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
import {defineAsyncComponent, onMounted, reactive, ref} from 'vue';
import {storeToRefs} from 'pinia';
import { useAdminMenuStore } from "/@/stores/admin_logic/menuStore";
import { useI18n } from "vue-i18n";
// import { setBackEndControlRefreshRoutes } from "/@/router/backEnd";
const emit = defineEmits(['refresh']);
const IconSelector = defineAsyncComponent(() => import('/@/components/iconSelector/index.vue'));
const menuDialogFormRef = ref();
const menuStore = useAdminMenuStore()
const menuStoreData = storeToRefs(menuStore);
const {t} = useI18n()
const state = reactive({
  dialog: {
    isShowDialog: false,
    type: '',
    title: '',
  },
});


// 打开弹窗
const openDialog = (type: string, row?: any) => {
  if (type === 'edit') {
    row.menuType = 'menu';
    // menuStoreData.currentMenu.value = JSON.parse(JSON.stringify(row));
    menuStoreData.currentMenu.value = row;
    state.dialog.title = t('message.adminMenu.modify_menu');
  } else {
    state.dialog.title = t('message.adminMenu.add_menu');
  }
  state.dialog.type = type;
  state.dialog.isShowDialog = true;
};
// 关闭弹窗
const closeDialog = () => {
  state.dialog.isShowDialog = false;
};
// 取消
const onCancel = () => {
  closeDialog();
};
//获取父级ID
const changeCheckMenu = (checkValue: any) => {
  // console.log("checkValue:", checkValue)
  menuStoreData.currentMenu.value.parent_id = checkValue[checkValue.length - 1]
}
// 提交
const onSubmit = () => {
  if (state.dialog.type === 'add') {
    menuStore.newMenu()
    setTimeout(() => {
      emit('refresh');
    }, 1000);
    closeDialog();
  } else {
    //请求
    menuStore.updateMenu()
    setTimeout(() => {
      emit('refresh');
    }, 1000);
    closeDialog(); // 关闭弹窗
  }
};
// 页面加载时
onMounted(() => {
  menuStore.getAllMenuList()
});

// 暴露变量
defineExpose({
  openDialog,
});
</script>
