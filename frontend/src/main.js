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

// wails 调用封装
Vue.prototype.wails = function (func, param, success, fail) {
    window.backend.App[func](param).then((resp) => {
        console.log('wails', func, param, resp)
        try {
            var result = JSON.parse(resp);
            if (result.code === 0) {
                success(result.data)
            } else {
                fail(result.data)
            }
        } catch (err) {
            fail(resp)
        }
    });
}

Wails.Init(() => {
    new Vue({
        render: h => h(App),
    }).$mount('#app');
});

