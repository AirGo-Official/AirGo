<template>
  <el-dialog v-model="state.isShowDialog" :title="state.title" width="80%" destroy-on-close align-center>
    <div v-if="state.title==='新建节点'">
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

    <div v-if="state.noteType === constantStore.NODE_TYPE_TRANSFER">
      <el-form :model="nodeStoreData.dialogData.value.transferInfo" label-width="100px">
        <el-form-item label="remarks">
          <el-input v-model="nodeStoreData.dialogData.value.transferInfo.remarks"/>
        </el-form-item>
        <el-form-item label="address">
          <el-input v-model="nodeStoreData.dialogData.value.transferInfo.transfer_address"/>
        </el-form-item>
        <el-form-item label="port">
          <el-input-number v-model="nodeStoreData.dialogData.value.transferInfo.transfer_port"/>
        </el-form-item>
        <el-form-item label="node">
          <el-select v-model="nodeStoreData.dialogData.value.transferInfo.transfer_node_id" class="m-2" placeholder="Select">
            <el-option
                v-for="item in nodeStoreData.nodeList.value.data"
                :key="item.id"
                :label="item.remarks"
                :value="item.id"
            />
          </el-select>
        </el-form-item>
      </el-form>

    </div>

    <div v-if="state.noteType === constantStore.NODE_TYPE_VLESS">
      <el-descriptions
        style="margin-bottom: 20px"
        :column="1"
        border
      >
        <el-descriptions-item label="ID">{{ nodeStoreData.dialogData.value.vlessInfo.id }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ DateStrToTime(nodeStoreData.dialogData.value.vlessInfo.created_at) }}</el-descriptions-item>
      </el-descriptions>
      <el-form :model="nodeStoreData.dialogData.value.vlessInfo" label-width="100px">
        <el-form-item label="remarks">
          <el-input v-model="nodeStoreData.dialogData.value.vlessInfo.remarks"/>
        </el-form-item>
        <el-form-item label="address">
          <el-input v-model="nodeStoreData.dialogData.value.vlessInfo.address"/>
        </el-form-item>
        <el-form-item label="port">
          <el-input v-model.number="nodeStoreData.dialogData.value.vlessInfo.port"/>
        </el-form-item>
        <el-form-item label="flow">
          <el-select
              v-model="nodeStoreData.dialogData.value.vlessInfo.flow"
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
              v-model="nodeStoreData.dialogData.value.vlessInfo.network"
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
        <el-form-item label="type" v-if="nodeStoreData.dialogData.value.vlessInfo.network === 'tcp'">
          <el-select
              v-model="nodeStoreData.dialogData.value.vlessInfo.type"
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

        <el-form-item label="mode" v-if="nodeStoreData.dialogData.value.vlessInfo.network==='grpc'">
          <el-radio-group v-model="nodeStoreData.dialogData.value.vlessInfo.mode">
            <el-radio label="gun">gun</el-radio>
            <el-radio label="multi">multi</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="serviceName" v-if="nodeStoreData.dialogData.value.vlessInfo.network==='grpc'">
          <el-input v-model="nodeStoreData.dialogData.value.vlessInfo.service_name"/>
        </el-form-item>

        <el-form-item label="host"
                      v-if="nodeStoreData.dialogData.value.vlessInfo.network==='ws' || 
                      (nodeStoreData.dialogData.value.vlessInfo.network==='tcp' && 
                      nodeStoreData.dialogData.value.vlessInfo.type ==='http')">
          <el-input v-model="nodeStoreData.dialogData.value.vlessInfo.host"/>
        </el-form-item>
        <el-form-item label="path"
                      v-if="nodeStoreData.dialogData.value.vlessInfo.network==='ws' || 
                      (nodeStoreData.dialogData.value.vlessInfo.network==='tcp' && 
                      nodeStoreData.dialogData.value.vlessInfo.type ==='http')">
          <el-input v-model="nodeStoreData.dialogData.value.vlessInfo.path"/>
        </el-form-item>
        <el-form-item label="security">
          <el-radio-group v-model="nodeStoreData.dialogData.value.vlessInfo.security">
            <el-radio label="none">none</el-radio>
            <el-radio label="tls">tls</el-radio>
            <el-radio label="reality">reality</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="dest" v-if="nodeStoreData.dialogData.value.vlessInfo.security==='reality'">
          <el-select
              v-model="nodeStoreData.dialogData.value.vlessInfo.dest"
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
                      v-if="nodeStoreData.dialogData.value.vlessInfo.security==='tls' || nodeStoreData.dialogData.value.vlessInfo.security==='reality'">
          <el-input v-model="nodeStoreData.dialogData.value.vlessInfo.sni"/>
        </el-form-item>
        <el-form-item label="fp"
                      v-if="nodeStoreData.dialogData.value.vlessInfo.security==='tls' ||nodeStoreData.dialogData.value.vlessInfo.security==='reality'">
          <el-select
              v-model="nodeStoreData.dialogData.value.vlessInfo.fp"
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
        <el-form-item label="public_key" v-if="nodeStoreData.dialogData.value.vlessInfo.security==='reality'">
          <el-input v-model="nodeStoreData.dialogData.value.vlessInfo.pbk">
            <template #append>
              <el-button @click="setReality('vless')">
                <el-icon>
                  <Refresh/>
                </el-icon>
              </el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="private_key" v-if="nodeStoreData.dialogData.value.vlessInfo.security==='reality'">
          <el-input v-model="nodeStoreData.dialogData.value.vlessInfo.private_key"/>
        </el-form-item>
        <el-form-item label="allowInsecure">
          <el-switch
              size="small"
              v-model="nodeStoreData.dialogData.value.vlessInfo.allowInsecure"
              style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
          />
        </el-form-item>
        <el-divider content-position="left">其他参数</el-divider>
        <el-form :model="nodeStoreData.dialogData.value.vlessInfo" label-width="100px">
          <el-form-item label="是否启用">
            <el-switch
                size="small"
                v-model="nodeStoreData.dialogData.value.vlessInfo.enabled"
                style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
            />
          </el-form-item>
          <el-form-item label="启用中转">
            <el-switch
                size="small"
                v-model="nodeStoreData.dialogData.value.vlessInfo.enable_transfer"
                style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
            />
          </el-form-item>
          <el-form-item label="中转ip" v-if="nodeStoreData.dialogData.value.vlessInfo.enable_transfer">
            <el-input v-model="nodeStoreData.dialogData.value.vlessInfo.transfer_address"/>
          </el-form-item>
          <el-form-item label="中转端口" v-if="nodeStoreData.dialogData.value.vlessInfo.enable_transfer">
            <el-input v-model.number="nodeStoreData.dialogData.value.vlessInfo.transfer_port"/>
          </el-form-item>
          <el-form-item label="节点限速">
            <el-input type="number" v-model.number="nodeStoreData.dialogData.value.vlessInfo.node_speed_limit"/>
          </el-form-item>
          <el-form-item label="节点倍率">
            <el-input type="number" v-model.number="nodeStoreData.dialogData.value.vlessInfo.traffic_rate"/>
          </el-form-item>
          <el-form-item label="访问控制">
            <el-transfer
                :data="nodeStoreData.accessList.value.data"
                v-model="nodeStoreData.dialogData.value.checkedAccessIDs"
                :props="{key: 'id',label: 'name',}"
                :titles="['全部', '选中']"
            />
          </el-form-item>
        </el-form>
      </el-form>
    </div>

    <div v-if="state.noteType === constantStore.NODE_TYPE_VMESS">
      <el-descriptions
        style="margin-bottom: 20px"
        :column="1"
        border
      >
        <el-descriptions-item label="ID">{{ nodeStoreData.dialogData.value.vmessInfo.id }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ DateStrToTime(nodeStoreData.dialogData.value.vmessInfo.created_at) }}</el-descriptions-item>
      </el-descriptions>
      <el-form :model="nodeStoreData.dialogData.value.vmessInfo" label-width="100px">
        <el-form-item label="remarks">
          <el-input v-model="nodeStoreData.dialogData.value.vmessInfo.remarks"/>
        </el-form-item>
        <el-form-item label="address">
          <el-input v-model="nodeStoreData.dialogData.value.vmessInfo.address"/>
        </el-form-item>
        <el-form-item label="port">
          <el-input v-model.number="nodeStoreData.dialogData.value.vmessInfo.port"/>
        </el-form-item>
        <el-form-item label="scy">
          <el-select
              v-model="nodeStoreData.dialogData.value.vmessInfo.scy"
              filterable
              allow-create
              default-first-option
              :reserve-keyword="false"
              style="width: 100%"
          >
            <el-option
                v-for="(v,k) in state.scyArrForVmess"
                :key="k"
                :label="v"
                :value="v">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="network">
          <el-select
              v-model="nodeStoreData.dialogData.value.vmessInfo.network"
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
        <el-form-item label="type" v-if="nodeStoreData.dialogData.value.vmessInfo.network === 'tcp'">
          <el-select
              v-model="nodeStoreData.dialogData.value.vmessInfo.type"
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

        <el-form-item label="mode" v-if="nodeStoreData.dialogData.value.vmessInfo.network==='grpc'">
          <el-radio-group v-model="nodeStoreData.dialogData.value.vmessInfo.mode">
            <el-radio label="gun">gun</el-radio>
            <el-radio label="multi">multi</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="serviceName" v-if="nodeStoreData.dialogData.value.vmessInfo.network==='grpc'">
          <el-input v-model="nodeStoreData.dialogData.value.vmessInfo.service_name"/>
        </el-form-item>
        <el-form-item label="host"
                      v-if="nodeStoreData.dialogData.value.vmessInfo.network==='ws' || (nodeStoreData.dialogData.value.vmessInfo.network==='tcp' && nodeStoreData.dialogData.value.vmessInfo.type ==='http')">
          <el-input v-model="nodeStoreData.dialogData.value.vmessInfo.host"/>
        </el-form-item>
        <el-form-item label="path"
                      v-if="nodeStoreData.dialogData.value.vmessInfo.network==='ws' || (nodeStoreData.dialogData.value.vmessInfo.network==='tcp' && nodeStoreData.dialogData.value.vmessInfo.type ==='http')">
          <el-input v-model="nodeStoreData.dialogData.value.vmessInfo.path"/>
        </el-form-item>

        <el-form-item label="security">
          <el-radio-group v-model="nodeStoreData.dialogData.value.vmessInfo.security">
            <el-radio label="none">none</el-radio>
            <el-radio label="tls">tls</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="sni" v-if="nodeStoreData.dialogData.value.vmessInfo.security==='tls'">
          <el-input v-model="nodeStoreData.dialogData.value.vmessInfo.sni"/>
        </el-form-item>
        <el-form-item label="fp" v-if="nodeStoreData.dialogData.value.vmessInfo.security==='tls'">
          <el-select
              v-model="nodeStoreData.dialogData.value.vmessInfo.fp"
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
        <el-form-item label="allowInsecure">
          <el-switch
              size="small"
              v-model="nodeStoreData.dialogData.value.vmessInfo.allowInsecure"
              style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
          />
        </el-form-item>

        <el-divider content-position="left">其他参数</el-divider>
        <el-form :model="nodeStoreData.dialogData.value.vmessInfo" label-width="100px">
          <el-form-item label="是否启用">
            <el-switch
                size="small"
                v-model="nodeStoreData.dialogData.value.vmessInfo.enabled"
                style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
            />
          </el-form-item>
          <el-form-item label="启用中转">
            <el-switch
                size="small"
                v-model="nodeStoreData.dialogData.value.vmessInfo.enable_transfer"
                style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
            />
          </el-form-item>
          <el-form-item label="中转ip" v-if="nodeStoreData.dialogData.value.vmessInfo.enable_transfer">
            <el-input v-model="nodeStoreData.dialogData.value.vmessInfo.transfer_address"/>
          </el-form-item>
          <el-form-item label="中转端口" v-if="nodeStoreData.dialogData.value.vmessInfo.enable_transfer">
            <el-input v-model.number="nodeStoreData.dialogData.value.vmessInfo.transfer_port"/>
          </el-form-item>
          <el-form-item label="节点限速">
            <el-input type="number" v-model.number="nodeStoreData.dialogData.value.vmessInfo.node_speed_limit"/>
          </el-form-item>
          <el-form-item label="节点倍率">
            <el-input type="number" v-model.number="nodeStoreData.dialogData.value.vmessInfo.traffic_rate"/>
          </el-form-item>
          <el-form-item label="访问控制">
            <el-transfer
                :data="nodeStoreData.accessList.value.data"
                v-model="nodeStoreData.dialogData.value.checkedAccessIDs"
                :props="{
                  key: 'id',
                  label: 'name',
                  }"
                :titles="['全部', '选中']"
            />
          </el-form-item>
        </el-form>
      </el-form>
    </div>

    <div v-if="state.noteType === constantStore.NODE_TYPE_SHADOWSOCKS">
      <el-descriptions
        style="margin-bottom: 20px"
        :column="1"
        border
      >
        <el-descriptions-item label="ID">{{ nodeStoreData.dialogData.value.shadowsocksInfo.id }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ DateStrToTime(nodeStoreData.dialogData.value.shadowsocksInfo.created_at) }}</el-descriptions-item>
      </el-descriptions>
      <el-form :model="nodeStoreData.dialogData.value.shadowsocksInfo" label-width="100px">
        <el-form-item label="remarks">
          <el-input v-model="nodeStoreData.dialogData.value.shadowsocksInfo.remarks"/>
        </el-form-item>
        <el-form-item label="address">
          <el-input v-model="nodeStoreData.dialogData.value.shadowsocksInfo.address"/>
        </el-form-item>
        <el-form-item label="port">
          <el-input v-model.number="nodeStoreData.dialogData.value.shadowsocksInfo.port"/>
        </el-form-item>
        <el-form-item label="scy">
          <el-select
              v-model="nodeStoreData.dialogData.value.shadowsocksInfo.scy"
              filterable
              allow-create
              default-first-option
              :reserve-keyword="false"
              style="width: 100%"
          >
            <el-option
                v-for="(v,k) in state.scyArrForSS"
                :key="k"
                :label="v"
                :value="v">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="type">
          <el-select
              v-model="nodeStoreData.dialogData.value.shadowsocksInfo.type"
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
        <el-form-item label="host" v-if="nodeStoreData.dialogData.value.shadowsocksInfo.type === 'http'">
          <el-input v-model="nodeStoreData.dialogData.value.shadowsocksInfo.host"/>
        </el-form-item>
        <el-form-item label="path" v-if="nodeStoreData.dialogData.value.shadowsocksInfo.type === 'http'">
          <el-input v-model="nodeStoreData.dialogData.value.shadowsocksInfo.path"/>
        </el-form-item>

        <el-divider content-position="left">其他参数</el-divider>
        <el-form-item label="是否启用">
          <el-switch
              size="small"
              v-model="nodeStoreData.dialogData.value.shadowsocksInfo.enabled"
              style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
          />
        </el-form-item>
        <el-form-item label="启用中转">
          <el-switch
              size="small"
              v-model="nodeStoreData.dialogData.value.shadowsocksInfo.enable_transfer"
              style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
          />
        </el-form-item>
        <el-form-item label="中转ip" v-if="nodeStoreData.dialogData.value.shadowsocksInfo.enable_transfer">
          <el-input v-model="nodeStoreData.dialogData.value.shadowsocksInfo.transfer_address"/>
        </el-form-item>
        <el-form-item label="中转端口" v-if="nodeStoreData.dialogData.value.shadowsocksInfo.enable_transfer">
          <el-input v-model.number="nodeStoreData.dialogData.value.shadowsocksInfo.transfer_port"/>
        </el-form-item>
        <el-form-item label="节点限速">
          <el-input type="number" v-model.number="nodeStoreData.dialogData.value.shadowsocksInfo.node_speed_limit"/>
        </el-form-item>
        <el-form-item label="节点倍率">
          <el-input type="number" v-model.number="nodeStoreData.dialogData.value.shadowsocksInfo.traffic_rate"/>
        </el-form-item>
        <el-form-item label="访问控制">
          <el-transfer
              :data="nodeStoreData.accessList.value.data"
              v-model="nodeStoreData.dialogData.value.checkedAccessIDs"
              :props="{
                  key: 'id',
                  label: 'name',
                  }"
              :titles="['全部', '选中']"
          />
        </el-form-item>
      </el-form>
    </div>

    <div v-if="state.noteType === constantStore.NODE_TYPE_HYSTERIA">
      <el-descriptions
        style="margin-bottom: 20px"
        :column="1"
        border
      >
        <el-descriptions-item label="ID">{{ nodeStoreData.dialogData.value.hysteriaInfo.id }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ DateStrToTime(nodeStoreData.dialogData.value.hysteriaInfo.created_at) }}</el-descriptions-item>
      </el-descriptions>
      <el-form :model="nodeStoreData.dialogData.value.hysteriaInfo" label-width="100px">
        <el-form-item label="remarks">
          <el-input v-model="nodeStoreData.dialogData.value.hysteriaInfo.remarks"/>
        </el-form-item>
        <el-form-item label="address">
          <el-input v-model="nodeStoreData.dialogData.value.hysteriaInfo.address"/>
        </el-form-item>
        <el-form-item label="port">
          <el-input v-model.number="nodeStoreData.dialogData.value.hysteriaInfo.port"/>
        </el-form-item>
        <el-form-item label="sni">
          <el-input v-model.number="nodeStoreData.dialogData.value.hysteriaInfo.sni"/>
        </el-form-item>
        <el-form-item label="allowInsecure">
          <el-switch
              size="small"
              v-model="nodeStoreData.dialogData.value.hysteriaInfo.allowInsecure"
              style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
          />
        </el-form-item>
        <el-divider content-position="left">其他参数</el-divider>
        <el-form-item label="是否启用">
          <el-switch
              size="small"
              v-model="nodeStoreData.dialogData.value.hysteriaInfo.enabled"
              style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
          />
        </el-form-item>
        <el-form-item label="启用中转">
          <el-switch
              size="small"
              v-model="nodeStoreData.dialogData.value.hysteriaInfo.enable_transfer"
              style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
          />
        </el-form-item>
        <el-form-item label="中转ip" v-if="nodeStoreData.dialogData.value.hysteriaInfo.enable_transfer">
          <el-input v-model="nodeStoreData.dialogData.value.hysteriaInfo.transfer_address"/>
        </el-form-item>
        <el-form-item label="中转端口" v-if="nodeStoreData.dialogData.value.hysteriaInfo.enable_transfer">
          <el-input v-model.number="nodeStoreData.dialogData.value.hysteriaInfo.transfer_port"/>
        </el-form-item>
        <el-form-item label="节点限速">
          <el-input-number v-model.number="nodeStoreData.dialogData.value.hysteriaInfo.node_speed_limit"/>
        </el-form-item>
        <el-form-item label="节点倍率">
          <el-input-number v-model.number="nodeStoreData.dialogData.value.hysteriaInfo.traffic_rate"/>
        </el-form-item>
        <el-form-item label="访问控制">
          <el-transfer
              :data="nodeStoreData.accessList.value.data"
              v-model="nodeStoreData.dialogData.value.checkedAccessIDs"
              :props="{
                  key: 'id',
                  label: 'name',
                  }"
              :titles="['全部', '选中']"
          />
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
import {useAdminNodeStore} from "/@/stores/admin_logic/nodeStore";
import {reactive, watch} from "vue";
import {useApiStore} from "/@/stores/apiStore";
import {request} from "/@/utils/request";
import { useConstantStore } from "/@/stores/constantStore";
import { DateStrToTime } from "../../../utils/formatTime";
const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
const nodeStore = useAdminNodeStore()
const nodeStoreData = storeToRefs(nodeStore)
const emit = defineEmits(['refresh']);
const constantStore = useConstantStore()
const state = reactive({
  title: "",
  noteType: constantStore.NODE_TYPE_VLESS,
  isShowDialog: false,
  nodeTypelist: [
    constantStore.NODE_TYPE_VLESS,
    constantStore.NODE_TYPE_VMESS,
    constantStore.NODE_TYPE_SHADOWSOCKS,
    constantStore.NODE_TYPE_HYSTERIA,
    constantStore.NODE_TYPE_TRANSFER],
  realityDefaultArr: [
    {dest: "www.speedtest.org:443", sni: "www.speedtest.org"},
    {dest: "www.lovelive-anime.jp:443", sni: "www.lovelive-anime.jp"},
    {dest: "swdist.apple.com:443", sni: "swdist.apple.com"},
    {dest: "blog.api.www.cloudflare.com:443", sni: "blog.api.www.cloudflare.com"},
    {dest: "www.icloud.com:443", sni: "www.icloud.com"},
  ] as RealityItem[],
  scyArrForVmess: [
    "auto",
    "none",
    "aes-128-gcm",
    "aes-256-gcm",
    "chacha20-ietf-poly1305",],
  scyArrForSS: [
    "aes-128-gcm",
    "aes-256-gcm",
    "chacha20-ietf-poly1305",
    "2022-blake3-aes-128-gcm",
    "2022-blake3-aes-256-gcm",
    "2022-blake3-chacha20-poly1305",],
  scyArrForClashMeta: [
    "aes-128-gcm",
    "aes-192-gcm",
    "aes-256-gcm",
    "chacha20-ietf-poly1305",
    "2022-blake3-aes-128-gcm",
    "2022-blake3-aes-256-gcm",
    "2022-blake3-chacha20-poly1305"],
  flowArr: ["", "xtls-rprx-vision", "xtls-rprx-vision-udp443",],
  networkArr: ["ws", "tcp", "grpc",],
  typeArr1: ["none", "http"],
  typeArr2: ["none", "srtp", "utp", "wechat-video", "dtls", "wireguard"],
  fpArr: ["chrome", "firefox", "safari", "ios", "android", "edge", "360", "qq", "random", "randomized",],
})

