<template>
  <div style="padding: 15px;">
    <el-card>
      <el-tabs stretch style="height: 100%" @tab-change="tap" v-model="state.currentTapName">
        <el-tab-pane :label="$t('message.adminServer.tabWebsite')" name="1">
          <el-row style="margin-bottom: 20px">
            <el-col :span="12">
              <div>
                <el-icon style="margin-right: 4px" :size="12">
                  <InfoFilled />
                </el-icon>
                {{$t('message.adminServer.Server.current_version')}}
              </div>
              <div style="font-size: 20px">{{ serverConfig.version.value.currentVersion.version }}</div>
            </el-col>
            <el-col :span="12">
              <div>
                <el-icon style="margin-right: 4px" :size="12">
                  <UploadFilled />
                </el-icon>
                {{$t('message.adminServer.Server.latest_version')}}
              </div>
              <div style="color: red;font-size: 20px">{{ serverConfig.version.value.latestVersion.version }}</div>
            </el-col>
          </el-row>

          <el-row>
            <el-button type="primary" :disabled="state.inShowUpdateButton" @click="openUpdateDialog">{{$t('message.adminServer.Server.starting_upgrade')}}</el-button>
          </el-row>
          <el-divider></el-divider>

          <el-form :model="serverConfig.serverConfig.value.website" label-position="top">
            <el-form-item :label="$t('message.adminServer.Server.enable_register')" >

              <el-switch v-model="serverConfig.serverConfig.value.website.enable_register" inline-prompt
                         :active-text="$t('message.common.enable')"
                         :inactive-text="$t('message.common.disable')"
                         style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>


            </el-form-item>
            <el-form-item :label="$t('message.adminServer.Server.acceptable_email_suffixes')" >
              <el-input v-model="serverConfig.serverConfig.value.website.acceptable_email_suffixes" type="textarea"
                        autosize />
            </el-form-item>
            <el-form-item :label="$t('message.adminServer.Server.enable_email_code')" >
              <el-switch v-model="serverConfig.serverConfig.value.website.enable_email_code" inline-prompt
                         :active-text="$t('message.common.enable')"
                         :inactive-text="$t('message.common.disable')"
                         style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
            </el-form-item>
            <el-form-item :label="$t('message.adminServer.Server.enable_base64_captcha')" >
              <el-switch v-model="serverConfig.serverConfig.value.website.enable_base64_captcha" inline-prompt
                         :active-text="$t('message.common.enable')"
                         :inactive-text="$t('message.common.disable')"
                         style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
            </el-form-item>
            <el-form-item :label="$t('message.adminServer.Server.enable_swagger_api')" >
              <el-switch v-model="serverConfig.serverConfig.value.website.enable_swagger_api" inline-prompt
                         :active-text="$t('message.common.enable')"
                         :inactive-text="$t('message.common.disable')"
                         style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
            </el-form-item>
            <el-form-item :label="$t('message.adminServer.Server.enable_assets_api')" >
              <el-switch v-model="serverConfig.serverConfig.value.website.enable_assets_api" inline-prompt
                         :active-text="$t('message.common.enable')"
                         :inactive-text="$t('message.common.disable')"
                         style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
            </el-form-item>
            <el-form-item :label="$t('message.adminServer.Server.frontend_url')" >
              <el-input v-model="serverConfig.serverConfig.value.website.frontend_url" placeholder="http://xxx.com" />
            </el-form-item>
            <el-divider></el-divider>
            <el-form-item>
              <el-button @click="onSubmit()" type="primary">{{ $t("message.common.button_confirm") }}</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        <el-tab-pane :label="$t('message.adminServer.tabSubscribe')" name="2">
          <el-form-item :label="$t('message.adminServer.Server.tek')" >
            <el-input v-model="serverConfig.serverConfig.value.subscribe.tek" />
          </el-form-item>
          <el-form-item :label="$t('message.adminServer.Server.sub_name')" >
            <el-input v-model="serverConfig.serverConfig.value.subscribe.sub_name" />
          </el-form-item>
          <el-form-item :label="$t('message.adminServer.Server.backend_url')" >
            <el-input v-model="serverConfig.serverConfig.value.subscribe.backend_url" placeholder="http://xxx.com"
                      type="textarea" autosize />
            <span style="color: #9b9da1">*{{ $t("message.adminServer.Server.sub_prefix_msg") }}</span>
          </el-form-item>
          <el-form-item :label="$t('message.adminServer.Server.subscribe_domain_bind_request')" >
            <el-switch v-model="serverConfig.serverConfig.value.subscribe.subscribe_domain_bind_request" inline-prompt
                       :active-text="$t('message.common.enable')"
                       :inactive-text="$t('message.common.disable')"
                       style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
          </el-form-item>
          <el-form-item :label="$t('message.adminServer.Server.clash_rule')" >
            <el-input v-model="serverConfig.serverConfig.value.subscribe.clash_rule"
                      type="textarea" :rows="10" />
          </el-form-item>
          <el-form-item :label="$t('message.adminServer.Server.surge_rule')" >
            <el-input v-model="serverConfig.serverConfig.value.subscribe.surge_rule"
                      type="textarea" :rows="10" />
          </el-form-item>
          <el-divider></el-divider>
          <el-form-item>
            <el-button @click="onSubmit()" type="primary">{{ $t("message.common.button_confirm") }}</el-button>
          </el-form-item>
        </el-tab-pane>

        <el-tab-pane :label="$t('message.adminServer.tabPayment')" name="3">
          <div>
            <el-button size="default" type="primary" class="ml10" @click="openPayDialog('add')">
              <el-icon>
                <ele-FolderAdd />
              </el-icon>
              {{ $t("message.adminServer.addPay") }}
            </el-button>
          </div>
          <div>
            <el-table :data="adminShopStoreData.payList.value" stripe style="width: 100%;flex: 1;">
              <el-table-column type="index" :label="$t('message.adminShop.PayInfo.index')" fixed show-overflow-tooltip
                               width="60px" />
              <el-table-column prop="id" :label="$t('message.adminShop.PayInfo.id')" show-overflow-tooltip
                               width="60px" />
              <el-table-column prop="name" :label="$t('message.adminShop.PayInfo.name')" show-overflow-tooltip
                               width="200px" />
              <el-table-column prop="pay_type" :label="$t('message.adminShop.PayInfo.pay_type')" show-overflow-tooltip
                               width="120px" />
              <el-table-column prop="pay_logo_url" :label="$t('message.adminShop.PayInfo.pay_logo_url')"
                               show-overflow-tooltip width="120px">
                <template #default="{row}">
                  <el-image :src="row.pay_logo_url" style="width: 40px;height: 40px"></el-image>
                </template>
              </el-table-column>
              <el-table-column prop="status" :label="$t('message.adminShop.PayInfo.status')" show-overflow-tooltip
                               width="80px">
                <template #default="{row}">
                  <el-button v-if="row.status" type="warning">{{ $t("message.common.enable") }}</el-button>
                  <el-button v-else type="info">{{ $t("message.common.disable") }}</el-button>
                </template>
              </el-table-column>

              <el-table-column :label="$t('message.common.operate')">
                <template #default="scope">
                  <el-button text @click="openPayDialog('edit',scope.row)" type="primary">
                    {{ $t("message.common.modify") }}
                  </el-button>
                  <el-button text @click="deletePay(scope.row)" type="primary">{{ $t("message.common.delete") }}
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-tab-pane>

        <el-tab-pane :label="$t('message.adminServer.tabEmail')" name="4">
          <el-form :model="serverConfig.serverConfig.value.email" label-position="top">
            <el-form-item :label="$t('message.adminServer.Server.email_host')" >
              <el-input v-model="serverConfig.serverConfig.value.email.email_host" placeholder="mail.example.com" />
            </el-form-item>
            <el-form-item :label="$t('message.adminServer.Server.email_port')" >
              <el-input v-model.number="serverConfig.serverConfig.value.email.email_port" type="number" />
            </el-form-item>
            <el-form-item :label="$t('message.adminServer.Server.email_from')" >
              <el-input v-model="serverConfig.serverConfig.value.email.email_from" placeholder="admin@qq.com" />
            </el-form-item>
            <el-form-item :label="$t('message.adminServer.Server.email_from_alias')" >
              <el-input v-model="serverConfig.serverConfig.value.email.email_from_alias" placeholder="admin@qq.com" />
            </el-form-item>
            <el-form-item :label="$t('message.adminServer.Server.email_nickname')" >
              <el-input v-model="serverConfig.serverConfig.value.email.email_nickname" placeholder="Admin" />
            </el-form-item>
            <el-form-item :label="$t('message.adminServer.Server.email_secret')" >
              <el-input v-model="serverConfig.serverConfig.value.email.email_secret" />
            </el-form-item>
            <el-form-item :label="$t('message.adminServer.Server.email_subject')" >
              <el-input v-model="serverConfig.serverConfig.value.email.email_subject" />
            </el-form-item>
            <el-form-item :label="$t('message.adminServer.Server.email_content')" >
              <el-input v-model="serverConfig.serverConfig.value.email.email_content" type="textarea" autosize />
              <el-text style="color: #9b9da1">{{ $t("message.adminServer.emailCodeTip") }}
              </el-text>
            </el-form-item>
            <el-divider></el-divider>
            <el-form-item>
              <el-button @click="onSubmit()" type="primary">{{ $t("message.common.button_confirm") }}</el-button>
              <el-button @click="onTestEmail">{{ $t("message.adminServer.emailTesting") }}</el-button>
            </el-form-item>
          </el-form>

        </el-tab-pane>

        <el-tab-pane :label="$t('message.adminServer.tabSecurity')" name="5">
          <el-form :model="serverConfig.serverConfig.value.security" label-position="top">
            <el-form-item :label="$t('message.adminServer.Server.ip_role_param')" >
              <el-col :span="2">
                <el-input-number v-model="serverConfig.serverConfig.value.security.rate_limit_params.ip_role_param"
                                 :precision="0" :step="10" :min="0" :max="10000000" />
              </el-col>
              <el-col :span="2" style="text-align: center">
                <span>-</span>
              </el-col>
              <el-col :span="18">
                <span class="text-gray-500">{{ $t("message.adminServer.RequestMinute") }}</span>
              </el-col>
            </el-form-item>

            <el-form-item :label="$t('message.adminServer.Server.visit_param')" >
              <el-col :span="2">
                <el-input-number v-model="serverConfig.serverConfig.value.security.rate_limit_params.visit_param"
                                 :precision="0" :step="10" :min="0" :max="10000000" />
              </el-col>
              <el-col :span="2" style="text-align: center">
                <span>-</span>
              </el-col>
              <el-col :span="18">
                <span class="text-gray-500">{{ $t("message.adminServer.RequestMinute") }}</span>
              </el-col>
            </el-form-item>
            <el-divider></el-divider>
            <el-form-item :label="$t('message.adminServer.Server.signing_key')" >
              <el-input v-model="serverConfig.serverConfig.value.security.jwt.signing_key" />
            </el-form-item>
            <el-form-item :label="$t('message.adminServer.Server.issuer')" >
              <el-input v-model="serverConfig.serverConfig.value.security.jwt.issuer" />
            </el-form-item>
            <el-form-item :label="$t('message.adminServer.Server.expires_time')" >
              <el-input v-model="serverConfig.serverConfig.value.security.jwt.expires_time" />
            </el-form-item>
            <el-form-item>
              <el-button @click="onSubmit()" type="primary">{{ $t("message.common.button_confirm") }}</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane :label="$t('message.adminServer.tabNotice')" name="6">
          <div style="margin-bottom: 50px">
            <el-divider content-position="left"><span style="color: var(--el-color-primary) !important">{{ $t("message.adminServer.Server.admin_id") }}</span></el-divider>
            <div>
              <el-form :model="serverConfig.serverConfig.value.notice" label-width="150px"
                       label-position="left">
                <el-form-item :label="$t('message.adminServer.Server.admin_id')" >
                  <el-input v-model="serverConfig.serverConfig.value.notice.admin_id" type="textarea" autosize />
                  <span style="color: #9b9da1">*{{$t('message.adminServer.Server.admin_id_msg')}}</span>
                </el-form-item>
              </el-form>
            </div>
          </div>

          <div style="margin-bottom: 50px">
            <el-divider content-position="left"><span style="color: var(--el-color-primary) !important">{{ $t("message.adminServer.Server.push_method") }}</span></el-divider>
            <div>
              <el-form :model="serverConfig.serverConfig.value.notice" label-width="100px"
                       label-position="left">
                <el-form-item :label="$t('message.adminServer.Server.enable_tg_bot')" >
                  <el-switch v-model="serverConfig.serverConfig.value.notice.enable_tg_bot" inline-prompt
                             :active-text="$t('message.common.enable')"
                             :inactive-text="$t('message.common.disable')"
                             style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
                </el-form-item>

                <el-form-item :label="$t('message.adminServer.Server.bot_token')" >
                  <el-input v-model="serverConfig.serverConfig.value.notice.bot_token"
                            placeholder="1234567890:AAAAABBBBBCCCCCDDDDFFFFGGGHHHJJKKLL" />
                </el-form-item>
                <el-form-item :label="$t('message.adminServer.Server.tg_socks5')" >
                  <el-input v-model="serverConfig.serverConfig.value.notice.tg_socks5"
                            placeholder="127.0.0.1:1080" />
                </el-form-item>
              </el-form>
            </div>
          </div>

          <div style="margin-bottom: 50px">
            <el-divider content-position="left"><span style="color: var(--el-color-primary) !important">{{ $t("message.adminServer.Server.trigger_condition") }}</span></el-divider>
            <div>
              <el-form :model="serverConfig.serverConfig.value.notice" label-width="150px"
                       label-position="left">
                <el-form-item :label="$t('message.adminServer.Server.when_node_offline')" >
                  <el-switch v-model="serverConfig.serverConfig.value.notice.when_node_offline" inline-prompt
                             :active-text="$t('message.common.enable')"
                             :inactive-text="$t('message.common.disable')"
                             style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
                </el-form-item>
                <el-form-item :label="$t('message.adminServer.Server.when_user_registered')" >
                  <el-switch v-model="serverConfig.serverConfig.value.notice.when_user_registered" inline-prompt
                             :active-text="$t('message.common.enable')"
                             :inactive-text="$t('message.common.disable')"
                             style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
                </el-form-item>
                <el-form-item :label="$t('message.adminServer.Server.when_user_purchased')" >
                  <el-switch v-model="serverConfig.serverConfig.value.notice.when_user_purchased" inline-prompt
                             :active-text="$t('message.common.enable')"
                             :inactive-text="$t('message.common.disable')"
                             style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
                </el-form-item>
                <el-form-item :label="$t('message.adminServer.Server.when_new_ticket')" >
                  <el-switch v-model="serverConfig.serverConfig.value.notice.when_new_ticket" inline-prompt
                             :active-text="$t('message.common.enable')"
                             :inactive-text="$t('message.common.disable')"
                             style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
                </el-form-item>
                <el-form-item style="margin-top: 20px">
                  <el-button @click="onSubmit()" type="primary">{{ $t("message.common.button_confirm") }}
                  </el-button>
                </el-form-item>
              </el-form>
            </div>
          </div>

        </el-tab-pane>

        <el-tab-pane :label="$t('message.adminServer.tabMigration')" name="7">
          <div style="margin-bottom: 20px">
            <el-alert :title="$t('message.adminServer.migrationTip')" type="warning" effect="dark" />
          </div>

          <el-form v-model="state.migrationParams" label-position="top">
            <el-form-item :label="$t('message.adminServer.Migration.panel_type')" >
              <el-select v-model="state.migrationParams.panel_type" placeholder="Select">
                <el-option
                  v-for="item in state.panels"
                  :key="item"
                  :label="item"
                  :value="item"
                />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('message.adminServer.Migration.db_address')" >
              <el-input v-model="state.migrationParams.db_address"></el-input>
            </el-form-item>
            <el-form-item :label="$t('message.adminServer.Migration.db_port')" >
              <el-input-number v-model="state.migrationParams.db_port"></el-input-number>
            </el-form-item>
            <el-form-item :label="$t('message.adminServer.Migration.db_name')" >
              <el-input v-model="state.migrationParams.db_name"></el-input>
            </el-form-item>
            <el-form-item :label="$t('message.adminServer.Migration.db_username')" >
              <el-input v-model="state.migrationParams.db_username"></el-input>
            </el-form-item>
            <el-form-item :label="$t('message.adminServer.Migration.db_password')" >
              <el-input v-model="state.migrationParams.db_password"></el-input>
            </el-form-item>
            <el-form-item>
              <el-button color="blue" @click="onSubmitMigration">{{ $t("message.common.button_confirm") }}</el-button>
            </el-form-item>
          </el-form>

        </el-tab-pane>

        <el-tab-pane :label="$t('message.adminServer.tabFinance')" name="8">
          <div>
            <el-form :model="serverConfig.serverConfig.value.finance" label-width="150px"
                     label-position="left">
              <el-form-item :label="$t('message.adminServer.Server.enable_invitation_commission')" >
                <el-switch v-model="serverConfig.serverConfig.value.finance.enable_invitation_commission" inline-prompt
                           :active-text="$t('message.common.enable')"
                           :inactive-text="$t('message.common.disable')"
                           style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
              </el-form-item>
              <el-form-item :label="$t('message.adminServer.Server.commission_rate')" >
                <el-input-number v-model="serverConfig.serverConfig.value.finance.commission_rate" :min="0" :step="0.01" :max="1"/>
              </el-form-item>
              <el-form-item :label="$t('message.adminServer.Server.withdraw_threshold')" >
                <el-input-number v-model="serverConfig.serverConfig.value.finance.withdraw_threshold" :min="0" :step="1"/>
              </el-form-item>

              <el-divider></el-divider>

              <el-form-item :label="$t('message.adminServer.Server.enable_lottery')" >
                <el-switch v-model="serverConfig.serverConfig.value.finance.enable_lottery" inline-prompt
                           :active-text="$t('message.common.enable')"
                           :inactive-text="$t('message.common.disable')"
                           style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949"></el-switch>
              </el-form-item>
              <el-form-item :label="$t('message.adminServer.Server.jackpot')" >
              </el-form-item>
            </el-form>
            <el-form :model="serverConfig.serverConfig.value.finance" label-width="100px"
                     label-position="right">
              <div v-for="(v,k) in serverConfig.serverConfig.value.finance.jackpot" :key="k">
                <el-row>
                  <el-col :span="12">
                    <el-form-item :label="$t('message.adminServer.Server.prize')">
                      <el-input-number v-model.number="v.balance" :precision="2" :min="0" :step="0.01"></el-input-number>
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item :label="$t('message.adminServer.Server.weight')">
                      <el-input-number v-model.number="v.weight" :precision="0" :min="0" :step="1"></el-input-number>
                    </el-form-item>
                  </el-col>
                </el-row>
              </div>
            </el-form>
            <el-form-item style="margin-top: 20px">
              <el-button @click="onSubmit()" type="primary">{{ $t("message.common.button_confirm") }}
              </el-button>
            </el-form-item>
          </div>
        </el-tab-pane>

      </el-tabs>
      <PayDialog ref="PayDialogRef" @refresh="adminShopStore.getPayList()"></PayDialog>
    </el-card>
    <el-dialog v-model="state.isShowTestEmailDialog" :title="$t('message.adminServer.emailTesting')" width="400px"
               destroy-on-close center>
      <el-form :model="state.emailParams" label-position="top">
        <el-form-item>
          <el-input v-model="state.emailParams.target_email" placeholder="xxx@xxx.com" />
        </el-form-item>
      </el-form>
      <template #footer>
      <span class="dialog-footer">
        <el-button @click="onGetEmailCode">{{ $t("message.common.button_confirm") }}</el-button>
      </span>
      </template>
    </el-dialog>
    <el-dialog v-model="state.isShowUpdateDialog" :title="$t('message.adminServer.Server.starting_upgrade')" width="80%"
               destroy-on-close center @close="closeSEE">
      <div
        v-if="serverConfig.version.value.currentVersion.version === serverConfig.version.value.latestVersion.version">
        <span>{{$t('message.adminServer.Server.is_latest')}}</span>
      </div>
      <div v-else>
        <el-row style="margin-bottom: 20px">
          <el-col :span="12" style="text-align: center;">
            <div>
              <el-icon style="margin-right: 4px" :size="12">
                <InfoFilled />
              </el-icon>
              <span>{{$t('message.adminServer.Server.current_version')}}</span>
            </div>
            <div style="font-size: 20px">{{ serverConfig.version.value.currentVersion.version }}</div>
          </el-col>
          <el-col :span="12" style="text-align: center;">
            <div>
              <el-icon style="margin-right: 4px" :size="12">
                <UploadFilled />
              </el-icon>
              <span>{{$t('message.adminServer.Server.latest_version')}}</span>
            </div>
            <div style="color: red;font-size: 20px">{{ serverConfig.version.value.latestVersion.version }}</div>
          </el-col>
        </el-row>
        <el-result icon="warning" :title="$t('message.adminServer.Server.upgrade_warning')"></el-result>
        <div v-if="state.isShowLogData">
          <div>
            <codemirror ref="codemirrorRef" v-model="state.logContent"></codemirror>
          </div>
        </div>

      </div>
      <template #footer>
      <span class="dialog-footer">
        <el-button @click="state.isShowUpdateDialog = false">{{ $t("message.common.button_cancel") }}</el-button>
        <el-button type="primary" @click="SSE" :disabled="state.isShowLogData">{{$t('message.adminServer.Server.starting_upgrade')}}</el-button>
      </span>
      </template>

    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import { defineAsyncComponent, onMounted, reactive, ref } from "vue";
