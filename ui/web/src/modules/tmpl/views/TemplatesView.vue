<template>
    <v-container no-gutters>
        <v-container fill-height no-gutters class="ma-0 mp-0">
            <v-row no-gutters>
                <v-col no-gutters cols="12" sm="12" md="12">
                    <v-btn-toggle v-model="refs.page" mandatory>
                        <v-btn value="/templates/manage" to="/templates/manage">Templates</v-btn>
                        <v-btn value="/templates/images" to="/templates/images">Images</v-btn>
                    </v-btn-toggle>
                </v-col>
            </v-row>
            <v-row>
                <v-col no-gutters cols="12" sm="12" md="12">
                    <RouterView />
                </v-col>
            </v-row>
        </v-container>
    </v-container>
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