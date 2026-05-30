import { createApp } from 'vue'
import './style/knife4j.less'
import App from './App.vue'
import { setupStore } from './store/index.js'
import router from '@/router/index.js'
import { setupI18n } from '@/lang/index.js'
import { createFromIconfontCN } from '@ant-design/icons-vue'
import { LITESTAR_CONFIG } from './config/litestar.js'

// LiteStar 标识符
window.KNIFE4J_FRAMEWORK = 'LiteStar';
window.KNIFE4J_FRAMEWORK_VERSION = '2.17.0';
window.LITESTAR_CONFIG = LITESTAR_CONFIG;

String.prototype.gblen = function () {
  let len = 0
  for (let i = 0; i < this.length; i++) {
    if (this.charCodeAt(i) > 127 || this.charCodeAt(i) == 94) {
      len += 2;
    } else {
      len++;
    }
  }
  return len;
}

String.prototype.startWith = function (str) {
  const reg = new RegExp("^" + str)
  return reg.test(this);
}

/***
 * 自定义图标
 */
import iconFront from './assets/iconfonts/iconfont.js'
const MyIcon = createFromIconfontCN({
  scriptUrl: iconFront
})

const app = createApp(App)
app.use(router)
app.component('my-icon', MyIcon)

// 添加 LiteStar 配置到全局属性
app.config.globalProperties.$litestar = LITESTAR_CONFIG;

setupStore(app)
setupI18n(app)

// LiteStar 初始化完成日志
console.log('🚀 Knife4j Vue3 + LiteStar initialized');
console.log('📚 Framework:', window.KNIFE4J_FRAMEWORK, window.KNIFE4J_FRAMEWORK_VERSION);

app.mount('#app')
