<template>
<div>
  <el-row>
    <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
      <el-card >
        <el-descriptions
          :column="1"
          border
          :title="$t('message.adminUser.userInfo')"
        >
          <el-descriptions-item :label="$t('message.adminUser.SysUser.created_at')">{{ DateStrToTime(userStoreData.currentUser.value.created_at) }}</el-descriptions-item>
          <el-descriptions-item :label="$t('message.adminUser.SysUser.id')">{{ userStoreData.currentUser.value.id }}</el-descriptions-item>
          <el-descriptions-item :label="$t('message.adminUser.SysUser.user_name')">{{ userStoreData.currentUser.value.user_name }}</el-descriptions-item>
          <el-descriptions-item :label="$t('message.adminUser.SysUser.enable')">
            <el-button v-if="userStoreData.currentUser.value.enable" type="success">{{$t('message.common.enable')}}</el-button>
            <el-button v-else type="info">{{$t('message.common.disable')}}</el-button>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('message.adminUser.SysUser.balance')">{{ userStoreData.currentUser.value.balance }}</el-descriptions-item>
          <el-descriptions-item :label="$t('message.adminUser.SysUser.tg_id')">{{ userStoreData.currentUser.value.tg_id }}</el-descriptions-item>

        </el-descriptions>
      </el-card>
    </el-col>
    <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
      <el-card >
        <template #header>
          <div class="card-header">
            <span>{{$t('message.adminTicket.Ticket.title')}}：{{ticketStoreData.currentTicket.value.title}}</span>
          </div>
        </template>
        <div id="chat">
          <div class="chatBox">
            <div >
              <span>{{$t('message.adminTicket.Ticket.details')}}：{{ticketStoreData.currentTicket.value.details}}</span>
            </div>
            <div class="chatBox-middle">
              <div class="chatInfo" id="chatInfo">
                <div class="chatUser-box" v-for="(item,index) in ticketStoreData.currentTicket.value.ticket_message" :key="index"
                     :class="[item.is_admin?'chatUser-box1':'chatUser-box']">
                  <div class="chatUser-info">
                    <div class="chatUser-info-name" :class="[!item.is_admin?'chatUser-info-name1':'chatUser-info-name']">
                      <span class="nowDate">{{ DateStrToTime(item.created_at) }}</span>
                    </div>
                    <div class="chatUser-info-text" :class="[!item.is_admin?'chatUser-info-text1':'chatUser-info-text']">
                      <span>{{ item.message }}</span>
                    </div>
                  </div>
                </div>

              </div>
            </div>
            <div class="chatBox-infoDesk">
              <div class="chatBox-textarea">
                <el-input v-model="ticketStoreData.newTicketMessage.value.message" type="textarea"/>
              </div>
              <div class="chatBox-sendOut">
                <el-button type="primary" :disabled="ticketStoreData.currentTicket.value.status === constantStore.TICKET_CLOSED" @click="sendMessage">{{$t('message.common.send')}}</el-button>
              </div>
            </div>
          </div>
        </div>
      </el-card>
    </el-col>
  </el-row>
</div>
</template>

<script setup lang="ts">
import { useAdminUserStore } from "/@/stores/admin_logic/userStore";
import { storeToRefs } from "pinia";
import {DateStrToTime} from "/@/utils/formatTime";
import { onMounted } from "vue";
import { Session } from "/@/utils/storage";
import { useAdminTicketStore } from "/@/stores/admin_logic/ticketStore";
import { useConstantStore } from "/@/stores/constantStore";

const userStore = useAdminUserStore()
const userStoreData = storeToRefs(userStore)
const ticketStore = useAdminTicketStore()
const ticketStoreData = storeToRefs(ticketStore)
const constantStore = useConstantStore()

const firstUser=()=>{
 const userID:number = Session.get('ticketUserID')
 const ticketID:number = Session.get('ticketID')
  userStore.firstUserByID(userID)
  ticketStore.firstTicketByID(ticketID)
}

const sendMessage=()=>{
  ticketStoreData.newTicketMessage.value.ticket_id=ticketStoreData.currentTicket.value.id
  ticketStore.sendTicketMessage().then((res)=>{
    ticketStoreData.newTicketMessage.value.message=''
    ticketStore.firstTicketByID(ticketStoreData.currentTicket.value.id)
  })
}
onMounted(()=>{
  firstUser()
});

</script>

<style scoped lang="scss">

.chatBox {
  height: 100%;
  background-color: #fff;
  overflow: hidden;
  border-radius: 1px;
}

.chatBox-top {
  width: 100%;
  height: 100px;
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
  height: 400px;
  background-color: #fff;
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
  max-width: 20rem;
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