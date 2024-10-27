<template>
    <v-main>
        <v-card outlined justify="right">
            <v-card-title>
                <div d-flex flex-column align-center bg-grey-lighten-4 pa-6>

                    <v-btn-toggle v-model="refs.page" mandatory>
                        <v-btn value="/templates/images" to="/templates/images">Images</v-btn>

                        <v-btn value="/templates/manage" to="/templates/manage">Templates</v-btn>
                    </v-btn-toggle>
                </div>
            </v-card-title>
            {{ refs }}
            <RouterView>
            </RouterView>
        </v-card>
    </v-main>
</template>

<script setup>
import { ref } from 'vue';
import { useRoute } from 'vue-router'
const route = useRoute()

const refs = ref({
    page: route.fullPath,
}
)
console.log("route", route)
</script>