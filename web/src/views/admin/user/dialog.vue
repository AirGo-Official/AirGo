<template>
  <div class="system-user-dialog-container">
    <el-dialog :title="state.title" v-model="state.isShowDialog" width="769px">
      <el-form ref="userDialogFormRef" :model="userManageData.dialog.user" size="default" label-width="90px">
        <el-row :gutter="35">
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="账户邮箱">
              <el-input v-model="userManageData.dialog.user.user_name" placeholder="请输入账户邮箱"
                        clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="账户密码">
              <el-input v-model="userManageData.dialog.user.password" placeholder="请输入"
                        clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col v-if="userInfos.id===1" :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20">
            <el-form-item label="关联角色">
              <el-checkbox-group v-model="userManageData.dialog.check_list">
                <el-checkbox :label="v.role_name" v-for="(v,index) in roleManageData.roles.role_list"
                             :key="index"></el-checkbox>
              </el-checkbox-group>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="用户状态">
              <el-switch v-model="userManageData.dialog.user.enable" inline-prompt active-text="开启"
                         inactive-text="关闭"></el-switch>
            </el-form-item>
          </el-col>

          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="订阅状态">
              <el-switch v-model="userManageData.dialog.user.subscribe_info.sub_status" inline-prompt active-text="开启"
                         inactive-text="关闭"></el-switch>
            </el-form-item>
          </el-col>

          <!--          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">-->
          <!--            <el-form-item label="订阅url">-->
          <!--              <el-input v-model="userManageData.dialog.user.subscribe_info.subscribe_url"></el-input>-->
          <!--            </el-form-item>-->
          <!--          </el-col>-->

          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="分配套餐">
              <el-select v-model="userManageData.dialog.user.subscribe_info.goods_id" class="m-2"
                         placeholder="选择套餐">
                <el-option
                    v-for="item in goodsList"
                    :key="item.id"
                    :label="item.subject"
                    :value="item.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="到期时间">
              <el-date-picker
                  v-model="userManageData.dialog.user.subscribe_info.expired_at"
                  type="datetime"
                  placeholder="选择到期时间"
                  size="default"
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="总流量">
              <el-row :gutter="0">
                <el-col :span="12">
                  <el-input v-model.number="state.subParams.t" type="number"/>
                </el-col>
                <el-col :span="2" style="text-align: center"><span>-</span></el-col>
                <el-col :span="10">
                  <span>GB</span>
                </el-col>
              </el-row>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="限速">
              <el-row :gutter="0">
                <el-col :span="12">
                  <el-input v-model.number="userManageData.dialog.user.subscribe_info.node_speedlimit" type="number"/>
                </el-col>
                <el-col :span="2" style="text-align: center"><span>-</span></el-col>
                <el-col :span="10">
                  <span>Mbps</span>
                </el-col>
              </el-row>

            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="连接数">
              <el-row :gutter="0">
                <el-col :span="12">
                  <el-input v-model.number="userManageData.dialog.user.subscribe_info.node_connector" type="number"/>
                </el-col>
                <el-col :span="2" style="text-align: center"><span>-</span></el-col>
                <el-col :span="10">
                  <span>个</span>
                </el-col>
              </el-row>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="我的邀请人">
              <el-input v-model="userManageData.dialog.user.referrer_code"></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="我的邀请码">
              <el-input v-model="userManageData.dialog.user.invitation_code"></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12" class="mb20">
            <el-form-item label="我的余额">
              <el-input v-model.number="userManageData.dialog.user.remain" type="number"></el-input>
            </el-form-item>
          </el-col>


        </el-row>
      </el-form>
      <template #footer>
				<span class="dialog-footer">
					<el-button @click="closeDialog" size="default">取 消</el-button>
					<el-button type="primary" @click="onSubmit" size="default">{{ state.submitTxt }}</el-button>
				</span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import {reactive, ref} from 'vue';
import {useUserStore} from '/@/stores/userStore'
import {storeToRefs} from 'pinia';
import {useRoleStore} from '/@/stores/roleStore'
import {useShopStore} from "/@/stores/shopStore";
import {GetLocalTime} from "/@/utils/formatTime";

const state = reactive({
  isShowDialog: false,
  type: '',
  title: '',
  submitTxt: '',
  params: {
    search: '',
    page_num: 1,
    page_size: 10,
  },
  subParams: {
    t: 0,
    u: 0,
    d: 0,
  }
})

const userStore = useUserStore()
const {userManageData, userInfos} = storeToRefs(userStore)
const roleStore = useRoleStore()
const {roleManageData} = storeToRefs(roleStore)
const shopStore = useShopStore()
const {goodsList} = storeToRefs(shopStore)

// 定义子组件向父组件传值/事件
const emit = defineEmits(['refresh']);
const userDialogFormRef = ref();

// 打开弹窗
const openDialog = (type: string, row: SysUser) => {
  state.isShowDialog = true;
  shopStore.getAllGoods() //获取全部套餐
  if (type === 'edit') {
    state.title = '修改用户';
    state.submitTxt = '修 改';
    userManageData.value.dialog.user = JSON.parse(JSON.stringify(row)) //深拷贝,防止修改时间报错
    //计算流量
    state.subParams.t = userManageData.value.dialog.user.subscribe_info.t / 1024 / 1024 / 1024
    //计算用户的角色
    let currentUserRoleIds: string[] = []
    userManageData.value.dialog.user.role_group.forEach((item: any) => {
      currentUserRoleIds.push(item.role_name)
    })
    userManageData.value.dialog.check_list = currentUserRoleIds
  } else {
    state.title = '新增用户';
    state.submitTxt = '新 增';
    userStore.resetData();
    (userManageData.value.dialog.user.subscribe_info.expired_at as Date) = GetLocalTime(8)
  }
//打开时加载全部角色，用来设置用户角色
  roleStore.getRoleList({page_num: 1, page_size: 10000})
};
// 关闭弹窗
const closeDialog = () => {
  state.isShowDialog = false;
};
// 提交
const onSubmit = () => {
  //处理角色
  userManageData.value.dialog.user.role_group = [];
  userManageData.value.dialog.check_list.forEach((value: string, index: number, array: string[]) => {
    userManageData.value.dialog.user.role_group.push({
      role_name: value,
    } as RowRoleType)
  })
  //处理流量
  if (state.subParams.t !== 0) {     //计算流量
    userManageData.value.dialog.user.subscribe_info.t = state.subParams.t * 1024 * 1024 * 1024
  }
  if (state.title === '新增用户') {
    userStore.newUser(userManageData.value.dialog.user)
  } else {
    userStore.updateUser(userManageData.value.dialog.user)
  }
  setTimeout(() => {
    userStore.getUserList(state.params)
  }, 500)
  closeDialog();
};
// 暴露变量
defineExpose({
  openDialog,
});
</script>