import { useAdminServerStore } from "/@/stores/admin_logic/serverStore";
import { storeToRefs } from "pinia";
import { ElMessage, ElMessageBox } from "element-plus";
import { getApiPrefixAddress, request } from "/@/utils/request";
import { useApiStore } from "/@/stores/apiStore";
import { useUserStore } from "/@/stores/user_logic/userStore";
import { useAdminShopStore } from "/@/stores/admin_logic/shopStore";
import { usePublicStore } from "/@/stores/publicStore";
import { useI18n } from "vue-i18n";
import { useConstantStore } from "/@/stores/constantStore";
import { EventSourcePolyfill } from "event-source-polyfill";
import { Local } from "/@/utils/storage";
import { Codemirror } from "vue-codemirror";
import {apiUrl} from "/@/utils/request"

const apiStore = useApiStore();
const apiStoreData = storeToRefs(apiStore);
const PayDialog = defineAsyncComponent(() => import("/@/views/admin/system/dialog_pay.vue"));
const PayDialogRef = ref();
const serverStore = useAdminServerStore();
const serverConfig = storeToRefs(serverStore);
const adminShopStore = useAdminShopStore();
const adminShopStoreData = storeToRefs(adminShopStore);
const publicStore = usePublicStore();
const userStore = useUserStore();
const { t } = useI18n();
const constantStore = useConstantStore();
const codemirrorRef = ref()

