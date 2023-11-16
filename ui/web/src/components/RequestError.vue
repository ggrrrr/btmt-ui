<template>
    <!-- <v-flex xs12 class=""> -->

    <v-dialog v-model="errorStore.show" width="500">
        <v-alert :type="errorStore.alertType()" :text="errorStore.message">
            <div v-for="(item, index) in errorStore.errors" :key="index">{{ index }} {{ item }}</div>
        </v-alert>
    </v-dialog>
    <!-- </v-flex> -->
</template>

<script setup>
import { useErrorStore } from "@/store/error";
// import { ref } from "vue";


let errorStore = useErrorStore()
</script>

