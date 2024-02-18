<template>
    <v-card max-width="800">
        <v-card-title>
            <v-row no-gutters>
                <v-col class="d-flex justify-start">
                    <v-chip color="primary">Person </v-chip>
                    <v-chip v-if="refs.person.id.length > 0" color="">{{ refs.person.id.substring(0, 10) }} </v-chip>
                </v-col>
                <v-col class="d-flex justify-end"><v-btn rounded @click="submit">{{ refs.submitTitle
                }}</v-btn></v-col>
            </v-row>
        </v-card-title>
        <v-container>
            <v-row no-gutters>
                <v-col class="mr-0 pr-1" cols="3">
                    <v-text-field v-model="refs.person.name" label="Name" hint="Peter"></v-text-field>
                </v-col>
                <v-col class="ml-0 pl-1" cols="">
                    <v-text-field v-model="refs.person.full_name" label="Full names" hint="Varban Krushev"></v-text-field>
                </v-col>
            </v-row>
            <v-row no-gutters>
                <v-col>
                    <InputTypeValue v-model="refs.person.id_numbers" label="Number" :typeItems="refs.idTypes">
                    </InputTypeValue>
                </v-col>
                <v-col>
                    <v-chip size="x-small" rounded v-for="(val, index) in refs.person.id_numbers" :key="index">
                        {{ index }}
                        {{ val }}
                    </v-chip>
                </v-col>
            </v-row>
            <v-row no-gutters>
                <v-col>
                    <InputTypeValue :rules="phoneRules" v-model="refs.person.phones" label="Phone"
                        :typeItems="refs.phoneTypes">
                    </InputTypeValue>
                </v-col>
                <v-col>
                    <v-chip size="x-small" rounded v-for="(val, index) in refs.person.phones" :key="index">
                        {{ index }}
                        {{ val }}
                    </v-chip>
                </v-col>
            </v-row>
            <v-row no-gutters>
                <v-col>
                    <InputTypeValue :rules="emailRules" v-model="refs.person.emails" label="Email"
                        :typeItems="refs.emailTypes">
                    </InputTypeValue>
                </v-col>
                <v-col>
                    <v-chip size="x-small" rounded v-for="(val, index) in refs.person.emails" :key="index">
                        {{ index }}
                        {{ val }}
                    </v-chip>
                </v-col>
            </v-row>
            <v-row>
                <v-col>
                    <InputBirthday v-model="refs.person.dob" />
                </v-col>
                <v-col>
                    <InputGender v-model="refs.person.gender" />
                </v-col>
            </v-row>
        </v-container>
    </v-card>
</template>


<script setup>
import { ref, onMounted } from 'vue';
// import { useDisplay } from 'vuetify'
import { fetchAPI } from "@/store/auth";
import { useConfig } from "@/store/app";
import { Person, EditPerson } from "@/store/people.js";

import InputTypeValue from '@/components/InputTypeValue.vue';
import InputGender from '@/components/InputGender.vue';
import InputBirthday from '@/components/InputBirthday.vue';

const config = useConfig;

onMounted(() => {
    console.log("onMounted.", props.modelValue.person.id)
    refs.value.person = props.modelValue.person
    if (props.modelValue.person.id.length > 0) {
        refs.value.submitTitle = "Save"
    } else {
        refs.value.submitTitle = "Add"
    }
})

const props = defineProps({
    modelValue: EditPerson,
})

const emits = defineEmits(['save', 'add', 'update:modelValue'])

const refs = ref({
    submitTitle: "Add",
    person: new Person(),
    phoneTypes: [
        "main",
        "home",
    ],
    idTypes: [
        "EGN",
        "Passport",
    ],
    emailTypes: [
        "main",
        "home",
    ]
})

const phoneRules = [
    value => !!value || 'Required.',
]


const emailRules = [
    value => !!value || 'Required.',
    value => {
        const pattern = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
        return pattern.test(value) || 'Invalid e-mail.'
    },
]

async function addPerson(person) {
    console.log("savePerson", person)
    const data = {
        data: person
    }
    const request = {
        // mode: "no-cors",
        method: "POST",
        body: JSON.stringify(data),
    };
    refs.value.loading = true;
    const url = config.BASE_URL + "/v1/people/save";
    console.log("url", url)
    await fetchAPI(url, request)
        .then((result) => {
            console.log("people.result", result)
        }).catch((err) => {
            console.log("err", err)
        }).finally(() => {
            refs.value.loading = false;
        });
}

async function updatePerson(person) {
    console.log("updatePerson", person)
    const data = {
        data: person
    }
    const request = {
        // mode: "no-cors",
        method: "POST",
        body: JSON.stringify(data),
    };
    const url = config.BASE_URL + "/v1/people/update";
    console.log("url", url)
    refs.value.loading = true;
    await fetchAPI(url, request)
        .then((result) => {
            console.log("people.result", result)
        }).catch((err) => {
            console.log("err", err)
        }).finally(() => {
            refs.value.loading = false;
        });
}



async function submit() {
    let emit = "add"
    if (refs.value.person.id.length == 0) {
        addPerson(refs.value.person)
        emit = "save"
    } else {
        updatePerson(refs.value.person)
    }
    let temp = props.modelValue
    temp.person = refs.value.person
    emits('update:modelValue', temp)
    emits(emit, refs.value.person)
    console.log("submit", emit, refs.value.person)
}
</script>