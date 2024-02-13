<template>
  <div class="container layout-padding">
    <el-card shadow="hover" class="layout-padding-auto">
      <div class="system-user-search mb15">
        <el-button size="default" type="success" class="ml10" @click="onOpenAddGoods">
          <el-icon>
            <ele-FolderAdd/>
          </el-icon>
          新增套餐
        </el-button>
        <el-button size="default" type="warning" class="ml10" @click="onOpenSortDialog">
          <el-icon>
            <DCaret/>
          </el-icon>
          排序
        </el-button>
      </div>
      <el-table :data="shopStoreData.goodsList.value" height="100%" style="width: 100%;flex: 1;">
        <el-table-column type="index" label="序号" width="60" fixed/>
        <el-table-column prop="subject" label="套餐名称" show-overflow-tooltip width="200"></el-table-column>
        <el-table-column prop="id" label="套餐ID" show-overflow-tooltip width="60"></el-table-column>
        <el-table-column prop="status" label="是否显示" show-overflow-tooltip width="100">
          <template #default="scope">
            <el-tag class="ml-2" v-if="scope.row.is_show" type="success">启用</el-tag>
            <el-tag class="ml-2" v-if="!scope.row.is_show" type="danger">禁用</el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="goods_type" label="类型" show-overflow-tooltip width="100">
          <template #default="{ row }">
            <el-tag class="ml-2" v-if="row.goods_type === 'subscribe'" type="success">订阅</el-tag>
            <el-tag class="ml-2" v-if="row.goods_type === 'recharge'" type="danger">充值</el-tag>
            <el-tag class="ml-2" v-if="row.goods_type === 'general'" type="warning">普通商品</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="deliver_type" label="发货类型" show-overflow-tooltip width="100">
          <template #default="{ row }">
            <el-tag class="ml-2" v-if="row.deliver_type === 'none'" type="success">不发货</el-tag>
            <el-tag class="ml-2" v-if="row.deliver_type === 'auto'" type="success">自动发货</el-tag>
            <el-tag class="ml-2" v-if="row.deliver_type === 'manual'" type="success">手动发货</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="price" label="价格(元)" show-overflow-tooltip width="200"></el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button size="small" text type="primary"
                       @click="onOpenEditGoods(scope.row)">修改
            </el-button>
            <el-button size="small" text type="primary"
                       @click="onRowDel(scope.row)">删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <!-- 引入弹窗组件 -->
    <ShopDialog ref="shopDialogRef" @refresh="shopStore.getGoodsList()"></ShopDialog>
    <SortDialog ref="sortDialogRef"></SortDialog>
  </div>
</template>

<script setup lang="ts">
import {defineAsyncComponent, onMounted, ref} from 'vue';
import {storeToRefs} from 'pinia';
import {useShopStore} from "/@/stores/user_logic/shopStore";
import {useAdminNodeStore} from "/@/stores/admin_logic/nodeStore";
import {ElMessageBox} from "element-plus";
import { useAdminShopStore } from "/@/stores/admin_logic/shopStore";

const shopStore = useAdminShopStore()
const shopStoreData = storeToRefs(shopStore)

const nodeStore = useAdminNodeStore()
const ShopDialog = defineAsyncComponent(() => import('/@/views/admin/shop/dialog_edit.vue'))
const SortDialog = defineAsyncComponent(() => import('/@/views/admin/shop/dialog_sort.vue'))
const shopDialogRef = ref()
const sortDialogRef = ref()


//修改套餐弹窗
const onOpenEditGoods = (row: Goods) => {
  shopDialogRef.value.openDialog('edit', row)
}
//添加套餐弹窗
const onOpenAddGoods = () => {
  shopDialogRef.value.openDialog('add')
}
//删除套餐
const onRowDel = (row: Goods) => {
  ElMessageBox.confirm(`此操作将永久删除商品：${row.subject}, 是否继续?`, '提示', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
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
  nodeStore.getAllNode()  //获取全部节点
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