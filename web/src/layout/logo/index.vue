<template>
<div class="logo_all" >
<div>
  <div class="layout-logo" v-if="setShowLogo" @click="onThemeConfigChange">
    <img class="memu-logo" :src="themeConfig.logo_link"/>
    <el-text class="menu-title" >{{ themeConfig.globalTitle }}</el-text>
  </div>
  
  <div class="layout-logo-size" v-else @click="onThemeConfigChange">
    <el-image style="width: 65%;" :src="themeConfig.logo_link" fit="cover" />
  </div>

</div>
</div>
</template>

<script setup lang="ts" name="layoutLogo">
import {computed} from 'vue';
import {storeToRefs} from 'pinia';
//store
import {useThemeConfig} from '/@/stores/themeConfig';

const storesThemeConfig = useThemeConfig();
const {themeConfig} = storeToRefs(storesThemeConfig);
//import logoMini from '/@/assets/logo-mini.svg';


// 设置 logo 的显示。classic 经典布局默认显示 logo
const setShowLogo = computed(() => {
  let {isCollapse, layout} = themeConfig.value;
  return !isCollapse || layout === 'classic' || document.body.clientWidth < 1000;
});
// logo 点击实现菜单展开/收起
const onThemeConfigChange = () => {
  if (themeConfig.value.layout === 'transverse') return false;
  themeConfig.value.isCollapse = !themeConfig.value.isCollapse;
};
</script>

<style scoped lang="scss">
.layout-logo {
  width: auto;
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: left;
  color: var(--el-color-primary);
  font-size: 16px;
  cursor: pointer;
  animation: logoAnimation 0.3s ease-in-out;
  span {
    white-space: nowrap;
    display: inline-block;
  }

  &:hover {
    span {
      color: var(--color-primary-light-2);
    }
  }

  &-medium-img {
    width: 30px;
    margin-right: 5px;
  }
}

.layout-logo-size {
  vertical-align: middle;
  display: table-cell;
  text-align: center;
  margin-top:6px;
}

.logo{
  margin-right: 0.8em;
}

.logo_all{
  margin: 0.5em 0 0.4em 2.1em;
}

</style>