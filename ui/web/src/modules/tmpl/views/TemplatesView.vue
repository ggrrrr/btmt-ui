<template>
    <v-main>
        <v-container no-gutters>
            <div d-flex flex-column align-center bg-grey-lighten-4 pa-6>
                <v-btn-toggle v-model="refs.page" mandatory>
                    <v-btn value="/templates/images" to="/templates/images">Images</v-btn>

                    <v-btn value="/templates/manage" to="/templates/manage">Templates</v-btn>
                </v-btn-toggle>
            </div>
            <div>
                <RouterView>
                </RouterView>
            </div>

        </v-container>
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
</script>