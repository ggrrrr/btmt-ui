<template>
    <v-container>
        <div>
            {{ formatDoB() }}
        </div>
    </v-container>
</template>

<script setup>
const props = defineProps(['dob'])

function formatDoB() {
    const ts = props['dob']
    if (ts) {
        return `${ts.day} / ${ts.month} / ${ts.year}`
    }
    return ''
}

</script>