import Vue from 'vue'
import App from './App.vue'
import VueResource from 'vue-resource'
import {
    Button,
    Input,
    Aside,
    Form,
    FormItem,
    Row,
    Col,
    Container,
    Main
} from 'element-ui'
import VueRouter from 'vue-router'
import 'element-ui/lib/theme-chalk/index.css'
import '@/assets/var.scss'

Vue.config.productionTip = false
Vue.use(VueRouter)
Vue.use(VueResource)
Vue.use(Button)
Vue.use(Input)
Vue.use(Aside)
Vue.use(Form)
Vue.use(FormItem)
Vue.use(Row)
Vue.use(Col)
Vue.use(Container)
Vue.use(Main)

const router = new VueRouter({
    routes: [

        {
            path: '/',
            component: () =>
                import ('@/views/login')
        },

        {
            path: '/chat',
            component: () =>
                import ('@/views/chat')
        }
    ]
})

new Vue({
    el: '#app',
    router,
    render: h => h(App)
})