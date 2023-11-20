<template>
    <v-btn variant="tonal" rounded="xl" color="#5865f2" @click="emits('click')">
        {{ props['text'] }}
    </v-btn>
</template>

<script setup>
import { defineEmits, defineProps } from 'vue'
const props = defineProps(["text"])
const emits = defineEmits(['click'])
</script>


