<template>
  <div v-loading="state.isShowLoading">
    <div class="home-card-one mb15" style="margin-top: 15px;margin-left: 20px">
      <el-radio-group v-model="state.goods_type" size="large" border @change="getAllEnabledGoods">
        <el-radio-button :label="constantStore.GOODS_TYPE_SUBSCRIBE">订阅</el-radio-button>
        <el-radio-button :label="constantStore.GOODS_TYPE_GENERAL">普通商品</el-radio-button>
        <el-radio-button :label="constantStore.GOODS_TYPE_RECHARGE">充值卡</el-radio-button>
      </el-radio-group>
    </div>
    <el-row :gutter="15" class="home-card-one mb15">
      <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" v-for="(v, k) in shopStoreData.goodsList.value" :key="k">
        <div class="home-card-item">
          <el-card>
            <template #header>
              <div>
                <el-text class="card-header-left">{{ v.subject }}</el-text>
              </div>
            </template>
            <!--            订阅商品开始-->
            <div v-if="v.goods_type === 'subscribe'">
              <div class="card-text">
                <el-tag class="card-text-left">套餐流量</el-tag>
                <span class="card-text-right">{{ v.total_bandwidth }}GB</span>
              </div>
              <div class="card-text">
                <el-tag class="card-text-left" type="warning">节点最大连接数</el-tag>
                <span class="card-text-right">{{ v.node_connector }}</span>
              </div>
            </div>
            <!--            订阅商品结束-->
            <!--            充值商品开始-->
            <div v-if="v.goods_type === 'recharge'">
              <div class="card-text">
                <el-tag class="card-text-left">充值金额</el-tag>
                <span class="card-text-right">{{ v.recharge_amount }}</span>
              </div>
            </div>
            <!--            充值商品结束-->
            <div class="card-text">
              <el-tag class="card-text-left" type="warning">价格</el-tag>
              <span class="card-text-right">¥{{ v.price }} / 月</span>
            </div>
            <div v-html="v.des"></div>
            <div
              style="margin-top: 10px;margin-bottom: 10px;display: flex; align-items: center;justify-content: space-between;">
              <el-button size="large" @click="openPurchase(v)" color="#FC3D08">立即购买</el-button>
              <span>限购：{{ v.quota }} / 库存：{{ v.stock }}</span>
            </div>
          </el-card>
        </div>

      </el-col>
    </el-row>
    <Purchase ref="PurchaseRef"></Purchase>
  </div>

</template>

<script setup lang="ts">
import { defineAsyncComponent, onMounted, reactive, ref } from "vue";
import { storeToRefs } from "pinia";
import { useShopStore } from "/@/stores/user_logic/shopStore";
import { useConstantStore } from "/@/stores/constantStore";

const shopStore = useShopStore();
const shopStoreData = storeToRefs(shopStore);
const Purchase = defineAsyncComponent(() => import("/@/views/shop/purchase.vue"));
const PurchaseRef = ref();
const constantStore = useConstantStore()
const state = reactive({
  isShowLoading:false,
  goods_type: constantStore.GOODS_TYPE_SUBSCRIBE,
});
const getAllEnabledGoods = () => {
  shopStore.getAllEnabledGoods({ goods_type: state.goods_type });
};
//加载时获取全部已启用商品
onMounted(() => {
  getAllEnabledGoods();
});

const openPurchase = (goood: Goods) => {
  state.isShowLoading = true
  //初始化购买商品信息
  shopStoreData.currentGoods.value = goood;
  shopStoreData.currentOrder.value.duration = 1; // 默认订购时长
  shopStoreData.currentOrder.value.order_type = constantStore.ORDER_TYPE_NEW //订单类型：新购入
  shopStoreData.currentOrder.value.goods_id = shopStoreData.currentGoods.value.id //订购商品ID

  //获取计算价格后的订单信息
  shopStore.getOrderInfo(shopStoreData.currentOrder.value).then((res) => {
    state.isShowLoading = false
    shopStoreData.currentOrder.value = {} as Order
    shopStoreData.currentOrder.value = res.data
    PurchaseRef.value.openDialog();
  }).catch(()=>{
    state.isShowLoading = false
  })
};

</script>

<style scoped>
.home-card-item {
  width: 100%;
  height: 100%;
  border-radius: 4px;
  transition: all ease 0.3s;
  padding: 20px;
  overflow: hidden;
  /*background: var(--el-color-white);*/
  color: var(--el-text-color-primary);
  /*border: 1px solid var(--next-border-color-light);*/
}

.el-card {
  background-repeat: no-repeat;
  background-position: 100%, 100%;
}

.card-text {
  display: flex;
  justify-content: space-between;
  height: 35px
}

.card-text-left {
  margin-top: auto;
  margin-bottom: auto;
}

.card-text-right {
  margin-top: auto;
  margin-bottom: auto;
  font-size: 20px;
}

.card-header-left {
  font-size: 30px;
  color: #FC3D08;
}
</style>