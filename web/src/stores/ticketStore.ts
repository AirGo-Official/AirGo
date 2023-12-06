import {defineStore, storeToRefs} from "pinia";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
const apiStore = useApiStore()
export const useTicketStore = defineStore("ticketStore", {
    state: () => ({
        ticketList: {
            total: 0,
            data: [] as Ticket[],
        },
        userTicketList: {
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
        newTicketInfo: {
            user_id: 0,
            title: '',
            details: '',
            status: 'TicketProcessing',
            ticket_message: [] as TicketMessage[],
        } as Ticket,
        newTicketMessage: {
            ticket_id: 0,
            is_admin: false,
            message: '',
        } as TicketMessage,
    }),
    actions: {
        async getTicketList(params: object) {
            const res = await request(apiStore.api.ticket_getTicketList, params)
            this.ticketList = res.data
        },
        async getUserTicketList(params: object) {
            const res = await request(apiStore.api.ticket_getUserTicketList, params)
            this.userTicketList = res.data
        },
        async newTicket() {
            return await request(apiStore.api.ticket_newTicket, this.newTicketInfo)
        },
        async deleteTicket(params: object) {
            return await request(apiStore.api.ticket_deleteTicket, params)
        },
        async sendTicketMessage() {
            return await request(apiStore.api.ticket_sendTicketMessage, this.newTicketMessage)
        },
        async updateTicket(params: object){
            return await request(apiStore.api.ticket_updateTicket, params)
        },
        async updateUserTicket(params: object){
            return await request(apiStore.api.ticket_updateUserTicket, params)
        },
        async getTicketMessage(){
            const res = await request(apiStore.api.ticket_getTicketMessage, this.currentTicket)
            this.currentTicket.ticket_message=res.data
        },
    }
})