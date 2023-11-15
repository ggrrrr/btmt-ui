<template>
    <div>
        {{ errorStore.type }}
    </div>
    <div>
        {{ errorStore.message }}
    </div>
    <div>
        {{ errorStore.error }}
    </div>
</template>

<script setup>
import { useErrorStore } from "@/store/error";
let errorStore = useErrorStore()
</script>