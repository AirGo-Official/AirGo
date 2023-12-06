<template>
  <div>
    <el-dialog v-model="state.isShowDialog" width="769px">
      <div id="chat">
        <div class="chatBox">
          <div class="chatBox-top">
            <div style="font-size: 20px;color: white;">
              <span style="margin-right: 20px">标题: {{ticketStoreData.currentTicket.value.title}}</span>
            </div>
            <div style="font-size: 16px;color: #c2c0c0;">
              <span>详情：{{ticketStoreData.currentTicket.value.details}}</span>
            </div>
          </div>
          <div class="chatBox-middle">
            <div class="chatInfo" id="chatInfo">
              <div class="chatUser-box" v-for="(item,index) in ticketStoreData.currentTicket.value.ticket_message" :key="index"
                   :class="[!item.is_admin?'chatUser-box1':'chatUser-box']">
                <div class="chatUser-info">
                  <div class="chatUser-info-name" :class="[!item.is_admin?'chatUser-info-name1':'chatUser-info-name']">
                    <span class="nowDate">{{ DateStrtoTime(item.created_at) }}</span>
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
              <el-input v-model="ticketStoreData.newTicketMessage.value.message" type="textarea" placeholder="请输入咨询信息"/>
            </div>
            <div class="chatBox-sendOut">
              <el-button type="primary" :disabled="ticketStoreData.currentTicket.value.status === 'TicketClosed'" @click="sendMessage">发送</el-button>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import {reactive} from "vue";
import {useThemeConfig} from "/@/stores/themeConfig";
import {storeToRefs} from "pinia";
import userlogo from "/@/assets/icon/userlogo.jpeg"
import {useTicketStore} from "/@/stores/ticketStore";
import {DateStrtoTime} from "/@/utils/formatTime";

const themeConfig = useThemeConfig()
const themeConfigData = storeToRefs(themeConfig)
const ticketStore = useTicketStore()
const ticketStoreData = storeToRefs(ticketStore)


const state = reactive({
  isShowDialog: false,
})

// 打开弹窗
const openDialog = (row: Ticket) => {
  state.isShowDialog = true;
  ticketStoreData.currentTicket.value = row
  ticketStore.getTicketMessage()
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
    ticketStore.getTicketMessage()
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