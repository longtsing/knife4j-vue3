import { createApp } from 'vue'
import './style/knife4j.less'
import App from './App.vue'
import { setupStore } from './store/index.js'
import router from '@/router/index.js'
import { setupI18n } from '@/lang/index.js'
import { createFromIconfontCN } from '@ant-design/icons-vue'
import { FRAMEWORK_CONFIG } from './config/framework.js'

// 框架标识符（通用）
window.KNIFE4J_FRAMEWORK = 'Generic';
window.KNIFE4J_FRAMEWORK_VERSION = '1.0.0';
// 合并配置：保留 doc.html 中设置的 apiBasePath，添加其他框架配置
window.KNIFE4J_CONFIG = {
  ...window.KNIFE4J_CONFIG,  // 保留 doc.html 中设置的 apiBasePath
  ...FRAMEWORK_CONFIG        // 添加框架配置
};
console.log('🔍 main.js - KNIFE4J_CONFIG after merge:', window.KNIFE4J_CONFIG);

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

// 添加框架配置到全局属性
app.config.globalProperties.$framework = FRAMEWORK_CONFIG;

setupStore(app)
setupI18n(app)

// Knife4j 初始化完成日志
console.log('🚀 Knife4j Vue3 initialized');
console.log('📚 Framework:', window.KNIFE4J_FRAMEWORK, window.KNIFE4J_FRAMEWORK_VERSION);

app.mount('#app')
