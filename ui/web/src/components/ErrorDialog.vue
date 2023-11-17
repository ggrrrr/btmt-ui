<template>
    <v-dialog v-model="errorStore.show" width="500">
        <v-alert :type="errorStore.alertType()" :text="errorStore.message">
            <div v-for="(item, index) in errorStore.errors" :key="index">{{ index }} {{ item }}</div>
        </v-alert>
    </v-dialog>
</template>

<script setup>
import { useErrorStore } from "@/store/error";

let errorStore = useErrorStore()
</script>

