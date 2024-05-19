<template>
  <div class="personal layout-pd">
    <el-row gutter="20">
      <!-- 个人信息 -->
      <el-col :xs="24" :sm="24">
        <el-card shadow="hover">
          <div class="personal-user">
            <div class="personal-user-left">
              <div class="h400 personal-user-left-upload" @click="state.isShowChangeAvatarDialog = true">
                <img :src="userInfos.avatar" />
              </div>
            </div>
            <div class="personal-user-right">
              <el-row>
                <el-col :span="24" class="personal-title mb18">{{ currentTime }}</el-col>
                <el-col :span="24">
                  <el-row>
                    <el-col :xs="24" :sm="24" class="personal-item mb6">
                      <div class="personal-item-label">{{$t('message.personal.nickname')}}：</div>
                      <span>{{ userInfos.nick_name }}</span>
                    </el-col>
                  </el-row>
                </el-col>
              </el-row>
            </div>
          </div>
          <div>
            <el-button type="primary" :disabled="!publicStoreData.publicSetting.value.enable_lottery" @click="lottery">{{ $t("message.personal.lottery") }}</el-button>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12">
        <el-card shadow="hover" class="mt15 personal-edit">
          <div class="personal-edit-title">{{$t('message.personal.message_setting')}}</div>
          <el-divider content-position="left"><span>{{ $t("message.personal.push_setting") }}</span>
          </el-divider>
          <el-text class="mx-1">{{ $t("message.adminServer.Server.push_method") }} :</el-text>
          <el-form :model="userInfos" label-width="100px"
                   label-position="left">
            <el-form-item :label="$t('message.adminServer.Server.enable_tg_bot')" class="label">
              <el-switch v-model="userInfos.enable_tg_bot" inline-prompt
                         :active-text="$t('message.common.enable')"
                         :inactive-text="$t('message.common.disable')"
                         style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
            </el-form-item>
            <el-form-item :label="$t('message.adminUser.SysUser.tg_id')" class="label">
              <el-input v-model.number="userInfos.tg_id"/>
            </el-form-item>
            <el-divider content-position="left"><span>{{ $t("message.adminServer.Server.trigger_condition") }}</span>
            </el-divider>
            <el-form-item :label="$t('message.adminUser.SysUser.when_service_almost_expired')" class="label">
              <el-switch v-model="userInfos.when_service_almost_expired" inline-prompt
                         :active-text="$t('message.common.enable')"
                         :inactive-text="$t('message.common.disable')"
                         style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
            </el-form-item>
            <el-form-item :label="$t('message.adminUser.SysUser.when_purchased')" class="label">
              <el-switch v-model="userInfos.when_purchased" inline-prompt
                         :active-text="$t('message.common.enable')"
                         :inactive-text="$t('message.common.disable')"
                         style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
            </el-form-item>
            <el-form-item :label="$t('message.adminUser.SysUser.when_balance_changed')" class="label">
              <el-switch v-model="userInfos.when_balance_changed" inline-prompt
                         :active-text="$t('message.common.enable')"
                         :inactive-text="$t('message.common.disable')"
                         style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
            </el-form-item>
            <el-form-item style="margin-top: 20px">
              <el-button @click="onSubmitForNotice()" type="primary" style="margin-left: auto">{{ $t("message.common.button_confirm") }}
              </el-button>
            </el-form-item>

          </el-form>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12">
        <el-card shadow="hover" class="mt15 personal-edit">
          <div class="personal-edit-title">{{$t('message.personal.account_setting')}}</div>
          <div class="personal-edit-safe-box">
            <div class="personal-edit-safe-item">
              <div class="personal-edit-safe-item-left">
                <div class="personal-edit-safe-item-left-label">{{$t('message.personal.login_password')}}</div>
              </div>
              <div class="personal-edit-safe-item-right">
                <el-button type="primary" @click="state.isShowChangePasswordDialog = true">{{$t('message.common.modify')}}</el-button>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    <el-dialog v-model="state.isShowChangePasswordDialog" :title="$t('message.personal.change_password')" width="500px" destroy-on-close>
      <el-form size="default" label-position="left">
        <el-form-item :label="$t('message.personal.new_password')">
          <el-input v-model="registerData.password":placeholder="$t('message.personal.please_enter_password')" clearable></el-input>
        </el-form-item>
        <el-form-item :label="$t('message.personal.confirm_password')">
          <el-input v-model="registerData.re_password" :placeholder="$t('message.personal.please_enter_password')" clearable></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
				<span class="dialog-footer">
					<el-button @click="state.isShowChangePasswordDialog = false" size="default">{{$t('message.common.button_cancel')}}</el-button>
					<el-button type="primary" @click="changePassword" size="default">{{$t('message.common.button_confirm')}}</el-button>
				</span>
      </template>
    </el-dialog >
    <el-dialog v-model="state.isShowChangeAvatarDialog" :title="$t('message.personal.modify_avatar')" width="500px" destroy-on-close>
      <el-form size="default" label-position="left">
        <el-form-item :label="$t('message.personal.avatar_link')">
          <el-input v-model="userAvatar.avatar"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
				<span class="dialog-footer">
					<el-button @click="state.isShowChangeAvatarDialog = false" size="default">{{$t('message.common.button_cancel')}}</el-button>
					<el-button type="primary" @click="changeAvatar" size="default">{{$t('message.common.button_confirm')}}</el-button>
				</span>
      </template>
    </el-dialog>

    <div class="dialog">
      <el-dialog  v-model="state.isShowLotteryDialog" destroy-on-close width="350px">
        <LuckyWheel
          ref="myLuckyRef"
          width="300px"
          height="300px"
          :prizes="state.lotteryData.prizes"
          :blocks="state.lotteryData.blocks"
          :buttons="state.lotteryData.buttons"
          @start="startCallback"
          @end="endCallback"
        />
      </el-dialog>
    </div>


  </div>
