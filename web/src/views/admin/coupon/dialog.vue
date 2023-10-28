<template>
  <div>
    <el-dialog v-model="state.isShowDialog" title="折扣码" width="80%">
      <el-form :model="couponStoreData.coupon.value" label-width="120px">
        <el-form-item label="名称">
          <el-input v-model="couponStoreData.coupon.value.name"></el-input>
        </el-form-item>
        <el-form-item label="折扣率">
          <el-col :span="6">
            <el-input v-model.number="couponStoreData.coupon.value.discount_rate" type="number"></el-input>
            <div style="color: #9b9da1">实际价格=原价*(1-折扣率)</div>
          </el-col>

        </el-form-item>
        <el-form-item label="限制次数">
          <el-col :span="6">
            <el-input v-model.number="couponStoreData.coupon.value.limit" type="number"></el-input>
          </el-col>
          <el-col :span="20"></el-col>
        </el-form-item>
        <el-form-item label="到期时间">
          <el-date-picker
              v-model="couponStoreData.coupon.value.expired_at"
              type="datetime"
              placeholder="选择到期时间"
              size="default"
          />
        </el-form-item>
        <el-form-item label="关联商品">
          <el-transfer
              :data="shopStoreData.goodsList.value"
              v-model="couponStoreData.coupon.value.checked_goods"
              :right-default-checked="couponStoreData.coupon.value.checked_goods"
              :props="{key: 'id',label: 'subject'}"
              :titles="['全部商品', '选中商品']"
          />
        </el-form-item>
      </el-form>
      <template #footer>
            <span class="dialog-footer">
                <el-button @click="closeDialog">取消</el-button>
                <el-button type="primary" @click="onSubmit(couponStoreData.coupon.value)" color="#FC3D08">
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
import {useShopStore} from "/@/stores/shopStore";
import {useCouponStore} from "/@/stores/couponStore";
const shopStore = useShopStore()
const shopStoreData = storeToRefs(shopStore)

const couponStore = useCouponStore()
const couponStoreData = storeToRefs(couponStore)

const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
const state = reactive({
  type: '',
  isShowDialog: false,
})

//打开弹窗
const openDialog = (type: string, data?: Coupon) => {
  //获取全部商品
  shopStore.getAllGoods() //获取全部商品
  state.isShowDialog = true
  state.type = type
  switch (type) {
    case "add":
      couponStoreData.coupon.value = {} as Coupon
      break
    case "edit":
      couponStoreData.coupon.value = JSON.parse(JSON.stringify(data))
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
      request(apiStoreData.api.value.coupon_newCoupon, params).then((res) => {
        ElMessage.success(res.msg)
        emits('refresh')
      })
      break
    case "edit":
      request(apiStoreData.api.value.coupon_updateCoupon, params).then((res) => {
        ElMessage.success(res.msg)
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