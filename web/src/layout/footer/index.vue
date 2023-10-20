<template>
  <div class="layout-footer pb15">
    <div class="layout-footer-warp">
      <div class="mt5">
        <el-link>Version old-version</el-link>
        <span> | </span>
        <el-link href="https://github.com/ppoonk/AirGo">© AirGo</el-link>
        <span> | </span>
        <el-link >By {{ themeConfig.globalTitle }}</el-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts" name="layoutFooter">

//system store
import {storeToRefs} from "pinia";
import {useThemeConfig} from "/@/stores/themeConfig";

const themeStore = useThemeConfig();
const {themeConfig} = storeToRefs(themeStore);


// 此处需有内容（注释也得），否则缓存将失败
</script>

<style scoped lang="scss">
.layout-footer {
  width: 100%;
  display: flex;

  &-warp {
    margin: auto;
    color: var(--el-text-color-secondary);
    text-align: center;
    animation: error-num 0.3s ease;
  }
}
</style>
