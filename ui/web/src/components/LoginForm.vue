<template>
    <v-sheet class="pa-0 ma-0" rounded width="100%">
        <v-card class="mx-auto px-6 py-8" width="100%">
            <v-form>
                <v-text-field :rules="[required]" :disabled="!props.enabled" class="mb-2" autocomplete="current-username"
                    v-model="email" label="Email" type="email" prepend-icon="mdi-account-circle" required
                    hide-details></v-text-field>
                <v-text-field :rules="[required]" placeholder="Enter your password" :disabled="!props.enabled" class="mb-2"
                    autocomplete="current-password" v-model="password" label="Password" hide-details prepend-icon="mdi-lock"
                    @keydown.enter="handleSubmit" :append-inner-icon="!showPassword ? 'mdi-eye-off' : 'mdi-eye'"
                    :type="showPassword ? 'text' : 'password'" @click:append-inner="showPassword = !showPassword"
                    required></v-text-field>
                <v-card-actions justify="center">
                    <v-btn color="primary" :disabled="!props.enabled" variant="tonal" block rounded="xl"
                        density="comfortable" @click="handleSubmit">Login</v-btn>
                </v-card-actions>
            </v-form>
        </v-card>
    </v-sheet>
</template>

<script setup>
import { ref } from "vue";

const email = ref("")
const password = ref("")
const showPassword = ref(false)

const emits = defineEmits(['submit'])
const props = defineProps(["enabled"])

function required(v) {
    return !!v || 'Field is required'
}

console.log("props.enabled:", props.enabled)

function handleSubmit() {
    console.log("handleSubmit", props.enabled)
    emits("submit", email.value, password.value)
}

</script>