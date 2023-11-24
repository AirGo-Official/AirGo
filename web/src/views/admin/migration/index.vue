<template>
  <div class="container layout-padding">
    <el-card shadow="hover" class="layout-padding-auto">
      <el-row :gutter="10" style="width: 768px">
        <el-col :span="4">
          <el-input size="default" placeholder="请输入名称"
                    style="max-width: 180px"></el-input>
        </el-col>
        <el-col :span="3">
          <el-button size="default" type="primary" class="ml10">
            <el-icon>
              <ele-Search/>
            </el-icon>
            查询
          </el-button>
        </el-col>
        <el-col :span="4">
          <el-button size="default" type="success" class="ml10" @click="openDialog('新建')">
            <el-icon>
              <ele-FolderAdd/>
            </el-icon>
            新增
          </el-button>
        </el-col>
      </el-row>
      <el-row style="margin-top: 20px;margin-bottom: 20px">
        <el-button type="danger">
          迁移之前请做好数据备份。由于面板之间数据不兼容，迁移仅保留用户最基本的账号email和uuid数据，并将用户密码初始为123456，请引导用户及时修改密码
        </el-button>
      </el-row>
      <el-row v-loading="state.loading">
        <el-card shadow="hover" class="layout-padding-auto">
          <el-form v-model="state.migrationParams" label-position="top">
            <el-form-item label="面板类型">
              <el-select v-model="state.migrationParams.panel_type" placeholder="Select">
                <el-option
                    v-for="item in state.panels"
                    :key="item"
                    :label="item"
                    :value="item"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="数据库地址">
              <el-input v-model="state.migrationParams.db_address"></el-input>
            </el-form-item>
            <el-form-item label="数据库端口">
              <el-input-number v-model="state.migrationParams.db_port"></el-input-number>
            </el-form-item>
            <el-form-item label="数据库名">
              <el-input v-model="state.migrationParams.db_name"></el-input>
            </el-form-item>
            <el-form-item label="数据库用户">
              <el-input v-model="state.migrationParams.db_username"></el-input>
            </el-form-item>
            <el-form-item label="数据库密码">
              <el-input v-model="state.migrationParams.db_password"></el-input>
            </el-form-item>
            <el-form-item >
              <el-button color="blue" @click="onSubmit">提交</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-row>
      <el-row>
<!--        <el-button  color="blue" @click="onStop">终止</el-button>-->
        <el-card shadow="hover" style="width: 100%;margin-top: 20px">
          迁移结果: {{state.migrationResult}}
        </el-card>
      </el-row>


    </el-card>
    <Dialog ref="DialogRef"></Dialog>
  </div>
</template>

<script lang="ts" setup>

//定义参数
import {defineAsyncComponent, reactive, ref} from "vue";
import {ElMessageBox} from "element-plus";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";

const Dialog = defineAsyncComponent(() => import('/@/views/admin/migration/dialog.vue'))
const DialogRef = ref()
const apiStore=useApiStore()

const state = reactive({
  loading: false,
  html: "",
  panels: ["v2board", "sspanel"],
  migrationParams: {
    "panel_type": "",
    "db_address": "",
    "db_port": 3306,
    "db_username": "",
    "db_password": "",
    "db_name": "",
  },
  migrationResult:"等待迁移",

})

//打开弹窗
function openDialog(title: string, row?: any) {
  DialogRef.value.openDialog(title, row)
}
const onSubmit=()=>{
  ElMessageBox.confirm(`请做好数据备份，并填写正确的数据库信息，是否继续?`, '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning',
  })
      .then(() => {
        state.loading=true
        console.log(1)
        request(apiStore.api.migration_fromOther,state.migrationParams).then((res)=>{
          state.migrationResult=res.data
          state.loading=false
        }).catch(()=>{
          state.loading=false
        })
      })
      .catch(() => {
      });
}

</script>

<style scoped lang="scss">
.container {
  :deep(.el-card__body) {
    display: flex;
    flex-direction: column;
    flex: 1;
    overflow: auto;

    .el-table {
      flex: 1;
    }
  }
}

</style>