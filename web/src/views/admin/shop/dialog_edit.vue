<template>
  <el-dialog v-model="state.isShowDialog" :title="state.title" width="80%" destroy-on-close>
    <el-form :model="shopStoreData.currentGoods.value" label-position="right" label-width="120">
      <el-form-item :label="$t('message.adminShop.Goods.goods_type')">
        <el-radio-group v-model="shopStoreData.currentGoods.value.goods_type">
          <el-radio :label="constantStore.GOODS_TYPE_GENERAL">{{ $t("message.constant.GOODS_TYPE_GENERAL") }}</el-radio>
          <el-radio :label="constantStore.GOODS_TYPE_SUBSCRIBE">{{ $t("message.constant.GOODS_TYPE_SUBSCRIBE") }}
          </el-radio>
          <el-radio :label="constantStore.GOODS_TYPE_RECHARGE">{{ $t("message.constant.GOODS_TYPE_RECHARGE") }}
          </el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item :label="$t('message.adminShop.Goods.subject')">
        <el-input v-model="shopStoreData.currentGoods.value.subject"/>
      </el-form-item>
      <el-form-item :label="$t('message.adminShop.Goods.cover_image')">
        <el-image :src="shopStoreData.currentGoods.value.cover_image" style="height: 100px">
          <template #error>
            <div class="image-slot">
              <i class="ri-signal-wifi-error-line"></i>
            </div>
          </template>
        </el-image>
        <el-input v-model="shopStoreData.currentGoods.value.cover_image"/>
      </el-form-item>
      <el-form-item :label="$t('message.adminShop.Goods.price')">
        <el-input v-model="shopStoreData.currentGoods.value.price"/>
      </el-form-item>
      <div v-if="shopStoreData.currentGoods.value.goods_type === constantStore.GOODS_TYPE_SUBSCRIBE">
        <el-form-item :label="$t('message.adminShop.Goods.price_3_month')">
          <el-input v-model="shopStoreData.currentGoods.value.price_3_month"/>
        </el-form-item>
        <el-form-item :label="$t('message.adminShop.Goods.price_6_month')">
          <el-input v-model="shopStoreData.currentGoods.value.price_6_month"/>
        </el-form-item>
        <el-form-item :label="$t('message.adminShop.Goods.price_12_month')">
          <el-input v-model="shopStoreData.currentGoods.value.price_12_month"/>
        </el-form-item>
        <el-form-item :label="$t('message.adminShop.Goods.price_unlimited_duration')">
          <el-input v-model="shopStoreData.currentGoods.value.price_unlimited_duration"/>
        </el-form-item>
      </div>
      <el-row class="mt15 mb15">
        <el-col :span="12">
          <el-form-item :label="$t('message.adminShop.Goods.quota')">
            <el-input-number v-model="shopStoreData.currentGoods.value.quota" :min="0" :step="1"/>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('message.adminShop.Goods.stock')">
            <el-input-number v-model="shopStoreData.currentGoods.value.stock" :min="0" :step="1"/>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row class="mt15 mb15">
        <el-col :span="8">
          <el-form-item :label="$t('message.adminShop.Goods.is_show')">
            <el-switch v-model="shopStoreData.currentGoods.value.is_show" inline-prompt
                       :active-text="$t('message.common.enable')"
                       :inactive-text="$t('message.common.disable')"
                       style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>

          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item :label="$t('message.adminShop.Goods.is_sale')">
            <el-switch v-model="shopStoreData.currentGoods.value.is_sale" inline-prompt
                       :active-text="$t('message.common.enable')"
                       :inactive-text="$t('message.common.disable')"
                       style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item :label="$t('message.adminShop.Goods.is_renew')">
            <el-switch v-model="shopStoreData.currentGoods.value.is_renew" inline-prompt
                       :active-text="$t('message.common.enable')"
                       :inactive-text="$t('message.common.disable')"
                       style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
          </el-form-item>
        </el-col>
      </el-row>

      <!--      发货参数开始-->
      <div v-if="shopStoreData.currentGoods.value.goods_type === constantStore.GOODS_TYPE_GENERAL">
        <el-form-item :label="$t('message.adminShop.Goods.deliver_type')">
          <el-radio-group v-model="shopStoreData.currentGoods.value.deliver_type">
            <el-radio label="none">{{ $t("message.constant.DELIVER_TYPE_NONE") }}</el-radio>
            <el-radio label="manual">{{ $t("message.constant.DELIVER_TYPE_MANUAL") }}</el-radio>
            <el-radio label="auto">{{ $t("message.constant.DELIVER_TYPE_AUTO") }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="$t('message.adminShop.Goods.deliver_text')"
                      v-if="shopStoreData.currentGoods.value.deliver_type === constantStore.DELIVER_TYPE_AUTO">
          <v-md-editor v-model="shopStoreData.currentGoods.value.deliver_text" height="400px"></v-md-editor>
        </el-form-item>
      </div>
      <!--      发货参数结束-->
      <!--      订阅商品开始-->
      <div v-if="shopStoreData.currentGoods.value.goods_type === constantStore.GOODS_TYPE_SUBSCRIBE">
        <el-form-item :label="$t('message.adminShop.Goods.enable_traffic_reset')">
          <el-switch v-model="shopStoreData.currentGoods.value.enable_traffic_reset" inline-prompt
                     :active-text="$t('message.common.enable')"
                     :inactive-text="$t('message.common.disable')"
                     style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
        </el-form-item>
        <el-form-item :label="$t('message.adminShop.Goods.total_bandwidth')">
          <el-col :span="4">
            <el-input-number v-model.number="shopStoreData.currentGoods.value.total_bandwidth"/>
          </el-col>
          <el-col :span="2" style="text-align: center">
            <span>-</span>
          </el-col>
          <el-col :span="18">
            <span class="text-gray-500">GB</span>
          </el-col>
        </el-form-item>
        <el-form-item :label="$t('message.adminShop.Goods.nodes')">
          <el-tree ref="nodes_tree_ref" node-key="id"
                   :data="nodeManageData.nodeList.value.data"
                   :props="{label:'remarks'}"
                   :default-checked-keys="shopStoreData.checkedNodeIDs.value"
                   show-checkbox class="menu-data-tree"/>
        </el-form-item>
      </div>
      <!--      订阅商品结束-->
      <!--      充值商品开始-->
      <div v-if="shopStoreData.currentGoods.value.goods_type === constantStore.GOODS_TYPE_RECHARGE">
        <el-form-item :label="$t('message.adminShop.Goods.recharge_amount')">
          <el-input v-model="shopStoreData.currentGoods.value.recharge_amount" placeholder="100.00"></el-input>
        </el-form-item>
      </div>
      <!--      充值商品结束-->
      <div class="mt15 mb15">
        <el-form-item :label="$t('message.adminShop.Goods.des')">
          <v-md-editor  v-model="shopStoreData.currentGoods.value.des" height="400px"></v-md-editor>
        </el-form-item>
      </div>

    </el-form>
    <template #footer>
            <span class="dialog-footer">
                <el-button @click="state.isShowDialog = false">{{ $t("message.common.button_cancel") }}</el-button>
                <el-button type="primary" @click="onSubmit">
                    {{ $t("message.common.button_confirm") }}
                </el-button>
            </span>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import {storeToRefs} from "pinia";
import {useAdminNodeStore} from "/@/stores/admin_logic/nodeStore";
import {reactive, ref} from "vue";
import {useAdminShopStore} from "/@/stores/admin_logic/shopStore";
import {useConstantStore} from "/@/stores/constantStore";
import {useI18n} from "vue-i18n";
import {ElMessage} from "element-plus";

const shopStore = useAdminShopStore();
const shopStoreData = storeToRefs(shopStore);
const nodeStore = useAdminNodeStore();
const nodeManageData = storeToRefs(nodeStore);
const constantStore = useConstantStore();
const nodes_tree_ref = ref();
const {t} = useI18n();
// 定义子组件向父组件传值/事件
const emit = defineEmits(["refresh"]);

//定义参数
const state = reactive({
  isShowDialog: false,
  type: "",
  title: "",
  queryParams: {
    table_name: "node",
    field_params_list: [
      {field: "id", field_chinese_name: "", field_type: "", condition: "<>", condition_value: "", operator: ""}
      // {field: 'created_at', field_chinese_name: '', field_type: '', condition: '<', condition_value: "", operator: 'AND',}
    ] as FieldParams[],
    pagination: {page_num: 1, page_size: 9999, order_by: "node_order"} as Pagination//设为9999，理论能获取全部节点，暂时取消详细的分页设置
  } as QueryParams,
});
//查询节点
const getNodeList = () => {
  nodeStore.getNodeList(state.queryParams);
}

// 打开弹窗
const openDialog = (type: string, row?: any) => {
  state.type = type;
  getNodeList()
  if (type == "add") {
    state.title = t("message.adminShop.addGoods");
    state.isShowDialog = true;
    shopStoreData.currentGoods.value.id = 0; //清空上次编辑的id，否则无法新建
  } else {
    state.type = type;
    state.title = t("message.adminShop.modifyGoods");
    shopStoreData.currentGoods.value = row;
    shopStore.nodeIDsHandler();
    state.isShowDialog = true;
  }
};
// 关闭弹窗
const closeDialog = () => {
  state.isShowDialog = false;
};

//确认提交
function onSubmit() {
  //处理商品的关联节点
  if (shopStoreData.currentGoods.value.goods_type === constantStore.GOODS_TYPE_SUBSCRIBE) {
    shopStoreData.checkedNodeIDs.value = [...nodes_tree_ref.value.getCheckedKeys()];
  }
  if (state.type === "add") {
    shopStore.newGoods().then((res) => {
      ElMessage.success(res.msg);
      emit("refresh");
    });
  } else {
    shopStore.updateGoods().then((res) => {
      ElMessage.success(res.msg);
      emit("refresh");
    });
  }
  closeDialog();
}

// 暴露变量
defineExpose({
  openDialog
});

</script>


<style>
/* 定义两边的el-transfer-panel大小的方法,直接设置是没有用的,需要去掉scoped即可。才能成功覆盖原生的样式 */
.el-transfer-panel {
  width: 400px;
  height: 600px;
}

.el-transfer-panel__body {
  height: 600px;
}

.el-transfer-panel__list {
  height: 550px;
}

/*穿梭框内部展示列表的高宽度*/
/*:deep(.el-transfer-panel__list.is-filterable){*/
/*  width:280px;*/
/*  height:500px;*/
/*}*/
.menu-data-tree {
  width: 100%;
  border: 1px solid var(--el-border-color);
  border-radius: var(--el-input-border-radius, var(--el-border-radius-base));
  padding: 5px;
}

</style>
  