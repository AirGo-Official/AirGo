<template>
  <div>
    <el-drawer v-model="state.isShowDrawer" size="90%">
      <div style="margin: 15px">
        <el-button size="default" type="primary" class="ml10" @click="openDialog('add')">{{$t("message.common.add")}}</el-button>
      </div>

      <el-table :data="adminCustomerServiceStoreData.customerServiceList.value.data" stripe
                style="width: 100%;flex: 1;">
        <!--        <el-table-column type="index" :label="$t('message.adminUser.CustomerService.index')" width="60" fixed/>-->

        <el-table-column prop="id" :label="$t('message.adminUser.CustomerService.id')" width="60" fixed />
        <!--        <el-table-column prop="user_name" :label="$t('message.adminUser.CustomerService.user_name')" width="120" fixed/>-->
        <el-table-column prop="service_status" :label="$t('message.adminUser.CustomerService.service_status')"
                         width="80">
          <template #default="{row}">
            <el-button v-if="row.service_status" type="success">{{ $t("message.common.enable") }}</el-button>
            <el-button v-else type="info">{{ $t("message.common.disable") }}</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="service_start_at" :label="$t('message.adminUser.CustomerService.service_start_at')"
                         width="150">
          <template #default="{row}">
            {{ DateStrToTime(row.service_start_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="service_end_at" :label="$t('message.adminUser.CustomerService.service_end_at')"
                         width="150">
          <template #default="{row}">
            {{ DateStrToTime(row.service_end_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="is_renew" :label="$t('message.adminUser.CustomerService.is_renew')" width="100">
          <template #default="{row}">
            <el-button v-if="row.is_renew" type="success">{{ $t("message.common.yes") }}</el-button>
            <el-button v-else type="info">{{ $t("message.common.no") }}</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="renewal_amount" :label="$t('message.adminUser.CustomerService.renewal_amount')"
                         width="100" />
        <el-table-column prop="goods_id" :label="$t('message.adminUser.CustomerService.goods_id')" width="100" />
        <el-table-column prop="subject" :label="$t('message.adminUser.CustomerService.subject')" width="200"
                         show-overflow-tooltip />
        <el-table-column prop="price" :label="$t('message.adminUser.CustomerService.price')" width="100" />
        <el-table-column prop="goods_type" :label="$t('message.adminUser.CustomerService.goods_type')" width="100" />
        <el-table-column prop="duration" :label="$t('message.adminUser.CustomerService.duration')" width="100" />
        <div>
          <el-table-column prop="total_bandwidth" :label="$t('message.adminUser.CustomerService.total_bandwidth')"
                           width="60" />
          <el-table-column prop="node_speed_limit" :label="$t('message.adminUser.CustomerService.node_speed_limit')"
                           width="60" />
          <el-table-column prop="node_connector" :label="$t('message.adminUser.CustomerService.node_connector')"
                           width="60" />
          <el-table-column prop="traffic_reset_day" :label="$t('message.adminUser.CustomerService.traffic_reset_day')"
                           width="60" />
          <el-table-column prop="sub_uuid" :label="$t('message.adminUser.CustomerService.sub_uuid')" width="60" />
          <el-table-column prop="used_up" :label="$t('message.adminUser.CustomerService.used_up')" width="60" />
          <el-table-column prop="used_down" :label="$t('message.adminUser.CustomerService.used_down')" width="60" />

        </div>
        <el-table-column :label="$t('message.common.operate')">
          <template #default="{row}">
            <el-button size="small" text type="primary" @click="openDialog('edit',row)">{{ $t("message.common.modify") }}
            </el-button>
            <el-button size="small" text type="danger" @click="deleteCustomerService(row)">{{ $t("message.common.delete") }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-drawer>
    <CustomerServiceEditDialog ref="customerServiceEditDialogRef"></CustomerServiceEditDialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, defineAsyncComponent } from "vue";
import { useAdminCustomerServiceStore } from "/@/stores/admin_logic/customerServiceStore";
import { storeToRefs } from "pinia";
import { DateStrToTime } from "/@/utils/formatTime";
import { ElMessageBox } from "element-plus";
import { useI18n } from "vue-i18n";
import { useAdminUserStore } from "/@/stores/admin_logic/userStore";
const { t } = useI18n();

const adminCustomerServiceStore = useAdminCustomerServiceStore();
const adminCustomerServiceStoreData = storeToRefs(adminCustomerServiceStore);
const adminUserStore = useAdminUserStore()
const adminUserStoreData = storeToRefs(adminUserStore)


const CustomerServiceEditDialog = defineAsyncComponent(() => import("/@/views/admin/user/dialog_customer_service_edit.vue"));
const customerServiceEditDialogRef = ref();
const state = reactive({
  isShowDialog: false,
  isShowDrawer: false,
  user_id:0,
});
const openDrawer = (row: SysUser) => {
  state.isShowDrawer = true;
  //存储当前用户信息
  adminUserStoreData.currentUser.value = row
  getData()
};
const getData=()=>{
  adminCustomerServiceStore.getCustomerServiceList({ user_id: adminUserStoreData.currentUser.value.id} as CustomerService);
}
const openDialog = (type: string, row?: CustomerService) => {
  state.isShowDrawer = true;
  customerServiceEditDialogRef.value.openDialog(type,row);
};

const deleteCustomerService=(row: CustomerService)=>{
  ElMessageBox.confirm(t("message.common.message_confirm_delete"), t("message.common.tip"), {
    confirmButtonText: t("message.common.button_confirm"),
    cancelButtonText: t("message.common.button_cancel"),
    type: "warning"
  })
    .then(() => {
      adminCustomerServiceStore.deleteCustomerService(row).then(()=>{
        getData()
      })
    })
    .catch(() => {
    });
}
// 暴露变量
defineExpose({
  openDrawer
});
</script>

<style scoped>

</style>