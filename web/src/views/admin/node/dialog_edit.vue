<template>
  <el-dialog v-model="state.isShowDialog" :title="state.title" width="80%" destroy-on-close align-center>
    <div>
      <el-divider content-position="left">节点类型</el-divider>
      <el-form label-width="100px">
        <el-form-item label="节点类型">
          <el-select
              v-model="state.noteType"
              filterable
              allow-create
              default-first-option
              :reserve-keyword="false"
              style="width: 100%"
          >
            <el-option
                v-for="(v,k) in state.nodeTypelist"
                :key="k"
                :label="v"
                :value="v">
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
    </div>
    <el-divider content-position="left">基础参数</el-divider>
    <div v-if="state.noteType === 'vless'">
      <el-form :model="dialogData.vlessInfo" label-width="100px">
        <el-form-item label="remarks">
          <el-input v-model="dialogData.vlessInfo.remarks"/>
        </el-form-item>
        <el-form-item label="address">
          <el-input v-model="dialogData.vlessInfo.address"/>
        </el-form-item>
        <el-form-item label="port">
          <el-input v-model.number="dialogData.vlessInfo.port"/>
        </el-form-item>
        <el-form-item label="flow">
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
        <el-form-item label="network">
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
        <el-form-item label="type" v-if="dialogData.vlessInfo.network === 'tcp'">
          <el-select
              v-model="dialogData.vlessInfo.type"
              filterable
              allow-create
              default-first-option
              :reserve-keyword="false"
              style="width: 100%"
          >
            <el-option
                v-for="(v,k) in state.typeArr1"
                :key="k"
                :label="v"
                :value="v">
            </el-option>
          </el-select>
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

        <el-form-item label="host" v-if="dialogData.vlessInfo.network==='ws' || (dialogData.vlessInfo.network==='tcp' && dialogData.vlessInfo.type ==='http')">
          <el-input v-model="dialogData.vlessInfo.host"/>
        </el-form-item>
        <el-form-item label="path" v-if="dialogData.vlessInfo.network==='ws' || (dialogData.vlessInfo.network==='tcp' && dialogData.vlessInfo.type ==='http')">
          <el-input v-model="dialogData.vlessInfo.path"/>
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
        <el-form-item label="sni"
                      v-if="dialogData.vlessInfo.security==='tls' || dialogData.vlessInfo.security==='reality'">
          <el-input v-model="dialogData.vlessInfo.sni"/>
        </el-form-item>
        <el-form-item label="fp"
                      v-if="dialogData.vlessInfo.security==='tls' ||dialogData.vlessInfo.security==='reality'">
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
        <el-form-item label="public_key" v-if="dialogData.vlessInfo.security==='reality'">
          <el-input v-model="dialogData.vlessInfo.pbk">
            <template #append>
              <el-button @click="setReality('vless')">
                <el-icon>
                  <Refresh/>
                </el-icon>
              </el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="private_key" v-if="dialogData.vlessInfo.security==='reality'">
          <el-input v-model="dialogData.vlessInfo.private_key"/>
        </el-form-item>
        <el-divider content-position="left">其他参数</el-divider>
        <el-form :model="dialogData.vlessInfo" label-width="100px">
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
      </el-form>
    </div>

    <div v-if="state.noteType === 'vmess'">
      <el-form :model="dialogData.vmessInfo" label-width="100px">
        <el-form-item label="remarks">
          <el-input v-model="dialogData.vmessInfo.remarks"/>
        </el-form-item>
        <el-form-item label="address">
          <el-input v-model="dialogData.vmessInfo.address"/>
        </el-form-item>
        <el-form-item label="port">
          <el-input v-model.number="dialogData.vmessInfo.port"/>
        </el-form-item>

        <el-form-item label="scy">
          <el-select
              v-model="dialogData.vmessInfo.scy"
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
        <el-form-item label="network">
          <el-select
              v-model="dialogData.vmessInfo.network"
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
        <el-form-item label="type" v-if="dialogData.vmessInfo.network === 'tcp'">
          <el-select
              v-model="dialogData.vmessInfo.type"
              filterable
              allow-create
              default-first-option
              :reserve-keyword="false"
              style="width: 100%"
          >
            <el-option
                v-for="(v,k) in state.typeArr1"
                :key="k"
                :label="v"
                :value="v">
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="mode" v-if="dialogData.vmessInfo.network==='grpc'">
          <el-radio-group v-model="dialogData.vmessInfo.mode">
            <el-radio label="gun">gun</el-radio>
            <el-radio label="multi">multi</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="serviceName" v-if="dialogData.vmessInfo.network==='grpc'">
          <el-input v-model="dialogData.vmessInfo.service_name"/>
        </el-form-item>
        <el-form-item label="host" v-if="dialogData.vmessInfo.network==='ws' || (dialogData.vmessInfo.network==='tcp' && dialogData.vmessInfo.type ==='http')">
          <el-input v-model="dialogData.vmessInfo.host"/>
        </el-form-item>
        <el-form-item label="path" v-if="dialogData.vmessInfo.network==='ws' || (dialogData.vmessInfo.network==='tcp' && dialogData.vmessInfo.type ==='http')">
          <el-input v-model="dialogData.vmessInfo.path"/>
        </el-form-item>

        <el-form-item label="security">
          <el-radio-group v-model="dialogData.vmessInfo.security">
            <el-radio label="none">none</el-radio>
            <el-radio label="tls">tls</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="sni" v-if="dialogData.vmessInfo.security==='tls'">
          <el-input v-model="dialogData.vmessInfo.sni"/>
        </el-form-item>
        <el-form-item label="fp" v-if="dialogData.vmessInfo.security==='tls'">
          <el-select
              v-model="dialogData.vmessInfo.fp"
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

        <el-divider content-position="left">其他参数</el-divider>
        <el-form :model="dialogData.vmessInfo" label-width="100px">
          <el-form-item label="是否启用">
            <el-switch
                size="small"
                v-model="dialogData.vmessInfo.enabled"
                style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
            />
          </el-form-item>
          <el-form-item label="节点限速">
            <el-input type="number" v-model.number="dialogData.vmessInfo.node_speedlimit"/>
          </el-form-item>
          <el-form-item label="节点倍率">
            <el-input type="number" v-model.number="dialogData.vmessInfo.traffic_rate"/>
          </el-form-item>
          <el-form-item label="启用中转">
            <el-switch
                size="small"
                v-model="dialogData.vmessInfo.enable_transfer"
                style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
            />
          </el-form-item>
          <el-form-item label="中转ip" v-if="dialogData.vmessInfo.enable_transfer">
            <el-input v-model="dialogData.vmessInfo.transfer_address" placeholder=""/>
          </el-form-item>
          <el-form-item label="中转端口" v-if="dialogData.vmessInfo.enable_transfer">
            <el-input v-model="dialogData.vmessInfo.transfer_port" placeholder=""/>
          </el-form-item>
        </el-form>
      </el-form>
    </div>

    <div v-if="state.noteType === 'shadowsocks'">
      <el-form :model="dialogData.shadowsocksInfo" label-width="100px">
        <el-form-item label="remarks">
          <el-input v-model="dialogData.shadowsocksInfo.remarks"/>
        </el-form-item>
        <el-form-item label="address">
          <el-input v-model="dialogData.shadowsocksInfo.address"/>
        </el-form-item>
        <el-form-item label="port">
          <el-input v-model.number="dialogData.shadowsocksInfo.port"/>
        </el-form-item>
        <el-form-item label="scy">
          <el-select
              v-model="dialogData.shadowsocksInfo.scy"
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
        <el-form-item label="type">
          <el-select
              v-model="dialogData.shadowsocksInfo.type"
              filterable
              allow-create
              default-first-option
              :reserve-keyword="false"
              style="width: 100%"
          >
            <el-option
                v-for="(v,k) in state.typeArr1"
                :key="k"
                :label="v"
                :value="v">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="host" v-if="dialogData.shadowsocksInfo.type === 'http'">
          <el-input v-model="dialogData.shadowsocksInfo.host"/>
        </el-form-item>
        <el-form-item label="path" v-if="dialogData.shadowsocksInfo.type === 'http'">
          <el-input v-model="dialogData.shadowsocksInfo.path"/>
        </el-form-item>

        <el-divider content-position="left">其他参数</el-divider>
        <el-form-item label="是否启用">
          <el-switch
              size="small"
              v-model="dialogData.shadowsocksInfo.enabled"
              style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
          />
        </el-form-item>
        <el-form-item label="节点限速">
          <el-input type="number" v-model.number="dialogData.shadowsocksInfo.node_speedlimit"/>
        </el-form-item>
        <el-form-item label="节点倍率">
          <el-input type="number" v-model.number="dialogData.shadowsocksInfo.traffic_rate"/>
        </el-form-item>
        <el-form-item label="启用中转">
          <el-switch
              size="small"
              v-model="dialogData.shadowsocksInfo.enable_transfer"
              style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
          />
        </el-form-item>
        <el-form-item label="中转ip" v-if="dialogData.shadowsocksInfo.enable_transfer">
          <el-input v-model="dialogData.shadowsocksInfo.transfer_address" placeholder=""/>
        </el-form-item>
        <el-form-item label="中转端口" v-if="dialogData.shadowsocksInfo.enable_transfer">
          <el-input v-model="dialogData.shadowsocksInfo.transfer_port" placeholder=""/>
        </el-form-item>
      </el-form>
    </div>

    <template #footer>
      <el-button @click="closeDialog">取消</el-button>
      <el-button @click="onSubmit" type="danger">保存</el-button>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>

