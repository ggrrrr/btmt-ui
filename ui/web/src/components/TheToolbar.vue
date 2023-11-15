<template>
    <v-toolbar app dense flat class="ma-0 pa-0">
        <v-toolbar-title small>Vuetify</v-toolbar-title>
        <v-spacer flat></v-spacer>
        <v-label>{{ loginStore.email }}</v-label>
        <v-icon class="mr-1">mdi-export</v-icon>
    </v-toolbar>
</template>

<script setup>
import { useLoginStore } from "@/store/auth";

let loginStore = useLoginStore()

</script>
