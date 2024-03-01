<template>
  <div class="lazy-img-container layout-pd">
    <div class="home-card-one mb15">
      <el-radio-group v-model="state.goods_type" size="large" border @change="getAllEnabledGoods">
        <el-radio-button :label="constantStore.GOODS_TYPE_SUBSCRIBE">{{$t('message.constant.GOODS_TYPE_SUBSCRIBE')}}</el-radio-button>
        <el-radio-button :label="constantStore.GOODS_TYPE_GENERAL">{{$t('message.constant.GOODS_TYPE_GENERAL')}}</el-radio-button>
        <el-radio-button :label="constantStore.GOODS_TYPE_RECHARGE">{{$t('message.constant.GOODS_TYPE_RECHARGE')}}</el-radio-button>
      </el-radio-group>
    </div>
    <el-card shadow="hover" >
      <div class="flex-warp" v-if="shopStoreData.goodsList.value.length > 0">
        <el-row :gutter="15">
          <el-col :xs="24" :sm="12" :md="8" :lg="6" :xl="4" class="mb15"
                  v-for="(v, k) in shopStoreData.goodsList.value"
                  :key="k" @click="showGoodsDetails(v)">
            <div class="flex-warp-item">
              <div class="flex-warp-item-box">
                <div class="item-img">
                  <el-image :src="v.cover_image" lazy>
                    <template #error>
                      <div class="image-slot">
                        <el-icon><icon-picture /></el-icon>
                      </div>
                    </template>
                  </el-image>
                </div>
                <div class="item-txt">
                  <div class="item-txt-title">{{ v.subject }}</div>
                  <div class="item-txt-other">
                    <div style="width: 100%">
                      <div class="item-txt-msg mb10">
                        <span>{{$t('message.adminShop.Goods.quota')}} {{ v.quota }}</span>
                        <span class="ml10">{{$t('message.adminShop.Goods.stock')}} {{ v.stock }}</span>
                      </div>
                      <div class="item-txt-msg item-txt-price">
												<span class="font-price">
													<span>￥</span>
													<span class="font">{{ v.price }}</span>
												</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
      <el-empty v-else :description="$t('message.common.noData')"></el-empty>
    </el-card>
    <el-dialog v-model="state.isShowGoodsDetails" width="80%"
               :title="$t('message.common.details')"
               destroy-on-close>
      <el-row :gutter="50">
        <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
          <div style="border-radius:10px;background: rgba(224,224,224,0.29);padding: 10px">
            <div style="margin-top: 10px;text-align: center">
              <el-image :src="shopStoreData.currentGoods.value.cover_image"
                        lazy
                        style="height: 150px"
                        fit="cover"
                        :preview-src-list="[shopStoreData.currentGoods.value.cover_image]">
                <template #error>
                  <div class="image-slot">
                    <el-icon><icon-picture/></el-icon>
                  </div>
                </template>
              </el-image>
            </div>
            <el-divider></el-divider>
            <div style="margin-top: 10px;">
              {{shopStoreData.currentGoods.value.subject}}
            </div>
          </div>

          <div style="margin-top: 10px;border-radius:10px;background: rgba(224,224,224,0.29);padding: 10px">
            <div style="margin-top: 10px;margin-bottom: 10px">
              <el-tag v-if="shopStoreData.currentGoods.value.goods_type === constantStore.GOODS_TYPE_SUBSCRIBE">{{$t('message.adminShop.Goods.goods_type')}}: {{$t('message.constant.GOODS_TYPE_SUBSCRIBE')}}</el-tag>
              <el-tag v-if="shopStoreData.currentGoods.value.goods_type === constantStore.GOODS_TYPE_RECHARGE" >{{$t('message.adminShop.Goods.goods_type')}}: {{$t('message.constant.GOODS_TYPE_RECHARGE')}}</el-tag>
              <el-tag v-if="shopStoreData.currentGoods.value.goods_type === constantStore.GOODS_TYPE_GENERAL">{{$t('message.adminShop.Goods.goods_type')}}: {{$t('message.constant.GOODS_TYPE_GENERAL')}}</el-tag>

              <el-tag v-if="shopStoreData.currentGoods.value.deliver_type === constantStore.DELIVER_TYPE_NONE">{{$t('message.adminShop.Goods.deliver_type')}}: {{$t('message.constant.DELIVER_TYPE_NONE')}}</el-tag>
              <el-tag v-if="shopStoreData.currentGoods.value.deliver_type === constantStore.DELIVER_TYPE_AUTO">{{$t('message.adminShop.Goods.deliver_type')}}: {{$t('message.constant.DELIVER_TYPE_AUTO')}}</el-tag>
              <el-tag v-if="shopStoreData.currentGoods.value.deliver_type === constantStore.DELIVER_TYPE_MANUAL">{{$t('message.adminShop.Goods.deliver_type')}}: {{$t('message.constant.DELIVER_TYPE_MANUAL')}}</el-tag>
            </div>
            <div style="margin-top: 10px;margin-bottom: 10px">
              <el-tag type="warning">{{$t('message.adminShop.Goods.quota')}}：{{ shopStoreData.currentGoods.value.quota }} / {{$t('message.adminShop.Goods.stock')}}：{{ shopStoreData.currentGoods.value.stock }}</el-tag>
            </div>
            <el-descriptions
              :column="1"
              border
              size="small"
              direction="horizontal"
            >
              <div v-if="shopStoreData.currentGoods.value.goods_type === constantStore.GOODS_TYPE_SUBSCRIBE">
                <el-descriptions-item :label="$t('message.adminShop.Goods.total_bandwidth')">{{ shopStoreData.currentGoods.value.total_bandwidth }}GB</el-descriptions-item>
                <el-descriptions-item :label="$t('message.adminShop.Goods.node_connector')">{{ shopStoreData.currentGoods.value.node_connector }}</el-descriptions-item>
                <el-descriptions-item :label="$t('message.adminShop.Goods.node_speed_limit')">{{ shopStoreData.currentGoods.value.node_speed_limit }}</el-descriptions-item>
              </div>
              <div v-if="shopStoreData.currentGoods.value.goods_type === constantStore.GOODS_TYPE_RECHARGE">
                <el-descriptions-item :label="$t('message.adminShop.Goods.recharge_amount')">{{ shopStoreData.currentGoods.value.recharge_amount }}</el-descriptions-item>
              </div>
            </el-descriptions>
            <div style="margin-top: 10px">
                <span>
                  <span style="color: red;">￥</span>
                  <span style="color: red;font-size: 30px;">{{ shopStoreData.currentGoods.value.price }}</span>
                  <span v-if="shopStoreData.currentGoods.value.goods_type === constantStore.GOODS_TYPE_SUBSCRIBE"> / {{$t('message.common.month')}}</span>
                </span>
            </div>
          </div>
        </el-col>

        <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
          <div style="margin-top: 10px;border-radius:10px;background:rgba(224,224,224,0.29)">
            <div style="margin-top: 10px;padding: 10px" v-html="shopStoreData.currentGoods.value.des"></div>
          </div>
        </el-col>
      </el-row>
      <template #footer>
				<span >
					<el-button @click="state.isShowGoodsDetails = false"
                     size="default">{{$t('message.common.button_cancel')}}
          </el-button>
					<el-button type="primary"
                     :disabled="shopStoreData.currentGoods.value.stock <= 0"
                     @click="openPurchase(shopStoreData.currentGoods.value)"
                     size="default">{{$t('message.adminShop.purchase')}}
          </el-button>
				</span>
      </template>
    </el-dialog>
    <Purchase ref="PurchaseRef"></Purchase>
  </div>