</template>

<script setup lang="ts">
import { computed, reactive,ref } from "vue";
import { formatAxis } from "/@/utils/formatTime";
import { useUserStore } from "/@/stores/user_logic/userStore";
import { storeToRefs } from "pinia";
import { useApiStore } from "/@/stores/apiStore";
import { ElMessage, ElMessageBox } from "element-plus";
import { usePublicStore } from "/@/stores/publicStore";
import { useI18n } from "vue-i18n";
import giftIcon  from "/@/assets/icon/gift.svg"

const userStore = useUserStore();
const { userInfos, registerData, userAvatar } = storeToRefs(userStore);

const publicStore = usePublicStore();
const publicStoreData = storeToRefs(publicStore);
const apiStore = useApiStore();
const myLuckyRef = ref()
const {t} = useI18n()



const state = reactive({
  isShowChangePasswordDialog: false,
  isShowChangeAvatarDialog: false,
  isShowLotteryDialog:false,
  lotteryData:{
    blocks: [{ padding: '13px', background: '#617df2' }],
    prizes: [
      { background: '#e9e8fe', imgs: [{ src: '', width: '20%', top: '40%' }], fonts:[{ text: '0000', top: '10%' }]}
    ],
    buttons: [{
      radius: '35%',
      background: '#8a9bf3',
      pointer: true,
      fonts: [{ text: 'Go', top: '-10px' }]
    }],
  },
  prizeImg:{
    src: giftIcon,
    width: '20%',
    top: '40%'
  },
});
//设置奖池
const setJackpot=()=>{
  state.lotteryData.prizes = []
  publicStoreData.publicSetting.value.jackpot.forEach((value: JackpotItem, index: number, array: JackpotItem[])=>{
    let background = '#e9e8fe'
    if(index%2 === 0){
      background = '#b8c5f2'
    }
    state.lotteryData.prizes.push({imgs: [state.prizeImg], fonts: [{ text: '+ '+value.balance.toString(), top: '10%' }],  background: background })
  })
}

//开始抽奖
const startCallback =() =>{
  // 调用抽奖组件的play方法开始游戏
  userStore.clockIn().then((res)=>{
    myLuckyRef.value.play()
    const index = res.data.data //res.data 格式： {"total":1,"data":1}
    myLuckyRef.value.stop(index)
    userStore.getUserInfo()
  }).then(()=>{

  })
}
// 抽奖结束会触发end回调
const endCallback =(prize:any)=> {
  state.isShowLotteryDialog = false
  ElMessageBox.alert(t('message.personal.balance')+prize.fonts[0].text, t('message.personal.congratulations'), {
    confirmButtonText: 'OK',
  })
}

// 当前时间提示语
const currentTime = computed(() => {
  return formatAxis(new Date());
});


