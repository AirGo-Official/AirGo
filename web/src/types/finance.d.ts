declare interface BalanceStatement{
  created_at:string
  updated_at:string
  id:number
  user_id:number
  title:string
  type:string
  amount:string
  final_amount:string
}
declare interface CommissionStatement{
  created_at:string
  updated_at:string
  id:number
  user_id:number
  order_user_id:number
  order_user_name:string
  order_id:number
  subject:string
  total_amount:string
  commission_rate:number
  commission:string
  is_withdrew:boolean
}