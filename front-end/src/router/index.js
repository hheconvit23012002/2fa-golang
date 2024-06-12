import { createRouter, createWebHistory } from 'vue-router'
import Login from '@/views/Login.vue'
import AuthTwoFA from "@/views/AuthTwoFA.vue";
import SuccessLogin from "@/views/SuccessLogin.vue";
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'login',
      component: Login
    },
    {
      path: '/2fa',
      name: '2fa',
      component: AuthTwoFA,
    },
    {
      path: '/loginSuccess',
      name: 'loginSuccess',
      component: SuccessLogin
    }
  ]
})

router.beforeEach(async (to, from, next) => {
  // redirect to login page if not logged in and trying to access a restricted page
  const token = localStorage.getItem('token');
  if(to.matched.some(record => record.meta.notLogin) && token){
    next(from.fullPath);
  }
  else if (to.matched.some(record => record.meta.requiresAuth) && !token) {
    next('/login');
  } else {
    next();
  }
});
export default router

