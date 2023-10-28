<template>
  <div style="padding: 15px;">
    <el-card>
      <el-tabs stretch style="height: 100%" class="demo-tabs">
        <el-tab-pane label="登录/注册">
          <el-form :model="serverConfig" label-width="120px">
            <el-form-item label="是否开启注册">
              <el-switch v-model="serverConfig.system.enable_register" inline-prompt active-text="开启"
                         inactive-text="关闭"
                         style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
            </el-form-item>
            <el-form-item label="注册邮箱验证码">
              <el-switch v-model="serverConfig.system.enable_email_code" inline-prompt active-text="开启"
                         inactive-text="关闭"
                         style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
            </el-form-item>
            <!--            <el-form-item label="登录邮箱验证码">-->
            <!--              <el-switch v-model="serverConfig.system.enable_login_email_code" inline-prompt active-text="开启"-->
            <!--                         inactive-text="关闭"-->
            <!--                         style="&#45;&#45;el-switch-on-color: #13ce66; &#45;&#45;el-switch-off-color: #ff4949"></el-switch>-->
            <!--              <el-tag type="info" style="margin-left: 10px">最好不要开启，配置不正确你自己都登录不上</el-tag>-->
            <!--            </el-form-item>-->


            <el-divider></el-divider>
            <el-form-item label="IP限流">
              <el-col :span="4">
                <el-input v-model.number="serverConfig.rate_limit_params.ip_role_param" type="number"/>
              </el-col>
              <el-col :span="2" style="text-align: center">
                <span>-</span>
              </el-col>
              <el-col :span="18">
                <span class="text-gray-500">请求/分钟</span>
              </el-col>
            </el-form-item>

            <el-form-item label="用户限流">
              <el-col :span="4">
                <el-input v-model.number="serverConfig.rate_limit_params.visit_param" type="number"/>
              </el-col>
              <el-col :span="2" style="text-align: center">
                <span>-</span>
              </el-col>
              <el-col :span="18">
                <span class="text-gray-500">请求/分钟</span>
              </el-col>
            </el-form-item>
            <el-divider></el-divider>
            <el-form-item label="通信密钥">
              <el-input v-model="serverConfig.system.muKey" placeholder="务必前后端保持一致！"/>
              <div style="color: #9b9da1;display:block">XrayR等配置的密钥</div>
            </el-form-item>
            <el-form-item label="订阅名称">
              <el-input v-model="serverConfig.system.sub_name"/>
              <div style="color: #9b9da1;display:block">更新订阅时显示的名字</div>
            </el-form-item>
            <el-form-item label="AirGo后端地址">
              <el-input v-model="serverConfig.system.backend_url"/>
              <div style="color: #9b9da1">
                该地址与更新订阅、支付回调有关，请认真填写。前后分离时一般和前端.env中的VITE_API_URL保持一致即可；前后不分离时填公网可访问的后端地址。例如：http://abc.com:8899
              </div>
            </el-form-item>
            <el-divider></el-divider>
            <el-form-item label="新注册分配套餐">
              <el-select v-model="serverConfig.system.default_goods" placeholder="选择套餐" style="width: 30%">
                <el-option
                    v-for="item in goodsList"
                    :key="item.id"
                    :label="item.subject"
                    :value="item.subject"
                />
              </el-select>
            </el-form-item>

            <el-form-item label="邀请返利">
              <el-switch v-model="serverConfig.system.enabled_rebate" inline-prompt active-text="开启"
                         inactive-text="关闭"
                         style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
            </el-form-item>
            <el-form-item label="返利率">
              <el-input v-model.number="serverConfig.system.rebate_rate" type="number"></el-input>
              <div style="color: #9b9da1">(范围0~1)邀请收入=其他用户套餐实际支付价格*返利率</div>
            </el-form-item>
            <el-form-item label="旧套餐抵扣">
              <el-switch v-model="serverConfig.system.enabled_deduction" inline-prompt active-text="开启"
                         inactive-text="关闭"
                         style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
            </el-form-item>
            <el-form-item label="旧套餐抵扣阈值">
              <el-input v-model.number="serverConfig.system.deduction_threshold" type="number"></el-input>
              <div style="color: #9b9da1">
                (范围0~1)原套餐100G，用50G，剩余比例0.5，小于该阈值，则不会进行抵扣；原套餐实际付款为0也不抵扣
              </div>
            </el-form-item>

            <el-divider></el-divider>
            <el-form-item>
              <el-button @click="onSubmit" type="primary">保存</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>


        <el-tab-pane label="支付">
          <div>
            <el-button size="default" type="primary" class="ml10" @click="openPayDialog('add')">
              <el-icon>
                <ele-FolderAdd/>
              </el-icon>
              新增支付
            </el-button>
          </div>
          <div>
            <el-table :data="payStoreData.payList.value" stripe style="width: 100%;flex: 1;">
              <el-table-column type="index" label="序号" fixed show-overflow-tooltip width="60px"/>
              <!--              <el-table-column prop="id" label="ID" fixed show-overflow-tooltip width="60px"/>-->
              <el-table-column prop="name" label="别名" fixed show-overflow-tooltip width="120px"/>
              <el-table-column prop="pay_type" label="类型" show-overflow-tooltip fixed width="80px"/>
              <el-table-column prop="pay_logo_url" label="logo" fixed show-overflow-tooltip width="120px">
                <template #default="{row}">
                  <el-image :src="row.pay_logo_url" style="width: 40px;height: 40px"></el-image>
                </template>
              </el-table-column>
              <el-table-column prop="status" label="状态" fixed show-overflow-tooltip width="80px">
                <template #default="{row}">
                  <el-button v-if="row.status" type="warning">启用</el-button>
                  <el-button v-else type="info">禁用</el-button>
                </template>
              </el-table-column>

              <el-table-column label="操作">
                <template #default="scope">
                  <el-button text @click="openPayDialog('edit',scope.row)" type="primary">编辑</el-button>
                  <el-button text @click="deletePay(scope.row)" type="primary">删除</el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>

        </el-tab-pane>
        <el-tab-pane label="邮件">
          <el-form :model="serverConfig" label-width="100px">
            <el-form-item label="服务器地址">
              <el-input v-model="serverConfig.email.email_host" placeholder="mail.example.com"/>
            </el-form-item>
            <el-form-item label="端口">
              <el-input v-model.number="serverConfig.email.email_port" type="number"/>
            </el-form-item>
            <el-form-item label="邮箱账户名">
              <el-input v-model="serverConfig.email.email_from" placeholder="10010@qq.com"/>
            </el-form-item>
            <el-form-item label="邮箱别名">
              <el-input v-model="serverConfig.email.email_from_alias" placeholder="10010@foxmail.com"/>
              <div style="color: #9b9da1">*例如：qq邮箱可以设置foxmil别名。发送邮件时优先显示别名。无特殊情况可忽略本项</div>
            </el-form-item>
            <el-form-item label="账户昵称">
              <el-input v-model="serverConfig.email.email_nickname" placeholder="吊炸天机场管理员"/>
            </el-form-item>
            <el-form-item label="密码">
              <el-input v-model="serverConfig.email.email_secret" type="password"/>
            </el-form-item>
            <el-form-item label="默认邮件主题">
              <el-input v-model="serverConfig.email.email_subject"/>
            </el-form-item>
            <el-form-item label="邮件内容格式">
              <el-input v-model="serverConfig.email.email_content" type="textarea" autosize/>
              <el-text style="color: #9b9da1">*自定义邮件验证码内容样式，支持HTML，`emailcode`为验证码字段，不可删除！
              </el-text>
            </el-form-item>
            <el-divider></el-divider>
            <el-form-item>
              <el-button @click="onSubmit" type="primary">保存</el-button>
              <el-button @click="onTestEmail" >测试</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="json web token">
          <el-form :model="serverConfig" label-width="100px">
            <el-form-item label="jwt签名">
              <el-input v-model="serverConfig.jwt.signing_key"/>
            </el-form-item>
            <el-form-item label="签发者">
              <el-input v-model="serverConfig.jwt.issuer"/>
            </el-form-item>
            <el-form-item label="过期时间">
              <el-input v-model="serverConfig.jwt.expires_time"/>
            </el-form-item>
            <el-divider></el-divider>
            <el-form-item>
              <el-button @click="onSubmit" type="primary">保存</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
      <PayDialog ref="PayDialogRef" @refresh="payStore.getPayList()"></PayDialog>
    </el-card>
    <el-dialog v-model="state.isShowTestEmailDialog" :title="state.title" width="80%" destroy-on-close center>
      <el-input v-model="registerData.user_name" placeholder="输入电子邮件地址"/>
      <template #footer>
      <span class="dialog-footer">
        <el-button @click="onGetEmailCode">发送</el-button>
      </span>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {defineAsyncComponent, onMounted, reactive, ref} from "vue";

