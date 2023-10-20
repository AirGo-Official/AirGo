<template>
  <div class="container layout-padding">
    <el-card shadow="hover" class="layout-padding-auto">
      <div class="mb15">
        <el-input v-model="state.params.search" placeholder="请输入订单号"
                  style="max-width: 180px" size="default"></el-input>
        <el-button size="default" type="primary" class="ml10" @click="onSearch(state.params)">
          <el-icon>
            <ele-Search/>
          </el-icon>
          查询
        </el-button>
        <el-button size="default" type="success" class="ml10" @click="onShowCollapse">
          <el-icon>
            <ele-Search/>
          </el-icon>
          高级查询
        </el-button>

        <el-collapse v-if="state.isShowCollapse" v-model="state.activeCollapseNames">
          <el-collapse-item name="1">
            <!--          report组件-->
            <ReportComponent ref="reportRef" @getReportData="getReportDataHandler"></ReportComponent>
          </el-collapse-item>
        </el-collapse>
      </div>
      <el-table :data="orderManageData.allOrders.order_list" stripe style="width: 100%;flex: 1;">
        <el-table-column type="index" label="序号" fixed width="60px"/>
        <el-table-column prop="out_trade_no" label="订单号" width="200"/>
        <el-table-column prop="id" label="订单ID" width="60px"/>
        <el-table-column prop="created_at" label="下单日期" width="150">
          <template #default="scope">
            <el-tag type="success">{{ DateStrtoTime(scope.row.created_at) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="user_name" label="用户" width="180"/>
        <el-table-column prop="goods_id" label="商品ID" show-overflow-tooltip width="60"/>
        <el-table-column prop="subject" label="商品标题" show-overflow-tooltip width="200"/>
        <el-table-column prop="total_amount" label="订单金额" show-overflow-tooltip width="80"/>
        <el-table-column prop="receipt_amount" label="实收金额" show-overflow-tooltip width="80"/>
        <el-table-column prop="trade_status" label="交易状态" show-overflow-tooltip>
          <template #default="scope">
            <el-tag type="success" v-if="scope.row.trade_status==='TRADE_SUCCESS'">支付成功</el-tag>
            <el-tag type="warning" v-else-if="scope.row.trade_status==='WAIT_BUYER_PAY'">等待买家付款</el-tag>
            <el-tag type="danger" v-else-if="scope.row.trade_status==='TRADE_CLOSED'">交易超时关闭</el-tag>
            <el-tag type="success" v-else-if="scope.row.trade_status==='TRADE_FINISHED'">交易结束</el-tag>
            <el-tag type="info" v-else-if="scope.row.trade_status==='Created'">订单已创建</el-tag>
            <el-tag type="success" v-else-if="scope.row.trade_status==='Completed'">订单已完成</el-tag>
            <el-tag type="danger" v-else>未知状态</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="scope">
            <el-button v-if="scope.row.trade_status === 'WAIT_BUYER_PAY' || scope.row.trade_status ==='Created'"
                       size="small" text type="primary"
                       @click="onCompleteOrder(scope.row)">完成
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination background
                     class="mt15"
                     layout="total, sizes, prev, pager, next, jumper"
                     :page-sizes="[10, 30, 50]"
                     v-model:current-page="state.params.page_num"
                     v-model:page-size="state.params.page_size"
                     :total="orderManageData.allOrders.total"
                     @size-change="onHandleSizeChange"
                     @current-change="onHandleCurrentChange">
      </el-pagination>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import {defineAsyncComponent, onMounted, reactive, ref} from "vue";
import {useOrderStore} from "/@/stores/orderStore";
import {DateStrtoTime} from "/@/utils/formatTime"
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
import {storeToRefs} from "pinia";

const orderStore = useOrderStore()
const {orderManageData} = storeToRefs(orderStore)
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
const ReportComponent = defineAsyncComponent(() => import('/@/components/report/index.vue'))
const reportRef = ref()

//定义参数
const state = reactive({
  params: {
    page_num: 1,
    page_size: 30,
    search: '',
    date: [],
  } as PaginationParams,
  activeCollapseNames: '1', //当前激活的折叠面板
  isShowCollapse: false,
  fieldConditionList: {},
})
//
const onSearch = (params?: object) => {
  orderStore.getAllOrder(params)
}
onMounted(() => {
  onSearch(state.params)
})
//完成未支付订单
const onCompleteOrder = (row: Order) => {
  orderStore.completedOrder(row)
}
// 分页改变
const onHandleSizeChange = (val: number) => {
  if (state.isShowCollapse) {
    getReportDataHandler(state.fieldConditionList)
  } else {
    state.params.page_size = val;
    orderStore.getAllOrder(state.params)
  }

};
// 分页改变
const onHandleCurrentChange = (val: number) => {
  if (state.isShowCollapse) {
    getReportDataHandler(state.fieldConditionList)
  } else {
    state.params.page_num = val;
    orderStore.getAllOrder(state.params)
  }
};
//开启高级查询折叠面板
const onShowCollapse = () => {
  state.isShowCollapse = !state.isShowCollapse
  //防止子组件渲染太慢，导致undefined问题
  setTimeout(() => {
    if (state.isShowCollapse) {
      reportRef.value.openReportComponent("orders")
    }
  }, 500)
}
//请求数据
const getReportDataHandler = (data: any) => {
  //拼接分页参数
  (data as any).pagination_params = state.params;
  state.fieldConditionList = data
  request(apiStoreData.api.value.report_reportSubmit, data).then((res) => {
    orderManageData.value.allOrders.order_list = res.data.table_data
    orderManageData.value.allOrders.total = res.data.total
  })
}

//时间范围
const shortcuts = [
  {
    text: '上周',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
      return [start, end]
    },
  },
  {
    text: '上月',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
      return [start, end]
    },
  },
  {
    text: '最近3个月',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 90)
      return [start, end]
    },
  },
]

</script>

<style scoped lang="scss">
.container {
  :deep(.el-card__body) {
    display: flex;
    flex-direction: column;
    flex: 1;
    overflow: auto;

    .el-table {
      flex: 1;
    }
  }
}
</style>