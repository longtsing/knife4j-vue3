<template>
  <a-menu :style="style" class="contextmenu" v-show="visible" @click="handleClick" :selectedKeys="selectedKeys">
    <a-menu-item :key="item.key" v-for="item in itemList" role="menuitem">
      <LeftOutlined v-if="item.icon === 'LeftOutlined'" role="menuitemicon" style="margin-right: 8px;" />
      <RightOutlined v-if="item.icon === 'RightOutlined'" role="menuitemicon" style="margin-right: 8px;" />
      <CloseOutlined v-if="item.icon === 'CloseOutlined'" role="menuitemicon" style="margin-right: 8px;" />
      <span>{{item.text}}</span>
    </a-menu-item>
  </a-menu>
</template>

<script>
import { LeftOutlined, RightOutlined, CloseOutlined } from '@ant-design/icons-vue'
export default {
  name: "Contextmenu",
  components: {
    LeftOutlined,
    RightOutlined,
    CloseOutlined,
  },
  props: {
    visible: {
      type: Boolean,
      required: false,
      default: false
    },
    itemList: {
      type: Array,
      required: true,
      default: () => []
    },
    x: {
      type: Number,
      required: false,
      default: 0
    },
    y: {
      type: Number,
      required: false,
      default: 0
    }
  },
  emits: ['update:visible', 'select'],
  data() {
    return {
      left: 0,
      top: 0,
      target: null,
      selectedKeys: []
    };
  },
  computed: {
    style() {
      return {
        left: (this.x || this.left) + "px",
        top: (this.y || this.top) + "px"
      };
    }
  },
  created() {
    // 延迟添加mousedown监听器，避免立即关闭菜单
    setTimeout(() => {
      window.addEventListener("mousedown", e => this.closeMenu(e));
    }, 300); // 增加延迟时间
    window.addEventListener("contextmenu", e => this.setPosition(e));
  },
  methods: {
    closeMenu(e) {
      // 如果菜单当前不可见，不需要处理关闭逻辑
      if (!this.visible) {
        return;
      }
      
      console.log("ContextMenu closeMenu called, target:", e.target);
      console.log("Target class:", e.target.className);
      console.log("Target role:", e.target.getAttribute("role"));
      
      // 更精确的判断：如果点击的是菜单相关元素，不关闭菜单
      const isMenuElement = e.target.closest('.contextmenu') || 
                           e.target.closest('.ant-menu') ||
                           e.target.closest('.ant-menu-item') ||
                           e.target.hasAttribute('data-v-6cc60603') || // Vue scoped CSS
                           ["menuitemicon", "menuitem"].indexOf(e.target.getAttribute("role")) >= 0;
      
      if (isMenuElement) {
        console.log("Not closing menu - clicked on menu or menu item");
        return;
      }
      
      console.log("Closing menu via mousedown");
      this.$emit("update:visible", false);
    },
    setPosition(e) {
      this.left = e.clientX;
      this.top = e.clientY;
      this.target = e.target;
    },
    handleClick({ key }) {
      console.log("ContextMenu handleClick triggered with key:", key);
      // 立即发送选择事件
      this.$emit("select", key, this.target);
      // 延迟关闭菜单，确保事件正确传递
      setTimeout(() => {
        this.$emit("update:visible", false);
      }, 5);
    }
  }
};
</script>

<style lang="less" scoped>
.contextmenu {
  position: fixed;
  z-index: 10000;
  border: 1px solid #9e9e9e;
  border-radius: 4px;
  box-shadow: 2px 2px 10px #aaaaaa !important;
}
</style>
