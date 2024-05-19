<template>
  <div class="personal layout-pd">
      <el-card style="border-radius:10px;padding: 10px">
        <el-row>
          <el-col :span="12">
            <el-row>
              <el-col :span="6">
                <div style="width: 100%">
                  <el-image :src="paperMoneyIcon" fit="cover" style="height: 50px;width: 50px"></el-image>
                </div>
              </el-col>
              <el-col :span="2"></el-col>
              <el-col :span="16">
                <div class="greyText">{{ $t("message.finance.balance") }}</div>
                <div><span style="font-size: 20px;font-weight: bold">{{ userInfos.balance }} </span></div>
              </el-col>
            </el-row>
          </el-col>
          <el-col :span="12">
            <el-row>
              <el-col :span="6">
                <div style="width: 100%">
                  <el-image :src="expensesIcon" fit="cover" style="height: 50px;width: 50px"></el-image>
                </div>
              </el-col>
              <el-col :span="2"></el-col>
              <el-col :span="16">
                <div class="greyText">{{$t('message.finance.total_consumption_amount')}}</div>
                <div><span style="font-size: 20px;font-weight: bold">{{ financeStoreData.commissionSummary.value.total_consumption_amount }}</span></div>
              </el-col>
            </el-row>
          </el-col>
        </el-row>
      </el-card>

      <el-card style="margin-top: 10px;border-radius:10px;padding: 10px">
        <el-row>
          <el-col :span="12">
            <el-row>
              <el-col :span="8">
                <div style="width: 100%">
                  <el-image :src="giftIcon" fit="cover" style="height: 50px;width: 50px"></el-image>
                </div>
              </el-col>
              <el-col :span="12">
                <div class="greyText">{{$t('message.finance.total_commission_amount')}}</div>
                <div><span style="font-size: 20px;font-weight: bold">{{ financeStoreData.commissionSummary.value.total_commission_amount }}</span></div>
              </el-col>
            </el-row>
          </el-col>
          <el-col :span="12">
            <el-row>
              <el-col :span="8">
                <div style="width: 100%">
                  <el-image :src="incomeOneIcon" fit="cover" style="height: 50px;width: 50px"></el-image>
                </div>
              </el-col>
              <el-col :span="12">
                <div class="greyText">{{$t('message.finance.pending_withdrawal_amount')}}</div>
                <div><span style="font-size: 20px;font-weight: bold">{{ financeStoreData.commissionSummary.value.pending_withdrawal_amount }}</span></div>
              </el-col>
            </el-row>
          </el-col>
        </el-row>
        <el-row class="mt15">
          <el-col :span="12">
            <el-row>
              <el-col :span="8">
                <div style="width: 100%">
                  <el-image :src="peoplePlusIcon" fit="cover" style="height: 50px;width: 50px"></el-image>
                </div>
              </el-col>
              <el-col :span="12">
                <div class="greyText">{{$t('message.finance.total_invitation')}}</div>
                <div><span style="font-size: 20px;font-weight: bold">{{ financeStoreData.commissionSummary.value.total_invitation }}</span></div>
              </el-col>
            </el-row>
          </el-col>
          <el-col :span="12">
            <el-row>
              <el-col :span="8">
                <div style="width: 100%">
                  <el-image :src="redEnvelopesIcon" fit="cover" style="height: 50px;width: 50px"></el-image>
                </div>
              </el-col>
              <el-col :span="12">
                <div class="greyText">{{$t('message.finance.commission_rate')}}</div>
                <div><span style="font-size: 20px;font-weight: bold">{{pubicStoreData.publicSetting.value.commission_rate*100}}%</span></div>
              </el-col>
            </el-row>
          </el-col>
        </el-row>
        <div class="mt30">
          <div style="margin-bottom: 15px">
            <el-button style="width: 150px;margin-right: 20px;" type="primary"
                       :disabled="Number(financeStoreData.commissionSummary.value.pending_withdrawal_amount) < pubicStoreData.publicSetting.value.withdraw_threshold "
                       @click="withdrawToBalance()">{{$t('message.finance.withdraw')}}</el-button>
            <span style="color: #9b9da1;font-size: 12px">{{$t('message.finance.withdraw_threshold')}}: >= {{pubicStoreData.publicSetting.value.withdraw_threshold}}</span>
          </div>
          <div>
            <el-button  style="width: 150px;margin-right: 20px" type="primary" @click="copyText(state.text);">{{$t('message.finance.copy_invitation_link')}}</el-button>
            <span style="color: #9b9da1;font-size: 12px">{{$t('message.finance.invitation_code')}}: {{userInfos.invitation_code}}</span>
          </div>
        </div>
      </el-card>

      <el-card style="margin-top: 10px;border-radius:10px;padding: 10px">        
        <el-tabs style="color: #6b778c;" v-model="state.tabName" @tab-change="changeTabName()">
          <el-tab-pane name="1">
            <template #label>
              <span class="custom-tabs-label center">
                <el-image :src="paperMoneyIcon" fit="cover" style="height: 20px;width: 20px"></el-image>
                <span>{{$t('message.finance.balance')}}</span>
              </span>
            </template>
            <div>
              <el-table :data="financeStoreData.balanceStatementList.value.data" row-key="id" height="100%" stripe @sort-change="sortChange">
                <el-table-column prop="created_at" :label="$t('message.finance.created_at')" show-overflow-tooltip width="150" fixed sortable="custom">
                  <template #default="{row}">
                    {{DateStrToTime(row.created_at)}}
                  </template>
                </el-table-column>
                <el-table-column prop="title" :label="$t('message.finance.title')" show-overflow-tooltip width="80" sortable="custom"/>
                <el-table-column prop="amount" :label="$t('message.finance.amount')" show-overflow-tooltip width="150" sortable="custom">
                  <template #default="{row}">
                    <span v-if="row.type === constantStore.BALANCE_STATEMENT_TYPE_PLUS">
                      <el-icon color="green"><Plus /></el-icon><span style="color:green;">{{row.amount}}</span>
                    </span>
                    <span v-else>
                      <el-icon color="red"><Minus /></el-icon><span style="color:red;">{{row.amount}}</span>
                    </span>
                  </template>
                </el-table-column>
                <el-table-column prop="final_amount" :label="$t('message.finance.final_amount')" show-overflow-tooltip width="150" sortable="custom"/>
              </el-table>
              <el-pagination background
                             class="mt15"
                             layout="total, sizes, prev, pager, next, jumper"
                             :page-sizes="[10, 30, 50]"
                             v-model:current-page="state.queryParams.pagination.page_num"
                             v-model:page-size="state.queryParams.pagination.page_size"
                             :total="financeStoreData.balanceStatementList.value.total"
                             @size-change="onHandleSizeChange"
                             @current-change="onHandleCurrentChange">
              </el-pagination>
            </div>
          </el-tab-pane>
          <el-tab-pane name="2">
            <template #label>
              <span class="custom-tabs-label center">
                <el-image :src="incomeOneIcon" fit="cover" style="height: 20px;width: 20px"></el-image>
                <span>{{$t('message.finance.commission')}}</span>
              </span>
            </template>
            <div>
              <el-table :data="financeStoreData.commissionStatementList.value.data" row-key="id" height="100%" stripe @sort-change="sortChange">
                <el-table-column prop="created_at" :label="$t('message.finance.created_at')" show-overflow-tooltip width="150" fixed sortable="custom">
                  <template #default="{row}">
                    {{DateStrToTime(row.created_at)}}
                  </template>
                </el-table-column>
                <el-table-column prop="order_user_name" :label="$t('message.finance.order_user_name')" show-overflow-tooltip width="150" sortable="custom">
                  <template #default="{row}">
                    {{replaceCharacters(row.order_user_name)}}
                  </template>
                </el-table-column>
                <el-table-column prop="subject" :label="$t('message.finance.subject')" show-overflow-tooltip width="100" sortable="custom"/>
                <el-table-column prop="commission" :label="$t('message.finance.commission')" show-overflow-tooltip width="150" sortable="custom"/>
                <el-table-column prop="is_withdrew" :label="$t('message.finance.is_withdrew')" show-overflow-tooltip width="150" sortable="custom">
                  <template #default="{row}">
                    <el-button type="info" v-if="row.is_withdrew">{{$t('message.finance.withdrew')}}</el-button>
                    <el-button type="success" v-else>{{$t('message.finance.not_withdrew')}}</el-button>
                  </template>
                </el-table-column>
              </el-table>
              <el-pagination background
                             class="mt15"
                             layout="total, sizes, prev, pager, next, jumper"
                             :page-sizes="[10, 30, 50]"
                             v-model:current-page="state.queryParams.pagination.page_num"
                             v-model:page-size="state.queryParams.pagination.page_size"
                             :total="financeStoreData.commissionStatementList.value.total"
                             @size-change="onHandleSizeChange"
                             @current-change="onHandleCurrentChange">
              </el-pagination>
            </div>
          </el-tab-pane>
          <el-tab-pane name="3">
            <template #label>
              <span class="custom-tabs-label center">
                 <el-image :src="peoplePlusIcon" fit="cover" style="height: 20px;width: 20px"></el-image>
                <span>{{$t('message.finance.invitation_user')}}</span>
              </span>
            </template>
            <div>
              <el-table :data="financeStoreData.invitationUserList.value.data" row-key="id" height="100%" stripe @sort-change="sortChange">
                <el-table-column prop="created_at" :label="$t('message.finance.created_at')" show-overflow-tooltip width="150" fixed sortable="custom">
                  <template #default="{row}">
                    {{DateStrToTime(row.created_at)}}
                  </template>
                </el-table-column>
                <el-table-column prop="user_name" :label="$t('message.finance.order_user_name')" show-overflow-tooltip width="150" sortable="custom">
                  <template #default="{row}">
                    {{replaceCharacters(row.user_name)}}
                  </template>
                </el-table-column>
              </el-table>
              <el-pagination background
                             class="mt15"
                             layout="total, sizes, prev, pager, next, jumper"
                             :page-sizes="[10, 30, 50]"
                             v-model:current-page="state.queryParams.pagination.page_num"
                             v-model:page-size="state.queryParams.pagination.page_size"
                             :total="financeStoreData.invitationUserList.value.total"
                             @size-change="onHandleSizeChange"
                             @current-change="onHandleCurrentChange">
              </el-pagination>
            </div>
          </el-tab-pane>
        </el-tabs>
      </el-card>

  </div>

