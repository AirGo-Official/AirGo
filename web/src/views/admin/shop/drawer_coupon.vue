<template>
  <div class="container layout-padding">
    <el-drawer v-model="state.isShowDrawer" size="90%" :title="$t('message.adminShop.coupon')" destroy-on-close>
      <el-card shadow="hover" class="layout-padding-auto">
        <div class="mb15">
          <el-button size="default" type="primary" class="ml10" @click="openDialog('add')">
            <el-icon>
              <ele-FolderAdd/>
            </el-icon>
            {{$t('message.adminShop.addCoupon')}}
          </el-button>
        </div>
        <div>
          <el-table :data="shopStoreData.couponList.value.data" stripe style="width: 100%;flex: 1;">
            <el-table-column type="index" :label="$t('message.adminShop.Coupon.index')" width="60"/>
            <el-table-column prop="name" :label="$t('message.adminShop.Coupon.name')" width="180"/>
            <el-table-column prop="id" :label="$t('message.adminShop.Coupon.id')" width="60"/>
            <el-table-column prop="discount_rate" :label="$t('message.adminShop.Coupon.discount_rate')" width="60"/>
            <el-table-column prop="limit" :label="$t('message.adminShop.Coupon.limit')" width="100"/>
            <el-table-column prop="expired_at" :label="$t('message.adminShop.Coupon.expired_at')" width="180">
              <template #default="{row}">
                {{ DateStrToTime(row.expired_at) }}
              </template>
            </el-table-column>
            <el-table-column prop="min_amount" :label="$t('message.adminShop.Coupon.min_amount')" width="180"></el-table-column>
            <el-table-column :label="$t('message.common.operate')">
              <template #default="scope">
                <el-button text type="primary" @click="openDialog('edit',scope.row)">{{$t('message.common.modify')}}</el-button>
                <el-button text type="danger" @click="deleteCoupon(scope.row)">{{$t('message.common.delete')}}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-card>
      <el-dialog v-model="state.isShowDialog" :title="state.title" width="80%">
        <el-form :model="shopStoreData.currentCoupon.value" label-position="top">
          <el-form-item :label="$t('message.adminShop.Coupon.name')">
            <el-input v-model="shopStoreData.currentCoupon.value.name"></el-input>
          </el-form-item>
          <el-form-item :label="$t('message.adminShop.Coupon.discount_rate')">
            <el-col :span="6">
              <el-input-number v-model.number="shopStoreData.currentCoupon.value.discount_rate" :min="0" :max="1" :step="0.1"></el-input-number>
              <div style="color: #9b9da1">{{$t('message.adminShop.couponRateTip')}}</div>
            </el-col>

          </el-form-item>
          <el-form-item :label="$t('message.adminShop.Coupon.limit')">
            <el-col :span="6">
              <el-input-number v-model.number="shopStoreData.currentCoupon.value.limit" :min="0" :step="1"></el-input-number>
            </el-col>
            <el-col :span="20"></el-col>
          </el-form-item>
          <el-form-item :label="$t('message.adminShop.Coupon.expired_at')">
            <el-date-picker
              v-model="shopStoreData.currentCoupon.value.expired_at"
              type="datetime"
              size="default"
            />
          </el-form-item>
          <el-form-item :label="$t('message.adminShop.Coupon.min_amount')">
            <el-input-number v-model.number="shopStoreData.currentCoupon.value.min_amount" :min="0.00" :step="1" :precision="2"></el-input-number>
          </el-form-item>

          <el-form-item :label="$t('message.adminShop.Coupon.goods')">
            <el-tree ref="goods_tree_ref" node-key="id"
                     :data="shopStoreData.goodsList.value"
                     :props="{label:'subject'}"
                     :default-checked-keys="shopStoreData.checkedGoodsIDs.value"
                     show-checkbox class="menu-data-tree"/>
          </el-form-item>
        </el-form>
        <template #footer>
            <span class="dialog-footer">
                <el-button @click="closeDialog">{{$t('message.common.button_cancel')}}</el-button>
                <el-button type="primary" @click="onSubmit(shopStoreData.currentCoupon.value)" color="#FC3D08">
                    {{$t('message.common.button_confirm')}}
                </el-button>
            </span>
        </template>
      </el-dialog>
    </el-drawer>
  </div>
</template>

<script lang="ts" setup>

import {defineAsyncComponent, onMounted, reactive, ref} from "vue";
import {DateStrToTime} from "/@/utils/formatTime";
import {ElMessage, ElMessageBox} from "element-plus";
import {useApiStore} from "/@/stores/apiStore";
import {storeToRefs} from "pinia";
import { useAdminShopStore } from "/@/stores/admin_logic/shopStore";
import { useI18n } from "vue-i18n";

const shopStore = useAdminShopStore()
const shopStoreData = storeToRefs(shopStore)
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
const goods_tree_ref = ref()
const {t} = useI18n()
const state = reactive({
  isShowDrawer:false,
  isShowDialog: false,
  title:'',
  type:'',
})
const openDrawer=()=>{
  state.isShowDrawer = true
  getCouponList()
}
//
const getCouponList = () => {
  shopStore.getCouponList()
}
//
const openDialog = (type: string, row?: any) => {
  //获取全部商品
  shopStore.getGoodsList() //获取全部商品
  state.isShowDialog = true
  state.type = type
  switch (type) {
    case "add":
      state.title = t('message.adminShop.addCoupon')
      // shopStoreData.currentCoupon.value = {} as Coupon
      shopStoreData.currentCoupon.value.id = 0
      break
    case "edit":
      state.title = t('message.adminShop.modifyCoupon')
      shopStoreData.currentCoupon.value = row
      shopStore.goodsIDsHandler()
      break
  }
}
//关闭弹窗
const closeDialog = () => {
  state.isShowDialog = false
}
//删除折扣
const deleteCoupon = (row: Coupon) => {
  ElMessageBox.confirm(t('message.common.message_confirm_delete'), t('message.common.tip'), {
    confirmButtonText:  t('message.common.button_confirm'),
    cancelButtonText:  t('message.common.button_cancel'),
    type: 'warning',
  })
      .then(() => {
        shopStore.deleteCoupon(row).then((res) => {
          getCouponList()
        })
      })
      .catch(() => {
      });
}
//提交按钮
const onSubmit = (params: object) => {
  //处理折扣码关联的商品
  shopStoreData.checkedGoodsIDs.value = [...goods_tree_ref.value.getCheckedKeys()];
  switch (state.type) {
    case "add":
      shopStore.newCoupon().then((res) => {
        getCouponList()
      })
      break
    case "edit":
      shopStore.updateCoupon().then((res) => {
        getCouponList()
      })
      break
    default:
      break
  }
  closeDialog()
}
//
onMounted(() => {
  getCouponList()
});
// 暴露变量
defineExpose({
  openDrawer,
});

</script>

<style scoped lang="scss">

.container {
  :deep(.el-card__body) {
    display: flex;
    flex-direction: column;
    flex: 1;
    overflow: auto;

    .el-table {
      flex: 1;
    }
  }
}

</style>