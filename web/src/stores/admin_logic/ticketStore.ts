import {defineStore, storeToRefs} from "pinia";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
const apiStore = useApiStore()
export const useAdminTicketStore = defineStore("adminTicketStore", {
  state: () => ({
    ticketList: {
      total: 0,
      data: [] as Ticket[],
    },
    currentTicket: {
      user_id: 0,
      title: '',
      details: '',
      status: 'TicketProcessing',
      ticket_message: [] as TicketMessage[],
    } as Ticket,
    newTicketMessage: {
      ticket_id: 0,
      is_admin: true,
      message: '',
    } as TicketMessage,
  }),
  actions: {
    async firstTicketByID(id:number){
      const res = await request(apiStore.adminApi.firstTicket, {id:id} as Ticket)
      this.currentTicket = res.data
    },
    async getTicketList(params: QueryParams) {
      const res = await request(apiStore.adminApi.getTicketList, params)
      this.ticketList = res.data
    },
    async deleteTicket(params: Ticket) {
      return request(apiStore.adminApi.deleteTicket, params)
    },
    async updateTicket(params: Ticket){
      return request(apiStore.adminApi.updateTicket, params)
    },
    async sendTicketMessage(){
      return request(apiStore.adminApi.sendTicketMessage, this.newTicketMessage)
    }
  }
})