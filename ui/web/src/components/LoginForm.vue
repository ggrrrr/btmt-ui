<template>
    <v-form>
        <v-container>
            <v-row>
                <v-col cols="4" md="4">
                    <v-text-field autocomplete="current-username" v-model="email" label="Email" type="text"
                        prepend-icon="mdi-account-circle" required hide-details></v-text-field>
                </v-col>
            </v-row>
            <v-row>
                <v-col cols="4" md="4">
                    <v-text-field autocomplete="current-password" v-model="password" label="Password" hide-details
                        prepend-icon="mdi-lock" :append-icon="!showPassword ? 'mdi-eye-off' : 'mdi-eye'"
                        :type="showPassword ? 'text' : 'password'" @click:append="showPassword = !showPassword"
                        required></v-text-field>
                </v-col>
            </v-row>
            <v-row>
                <v-col cols="4" md="4">
                    <v-btn block class="mt-1" @click="loginClick">Login</v-btn>
                </v-col>
            </v-row>
            <v-row>
                <v-col cols="4" md="4">
                    <v-btn block class="mt-1" @click="validateClick">validate</v-btn>
                </v-col>
            </v-row>
        </v-container>
    </v-form>
</template>

<script setup>
import { useLoginStore } from "@/store/auth";
import { ref } from "vue";

const email = ref("")
const password = ref("")
const showPassword = ref(false)

let loginStore = useLoginStore()

function loginClick() {
    console.log("email", email.value)
    loginStore.loginRequest(email.value, password.value)
}
function validateClick() {
    console.log("email", email.value)
    loginStore.validateRequest()
}

</script>