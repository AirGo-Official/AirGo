<template>
  <el-dialog v-model="state.isShowDialog" :title="state.title" width="90%" destroy-on-close align-center>
    <!--    节点类型、节点协议类型-->
    <div v-if="state.type === 'add'">
      <el-divider content-position="left">{{ $t("message.adminNode.NodeInfo.node_type") }}</el-divider>
      <el-form label-position="top">
        <el-form-item :label="$t('message.adminNode.NodeInfo.node_type')">
          <el-select
            v-model="state.nodeType"
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
        <el-form-item :label="$t('message.adminNode.NodeInfo.protocol')"
                      v-if="state.nodeType === constantStore.NODE_TYPE_NORMAL"
        >
          <el-select
            v-model="state.nodeProtocol"
            filterable
            allow-create
            default-first-option
            :reserve-keyword="false"
            style="width: 100%"
          >
            <el-option
              v-for="(v,k) in state.nodeProtocolList"
              :key="k"
              :label="v"
              :value="v">
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
    </div>

    <el-divider content-position="left">{{ $t("message.adminNode.basicParameters") }}</el-divider>

    <!--    具体的参数-->

    <div v-if="state.nodeType === constantStore.NODE_TYPE_TRANSFER">
      <el-form :model="nodeStoreData.dialogData.value.transferInfo" label-position="top">
        <el-form-item :label="$t('message.adminNode.NodeInfo.remarks')">
          <el-input v-model="nodeStoreData.dialogData.value.transferInfo.remarks" />
        </el-form-item>
        <el-form-item :label="$t('message.adminNode.NodeInfo.transfer_address')">
          <el-input v-model="nodeStoreData.dialogData.value.transferInfo.transfer_address" />
        </el-form-item>
        <el-form-item :label="$t('message.adminNode.NodeInfo.transfer_port')">
          <el-input-number v-model="nodeStoreData.dialogData.value.transferInfo.transfer_port" />
        </el-form-item>
        <el-form-item :label="$t('message.adminNode.NodeInfo.transfer_node_id')">
          <el-select v-model="nodeStoreData.dialogData.value.transferInfo.transfer_node_id" class="m-2"
                     placeholder="Select">
            <el-option
              v-for="item in nodeStoreData.nodeList.value.data"
              :key="item.id"
              :label="item.remarks"
              :value="item.id">
              <span style="float: left">id: </span>
              <span>{{ item.id }}</span>
              <span style="margin-left: 20px">{{ item.remarks }}</span>
              <span style="float: right">
            <el-tag type="success" style="width: 100px"
                    v-if="item.node_type ===constantStore.NODE_TYPE_NORMAL">{{ $t("message.constant.NODE_TYPE_NORMAL") }}</el-tag>
            <el-tag type="warning" style="width: 100px" v-else-if="item.node_type ===constantStore.NODE_TYPE_TRANSFER">{{ $t("message.constant.NODE_TYPE_TRANSFER") }}</el-tag>
            <el-tag type="info" style="width: 100px"
                    v-else-if="item.node_type ===constantStore.NODE_TYPE_SHARED">{{ $t("message.constant.NODE_TYPE_SHARED") }}</el-tag>
              </span>
              <span style="float: right">
                <el-button type="success" v-if="item.protocol ===constantStore.NODE_PROTOCOL_VMESS">{{ $t("message.constant.NODE_PROTOCOL_VMESS") }}</el-button>
                <el-button type="warning"
                           v-else-if="item.protocol ===constantStore.NODE_PROTOCOL_VLESS">{{ $t("message.constant.NODE_PROTOCOL_VLESS") }}</el-button>
                <el-button type="info"
                           v-else-if="item.protocol ===constantStore.NODE_PROTOCOL_TROJAN">{{ $t("message.constant.NODE_PROTOCOL_TROJAN") }}</el-button>
                <el-button type="danger" v-else-if="item.protocol ===constantStore.NODE_PROTOCOL_SHADOWSOCKS">{{ $t("message.constant.NODE_PROTOCOL_SHADOWSOCKS") }}</el-button>
                <el-button type="primary"
                           v-else-if="item.protocol ===constantStore.NODE_PROTOCOL_HYSTERIA">{{ $t("message.constant.NODE_PROTOCOL_HYSTERIA") }}</el-button>
              </span>
            </el-option>
          </el-select>
        </el-form-item>
        <el-divider content-position="left">{{ $t("message.adminNode.otherParameters") }}</el-divider>
        <el-form-item :label="$t('message.adminNode.NodeInfo.enabled')">
          <el-switch
            size="small"
            inline-prompt
            v-model="nodeStoreData.dialogData.value.transferInfo.enabled"
            style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
            :active-text="$t('message.common.display')"
            :inactive-text="$t('message.common.hide')"
          />
        </el-form-item>
      </el-form>
    </div>
    <div v-else>
      <div v-if="state.nodeType === constantStore.NODE_TYPE_SHARED && state.type==='add'">
        <el-input v-model="nodeStoreData.nodeSharedData.value.newNodeSharedUrl.url"
                  size="default"
                  type="textarea"
                  autosize
                  :placeholder="$t('message.adminNode.placeholder1')">
        </el-input>
        <el-button style="margin-top: 10px;margin-bottom: 10px;width: 100px" type="primary" @click="parseUrl">
          <el-icon>
            <Search />
          </el-icon>
          {{ $t("message.adminNode.parse") }}
        </el-button>

        <el-table :data="nodeStoreData.nodeSharedData.value.nodeList.data" height="100%" stripe border
                  style="width: 100%;flex: 1;">
          <!--        <el-table-column prop="id" label="ID" show-overflow-tooltip width="40"></el-table-column>-->
          <el-table-column prop="protocol" :label="$t('message.adminNode.NodeInfo.protocol')" show-overflow-tooltip
                           width="80"></el-table-column>
          <el-table-column prop="remarks" :label="$t('message.adminNode.NodeInfo.remarks')" show-overflow-tooltip
                           width="200"></el-table-column>
          <el-table-column prop="address" :label="$t('message.adminNode.NodeInfo.address')" show-overflow-tooltip
                           width="200"></el-table-column>
          <el-table-column prop="port" :label="$t('message.adminNode.NodeInfo.port')" show-overflow-tooltip
                           width="80"></el-table-column>
          <el-table-column prop="network" :label="$t('message.adminNode.NodeInfo.network')" show-overflow-tooltip
                           width="80"></el-table-column>
          <el-table-column :label="$t('message.common.operate')">
            <template #default="{row}">
              <el-button size="small" text type="primary" @click="deleteOneSharedNode(row)">
                {{ $t("message.common.delete") }}
              </el-button>
            </template>
          </el-table-column>

        </el-table>
      </div>
      <div v-else>
        <div v-if="state.nodeProtocol === constantStore.NODE_PROTOCOL_VLESS">
          <el-form :model="nodeStoreData.dialogData.value.vlessInfo" label-position="top">
            <el-form-item :label="$t('message.adminNode.NodeInfo.remarks')">
              <el-input v-model="nodeStoreData.dialogData.value.vlessInfo.remarks" />
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.address')">
              <el-input v-model="nodeStoreData.dialogData.value.vlessInfo.address" />
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.port')">
              <el-input-number v-model.number="nodeStoreData.dialogData.value.vlessInfo.port" />
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.flow')">
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
            <el-form-item :label="$t('message.adminNode.NodeInfo.network')">
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
            <el-form-item :label="$t('message.adminNode.NodeInfo.type')"
                          v-if="nodeStoreData.dialogData.value.vlessInfo.network === 'tcp'">
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

            <el-form-item :label="$t('message.adminNode.NodeInfo.mode')"
                          v-if="nodeStoreData.dialogData.value.vlessInfo.network==='grpc'">
              <el-radio-group v-model="nodeStoreData.dialogData.value.vlessInfo.mode">
                <el-radio label="gun">gun</el-radio>
                <el-radio label="multi">multi</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.service_name')"
                          v-if="nodeStoreData.dialogData.value.vlessInfo.network==='grpc'">
              <el-input v-model="nodeStoreData.dialogData.value.vlessInfo.service_name" />
            </el-form-item>

            <el-form-item :label="$t('message.adminNode.NodeInfo.host')"
                          v-if="nodeStoreData.dialogData.value.vlessInfo.network==='ws' ||
                      (nodeStoreData.dialogData.value.vlessInfo.network==='tcp' &&
                      nodeStoreData.dialogData.value.vlessInfo.type ==='http')">
              <el-input v-model="nodeStoreData.dialogData.value.vlessInfo.host" />
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.path')"
                          v-if="nodeStoreData.dialogData.value.vlessInfo.network==='ws' ||
                      (nodeStoreData.dialogData.value.vlessInfo.network==='tcp' &&
                      nodeStoreData.dialogData.value.vlessInfo.type ==='http')">
              <el-input v-model="nodeStoreData.dialogData.value.vlessInfo.path" />
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.security')">

              <el-radio-group v-model="nodeStoreData.dialogData.value.vlessInfo.security">
                <el-radio label="">none</el-radio>
                <el-radio label="tls">tls</el-radio>
                <el-radio label="reality">reality</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.dest')"
                          v-if="nodeStoreData.dialogData.value.vlessInfo.security==='reality'">
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
            <el-form-item :label="$t('message.adminNode.NodeInfo.sni')"
                          v-if="nodeStoreData.dialogData.value.vlessInfo.security==='tls' || nodeStoreData.dialogData.value.vlessInfo.security==='reality'">
              <el-input v-model="nodeStoreData.dialogData.value.vlessInfo.sni" />
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.fp')"
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
            <el-form-item :label="$t('message.adminNode.NodeInfo.pbk')"
                          v-if="nodeStoreData.dialogData.value.vlessInfo.security==='reality'">
              <el-input v-model="nodeStoreData.dialogData.value.vlessInfo.pbk">
                <template #append>
                  <el-button @click="setReality('vless')">
                    <el-icon>
                      <Refresh />
                    </el-icon>
                  </el-button>
                </template>
              </el-input>
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.private_key')"
                          v-if="nodeStoreData.dialogData.value.vlessInfo.security==='reality'">
              <el-input v-model="nodeStoreData.dialogData.value.vlessInfo.private_key" />
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.allowInsecure')">
              <el-switch
                size="small"
                v-model="nodeStoreData.dialogData.value.vlessInfo.allowInsecure"
                style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
              />
            </el-form-item>
            <el-divider content-position="left">{{ $t("message.adminNode.otherParameters") }}</el-divider>
            <el-form :model="nodeStoreData.dialogData.value.vlessInfo" label-position="top">
              <el-form-item :label="$t('message.adminNode.NodeInfo.enabled')">
                <el-switch
                  size="small"
                  inline-prompt
                  v-model="nodeStoreData.dialogData.value.vlessInfo.enabled"
                  style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
                  :active-text="$t('message.common.display')"
                  :inactive-text="$t('message.common.hide')"
                />
              </el-form-item>
              <div v-if="state.nodeType===constantStore.NODE_TYPE_NORMAL">
                <el-form-item :label="$t('message.adminNode.NodeInfo.transfer_address')">
                  <el-input v-model="nodeStoreData.dialogData.value.vlessInfo.transfer_address" />
                </el-form-item>
                <el-form-item :label="$t('message.adminNode.NodeInfo.transfer_port')">
                  <el-input-number v-model.number="nodeStoreData.dialogData.value.vlessInfo.transfer_port" />
                </el-form-item>
                <el-form-item :label="$t('message.adminNode.NodeInfo.node_speed_limit')">
                  <el-input-number v-model.number="nodeStoreData.dialogData.value.vlessInfo.node_speed_limit" />
                </el-form-item>
                <el-form-item :label="$t('message.adminNode.NodeInfo.traffic_rate')">
                  <el-input-number v-model.number="nodeStoreData.dialogData.value.vlessInfo.traffic_rate" />
                </el-form-item>
                <el-form-item :label="$t('message.adminNode.NodeInfo.access')">
                  <el-tree ref="access_tree_ref" node-key="id"
                           :data="nodeStoreData.accessList.value.data"
                           :props="{label:'name'}"
                           :default-checked-keys="nodeStoreData.dialogData.value.checkedAccessIDs"
                           show-checkbox class="data-tree">
                    <template #default="{ node, data }">
                  <span class="custom-tree-node">
                    <span>{{ $t(data.name) }}</span>
                  </span>
                    </template>
                  </el-tree>
                </el-form-item>
              </div>

            </el-form>
          </el-form>
        </div>

        <div v-else-if="state.nodeProtocol === constantStore.NODE_PROTOCOL_VMESS">
          <el-form :model="nodeStoreData.dialogData.value.vmessInfo" label-position="top">
            <el-form-item :label="$t('message.adminNode.NodeInfo.remarks')">
              <el-input v-model="nodeStoreData.dialogData.value.vmessInfo.remarks" />
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.address')">
              <el-input v-model="nodeStoreData.dialogData.value.vmessInfo.address" />
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.port')">
              <el-input-number v-model.number="nodeStoreData.dialogData.value.vmessInfo.port" />
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.scy')">
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
            <el-form-item :label="$t('message.adminNode.NodeInfo.network')">
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
            <el-form-item :label="$t('message.adminNode.NodeInfo.type')"
                          v-if="nodeStoreData.dialogData.value.vmessInfo.network === 'tcp'">
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

            <el-form-item :label="$t('message.adminNode.NodeInfo.mode')"
                          v-if="nodeStoreData.dialogData.value.vmessInfo.network==='grpc'">
              <el-radio-group v-model="nodeStoreData.dialogData.value.vmessInfo.mode">
                <el-radio label="gun">gun</el-radio>
                <el-radio label="multi">multi</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.service_name')"
                          v-if="nodeStoreData.dialogData.value.vmessInfo.network==='grpc'">
              <el-input v-model="nodeStoreData.dialogData.value.vmessInfo.service_name" />
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.host')"
                          v-if="nodeStoreData.dialogData.value.vmessInfo.network==='ws' || (nodeStoreData.dialogData.value.vmessInfo.network==='tcp' && nodeStoreData.dialogData.value.vmessInfo.type ==='http')">
              <el-input v-model="nodeStoreData.dialogData.value.vmessInfo.host" />
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.path')"
                          v-if="nodeStoreData.dialogData.value.vmessInfo.network==='ws' || (nodeStoreData.dialogData.value.vmessInfo.network==='tcp' && nodeStoreData.dialogData.value.vmessInfo.type ==='http')">
              <el-input v-model="nodeStoreData.dialogData.value.vmessInfo.path" />
            </el-form-item>

            <el-form-item :label="$t('message.adminNode.NodeInfo.security')">
              <el-radio-group v-model="nodeStoreData.dialogData.value.vmessInfo.security">
                <el-radio label="">none</el-radio>
                <el-radio label="tls">tls</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.sni')"
                          v-if="nodeStoreData.dialogData.value.vmessInfo.security==='tls'">
              <el-input v-model="nodeStoreData.dialogData.value.vmessInfo.sni" />
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.fp')"
                          v-if="nodeStoreData.dialogData.value.vmessInfo.security==='tls'">
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
            <el-form-item :label="$t('message.adminNode.NodeInfo.allowInsecure')">
              <el-switch
                size="small"
                v-model="nodeStoreData.dialogData.value.vmessInfo.allowInsecure"
                style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
              />
            </el-form-item>

            <el-divider content-position="left">{{ $t("message.adminNode.otherParameters") }}</el-divider>
            <el-form :model="nodeStoreData.dialogData.value.vmessInfo" label-position="top">
              <el-form-item :label="$t('message.adminNode.NodeInfo.enabled')">
                <el-switch
                  size="small"
                  inline-prompt
                  v-model="nodeStoreData.dialogData.value.vmessInfo.enabled"
                  style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
                  :active-text="$t('message.common.display')"
                  :inactive-text="$t('message.common.hide')"
                />
              </el-form-item>
              <div v-if="state.nodeType===constantStore.NODE_TYPE_NORMAL">
                <el-form-item :label="$t('message.adminNode.NodeInfo.transfer_address')">
                  <el-input v-model="nodeStoreData.dialogData.value.vmessInfo.transfer_address" />
                </el-form-item>
                <el-form-item :label="$t('message.adminNode.NodeInfo.transfer_port')">
                  <el-input-number v-model.number="nodeStoreData.dialogData.value.vmessInfo.transfer_port" />
                </el-form-item>
                <el-form-item :label="$t('message.adminNode.NodeInfo.node_speed_limit')">
                  <el-input-number  v-model.number="nodeStoreData.dialogData.value.vmessInfo.node_speed_limit" />
                </el-form-item>
                <el-form-item :label="$t('message.adminNode.NodeInfo.traffic_rate')">
                  <el-input-number v-model.number="nodeStoreData.dialogData.value.vmessInfo.traffic_rate" />
                </el-form-item>
                <el-form-item :label="$t('message.adminNode.NodeInfo.access')">
                  <el-tree ref="access_tree_ref" node-key="id"
                           :data="nodeStoreData.accessList.value.data"
                           :props="{label:'name'}"
                           :default-checked-keys="nodeStoreData.dialogData.value.checkedAccessIDs"
                           show-checkbox class="data-tree">
                    <template #default="{ node, data }">
                  <span class="custom-tree-node">
                    <span>{{ $t(data.name) }}</span>
                  </span>
                    </template>
                  </el-tree>
                </el-form-item>
              </div>

            </el-form>
          </el-form>
        </div>

        <div v-else-if="state.nodeProtocol === constantStore.NODE_PROTOCOL_SHADOWSOCKS">
          <el-form :model="nodeStoreData.dialogData.value.shadowsocksInfo" label-position="top">
            <el-form-item :label="$t('message.adminNode.NodeInfo.remarks')">
              <el-input v-model="nodeStoreData.dialogData.value.shadowsocksInfo.remarks" />
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.address')">
              <el-input v-model="nodeStoreData.dialogData.value.shadowsocksInfo.address" />
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.port')">
              <el-input-number v-model.number="nodeStoreData.dialogData.value.shadowsocksInfo.port" />
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.scy')">
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
            <el-form-item :label="$t('message.adminNode.NodeInfo.server_key')">
              <el-input v-model="nodeStoreData.dialogData.value.shadowsocksInfo.server_key">
                <template #append>
                  <el-button @click="setShadowsocksServerKey()">
                    <el-icon>
                      <Refresh />
                    </el-icon>
                  </el-button>
                </template>
              </el-input>
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.type')">
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
            <el-form-item :label="$t('message.adminNode.NodeInfo.host')"
                          v-if="nodeStoreData.dialogData.value.shadowsocksInfo.type === 'http'">
              <el-input v-model="nodeStoreData.dialogData.value.shadowsocksInfo.host" />
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.path')"
                          v-if="nodeStoreData.dialogData.value.shadowsocksInfo.type === 'http'">
              <el-input v-model="nodeStoreData.dialogData.value.shadowsocksInfo.path" />
            </el-form-item>

            <el-divider content-position="left">{{ $t("message.adminNode.otherParameters") }}</el-divider>
            <el-form-item :label="$t('message.adminNode.NodeInfo.enabled')">
              <el-switch
                size="small"
                inline-prompt
                v-model="nodeStoreData.dialogData.value.shadowsocksInfo.enabled"
                style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
                :active-text="$t('message.common.display')"
                :inactive-text="$t('message.common.hide')"
              />
            </el-form-item>
            <div v-if="state.nodeType===constantStore.NODE_TYPE_NORMAL">
              <el-form-item :label="$t('message.adminNode.NodeInfo.transfer_address')">
                <el-input v-model="nodeStoreData.dialogData.value.shadowsocksInfo.transfer_address" />
              </el-form-item>
              <el-form-item :label="$t('message.adminNode.NodeInfo.transfer_port')">
                <el-input-number v-model.number="nodeStoreData.dialogData.value.shadowsocksInfo.transfer_port" />
              </el-form-item>
              <el-form-item :label="$t('message.adminNode.NodeInfo.node_speed_limit')">
                <el-input-number
                          v-model.number="nodeStoreData.dialogData.value.shadowsocksInfo.node_speed_limit" />
              </el-form-item>
              <el-form-item :label="$t('message.adminNode.NodeInfo.traffic_rate')">
                <el-input-number v-model.number="nodeStoreData.dialogData.value.shadowsocksInfo.traffic_rate" />
              </el-form-item>
              <el-form-item :label="$t('message.adminNode.NodeInfo.access')">
                <el-tree ref="access_tree_ref" node-key="id"
                         :data="nodeStoreData.accessList.value.data"
                         :props="{label:'name'}"
                         :default-checked-keys="nodeStoreData.dialogData.value.checkedAccessIDs"
                         show-checkbox class="data-tree">
                  <template #default="{ node, data }">
                  <span class="custom-tree-node">
                    <span>{{ $t(data.name) }}</span>
                  </span>
                  </template>
                </el-tree>
              </el-form-item>
            </div>
          </el-form>
        </div>

        <div v-else-if="state.nodeProtocol === constantStore.NODE_PROTOCOL_HYSTERIA">
          <el-form :model="nodeStoreData.dialogData.value.hysteriaInfo" label-position="top">
            <el-form-item :label="$t('message.adminNode.NodeInfo.remarks')">
              <el-input v-model="nodeStoreData.dialogData.value.hysteriaInfo.remarks" />
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.address')">
              <el-input v-model="nodeStoreData.dialogData.value.hysteriaInfo.address" />
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.port')">
              <el-input-number v-model.number="nodeStoreData.dialogData.value.hysteriaInfo.port" />
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.hy_ports')">
              <el-input v-model="nodeStoreData.dialogData.value.hysteriaInfo.hy_ports" />
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.hy_up_mbps')">
              <el-input-number v-model.number="nodeStoreData.dialogData.value.hysteriaInfo.hy_up_mbps" />
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.hy_down_mbps')">
              <el-input-number v-model.number="nodeStoreData.dialogData.value.hysteriaInfo.hy_down_mbps" />
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.hy_obfs')">
              <el-select
                v-model="nodeStoreData.dialogData.value.hysteriaInfo.hy_obfs"
                filterable
                default-first-option
                :reserve-keyword="false"
                style="width: 100%"
              >
                <el-option
                  v-for="(v,k) in state.hy_obfs"
                  :key="k"
                  :label="v"
                  :value="v">
                </el-option>
              </el-select>
            </el-form-item>
            <el-form-item v-if="nodeStoreData.dialogData.value.hysteriaInfo.hy_obfs !==''" :label="$t('message.adminNode.NodeInfo.hy_obfs_password')">
              <el-input v-model="nodeStoreData.dialogData.value.hysteriaInfo.hy_obfs_password">
                <template #append>
                  <el-button @click="setHyObfsPassword()">
                    <el-icon>
                      <Refresh />
                    </el-icon>
                  </el-button>
                </template>
              </el-input>
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.sni')">
              <el-input v-model="nodeStoreData.dialogData.value.hysteriaInfo.sni" />
            </el-form-item>
            <el-form-item :label="$t('message.adminNode.NodeInfo.allowInsecure')">
              <el-switch
                size="small"
                v-model="nodeStoreData.dialogData.value.hysteriaInfo.allowInsecure"
                style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
              />
            </el-form-item>
            <el-divider content-position="left">{{ $t("message.adminNode.otherParameters") }}</el-divider>
            <el-form-item :label="$t('message.adminNode.NodeInfo.enabled')">
              <el-switch
                size="small"
                inline-prompt
                v-model="nodeStoreData.dialogData.value.hysteriaInfo.enabled"
                style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"
                :active-text="$t('message.common.display')"
                :inactive-text="$t('message.common.hide')"
              />
            </el-form-item>
            <div v-if="state.nodeType===constantStore.NODE_TYPE_NORMAL">
              <el-form-item :label="$t('message.adminNode.NodeInfo.transfer_address')">
                <el-input v-model="nodeStoreData.dialogData.value.hysteriaInfo.transfer_address" />
              </el-form-item>
              <el-form-item :label="$t('message.adminNode.NodeInfo.transfer_port')">
                <el-input-number v-model="nodeStoreData.dialogData.value.hysteriaInfo.transfer_port" />
              </el-form-item>
              <el-form-item :label="$t('message.adminNode.NodeInfo.node_speed_limit')">
                <el-input-number v-model.number="nodeStoreData.dialogData.value.hysteriaInfo.node_speed_limit" />
              </el-form-item>
              <el-form-item :label="$t('message.adminNode.NodeInfo.traffic_rate')">
                <el-input-number v-model.number="nodeStoreData.dialogData.value.hysteriaInfo.traffic_rate" />
              </el-form-item>
              <el-form-item :label="$t('message.adminNode.NodeInfo.access')">
                <el-tree ref="access_tree_ref" node-key="id"
                         :data="nodeStoreData.accessList.value.data"
                         :props="{label:'name'}"
                         :default-checked-keys="nodeStoreData.dialogData.value.checkedAccessIDs"
                         show-checkbox class="data-tree">
                  <template #default="{ node, data }">
                  <span class="custom-tree-node">
                    <span>{{ $t(data.name) }}</span>
                  </span>
                  </template>
                </el-tree>
              </el-form-item>
            </div>

          </el-form>
        </div>
      </div>
    </div>


    <template #footer>
      <el-button @click="closeDialog">{{ $t("message.common.button_cancel") }}</el-button>
      <el-button v-if="state.type === 'add'" @click="submitForAdd" type="primary">
        {{ $t("message.common.button_confirm") }}
      </el-button>
      <el-button v-else @click="submitForModify" type="primary">{{ $t("message.common.button_confirm") }}</el-button>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>

