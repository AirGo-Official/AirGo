<template>
  <div>
    <div>
      <el-row :gutter="15">
        <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8">
          <div class="home-card-item">
            <el-card class="box-card">
              <el-row :gutter="10" justify="space-around" align="middle">
                <el-col :xs="12" :sm="12" :md="12" :lg="12" :xl="12">
                  今日订单：{{ state.statisticsData.dayOrderDataCurrent.total }}
                </el-col>
                <el-col :xs="12" :sm="12" :md="12" :lg="12" :xl="12">
                  今日收入：{{ state.statisticsData.dayOrderDataCurrent.total_amount }} ¥
                </el-col>
              </el-row>
            </el-card>
          </div>
        </el-col>
        <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8">
          <div class="home-card-item">
            <el-card class="box-card">
              <el-row :gutter="10" justify="space-around" align="middle">
                <el-col :xs="12" :sm="12" :md="12" :lg="12" :xl="12">
                  当月订单：{{ state.statisticsData.monthOrderDataCurrent.total }}
                </el-col>
                <el-col :xs="12" :sm="12" :md="12" :lg="12" :xl="12">
                  当月收入：{{ state.statisticsData.monthOrderDataCurrent.total_amount }} ¥
                </el-col>
              </el-row>
            </el-card>
          </div>
        </el-col>
        <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8">
          <div class="home-card-item">
            <el-card class="box-card">
              <el-row :gutter="10" justify="space-around" align="middle">
                <el-col :xs="12" :sm="12" :md="12" :lg="12" :xl="12">
                  上月订单：{{ state.statisticsData.monthOrderDataLast.total }}
                </el-col>
                <el-col :xs="12" :sm="12" :md="12" :lg="12" :xl="12">
                  上月收入：{{ state.statisticsData.monthOrderDataLast.total_amount }} ¥
                </el-col>
              </el-row>
            </el-card>
          </div>
        </el-col>
      </el-row>
    </div>
    <div class="home-card-item">
      <el-divider>
        当月数据
      </el-divider>
      <el-table :data="state.statisticsData.monthNodeCurrent.node_list" height="100%" style="width: 100%;flex: 1;"
                stripe fit show-summary :summary-method="getSummaries">
        <el-table-column fixed type="index" label="序号" width="60"/>
        <el-table-column prop="remarks" label="节点名称" show-overflow-tooltip width="200"></el-table-column>
        <el-table-column prop="address" label="节点地址" show-overflow-tooltip width="200"></el-table-column>
        <el-table-column prop="total_up" label="上行流量(GB)" show-overflow-tooltip width="200">
          <template #default="scope">
            <el-tag type="warning">{{ scope.row.total_up / 1024 / 1024 / 1024 }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="total_down" label="下行流量(GB)" show-overflow-tooltip>
          <template #default="scope">
            <el-tag type="warning">{{ scope.row.total_down / 1024 / 1024 / 1024 }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <div class="home-card-item">
      <el-divider>
        上月数据
      </el-divider>
      <el-table :data="state.statisticsData.monthNodeLast.node_list" height="100%" style="width: 100%;flex: 1;" stripe
                fit show-summary :summary-method="getSummaries">
        <el-table-column fixed type="index" label="序号" width="60"/>
        <el-table-column prop="remarks" label="节点名称" show-overflow-tooltip width="200"></el-table-column>
        <el-table-column prop="address" label="节点地址" show-overflow-tooltip width="200"></el-table-column>
        <el-table-column prop="total_up" label="上行流量(GB)" show-overflow-tooltip width="200">
          <template #default="scope">
            <el-tag type="warning">{{ scope.row.total_up / 1024 / 1024 / 1024 }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="total_down" label="下行流量(GB)" show-overflow-tooltip>
          <template #default="scope">
            <el-tag type="warning">{{ scope.row.total_down / 1024 / 1024 / 1024 }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script lang="ts" setup>

import {onMounted, reactive} from 'vue'
import {TableColumnCtx} from "element-plus";
import {formatDate} from "/@/utils/formatTime";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
import {storeToRefs} from "pinia";

const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
const state = reactive({
  statisticsData: {
    dayOrderDataCurrent: {
      total_amount: 0,
      total: 0,
    } as OrdersWithTotal,
    monthOrderDataCurrent: {
      total_amount: 0,
      total: 0,
    } as OrdersWithTotal,
    monthOrderDataLast: {
      total_amount: 0,
      total: 0,
    } as OrdersWithTotal,
    monthNodeCurrent: {
      total: 0,
      node_list: [] as NodeInfo[],
    },
    monthNodeLast: {
      total: 0,
      node_list: [] as NodeInfo[],
    },
  },
  monthOptions: {
    dateDay: {
      page_num: 1,
      page_size: 200,
      date: ['', ''],
    },
    dateCurrent: {
      page_num: 1,
      page_size: 200,
      date: ['', ''],
    },
    dateLast: {
      page_num: 1,
      page_size: 200,
      date: ['', ''],
    },
  },
})

//获取订单统计
function getMonthOrder(params?: object) {
  request(apiStoreData.api.value.order_getMonthOrderStatistics, state.monthOptions.dateDay).then((res) => {
    state.statisticsData.dayOrderDataCurrent = res.data
  })
  request(apiStoreData.api.value.order_getMonthOrderStatistics, state.monthOptions.dateCurrent).then((res) => {
    state.statisticsData.monthOrderDataCurrent = res.data
  })
  request(apiStoreData.api.value.order_getMonthOrderStatistics, state.monthOptions.dateLast).then((res) => {
    state.statisticsData.monthOrderDataLast = res.data
  })
}

//获取节点流量统计
function getMonthNodeInfo(params?: object) {
  request(apiStoreData.api.value.node_getTraffic, state.monthOptions.dateCurrent).then((res) => {
    state.statisticsData.monthNodeCurrent = res.data
    //ElMessage.success(res.msg)

  })
  request(apiStoreData.api.value.node_getTraffic, state.monthOptions.dateLast).then((res) => {
    state.statisticsData.monthNodeLast = res.data
    // ElMessage.success(res.msg)
  })
}

function initDate() {
  // 时间范围格式 "2023-05-09 11:56:02"
  let currentDate = new Date();
  let currentY = currentDate.getFullYear();
  let currentM = currentDate.getMonth() + 1;
  // let MonthDayNum = new Date(currentY,currentM,0).getDate();  //计算当月的天数

  //当月
  let startDate = new Date(currentY, currentM - 1, 1);
  let endDate = new Date(currentY, currentM, 0, 23, 59, 59); // new Date(2020,11,0);//表示2020/11/30这天
  let start = formatDate(startDate, "YYYY-mm-dd HH:MM:SS")
  let end = formatDate(endDate, "YYYY-mm-dd HH:MM:SS")
  //上月
  let lastStartDate = new Date(currentY, currentM - 2, 1);
  let lastEndDate = new Date(currentY, currentM - 1, 0, 23, 59, 59); // new Date(2020,11,0);//表示2020/11/30这天
  let lastStart = formatDate(lastStartDate, "YYYY-mm-dd HH:MM:SS")
  let lastEnd = formatDate(lastEndDate, "YYYY-mm-dd HH:MM:SS")
  //当天范围
  let currentDayStartDate = new Date();
  currentDayStartDate.setHours(0);
  currentDayStartDate.setMinutes(0);
  currentDayStartDate.setSeconds(0);
  currentDayStartDate.setMilliseconds(0);
  let currentDayEndDate = new Date(currentDayStartDate.getTime() + 3600 * 1000 * 24 - 1000)
  let currentDayStart = formatDate(currentDayStartDate, "YYYY-mm-dd HH:MM:SS")
  let currentDayEnd = formatDate(currentDayEndDate, "YYYY-mm-dd HH:MM:SS")
  state.monthOptions.dateCurrent.date = [start, end]
  state.monthOptions.dateLast.date = [lastStart, lastEnd]
  state.monthOptions.dateDay.date = [currentDayStart, currentDayEnd]
}


// 页面加载前
onMounted(() => {
  initDate()
  getMonthOrder()
  getMonthNodeInfo()

});

//表格合计
interface SummaryMethodProps<T = any> {
  columns: TableColumnCtx<T>[]
  data: T[]
}

//合计
const getSummaries = (param: SummaryMethodProps) => {
  const {columns, data} = param
  const sums: string[] = []
  columns.forEach((column, index) => {
    if (index === 0) {
      sums[index] = '合计'
      return
    }
    const values = data.map((item) => Number(item[column.property]))
    if (!values.every((value) => Number.isNaN(value))) {
      sums[index] = `${values.reduce((prev, curr) => {
        const value = Number(curr)
        if (!Number.isNaN(value)) {
          return prev + curr / (1024 * 1024 * 1024)
        } else {
          return prev / (1024 * 1024 * 1024)
        }
      }, 0)} GB`
    } else {
      sums[index] = 'N/A'
    }
  })
  return sums
}

</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.home-card-item {
  width: 100%;
  /*height: 100%;*/
  border-radius: 4px;
  transition: all ease 0.3s;
  padding: 20px;
  overflow: hidden;
  background: var(--el-color-white);
  color: var(--el-text-color-primary);
  border: 1px solid var(--next-border-color-light);
}

.el-card {
  background-image: url("../../assets/bgc/1.png");
}
</style>