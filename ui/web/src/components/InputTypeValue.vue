<template>
    <v-row width="300">
        <v-col class="pr-1" cols="3">
            <v-select v-model="type" :label="props.typeLabel" :items="props.typeItems" density="compact"
                variant="underlined"></v-select>
        </v-col>
        <v-col class="pl-0" cols="6">
            <v-text-field :rules="props.rules" v-model="value" @update:focused="onFocus" :label="props.label"
                :hint="props.hint" @keydown.enter="handleAdd" density="compact" variant="underlined"></v-text-field>
        </v-col>
    </v-row>
</template>
<!--  https://webmound.com/use-v-model-custom-components-vue-3/ -->
<script setup>

import { ref, onMounted } from 'vue'

const props = defineProps({
    modelValue: Object,
    rules: Array,
    typeLabel: {
        type: String,
        default: "type",
    },
    typeItems: {
        type: Array,
        required: true,
    },
    hint: {
        type: String,
        default: "More info here"
    },
    label: {
        type: String,
        default: "Label",
    }
})

const type = ref("")
const value = ref("")
const emits = defineEmits(['update:modelValue'])

onMounted(() => {
    if (props) {
        if (props.typeItems) {
            if (type.value == "") {
                if (props.typeItems.length > 0) {
                    type.value = props.typeItems[0]
                }
            }
        }
    }
})

function onFocus(asd) {
    if (asd == false) {
        handleAdd()
    }
}

function handleAdd() {
    let rules = true
    if (props.rules) {
        props.rules.forEach(f => {
            const rule = f(value.value)
            if (rule !== true) {
                rules = false
            }
        })
    }
    if (rules) {
        let temp = props.modelValue
        temp[type.value] = value.value
        // temp.set(type.value, value.value)
        console.log("modelVelue", temp)
        emits('update:modelValue', temp)
    }
}

</script>


