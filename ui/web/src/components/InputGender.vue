<template>
    <v-btn-toggle @update:model-value="update" :model-value="props.modelValue" color="secondary" rounded="xl"
        variant="flat">
        <v-btn value="male">
            <v-icon size="large">mdi-human-male</v-icon>
        </v-btn>
        <v-btn value="">
            <v-icon size="large">mdi-empty</v-icon>
        </v-btn>
        <v-btn value="female">
            <v-icon size="large">mdi-human-female</v-icon>
        </v-btn>
    </v-btn-toggle>
</template>
<script setup>

import { ref } from 'vue'

const refs = ref({
    g: "male"
})

const emits = defineEmits(['update:modelValue'])

const props = defineProps({
    modelValue: String
})

function update(val) {
    // console.log("update:modelValue", val)
    emits('update:modelValue', val)
}

</script>