</template>

<script setup lang="ts" name="pagesLazyImg">
import { reactive, onMounted ,ref,defineAsyncComponent} from 'vue';
import { useRouter } from 'vue-router';
// import other from '/@/utils/other';

import { useShopStore } from "/@/stores/user_logic/shopStore";
import { storeToRefs } from "pinia";
import { useConstantStore } from "/@/stores/constantStore";
const shopStore = useShopStore();
const shopStoreData = storeToRefs(shopStore);
const constantStore = useConstantStore()
const Purchase = defineAsyncComponent(() => import("/@/views/shop/purchase.vue"));
const PurchaseRef = ref();

// 定义变量内容
const router = useRouter();
const state = reactive({
  isShowGoodsDetails:false,
  isShowLoading:false,
  goods_type: constantStore.GOODS_TYPE_SUBSCRIBE,
});

const showGoodsDetails = (v: Goods) => {
  shopStoreData.currentGoods.value = v
  state.isShowGoodsDetails = true
};
const openPurchase = () => {
  state.isShowGoodsDetails = false
  state.isShowLoading = true
  //初始化购买商品信息
  // if (goods){
  //   shopStoreData.currentGoods.value = goods;
  // }
  shopStoreData.currentOrder.value.duration = 1; // 默认订购时长
  shopStoreData.currentOrder.value.order_type = constantStore.ORDER_TYPE_NEW //订单类型：新购入
  shopStoreData.currentOrder.value.goods_id = shopStoreData.currentGoods.value.id //订购商品ID

  PurchaseRef.value.openDialog(constantStore.ORDER_TYPE_NEW);


  //获取计算价格后的订单信息
  // shopStore.getOrderInfo(shopStoreData.currentOrder.value).then((res) => {
  //   state.isShowLoading = false
  //   shopStoreData.currentOrder.value = {} as Order
  //   shopStoreData.currentOrder.value = res.data
  //   PurchaseRef.value.openDialog();
  // })
};
const getAllEnabledGoods = () => {
  shopStore.getAllEnabledGoods({ goods_type: state.goods_type });
};
// 页面加载时
onMounted(() => {
  getAllEnabledGoods();
});
</script>

