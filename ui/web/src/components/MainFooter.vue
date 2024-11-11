<template>
    <v-footer small padless>{{ loginStore.username }}</v-footer>
</template>

<script setup>
import { useLoginStore } from "@/store/auth";

let loginStore = useLoginStore()

</script>
