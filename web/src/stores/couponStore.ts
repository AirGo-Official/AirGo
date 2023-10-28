import {defineStore} from "pinia";

export const useCouponStore = defineStore("couponStore", {
    state: () => ({
        couponList: [] as Coupon[],
        //当前编辑的折扣
        coupon: {
            name: '',
            discount_rate: 0,
            limit: 0,
            expired_at: '',
            checked_goods:[0]
        } as Coupon,
    }),
    actions: {}

})