<style scoped lang="scss">
.lazy-img-container {
  .flex-warp {
    display: flex;
    flex-wrap: wrap;
    align-content: flex-start;
    margin: 0 -5px;
    .flex-warp-item {
      padding: 5px;
      width: 100%;
      height: 360px;
      .flex-warp-item-box {
        border: 1px solid var(--next-border-color-light);
        width: 100%;
        height: 100%;
        border-radius: 2px;
        display: flex;
        flex-direction: column;
        transition: all 0.3s ease;
        &:hover {
          cursor: pointer;
          border: 1px solid var(--el-color-primary);
          transition: all 0.3s ease;
          box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.03);
          .item-txt-title {
            color: var(--el-color-primary) !important;
            transition: all 0.3s ease;
          }
        }
        .item-img {
          width: 100%;
          height: 210px;
          overflow: hidden;
          img {
            transition: all 0.3s ease;
            width: 100%;
            height: 100%;
          }
        }
        .item-txt {
          flex: 1;
          padding: 15px;
          display: flex;
          flex-direction: column;
          overflow: hidden;
          .item-txt-title {
            text-overflow: ellipsis;
            overflow: hidden;
            -webkit-line-clamp: 2;
            -webkit-box-orient: vertical;
            display: -webkit-box;
            color: #666666;
            transition: all 0.3s ease;
            &:hover {
              color: var(--el-color-primary);
              text-decoration: underline;
              transition: all 0.3s ease;
            }
          }
          .item-txt-other {
            flex: 1;
            align-items: flex-end;
            display: flex;
            .item-txt-msg {
              font-size: 12px;
              color: #8d8d91;
            }
            .item-txt-price {
              display: flex;
              justify-content: space-between;
              align-items: center;
              .font-price {
                color: #ff5000;
                .font {
                  font-size: 22px;
                }
              }
            }
          }
        }
      }
    }
  }
}
</style>
