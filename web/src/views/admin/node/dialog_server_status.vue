<template>
  <el-dialog v-model="state.isShowDialog" width="90%" style="height: 90%" destroy-on-close align-center>
    <el-row :gutter="15">
      <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" v-for="(v,k) in serverStatusData.data" :key="k">
        <div class="home-card-item">
          <el-card>
            <el-row :gutter="10" justify="space-around" align="middle">
              <el-col :xs="12" :sm="12" :md="4" :lg="4" :xl="4" style="margin: auto">
                <el-text style="color: #F9A43D;font-size: 15px;font-weight: bolder">{{ v.name }}</el-text>
              </el-col>
              <el-col :xs="12" :sm="12" :md="4" :lg="4" :xl="4" style="margin: auto">
                <el-icon v-if="v.status" color="#90ee90" :size="20">
                  <SuccessFilled/>
                </el-icon>
                <el-icon v-else color="#ff4d00" :size="20">
                  <CircleCloseFilled/>
                </el-icon>
              </el-col>
              <el-col :xs="12" :sm="12" :md="4" :lg="4" :xl="4" style="margin-top: 10px;margin-bottom: 10px">
                rate：{{ v.traffic_rate }}
              </el-col>
              <el-col :xs="12" :sm="12" :md="4" :lg="4" :xl="4" style="margin-top: 10px;margin-bottom: 10px">
                online：{{ v.user_amount }}
              </el-col>
              <el-col :xs="12" :sm="12" :md="4" :lg="4" :xl="4" style="margin: auto">
                <el-icon color="#409EFC">
                  <Top/>
                </el-icon>
                <span v-if="v.u < 1024">{{ v.u.toFixed(2) }}B/s</span>
                <span v-else-if="1024 < v.u < 1048576">{{ (v.u/1024).toFixed(2) }}kB/s</span>
                <span v-else-if="1048576 < v.u < 1073741824">{{ (v.u/1024/1024).toFixed(2) }}MB/s</span>
                <span v-else>{{ v.u }}B/s</span>
              </el-col>
              <el-col :xs="12" :sm="12" :md="4" :lg="4" :xl="4" style="margin: auto">
                <el-icon color="#409EFC">
                  <Bottom/>
                </el-icon>
                <span v-if="v.d < 1024">{{ v.d.toFixed(2) }}B/s</span>
                <span v-else-if="1024 < v.d < 1048576">{{ (v.d/1024).toFixed(2) }}kB/s</span>
                <span v-else-if="1048576 < v.d < 1073741824">{{ (v.d/1024/1024).toFixed(2) }}MB/s</span>
                <span v-else>{{ v.d }}B/s</span>
              </el-col>
            </el-row>

            <el-row :gutter="10">
              <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8" style="margin-top: 5px;margin-bottom: 5px">
                <el-progress :text-inside="true" :stroke-width="20" :percentage="v.cpu.toFixed(0)" striped striped-flow :color="customColors">
                  <template #default="{ percentage }">
                    <span class="percentage-label">cpu：</span>
                    <span class="percentage-value">{{ percentage }}%</span>
                  </template>
                </el-progress>
              </el-col>
              <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8" style="margin-top: 5px;margin-bottom: 5px">
                <el-progress :text-inside="true" :stroke-width="20" :percentage="v.mem.toFixed(0)" striped striped-flow :color="customColors">
                  <template #default="{ percentage }">
                    <span class="percentage-label">memory：</span>
                    <span class="percentage-value">{{ percentage }}%</span>
                  </template>
                </el-progress>
              </el-col>
              <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8" style="margin-top: 5px;margin-bottom: 5px">
                <el-progress :text-inside="true" :stroke-width="20" :percentage="v.disk.toFixed(0)" striped striped-flow :color="customColors">
                  <template #default="{ percentage }">
                    <span class="percentage-label">disk：</span>
                    <span class="percentage-value">{{ percentage }}%</span>
                  </template>
                </el-progress>
              </el-col>
            </el-row>
          </el-card>
        </div>
      </el-col>
    </el-row>
  </el-dialog>

</template>

<script setup lang="ts">
import { onMounted, onUnmounted, reactive } from "vue";
import {useApiStore} from "/@/stores/apiStore";
import {storeToRefs} from "pinia";
import {useAdminNodeStore} from "/@/stores/admin_logic/nodeStore";
import {Local} from "/@/utils/storage";

const nodeStore = useAdminNodeStore()
const {serverStatusData} = storeToRefs(nodeStore)
// const token = Local.get('token')
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
const customColors = [
  { color: '#5cb87a', percentage: 1 },
  { color: '#19befa', percentage: 40 },
  { color: '#f1d85f', percentage: 60 },
  { color: '#ef8b41', percentage: 70 },
  { color: '#f56c6c', percentage: 80 },
  { color: '#000000', percentage: 100 },
]
const state = reactive({
  isShowDialog:false,
})

const openDialog = () => {
  state.isShowDialog = true
  loop()
}
const closeDialog = () => {
  state.isShowDialog = false
}


const loop=()=>{
  let i = 0;
  let timer = setInterval(() => {
    getNodeServerStatus(timer, i++);
  }, 3000);
}
const getNodeServerStatus=(timer: NodeJS.Timeout, i: number)=> {
  setTimeout(() => {
    if (i >= 100) {
      clearInterval(timer);
    }
    nodeStore.getNodeServerStatus()
  },0)
}


// 暴露变量
defineExpose({
  openDialog   // 打开弹窗
});

</script>

<style scoped>
.home-card-item {
  font-size: 16px;
  width: 100%;
  height: 100%;
  border-radius: 4px;
  transition: all ease 0.3s;
  padding: 10px;
  overflow: hidden;
  background: var(--el-color-white);
  color: var(--el-text-color-primary);
  border: 1px solid var(--next-border-color-light);
}

.el-card {
  background-image: url("../../../assets/bgc/3.png");
  background-repeat: no-repeat;
}
</style>