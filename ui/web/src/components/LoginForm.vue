<template>
    <v-sheet class="ma-2 pa-2" width="100%">
        <v-form>
            <v-text-field :disabled="!props.enabled" class="mr-2" autocomplete="current-username" v-model="email"
                label="Email" type="email" prepend-icon="mdi-account-circle" required hide-details></v-text-field>
            <v-text-field :disabled="!props.enabled" class="mr-2" autocomplete="current-password" v-model="password"
                label="Password" hide-details prepend-icon="mdi-lock"
                :append-inner-icon="!showPassword ? 'mdi-eye-off' : 'mdi-eye'" :type="showPassword ? 'text' : 'password'"
                @click:append-inner="showPassword = !showPassword" required></v-text-field>
            <v-card-actions justify="center">
                <v-btn :disabled="!props.enabled" block rounded density="comfortable" @click="handleSubmit">Submit</v-btn>
            </v-card-actions>
        </v-form>
    </v-sheet>
</template>

<script setup>
import { defineEmits, defineProps } from 'vue'
import { ref } from "vue";

const email = ref("")
const password = ref("")
const showPassword = ref(false)

const emits = defineEmits(['submit'])
const props = defineProps(["enabled"])



console.log("props.enabled:", props.enabled)

function handleSubmit() {
    console.log("handleSubmit", props.enabled)
    emits("submit", email.value, password.value)
}

</script>