const state = reactive({
  currentTapName: "1",
  isShowTestEmailDialog: false,
  emailParams: {
    email_type: constantStore.EMAIL_TYPE_TEST,
    target_email: ""
  },
  loading: false,
  panels: ["v2board", "sspanel", "AirGo"],
  migrationParams: {
    "panel_type": "",
    "db_address": "",
    "db_port": 3306,
    "db_username": "",
    "db_password": "",
    "db_name": ""
  },
  migrationResult: "",
  inShowUpdateButton: false,
  isShowUpdateDialog: false,
  isShowLogData: false,
  logContent:'',
  logOptions:{
    tabSize: 2, // 缩进格式
    // theme: 'rubyblue',
    lineNumbers: true, // 是否显示行号
    // mode: 'javascript',
    line: true,
    styleActiveLine: true, // 高亮选中行
    //是否为只读,如果为"nocursor" 不仅仅为只读 连光标都无法在区域聚焦
    readOnly: true,
    // hintOptions: {
    //   completeSingle: true // 当匹配只有一项的时候是否自动补全
    // },
    viewportMargin: 30
  },
});
const tap = (tapName: string) => {
  switch (tapName) {
    case "":
      break;
    default:
      break;
  }

};
//测试邮件
const onTestEmail = () => {
  state.isShowTestEmailDialog = true;
};
//获取邮件验证码
const onGetEmailCode = () => {
  if (state.emailParams.target_email === "") {
    return;
  }
  request(apiStoreData.publicApi.value.getEmailCode, state.emailParams).then((res) => {
    ElMessage.success(res.msg);
  });
};


