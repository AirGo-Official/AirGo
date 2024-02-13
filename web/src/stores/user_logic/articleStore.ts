import {defineStore, storeToRefs} from "pinia";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";

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
            const res = await request(apiStoreData.userApi.value.getDefaultArticleList)
            let temp: Article[] = res.data.data
            if (temp.length === 2) {
                this.articleID1 = temp[0]
                this.articleID2 = temp[1]
            }
        }
    }
})