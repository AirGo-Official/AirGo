<template>
  <el-dialog v-model="state.isShowDialog" :title="state.title" width="769px" destroy-on-close align-center>
    <el-divider content-position="left">节点参数</el-divider>
    <el-form :model="dialogData.vlessInfo" label-width="120px">
      <el-form-item label="node_type">
        <el-radio-group v-model="dialogData.vlessInfo.node_type">
          <el-radio label="vmess">vmess</el-radio>
          <el-radio label="vless">vless</el-radio>
          <el-radio label="trojan">trojan</el-radio>
          <el-radio label="shadowsocks">shadowsocks</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="remarks">
        <el-input v-model="dialogData.vlessInfo.remarks" placeholder="input node name"/>
      </el-form-item>
      <el-form-item label="address">
        <el-input v-model="dialogData.vlessInfo.address"/>
      </el-form-item>
      <el-form-item label="port">
        <el-input v-model.number="dialogData.vlessInfo.port"/>
      </el-form-item>

      <el-form-item label="scy" v-if="dialogData.vlessInfo.node_type!=='trojan'">
        <el-select
            v-model="dialogData.vlessInfo.scy"
            filterable
            allow-create
            default-first-option
            :reserve-keyword="false"
            style="width: 100%"
        >
          <el-option
              v-for="(v,k) in state.scyArr"
              :key="k"
              :label="v"
              :value="v">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="flow" v-if="dialogData.vlessInfo.node_type==='vless'">
          <el-select
              v-model="dialogData.vlessInfo.flow"
              filterable
              allow-create
              default-first-option
              :reserve-keyword="false"
              style="width: 100%"
          >
            <el-option
                v-for="(v,k) in state.flowArr"
                :key="k"
                :label="v"
                :value="v">
            </el-option>
          </el-select>
      </el-form-item>
      <el-form-item label="network" v-if="dialogData.vlessInfo.node_type !== 'shadowsocks'">
        <el-select
            v-model="dialogData.vlessInfo.network"
            filterable
            allow-create
            default-first-option
            :reserve-keyword="false"
            style="width: 100%"
        >
          <el-option
              v-for="(v,k) in state.networkArr"
              :key="k"
              :label="v"
              :value="v">
          </el-option>
        </el-select>

      </el-form-item>
      <el-form-item label="type" v-if="dialogData.vlessInfo.network==='tcp' || dialogData.vlessInfo.network==='kcp' || dialogData.vlessInfo.network=='quic' || dialogData.vlessInfo.node_type === 'shadowsocks'">
        <el-radio-group v-model="dialogData.vlessInfo.type">
          <el-radio label="none"
                    v-if="dialogData.vlessInfo.network==='tcp' || dialogData.vlessInfo.network==='kcp' || dialogData.vlessInfo.network=='quic' || dialogData.vlessInfo.node_type === 'shadowsocks'">
            none
          </el-radio>
          <el-radio label="http" v-if="dialogData.vlessInfo.network==='tcp' || dialogData.vlessInfo.node_type === 'shadowsocks'">http</el-radio>
          <el-radio label="srtp"
                    v-if="dialogData.vlessInfo.network==='kcp' || dialogData.vlessInfo.network=='quic'">
            srtp
          </el-radio>
          <el-radio label="utp"
                    v-if="dialogData.vlessInfo.network==='kcp' || dialogData.vlessInfo.network=='quic'">
            utp
          </el-radio>
          <el-radio label="wechat-video"
                    v-if="dialogData.vlessInfo.network==='kcp' || dialogData.vlessInfo.network=='quic'">
            wechat-video
          </el-radio>
          <el-radio label="dtls"
                    v-if="dialogData.vlessInfo.network==='kcp' || dialogData.vlessInfo.network=='quic'">
            dtls
          </el-radio>
          <el-radio label="wireguard"
                    v-if="dialogData.vlessInfo.network==='kcp' || dialogData.vlessInfo.network=='quic'">
            wireguard
          </el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="host" v-if="dialogData.vlessInfo.network === 'ws' || (dialogData.vlessInfo.node_type === 'shadowsocks' && dialogData.vlessInfo.type==='http') || dialogData.vlessInfo.node_type === 'trojan'">
        <el-input v-model="dialogData.vlessInfo.host"/>
      </el-form-item>
      <el-form-item label="path" v-if="dialogData.vlessInfo.network === 'ws' || (dialogData.vlessInfo.network === 'tcp' && dialogData.vlessInfo.type === 'http') ||  (dialogData.vlessInfo.node_type === 'shadowsocks' && dialogData.vlessInfo.type==='http') || dialogData.vlessInfo.node_type === 'trojan'">
        <el-input v-model="dialogData.vlessInfo.path"/>
      </el-form-item>
      <el-form-item label="mode" v-if="dialogData.vlessInfo.network==='grpc'">
        <el-radio-group v-model="dialogData.vlessInfo.mode">
          <el-radio label="gun">gun</el-radio>
          <el-radio label="multi">multi</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="serviceName" v-if="dialogData.vlessInfo.network==='grpc'">
        <el-input v-model="dialogData.vlessInfo.service_name"/>
      </el-form-item>

      <el-form-item label="security">
        <el-radio-group v-model="dialogData.vlessInfo.security">
          <el-radio label="none">none</el-radio>
          <el-radio label="tls">tls</el-radio>
          <el-radio label="reality">reality</el-radio>
        </el-radio-group>
      </el-form-item>
        <el-form-item label="dest" v-if="dialogData.vlessInfo.security==='reality'">
          <el-select
              v-model="dialogData.vlessInfo.dest"
              filterable
              allow-create
              default-first-option
              :reserve-keyword="false"
              style="width: 100%"
          >
            <el-option
                v-for="(v,k) in state.realityDefaultArr"
                :key="k"
                :label="v.dest"
                :value="v.dest">
            </el-option>
          </el-select>
        </el-form-item>

      <el-form-item label="sni" v-if="dialogData.vlessInfo.security==='tls' || dialogData.vlessInfo.security==='reality'">
        <el-input v-model="dialogData.vlessInfo.sni"/>
      </el-form-item>
      <el-form-item label="fp" v-if="dialogData.vlessInfo.security==='tls' ||dialogData.vlessInfo.security==='reality'">
        <el-select
            v-model="dialogData.vlessInfo.fp"
            filterable
            allow-create
            default-first-option
            :reserve-keyword="false"
            style="width: 100%"
        >
          <el-option
              v-for="(v,k) in state.fpArr"
              :key="k"
              :label="v"
              :value="v">
          </el-option>
        </el-select>
      </el-form-item>
