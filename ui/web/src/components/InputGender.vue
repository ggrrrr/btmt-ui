<template>
    <v-btn-toggle @update:model-value="update" :model-value="props.modelValue" color="primary" rounded="xl"
        variant="outlined">
        <v-btn title="male" value="male">
            <v-icon x size="large">mdi-gender-male</v-icon>
        </v-btn>
        <v-btn title="not set" value="">
            <v-icon size="large">mdi-empty</v-icon>
        </v-btn>
        <v-btn title="female" value="female">
            <v-icon size="large">mdi-gender-female</v-icon>
        </v-btn>
    </v-btn-toggle>
</template>
<script setup>

const emits = defineEmits(['update:modelValue'])

const props = defineProps({
    modelValue: String
})

function update(val) {
    emits('update:modelValue', val)
}

</script>


