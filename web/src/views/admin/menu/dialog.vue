<template>
  <div class="system-menu-dialog-container">
    <el-dialog :title="state.dialog.title" v-model="state.dialog.isShowDialog" width="769px" destroy-on-close>
      <el-form ref="menuDialogFormRef" :model="state.menuForm" size="default" label-width="80px">
        <el-row :gutter="35">
          <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
            <el-form-item label="父级菜单">
              <el-cascader :options="allRoutesList"
                           @change="changeCheckMenu"
                           :props="{ checkStrictly: true, value: 'id', label: 'name' }" placeholder="请选择父级菜单"
                           clearable class="w100">
                <template #default="{ node, data }">
                  <span>{{ data.meta.title }}</span>
                  <span v-if="!node.isLeaf"> ({{ data.children.length }}) </span>
                </template>
              </el-cascader>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
            <el-form-item label="父菜单ID">
              <el-input v-model="state.menuForm.parent_id"></el-input>
            </el-form-item>
          </el-col>
          <!--					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">-->
          <!--						<el-form-item label="菜单类型">-->
          <!--							<el-radio-group v-model="state.menuForm.menuType">-->
          <!--								<el-radio label="menu">菜单</el-radio>-->
          <!--								<el-radio label="btn">按钮</el-radio>-->
          <!--							</el-radio-group>-->
          <!--						</el-form-item>-->
          <!--					</el-col>-->
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="菜单名称">
              <el-input v-model="state.menuForm.meta.title" placeholder="例如：用户管理"
                        clearable></el-input>
            </el-form-item>
          </el-col>
          <template v-if="state.menuForm.menuType === 'menu'">
            <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
              <el-form-item label="路由名称">
                <el-input v-model="state.menuForm.name" placeholder="例如：userManage" clearable></el-input>
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
              <el-form-item label="路由路径">
                <el-input v-model="state.menuForm.path" placeholder="例如：/admin/userManage" clearable></el-input>
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
              <el-form-item label="重定向">
                <el-input v-model="state.menuForm.redirect" placeholder="请输入路由重定向" clearable></el-input>
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
              <el-form-item label="菜单图标">
                <IconSelector placeholder="请输入菜单图标" v-model="state.menuForm.meta.icon"/>
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
              <el-form-item label="组件路径">
                <el-input v-model="state.menuForm.component" placeholder="例如：/login/index.vue" clearable></el-input>
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
              <el-form-item label="外链地址">
                <el-input v-model="state.menuForm.meta.isLink" placeholder="外链/内嵌时链接地址（http:xxx.com）"
                          clearable>
                </el-input>
              </el-form-item>
            </el-col>
          </template>
          <template v-if="state.menuForm.menuType === 'menu'">
            <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
              <el-form-item label="外链内嵌">
                <el-radio-group v-model="state.menuForm.meta.isIframe">
                  <el-radio :label="true">是</el-radio>
                  <el-radio :label="false">否</el-radio>
                </el-radio-group>
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
              <el-form-item label="是否隐藏">
                <el-radio-group v-model="state.menuForm.meta.isHide">
                  <el-radio :label="true">隐藏</el-radio>
                  <el-radio :label="false">不隐藏</el-radio>
                </el-radio-group>
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
              <el-form-item label="页面缓存">
                <el-radio-group v-model="state.menuForm.meta.isKeepAlive">
                  <el-radio :label="true">缓存</el-radio>
                  <el-radio :label="false">不缓存</el-radio>
                </el-radio-group>
              </el-form-item>
            </el-col>
            <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
              <el-form-item label="是否固定">
                <el-radio-group v-model="state.menuForm.meta.isAffix">
                  <el-radio :label="true">固定</el-radio>
                  <el-radio :label="false">不固定</el-radio>
                </el-radio-group>
              </el-form-item>
            </el-col>

          </template>
        </el-row>
      </el-form>
      <template #footer>
				<span class="dialog-footer">
					<el-button @click="onCancel" size="default">取 消</el-button>
					<el-button type="primary" @click="onSubmit" size="default">{{ state.dialog.submitTxt }}</el-button>
				</span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import {defineAsyncComponent, onMounted, reactive, ref} from 'vue';
import {storeToRefs} from 'pinia';
import {useRoutesStore} from '/@/stores/routesStore';
// import { setBackEndControlRefreshRoutes } from "/@/router/backEnd";
const emit = defineEmits(['refresh']);
const IconSelector = defineAsyncComponent(() => import('/@/components/iconSelector/index.vue'));
const menuDialogFormRef = ref();
const stores = useRoutesStore();
const {allRoutesList} = storeToRefs(stores);
const state = reactive({
  // 参数请参考 `/src/router/route.ts` 中的 `dynamicRoutes` 路由菜单格式
  menuForm: {
    menuType: 'menu', // 菜单类型
    name: '', // 路由名称
    component: '', // 组件路径
    parent_id: 0,
    path: '', // 路由路径
    redirect: '', // 路由重定向，有子集 children 时
    roles: [],//角色
    meta: {
      title: '', // 菜单名称
      icon: '', // 菜单图标
      isHide: false, // 是否隐藏
      isKeepAlive: true, // 是否缓存
      isAffix: false, // 是否固定
      isLink: '', // 外链/内嵌时链接地址（http:xxx.com），开启外链条件，`1、isLink: 链接地址不为空`
      isIframe: false, // 是否内嵌，开启条件，`1、isIframe:true 2、isLink：链接地址不为空`
    },
  },
  //menuData: [] as RouteItems, // 上级菜单数据
  menuData: [] as Route[], // 上级菜单数据
  dialog: {
    isShowDialog: false,
    type: '',
    title: '',
    submitTxt: '',
  },
});


// 打开弹窗
const openDialog = (type: string, row?: any) => {
  if (type === 'edit') {
    row.menuType = 'menu';
    // row.menuSort = Math.floor(Math.random() * 100);
    state.menuForm = JSON.parse(JSON.stringify(row));
    state.dialog.title = '修改菜单';
    state.dialog.submitTxt = '修 改';
  } else {
    state.dialog.title = '新增菜单';
    state.dialog.submitTxt = '新 增';
    // 清空表单，此项需加表单验证才能使用
    // nextTick(() => {
    // 	menuDialogFormRef.value.resetFields();
    // });
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
  state.menuForm.parent_id = checkValue[checkValue.length - 1]
}
// 提交
const onSubmit = () => {
  if (state.dialog.type === 'add') {
    stores.newDynamicRoute(state.menuForm)
    setTimeout(() => {
      emit('refresh');
    }, 1000);
    closeDialog();
  } else {
    //请求
    stores.updateDynamicRoute(state.menuForm)
    setTimeout(() => {
      emit('refresh');
    }, 1000);
    closeDialog(); // 关闭弹窗
  }
};
// 页面加载时
onMounted(() => {
  stores.setAllRoutesList()
  stores.setRoutesTree()
});

// 暴露变量
defineExpose({
  openDialog,
});
</script>