<!--      <el-form-item label="alpn" v-if="dialogData.vlessInfo.security==='tls'">-->
<!--        <el-input v-model="dialogData.vlessInfo.alpn"/>-->
<!--      </el-form-item>-->
<!--      <el-form-item label="allowInsecure" v-if="dialogData.vlessInfo.security==='tls'">-->
<!--        <el-switch-->
<!--            size="small"-->
<!--            v-model="dialogData.vlessInfo.allowInsecure"-->
<!--            style="&#45;&#45;el-switch-on-color: #13ce66; &#45;&#45;el-switch-off-color: #ff4949"-->
<!--        />-->
<!--      </el-form-item>-->

      <el-form-item label="public_key" v-if="dialogData.vlessInfo.security==='reality'">
        <el-input v-model="dialogData.vlessInfo.pbk">
          <template #append>
            <el-button @click="setReality()"><el-icon><Refresh /></el-icon></el-button>
          </template>
        </el-input>
      </el-form-item>
      <el-form-item label="private_key" v-if="dialogData.vlessInfo.security==='reality'">
        <el-input v-model="dialogData.vlessInfo.private_key"/>
      </el-form-item>
<!--      <el-form-item label="sid" v-if="dialogData.vlessInfo.security==='reality'">-->
<!--        <el-input v-model="dialogData.vlessInfo.sid"/>-->
<!--      </el-form-item>-->
<!--      <el-form-item label="spx" v-if="dialogData.vlessInfo.security==='reality'">-->
<!--        <el-input v-model="dialogData.vlessInfo.spx"/>-->
<!--      </el-form-item>-->
    </el-form>
    <el-divider content-position="left">其他参数</el-divider>
    <el-form :model="dialogData.vlessInfo" label-width="120px">
      <el-form-item label="是否启用">
        <el-switch
            size="small"
            v-model="dialogData.vlessInfo.enabled"
            style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
        />
      </el-form-item>
      <el-form-item label="节点限速">
        <el-input type="number" v-model.number="dialogData.vlessInfo.node_speedlimit"/>
      </el-form-item>
      <el-form-item label="节点倍率">
        <el-input type="number" v-model.number="dialogData.vlessInfo.traffic_rate"/>
      </el-form-item>
      <el-form-item label="启用中转">
        <el-switch
            size="small"
            v-model="dialogData.vlessInfo.enable_transfer"
            style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
        />
      </el-form-item>
      <el-form-item label="中转ip" v-if="dialogData.vlessInfo.enable_transfer">
        <el-input v-model="dialogData.vlessInfo.transfer_address" placeholder=""/>
      </el-form-item>
      <el-form-item label="中转端口" v-if="dialogData.vlessInfo.enable_transfer">
        <el-input v-model="dialogData.vlessInfo.transfer_port" placeholder=""/>
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

