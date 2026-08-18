import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import WorkDetailView from './views/WorkDetailView.vue'
import './style.css'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'home', component: WorkDetailView },
    { path: '/works/:id', name: 'work', component: WorkDetailView, props: true },
    { path: '/works/:id/files/:path(.*)', name: 'file', component: WorkDetailView, props: true },
  ],
})

createApp(App).use(router).mount('#app')