import { storeToRefs } from "pinia";
import { useAdminNodeStore } from "/@/stores/admin_logic/nodeStore";
import { reactive, watch, ref } from "vue";
import { useApiStore } from "/@/stores/apiStore";
import { request } from "/@/utils/request";
import { useConstantStore } from "/@/stores/constantStore";
import { DateStrToTime } from "../../../utils/formatTime";
import { useI18n } from "vue-i18n";
import {ElMessage} from "element-plus";
import { randomStringNoUpper } from "/@/utils/encrypt";

const apiStore = useApiStore();
const apiStoreData = storeToRefs(apiStore);
const nodeStore = useAdminNodeStore();
const nodeStoreData = storeToRefs(nodeStore);
const emit = defineEmits(["refresh"]);
const constantStore = useConstantStore();
const { t } = useI18n();
const state = reactive({
  type: "",
  title: "",
  isShowDialog: false,
  nodeType: constantStore.NODE_TYPE_NORMAL,
  nodeProtocol: constantStore.NODE_PROTOCOL_VLESS,
  nodeTypelist: [
    constantStore.NODE_TYPE_NORMAL,
    constantStore.NODE_TYPE_TRANSFER,
    constantStore.NODE_TYPE_SHARED
  ],
  nodeProtocolList: [
    constantStore.NODE_PROTOCOL_VLESS,
    constantStore.NODE_PROTOCOL_VMESS,
    constantStore.NODE_PROTOCOL_SHADOWSOCKS,
    constantStore.NODE_PROTOCOL_HYSTERIA
  ],
  realityDefaultArr: [
    { dest: "www.speedtest.org:443", sni: "www.speedtest.org" },
    { dest: "www.lovelive-anime.jp:443", sni: "www.lovelive-anime.jp" },
    { dest: "swdist.apple.com:443", sni: "swdist.apple.com" },
    { dest: "blog.api.www.cloudflare.com:443", sni: "blog.api.www.cloudflare.com" },
    { dest: "www.icloud.com:443", sni: "www.icloud.com" }
  ] as RealityItem[],
  scyArrForVmess: [
    "auto",
    "none",
    "aes-128-gcm",
    "aes-256-gcm",
    "chacha20-ietf-poly1305"],
  scyArrForSS: [
    "aes-128-gcm",
    "aes-256-gcm",
    "chacha20-ietf-poly1305",
    "2022-blake3-aes-128-gcm",
    "2022-blake3-aes-256-gcm",
    "2022-blake3-chacha20-poly1305"],
  scyArrForClashMeta: [
    "aes-128-gcm",
    "aes-192-gcm",
    "aes-256-gcm",
    "chacha20-ietf-poly1305",
    "2022-blake3-aes-128-gcm",
    "2022-blake3-aes-256-gcm",
    "2022-blake3-chacha20-poly1305"],
  flowArr: ["", "xtls-rprx-vision", "xtls-rprx-vision-udp443"],
  networkArr: ["ws", "tcp", "grpc"],
  typeArr1: ["none", "http"],
  typeArr2: ["none", "srtp", "utp", "wechat-video", "dtls", "wireguard"],
  fpArr: ["chrome", "firefox", "safari", "ios", "android", "edge", "360", "qq", "random", "randomized"],
  hy_obfs:["","salamander"],
});
const access_tree_ref = ref();
// 打开弹窗
const openDialog = (type: string, row?: NodeInfo) => {
  state.type = type;
  if (type === "add") {
    nodeStoreData.dialogData.value.vlessInfo.id = 0; //编辑和添加公用一个store，清空id,否则服务器无法插入
    nodeStoreData.dialogData.value.vmessInfo.id = 0; //编辑和添加公用一个store，清空id,否则服务器无法插入
    nodeStoreData.dialogData.value.shadowsocksInfo.id = 0; //编辑和添加公用一个store，清空id,否则服务器无法插入
    nodeStoreData.dialogData.value.hysteriaInfo.id = 0; //编辑和添加公用一个store，清空id,否则服务器无法插入
    nodeStoreData.dialogData.value.transferInfo.id = 0; //编辑和添加公用一个store，清空id,否则服务器无法插入
    state.title = t("message.adminNode.addNode");
    state.isShowDialog = true;
  } else {
    state.title = t("message.adminNode.modifyNode");

    switch (row?.node_type) {
      case constantStore.NODE_TYPE_SHARED:
        state.nodeType = constantStore.NODE_TYPE_SHARED;
        saveRow(row);
        break;
      case constantStore.NODE_TYPE_NORMAL:
        state.nodeType = constantStore.NODE_TYPE_NORMAL;
        saveRow(row);
        break;
      case constantStore.NODE_TYPE_TRANSFER:
        state.nodeType = constantStore.NODE_TYPE_TRANSFER;
        nodeStoreData.dialogData.value.transferInfo = row;
        state.isShowDialog = true;
        break;

    }
    nodeStore.accessHandler(row);
    state.isShowDialog = true;
  }
};
//根据节点类型保存数据
const saveRow = (row: NodeInfo) => {
  switch (row.protocol) {
    case constantStore.NODE_PROTOCOL_VLESS:
      state.nodeProtocol = constantStore.NODE_PROTOCOL_VLESS;
      nodeStoreData.dialogData.value.vlessInfo = row;
      break;
    case constantStore.NODE_PROTOCOL_VMESS:
      state.nodeProtocol = constantStore.NODE_PROTOCOL_VMESS;
      nodeStoreData.dialogData.value.vmessInfo = row;
      break;
    case constantStore.NODE_PROTOCOL_SHADOWSOCKS:
      state.nodeProtocol = constantStore.NODE_PROTOCOL_SHADOWSOCKS;
      nodeStoreData.dialogData.value.shadowsocksInfo = row;
      break;
    case constantStore.NODE_PROTOCOL_HYSTERIA:
      state.nodeProtocol = constantStore.NODE_PROTOCOL_HYSTERIA;
      nodeStoreData.dialogData.value.hysteriaInfo = row;
      break;
  }
};
//根据节点类型返回节点对象
const returnNodeInfo = (protocol: string) => {
  let n = {} as NodeInfo;
  switch (protocol) {
    case constantStore.NODE_PROTOCOL_VLESS:
      n = nodeStoreData.dialogData.value.vlessInfo;
      break;
    case constantStore.NODE_PROTOCOL_VMESS:
      n = nodeStoreData.dialogData.value.vmessInfo;
      break;
    case constantStore.NODE_PROTOCOL_SHADOWSOCKS:
      n = nodeStoreData.dialogData.value.shadowsocksInfo;
      break;
    case constantStore.NODE_PROTOCOL_HYSTERIA:
      n = nodeStoreData.dialogData.value.hysteriaInfo;
      break;
  }
  return n;
};

