import {createApp} from 'vue';
import pinia from '/@/stores/index';
import App from '/@/App.vue';
import router from '/@/router';
import {directive} from '/@/directive/index';
import other from '/@/utils/other';
import { i18n } from '/@/i18n/index';
import VueLuckyCanvas from '@lucky-canvas/vue'
//vue3 中注册markdown
// @ts-ignore
import VueMarkdownEditor from '@kangc/v-md-editor';
import '@kangc/v-md-editor/lib/style/base-editor.css';
import vuepressTheme from '@kangc/v-md-editor/lib/theme/vuepress.js';
import '@kangc/v-md-editor/lib/theme/style/vuepress.css';
// @ts-ignore
import Prism from 'prismjs';

VueMarkdownEditor.use(vuepressTheme, {
    Prism,
});
//vue3 中注册markdown预览
// @ts-ignore
import VMdPreview from '@kangc/v-md-editor/lib/preview';
import '@kangc/v-md-editor/lib/style/preview.css';
import githubTheme from '@kangc/v-md-editor/lib/theme/github.js';
import '@kangc/v-md-editor/lib/theme/style/github.css';

// highlightjs
import hljs from 'highlight.js';

VMdPreview.use(githubTheme, {
    Hljs: hljs,
});


import ElementPlus from 'element-plus';
import * as ElementPlusIconsVue from '@element-plus/icons-vue' //Element Plus 常用的图标集合

import '/@/theme/index.scss'; //样式


const app = createApp(App);
directive(app);
other.elSvg(app);
app.use(pinia).use(i18n).use(router).use(ElementPlus).use(VueMarkdownEditor).use(VMdPreview).use(VueLuckyCanvas).mount('#app');

for (const [key, component] of Object.entries(ElementPlusIconsVue)) { //Element Plus 常用的图标集合
    app.component(key, component)
}
