declare interface Coupon {
    created_at: string;
    updated_at: string;
    id: number;
    name: string;
    discount_rate: number;
    limit: number;
    expired_at: string;
    checked_goods: number[];
}