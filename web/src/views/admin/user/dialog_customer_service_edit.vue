<template>
  <div>
    <el-dialog :title="state.title" v-model="state.isShowDialog" width="80%">
      <el-form ref="userDialogFormRef" :model="adminCustomerServiceStoreData.currentCustomerService.value" size="default" label-width="90px">
        <el-form-item label="结束时间">
          <el-date-picker
            v-model="adminCustomerServiceStoreData.currentCustomerService.value.service_end_at"
            type="datetime"
            placeholder="选择结束时间"
            size="default"
          />
        </el-form-item>
        <el-form-item label="服务状态">
          <el-switch v-model="adminCustomerServiceStoreData.currentCustomerService.value.service_status" inline-prompt active-text="开启"
                     inactive-text="关闭"></el-switch>
        </el-form-item>
        <el-form-item label="是否可续费">
          <el-switch v-model="adminCustomerServiceStoreData.currentCustomerService.value.is_renew" inline-prompt active-text="是"
                     inactive-text="否"></el-switch>
        </el-form-item>
        <el-form-item label="续费价格">
          <el-input v-model="adminCustomerServiceStoreData.currentCustomerService.value.renewal_amount"></el-input>
        </el-form-item>
        <el-form-item label="订购时长">
          <el-input-number v-model.number="adminCustomerServiceStoreData.currentCustomerService.value.duration"></el-input-number>
        </el-form-item>
        <el-form-item label="订阅状态">
          <el-switch v-model="adminCustomerServiceStoreData.currentCustomerService.value.sub_status" inline-prompt active-text="有效"
                     inactive-text="失效"></el-switch>
        </el-form-item>
        <div v-if="adminCustomerServiceStoreData.currentCustomerService.value.goods_type === constantStore.GOODS_TYPE_SUBSCRIBE">
          <el-form-item label="总流量">
            <el-input-number v-model.number="adminCustomerServiceStoreData.currentCustomerService.value.total_bandwidth"></el-input-number>
          </el-form-item>
          <el-form-item label="限速">
            <el-input-number v-model.number="adminCustomerServiceStoreData.currentCustomerService.value.node_speed_limit"></el-input-number>
          </el-form-item>
          <el-form-item label="连接数">
            <el-input-number v-model.number="adminCustomerServiceStoreData.currentCustomerService.value.node_connector"></el-input-number>
          </el-form-item>
          <el-form-item label="流量重置日">
            <el-input-number v-model.number="adminCustomerServiceStoreData.currentCustomerService.value.traffic_reset_day"></el-input-number>
          </el-form-item>
          <el-form-item label="订阅UUID">
            <el-input v-model="adminCustomerServiceStoreData.currentCustomerService.value.sub_uuid">
              <template #append>
                <el-button @click="resetSubUUID">重置</el-button>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item label="已用上行">
            <el-input-number v-model.number="adminCustomerServiceStoreData.currentCustomerService.value.used_up"></el-input-number>
          </el-form-item>
          <el-form-item label="已用下行">
            <el-input-number v-model.number="adminCustomerServiceStoreData.currentCustomerService.value.used_down"></el-input-number>
          </el-form-item>
        </div>
      </el-form>
      <template #footer>
				<span class="dialog-footer">
					<el-button @click="closeDialog" size="default">取消</el-button>
					<el-button type="primary" @click="onSubmit" size="default">提交</el-button>
				</span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive } from "vue";
import { useAdminCustomerServiceStore } from "/@/stores/admin_logic/customerServiceStore";
import { storeToRefs } from "pinia";
import { DateStrToTime } from "/@/utils/formatTime"
import { useConstantStore } from "/@/stores/constantStore";
import { v4 as uuid } from 'uuid';

const adminCustomerServiceStore = useAdminCustomerServiceStore()
const adminCustomerServiceStoreData = storeToRefs(adminCustomerServiceStore)
const constantStore = useConstantStore()


const state = reactive({
  title:'客户服务',
  isShowDialog:false,
})
const openDialog = (row:CustomerService)=>{
  state.isShowDialog = true
  adminCustomerServiceStoreData.currentCustomerService.value=row
  //处理一下流量
  adminCustomerServiceStoreData.currentCustomerService.value.total_bandwidth =   Number((adminCustomerServiceStoreData.currentCustomerService.value.total_bandwidth/1024/1024/1024).toFixed(2))
  adminCustomerServiceStoreData.currentCustomerService.value.used_up =   Number((adminCustomerServiceStoreData.currentCustomerService.value.used_up/1024/1024/1024).toFixed(2))
  adminCustomerServiceStoreData.currentCustomerService.value.used_down =   Number((adminCustomerServiceStoreData.currentCustomerService.value.used_down/1024/1024/1024).toFixed(2))
}

const closeDialog=()=>{
  state.isShowDialog = false
}

const onSubmit=()=>{
//处理一下流量
  adminCustomerServiceStoreData.currentCustomerService.value.total_bandwidth =   Number((adminCustomerServiceStoreData.currentCustomerService.value.total_bandwidth*1024*1024*1024).toFixed(0))
  adminCustomerServiceStoreData.currentCustomerService.value.used_up = Number((adminCustomerServiceStoreData.currentCustomerService.value.used_up*1024*1024*1024).toFixed(0))
  adminCustomerServiceStoreData.currentCustomerService.value.used_down =  Number((adminCustomerServiceStoreData.currentCustomerService.value.used_down*1024*1024*1024).toFixed(0))
  adminCustomerServiceStore.updateCustomerService( adminCustomerServiceStoreData.currentCustomerService.value)
  closeDialog()
}

const resetSubUUID=()=>{
  adminCustomerServiceStoreData.currentCustomerService.value.sub_uuid = uuid()
}

defineExpose({
  openDialog,
});
</script>

<style scoped>

</style>