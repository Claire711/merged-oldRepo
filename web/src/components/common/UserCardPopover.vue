<template>
<div class="user-card-popover">
  <el-popover width="280">
    <template #reference>
      <el-avatar :src="childMessage.avatar"/>
    </template>
    <template #default>
      <user-card :uuid="props.uuid" :user-card-info="null" @userCardInfo="handleDataFromChild"/>
    </template>
  </el-popover>
</div>
</template>

<script setup lang="ts">
import UserCard from "@/components/widgets/UserCard.vue";
import {defineProps, ref} from "vue";
import type {UserCardResponse} from "@/api/user";

const props = defineProps<{
  uuid: string;
}>();

const childMessage=ref<UserCardResponse>({
  uuid: "",
  username: "",
  avatar: "",
  address: "",
  signature: "",
})
//获取到用户信息后，UserCard 通过 emit('userCardInfo', data) 触发事件，将数据传递给 UserCardPopover。
//处理子组件事件
const handleDataFromChild = (message:UserCardResponse) => {
  childMessage.value = message// 更新本地数据
}
</script>


<style scoped lang="scss">

</style>
