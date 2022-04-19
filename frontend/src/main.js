import Vue from 'vue'
import App from './app.vue'
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';

Vue.config.productionTip = false
Vue.config.devtools = false;

Vue.use(ElementUI);

new Vue({
    render: h => h(App),
}).$mount('#app')
