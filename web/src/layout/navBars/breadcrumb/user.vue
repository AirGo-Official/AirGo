<template>
  <div class="layout-navbars-breadcrumb-user pr15" :style="{ flex: layoutUserFlexNum }">
    <!--    深色模式切换-->
    
    <el-switch  style="margin-right: 0.5em;" v-model="getThemeConfig.isIsDark" @change="onAddDarkChange" size="large" inline-prompt active-text="Dark" inactive-text="Light">
      <template #active-action>
      <span class="custom-active-action"><i class="ri-moon-line"></i></span>
      </template>
      <template #inactive-action>
      <span class="custom-inactive-action"><i class="ri-sun-line"></i></span>
    </template>
    </el-switch>
    <!--    组件大小切换-->
    <el-dropdown  :show-timeout="70" :hide-timeout="50" trigger="click" @command="onComponentSizeChange">
       <div class="layout-navbars-breadcrumb-user-icon">
        <i style="font-size: 1.2rem; font-weight: 600;" class="ri-font-size-2" :title="$t('message.user.title0')"></i>
      </div> 
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item command="large" :disabled="state.disabledSize === 'large'">
            {{ $t("message.user.dropdownLarge") }}
          </el-dropdown-item>
          <el-dropdown-item command="default" :disabled="state.disabledSize === 'default'">
            {{ $t("message.user.dropdownDefault") }}
          </el-dropdown-item>
          <el-dropdown-item command="small" :disabled="state.disabledSize === 'small'">
            {{ $t("message.user.dropdownSmall") }}
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
    <!--    语言切换-->
    <el-dropdown :show-timeout="70" :hide-timeout="50" trigger="click" @command="onLanguageChange">
      <div class="layout-navbars-breadcrumb-user-icon" style="font-size: 1.2rem;font-weight: 600;">
      <svg preserveAspectRatio="xMidYMid meet" viewBox="0 0 24 24" width="1.2em" height="1.2em" data-v-12008bb2=""><path fill="currentColor" d="m18.5 10l4.4 11h-2.155l-1.201-3h-4.09l-1.199 3h-2.154L16.5 10h2zM10 2v2h6v2h-1.968a18.222 18.222 0 0 1-3.62 6.301a14.864 14.864 0 0 0 2.336 1.707l-.751 1.878A17.015 17.015 0 0 1 9 13.725a16.676 16.676 0 0 1-6.201 3.548l-.536-1.929a14.7 14.7 0 0 0 5.327-3.042A18.078 18.078 0 0 1 4.767 8h2.24A16.032 16.032 0 0 0 9 10.877a16.165 16.165 0 0 0 2.91-4.876L2 6V4h6V2h2zm7.5 10.885L16.253 16h2.492L17.5 12.885z"></path></svg>
      </div>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item command="zh-cn" :disabled="state.disabledI18n === 'zh-cn'">简体中文</el-dropdown-item>
          <el-dropdown-item command="en" :disabled="state.disabledI18n === 'en'">English</el-dropdown-item>
          <el-dropdown-item command="zh-tw" :disabled="state.disabledI18n === 'zh-tw'">繁體中文</el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
    <!--    搜索-->
    <!--    <div v-if="userInfos.id===1" class="layout-navbars-breadcrumb-user-icon" @click="onSearchClick">-->
    <!--      <el-icon :title="$t('message.user.title2')">-->
    <!--        <ele-Search/>-->
    <!--      </el-icon>-->
    <!--    </div>-->
    <!--    布局设置-->
    <div v-if="userInfos.id===1" class="layout-navbars-breadcrumb-user-icon" @click="onLayoutSetingClick" style="font-size: 1.2rem;font-weight: 600;">
      <i class="ri-t-shirt-2-line" :title="$t('message.user.title3')"></i>
    </div>
    <!--    全屏
    <div class="layout-navbars-breadcrumb-user-icon mr10" @click="onScreenfullClick">
      <i
        class="iconfont"
        :title="state.isScreenfull ? $t('message.user.title6') : $t('message.user.title5')"
        :class="!state.isScreenfull ? 'icon-fullscreen' : 'icon-tuichuquanping'"
      ></i>
    </div>-->
    <!--    个人头像-->
    <el-dropdown style="margin-left: 1em;" :show-timeout="70" :hide-timeout="50" @command="onHandleCommandClick">
			<span class="layout-navbars-breadcrumb-user-link">
        
				<img :src="userInfos.avatar" class="layout-navbars-breadcrumb-user-link-photo mr5" />
				<div class="layout-navbars-username-hide">{{ userInfos.user_name }}</div>
				<el-icon class="el-icon--right">
					<ele-ArrowDown />
				</el-icon>
			</span>
      <template #dropdown>
        <el-dropdown-menu
        >
          <el-dropdown-item command="/home">{{ $t("message.user.dropdown1") }}</el-dropdown-item>
          <el-dropdown-item divided command="logOut">{{ $t("message.user.dropdown5") }}</el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
    <!--    搜索组件-->
    <Search ref="searchRef" />
  </div>
</template>

<script setup lang="ts" name="layoutBreadcrumbUser">
import { defineAsyncComponent, ref, computed, reactive, onMounted } from "vue";
import { useRouter } from "vue-router";
import { ElMessageBox, ElMessage } from "element-plus";
import screenfull from "screenfull";
import { storeToRefs } from "pinia";
import { useUserStore } from "/@/stores/user_logic/userStore";
import { useThemeConfig } from "/@/stores/themeConfig";
import other from "/@/utils/other";
import mittBus from "/@/utils/mitt";
import { Session, Local } from "/@/utils/storage";
import { useI18n } from "vue-i18n";