const changePassword = () => {
  userStore.changePassword().then(() => {
  });
  state.isShowChangePasswordDialog = false;
};
const changeAvatar = () => {
  userStore.changeAvatar().then(() => {
    userStore.setUserNotice().then((res) => {
      ElMessage.success(res.msg);
    });
    userStore.getUserInfo();
  });
  state.isShowChangeAvatarDialog = false;
};

const lottery=()=>{
  state.isShowLotteryDialog = true
  setJackpot()
}
const onSubmitForNotice = () => {
  userStore.setUserNotice().then((res) => {
    ElMessage.success(res.msg);
  });
};

</script>

<style scoped lang="scss">
@import '../../theme/mixins/index.scss';

.personal {
  .personal-user {
    height: 130px;
    display: flex;
    align-items: center;

    .personal-user-left {
      width: 100px;
      height: 130px;
      border-radius: 3px;

      :deep(.el-upload) {
        height: 100%;
      }

      .personal-user-left-upload {
        img {
          width: 100%;
          height: 100%;
          border-radius: 3px;
        }

        &:hover {
          img {
            animation: logoAnimation 0.3s ease-in-out;
          }
        }
      }
    }

    .personal-user-right {
      flex: 1;
      padding: 0 15px;

      .personal-title {
        font-size: 18px;
        @include text-ellipsis(1);
      }

      .personal-item {
        display: flex;
        align-items: center;
        font-size: 13px;

        .personal-item-label {
          color: var(--el-text-color-secondary);
          @include text-ellipsis(1);
        }

        .personal-item-value {
          @include text-ellipsis(1);
        }
      }
    }
  }

  .personal-info {
    .personal-info-more {
      float: right;
      color: var(--el-text-color-secondary);
      font-size: 13px;

      &:hover {
        color: var(--el-color-primary);
        cursor: pointer;
      }
    }

    .personal-info-box {
      height: 130px;
      overflow: hidden;

      .personal-info-ul {
        list-style: none;

        .personal-info-li {
          font-size: 13px;
          padding-bottom: 10px;

          .personal-info-li-title {
            display: inline-block;
            @include text-ellipsis(1);
            color: var(--el-text-color-secondary);
            text-decoration: none;
          }

          & a:hover {
            color: var(--el-color-primary);
            cursor: pointer;
          }
        }
      }
    }
  }

  .personal-recommend-row {
    .personal-recommend-col {
      .personal-recommend {
        position: relative;
        height: 100px;
        border-radius: 3px;
        overflow: hidden;
        cursor: pointer;

        &:hover {
          i {
            right: 0px !important;
            bottom: 0px !important;
            transition: all ease 0.3s;
          }
        }

        i {
          position: absolute;
          right: -10px;
          bottom: -10px;
          font-size: 70px;
          transform: rotate(-30deg);
          transition: all ease 0.3s;
        }

        .personal-recommend-auto {
          padding: 15px;
          position: absolute;
          left: 0;
          top: 5%;
          color: var(--next-color-white);

          .personal-recommend-msg {
            font-size: 12px;
            margin-top: 10px;
          }
        }
      }
    }
  }

  .personal-edit {
    .personal-edit-title {
      position: relative;
      padding-left: 10px;
      color: var(--el-text-color-regular);

      &::after {
        content: '';
        width: 2px;
        height: 10px;
        position: absolute;
        left: 0;
        top: 50%;
        transform: translateY(-50%);
        background: var(--el-color-primary);
      }
    }

    .personal-edit-safe-box {
      border-bottom: 1px solid var(--el-border-color-light, #ebeef5);
      padding: 15px 0;

      .personal-edit-safe-item {
        width: 100%;
        display: flex;
        align-items: center;
        justify-content: space-between;

        .personal-edit-safe-item-left {
          flex: 1;
          overflow: hidden;

          .personal-edit-safe-item-left-label {
            color: var(--el-text-color-regular);
            margin-bottom: 5px;
          }

          .personal-edit-safe-item-left-value {
            color: var(--el-text-color-secondary);
            @include text-ellipsis(1);
            margin-right: 15px;
          }
        }
      }

      &:last-of-type {
        padding-bottom: 0;
        border-bottom: none;
      }
    }
  }
  .center {
    text-align:center;
    display:flex;justify-content: center; align-items:center;
  }
}

.dialog {
  :deep(.el-dialog) {
    box-shadow: 0 0px 0px rgb(0 0 0 / 0%);
    background: transparent;
    //.el-dialog__body {
    //  padding: 0 !important;
    //}
    //.el-dialog__header {
    //  display: none !important;
    //}
  }
}

</style>

