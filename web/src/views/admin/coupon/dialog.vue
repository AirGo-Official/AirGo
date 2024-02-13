<template>
  <div>
    <el-dialog v-model="state.isShowDialog" title="折扣码" width="80%">
      <el-form :model="shopStoreData.currentCoupon.value" label-width="120px">
        <el-form-item label="名称">
          <el-input v-model="shopStoreData.currentCoupon.value.name"></el-input>
        </el-form-item>
        <el-form-item label="折扣率">
          <el-col :span="6">
            <el-input-number v-model.number="shopStoreData.currentCoupon.value.discount_rate" :min="0" :max="1" :step="0.1"></el-input-number>
            <div style="color: #9b9da1">价格 = 原价 * ( 1 - 折扣率)</div>
          </el-col>

        </el-form-item>
        <el-form-item label="限制次数">
          <el-col :span="6">
            <el-input-number v-model.number="shopStoreData.currentCoupon.value.limit" :min="0" :step="1"></el-input-number>
          </el-col>
          <el-col :span="20"></el-col>
        </el-form-item>
        <el-form-item label="到期时间">
          <el-date-picker
              v-model="shopStoreData.currentCoupon.value.expired_at"
              type="datetime"
              placeholder="选择到期时间"
              size="default"
          />
        </el-form-item>
        <el-form-item label="关联商品">
          <el-transfer
              :data="shopStoreData.goodsList.value"
              v-model="shopStoreData.checkedGoodsIDs.value"
              :props="{key: 'id',label: 'subject'}"
              :titles="['全部商品', '选中商品']"
          />
        </el-form-item>
      </el-form>
      <template #footer>
            <span class="dialog-footer">
                <el-button @click="closeDialog">取消</el-button>
                <el-button type="primary" @click="onSubmit(shopStoreData.currentCoupon.value)" color="#FC3D08">
                    提交
                </el-button>
            </span>
      </template>
    </el-dialog>

  </div>
</template>

<script lang="ts" setup>
import {reactive} from "vue";
import {ElMessage} from "element-plus";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
import {storeToRefs} from "pinia";
import {useAdminShopStore} from "/@/stores/admin_logic/shopStore";

const shopStore = useAdminShopStore()
const shopStoreData = storeToRefs(shopStore)
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
const state = reactive({
  type: '',
  isShowDialog: false,
})

//打开弹窗
const openDialog = (type: string, row?: Coupon) => {
  //获取全部商品
  shopStore.getGoodsList() //获取全部商品
  state.isShowDialog = true
  state.type = type
  switch (type) {
    case "add":
      // shopStoreData.currentCoupon.value = {} as Coupon
      shopStoreData.currentCoupon.value.id = 0
      break
    case "edit":
      // shopStoreData.currentCoupon.value = JSON.parse(JSON.stringify(row))
      shopStoreData.currentCoupon.value = row
      shopStore.goodsIDsHandler()
      break
  }
}
//关闭弹窗
const closeDialog = () => {
  state.isShowDialog = false
}
//提交按钮
const onSubmit = (params: object) => {
  switch (state.type) {
    case "add":
      shopStore.newCoupon().then((res) => {
        emits('refresh')
      })
      break
    case "edit":
      shopStore.updateCoupon().then((res) => {
        emits('refresh')
      })
      break
  }
  closeDialog()
}
//子组件调用父组件
const emits = defineEmits(['refresh'])
//暴露变量
defineExpose({
  openDialog,
})
</script>

<style scoped>

</style>