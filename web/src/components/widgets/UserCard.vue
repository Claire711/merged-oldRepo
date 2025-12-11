<template>
    <div class="user-card">
        <div class="user-info">
            <el-avatar :size="50" :src="userCardInfo.avatar" />
            <el-row>{{ userCardInfo.username }}</el-row>
            <el-row>{{ userCardInfo.address }}</el-row>
        </div>
        <el-row class="uuid">uuid：{{ userCardInfo.uuid }}</el-row>
        <div class="container">
            <el-row>签名：{{ userCardInfo.signature }}</el-row>
        </div>
    </div>
</template>
<script setup lang="ts">
import { userCard, type UserCardRequest, type UserCardResponse } from '@/api/user';
import { defineProps, ref } from "vue";

const props = defineProps<{
    uuid: string,
    userCardInfo: UserCardResponse | null;
}>();

const userCardInfo = ref<UserCardResponse>({
    uuid: '',
    username: '',
    avatar: '',
    address: '',
    signature: '',
})

const getUserCard = async () => {
    if (props.userCardInfo) {
        userCardInfo.value = props.userCardInfo;
    } else {
        const userCardRequest: UserCardRequest = {
            uuid: props.uuid,
        }
        const res = await userCard(userCardRequest)
        if (res.code === 0) {
            userCardInfo.value = res.data;
            sendDataToParent(res.data)
        }
    }
}

getUserCard()
//定义事件
const emit = defineEmits<{
    (event: 'userCardInfo', message: UserCardResponse): void;
}>();
// 触发事件
const sendDataToParent = (userCardInfo: UserCardResponse) => {
    emit('userCardInfo', userCardInfo)
}
</script>

<style scoped lang="scss">
.user-card {
    max-width: 320px;
    padding-top: 40px;

    .user-info {
        text-align: center;

        .el-row {
            display: flex;
            justify-content: center;
            margin: 5px;
        }
    }

    .uuid {
        margin: 10px 20px;
    }

    .container {
        background-color: white;

        .el-row {
            padding: 20px;
        }
    }
}
</style>