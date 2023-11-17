// Composables
import { createRouter, createWebHistory } from "vue-router";

const routes = [
  {
    path: "/",
    component: () => import("@/layouts/default/MainView.vue"),
    meta: {
      title: "Main",
    },
    children: [
      {
        path: "",
        name: "Home",
        // route level code-splitting
        // this generates a separate chunk (Home-[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () => import("@/views/HomeView.vue"),
      },
      {
        path: "/people",
        name: "People",
        // route level code-splitting
        // this generates a separate chunk (Home-[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () => import("@/views/PeopleView.vue"),
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

// router.beforeEach((to, from, next) => {
//   console.log("beforeEach", to.fullPath, from.fullPath);
//   next();
// });

// router.afterEach((to, from) => {
//   console.log("afterEach", to.meta, from.meta);
//   if (to.meta.title) {
//     document.title = to.meta.title;
//   }
// });

export default router;