// 关闭弹窗
const closeDialog = () => {
  state.isShowDialog = false;
};

//确认提交
const submitForAdd = () => {
  //处理access
  if (state.nodeType === constantStore.NODE_TYPE_NORMAL) {
    nodeStoreData.dialogData.value.checkedAccessIDs = [...access_tree_ref.value.getCheckedKeys()];
  }
  //根据不同的参数，获取对应的请求数据
  let data: any;
  switch (state.nodeType) {
    case constantStore.NODE_TYPE_SHARED:
      //在pinia中处理
      break;
    case constantStore.NODE_TYPE_NORMAL:
      data = returnNodeInfo(state.nodeProtocol);
      data.node_type = state.nodeType;
      data.protocol = state.nodeProtocol;

      break;
    case constantStore.NODE_TYPE_TRANSFER:
      data = nodeStoreData.dialogData.value.transferInfo;
      data.node_type = state.nodeType;
      break;
  }
  // console.log("data:",data)
  //发送请求
  if (state.nodeType === constantStore.NODE_TYPE_SHARED) {
    nodeStore.newNodeShared().then((res) => {
      ElMessage.success(res.msg)
      emit("refresh");
    });
  } else {
    nodeStore.newNode(data).then((res) => {
      ElMessage.success(res.msg)
      emit("refresh");
    });
  }


  closeDialog();
};
const submitForModify = () => {
  //处理access
  if (state.nodeType === constantStore.NODE_TYPE_NORMAL) {
    nodeStoreData.dialogData.value.checkedAccessIDs = [...access_tree_ref.value.getCheckedKeys()];
  }
  //根据不同的参数，获取对应的请求数据
  let data: any;
  switch (state.nodeType) {
    case constantStore.NODE_TYPE_SHARED:
    case constantStore.NODE_TYPE_NORMAL:
      data = returnNodeInfo(state.nodeProtocol);
      data.node_type = state.nodeType;
      data.protocol = state.nodeProtocol;
      break;
    case constantStore.NODE_TYPE_TRANSFER:
      data = nodeStoreData.dialogData.value.transferInfo;
      data.node_type = state.nodeType;
      break;
  }
  nodeStore.updateNode(data).then((res) => {
    ElMessage.success(res.msg)
    emit("refresh");
  });
  closeDialog();
};

