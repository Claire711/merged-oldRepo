<template>
  <div :class="{ 'web-navbar': true, show: isShow }">
    <div class="container">
      <logo />
      <div class="web-menu">
      <!--:ellipsis="false":控制当菜单项过多、超出容器宽度时，是否自动省略溢出的菜单项。-->
       <!--:router="true" 启用「路由模式」，让菜单与 Vue Router 联动，实现点击菜单项自动跳转路由-->
        <el-menu mode="horizontal" :ellipsis="false" :router="true" :default-active="$route.path">
          <template v-for="item in menuList">
            <el-menu-item :index="item.name"><span>{{ item.title }}</span></el-menu-item>
          </template>
        </el-menu>
      </div>
      <auth-popover />
    </div>
  </div>
</template>

<script setup lang="ts">
import AuthPopover from "@/components/common/AuthPopover.vue";
import Logo from "@/components/widgets/Logo.vue";
import { ref } from "vue";
import { onUnmounted } from "vue";

const isShow = ref(true)

const props = defineProps<{
  noScroll?: boolean
}>()

if (!props.noScroll) {
  isShow.value = false// 初始化：先把元素设为隐藏（覆盖初始的 true）
  window.addEventListener("scroll", scroll)
  scroll()
}

function scroll() {
  let top = document.documentElement.scrollTop// 获取页面当前的「滚动距离」（从顶部开始算）
  isShow.value = top >= 100;// 核心判断：滚动距离 > 100px → 显示（true）；否则隐藏（false）
}

onUnmounted(() => {
  if (!props.noScroll) {
    window.removeEventListener("scroll", scroll)
  }
})

interface MenuItem {
  title: string;
  name: string;
}

const menuList: MenuItem[] = [
  {
    title: "首页",
    name: "/",
  },
  {
    title: "搜索",
    name: "/search",
  },
  {
    title: "新闻",
    name: "/news",
  },
  {
    title: "友链",
    name: "/friend-link",
  },
  {
    title: "关于",
    name: "/about",
  }
]

</script>


<style scoped lang="scss">
.web-navbar {
  display: flex;
  justify-content: center;
  width: 100%;
  position: fixed;
  z-index: 6;
  color: dimgray;
  --el-menu-text-color: dimgray;//<el-menu> 里的菜单文字（比如 “首页”“文章” 这些选项）会变成暗灰色。
  --color: dimgray;//控制导航栏里除了 Element 的 <el-menu>，自己写的元素（比如 Logo 旁边的文字、小图标等）颜色。

  &.show {
    top: 0;
    background-color: white;
    color: black;
    --el-menu-text-color: black;
    --color: black;

    .container {
      margin-top: 8px;
    }
  }

  .container {
    display: flex;
    max-width: 1400px;
    width: 100%;

    .logo {
      height: 60px;
      width: 200px;
    }

    .web-menu {
      margin-left: 20px;

      .el-menu {
        background-color: transparent;
        border-bottom: none;
        --el-menu-item-font-size: 20px;

        .el-menu-item {
          border-bottom: none;
          background-color: transparent;
        }
      }
    }

    .auth-popover {
      margin-left: auto;
      margin-top: auto;
      margin-bottom: auto;
      padding-right: 20px;
    }
  }
}
</style>
