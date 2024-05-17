<template>
  <div v-if="isShowBreadcrumb" class="layout-navbars-breadcrumb">
    <div class="layout-memu-bottom-display" style="margin:1.8vh 0 1.8vh 1.8vh; font-size: 1.3rem;color: var(--next-bg-topBarColor);font-weight: 600;display: none;">
    <i class="ri-menu-line" @click="onThemeConfigChange"></i></div>

    <el-breadcrumb>
      <transition-group name="breadcrumb">
        <el-breadcrumb-item class="layout-breadcrumb-mobile" style="font-size: 1.4rem;margin-left: 1.6em;" v-for="(v, k) in state.breadcrumbList"
                            :key="!v.meta.tagsViewName ? v.meta.title : v.meta.tagsViewName">
					<span v-if="k === state.breadcrumbList.length - 1" class="layout-navbars-breadcrumb-span">

<!--						<div v-if="!v.meta.tagsViewName">{{ $t(v.meta.title) }}</div>-->
<!--						<div v-else>{{ v.meta.tagsViewName }}</div>-->
						<div style="font-weight: 600;color: var(--el-text-color-primary) !important;">{{ $t(v.meta.title) }}</div>
					</span>
          <a v-else>
            {{ $t(v.meta.title) }}
          </a>
        </el-breadcrumb-item>
      </transition-group>
    </el-breadcrumb>
  </div>
</template>

<script setup lang="ts" name="layoutBreadcrumb">
import {reactive, computed, onMounted} from 'vue';
import {onBeforeRouteUpdate, useRoute, useRouter} from 'vue-router';
import {Local,Session} from '/@/utils/storage';
import other from '/@/utils/other';
import {storeToRefs} from 'pinia';
import {useThemeConfig} from '/@/stores/themeConfig';
import {useMenuStore} from '/@/stores/menuStore';

// 定义变量内容
const stores = useMenuStore();
const storesThemeConfig = useThemeConfig();
const {themeConfig} = storeToRefs(storesThemeConfig);
const {routesListSate} = storeToRefs(stores);
const route = useRoute();
const router = useRouter();
const state = reactive<BreadcrumbState>({
  breadcrumbList: [],
  routeSplit: [],
  routeSplitFirst: '',
  routeSplitIndex: 1,
});

// 动态设置经典、横向布局不显示
const isShowBreadcrumb = computed(() => {
  initRouteSplit(route.path);
  const {layout, isBreadcrumb} = themeConfig.value;
  if (layout === 'classic' || layout === 'transverse') return false;
  else return isBreadcrumb ? true : false;
});

// 展开/收起左侧菜单点击
const onThemeConfigChange = () => {
  themeConfig.value.isCollapse = !themeConfig.value.isCollapse;
  setLocalThemeConfig();
};
// 存储布局配置
const setLocalThemeConfig = () => {
  Session.remove('themeConfig');
  Session.set('themeConfig', themeConfig.value);
};
// 处理面包屑数据
const getBreadcrumbList = (arr: RouteItems) => {
  arr.forEach((item: RouteItem) => {
    state.routeSplit.forEach((v: string, k: number, arrs: string[]) => {
      if (state.routeSplitFirst === item.path) {
        state.routeSplitFirst += `/${arrs[state.routeSplitIndex]}`;
        state.breadcrumbList.push(item);
        state.routeSplitIndex++;
        if (item.children) getBreadcrumbList(item.children);
      }
    });
  });
};
// 当前路由字符串切割成数组，并删除第一项空内容
const initRouteSplit = (path: string) => {
  if (!themeConfig.value.isBreadcrumb) return false;
  state.breadcrumbList = [routesListSate.value.routesList[0]];
  state.routeSplit = path.split('/');
  state.routeSplit.shift();
  state.routeSplitFirst = `/${state.routeSplit[0]}`;
  state.routeSplitIndex = 1;
  getBreadcrumbList(routesListSate.value.routesList);
  if (route.name === 'home' || (route.name === 'notFound' && state.breadcrumbList.length > 1)) state.breadcrumbList.shift();
  if (route.name === 'shop' || (route.name === 'notFound' && state.breadcrumbList.length > 1)) state.breadcrumbList.shift();
  if (route.name === 'myOrder' || (route.name === 'notFound' && state.breadcrumbList.length > 1)) state.breadcrumbList.shift();
  if (route.name === 'personal' || (route.name === 'notFound' && state.breadcrumbList.length > 1)) state.breadcrumbList.shift();
  if (route.name === 'documents' || (route.name === 'notFound' && state.breadcrumbList.length > 1)) state.breadcrumbList.shift();
  if (route.name === 'ticket' || (route.name === 'notFound' && state.breadcrumbList.length > 1)) state.breadcrumbList.shift();
  if (route.name === 'finance' || (route.name === 'notFound' && state.breadcrumbList.length > 1)) state.breadcrumbList.shift();


};
// 页面加载时
onMounted(() => {
  initRouteSplit(route.path);
});
// 路由更新时
onBeforeRouteUpdate((to) => {
  initRouteSplit(to.path);
});
</script>

<style scoped lang="scss">
.layout-navbars-breadcrumb {
  flex: 1;
  height: inherit;
  display: flex;
  align-items: center;



  .layout-navbars-breadcrumb-span {
    display: flex;
    opacity: 0.7;
    color: var(--next-bg-topBarColor);
  }

  .layout-navbars-breadcrumb-iconfont {
    margin-right: 5px;
  }

  :deep(.el-breadcrumb__separator) {
    opacity: 0.7;
    color: var(--next-bg-topBarColor);
  }

  :deep(.el-breadcrumb__inner a, .el-breadcrumb__inner.is-link) {
    font-weight: unset !important;
    color: var(--next-bg-topBarColor);

    &:hover {
      color: var(--el-color-primary) !important;
    }
  }
}
</style>
