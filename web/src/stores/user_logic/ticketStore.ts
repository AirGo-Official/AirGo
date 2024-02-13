import {defineStore, storeToRefs} from "pinia";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
const apiStore = useApiStore()
export const useTicketStore = defineStore("ticketStore", {
    state: () => ({
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
        async newTicket() {
            return  request(apiStore.userApi.newTicket, this.newTicketInfo)
        },
        async updateUserTicket(params: object){
            return  request(apiStore.userApi.updateUserTicket, params)
        },
        async getUserTicketList(params: object) {
            const res = await request(apiStore.userApi.getUserTicketList, params)
            this.userTicketList = res.data
        },
        async firstTicket(){
            const res = await request(apiStore.userApi.firstTicket, this.currentTicket)
            this.currentTicket = res.data
        },
        async sendTicketMessage() {
            return  request(apiStore.userApi.sendTicketMessage, this.newTicketMessage)
        },


    }
})