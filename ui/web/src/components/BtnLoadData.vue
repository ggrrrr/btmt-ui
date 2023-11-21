<template>
    <v-btn variant="tonal" rounded="xl" color="primary" @click="emits('click')">
        {{ props['text'] }}
    </v-btn>
</template>

<script setup>
// import { } from 'vue'
const props = defineProps(["text"])
const emits = defineEmits(['click'])
</script>


