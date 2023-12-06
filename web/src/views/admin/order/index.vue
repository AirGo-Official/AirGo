<template>
  <div class="container layout-padding">
    <el-card shadow="hover" class="layout-padding-auto">
      <div class="mb15">
        <el-input v-model="reportStoreData.reportParams.value.field_params_list[0].condition_value" placeholder="请输入订单号"
                  style="max-width: 180px" size="default"></el-input>
        <el-button size="default" type="primary" class="ml10" @click="onSearch()">
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
            <ReportComponent ref="reportRef" @getReportData="getReportData"></ReportComponent>
          </el-collapse-item>
        </el-collapse>
      </div>
      <el-table :data="orderManageData.allOrders.order_list" stripe style="width: 100%;flex: 1;" @sort-change="sortChange">
        <el-table-column type="index" label="序号" fixed width="60px"/>
        <el-table-column prop="out_trade_no" label="订单号" width="200" sortable="custom"/>
        <el-table-column prop="id" label="订单ID" width="100px" sortable="custom"/>
        <el-table-column prop="created_at" label="下单日期" width="150" sortable="custom">
          <template #default="scope">
            <el-tag type="success">{{ DateStrtoTime(scope.row.created_at) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="user_name" label="用户" width="180" sortable="custom"/>
        <el-table-column prop="goods_id" label="商品ID" show-overflow-tooltip width="100" sortable="custom"/>
        <el-table-column prop="subject" label="商品标题" show-overflow-tooltip width="200" sortable="custom"/>
        <el-table-column prop="total_amount" label="订单金额" show-overflow-tooltip width="100" sortable="custom"/>
        <el-table-column prop="receipt_amount" label="实收金额" show-overflow-tooltip width="100" sortable="custom"/>
        <el-table-column prop="trade_status" label="交易状态" show-overflow-tooltip sortable="custom" width="100">
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
                     v-model:current-page="reportStoreData.reportParams.value.pagination.page_num"
                     v-model:page-size="reportStoreData.reportParams.value.pagination.page_size"
                     :total="orderManageData.allOrders.total"
                     @size-change="onHandleSizeChange"
                     @current-change="onHandleCurrentChange">
      </el-pagination>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import {defineAsyncComponent, onBeforeMount, onMounted, reactive, ref} from "vue";
import {useOrderStore} from "/@/stores/orderStore";
import {DateStrtoTime} from "/@/utils/formatTime"
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
import {storeToRefs} from "pinia";
import {ElMessage, ElMessageBox} from "element-plus";
import {useReportStore} from "/@/stores/reportStore";

const orderStore = useOrderStore()
const {orderManageData} = storeToRefs(orderStore)
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
const ReportComponent = defineAsyncComponent(() => import('/@/components/report/index.vue'))
const reportRef = ref()
const reportStore = useReportStore()
const reportStoreData = storeToRefs(reportStore)

//定义参数
const state = reactive({
  activeCollapseNames: '1', //当前激活的折叠面板
  isShowCollapse: false,
})
//
const onSearch = () => {
  orderStore.getAllOrder(reportStoreData.reportParams.value)
}
//完成未支付订单
const onCompleteOrder=(row: Order)=> {
  ElMessageBox.confirm(`此操作将永久完成用户未支付订单（${row.subject}），并使之有效，, 是否继续?`, '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        orderStore.completedOrder(row)
        setTimeout(()=>{
          onSearch()
        },500)
      })
      .catch(() => {
      });
}

// 分页改变
const onHandleSizeChange = (val: number) => {
  reportStoreData.reportParams.value.pagination.page_size = val;
    onSearch()
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
  reportStoreData.reportParams.value.pagination.page_num = val;
    onSearch()
};
//排序监听
const sortChange = (column: any) => {
  //处理嵌套字段
  let p = (column.prop as string)
  if (p.indexOf('.') !== -1) {
    p = p.slice(p.indexOf('.')+1)
  }
  switch (column.order){
    case 'ascending':
      reportStoreData.reportParams.value.pagination.order_by=p+" ASC"
      break
    default:
      reportStoreData.reportParams.value.pagination.order_by=p+" DESC"
      break
  }
  onSearch()
}
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
//
const getReportData = (data: any) => {
  onSearch()
}
//初始化查询参数
const defaultFieldParams = () => {
  reportStoreData.reportParams.value.table_name = 'orders'
  reportStoreData.reportParams.value.field_params_list = [
    {field: 'out_trade_no', field_chinese_name: '', field_type: '', condition: 'like', condition_value: '', operator: '',} as FieldParams]
  reportStoreData.reportParams.value.pagination = {page_num: 1, page_size: 30, order_by: 'id DESC',} as Pagination
}
//
onBeforeMount(() => {
  defaultFieldParams()
});
onMounted(() => {
  onSearch()
});

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