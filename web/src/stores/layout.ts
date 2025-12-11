import {defineStore} from "pinia";
import {ref, watch} from 'vue';

function initState() {
    const savedIsCollapse = localStorage.getItem('isCollapse');
    return {
        isCollapse: savedIsCollapse === 'true',
        popoverVisible: false,
        loginVisible: false,
        registerVisible: false,
        forgotPasswordVisible: false,
        passwordResetVisible:false,
        shouldRefreshUserTable:false,
        shouldRefreshImageTable:false,
        articleCreateVisible:false,
        articleUpdateVisible:false,
        shouldRefreshArticleTable:false,
        friendLinkCreateVisible: false,
        friendLinkUpdateVisible: false,
        shouldRefreshFriendLinkTable: false,
        advertisementCreateVisible:false,
        advertisementUpdateVisible:false,
        shouldRefreshAdvertisementTable: false,
        feedbackReplyVisible:false,
        shouldRefreshFeedbackTable:false,
        shouldRefreshCommentTable:false,
        shouldRefreshCommentList:false,
    };
}

export const useLayoutStore = defineStore('layout', () => {
    const state = ref(initState());

    // 监视 isCollapse 的变化并更新 localStorage
    /*watch(监听的目标, 变化时的回调函数, [配置项])
    第一个参数：() => state.value.isCollapse（监听目标）
    第二个参数：(newIsCollapse) => { ... }（变化时的回调）
    */
    watch(() => state.value.isCollapse, (newIsCollapse) => {
        localStorage.setItem('isCollapse', String(newIsCollapse));
    });

    return {state};
});