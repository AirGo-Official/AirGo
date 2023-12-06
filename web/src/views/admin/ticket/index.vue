<template>
  <div>
    <el-row>
      <el-col :xs="24" :sm="24" :md="24" :lg="8" :xl="8" style="padding-right: 10px">
        <el-card style="height: 88vh;overflow: auto;">
          <el-button type="primary" style="margin-bottom: 10px" @click="getTicketList">
            <el-icon><Refresh /></el-icon>刷新列表</el-button>
          <!--          列表开始-->
          <el-row style="height: 50px;margin-bottom: 10px" :class="[v.id === state.isCheckedTicket?'isCheckedTicket':'noCheckedTicket']" v-for="(v, k) in ticketStoreData.ticketList.value.data"
                  :key="k">
            <el-col :span="18" @click="toChat(v)">
              <div ><el-button type="primary" size="small">工单ID: {{ v.id }}</el-button></div>
              <div >> {{ v.title }}</div>
            </el-col>
            <el-col :span="6">
              <div style="text-align: center;display: flex;flex-direction: column;justify-content: space-between; ">
                <div >
                  <el-button style="width:80px" type="success" v-if="v.status === 'TicketProcessing'">进行中</el-button>
                  <el-button style="width:80px" type="info" v-else>已关闭</el-button>
                </div>
                <div>
                  <el-button style="width:80px" :disabled="v.status === 'TicketClosed'" @click="closeTicket(v)">关闭工单</el-button>
                </div>
              </div>
            </el-col>
          </el-row>
          <!--            列表结束-->
          <el-pagination
              small
              background
              layout="prev, pager, next"
              :total="ticketStoreData.ticketList.value.total"
              v-model:page-size="reportStoreData.reportParams.value.pagination.page_size"
              v-model:current-page="reportStoreData.reportParams.value.pagination.page_num"
              @current-change="onHandleCurrentChange"
          />
        </el-card>
      </el-col>

      <!--      右侧开始-->
      <el-col :xs="24" :sm="24" :md="24" :lg="16" :xl="16">
        <!--        聊天组件开始-->
        <div>
          <el-card style="height: 44vh">
            <div id="chat">
              <div class="chatBox">
                <div class="chatBox-top">
                  <div class="chatBox-top-text">
                    <el-row>
                      <el-col :span="4">
                        <span>工单ID: {{ ticketStoreData.currentTicket.value.id }}</span>
                      </el-col>
                      <el-col :span="20">
                        <span>标题: {{ ticketStoreData.currentTicket.value.title }}</span>
                      </el-col>
                    </el-row>
                  </div>
                  <div class="chatBox-top-text"><span>详情: {{ ticketStoreData.currentTicket.value.details }}</span>
                  </div>
                </div>
                <div class="chatBox-middle">
                  <div class="chatInfo" id="chatInfo">
                    <div class="chatUser-box" v-for="(item,index) in ticketStoreData.currentTicket.value.ticket_message"
                         :key="index" :class="[item.is_admin?'chatUser-box1':'chatUser-box']">
                      <div class="chatUser-info">
                        <div class="chatUser-info-name"
                             :class="[item.is_admin?'chatUser-info-name1':'chatUser-info-name']"><span class="nowDate">{{ DateStrtoTime(item.created_at) }}</span>
                        </div>
                        <div class="chatUser-info-text"
                             :class="[item.is_admin?'chatUser-info-text1':'chatUser-info-text']">
                          <span>{{ item.message }}</span>
                        </div>
                      </div>
                    </div>

                  </div>
                </div>
                <div class="chatBox-infoDesk">
                  <div class="chatBox-textarea">
                    <el-input v-model="ticketStoreData.newTicketMessage.value.message" type="textarea"
                              placeholder="请输入"/>
                  </div>
                  <div class="chatBox-sendOut">
                    <el-button class="sendOut" type="primary" @click="sendMessage">发送</el-button>
                  </div>
                </div>
              </div>
            </div>
          </el-card>
        </div>
        <!--        聊天组件结束-->
        <!--        数据展示开始-->
        <div>
          <el-card style="height: 44vh">
            <div style="margin-bottom: 10px">
              <el-button type="primary" @click="state.showDataType = 'user'"><el-icon><Avatar /></el-icon></el-button>
              <el-button type="primary" @click="state.showDataType = 'order'"><el-icon><Histogram /></el-icon></el-button>
            </div>
            <div v-if="state.showDataType === 'user'">
              <div><el-button style="width: 80px">用户ID: </el-button>{{userStoreData.userManageData.value.dialog.user.id}}</div>
              <div><el-button style="width: 80px">账户: </el-button>{{userStoreData.userManageData.value.dialog.user.user_name}}</div>
              <div><el-button style="width: 80px">注册日期: </el-button>{{userStoreData.userManageData.value.dialog.user.created_at}}</div>
              <div><el-button style="width: 80px">账户状态: </el-button>{{userStoreData.userManageData.value.dialog.user.enable}}</div>
              <div><el-button style="width: 80px">订阅状态: </el-button>{{userStoreData.userManageData.value.dialog.user.subscribe_info.sub_status}}</div>
              <div><el-button style="width: 80px">订阅到期: </el-button>{{userStoreData.userManageData.value.dialog.user.subscribe_info.expired_at}}</div>
              <div><el-button style="width: 80px">订阅套餐: </el-button>{{userStoreData.userManageData.value.dialog.user.subscribe_info.goods_subject}}</div>
              <div><el-button style="width: 80px">订阅URL: </el-button>{{serverStoreData.publicServerConfig.value.backend_url + apiStoreData.staticApi.value.user_getSub.path +"?link=" + userStoreData.userManageData.value.dialog.user.subscribe_info.subscribe_url}}</div>
              <div><el-button style="width: 80px">总流量: </el-button>{{(userStoreData.userManageData.value.dialog.user.subscribe_info.t / 1024 / 1024 / 1024).toFixed(2)}}GB</div>
              <div><el-button style="width: 80px">上行流量: </el-button>{{(userStoreData.userManageData.value.dialog.user.subscribe_info.u / 1024 / 1024 / 1024).toFixed(2)}}GB</div>
              <div><el-button style="width: 80px">下行流量: </el-button>{{(userStoreData.userManageData.value.dialog.user.subscribe_info.d / 1024 / 1024 / 1024).toFixed(2)}}GB</div>
            </div>
            <div v-if="state.showDataType === 'order'">
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

            </div>
          </el-card>
        </div>
        <!--        数据展示结束-->
      </el-col>
      <!--      右侧结束-->
    </el-row>


  </div>

