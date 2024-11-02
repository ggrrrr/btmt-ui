<template>
    <v-container>
        <div v-for="(phone, index) in formatList()" :key="index">
            {{ phone }}
        </div>
    </v-container>
</template>

<script setup>
const props = defineProps(['labels'])

function formatList() {
    let out = []
    if (props['labels']) {
        props['labels'].forEach((k) => {
            out.push(k)
        });
    }
    return out
}

</script>