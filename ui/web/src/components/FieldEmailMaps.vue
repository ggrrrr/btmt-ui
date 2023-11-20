<template>
    <v-container>
        <div v-for="(email, index) in formatEmails()" :key="index">
            {{ email }}
        </div>
    </v-container>
</template>

<script setup>
const props = defineProps(['emails'])

function formatEmails() {
    let out = []
    const emails = props['emails']
    Object.keys(emails).forEach((k) => {
        out.push(emails[k])
    });
    return out
}

</script>