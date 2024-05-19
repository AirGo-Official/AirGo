<template>
  <div class="container layout-padding">
    <el-card shadow="hover" class="layout-padding-auto">
      <div class="mb15">
        <el-input v-model="state.queryParams.field_params_list[0].condition_value" size="default"
                  style="max-width: 180px"></el-input>
        <el-button size="default" type="primary" class="ml10" @click="getTicketList()">
          <el-icon>
            <ele-Search/>
          </el-icon>
          {{$t('message.common.query')}}
        </el-button>
      </div>
      <el-table :data="ticketStoreData.ticketList.value.data" height="100%" stripe style="width: 100%;flex: 1;" @sort-change="sortChange">
        <el-table-column fixed type="index" :label="$t('message.adminTicket.Ticket.index')" width="60"></el-table-column>
        <el-table-column prop="id" :label="$t('message.adminTicket.Ticket.id')" show-overflow-tooltip width="60px" sortable="custom"></el-table-column>
        <el-table-column prop="user_id" :label="$t('message.adminTicket.Ticket.user_id')" show-overflow-tooltip width="80px" sortable="custom"></el-table-column>
        <el-table-column prop="title" :label="$t('message.adminTicket.Ticket.title')" show-overflow-tooltip width="200px" sortable="custom"></el-table-column>
        <el-table-column prop="details" :label="$t('message.adminTicket.Ticket.details')" show-overflow-tooltip width="200px" sortable="custom"></el-table-column>
        <el-table-column prop="status" :label="$t('message.adminTicket.Ticket.status')" show-overflow-tooltip width="200px" sortable="custom">
          <template #default="{row}">
            <el-button v-if="row.status === constantStore.TICKET_PROCESSING" size="small" type="success">{{$t('message.constant.TICKET_PROCESSING')}}</el-button>
            <el-button v-else size="small" type="info" >{{$t('message.constant.TICKET_CLOSED')}}</el-button>
          </template>
        </el-table-column>
        <el-table-column :label="$t('message.common.operate')" show-overflow-tooltip width="200px">
          <template #default="{row}">
            <el-button size="small" text type="primary"
                       @click="openChat(row)">{{$t('message.common.reply')}}
            </el-button>
            <el-button size="small" text type="primary" @click="closeTicket(row)">{{$t('message.common.close')}}</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        background
        :page-sizes="[10, 30, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="onHandleSizeChange" @current-change="onHandleCurrentChange"
        v-model:current-page="state.queryParams.pagination.page_num"
        v-model:page-size="state.queryParams.pagination.page_size"
        :total="ticketStoreData.ticketList.value.total"
      />
    </el-card>
  </div>

</template>

<script setup lang="ts">
import {onBeforeMount, onMounted, reactive} from "vue";
import {storeToRefs} from "pinia";
import {useThemeConfig} from "/@/stores/themeConfig";
import {DateStrToTime} from "/@/utils/formatTime";


import { useAdminTicketStore } from "/@/stores/admin_logic/ticketStore";
import { useRouter } from "vue-router";
import { useAdminUserStore } from "/@/stores/admin_logic/userStore";
import { Session } from "/@/utils/storage"
import { useConstantStore } from "/@/stores/constantStore";
import { getCurrentAddress } from "/@/utils/request";
import { useApiStore } from "/@/stores/apiStore";
const ticketStore = useAdminTicketStore()
const ticketStoreData = storeToRefs(ticketStore)
const themeConfig = useThemeConfig()
const userStore = useAdminUserStore()
const userStoreData = storeToRefs(userStore)
const router = useRouter();
const constantStore = useConstantStore()
const apiStore = useApiStore();



const state = reactive({
  isCheckedTicket: 0,
  showDataType:'',
  queryParams:{
    table_name: 'ticket',
    field_params_list: [
      {field: 'title', field_chinese_name: '', field_type: '', condition: 'like', condition_value: '', operator: '',}
    ] as FieldParams[],
    pagination: {
      page_num: 1,
      page_size: 30,
      order_by: 'id DESC',
    } as Pagination,
  } as QueryParams,
})

const openChat=(row:Ticket)=>{
  Session.set('ticketUserID',row.user_id)
  Session.set('ticketID',row.id)
  let url = getCurrentAddress()
  window.open(url+'/#/static/ticketToChat',"_blank")
}
//关闭工单
const closeTicket=(row:Ticket)=>{
  row.status = constantStore.TICKET_CLOSED
  ticketStore.updateTicket(row).then((res)=>{
    getTicketList()
  })
}

const getTicketList = () => {
  ticketStore.getTicketList(state.queryParams)
}
const onHandleSizeChange=(val: number)=>{
  state.queryParams.pagination.page_size = val;
  getTicketList()
}
const onHandleCurrentChange = (val: number) => {
 state.queryParams.pagination.page_num = val;
  getTicketList()
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
     state.queryParams.pagination.order_by = p + " ASC"
      break
    default:
     state.queryParams.pagination.order_by = p + " DESC"
      break
  }
  getTicketList()
}
onMounted(() => {
  getTicketList()
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