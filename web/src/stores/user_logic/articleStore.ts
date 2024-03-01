import {defineStore, storeToRefs} from "pinia";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";
import { Session } from "/@/utils/storage";

const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
export const useArticleStore = defineStore("articleStore", {
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
        defaultArticles:[] as Article[],
    }),
    actions: {
        //获取article列表
        async getArticleList(params: object) {
            const res = await request(apiStoreData.userApi.value.getArticleList, params)
            this.articleList = res.data
        },
        async getDefaultArticles() {
            //尝试从session中获取
            if (Session.get('defaultArticles')){
                let temp: Article[] = Session.get('defaultArticles')
                this.defaultArticles = temp
            } else {
                const res = await request(apiStoreData.publicApi.value.getDefaultArticleList)
                let temp: Article[] = res.data.data
                this.defaultArticles = temp
                Session.set('defaultArticles',temp)
            }
            return
        }
    }
})