</template>

<script setup lang="ts">
import paperMoneyIcon from "/@/assets/icon/currency.svg";
import expensesIcon from "/@/assets/icon/expenses.svg";
import giftIcon from "/@/assets/icon/gift-bag.svg";
import peoplePlusIcon from "/@/assets/icon/people-plus.svg";
import incomeOneIcon from "/@/assets/icon/income-one.svg";
import redEnvelopesIcon from "/@/assets/icon/red-envelopes.svg";
import { getCurrentAddress, request } from "/@/utils/request";
import {DateStrToTime} from "/@/utils/formatTime";


import { storeToRefs } from "pinia";
import { useUserStore } from "/@/stores/user_logic/userStore";
import { onMounted, reactive } from "vue";
import { useFinanceStore } from "/@/stores/user_logic/financeStore";
import { usePublicStore } from "/@/stores/publicStore";
import commonFunction from "/@/utils/commonFunction";
import { useConstantStore } from "/@/stores/constantStore";
import { ElMessageBox } from "element-plus";
import { useI18n } from "vue-i18n";

const { copyText } = commonFunction();
const userStore = useUserStore();
const { userInfos } = storeToRefs(userStore);
const pubicStore = usePublicStore()
const pubicStoreData = storeToRefs(pubicStore)
const financeStore = useFinanceStore()
const financeStoreData = storeToRefs(financeStore)
const constantStore = useConstantStore()
const {t} = useI18n()



