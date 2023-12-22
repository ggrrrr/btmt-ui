<template>
    <v-container>
        <div>
            {{ formatTimeStamp() }}
        </div>
    </v-container>
</template>

<script setup>
const props = defineProps(['timeStamp'])

function formatTimeStamp() {
    const ts = props['timeStamp']
    if (ts.seconds) {
        const dateObj = new Date(ts.seconds * 1000);
        return dateObj.toLocaleString()
    }
    const dateObj = new Date(ts);
    return dateObj.toLocaleString()

}

</script>