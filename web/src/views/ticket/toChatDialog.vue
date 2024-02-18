<template>
  <div>
    <el-dialog v-model="state.isShowDialog" :title="ticketStoreData.currentTicket.value.title" width="769px">
      <div id="chat">
        <div class="chatBox">
          <div >
              <span>{{$t('message.adminTicket.Ticket.details')}}：{{ticketStoreData.currentTicket.value.details}}</span>
          </div>
          <div class="chatBox-middle">
            <div class="chatInfo" id="chatInfo">
              <div class="chatUser-box" v-for="(item,index) in ticketStoreData.currentTicket.value.ticket_message" :key="index"
                   :class="[!item.is_admin?'chatUser-box1':'chatUser-box']">
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
              <el-button type="primary" :disabled="ticketStoreData.currentTicket.value.status === 'TicketClosed'" @click="sendMessage">{{$t('message.common.send')}}</el-button>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import {reactive} from "vue";
import {storeToRefs} from "pinia";
import {useTicketStore} from "/@/stores/user_logic/ticketStore";
import {DateStrToTime} from "/@/utils/formatTime";

const ticketStore = useTicketStore()
const ticketStoreData = storeToRefs(ticketStore)


const state = reactive({
  isShowDialog: false,
})

// 打开弹窗
const openDialog = (row: Ticket) => {
  state.isShowDialog = true;
  ticketStoreData.currentTicket.value = row
  ticketStore.firstTicket()
}
// 关闭弹窗
const closeDialog = () => {
  state.isShowDialog = false;
};
//
const sendMessage=()=>{
  ticketStoreData.newTicketMessage.value.ticket_id=ticketStoreData.currentTicket.value.id
  ticketStore.sendTicketMessage().then((res)=>{
    ticketStoreData.newTicketMessage.value.message=''
    ticketStore.firstTicket()
  })
}

// 暴露变量
defineExpose({
  openDialog,
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