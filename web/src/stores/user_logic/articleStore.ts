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
        articleID1: {
            created_at: '',
            updated_at: '',
            // id: 0,
            status: false,
            type: '',
            title: '',
            introduction: '',
            content: '',
        } as Article,
        articleID2: {
            created_at: '',
            updated_at: '',
            // id: 0,
            status: false,
            type: '',
            title: '',
            introduction: '',
            content: '',
        } as Article,
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
                this.articleID1 = temp[0]
                this.articleID2 = temp[1]
            } else {
                const res = await request(apiStoreData.publicApi.value.getDefaultArticleList)
                const total = res.data.total
                let temp: Article[] = res.data.data
                if (total === 2) {
                    this.articleID1 = temp[0]
                    this.articleID2 = temp[1]
                    Session.set('defaultArticles',temp)
                }
            }
        }
    }
})