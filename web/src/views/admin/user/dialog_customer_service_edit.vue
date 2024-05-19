<template>
  <div>
    <el-dialog v-model="state.isShowDialog" width="80%">
      <el-form ref="userDialogFormRef" :model="adminCustomerServiceStoreData.currentCustomerService.value"
               size="default" label-position="top">
        <el-form-item :label="$t('message.adminUser.CustomerService.service_start_at')">
          <el-date-picker
            v-model="adminCustomerServiceStoreData.currentCustomerService.value.service_start_at"
            type="datetime"
            size="default"
          />
        </el-form-item>
        <el-form-item :label="$t('message.adminUser.CustomerService.service_end_at')">
          <el-date-picker
            v-model="adminCustomerServiceStoreData.currentCustomerService.value.service_end_at"
            type="datetime"
            size="default"
          />
        </el-form-item>
        <el-form-item :label="$t('message.adminUser.CustomerService.service_status')">
          <el-switch v-model="adminCustomerServiceStoreData.currentCustomerService.value.service_status" inline-prompt
                     :active-text="$t('message.common.enable')"
                     :inactive-text="$t('message.common.disable')"></el-switch>
        </el-form-item>
        <el-form-item :label="$t('message.adminUser.CustomerService.is_renew')">
          <el-switch v-model="adminCustomerServiceStoreData.currentCustomerService.value.is_renew" inline-prompt
                     :active-text="$t('message.common.enable')"
                     :inactive-text="$t('message.common.disable')"></el-switch>
        </el-form-item>
        <el-form-item :label="$t('message.adminUser.CustomerService.renewal_amount')">
          <el-input v-model="adminCustomerServiceStoreData.currentCustomerService.value.renewal_amount"></el-input>
        </el-form-item>
        <el-form-item :label="$t('message.adminUser.CustomerService.duration')">
          <el-input-number
            v-model.number="adminCustomerServiceStoreData.currentCustomerService.value.duration"></el-input-number>
        </el-form-item>
        <el-form-item :label="$t('message.adminUser.CustomerService.goods_id')">
          <el-select
            v-model="adminCustomerServiceStoreData.currentCustomerService.value.goods_id"
            filterable
            allow-create
            default-first-option
            :reserve-keyword="false"
            style="width: 100%"
            @change="goodsChanged()"
          >
            <el-option
              v-for="(v,k) in adminShopStoreData.goodsList.value"
              :key="k"
              :label="v.subject"
              :value="v.id">
              <span style="float: left">{{ v.subject }}</span>
              <span style="float: right;color: var(--el-text-color-secondary);font-size: 13px;">ID: {{ v.id }}</span
              >
            </el-option>
          </el-select>
        </el-form-item>
        <el-divider></el-divider>
        <div
          v-if="adminCustomerServiceStoreData.currentCustomerService.value.goods_type === constantStore.GOODS_TYPE_SUBSCRIBE">
          <el-form-item :label="$t('message.adminUser.CustomerService.sub_status')">
            <el-switch v-model="adminCustomerServiceStoreData.currentCustomerService.value.sub_status" inline-prompt
                       :active-text="$t('message.common.enable')"
                       :inactive-text="$t('message.common.disable')"></el-switch>
          </el-form-item>
          <el-form-item :label="$t('message.adminUser.CustomerService.total_bandwidth')">
            <el-input-number
              v-model.number="adminCustomerServiceStoreData.currentCustomerService.value.total_bandwidth"></el-input-number>
          </el-form-item>
          <el-form-item :label="$t('message.adminUser.CustomerService.node_speed_limit')">
            <el-input-number
              v-model.number="adminCustomerServiceStoreData.currentCustomerService.value.node_speed_limit"></el-input-number>
          </el-form-item>
          <el-form-item :label="$t('message.adminUser.CustomerService.node_connector')">
            <el-input-number
              v-model.number="adminCustomerServiceStoreData.currentCustomerService.value.node_connector"></el-input-number>
          </el-form-item>
          <el-form-item :label="$t('message.adminUser.CustomerService.traffic_reset_day')">
            <el-input-number
              v-model.number="adminCustomerServiceStoreData.currentCustomerService.value.traffic_reset_day"></el-input-number>
          </el-form-item>
          <el-form-item :label="$t('message.adminUser.CustomerService.sub_uuid')">
            <el-input v-model="adminCustomerServiceStoreData.currentCustomerService.value.sub_uuid">
              <template #append>
                <el-button @click="resetSubUUID">{{ $t("message.common.reset") }}</el-button>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item :label="$t('message.adminUser.CustomerService.used_up')">
            <el-input-number
              v-model.number="adminCustomerServiceStoreData.currentCustomerService.value.used_up"></el-input-number>
          </el-form-item>
          <el-form-item :label="$t('message.adminUser.CustomerService.used_down')">
            <el-input-number
              v-model.number="adminCustomerServiceStoreData.currentCustomerService.value.used_down"></el-input-number>
          </el-form-item>
        </div>
      </el-form>
      <template #footer>
				<span class="dialog-footer">
					<el-button @click="closeDialog" size="default">{{ $t("message.common.button_cancel") }}</el-button>
					<el-button type="primary" @click="onSubmit"
                     size="default">{{ $t("message.common.button_confirm") }}</el-button>
				</span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive } from "vue";
