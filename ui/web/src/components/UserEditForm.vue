<template>
    <v-card max-width="800">
        <v-card-title>
            <v-row no-gutters>
                <v-col class="d-flex justify-start">
                    <v-chip color="primary">User </v-chip>
                    <v-chip v-if="refs.isNew == false">{{ refs.user.username }} </v-chip>
                    <v-chip v-else>New</v-chip>
                </v-col>
                <v-col class="d-flex justify-end"><v-btn rounded @click="submit">{{ refs.submitTitle }}</v-btn></v-col>
            </v-row>
            <v-row no-gutters>
                <v-col v-if="refs.isNew" class="d-flex justify-start">
                    <v-text-field v-model="refs.user.username" label="Username"
                        hint="email/username address"></v-text-field>
                </v-col>
            </v-row>
        </v-card-title>
        <v-container>
            <v-row no-gutters>
                <v-col class="mr-0 pr-1" cols="3">
                    <v-text-field v-model="refs.user.status" label="Status" hint="active"></v-text-field>
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
import { User, EditUser } from "@/store/users.js";


const config = useConfig;

onMounted(() => {
    console.log("onMounted.", props.modelValue.user.username)
    refs.value.user = props.modelValue.user
    if (props.modelValue.user.username.length > 0) {
        refs.value.submitTitle = "Save"
        refs.value.isNew = false
    } else {
        refs.value.isNew = true
        refs.value.submitTitle = "Add"
    }
})

const props = defineProps({
    modelValue: EditUser,
})

const emits = defineEmits(['save', 'add', 'update:modelValue'])

const refs = ref({
    submitTitle: "Add",
    isNew: true,
    user: new User(),
    systemRoles: [
        "admin",
    ],

})

async function addUser(user) {
    console.log("addUser", user)
    const data = {
        data: user
    }
    const request = {
        // mode: "no-cors",
        method: "POST",
        body: JSON.stringify(data),
    };
    refs.value.loading = true;
    const url = config.BASE_URL + "/v1/auth/save";
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

async function updateUser(user) {
    console.log("updateUser", user)
    const data = {
        data: user
    }
    const request = {
        // mode: "no-cors",
        method: "POST",
        body: JSON.stringify(data),
    };
    const url = config.BASE_URL + "/v1/auth/update";
    console.log("url", url)
    refs.value.loading = true;
    await fetchAPI(url, request)
        .then((result) => {
            console.log("user.result", result)
        }).catch((err) => {
            console.log("err", err)
        }).finally(() => {
            refs.value.loading = false;
        });
}



async function submit() {
    let emit = "add"
    if (refs.value.user.username.length == 0) {
        addUser(refs.value.user)
        emit = "save"
    } else {
        updateUser(refs.value.user)
    }
    let temp = props.modelValue
    temp.person = refs.value.user
    emits('update:modelValue', temp)
    emits(emit, refs.value.user)
    console.log("submit", emit, refs.value.user)
}
</script>