import {defineStore, storeToRefs} from "pinia";
import {ElMessage} from "element-plus";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";

const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)


export const useAdminShopStore = defineStore("adminShopStore", {
  state: () => ({
    currentGoods: {
      cover_image:"https://telegraph-image.pages.dev/file/ee6166159a0e330f67bd6.png",
      subject: "新套餐",
      price: "0.00",
      good_order: 0,
      is_show: true,
      is_sale:true,
      is_renew:true,
      des: `<h3 style="color:#00BFFF">新套餐</h3>
<h3 style="color:#DDA0DD">新套餐</h3>`,
      //
      price_3_month:'0.00',
      price_6_month:'0.00',
      price_12_month:'0.00',
      price_unlimited_duration:'0.00',
      enable_traffic_reset:false,
      goods_type: 'subscribe',
      deliver_type: 'none',
      deliver_text: '',
      //
      total_bandwidth: 0,
      node_connector: 3,
      //
      recharge_amount: '0.00',
      // checked_nodes: [0], //套餐编辑时选中的节点
      // nodes: [],

    } as Goods,
    //商品列表
    goodsList: [] as Goods[],
    currentOrder: {
      user_id: 0,
      user_name: '',
      // user: any;

      order_type: "",
      out_trade_no: '',
      goods_id: 0,
      goods_type: '',  //类型
      deliver_type: '',//发货类型
      deliver_text: '',//发货内容
      subject: '',
      price: '',
      duration: 0,
      pay_id: 0,   //支付方式id
      pay_type: '', //支付方式，alipay,epay
      coupon_id: 0,    //
      coupon_name: '',
      coupon_amount: '',
      balance_amount: '',

      // pay_info: {} as PreCreatePayToFrontend, //支付信息，epay，alipay等
      trade_no: '',
      buyer_logon_id: '',
      trade_status: '',
      original_amount: "0.00",
      total_amount: '0.00',
      buyer_pay_amount: '0.00',
    } as Order,
    orderList: {
      data: [] as Order[],
      total: 0,
    },
    //当前编辑的pay
    payInfo: {
      id: 0,
      name: '',
      pay_type: '',
      pay_logo_url: '',
      status: true,
      alipay: {
        alipay_app_id: '',
        alipay_notify_url: '',
        alipay_app_private_key: '',
        alipay_ali_public_key: '',
        alipay_encrypt_key: '',
      } as Alipay,
      epay: {
        epay_pid: 0,
        epay_key: '',
        epay_api_url: '',
        epay_return_url: '',
        epay_notify_url: '',
        epay_type: '',
      } as Epay,
    } as PayInfo,
    payList: [] as PayInfo[],
    //折扣列表
    couponList:{
      total:0,
      data:[] as Coupon[],
    },
    //当前编辑的折扣
    currentCoupon: {
      name: '',
      discount_rate: 0,
      limit: 0,
      expired_at: '',
      min_amount:0.01,
    } as Coupon,
    checkedGoodsIDs:[] as number[],
    checkedNodeIDs:[] as number[],
    orderSummary:{
      lastMonth:[] as OrderSummary[],
      thisMonth:[] as OrderSummary[],
    }
  }),
  actions: {
    //获取全部商品
    async getGoodsList() {
      const res = await request(apiStoreData.adminApi.value.getGoodsList)
      this.goodsList = res.data
    },
    //添加商品
    async newGoods() {
      return request(apiStoreData.adminApi.value.newGoods,this.nodeIDsWhenSubmitHandler(this.currentGoods))
    },
    //修改商品
    async updateGoods() {
      return request(apiStoreData.adminApi.value.updateGoods,this.nodeIDsWhenSubmitHandler(this.currentGoods))
    },
    //删除商品
    async deleteGoods(data?: Goods) {
      return request(apiStoreData.adminApi.value.deleteGoods,data)
    },
    //获取订单
    async getOrderList(params?: object) {
      const res = await request(apiStoreData.adminApi.value.getOrderList, params)
      this.orderList = res.data
    },
    //更新用户订单
    async updateUserOrder(params:object){
      return request(apiStoreData.adminApi.value.updateOrder, params)
    },
    //获取支付列表
    async getPayList() {
      const res = await request(apiStoreData.adminApi.value.getPayList)
      this.payList = res.data
    },
    //清空payInfo
    wipePayInfo() {
      this.payInfo = {
        id: 0,
        name: '',
        pay_type: '',
        pay_logo_url: '',
        status: true,
        alipay: {
          alipay_app_id: '',
          alipay_notify_url: '',
          alipay_app_private_key: '',
          alipay_ali_public_key: '',
          alipay_encrypt_key: '',
        } as Alipay,
        epay: {
          epay_pid: 0,
          epay_key: '',
          epay_api_url: '',
          epay_return_url: '',
          epay_notify_url: '',
          epay_type: '',
        } as Epay,
      } as PayInfo
    },
    //新建折扣
    async newCoupon(){
      return request(apiStore.adminApi.newCoupon,this.goodsIDsHandlerWhenSubmit(this.currentCoupon))
    },
    //获取折扣列表
    async getCouponList(){
     const res = await request(apiStoreData.adminApi.value.getCouponList)
        this.couponList = res.data
    },
    //更新折扣
    async updateCoupon(){
      return request(apiStore.adminApi.updateCoupon,this.goodsIDsHandlerWhenSubmit(this.currentCoupon))
    },
    //删除折扣
    async deleteCoupon(params:Coupon){
      return request(apiStoreData.adminApi.value.deleteCoupon, params)
    },
    goodsIDsHandler(){
      this.checkedGoodsIDs = []
      this.currentCoupon.goods.forEach((value: Goods, index: number, array: Goods[])=>{
        this.checkedGoodsIDs.push(value.id)
      })
    },
    goodsIDsHandlerWhenSubmit(coupon:Coupon){
      coupon.goods = []
      this.checkedGoodsIDs.forEach((value: number, index: number, array: number[])=>{
        coupon.goods.push({id:value} as Goods)
      })
      return coupon
    },
    nodeIDsHandler(){
      this.checkedNodeIDs = []
      this.currentGoods.nodes.forEach((value: NodeInfo, index: number, array: NodeInfo[])=>{
        this.checkedNodeIDs.push(value.id)
      })
    },
    nodeIDsWhenSubmitHandler(goods:Goods){
      goods.nodes = []
      this.checkedNodeIDs.forEach((value: number, index: number, array: number[])=>{
        goods.nodes.push({id:value} as NodeInfo)
      })
      return goods
    },
    async getOrderSummary(params:QueryParams,m:number){
      let mm = new Date().getMonth()
      const res = await request(apiStore.adminApi.orderSummary,params)
      if (res.data === null){
        return
      }
      if (m === (mm+1)){ //本月=4，则mm=3;mm+1=当前月
        this.orderSummary.thisMonth = res.data
      } else {
        this.orderSummary.lastMonth = res.data
      }
    },
  },
})