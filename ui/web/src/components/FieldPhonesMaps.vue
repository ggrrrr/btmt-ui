<template>
    <v-container>
        <div v-for="(phone, index) in formatMap()" :key="index">
            {{ phone }}
        </div>
    </v-container>
</template>

<script setup>
const props = defineProps(['phones'])

function formatMap() {
    let out = []
    const phones = props['phones']
    Object.keys(phones).forEach((k) => {
        out.push(phones[k])

    });
    return out
}

</script>