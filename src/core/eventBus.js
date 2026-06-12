/**
 * 简单事件总线 - 用于 Vue 3 组件间通信
 * 替代 Vue 2 的 $on/$emit 模式
 */

class EventBus {
  constructor() {
    this.events = {};
  }

  /**
   * 监听事件
   * @param {string} event 事件名
   * @param {Function} callback 回调函数
   */
  on(event, callback) {
    if (!this.events[event]) {
      this.events[event] = [];
    }
    this.events[event].push(callback);
  }

  /**
   * 移除事件监听
   * @param {string} event 事件名
   * @param {Function} callback 回调函数
   */
  off(event, callback) {
    if (!this.events[event]) return;
    const index = this.events[event].indexOf(callback);
    if (index > -1) {
      this.events[event].splice(index, 1);
    }
  }

  /**
   * 触发事件
   * @param {string} event 事件名
   * @param  {...any} args 传递给回调的参数
   */
  emit(event, ...args) {
    if (!this.events[event]) return;
    this.events[event].forEach(callback => {
      try {
        callback(...args);
      } catch (e) {
        console.error(`Error in event handler for ${event}:`, e);
      }
    });
  }

  /**
   * 监听一次性事件
   * @param {string} event 事件名
   * @param {Function} callback 回调函数
   */
  once(event, callback) {
    const wrapper = (...args) => {
      callback(...args);
      this.off(event, wrapper);
    };
    this.on(event, wrapper);
  }
}

// 创建全局事件总线实例
export const eventBus = new EventBus();

// 全局参数更新事件
export const GLOBAL_PARAMETERS_UPDATED = 'global-parameters-updated';