</template>

<script setup lang="ts">
import {onBeforeMount, onMounted, reactive} from "vue";
import {useTicketStore} from "/@/stores/ticketStore";
import {storeToRefs} from "pinia";
import {useReportStore} from "/@/stores/reportStore";
import {useThemeConfig} from "/@/stores/themeConfig";
import {DateStrtoTime} from "/@/utils/formatTime";
import {useUserStore} from "/@/stores/userStore";
import {useApiStore} from "/@/stores/apiStore";
import {request} from "/@/utils/request";

import {useServerStore} from "/@/stores/serverStore";
import {defineAsyncComponent, ref} from "vue";
import {useOrderStore} from "/@/stores/orderStore";
const serverStore = useServerStore()
const serverStoreData = storeToRefs(serverStore)
const ticketStore = useTicketStore()
const ticketStoreData = storeToRefs(ticketStore)
const reportStore = useReportStore()
const reportStoreData = storeToRefs(reportStore)
const themeConfig = useThemeConfig()
const themeConfigData = storeToRefs(themeConfig)
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
const userStore = useUserStore()
const userStoreData = storeToRefs(userStore)
const orderStore = useOrderStore()
const {orderManageData} = storeToRefs(orderStore)




const state = reactive({
  isCheckedTicket: 0,
  showDataType:'',
})
//打开聊天框
const toChat = (row: Ticket) => {
  state.isCheckedTicket = row.id
  ticketStoreData.currentTicket.value = row
  ticketStore.getTicketMessage()
  getUserInfo(row)
  getUserOrder(row)
}
//关闭工单
const closeTicket=(row:Ticket)=>{
  row.status = 'TicketClosed'
  ticketStore.updateTicket(row).then((res)=>{
    getTicketList()
  })
}
//获取用户信息
const getUserInfo=(row:Ticket)=>{
  //初始化用户查询参数
  reportStoreData.reportParams.value.table_name = 'user'
  reportStoreData.reportParams.value.field_params_list = [
    {
      field: 'id',
      field_chinese_name: '',
      field_type: '',
      condition: '=',
      condition_value: row.user_id.toString(),
      operator: '',
    } as FieldParams]

  request(apiStore.api.user_getUserList, reportStoreData.reportParams.value).then((res)=>{
    userStoreData.userManageData.value.dialog.user=res.data.data[0]
  })

}
//获取用户订单
const getUserOrder=(row:Ticket)=>{
  //初始化订单查询参数
  reportStoreData.reportParams.value.table_name = 'orders'
  reportStoreData.reportParams.value.field_params_list = [
    {field: 'user_id', field_chinese_name: '', field_type: '', condition: '=', condition_value: row.user_id.toString(), operator: '',} as FieldParams]
  reportStoreData.reportParams.value.pagination = {page_num: 1, page_size: 8, order_by: 'id DESC',}
  orderStore.getAllOrder(reportStoreData.reportParams.value)
}
//初始化工单查询参数
const defaultFieldParams = () => {
  reportStoreData.reportParams.value.table_name = 'ticket'
  reportStoreData.reportParams.value.field_params_list = [
    {field: 'id', field_chinese_name: '', field_type: '', condition: '<>', condition_value: '', operator: '',} as FieldParams]
  reportStoreData.reportParams.value.pagination = {page_num: 1, page_size: 5, order_by: 'id DESC',
  }
}
const getTicketList = () => {
  ticketStore.getTicketList(reportStoreData.reportParams.value)
}
const sendMessage = () => {
  ticketStoreData.newTicketMessage.value.ticket_id = ticketStoreData.currentTicket.value.id
  ticketStoreData.newTicketMessage.value.is_admin = true
  ticketStore.sendTicketMessage().then((res) => {
    ticketStoreData.newTicketMessage.value.message = ''
    ticketStore.getTicketMessage()
  })
}

