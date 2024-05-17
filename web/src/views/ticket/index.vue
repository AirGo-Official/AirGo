<template>
  <div class="container layout-padding">
    <el-card style="border-radius:10px;margin-left: 0.2em;margin-right: 0.2em;margin-bottom: 1em;">
      <h3><i class="ri-list-unordered" style="margin-right: 2vh;"></i>{{ $t("message.ticket.total_ticket") }} : {{ ticketStoreData.userTicketList.value.total }}</h3>
    </el-card>
    <el-card shadow="hover" class="layout-padding-auto" style="border-radius:10px;margin-left: 0.2em;margin-right: 0.2em;margin-bottom: 1em;">
      <div class="mb15">
        <el-button size="default" type="success" class="ml10" @click="openTicketDialog">
          <el-icon>
            <ele-FolderAdd/>
          </el-icon>
          {{$t('message.adminTicket.addTicket')}}
        </el-button>
      </div>
      
      <el-table :data="ticketStoreData.userTicketList.value.data" stripe @sort-change="sortChange" height="100%">
        <el-table-column type="index" :label="$t('message.adminTicket.Ticket.index')" width="55" fixed/>
        <el-table-column prop="title" :label="$t('message.adminTicket.Ticket.title')"  show-overflow-tooltip 
                         sortable="custom" ></el-table-column>  
        <el-table-column prop="status" :label="$t('message.adminTicket.Ticket.status')"   show-overflow-tooltip width="130" max-width="200" fit
                         sortable="custom">
          <template #default="scope">
            <el-button v-if="scope.row.status === constantStore.TICKET_PROCESSING" type="success">{{$t('message.constant.TICKET_PROCESSING')}}</el-button>
            <el-button v-else type="info">{{$t('message.constant.TICKET_CLOSED')}}</el-button>
          </template>
        </el-table-column>
        <el-table-column :label="$t('message.common.operate') "width="150"  >
          <template #default="scope">
            <el-button size="small" type="primary" :disabled="scope.row.status === constantStore.TICKET_CLOSED" @click="toChat(scope.row)" >{{$t('message.common.reply')}}</el-button>
            <el-button size="small" type="primary" :disabled="scope.row.status === constantStore.TICKET_CLOSED" @click="closeTicket(scope.row)">{{$t('message.common.close')}}</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('message.adminTicket.Ticket.created_at')" width="150" show-overflow-tooltip 
                         sortable="custom">
          <template #default="scope">
            {{ DateStrToTime(scope.row.created_at) }}
          </template>
        </el-table-column>


      </el-table>
      <el-pagination
          background
          class="mt15"
          layout="total, sizes, prev, pager, next, jumper"
          :page-sizes="[10, 30, 50]"
          v-model:current-page="state.queryParams.pagination.page_num"
          v-model:page-size="state.queryParams.pagination.page_size"
          :total="ticketStoreData.userTicketList.value.total"
          @size-change="onHandleSizeChange"
          @current-change="onHandleCurrentChange"
      >
      </el-pagination>
    </el-card>
    <el-dialog v-model="state.isShowTicketDialog" :title="$t('message.adminTicket.addTicket')" width="80%" destroy-on-close align-center>
      <el-form v-model="ticketStoreData.newTicketInfo.value" size="default" label-position="top">
        <el-form-item :label="$t('message.adminTicket.Ticket.title')">
          <el-input v-model="ticketStoreData.newTicketInfo.value.title"></el-input>
        </el-form-item>
        <el-form-item :label="$t('message.adminTicket.Ticket.details')">
          <el-input v-model="ticketStoreData.newTicketInfo.value.details" type="textarea" autosize></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
				<span class="dialog-footer">
					<el-button @click="closeTicketDialog" size="default">{{$t('message.common.button_cancel')}}</el-button>
					<el-button type="primary" @click="newTicketSubmit" size="default">{{$t('message.common.button_confirm')}}</el-button>
				</span>
      </template>
    </el-dialog>
    <ToChatDialog ref="ToChatDialogRef"></ToChatDialog>
  </div>
</template>

<script setup lang="ts">
import {useTicketStore} from "/@/stores/user_logic/ticketStore";
import {storeToRefs} from "pinia";
import { defineAsyncComponent, onMounted, reactive, ref } from "vue";
import {DateStrToTime} from "/@/utils/formatTime";
import { useConstantStore } from "/@/stores/constantStore";

const constantStore = useConstantStore()
const ticketStore = useTicketStore()
const ticketStoreData = storeToRefs(ticketStore)
const ToChatDialog = defineAsyncComponent(() => import('/@/views/ticket/toChatDialog.vue'))
const ToChatDialogRef = ref()

const state = reactive({
  isShowTicketDialog:false,
  queryParams:{
    table_name: 'ticket',
    field_params_list: [
      { field: 'id', field_chinese_name: '', field_type: '', condition: '<>', condition_value: '', operator: '', }
    ] as FieldParams[],
    pagination: { page_num: 1, page_size: 30, order_by: 'id DESC',
    } as Pagination,//分页参数
  },
})
//
const openTicketDialog = () => {
  state.isShowTicketDialog = true
}
const closeTicketDialog=()=>{
  state.isShowTicketDialog = false
}
//确认提交
const newTicketSubmit = () => {
  ticketStore.newTicket().then(()=>{
    getUserTicketList()
  })
  closeTicketDialog()
}
const getUserTicketList = () => {
  ticketStore.getUserTicketList(state.queryParams)
}
const sortChange = (column:any) => {
  //处理嵌套字段
  let p = (column.prop as string)
  if (p.indexOf('.') !== -1) {
    p = p.slice(p.indexOf('.')+1)
  }
  switch (column.order){
    case 'ascending':
      state.queryParams.pagination.order_by=p+" ASC"
      break
    default:
      state.queryParams.pagination.order_by=p+" DESC"
      break
  }
  getUserTicketList()
}
//
const closeTicket=(row:Ticket)=>{
  row.status = constantStore.TICKET_CLOSED
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
  state.queryParams.pagination.page_size = val;
  getUserTicketList()
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
  state.queryParams.pagination.page_num = val;
  getUserTicketList()
};
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