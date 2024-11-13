/**
 * main.js
 *
 * Bootstraps Vuetify and other plugins then mounts the App`
 */

// Plugins
import { registerPlugins } from "@/plugins";

// Components
// import App from "./App.vue";
import App from "./App.vue";

import OpenLayersMap from "vue3-openlayers";

// console.log("asdasd", Vue.http.options.root);

// Composables
import { createApp } from "vue";

const app = createApp(App).use(OpenLayersMap);

registerPlugins(app);

app.mount("#app");