// 分页改变
const onHandleCurrentChange = (val: number) => {
  reportStoreData.reportParams.value.pagination.page_num = val;
  getTicketList()
};


const onHandleCurrentChangeForOrder=(val: number)=>{
  reportStoreData.reportParams.value.pagination.page_num = val;
}
onBeforeMount(() => {
  defaultFieldParams()
});
onMounted(() => {
  getTicketList()
});

</script>

<style scoped lang="scss">

.noCheckedTicket {
  background-color: rgba(181, 195, 210, 0.38)
}

.isCheckedTicket {
  background-color: rgba(245, 177, 108, 0.48)
}

.chatBox {
  height: 100%;
  background-color: #fff;
  overflow: hidden;
  border-radius: 1px;
}

.chatBox-top {
  width: 100%;
  //height: 100px;
  /*display: flex;*/
  flex-wrap: nowrap;
  align-items: center;
  background-color: #2B3D63;
}

.chatBox-top-imgBox {
  margin-left: 10px;
}

.chatBox-top-text {
  margin-left: 10px;
  font-size: 1rem;
  color: #fff;
}

.chatBox-middle {
  width: 100%;
  height: 200px;
  background-color: #fff;
  border-bottom: 1px solid #2B3D63;
}

.chatBox-infoDesk {
  width: 100%;
}

.chatBox-textarea {
  width: 100%;
}

.chatBox-sendOut {
  margin-top: 0.625rem;
  width: 100%;
  height: 3.125rem;
  text-align: right;
}

.sendOut {
  padding: 0 10px;
  height: 20px;
  margin: 0.3125rem 1.25rem 0 0;
}

.chatInfo {
  width: 94%;
  height: 94%;
  margin: 1.25rem auto;
  overflow: auto;
}

.chatUser-box {
  width: 100%;
  margin-bottom: 6px;
  display: flex;
  flex-direction: row;
}


.chatUser-box-img {
  display: flex;
}

.chatUser-info {
  margin: 0 1.25rem;
}

.chatUser-info-name {
  font-size: 0.875rem;
  color: #888;
  display: flex;
  flex-direction: row;
}

.nowDate {
  margin: 0 0.625rem;
}

.chatUser-info-text {
  margin-top: 0.3125rem;
  max-width: 40rem;
  padding: 0.5rem;
  background-color: #E8E8E8;
  border-radius: 0.5rem;
  float: left;
  table-layout: fixed;
  word-break: break-all;
  overflow: hidden;
}

.chatUser-info-text span {
  font-size: 0.9375rem;
  line-height: 1.5625rem;
}

.chatUser-box1 {
  flex-direction: row-reverse;
}

.chatUser-info-name1 {
  flex-direction: row-reverse;
}

.chatUser-info-text1 {
  background-color: #9bc6f3;
  float: right;
}

</style>