import {storeToRefs} from "pinia";
import {useNodeStore} from "/@/stores/nodeStore";
import {reactive, watch} from "vue";
import {useApiStore} from "/@/stores/apiStore";
import {request} from "/@/utils/request";
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)


const nodeStore = useNodeStore()
const {dialogData} = storeToRefs(nodeStore)
const emit = defineEmits(['refresh']);
const state = reactive({
  type: "",
  title: "",
  isShowDialog: false,
  realityDefaultArr:[
    {dest:"www.speedtest.org:443",sni:"www.speedtest.org"},
    {dest:"www.lovelive-anime.jp:443",sni:"www.lovelive-anime.jp"},
    {dest:"swdist.apple.com:443",sni:"swdist.apple.com"},
    {dest:"blog.api.www.cloudflare.com:443",sni:"blog.api.www.cloudflare.com"},
    {dest:"www.icloud.com:443",sni:"www.icloud.com"},
  ] as RealityItem[],
  scyArr:[
      "auto",
      "none",
      "chacha20-poly1305",
      "aes-128-gcm",
      "aes-256-gcm",
      "2022-blake3-aes-128-gcm",
      "2022-blake3-aes-256-gcm",
      "2022-blake3-chacha20-poly1305",
  ],
  flowArr:[
      "",
      "xtls-rprx-vision",
      "xtls-rprx-vision-udp443",
  ],
  networkArr:[
      "ws",
      "tcp",
      "kcp",
      "quic",
      "grpc",
  ],
  fpArr:[
      "chrome",
      "firefox",
      "safari",
      "ios",
      "android",
      "edge",
      "360",
      "qq",
      "random",
      "randomized",
  ],
})

// 打开弹窗
const openDialog = (type: string, row?: any) => {
  if (type == 'add') {
    dialogData.value.vlessInfo.id = 0 //编辑和添加公用一个store，清空id,否则服务器无法插入
    state.type = type
    state.title = "新建节点"
    state.isShowDialog = true
  } else {
    state.type = type
    state.title = "修改节点"
    dialogData.value.vlessInfo = row  //将当前row写入pinia
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
    //新建节点
    nodeStore.newNode(dialogData.value.vlessInfo)
    setTimeout(() => {
      emit('refresh');
    }, 500);
  } else {
    //更新节点
    nodeStore.updateNode(dialogData.value.vlessInfo)
    setTimeout(() => {
      emit('refresh');
    }, 500);
  }
  closeDialog()
}
const setReality = ()=>{
  request(apiStoreData.api.value.system_createx25519).then((res)=>{
  // console.log("res:",res)
  dialogData.value.vlessInfo.pbk=res.data.public_key
  dialogData.value.vlessInfo.private_key=res.data.private_key
})
}
//监听
watch(
    () => dialogData.value.vlessInfo.dest,
    () => {
      let temp = state.realityDefaultArr.filter(r => r.dest === dialogData.value.vlessInfo.dest)
      if (temp.length > 0){
        dialogData.value.vlessInfo.sni = temp[0].sni
      }
      if (dialogData.value.vlessInfo.private_key === '' || dialogData.value.vlessInfo.private_key === ''){
        setReality()
      }
    },
    {
      // deep: true,
    }
);

// 暴露变量
defineExpose({
  openDialog,   // 打开弹窗
});
</script>


<style scoped>
.dialog-footer button:first-child {
  margin-right: 10px;
}
</style>
