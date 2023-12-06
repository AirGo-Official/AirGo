<template>
  <div class="personal layout-pd">
    <el-row>
      <!-- 个人信息 -->
      <el-col :xs="24" :sm="24">
        <el-card shadow="hover" header="个人信息">
          <div class="personal-user">
            <div class="personal-user-left">
              <el-upload class="h400 personal-user-left-upload" accept=".bmg,.png,.jpg"
                         action="https://jsonplaceholder.typicode.com/posts/" multiple :limit="1">
                <img :src="userInfos.avatar"/>
              </el-upload>
            </div>
            <div class="personal-user-right">
              <el-row>
                <el-col :span="24" class="personal-title mb18">{{ currentTime }}</el-col>
                <el-col :span="24">
                  <el-row>
                    <el-col :xs="24" :sm="24" class="personal-item mb6">
                      <div class="personal-item-label">昵称：</div>
                      <span>{{ userInfos.nick_name }}</span>
                    </el-col>
                  </el-row>
                </el-col>
              </el-row>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 更新信息 -->
      <el-col :span="24">
        <el-card shadow="hover" class="mt15 personal-edit" header="个人信息">
          <div class="personal-edit-title">基本信息</div>
          <div class="personal-edit-safe-box">
            <div class="personal-edit-safe-item">
              <div class="personal-edit-safe-item-left">
                <div class="personal-edit-safe-item-left-label">账户密码</div>
                <!--                <div class="personal-edit-safe-item-left-value">当前密码强度：强</div>-->
              </div>
              <div class="personal-edit-safe-item-right">
                <el-button type="primary" @click="onOpenPWDialog">立即修改</el-button>
              </div>
            </div>
          </div>
          <div class="personal-edit-title">免流混淆：{{userInfos.subscribe_info.host}}</div>
          <div class="personal-edit-safe-box">
            <div class="personal-edit-safe-item">
              <div class="personal-edit-safe-item-left">
                <el-input class="personal-edit-safe-item-left-label" v-model="state.host.host" style="width: 90%"></el-input>
              </div>
              <div class="personal-edit-safe-item-right">
                <el-button type="primary" @click="onChangeHost">立即修改</el-button>
              </div>
            </div>
          </div>
          <div class="personal-edit-title">
            我的邀请（佣金率：{{ serverConfig.publicServerConfig.value.rebate_rate * 100 }}%）
          </div>
          <div class="personal-edit-safe-box">
            <div class="personal-edit-safe-item">
              <div class="personal-edit-safe-item-left">
<!--                <div class="personal-edit-safe-item-left-label">专属邀请链接</div>-->
                <div class="personal-edit-safe-item-left-value">
                  {{ state.url }}/#/login?i={{ userInfos.invitation_code }}
                </div>
              </div>
            </div>
          </div>
          <div class="personal-edit-title">
            余额
          </div>
          <div class="personal-edit-safe-box">
            <div class="personal-edit-safe-item">
              <div class="personal-edit-safe-item-left">
                <div class="personal-edit-safe-item-left-value">¥{{ userInfos.remain }}</div>
              </div>
            </div>
          </div>
          <div class="personal-edit-title">
            打卡
          </div>
          <div class="personal-edit-safe-box">
            <div class="personal-edit-safe-item">
              <div class="personal-edit-safe-item-left">
                <div class="personal-edit-safe-item-left-value">
                  <el-button type="primary" @click="clockin()">立即打卡</el-button>
                </div>
              </div>
            </div>
          </div>

        </el-card>
      </el-col>
    </el-row>
    <ChangePasswordDialog ref="changePasswordDialogRef"></ChangePasswordDialog>
  </div>
</template>

<script setup lang="ts" name="personal">
import {computed, defineAsyncComponent, onMounted, reactive, ref} from 'vue';
import {formatAxis} from '/@/utils/formatTime';
import * as imageConversion from 'image-conversion'
import {useUserStore} from "/@/stores/userStore";
import {storeToRefs} from 'pinia';
import {useServerStore} from "/@/stores/serverStore";
import {useApiStore} from "/@/stores/apiStore";
import {request} from "/@/utils/request";
import {ElMessage} from "element-plus";

const userStore = useUserStore()
const {userInfos} = storeToRefs(userStore)
const serverStore = useServerStore()
const serverConfig = storeToRefs(serverStore)
const ChangePasswordDialog = defineAsyncComponent(() => import('/@/views/personal/change_password_dialog.vue'));
const changePasswordDialogRef = ref()
const apiStore = useApiStore()

const state = reactive({
  url: '',
  host: {
    host: '',
  },
})

//打开修改密码弹窗
const onOpenPWDialog = () => {
  changePasswordDialogRef.value.openDialog()
}
//修改混淆
const onChangeHost = () => {
  userStore.changeHost(state.host)
  state.host.host = ''
}
//图片超过4M就压缩
function beforeUpload(file: any) {
  return new Promise((resolve, reject) => {
    let isLt2M = file.size / 1024 / 1024 < 4 // 判定图片大小是否小于4MB
    if (isLt2M) {
      resolve(file)
    }
    //  console.log(file) // 压缩到400KB,这里的400就是要压缩的大小,可自定义
    imageConversion.compressAccurately(file, 400).then(res => { // console.log(res)
      resolve(res)
    })
  })
}

//获取当前url
const getUrl = () => {
  state.url = window.location.host
}

// 当前时间提示语
const currentTime = computed(() => {
  return formatAxis(new Date());
});

//打卡
const clockin = () => {
  //是否开启打卡
  if (!serverConfig.publicServerConfig.value.enabled_clock_in){
    ElMessage.error("打卡已关闭")
    return
  }

  //判断订阅是否有效
  if (!userInfos.value.subscribe_info.sub_status){
    ElMessage.error("订阅已失效，请购买订阅后重试")
    return
  }
  request(apiStore.api.user_clockin).then((res)=>{
    ElMessage.success(res.data)
    userStore.getUserInfo()//刷新用户信息
  })
}

onMounted(() => {
  getUrl(); //获取专属邀请url

});
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
}
</style>
