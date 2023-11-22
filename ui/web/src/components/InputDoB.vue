<template>
    <div>
        <v-text-field class="input-day" v-model="value.d" @update:focused="onFocusDay" @keydown.enter="handleDay"
            density="compact" variant="underlined"></v-text-field>
    </div>
    <div>
        /
    </div>
    <div>
        <v-text-field class="input-month" v-model="value.m" @update:focused="onFocusMonth" @keydown.enter="handleMonth"
            density="compact" variant="underlined"></v-text-field>
    </div>
    <div>
        /
    </div>
    <div>
        <v-text-field class="input-year" v-model="value.y" @update:focused="onFocusYear" @keydown.enter="handleYear"
            density="compact" variant="underlined"></v-text-field>
    </div>
</template>
    <!--  https://webmound.com/use-v-model-custom-components-vue-3/ -->
<script setup>

import { ref, onMounted } from 'vue'

const props = defineProps({
    modelValue: Object,
    hint: {
        type: String,
        default: "Date of birth"
    },
    label: {
        type: String,
        default: "birthday",
    }
})

const value = ref({
    d: "",
    m: "",
    y: ""
})
const emits = defineEmits(['update:modelValue'])

onMounted(() => {
    if (props) {
        if (props.modelValue) {
            if (props.modelValue == "") {
                console.log("onMounted", props.modelValue)
            }
        }
    }
})

function onFocusDay(asd) {
    if (asd == false) {
        handleDay()
    }
}
function handleDay(val) {
    console.log("handleDay", val)
    handleAll()
}

function onFocusMonth(asd) {
    if (asd == false) {
        handleMonth()
    }
}
function handleMonth(val) {
    console.log("handleMonth", val)
    handleAll()
}

function onFocusYear(asd) {
    if (asd == false) {
        handleMonth()
    }
}
function handleYear(val) {
    console.log("handleYear", val)
    handleAll()
}

function handleAll() {
    emits('update:modelValue', value.value.d)
}

</script>

<style scoped>
.input-day {
    width: 20px !important;
    max-width: 20px;
}

.input-month {
    width: 20px !important;
    max-width: 20px;
}

.input-year {
    width: 40px !important;
    max-width: 40px;
}
</style>