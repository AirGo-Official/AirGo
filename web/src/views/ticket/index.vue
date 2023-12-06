<template>
  <div class="container layout-padding">
    <el-card shadow="hover" class="layout-padding-auto">
      <div class="mb15">
        <el-button size="default" type="success" class="ml10" @click="newTicket">
          <el-icon>
            <ele-FolderAdd/>
          </el-icon>
          新增工单
        </el-button>
      </div>
      <el-table :data="ticketStoreData.userTicketList.value.data" stripe @sort-change="sortChange" height="100%">
        <el-table-column type="index" label="序号" width="60" fixed/>
        <el-table-column prop="title" label="标题" show-overflow-tooltip width="300"
                         sortable="custom"></el-table-column>
        <el-table-column prop="created_at" label="创建日期" show-overflow-tooltip width="200"
                         sortable="custom">
          <template #default="scope">
            {{ DateStrtoTime(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" show-overflow-tooltip width="150"
                         sortable="custom">
          <template #default="scope">
            <el-button v-if="scope.row.status === 'TicketProcessing'" type="success">进行中</el-button>
            <el-button v-else type="info">关闭</el-button>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button size="small" text type="primary" @click="toChat(scope.row)">查看</el-button>
            <el-button size="small" text type="primary" @click="closeTicket(scope.row)">关闭</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
          background
          class="mt15"
          layout="total, sizes, prev, pager, next, jumper"
          :page-sizes="[10, 30, 50]"
          v-model:current-page="reportStoreData.reportParams.value.pagination.page_num"
          v-model:page-size="reportStoreData.reportParams.value.pagination.page_size"
          :total="ticketStoreData.userTicketList.value.total"
          @size-change="onHandleSizeChange"
          @current-change="onHandleCurrentChange"
      >
      </el-pagination>
    </el-card>
    <ToChatDialog ref="ToChatDialogRef"></ToChatDialog>
    <NewTicketDialog ref="NewTicketDialogRef" @refresh="getUserTicketList"></NewTicketDialog>
  </div>
</template>

<script setup lang="ts">
import {useTicketStore} from "/@/stores/ticketStore";
import {storeToRefs} from "pinia";
import {defineAsyncComponent, onBeforeMount, onMounted, ref} from "vue";
import {useReportStore} from "/@/stores/reportStore";
import {DateStrtoTime} from "/@/utils/formatTime";

const ticketStore = useTicketStore()
const ticketStoreData = storeToRefs(ticketStore)
const reportStore = useReportStore()
const reportStoreData = storeToRefs(reportStore)
const ToChatDialog = defineAsyncComponent(() => import('/@/views/ticket/toChatDialog.vue'))
const ToChatDialogRef = ref()

const NewTicketDialog = defineAsyncComponent(() => import('/@/views/ticket/newTicket.vue'))
const NewTicketDialogRef = ref()

//
const newTicket = () => {
  NewTicketDialogRef.value.openDialog()
}
const getUserTicketList = () => {
  ticketStore.getUserTicketList(reportStoreData.reportParams.value)
}
const sortChange = (column:any) => {
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
  getUserTicketList()
}
//
const closeTicket=(row:Ticket)=>{
  row.status = 'TicketClosed'
  ticketStore.updateUserTicket(row).then(()=>{
    getUserTicketList()
  })
}
//
const toChat = (row: Ticket) => {
  ToChatDialogRef.value.openDialog(row)
}
// 分页改变
const onHandleSizeChange = (val: number) => {
  reportStoreData.reportParams.value.pagination.page_size = val;
  getUserTicketList()
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
  reportStoreData.reportParams.value.pagination.page_num = val;
  getUserTicketList()
};
//初始化查询参数
const defaultFieldParams = () => {
  reportStoreData.reportParams.value.table_name = 'ticket'
  reportStoreData.reportParams.value.field_params_list = [
    {
      field: 'id',
      field_chinese_name: '',
      field_type: '',
      condition: '<>',
      condition_value: '',
      operator: '',
    } as FieldParams]
  reportStoreData.reportParams.value.pagination.order_by = "id DESC"
}
onBeforeMount(() => {
  defaultFieldParams()
});
onMounted(() => {
  getUserTicketList()
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