<template>
  <div class="logo">
    <!--/image/logo.png 毛笔墨水瓶Blog xxxxxxxxxxxx-->
    <el-image :src="websiteStore.state.websiteInfo.logo === '' ? '/image/logo.png' : websiteStore.state.websiteInfo.logo"
      alt="" />
    <div class="brand-info">
      <div>{{ websiteStore.state.websiteInfo.slogan }}</div>
      <div>{{ websiteStore.state.websiteInfo.slogan_en }}</div>
    </div>
  </div>
</template>


<script setup lang="ts">
//Vue 3 中用 <script setup> 语法时，默认支持组件导出
import { useWebsiteStore } from "@/stores/website";

const websiteStore = useWebsiteStore()
</script>

<style scoped lang="scss">
.dashboard {
  .logo {
    display: flex;
    padding: 10px;

    .el-image {
      height: 60px;
      width: 60px;
    }

    .brand-info {
      margin-top: 5px;
      margin-left: 20px;

      >div:nth-child(1) {
        font-size: 24px;
      }

      >div:nth-child(2) {
        font-size: 13px;
      }
    }
  }


  .collapsed {
    .logo {
      .el-image {
        height: 44px;
        width: 44px;
      }

      .brand-info {
        display: none;
      }
    }
  }
}


.web-navbar .logo {
  display: flex;

  .el-image {
    padding: 5px;
  }

  .brand-info {
    margin-top: 5px;

    >div:nth-child(1) {
      font-size: 22px;

    }

    >div:nth-child(2) {
      font-size: 12px;
    }
  }
}
</style>