//
const setReality = (nodeType: string) => {
  request(apiStoreData.adminApi.value.createx25519).then((res) => {
    switch (nodeType) {
      case constantStore.NODE_PROTOCOL_VLESS:
        nodeStoreData.dialogData.value.vlessInfo.pbk = res.data.public_key;
        nodeStoreData.dialogData.value.vlessInfo.private_key = res.data.private_key;
        break;
    }
  });
};
const setHyObfsPassword=()=>{
  nodeStoreData.dialogData.value.hysteriaInfo.hy_obfs_password = randomStringNoUpper(32)
}
const setShadowsocksServerKey=()=>{
  nodeStoreData.dialogData.value.shadowsocksInfo.server_key = randomStringNoUpper(32)
}
const parseUrl = () => {
  nodeStore.parseUrl();
};
const deleteOneSharedNode = (row: any) => {
  nodeStoreData.nodeSharedData.value.nodeList.data = nodeStoreData.nodeSharedData.value.nodeList.data.filter(item => item !== row);
};
//监听
watch(
  () => nodeStoreData.dialogData.value.vlessInfo.dest,
  () => {
    let temp = state.realityDefaultArr.filter(r => r.dest === nodeStoreData.dialogData.value.vlessInfo.dest);
    if (temp.length > 0) {
      nodeStoreData.dialogData.value.vlessInfo.sni = temp[0].sni;
    }
    if (nodeStoreData.dialogData.value.vlessInfo.private_key === "" || nodeStoreData.dialogData.value.vlessInfo.private_key === "") {
      setReality(constantStore.NODE_PROTOCOL_VLESS);
    }
  },
  {
    // deep: true,
  }
);
watch(
  () => nodeStoreData.dialogData.value.hysteriaInfo.hy_obfs,
  () => {
    if (nodeStoreData.dialogData.value.hysteriaInfo.hy_obfs_password === ''){
      setHyObfsPassword()
    }
  },
  {
    // deep: true,
  }
);

// 暴露变量
defineExpose({
  openDialog   // 打开弹窗
});
</script>


<style scoped>
.dialog-footer button:first-child {
  margin-right: 10px;
}

.data-tree {
  width: 100%;
  border: 1px solid var(--el-border-color);
  border-radius: var(--el-input-border-radius, var(--el-border-radius-base));
  padding: 5px;
}

.custom-tree-node {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 14px;
  padding-right: 8px;
}
</style>
