import {defineStore, storeToRefs} from 'pinia';
import {request} from "/@/utils/request";
import {useApiStore} from "/@/stores/apiStore";

const apiStore = useApiStore()
const apiStoreData = storeToRefs(apiStore)


export const usePublicStore = defineStore('publicStore', {
    state: () => ({
        base64CaptchaData: {
            id: '',
            b64s: '',
        } as Base64CaptchaInfo,

    }),
    actions: {
        //获取base64Captcha
        async getBase64Captcha() {
            const res = await request(apiStoreData.staticApi.value.public_getBase64Captcha)
            this.base64CaptchaData = res.data
        }

    }

})