// 引入组件
const Search = defineAsyncComponent(() => import("/@/layout/navBars/breadcrumb/search.vue"));

// 定义变量内容
const router = useRouter();
const stores = useUserStore();
const storesThemeConfig = useThemeConfig();
const { userInfos } = storeToRefs(stores);
const { themeConfig } = storeToRefs(storesThemeConfig);
const searchRef = ref();
const { locale, t } = useI18n();

const state = reactive({
  isScreenfull: false,
  disabledI18n: "zh-cn",
  disabledSize: "default"
});

// 获取布局配置信息
const getThemeConfig = computed(() => {
  return themeConfig.value;
});
// 界面显示 --> 深色模式
const onAddDarkChange = () => {
  const body = document.documentElement as HTMLElement;
  if (getThemeConfig.value.isIsDark) body.setAttribute('data-theme', 'dark');
  else body.setAttribute('data-theme', '');
  //更新session中的theme
  setLocalThemeConfig()
};
// 存储布局配置
const setLocalThemeConfig = () => {
  Session.remove('themeConfig');
  Session.set('themeConfig', getThemeConfig.value);
};

// 设置分割样式
const layoutUserFlexNum = computed(() => {
  let num: string | number = "";
  const { layout, isClassicSplitMenu } = themeConfig.value;
  const layoutArr: string[] = ["defaults", "columns"];
  if (layoutArr.includes(layout) || (layout === "classic" && !isClassicSplitMenu)) num = "1";
  else num = "";
  return num;
});
// 全屏点击时
const onScreenfullClick = () => {
  if (!screenfull.isEnabled) {
    ElMessage.warning("暂不不支持全屏");
    return false;
  }
  screenfull.toggle();
  screenfull.on("change", () => {
    if (screenfull.isFullscreen) state.isScreenfull = true;
    else state.isScreenfull = false;
  });
};
// 布局配置 icon 点击时
const onLayoutSetingClick = () => {
  mittBus.emit("openSetingsDrawer");
};
// 下拉菜单点击时
const onHandleCommandClick = (path: string) => {
  if (path === "logOut") {
    ElMessageBox({
      closeOnClickModal: false,
      closeOnPressEscape: false,
      title: t('message.user.logOutTitle'),
      message: t('message.user.logOutMessage'),
      showCancelButton: true,
      confirmButtonText: t('message.user.logOutConfirm'),
      cancelButtonText: t('message.user.logOutCancel'),
      buttonSize: 'default',
      beforeClose: (action, instance, done) => {
        if (action === 'confirm') {
          instance.confirmButtonLoading = true;
          instance.confirmButtonText = t('message.user.logOutExit');
          setTimeout(() => {
            done();
            setTimeout(() => {
              instance.confirmButtonLoading = false;
            }, 300);
          }, 700);
        } else {
          done();
        }
      },
    })
      .then(async () => {
        // 清除缓存/token等
        Local.clear();
        // 使用 reload 时，不需要调用 resetRoute() 重置路由
        window.location.reload();
      })
      .catch(() => {
      });
  } else if (path === "wareHouse") {
    window.open("https://gitee.com/lyt-top/vue-next-admin");
  } else {
    router.push(path);
  }
};

// 组件大小改变
const onComponentSizeChange = (size: string) => {
  Session.remove("themeConfig");
  themeConfig.value.globalComponentSize = size;
  Session.set("themeConfig", themeConfig.value);
  initI18nOrSize("globalComponentSize", "disabledSize");
  window.location.reload();
};
// 初始化组件大小/i18n
const initI18nOrSize = (value: string, attr: string) => {
  (<any>state)[attr] = Session.get("themeConfig")[value];
};
// 语言切换
const onLanguageChange = (lang: string) => {
  Session.remove("themeConfig");
  themeConfig.value.globalI18n = lang;
  Session.set("themeConfig", themeConfig.value);
  locale.value = lang;
  other.useTitle();
  initI18nOrSize("globalI18n", "disabledI18n");
};
// 页面加载时
onMounted(() => {
  if (Session.get("themeConfig")) {
    initI18nOrSize("globalComponentSize", "disabledSize");
    initI18nOrSize("globalI18n", "disabledI18n");
  }
  //session读取用户信息到pinia
  userInfos.value = Session.get("userInfos");

});
</script>

<style scoped lang="scss">
.layout-navbars-breadcrumb-user {
  display: flex;
  align-items: center;
  justify-content: flex-end;

  &-link {
    height: 100%;
    display: flex;
    align-items: center;
    white-space: nowrap;

    &-photo {
      width: 25px;
      height: 25px;
      border-radius: 100%;
    }
  }

  &-icon {
    padding: 0 10px;
    cursor: pointer;
    color: var(--next-bg-topBarColor);
    height: 50px;
    line-height: 50px;
    display: flex;
    align-items: center;

    &:hover {
      background: var(--next-color-user-hover);
      i {
        display: inline-block;
        animation: logoAnimation 0.3s ease-in-out;
      }
    }
  }

  :deep(.el-dropdown) {
    color: var(--next-bg-topBarColor);
  }

  :deep(.el-badge) {
    height: 40px;
    line-height: 40px;
    display: flex;
    align-items: center;
  }

  :deep(.el-badge__content.is-fixed) {
    top: 12px;
  }
}
</style>
