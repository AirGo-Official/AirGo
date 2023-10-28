<template>
  <div>
    <div class="img">
      <el-button color="#626aef" round style="height: 100%" @click="openDialog">我的图库列表</el-button>
    </div>
    <div>
      <iframe name="if1" src="https://telegraph-image.pages.dev/" height="800px" width="100%"></iframe>
    </div>

    <div>
      <el-input class="input1" v-model="state.galleryData.picture_url" placeholder="粘贴上方url"/>
      <el-input class="input2" v-model="state.galleryData.subject" placeholder="备注"/>
      <el-button class="button" color="#626aef" @click="savePicture" round>保 存</el-button>
    </div>
    <el-dialog
        v-model="state.galleryDialogData.isShowDialog"
        title="我的图库列表"
        width="50%"
        align-center
    >
      <el-input style="width: 200px" v-model="state.galleryDialogData.params.search" placeholder="输入图片备注"/>
      <el-button color="#626aef" @click="getPictureList(state.galleryDialogData.params)">查询</el-button>
      <el-table :data="state.galleryDialogData.galleryList" height="500px" stripe>
        <el-table-column property="picture_url" label="预览" width="120">
          <template #default="scope">
            <div style="display: flex; align-items: center">
              <el-image
                  style="width: 100px; height: 100px"
                  :src="scope.row.picture_url"
                  fit="cover"
              />
            </div>
          </template>
        </el-table-column>
        <el-table-column property="subject" label="备注" width="100"/>
        <el-table-column property="picture_url" label="操作">
          <template #default="scope">
            <div>
              <el-button color="#626aef" @click="copyText(scope.row.picture_url)">复制链接</el-button>
            </div>
          </template>
        </el-table-column>

      </el-table>
      <template #footer>
      <span class="dialog-footer">
        <el-button @click="state.galleryDialogData.isShowDialog = false">关闭</el-button>
      </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import {ElMessage} from 'element-plus';
import {onMounted, reactive} from 'vue';
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
import {storeToRefs} from "pinia";
import commonFunction from '/@/utils/commonFunction';

const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)

const {copyText} = commonFunction();

const state = reactive({
  galleryData: {
    picture_url: '',
    subject: '',
  },
  galleryDialogData: {
    isShowDialog: false,
    galleryList: [] as UploadPicture[],
    params: {
      search: '',
      page_num: 1,
      page_size: 30,
    } as PaginationParams,
  },

})

function savePicture() {
  request(apiStoreData.api.value.upload_newPictureUrl, state.galleryData).then((res) => {
    ElMessage.success(res.msg)
    state.galleryData.picture_url = '' //清空输入框
    state.galleryData.subject = '' //清空输入框
  }).catch()
}

function openDialog() {
  state.galleryDialogData.isShowDialog = true
  getPictureList(state.galleryDialogData.params)

}

function getPictureList(params?: object) {
  request(apiStoreData.api.value.upload_getPictureList, params).then((res) => {
    ElMessage.success(res.msg)
    state.galleryDialogData.galleryList = res.data
  }).catch()

}

// 页面加载时
onMounted(() => {

});
</script>

<style scoped lang="scss">
.img {
  height: 40px;
  position: absolute;
  right: 20px;
  top: 20px;
  //transform: translate(-50%, -50%); /* 50%为自身尺寸的一半 */
}

.button {
  //display: flex;
  //justify-content: space-between;
  margin-top: auto;
  margin-bottom: auto;

  width: 280px;
  height: 32px;

  position: absolute;
  left: 50%;
  top: 80%;
  transform: translate(-50%, -50%); /* 50%为自身尺寸的一半 */
}

.input1 {
  width: 280px;
  height: 32px;
  position: absolute;
  left: 50%;
  top: 70%;
  transform: translate(-50%, -50%); /* 50%为自身尺寸的一半 */
}

.input2 {
  width: 280px;
  height: 32px;
  position: absolute;
  left: 50%;
  top: 75%;
  transform: translate(-50%, -50%); /* 50%为自身尺寸的一半 */
}

//:deep(.el-input__wrapper) { //el-input__inner
//  background: rgba(0, 0, 0, 0.2); /*调整inner的背景色，透明*/
//  border: 0;
//  height: 32px; /*调整inner的高度*/
//  border-radius: 50px; /*输入框圆角值*/
//}

</style>