<template>
  <div class="container report_layout-padding">
    <div>
      <el-button type="primary" @click="addCondition">{{$t('message.report.addCondition')}}</el-button>
    </div>
    <div>
      <el-table :data="reportStoreData.reportParams.value.field_params_list" height="100%" style="width: 100%;flex: 1;" stripe>
        <el-table-column align="left" type="index" :label="$t('message.report.index')" width="60"/>
        <el-table-column align="left" prop="operator" :label="$t('message.report.operator')" width="100">
          <template #default="scope">
            <el-select v-model="scope.row.operator" class="m-2" v-if="scope.$index > 0">
              <el-option
                  v-for="(v,k) in operatorList"
                  :key="k"
                  :label="v"
                  :value="v"
              />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column align="left" prop="field" :label="$t('message.report.field')" width="160">
          <template #default="{row}">
            <el-select  v-model="row.field" class="m-2">
              <el-option
                  v-for="(v,k) in reportStoreData.fieldData.value.field_list"
                  :key="k"
                  :label="v"
                  :value="v"
              />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column align="left" prop="field_chinese_name " :label="$t('message.report.field_chinese_name')" width="160">
          <template #default="{row}">
            <el-text>{{ row.field_chinese_name }}</el-text>
          </template>
        </el-table-column>
        <el-table-column align="left" prop="field_type" :label="$t('message.report.field_type')" width="100">
          <template #default="{row}">
            <el-text>{{ row.field_type }}</el-text>
          </template>
        </el-table-column>
        <el-table-column align="left" prop="condition" :label="$t('message.report.condition')" width="100">
          <template #default="{row}">
            <el-select v-if="row.field_type === 'string'" v-model="row.condition" class="m-2">
              <el-option
                  v-for="(v,k) in stringConditionList"
                  :key="k"
                  :label="v"
                  :value="v"
              />
            </el-select>
            <el-select v-if="row.field_type === 'num'" v-model="row.condition" class="m-2">
              <el-option
                  v-for="(v,k) in numConditionList"
                  :key="k"
                  :label="v"
                  :value="v"
              />
            </el-select>
            <el-select v-if="row.field_type === 'date'" v-model="row.condition" class="m-2">
              <el-option
                  v-for="(v,k) in dateConditionList"
                  :key="k"
                  :label="v"
                  :value="v"
              />
            </el-select>
            <el-select v-if="row.field_type === 'bool'" v-model="row.condition" class="m-2">
              <el-option
                  v-for="(v,k) in boolConditionList"
                  :key="k"
                  :label="v"
                  :value="v"
              />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column align="left" prop="conditionValue" :label="$t('message.report.conditionValue')" width="200px">
          <template #default="{row}">
            <el-date-picker
                v-if="row.field_type==='date'"
                v-model="row.condition_value"
                type="datetime"
                value-format="YYYY-MM-DD HH:mm:ss"
            />
            <el-switch v-if="row.field_type==='bool'" v-model="row.condition_value" inline-prompt active-text="true"
                       inactive-text="false"
                       style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
            <el-input v-if="row.field_type === 'num'" v-model="row.condition_value" type="number"></el-input>
            <el-input v-if="row.field_type === 'string'" v-model="row.condition_value"></el-input>
          </template>
        </el-table-column>
        <el-table-column align="left" :label="$t('message.common.operate')">
          <template #default="{row}">
            <el-button type="primary" @click="deleteCurrrentCondition(row)">{{$t('message.common.delete')}}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <div style="margin-top: 20px">
      <el-button v-if="reportStoreData.reportParams.value.field_params_list.length!==0" @click="onSubmit()" type="primary">查询</el-button>
    </div>


  </div>
</template>

<script lang="ts" setup>


import {onMounted, reactive, watch} from "vue";
import {useAdminReportStore} from "/@/stores/admin_logic/reportStore"
import {storeToRefs} from "pinia";

const reportStore = useAdminReportStore()
const reportStoreData = storeToRefs(reportStore)

const state=reactive({
  data:[]
})

//搜索条件
const operatorList = ["AND", "OR"]
const stringConditionList = ["=", "<>", "like"]
const numConditionList = ['<', '>', "=", "<>"]
const boolConditionList = ["=", "<>"]
const dateConditionList = ['<', '>', "=", "<>"]
//子传父
const emits = defineEmits(['getReportData'])

//删除当前条件
const deleteCurrrentCondition = (row: any) => {
  // console.log("当前所有条件:",reportParams.value.field_params_list)
  reportStoreData.reportParams.value.field_params_list = reportStoreData.reportParams.value.field_params_list.filter(item => item !== row);
}
//增加新条件
const addCondition = () => {
  if (reportStoreData.reportParams.value.field_params_list.length===0){
    reportStoreData.reportParams.value.field_params_list.push({
      field: '',
      field_chinese_name: '',
      field_type: '',
      condition: '=',
      condition_value: '',
      operator: '',
    } as FieldParams)
  } else {
    reportStoreData.reportParams.value.field_params_list.push({
      field: '',
      field_chinese_name: '',
      field_type: '',
      condition: '=',
      condition_value: '',
      operator: 'AND',
    } as FieldParams)
  }

};

//
const onFind=()=>{
  reportStoreData.reportParams.value.table_name=reportStoreData.checkedDbInfo.value.table_name
  reportStore.getReport()
}
// 分页改变
const onHandleSizeChange = (val: number) => {
  reportStoreData.reportParams.value.pagination.page_size = val;
  onSubmit()
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
  reportStoreData.reportParams.value.pagination.page_num = val;
  onSubmit()
};
//提交
const onSubmit = (params?: object) => {
  if (reportStoreData.checkedDbInfo.value.table_name === '') {
    return
  }
  //调用父组件 getReportData()方法
  emits('getReportData',reportStoreData.reportParams.value)
}

//监听
watch(
    () => reportStoreData.reportParams.value,   //数据源有变化就开始处理
    () => {
      // console.log("reportParams.value.field_params_list:",reportStoreData.reportParams.value.field_params_list)
      reportStoreData.reportParams.value.field_params_list.forEach((value,index,array) => {
        value.field_chinese_name = reportStoreData.fieldData.value.field_chinese_name_list[value.field]
        value.field_type = reportStoreData.fieldData.value.field_type_list[value.field]
      })
    },
    {
      deep: true,
    }
);
// onMounted(() => {
//
// });

//打开时加载数据
const openReportComponent = (params: string) => {
  //设置需要操作的库表
  reportStoreData.checkedDbInfo.value.table_name = params
  reportStoreData.reportParams.value.table_name = params
  //加载字段信息
  reportStore.getColumn()
}

//暴露变量
defineExpose({
  openReportComponent,
});


</script>

<style scoped lang="scss">
.container {
  :deep(.el-card__body) {
    display: flex;
    flex-direction: column;
    flex: 1;
    overflow: auto;
  }

}

.home-card-item {
  font-size: 16px;
  width: 100%;
  height: 100%;
  border-radius: 4px;
  transition: all ease 0.3s;
  padding: 10px;
  overflow: hidden;
  background: var(--el-color-white);
  color: var(--el-text-color-primary);
  border: 1px solid var(--next-border-color-light);
}

.report_layout-padding {
  padding: 15px;
}

//.el-card {
//  background-image: url("../../assets/bgc/3.png");
//  background-repeat: no-repeat;
//}
</style>