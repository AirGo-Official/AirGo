<template>
  <div class="container layout-padding">
    <el-card shadow="hover" class="layout-padding-auto">
      <div class="system-user-search mb15">
        <el-button size="default" type="primary" class="ml10" @click="openGoodsDialog('add')">
          <el-icon>
            <ele-FolderAdd/>
          </el-icon>
          {{$t('message.adminShop.addGoods')}}
        </el-button>
        <el-button size="default" type="primary" class="ml10" @click="onOpenSortDialog">
          <el-icon><DCaret/></el-icon>{{$t('message.adminShop.sort')}}
        </el-button>
        <el-button size="default" type="primary" class="ml10" @click="openCouponDrawer">
          <el-icon><Ticket /></el-icon>{{$t('message.adminShop.coupon')}}
        </el-button>
      </div>
      <el-table :data="shopStoreData.goodsList.value" height="100%" style="width: 100%;flex: 1;">
        <el-table-column type="index" :label="$t('message.adminShop.Goods.index')" width="60" fixed/>
        <el-table-column prop="subject" :label="$t('message.adminShop.Goods.subject')" show-overflow-tooltip width="200"></el-table-column>
        <el-table-column prop="id" :label="$t('message.adminShop.Goods.id')" show-overflow-tooltip width="60"></el-table-column>
        <el-table-column prop="is_show" :label="$t('message.adminShop.Goods.is_show')" show-overflow-tooltip width="100">
          <template #default="scope">
            <el-tag class="ml-2" v-if="scope.row.is_show" type="success">{{$t('message.common.enable')}}</el-tag>
            <el-tag class="ml-2" v-else type="danger">{{$t('message.common.disable')}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="is_sale" :label="$t('message.adminShop.Goods.is_sale')" show-overflow-tooltip width="100">
          <template #default="scope">
            <el-tag class="ml-2" v-if="scope.row.is_sale" type="success">{{$t('message.common.enable')}}</el-tag>
            <el-tag class="ml-2" v-else type="danger">{{$t('message.common.disable')}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="is_renew" :label="$t('message.adminShop.Goods.is_renew')" show-overflow-tooltip width="100">
          <template #default="scope">
            <el-tag class="ml-2" v-if="scope.row.is_renew" type="success">{{$t('message.common.enable')}}</el-tag>
            <el-tag class="ml-2" v-else type="danger">{{$t('message.common.disable')}}</el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="goods_type" :label="$t('message.adminShop.Goods.goods_type')" show-overflow-tooltip width="100">
          <template #default="{ row }">
            <el-tag class="ml-2" v-if="row.goods_type === constantStore.GOODS_TYPE_SUBSCRIBE">{{$t('message.constant.GOODS_TYPE_SUBSCRIBE')}}</el-tag>
            <el-tag class="ml-2" v-if="row.goods_type === constantStore.GOODS_TYPE_RECHARGE" >{{$t('message.constant.GOODS_TYPE_RECHARGE')}}</el-tag>
            <el-tag class="ml-2" v-if="row.goods_type === constantStore.GOODS_TYPE_GENERAL">{{$t('message.constant.GOODS_TYPE_GENERAL')}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="deliver_type" :label="$t('message.adminShop.Goods.deliver_type')" show-overflow-tooltip width="100">
          <template #default="{ row }">
            <el-tag class="ml-2" v-if="row.deliver_type === constantStore.DELIVER_TYPE_NONE" type="success">{{$t('message.constant.DELIVER_TYPE_NONE')}}</el-tag>
            <el-tag class="ml-2" v-if="row.deliver_type === constantStore.DELIVER_TYPE_AUTO" type="success">{{$t('message.constant.DELIVER_TYPE_AUTO')}}</el-tag>
            <el-tag class="ml-2" v-if="row.deliver_type === constantStore.DELIVER_TYPE_MANUAL" type="success">{{$t('message.constant.DELIVER_TYPE_MANUAL')}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="price" :label="$t('message.adminShop.Goods.price')" show-overflow-tooltip width="200"></el-table-column>
        <el-table-column :label="$t('message.common.operate')">
          <template #default="scope">
            <el-button size="small" text type="primary"
                       @click="openGoodsDialog('edit',scope.row)">{{$t('message.common.modify')}}
            </el-button>
            <el-button size="small" text type="danger"
                       @click="onRowDel(scope.row)">{{$t('message.common.delete')}}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <!-- 引入弹窗组件 -->
    <ShopDialog ref="shopDialogRef" @refresh="shopStore.getGoodsList()"></ShopDialog>
    <SortDialog ref="sortDialogRef"></SortDialog>
    <CouponDrawer ref="CouponDrawerRef"></CouponDrawer>
  </div>
</template>

<script setup lang="ts">
import {defineAsyncComponent, onMounted, ref} from 'vue';
import {storeToRefs} from 'pinia';
import {useShopStore} from "/@/stores/user_logic/shopStore";
import {useAdminNodeStore} from "/@/stores/admin_logic/nodeStore";
import {ElMessageBox} from "element-plus";
import { useAdminShopStore } from "/@/stores/admin_logic/shopStore";
import { useI18n } from "vue-i18n";
import { useConstantStore } from "/@/stores/constantStore";

const shopStore = useAdminShopStore()
const shopStoreData = storeToRefs(shopStore)
const constantStore = useConstantStore()

const nodeStore = useAdminNodeStore()
const ShopDialog = defineAsyncComponent(() => import('/@/views/admin/shop/dialog_edit.vue'))
const SortDialog = defineAsyncComponent(() => import('/@/views/admin/shop/dialog_sort.vue'))
const CouponDrawer = defineAsyncComponent(() => import('/@/views/admin/shop/drawer_coupon.vue'))
const shopDialogRef = ref()
const sortDialogRef = ref()
const CouponDrawerRef = ref()
const {t} = useI18n()

//修改套餐弹窗
const openGoodsDialog = (type:string,row?: Goods) => {
  shopDialogRef.value.openDialog(type, row)
}

const openCouponDrawer=()=>{
  CouponDrawerRef.value.openDrawer()
}
//删除套餐
const onRowDel = (row: Goods) => {
  ElMessageBox.confirm(t('message.common.message_confirm_delete'), t('message.common.tip'), {
    confirmButtonText: t('message.common.button_confirm'),
    cancelButtonText: t('message.common.button_cancel'),
    type: 'warning',
  })
      .then(() => {
        shopStoreData.currentGoods.value = row
        shopStore.deleteGoods(row)
        //延迟2秒
        setTimeout(() => {
          shopStore.getGoodsList()
        }, 500)
      })
      .catch(() => {
      });
}

//排序弹窗
function onOpenSortDialog() {
  sortDialogRef.value.openDialog()
}

//加载时
onMounted(() => {
  shopStore.getGoodsList() //获取全部商品
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