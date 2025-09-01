import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import AdminLayout from '../views/admin/AdminLayout.vue'
import EventsAdmin from '../views/admin/EventsAdmin.vue'
import DateTemplateGroupsAdmin from '../views/admin/DateTemplateGroupsAdmin.vue'
import DateTemplatesAdmin from '../views/admin/DateTemplatesAdmin.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/admin',
    component: AdminLayout,
    children: [
      {
        path: 'events',
        name: 'AdminEvents',
        component: EventsAdmin
      },
      {
        path: 'date-template-groups',
        name: 'AdminDateTemplateGroups', 
        component: DateTemplateGroupsAdmin
      },
      {
        path: 'date-templates',
        name: 'AdminDateTemplates',
        component: DateTemplatesAdmin
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router