//打开支付编辑
const openPayDialog = (type: string, row?: PayInfo) => {
  PayDialogRef.value.openDialog(type, row);
};

//保存提交
const onSubmit = () => {
  serverStore.updateServerConfig(serverConfig.serverConfig.value).then((res) => {
    ElMessage.success(res.msg);
    serverStore.getServerConfig();
    publicStore.getPublicSetting();
  });
};
const onSubmitMigration = () => {
  ElMessageBox.confirm(t("message.adminServer.migrationTip2"), t("message.common.tip"), {
    confirmButtonText: t("message.common.button_confirm"),
    cancelButtonText: t("message.common.button_cancel"),
    type: "warning"
  })
    .then(() => {
      state.loading = true;
      request(apiStore.adminApi.migrationData, state.migrationParams).then((res) => {
        state.migrationResult = res.data;
        state.loading = false;
      }).catch(() => {
        state.loading = false;
      });
    })
    .catch(() => {
    });
};
//删除支付
const deletePay = (data: PayInfo) => {
  ElMessageBox.confirm(t("message.common.message_confirm_delete"), t("message.common.tip"), {
    confirmButtonText: t("message.common.button_confirm"),
    cancelButtonText: t("message.common.button_cancel"),
    type: "warning"
  })
    .then(() => {
      request(apiStoreData.adminApi.value.deletePay, data).then((res) => {
        adminShopStore.getPayList(); //获取支付列表
      });
    })
    .catch(() => {
    });

};
//获取版本
const getVersion = () => {
  serverStore.getCurrentVersion();
  serverStore.getLatestVersion();
};
//打开升级弹窗
const openUpdateDialog = () => {
  state.isShowUpdateDialog = true;
};

