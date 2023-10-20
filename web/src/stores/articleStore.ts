import {defineStore, storeToRefs} from "pinia";
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";

const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)
export const useArticleStore = defineStore("articleStore", {
    state: () => ({
        articleDate: {
            total: 0,
            article_list: [] as Article[],
        },
        currentArticle: {
            status: true,
            type: 'notice',
            title: '',
            introduction: '',
            content: '',} as Article,
        articleID1:{
            created_at: '',
            updated_at: '',
            // id: 0,
            status: false,
            type: '',
            title: '',
            introduction: '',
            content: '',
        } as Article,
        articleID2:{
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
            const res = await request(apiStoreData.api.value.article_getArticle, params)
            // this.articleDate = res.data
            return res
        }
    }
})