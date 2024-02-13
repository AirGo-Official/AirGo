declare interface CustomerService {
  created_at:string
  updated_at:string
  id:number
  user_id:number
  service_status:boolean
  service_start_at:string
  service_end_at:string
  is_renew:boolean
  renewal_amount:string
  goods_id:number
  subject:string
  des:string
  price:string
  goods_type:string
  duration:number
  total_bandwidth:number
  node_connector:number
  node_speed_limit:number
  traffic_reset_day:number
  sub_status: boolean
  sub_uuid:string
  used_up:number
  used_down:number
}
declare interface PushCustomerServiceRequest {
  customer_service_id:number
  to_user_name:string
}