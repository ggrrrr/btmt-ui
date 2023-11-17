<script setup>
import { useLoginStore } from "@/store/auth";

let loginStore = useLoginStore()

</script>

<template>
    <v-footer small padless>{{ loginStore.email }}</v-footer>
</template>
