<template>
    <v-chip v-for="(name, index) in props.list.list" :key="index" @click="delSearchField(index)" class="ga-2 ma-2" rounded
        density="comfortable" size="large" variant="tonal" :color="props.color">
        <template v-slot>
            {{ name }}
        </template>
        <template v-slot:append>
            <v-icon color="warning">mdi-delete</v-icon>
        </template>
    </v-chip>
</template>

<script setup>
let emits = defineEmits(['click'])
let props = defineProps({
    list: Object,
    color: String,
})

function delSearchField(index) {
    emits('click', index)
}


</script>