const SSE = () => {
  state.isShowLogData = true;
  const url = apiUrl + apiStore.adminApi.updateLatestVersion.path;
  let token = Local.get("token");
  if (window.EventSource) {
    // let sseSource = new EventSource(url, { withCredentials: true });
    const sseSource = new EventSourcePolyfill(url, {
      headers: {
        "Authorization": token
      },
      "heartbeatTimeout":200000,
    });
    sseSource.onopen = function(e: any) {
      console.log("建立连接", e);
    };
    sseSource.onmessage = function(e: any) {
      console.log("onmessage 收到数据", e.data, "消息类型", e.type);
      showMsg(e.data)
    };
    sseSource.onerror = function(e: any) {
      console.log("关闭连接", e);
      sseSource.close();
    };
    //自定义2个类型的监听，其中message error是为了和默认的 onerror 区别
    sseSource.addEventListener("success", function(e: any) {
      console.log("onmessage 收到数据", e.data, "消息类型", e.type);
      showMsg(e.data)
    });
    sseSource.addEventListener("message error", function(e: any) {
      console.log("onmessage 收到数据", e.data, "消息类型", e.type);
      showMsg(e.data)
    });
  } else {
    console.log("浏览器不支持SSE");
  }
};
const showMsg=(msg:string)=>{
  state.logContent+=msg+'\r'
  //滚动到最后一行
  // let sc = codemirrorRef.value.codemirror.codemirror.getScrollInfo();
  // codemirrorRef.value.codemirror.codemirror.scrollTo(sc.left,( sc.height + sc.top));
}
const closeSEE=()=>{
  state.isShowLogData = false
  state.logContent = ''
}

onMounted(() => {
  serverStore.getServerConfig(); //获取设置参数
  adminShopStore.getGoodsList(); //获取全部商品，用来设置新注册分配套餐
  adminShopStore.getPayList();   //获取支付列表
  getVersion();                  //获取版本
});

</script>

<style lang="scss">
.label .el-form-item__label {
  font-weight: bolder;
  //font-size: 15px;
}
.form-inline .el-input {
  --el-input-width: 100px;
}

</style>