<template>
  <div class="container layout-padding">
    <el-card shadow="hover" class="layout-padding-auto">
      <div>
        <el-radio-group v-model="state.tabValue" @change="getUserTrafficLog" style="margin-bottom: 30px">
          <el-radio-button label="today">今日流量</el-radio-button>
          <el-radio-button label="this_week">本周流量</el-radio-button>
          <el-radio-button label="this_month">本月流量</el-radio-button>
          <el-radio-button label="last_month">上月流量</el-radio-button>
        </el-radio-group>
      </div>
      <div style="flex: 1;overflow: auto">
        <el-table :data="trafficStoreData.userTrafficLog.value.data" stripe height="100%" style="width: 100%"
                  @sort-change="sortChange">
          <el-table-column type="index" label="序号" width="60" fixed/>
<!--          <el-table-column prop="user_id" label="用户ID" width="80" show-overflow-tooltip-->
<!--                           sortable="custom"></el-table-column>-->
          <el-table-column prop="user_name" label="用户" width="200" show-overflow-tooltip
                           sortable="custom">
            <template #default="{ row }">
              <el-tag type="info">{{ acturalname(row.user_name) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="u" label="上行GB" width="200" show-overflow-tooltip sortable="custom">
            <template #default="{ row }">
              <el-tag type="info">{{ (row.u / 1024 / 1024 / 1024).toFixed(2) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="d" label="下行GB" show-overflow-tooltip sortable="custom">
            <template #default="{ row }">
              <el-tag type="info">{{ (row.d / 1024 / 1024 / 1024).toFixed(2) }}</el-tag>
            </template>
          </el-table-column>
        </el-table>
      </div>

    </el-card>
  </div>
</template>

<script setup lang="ts">

import {onBeforeMount, onMounted, reactive} from "vue";
import {useAccessStore} from "/@/stores/accessStore";
import {storeToRefs} from "pinia";
import {useReportStore} from "/@/stores/reportStore";
import {GetDurationDay, GetDurationMonth, GetDurationWithNDay} from "/@/utils/formatTime";
import {useTrafficStore} from "/@/stores/trafficStore";

const accessStore = useAccessStore()
const accessStoreData = storeToRefs(accessStore)
const reportStore = useReportStore()
const reportStoreData = storeToRefs(reportStore)
const trafficStore = useTrafficStore()
const trafficStoreData = storeToRefs(trafficStore)


const state = reactive({
  isShowDialog: false,
  title: "",
  tabValue: 'today',
  duration: [''],
  todayDuration: [''],
  thisWeekDuration: [''],
  thisMonthDuration: [''],
  lastMonthDuration: [''],
})
//初始化时间间隔
const initDuration = () => {
  state.todayDuration = GetDurationDay()
  state.thisWeekDuration = GetDurationWithNDay(7)
  state.thisMonthDuration = GetDurationMonth(0)
  state.lastMonthDuration = GetDurationMonth(1)
}
//
//初始化查询参数
const defaultFieldParams = () => {
  switch (state.tabValue) {
    case "today":
      state.duration = state.todayDuration
      break
    case "this_week":
      state.duration = state.thisWeekDuration
      break
    case "this_month":
      state.duration = state.thisMonthDuration
      break
    case "last_month":
      state.duration = state.lastMonthDuration
      break
  }
  reportStoreData.reportParams.value.table_name = 'user_traffic_log'
  reportStoreData.reportParams.value.field_params_list = [
    {field: 'created_at', condition: '>', condition_value: state.duration[0], operator: '',} as FieldParams,
    {field: 'created_at', condition: '<', condition_value: state.duration[1], operator: 'AND',} as FieldParams,
  ]
  reportStoreData.reportParams.value.pagination.page_num = 1
  reportStoreData.reportParams.value.pagination.page_size = 100
}
//
const getUserTrafficLog = () => {
  defaultFieldParams()
  trafficStore.getAllUserTrafficLog(reportStoreData.reportParams.value)
}

// 分页改变
const onHandleSizeChange = (val: number) => {
  reportStoreData.reportParams.value.pagination.page_size = val;
  getUserTrafficLog()
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
  reportStoreData.reportParams.value.pagination.page_num = val;
  getUserTrafficLog()
};
//排序监听
const sortChange = (column: any) => {
  //处理嵌套字段
  let p = (column.prop as string)
  if (p.indexOf('.') !== -1) {
    p = p.slice(p.indexOf('.') + 1)
  }
  switch (column.order) {
    case 'ascending':
      reportStoreData.reportParams.value.pagination.order_by = p + " ASC"
      break
    default:
      reportStoreData.reportParams.value.pagination.order_by = p + " DESC"
      break
  }
  getUserTrafficLog()
}
//
const acturalname = (str: string) => {
    return str.substring(0, 4) + '****' + str.substring(8, str.length)
}

onBeforeMount(() => {
  initDuration()
});
onMounted(() => {
  getUserTrafficLog()
});
</script>
<style scoped lang="scss">
.container {
  :deep(.el-card__body) {
    display: flex;
    flex-direction: column;
    flex: 1;
    overflow: auto;
  }
}
</style>