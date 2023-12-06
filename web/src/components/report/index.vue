<template>
  <div class="container report_layout-padding">
    <div>
<!--      测试数据源-->
<!--      <el-select v-model="reportStoreData.checkedDbInfo.value.table_name" placeholder="选择数据源" @change="getColumn">-->
<!--        <el-option-->
<!--            v-for="item in reportStoreData.allDbTables.value"-->
<!--            :key="item.en_name"-->
<!--            :label="item.cn_name"-->
<!--            :value="item.en_name"-->
<!--        >-->
<!--          <span style="float: left">{{ item.en_name }}</span>-->
<!--          <span style="float: right;color: var(&#45;&#45;el-text-color-secondary);font-size: 13px;">{{ item.cn_name }}</span-->
<!--          >-->
<!--        </el-option>-->
<!--      </el-select>-->
      <el-button type="primary" @click="addCondition">新增条件</el-button>
    </div>
    <div>
      <el-table :data="reportStoreData.reportParams.value.field_params_list" height="100%" style="width: 100%;flex: 1;" stripe>
        <el-table-column align="left" type="index" label="序列" width="60"/>
        <el-table-column align="left" prop="operator" label="运算符" width="100">
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
        <el-table-column align="left" prop="field" label="字段" width="160">
          <template #default="{row}">
            <el-select  v-model="row.field" class="m-2" placeholder="选择字段">
              <el-option
                  v-for="(v,k) in reportStoreData.fieldData.value.field_list"
                  :key="k"
                  :label="v"
                  :value="v"
              />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column align="left" prop="field_chinese_name " label="字段中文名" width="160">
          <template #default="{row}">
            <el-text>{{ row.field_chinese_name }}</el-text>
          </template>
        </el-table-column>
        <el-table-column align="left" prop="field_type" label="字段类型" width="100">
          <template #default="{row}">
            <el-text>{{ row.field_type }}</el-text>
          </template>
        </el-table-column>
        <el-table-column align="left" prop="condition" label="搜索条件" width="100">
          <template #default="{row}">
            <el-select v-if="row.field_type === 'string'" v-model="row.condition" class="m-2" placeholder="搜索条件">
              <el-option
                  v-for="(v,k) in stringConditionList"
                  :key="k"
                  :label="v"
                  :value="v"
              />
            </el-select>
            <el-select v-if="row.field_type === 'num'" v-model="row.condition" class="m-2" placeholder="搜索条件">
              <el-option
                  v-for="(v,k) in numConditionList"
                  :key="k"
                  :label="v"
                  :value="v"
              />
            </el-select>
            <el-select v-if="row.field_type === 'date'" v-model="row.condition" class="m-2" placeholder="搜索条件">
              <el-option
                  v-for="(v,k) in dateConditionList"
                  :key="k"
                  :label="v"
                  :value="v"
              />
            </el-select>
            <el-select v-if="row.field_type === 'bool'" v-model="row.condition" class="m-2" placeholder="搜索条件">
              <el-option
                  v-for="(v,k) in boolConditionList"
                  :key="k"
                  :label="v"
                  :value="v"
              />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column align="left" prop="conditionValue" label="条件值" width="200px">
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
        <el-table-column align="left" label="操作">
          <template #default="{row}">
            <el-button type="primary" @click="deleteCurrrentCondition(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
<!--      <el-pagination-->
<!--          background-->
<!--          class="mt15"-->
<!--          layout="total, sizes, prev, pager, next, jumper"-->
<!--          :page-sizes="[10, 30, 50, 100, 200]"-->
<!--          v-model:current-page="reportStoreData.reportParams.value.pagination.page_num"-->
<!--          v-model:page-size="reportStoreData.reportParams.value.pagination.page_size"-->
<!--          :total="reportStoreData.reportData.value.total"-->
<!--          @size-change="onHandleSizeChange"-->
<!--          @current-change="onHandleCurrentChange"-->
<!--      >-->
<!--      </el-pagination>-->
    </div>
    <div style="margin-top: 20px">
      <el-button v-if="reportStoreData.reportParams.value.field_params_list.length!==0" @click="onSubmit()" type="primary">查询</el-button>
<!--      <el-button v-if="reportStoreData.reportParams.value.field_params_list.length!==0" @click="onFind()" type="primary">查询</el-button>-->
    </div>
<!--测试显示结果-->
<!--    <div>-->
<!--      <el-table :data="reportStoreData.reportData.value.data" stripe>-->
<!--        <el-table-column :label="item" :property="item" v-for="item in reportStoreData.fieldData.value.field_list" :key="item" align="center">-->
<!--          <template #default="scope">-->
<!--            <span>{{scope.row[item]}}</span>-->
<!--          </template>-->
<!--        </el-table-column>-->
<!--      </el-table>-->
<!--    </div>-->

  </div>
</template>

<script lang="ts" setup>


import {onMounted, reactive, watch} from "vue";
import {useReportStore} from "/@/stores/reportStore"
import {storeToRefs} from "pinia";

const reportStore = useReportStore()
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
  onFind()
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
  reportStoreData.reportParams.value.pagination.page_num = val;
  onFind()
};
//提交
const onSubmit = (params?: object) => {
  if (reportStoreData.checkedDbInfo.value.table_name === '') {
    return
  }
  //调用父组件 getReportData()方法
  emits('getReportData')
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