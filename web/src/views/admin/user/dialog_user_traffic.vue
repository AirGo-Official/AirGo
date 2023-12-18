<template>
  <el-dialog v-model="state.isShowDialog" :title="state.title" width="90%" style="height: 90%" destroy-on-close
             align-center  class="layout-padding-auto">

      <div class="mb15">
        <el-row>
          <el-col :xs="24" :sm="24" :md="24" :lg="8" :xl="8">
            <el-date-picker
                size="default"
                v-model="state.duration"
                type="datetimerange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                value-format="YYYY-MM-DD HH:mm:ss"
            />
          </el-col>
          <el-col :xs="24" :sm="24" :md="24" :lg="8" :xl="8">
            <el-input v-model="reportStoreData.reportParams.value.field_params_list[2].condition_value" size="default"
                      placeholder="请输入用户名称"
                      style="max-width: 180px"></el-input>
            <el-button @click="getUserTrafficLog" size="default" type="primary" class="ml10">
              <el-icon>
                <ele-Search/>
              </el-icon>
              查询
            </el-button>
          </el-col>
        </el-row>
      </div>

        <el-table :data="trafficStoreData.userTrafficLog.value.data" stripe height="550px" style="width: 100%;flex: 1;overflow: auto"
                  @sort-change="sortChange">
          <el-table-column type="index" label="序号" width="60" fixed/>
          <el-table-column prop="user_id" label="用户ID" width="80" show-overflow-tooltip
                           sortable="custom"></el-table-column>
          <el-table-column prop="user_name" label="用户" width="200" show-overflow-tooltip
                           sortable="custom"></el-table-column>
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
        <el-pagination
            background
            layout="total, sizes, prev, pager, next, jumper"
            :page-sizes="[10, 30, 50]"
            v-model:current-page="reportStoreData.reportParams.value.pagination.page_num"
            v-model:page-size="reportStoreData.reportParams.value.pagination.page_size"
            :total="trafficStoreData.userTrafficLog.value.total"
            @size-change="onHandleSizeChange"
            @current-change="onHandleCurrentChange"
        >
        </el-pagination>
  </el-dialog>
</template>

<script setup lang="ts">

import {reactive} from "vue";
import {useAccessStore} from "/@/stores/accessStore";
import {storeToRefs} from "pinia";
import {useReportStore} from "/@/stores/reportStore";
import {GetDurationMonth} from "/@/utils/formatTime";
import {useTrafficStore} from "/@/stores/trafficStore";

const accessStore = useAccessStore()
const accessStoreData = storeToRefs(accessStore)
const reportStore = useReportStore()
const reportStoreData = storeToRefs(reportStore)
const trafficStore = useTrafficStore()
const trafficStoreData = storeToRefs(trafficStore)
const emit = defineEmits(['defaultParams'])


const state = reactive({
  isShowDialog: false,
  title: "",
  duration: [''],
  search: '',
})
//
//初始化查询参数
const defaultFieldParams = () => {
  state.duration = GetDurationMonth(0)
  reportStoreData.reportParams.value.table_name = 'user_traffic_log'
  reportStoreData.reportParams.value.field_params_list = [
    {field: 'created_at', condition: '>', condition_value: state.duration[0], operator: '',} as FieldParams,
    {field: 'created_at', condition: '<', condition_value: state.duration[1], operator: 'AND',} as FieldParams,
    {field: 'user_name', condition: 'like', condition_value: state.search, operator: 'AND',} as FieldParams,
  ]
  reportStoreData.reportParams.value.pagination = {page_num: 1, page_size: 30, order_by: 'user_id',} as Pagination
}
//
const getUserTrafficLog = () => {
  reportStoreData.reportParams.value.field_params_list[0].condition_value=state.duration[0]
  reportStoreData.reportParams.value.field_params_list[1].condition_value=state.duration[1]
  trafficStore.getAllUserTrafficLog(reportStoreData.reportParams.value)
}

// 打开弹窗
const openDialog = () => {
  state.isShowDialog = true
  defaultFieldParams()
  getUserTrafficLog()
}
// 关闭弹窗
const closeDialog = () => {
  state.isShowDialog = false
  emit('defaultParams')

};

//确认提交
function onSubmit() {

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

// 暴露变量
defineExpose({
  openDialog,   // 打开弹窗
});

</script>

<style scoped lang="scss">

</style>