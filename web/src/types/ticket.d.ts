declare interface Ticket {
    created_at: string
    updated_at: string
    id: number
    user_id: number
    title: string
    details: string
    status: string
    ticket_message: TicketMessage[]
}

declare interface TicketMessage {
    created_at: string
    updated_at: string
    id: number
    ticket_id: number
    is_admin: boolean
    message: string
}