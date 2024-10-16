<template>
    <v-main>
        <v-card outlined justify="right">
            <v-card-title>
                <v-btn variant="tonal" rounded="xl" color="primary" :enabled="routeIs"
                    to="/templates/images">Images</v-btn>
            </v-card-title>
            <RouterView>
            </RouterView>
        </v-card>
    </v-main>
</template>

<script setup>
// import { useRoute } from 'vue-router'
// const route = useRoute()

function routeIs() {
    // console.log("asd", route.fullPath)
    return ""
}
</script>