import Vue from 'vue'
import Router from 'vue-router'
import dashboard from '@/components/pages/dashboard'
import chart from '@/components/pages/chart'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: dashboard
    },
    {
      path: '/chart',
      name: 'chart',
      component: chart
    }
  ]
})
