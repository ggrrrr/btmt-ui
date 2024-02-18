<template>
    <v-text-field density="compact" v-model="searchTextField" @keydown.enter="addSearchField" :label="props.label"
        type="text" :hint="props.hint">
        <template v-slot:append-inner>
            <v-icon @click="addSearchField" color="primary">mdi-plus</v-icon>
        </template>
    </v-text-field>
</template>

<script setup>
import { ref } from 'vue'
let emits = defineEmits(['click'])

const searchTextField = ref("")
let props = defineProps({
    list: {
        type: Array,
        required: true
    },
    label: {
        type: String,
        required: true
    },
    hint: {
        type: String,
        required: true
    }
})
function addSearchField() {
    let list = props.list
    console.log("addSearchField", list, searchTextField.value, searchTextField.value.length)
    if (searchTextField.value.length > 1) {
        let ok = true
        list.forEach((v, index) => {
            console.log("addSearchField", v, index, searchTextField.value)
            if (v == searchTextField.value) {
                ok = false
            }
        })
        if (ok) {
            // props['list'].push(searchTextField.value)
            list.push(searchTextField.value)
            searchTextField.value = ""
            emits('click')
        }
    }
}

</script>


