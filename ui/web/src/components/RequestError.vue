<template>
    <div>
        {{ errorStore.type }}
    </div>
    <div>
        {{ errorStore.message }}
    </div>
    <div>
        <li v-for="(item, index) in errorStore.errors" :key="index"> {{ index }} {{ item }} </li>
        <!-- {{ errorStore.errors }} -->
    </div>
</template>

<script setup>
import { useErrorStore } from "@/store/error";
let errorStore = useErrorStore()
</script>