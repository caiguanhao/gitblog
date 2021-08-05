import { createRouter, createWebHistory } from 'vue-router'

import RoutePosts from './posts/index.vue'
import RoutePostsNew from './posts/new.vue'
import RoutePostsEdit from './posts/edit.vue'
import RouteSync from './status/sync.vue'
import RouteError from './errors/index.vue'
import RouteBlank from './errors/blank.vue'

const router = createRouter({
  scrollBehavior (to, from, savedPosition) {
    if (to.hash) {
      return { el: to.hash }
    } else if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  },
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'RouteHome', component: RoutePosts, alias: [ '/posts' ] },
    { path: '/posts/new', name: 'RoutePostsNew', component: RoutePostsNew },
    { path: '/posts/:id/edit', name: 'RoutePostsEdit', component: RoutePostsEdit },
    { path: '/sync', name: 'RouteSync', component: RouteSync },
    { path: '/error', name: 'RouteError', component: RouteError },
    { path: '/\n', name: 'RouteBlank', component: RouteBlank },
    { path: '/:pathMatch(.*)*', name: 'RouteNotFound', component: RouteError }
  ],
})

router.$lastRoute = null
router.$lastError = null
router.$vm = null
router.setVM = vm => router.$vm = vm

router.beforeEach((to, from, next) => {
  if (to.name === 'RouteError' || to.name === 'RouteBlank') return next()
  router.$lastRoute = to
  if (router.$vm) {
    router.$vm.getCurrentStatus().then(next, next)
    return
  }
  next()
})

router.onError((err) => {
  router.$lastError = err
  router.push({ name: 'RouteError' })
})

export default router
