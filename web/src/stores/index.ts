import { createPinia } from 'pinia'

const pinia = createPinia()

// 注册插件
pinia.use(({ store })=>{
   if(store.$id ==='website'){
    store.initializeWebsite()
   }
})

export { pinia }