import {storeToRefs} from "pinia";
import {useNodeStore} from "/@/stores/nodeStore";
import {reactive, watch} from "vue";
import {useApiStore} from "/@/stores/apiStore";
import {request} from "/@/utils/request";
import {deepClone} from "/@/utils/other";

const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)


const nodeStore = useNodeStore()
const {dialogData} = storeToRefs(nodeStore)
const emit = defineEmits(['refresh']);
const state = reactive({
  title: "",
  noteType: 'vless',
  isShowDialog: false,
  nodeTypelist: ["vless", "vmess", "shadowsocks"],
  realityDefaultArr: [
    {dest: "www.speedtest.org:443", sni: "www.speedtest.org"},
    {dest: "www.lovelive-anime.jp:443", sni: "www.lovelive-anime.jp"},
    {dest: "swdist.apple.com:443", sni: "swdist.apple.com"},
    {dest: "blog.api.www.cloudflare.com:443", sni: "blog.api.www.cloudflare.com"},
    {dest: "www.icloud.com:443", sni: "www.icloud.com"},
  ] as RealityItem[],
  scyArr: ["auto", "none", "chacha20-poly1305", "aes-128-gcm", "aes-256-gcm", "2022-blake3-aes-128-gcm", "2022-blake3-aes-256-gcm", "2022-blake3-chacha20-poly1305",],
  flowArr: ["", "xtls-rprx-vision", "xtls-rprx-vision-udp443",],
  networkArr: ["ws", "tcp", "grpc",],
  typeArr1: ["none", "http"],
  typeArr2: ["none", "srtp", "utp", "wechat-video", "dtls", "wireguard"],
  fpArr: ["chrome", "firefox", "safari", "ios", "android", "edge", "360", "qq", "random", "randomized",],
})

