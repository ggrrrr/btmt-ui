<template>
    <v-row class="ma-0 pa-0">
        <div class="ma-0 pa-0">
            <v-text-field style="width: 2ch;" @update:model-value="updateDay" :value="modelValue.day" density="compact"
                variant="underlined"></v-text-field>
        </div>
        <div class="ma-0 pa-0" style="width: 1ch; height: 1ch; text-align: center;">
            /
        </div>
        <div class="ma-0 pa-0">
            <v-text-field style="width: 2ch;" @update:model-value="updateMonth" :value="modelValue.month" density="compact"
                variant="underlined"></v-text-field>
        </div>
        <div style="width: 1ch; text-align: center;">
            /
        </div>
        <div class="ma-0 pa-0">
            <v-text-field style="width: 4ch;" @update:model-value="updateYear" :value="modelValue.year" density="compact"
                variant="underlined"></v-text-field>
        </div>
    </v-row>
</template>

<!--  https://webmound.com/use-v-model-custom-components-vue-3/ -->

<script setup>

import { Dob } from '@/store/people'

const emits = defineEmits(['update:modelValue'])

const props = defineProps({
    modelValue: Dob,
    hint: {
        type: String,
        default: "Date of birth"
    },
    label: {
        type: String,
        default: "birthday",
    }
})

function updateDay(val) {
    let out = props.modelValue
    if (val == "") {
        delete out.day
    } else {
        const n = Number(val)
        if (!isNaN(n)) {
            if (n < 32) {
                out.day = n
            }
        }
    }
    console.log("updateDay", out)
    emits('update:modelValue', out)
}
function updateMonth(val) {
    let out = props.modelValue
    console.log("updateMonth", val)
    if (val == "") {
        delete out.month
    } else {
        const n = Number(val)
        if (!isNaN(n)) {
            if (n < 13) {
                out.month = n
            }
        }
    }
    console.log("updateDay", out)
    emits('update:modelValue', out)
}
function updateYear(val) {
    console.log("updateYear", val)
    let out = props.modelValue
    if (val == "") {
        delete out.year
    } else {
        const n = Number(val)
        if (!isNaN(n)) {
            out.year = n
        }
    }
    console.log("updateDay", out)
    emits('update:modelValue', out)
}

</script>
