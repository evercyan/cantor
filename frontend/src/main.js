import 'core-js/stable';
import 'regenerator-runtime/runtime';
import Vue from 'vue';
import App from './App.vue';

Vue.config.productionTip = true;
Vue.config.devtools = false;

import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
Vue.use(ElementUI);

import * as Wails from '@wailsapp/runtime';

Wails.Init(() => {
    new Vue({
        render: h => h(App),
    }).$mount('#app');
});