import {useServerStore} from "/@/stores/serverStore";
import {storeToRefs} from "pinia";
import {useShopStore} from "/@/stores/shopStore";
import {usePayStore} from "/@/stores/payStore";

import {ElMessage} from "element-plus";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
import {useUserStore} from "/@/stores/userStore";

const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
const PayDialog = defineAsyncComponent(() => import("/@/views/admin/system/dialog_pay.vue"))
const PayDialogRef = ref()
const serverStore = useServerStore()
const {serverConfig} = storeToRefs(serverStore)
const shopStore = useShopStore()
const {goodsList} = storeToRefs(shopStore)
const payStore = usePayStore()
const payStoreData = storeToRefs(payStore)
const userStore = useUserStore()
const {registerData} = storeToRefs(userStore)

const state = reactive({
  isShowTestEmailDialog: false,
  title: "电子邮件测试"
});

//测试邮件
const onTestEmail=()=>{
  state.isShowTestEmailDialog=true

}
//获取邮件验证码
const onGetEmailCode = () => {
  if (registerData.value.user_name === '') {
    return
  }
  request(apiStoreData.staticApi.value.public_getEmailCode, userStore.registerData).then((res) => {
    ElMessage.success(res.msg)
  })
};


//打开支付编辑
const openPayDialog = (type: string, row?: PayInfo) => {
  PayDialogRef.value.openDialog(type, row)
}

//保存提交
const onSubmit = () => {
  serverStore.updateServerConfig(serverConfig.value)
  setTimeout(() => {
    serverStore.getServerConfig()
  }, 500)
}
//删除支付
const deletePay = (data: PayInfo) => {
  // payApi.deletePayApi(data).then((res) => {
  request(apiStoreData.api.value.pay_deletePay, data).then((res) => {
    ElMessage.success(res.msg)
    setTimeout(() => {
      payStore.getPayList(); //获取支付列表
    }, 500);
  })

}
onMounted(() => {
  serverStore.getServerConfig() //获取设置参数
  shopStore.getAllGoods()       //获取全部商品，用来设置新注册分配套餐
  payStore.getPayList()         //获取支付列表
});

</script>

<style lang="scss" scoped>

</style>