const state = reactive({
  text: getCurrentAddress()+"/#/login?aff="+userInfos.value.invitation_code,
  tabName: "1",
  queryParams: {
    table_name: "balance_statement",
    pagination: {
      page_num: 1, page_size: 30, order_by: "id DESC"
    } as Pagination//分页参数
  } as QueryParams,
});
const getCommissionSummary=()=>{
  financeStore.getCommissionSummary()
}

const defaultQueryParams=()=>{
  state.queryParams = {
    table_name: "balance_statement",
    pagination: {
      page_num: 1, page_size: 30, order_by: "id DESC"
    } as Pagination//分页参数
  } as QueryParams
}
const changeTabName=()=>{
  switch (state.tabName) {
    case "1":
      defaultQueryParams()
      state.queryParams.table_name = 'balance_statement'
      getData()
      break;
    case "2":
      defaultQueryParams()
      state.queryParams.table_name = 'commission_statement'
      getData()
      break;
    case "3":
      defaultQueryParams()
      state.queryParams.table_name = 'user'
      getData()
      break;
    default:
      break;
  }
}
const getData = () => {
  switch (state.tabName) {
    case "1":
      financeStore.getBalanceStatementList(state.queryParams)
      break;
    case "2":
      financeStore.getCommissionStatementList(state.queryParams)
      break;
    case "3":
      financeStore.getInvitationUserList(state.queryParams)
      break;
    default:
      break;
  }
};

// 分页改变
const onHandleSizeChange = (val: number) => {
  state.queryParams.pagination.page_size = val;
  getData()
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
  state.queryParams.pagination.page_num = val;
  getData()
};
//排序监听
const sortChange = (column: any) => {
  //处理嵌套字段
  let p = (column.prop as string);
  if (p.indexOf(".") !== -1) {
    p = p.slice(p.indexOf(".") + 1);
  }
  switch (column.order) {
    case "ascending":
      state.queryParams.pagination.order_by = p + " ASC";
      break;
    default:
      state.queryParams.pagination.order_by = p + " DESC";
      break;
  }
  getData()
};

const replaceCharacters = (str: string) => {
  return str.substring(0, 4) + '****' + str.substring(8, str.length)
}
const withdrawToBalance=()=>{
  ElMessageBox.confirm("是否把佣金提现到余额？", t('message.common.tip'), {
    confirmButtonText: t('message.common.button_confirm'),
    cancelButtonText: t('message.common.button_cancel'),
    type: 'warning',
  })
    .then(() => {
      financeStore.withdrawToBalance().then(()=>{
        userStore.getUserInfo()
        getCommissionSummary()
        getData()
      })
    })
    .catch(() => {
    });
}

onMounted(()=>{
  getCommissionSummary()
  getData()
});
</script>

<style scoped>
.center {
  text-align: center;
  display: flex;
  justify-content: center;
  align-items: center;
}
.greyText {
  color: #9b9da1;
  font-size: 12px;
}

</style>