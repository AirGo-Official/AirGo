<template>
  <el-config-provider :size="getGlobalComponentSize" :locale="zhCn">
    <router-view v-show="setLockScreen"/>
    <Setings ref="setingsRef" v-show="setLockScreen"/>
    <CloseFull v-if="!themeConfig.isLockScreen"/>

  </el-config-provider>
</template>

<script setup lang="ts" name="app">
import {computed, defineAsyncComponent, nextTick, onBeforeMount, onMounted, onUnmounted, ref, watch} from 'vue';
import {useRoute} from 'vue-router';
import zhCn from 'element-plus/es/locale/lang/zh-cn';
import {storeToRefs} from 'pinia';
import {useTagsViewRoutes} from '/@/stores/tagsViewRoutes';
import {useThemeConfig} from '/@/stores/themeConfig';
import other from '/@/utils/other';
import {Local, Session} from '/@/utils/storage';
import mittBus from '/@/utils/mitt';
import setIntroduction from '/@/utils/setIconfont';
import {useServerStore} from "/@/stores/serverStore";
import {useUserStore} from "/@/stores/userStore";
const userStore = useUserStore()
// 引入组件
const Setings = defineAsyncComponent(() => import('/@/layout/navBars/breadcrumb/setings.vue'));
const CloseFull = defineAsyncComponent(() => import('/@/layout/navBars/breadcrumb/closeFull.vue'));


// 定义变量内容
const setingsRef = ref();
const route = useRoute();
const stores = useTagsViewRoutes();
const storesThemeConfig = useThemeConfig();
const {themeConfig} = storeToRefs(storesThemeConfig);
const serverStore = useServerStore()

// 设置锁屏时组件显示隐藏
const setLockScreen = computed(() => {
  // 防止锁屏后，刷新出现不相关界面
  // https://gitee.com/lyt-top/vue-next-admin/issues/I6AF8P
  return themeConfig.value.isLockScreen ? themeConfig.value.lockScreenTime > 1 : themeConfig.value.lockScreenTime >= 0;
});
// 获取全局组件大小
const getGlobalComponentSize = computed(() => {
  return other.globalComponentSize();
});
//设置初始化，防止刷新时恢复默认
onBeforeMount(() => {
  // 设置批量第三方 icon 图标
  setIntroduction.cssCdn();
  // 设置批量第三方 js
  setIntroduction.jsCdn();
});

//组件被挂载之前,获取布局配置,公共配置,
onBeforeMount(() => {
  storesThemeConfig.getThemeConfig()
  serverStore.getPublicServerConfig();//获取public config
})
// 页面加载时
onMounted(() => {
  nextTick(() => {
    // 监听布局配置弹窗点击打开
    mittBus.on('openSetingsDrawer', () => {
      setingsRef.value.openDrawer();
    });
    // 获取缓存中的全屏配置
    if (Session.get('isTagsViewCurrenFull')) {
      stores.setCurrenFullscreen(Session.get('isTagsViewCurrenFull'));
    }
    //如果存在token，刷新页面时获取用户信息
    const token = Local.get('token');
    if (token) {
      userStore.getUserInfo()//获取用户信息
    }
  });
});
// 页面销毁时，关闭监听布局配置/i18n监听
onUnmounted(() => {
  mittBus.off('openSetingsDrawer', () => {
  });
});
// 监听路由的变化，设置网站标题
watch(
    () => route.path,
    () => {
      other.useTitle();
    },
    {
      deep: true,
    }
);
</script>
