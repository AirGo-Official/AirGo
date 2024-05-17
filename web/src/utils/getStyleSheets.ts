import {nextTick} from 'vue';
import * as svg from '@element-plus/icons-vue';



// 初始化获取 css 样式，获取 element plus 自带 svg 图标，增加了 ele- 前缀，使用时：ele-Aim
const getElementPlusIconfont = () => {
    return new Promise((resolve, reject) => {
        nextTick(() => {
            const icons = svg as any;
            const sheetsIconList = [];
            for (const i in icons) {
                sheetsIconList.push(`ele-${icons[i].name}`);
            }
            if (sheetsIconList.length > 0) resolve(sheetsIconList);
            else reject('未获取到值，请刷新重试');
        });
    });
};



const initIconfont = {

    // element plus
    ele: () => {
        return getElementPlusIconfont();
    },


};

// 导出方法
export default initIconfont;
