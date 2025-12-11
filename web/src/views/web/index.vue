<template>
    <div class="web">
        <el-container>
            <router-view />
        </el-container>
        <el-footer>
            <web-footer />
        </el-footer>
    </div>
</template>

<script setup lang="ts">
import WebFooter from '@/components/layout/WebFooter.vue';
</script>
<style scoped lang="scss">
.web {
    .el-footer {
        height: auto;
        border: 1px solid #dcdFE6
    }
}
</style>