// 打开弹窗
const openDialog = (title: string, row?: NodeInfo) => {
  if (title === '新建节点') {
    nodeStoreData.dialogData.value.vlessInfo.id = 0 //编辑和添加公用一个store，清空id,否则服务器无法插入
    nodeStoreData.dialogData.value.vmessInfo.id = 0 //编辑和添加公用一个store，清空id,否则服务器无法插入
    nodeStoreData.dialogData.value.shadowsocksInfo.id = 0 //编辑和添加公用一个store，清空id,否则服务器无法插入
    nodeStoreData.dialogData.value.hysteriaInfo.id = 0 //编辑和添加公用一个store，清空id,否则服务器无法插入
    nodeStoreData.dialogData.value.transferInfo.id = 0 //编辑和添加公用一个store，清空id,否则服务器无法插入
    state.title = "新建节点"
    state.isShowDialog = true
  } else {
    state.title = "修改节点"
    if (row?.enable_transfer && row?.transfer_node_id!==0){
      state.noteType = constantStore.NODE_TYPE_TRANSFER
      nodeStoreData.dialogData.value.transferInfo = row
      state.isShowDialog = true
      return
    }
    switch (row?.node_type) {
      case constantStore.NODE_TYPE_VLESS:
        state.noteType = constantStore.NODE_TYPE_VLESS
        nodeStoreData.dialogData.value.vlessInfo = row
        break
      case constantStore.NODE_TYPE_VMESS:
        state.noteType = constantStore.NODE_TYPE_VMESS
        nodeStoreData.dialogData.value.vmessInfo = row
        break
      case constantStore.NODE_TYPE_SHADOWSOCKS:
        state.noteType = constantStore.NODE_TYPE_SHADOWSOCKS
        nodeStoreData.dialogData.value.shadowsocksInfo = row
        break
      case constantStore.NODE_TYPE_HYSTERIA:
        state.noteType = constantStore.NODE_TYPE_HYSTERIA
        nodeStoreData.dialogData.value.hysteriaInfo = row
        break
    }
    nodeStore.accessHandler(row)
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
  request(apiStoreData.adminApi.value.createx25519).then((res) => {
    switch (nodeType) {
      case constantStore.NODE_TYPE_VLESS:
        nodeStoreData.dialogData.value.vlessInfo.pbk = res.data.public_key
        nodeStoreData.dialogData.value.vlessInfo.private_key = res.data.private_key
        break
    }
  })
}
//监听
watch(
    () => nodeStoreData.dialogData.value.vlessInfo.dest,
    () => {
      let temp = state.realityDefaultArr.filter(r => r.dest === nodeStoreData.dialogData.value.vlessInfo.dest)
      if (temp.length > 0) {
        nodeStoreData.dialogData.value.vlessInfo.sni = temp[0].sni
      }
      if (nodeStoreData.dialogData.value.vlessInfo.private_key === '' || nodeStoreData.dialogData.value.vlessInfo.private_key === '') {
        setReality(constantStore.NODE_TYPE_VLESS)
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
