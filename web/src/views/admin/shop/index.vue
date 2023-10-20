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
      <el-table :data="goodsList" height="100%" style="width: 100%;flex: 1;">
        <el-table-column type="index" label="序号" width="60" fixed/>
        <el-table-column prop="subject" label="套餐名称" show-overflow-tooltip width="180"></el-table-column>
        <el-table-column prop="id" label="套餐ID" show-overflow-tooltip width="60"></el-table-column>
        <el-table-column prop="total_amount" label="套餐价格(元)" show-overflow-tooltip width="100"></el-table-column>
        <el-table-column prop="total_bandwidth" label="总流量(GB)" show-overflow-tooltip width="100"></el-table-column>
        <el-table-column prop="expiration_date" label="有效期(天)" show-overflow-tooltip></el-table-column>
        <el-table-column prop="status" label="是否显示" show-overflow-tooltip>
          <template #default="scope">
            <el-tag class="ml-2" v-if="scope.row.status" type="success">启用</el-tag>
            <el-tag class="ml-2" v-if="!scope.row.status" type="danger">禁用</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
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
    <ShopDialog ref="shopDialogRef" @refresh="shopStore.getAllGoods()"></ShopDialog>
    <SortDialog ref="sortDialogRef"></SortDialog>
  </div>
</template>

<script setup lang="ts">
import {defineAsyncComponent, onMounted, ref} from 'vue';
import {storeToRefs} from 'pinia';
import {useShopStore} from "/@/stores/shopStore";
import {useNodeStore} from "/@/stores/nodeStore";
import {ElMessageBox} from "element-plus";

const shopStore = useShopStore()
const {goodsList, goodsManageData} = storeToRefs(shopStore)
const nodeStore = useNodeStore()
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
        goodsManageData.value.currentGoods = row
        shopStore.deleteGoods(row)
        //延迟2秒
        setTimeout(() => {
          shopStore.getAllGoods()
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
  shopStore.getAllGoods() //获取全部商品
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