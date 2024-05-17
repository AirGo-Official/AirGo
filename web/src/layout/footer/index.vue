<template>
  <div class="layout-footer pb15">
    <div class="layout-footer-warp">
      <div class="mt5">    
        <span>@風嶼Link「正在不断进化中...」| </span>
        <span> build:2024/5/3</span>
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