// 打开弹窗
const openDialog = (title: string, row?: NodeInfo) => {

  if (title === '新建节点') {
    dialogData.value.vlessInfo.id = 0 //编辑和添加公用一个store，清空id,否则服务器无法插入
    dialogData.value.vmessInfo.id = 0 //编辑和添加公用一个store，清空id,否则服务器无法插入
    dialogData.value.shadowsocksInfo.id = 0 //编辑和添加公用一个store，清空id,否则服务器无法插入
    state.title = "新建节点"
    state.isShowDialog = true
  } else {
    state.title = "修改节点"
    switch (row?.node_type) {
      case "vless":
        state.noteType="vless"
        dialogData.value.vlessInfo = row
            break
      case "vmess":
        state.noteType="vmess"
        dialogData.value.vmessInfo = row
            break
      case "shadowsocks":
        state.noteType="shadowsocks"
        dialogData.value.shadowsocksInfo = row
            break
    }
    state.isShowDialog = true
  }
}
// 关闭弹窗
const closeDialog = () => {
  state.isShowDialog = false
};

//确认提交
function onSubmit() {
  if (state.title === '新建节点') {
    // let n = nodeStore.returnNodeInfo(state.noteType)
    // console.log("确认提交n:",n)
    nodeStore.newNode(nodeStore.returnNodeInfo(state.noteType))
    setTimeout(() => {
      emit('refresh');
    }, 500);
  } else {
    nodeStore.updateNode(nodeStore.returnNodeInfo(state.noteType))
    setTimeout(() => {
      emit('refresh');
    }, 500);
  }
  closeDialog()
}

//
const setReality = (nodeType: string) => {
  request(apiStoreData.api.value.system_createx25519).then((res) => {
    // console.log("res:",res)
    switch (nodeType) {
      case "vless":
        dialogData.value.vlessInfo.pbk = res.data.public_key
        dialogData.value.vlessInfo.private_key = res.data.private_key
        break
    }
  })
}
//监听
watch(
    () => dialogData.value.vlessInfo.dest,
    () => {
      let temp = state.realityDefaultArr.filter(r => r.dest === dialogData.value.vlessInfo.dest)
      if (temp.length > 0) {
        dialogData.value.vlessInfo.sni = temp[0].sni
      }
      if (dialogData.value.vlessInfo.private_key === '' || dialogData.value.vlessInfo.private_key === '') {
        setReality('vless')
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
