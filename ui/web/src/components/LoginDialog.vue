<template>
    <v-dialog v-model="loginStore.showLogin" width="500">
        <LoginForm :enabled="enabled" @submit="handleLogin"></LoginForm>
        <v-progress-linear :active="enabled == false" color="warning" indeterminate rounded height="10"></v-progress-linear>
    </v-dialog>
</template>

<script setup>
import LoginForm from '@/components/LoginForm'

import { useLoginStore } from "@/store/auth";
import { ref } from "vue";

// import { LoginForm } from "./LoginForm.vue";

const enabled = ref(true)

const loginStore = useLoginStore()

// function sleep(ms) {
//     return new Promise(resolve => setTimeout(resolve, ms));
// }
async function handleLogin(email, password) {
    enabled.value = false
    // await sleep(2000);
    await loginStore.loginRequest(email, password).then(() => {
    }).finally(() => {
        enabled.value = true
    })
}

</script>