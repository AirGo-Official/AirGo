<template>
  <el-dialog v-model="state.isShowDialog" :title="state.title" width="80%" destroy-on-close>
    <el-form :model="goodsManageData.currentGoods" label-width="80px" label-position="top">
      <el-form-item label="商品标题">
        <el-input v-model="goodsManageData.currentGoods.subject"/>
      </el-form-item>
      <el-form-item label="价格">
        <el-col :span="4">
          <el-input v-model="goodsManageData.currentGoods.total_amount"/>
        </el-col>
        <el-col :span="2" style="text-align: center">
          <span>-</span>
        </el-col>
        <el-col :span="18">
          <span class="text-gray-500">RMB</span>
        </el-col>
      </el-form-item>
      <el-form-item label="总流量">
        <el-col :span="4">
          <el-input-number v-model.number="goodsManageData.currentGoods.total_bandwidth" type="number"/>
        </el-col>
        <el-col :span="2" style="text-align: center">
          <span>-</span>
        </el-col>
        <el-col :span="18">
          <span class="text-gray-500">GB</span>
        </el-col>
      </el-form-item>
      <el-form-item label="有效期">
        <el-col :span="4">
          <el-input-number v-model.number="goodsManageData.currentGoods.expiration_date" type="number"/>
        </el-col>
        <el-col :span="2" style="text-align: center">
          <span>-</span>
        </el-col>
        <el-col :span="18">
          <span class="text-gray-500">天</span>
        </el-col>
      </el-form-item>
      <el-form-item label="流量重置日">
        <el-col :span="4">
          <el-input-number v-model.number="goodsManageData.currentGoods.reset_day" type="number"/>
        </el-col>
        <el-col :span="2" style="text-align: center">
          <span>-</span>
        </el-col>
        <el-col :span="18">
          <span class="text-gray-500">日</span>
        </el-col>
      </el-form-item>
      <el-form-item label="是否显示">
        <el-col :span="4">
          <el-switch v-model="goodsManageData.currentGoods.status" inline-prompt active-text="开启" inactive-text="关闭"
                     style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
        </el-col>
      </el-form-item>
      <el-form-item label="流量重置方式">
        <el-radio-group v-model="goodsManageData.currentGoods.traffic_reset_method" class="ml-4">
          <el-radio label="Stack" >叠加</el-radio>
          <el-radio label="NotStack" >不叠加</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="关联节点">
        <el-transfer
            :data="nodeManageData.nodes.node_list"
            v-model="goodsManageData.currentGoods.checked_nodes"
            :right-default-checked="goodsManageData.currentGoods.checked_nodes"
            :props="{
                  key: 'id',
                  label: 'remarks',
                  }"
            :titles="['全部节点', '选中节点']"
        />
      </el-form-item>
      <el-form-item label="描述">
        <v-md-editor v-model="goodsManageData.currentGoods.des" height="400px"></v-md-editor>
      </el-form-item>
    </el-form>
    <template #footer>
            <span class="dialog-footer">
                <el-button @click="state.isShowDialog = false">取消</el-button>
                <el-button type="primary" @click="onSubmit">
                    确认
                </el-button>
            </span>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import {storeToRefs} from 'pinia';
import {useShopStore} from "/@/stores/shopStore";
import {useNodeStore} from "/@/stores/nodeStore";
import {reactive} from "vue";

const shopStore = useShopStore()
const {goodsManageData} = storeToRefs(shopStore)

const nodeStore = useNodeStore();
const {nodeManageData} = storeToRefs(nodeStore);
// 定义子组件向父组件传值/事件
const emit = defineEmits(['refresh']);

//定义参数
const state = reactive({
  isShowDialog: false,
  type: '',
  title: '',
})

// 打开弹窗
const openDialog = (type: string, row?: any) => {
  if (type == 'add') {
    state.type = type
    state.title = "新建商品"
    state.isShowDialog = true
    goodsManageData.value.currentGoods.id = 0 //清空上次编辑的id，否则无法新建
  } else {
    state.type = type
    state.title = "修改商品"
    goodsManageData.value.currentGoods = row //将当前row写入pinia
    if (goodsManageData.value.currentGoods.checked_nodes === null || goodsManageData.value.currentGoods.checked_nodes === undefined) {
      goodsManageData.value.currentGoods.checked_nodes = [] //剔除null,否则ts报错
    }
    goodsManageData.value.currentGoods.nodes = [] //清空nodes
    state.isShowDialog = true
  }
}
// 关闭弹窗
const closeDialog = () => {
  state.isShowDialog = false
};

//确认提交
function onSubmit() {
  if (state.type === 'add') {
    shopStore.newGoods(goodsManageData.value.currentGoods)
    setTimeout(() => {
      emit('refresh');
    }, 500);
  } else {
    shopStore.updateGoods(goodsManageData.value.currentGoods)
    setTimeout(() => {
      emit('refresh');
    }, 500);
  }
  closeDialog();
}

// 暴露变量
defineExpose({
  openDialog,
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

</style>
  