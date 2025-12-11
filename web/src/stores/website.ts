import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { Website } from '@/api/config';
import {websiteInfo} from "@/api/website";

function initState() {
    const websiteInfo = ref<Website>({
    logo: '',
    full_logo: '',
    title: '',
    slogan: '',
    slogan_en: '',
    description: '',
    version: '',
    created_at: '',
    icp_filing: '',
    public_security_filing: '',
    bilibili_url: '',
    gitee_url: '',
    github_url: '',
    name: '',
    job: '',
    address: '',
    email: '',
    qq_image: '',
    wechat_image: '',
    })
    return {websiteInfo,websiteInfoInitiakized:false,}
}
//useXXXStore 是 Pinia 推荐的命名方式（以 use 开头，以 Store 结尾
export const useWebsiteStore = defineStore('website',() =>{
    const state = ref(initState())
//如果网站信息尚未初始化，则发送 API 请求获取数据；请求成功后，更新状态并标记为已初始化，避免重复请求。
    const initializeWebsite = async () => {
        if (!state.value.websiteInfoInitiakized){
            const res = await websiteInfo()
            if (res.code ===0){
                state.value.websiteInfo = res.data
                state.value.websiteInfoInitiakized = true
            }
        }
    }
    return{state,initializeWebsite}
})
//通过 defineStore 注册了一个 ID 为 website 的仓库，封装了网站信息相关的状态和操作逻辑。
//通过 export 将仓库的访问函数 useWebsiteStore 导出，供整个应用的组件使用。