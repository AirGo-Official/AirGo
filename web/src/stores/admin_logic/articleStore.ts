import {defineStore, storeToRefs} from "pinia";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";

const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
export const useAdminArticleStore = defineStore("adminArticleStore", {
  state: () => ({
    articleList: {
      total: 0,
      data: [] as Article[],
    },
    currentArticle: {
      status: true,
      type: 'notice',
      title: '',
      introduction: '',
      content: '',
    } as Article,
  }),
  actions: {
    async newArticle(){
      return request(apiStoreData.adminApi.value.newArticle,this.currentArticle)
    },
    async getArticleList(params: QueryParams) {
      const res = await request(apiStoreData.adminApi.value.getArticleList, params)
      this.articleList = res.data
    },
    async updateArticle(){
      return request(apiStoreData.adminApi.value.updateArticle,this.currentArticle)
    },
    async deleteArticle(params:Article){
      return request(apiStoreData.adminApi.value.newArticle,params)
    },
  }
})