import { useAdminCustomerServiceStore } from "/@/stores/admin_logic/customerServiceStore";
import { storeToRefs } from "pinia";
import { DateStrToTime } from "/@/utils/formatTime";
import { useConstantStore } from "/@/stores/constantStore";
import { v4 as uuid } from "uuid";
import { useAdminShopStore } from "/@/stores/admin_logic/shopStore";
import { useAdminUserStore } from "/@/stores/admin_logic/userStore";

const adminCustomerServiceStore = useAdminCustomerServiceStore();
const adminCustomerServiceStoreData = storeToRefs(adminCustomerServiceStore);
const constantStore = useConstantStore();
const adminShopStore = useAdminShopStore()
const adminShopStoreData = storeToRefs(adminShopStore)
const adminUserStore = useAdminUserStore()
const adminUserStoreData = storeToRefs(adminUserStore)

const state = reactive({
  isShowDialog: false,
  type:'',
});
const openDialog = (type: string, row?: CustomerService) => {
  state.type = type;
  state.isShowDialog = true;
  adminShopStore.getGoodsList()
  if (type === "edit") {
    adminCustomerServiceStoreData.currentCustomerService.value = row;
    //处理一下流量
    adminCustomerServiceStoreData.currentCustomerService.value.total_bandwidth = Number((adminCustomerServiceStoreData.currentCustomerService.value.total_bandwidth / 1024 / 1024 / 1024).toFixed(2));
    adminCustomerServiceStoreData.currentCustomerService.value.used_up = Number((adminCustomerServiceStoreData.currentCustomerService.value.used_up / 1024 / 1024 / 1024).toFixed(2));
    adminCustomerServiceStoreData.currentCustomerService.value.used_down = Number((adminCustomerServiceStoreData.currentCustomerService.value.used_down / 1024 / 1024 / 1024).toFixed(2));
  } else {
    adminCustomerServiceStoreData.currentCustomerService.value = {} as CustomerService
    adminCustomerServiceStoreData.currentCustomerService.value.user_id = adminUserStoreData.currentUser.value.id
    adminCustomerServiceStoreData.currentCustomerService.value.user_name = adminUserStoreData.currentUser.value.user_name
    adminCustomerServiceStoreData.currentCustomerService.value.duration = 1
  }

};

const closeDialog = () => {
  state.isShowDialog = false;
};

const onSubmit = () => {
//处理一下流量
  adminCustomerServiceStoreData.currentCustomerService.value.total_bandwidth = Number((adminCustomerServiceStoreData.currentCustomerService.value.total_bandwidth * 1024 * 1024 * 1024).toFixed(0));
  adminCustomerServiceStoreData.currentCustomerService.value.used_up = Number((adminCustomerServiceStoreData.currentCustomerService.value.used_up * 1024 * 1024 * 1024).toFixed(0));
  adminCustomerServiceStoreData.currentCustomerService.value.used_down = Number((adminCustomerServiceStoreData.currentCustomerService.value.used_down * 1024 * 1024 * 1024).toFixed(0));
  adminCustomerServiceStore.updateCustomerService(adminCustomerServiceStoreData.currentCustomerService.value);
  closeDialog();
};

const resetSubUUID = () => {
  adminCustomerServiceStoreData.currentCustomerService.value.sub_uuid = uuid();
};
const goodsChanged=()=>{
  adminShopStoreData.goodsList.value.forEach((item)=>{
    if(adminCustomerServiceStoreData.currentCustomerService.value.goods_id === item.id){
      adminCustomerServiceStoreData.currentCustomerService.value.subject = item.subject
      adminCustomerServiceStoreData.currentCustomerService.value.des = item.des
      adminCustomerServiceStoreData.currentCustomerService.value.price = item.price
      adminCustomerServiceStoreData.currentCustomerService.value.goods_type = item.goods_type
    }
  })

}

defineExpose({
  openDialog
});
</script>

<style scoped>

</style>