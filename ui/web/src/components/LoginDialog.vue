<template>
    <v-dialog v-model="loginStore.showLogin" width="500">
        <v-card>
            <v-card-title>Login</v-card-title>
            <LoginForm :enabled="enabled" @submit="handleLogin"></LoginForm>
            <v-progress-linear :active="enabled == false" color="deep-purple-accent-4" indeterminate rounded
                height="6"></v-progress-linear>
        </v-card>
    </v-dialog>
</template>

<script setup>
import { useLoginStore } from "@/store/auth";
import LoginForm from '@/components/LoginForm'
import { ref } from "vue";
// import { LoginForm } from "./LoginForm.vue";

const enabled = ref(true)

const loginStore = useLoginStore()

function handleLogin(email, password) {
    enabled.value = false
    loginStore.loginRequest(email, password).then(() => {
    }).finally(() => {
        enabled.value = true
    })
}

</script>