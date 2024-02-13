<template>
  <div>
    <el-drawer v-model="state.isShowDrawer" size="90%">
      <el-table :data="adminCustomerServiceStoreData.customerServiceList.value.data" stripe style="width: 100%;flex: 1;">
<!--        <el-table-column type="index" label="序号" width="60" fixed/>-->

        <el-table-column prop="id" label="服务ID" width="60" fixed/>
        <el-table-column prop="user_name" label="用户名" width="120" fixed/>
        <el-table-column prop="user_id" label="用户ID" width="60"/>
        <el-table-column prop="service_status" label="服务状态" width="80">
          <template #default="{row}">
            <el-button v-if="row.service_status" type="success">有效</el-button>
            <el-button v-else type="info">失效</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="service_start_at" label="开始时间" width="150">
          <template #default="{row}">
            {{ DateStrToTime(row.service_start_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="service_end_at" label="结束时间" width="150">
          <template #default="{row}">
            {{ DateStrToTime(row.service_end_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="is_renew" label="是否可续费" width="100">
          <template #default="{row}">
            <el-button v-if="row.is_renew" type="success">是</el-button>
            <el-button v-else type="info">否</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="renewal_amount" label="续费价格" width="100"/>
        <el-table-column prop="goods_id" label="商品ID" width="100"/>
        <el-table-column prop="subject" label="商品标题" width="100"/>
        <el-table-column prop="price" label="商品价格" width="100"/>
        <el-table-column prop="goods_type" label="商品类型" width="100"/>
        <el-table-column prop="duration" label="订购时长" width="100"/>
        <div>
          <el-table-column prop="total_bandwidth" label="总流量" width="60"/>
          <el-table-column prop="node_speed_limit" label="限速" width="60"/>
          <el-table-column prop="node_connector" label="连接数" width="60"/>
          <el-table-column prop="traffic_reset_day" label="流量重置日" width="60"/>
          <el-table-column prop="sub_uuid" label="订阅UUID" width="60"/>
          <el-table-column prop="used_up" label="已用上行" width="60"/>
          <el-table-column prop="used_down" label="已用下行" width="60"/>

        </div>
        <el-table-column label="操作">
          <template #default="{row}">
            <el-button size="small" text type="primary" @click="openDialog(row)">修改</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-drawer>
    <CustomerServiceEditDialog ref="customerServiceEditDialogRef"></CustomerServiceEditDialog>
  </div>
</template>

<script setup lang="ts">
import { reactive,ref,defineAsyncComponent } from "vue";
import { useAdminCustomerServiceStore } from "/@/stores/admin_logic/customerServiceStore";
import { storeToRefs } from "pinia";
import { DateStrToTime } from "/@/utils/formatTime"
const adminCustomerServiceStore = useAdminCustomerServiceStore()
const adminCustomerServiceStoreData = storeToRefs(adminCustomerServiceStore)


const CustomerServiceEditDialog =  defineAsyncComponent(() => import('/@/views/admin/user/dialog_customer_service_edit.vue'));
const customerServiceEditDialogRef = ref()
const state = reactive({
  // title:'客户服务',
  isShowDialog:false,
  isShowDrawer: false,
})
const onOpenDrawer=(row:SysUser)=>{
  state.isShowDrawer = true
  adminCustomerServiceStore.getCustomerServiceList({user_id:row.id} as CustomerService)
}
const openDialog = (row:CustomerService)=>{
  state.isShowDrawer = true
  customerServiceEditDialogRef.value.openDialog(row)
}
// 暴露变量
defineExpose({
  onOpenDrawer,
});
</script>

<style scoped>

</style>