<template>
    <v-container no-gutters>
        <v-row>
            <v-col no-gutters cols="6" sm="6" md="6">
                <BtnLoadData @click="renderRequest" text="Render"></BtnLoadData>
                <BtnLoadData @click="saveRequest" text="Save"></BtnLoadData>
                <BtnLoadData to="/templates/manage" text="Cancel"></BtnLoadData>
                <v-text-field v-model="refs.tmpl.name" label="Name"></v-text-field>
                <v-text-field v-model="refs.tmpl.content_type" label="Type"></v-text-field>
                <v-textarea label="New html" v-model="refs.tmpl.body" variant="outlined"></v-textarea>
            </v-col>
            <v-col no-gutters cols="6" sm="6" md="4">
                <v-container no-gutters>
                    <span v-html="refs.render"></span>
                </v-container>
            </v-col>
        </v-row>
    </v-container>
</template>

<script setup>
import BtnLoadData from '@/components/BtnLoadData.vue';

import { ref } from 'vue';
import { useRoute } from 'vue-router'

import { fetchAPI } from "@/store/auth";
import { useConfig } from "@/store/app";
import { Template } from '../svc/tmpl';

const config = useConfig;
const route = useRoute();
if (route.params.id === undefined) {
    console.log("bad call")
}

if (route.params.id === "new") {
    console.log("new template")
} else {
    console.log("loading template")
    getRequest(route.params.id)
}


const refs = ref({
    tmpl: new Template(),
    render: ``,
})

async function getRequest(id) {
    const request = {
        method: "GET",
    };

    const url = config.BASE_URL + "/v1/tmpl/manage/" + id;
    console.log("url", url)
    await fetchAPI(url, request)
        .then((result) => {
            console.log("get /v1/tmpl/manage:", result)
            refs.value.tmpl = new Template(result.result)
        }).finally(() => {
            refs.value.loading = false;
        });
}


async function renderRequest() {
    const reqest = {
        items: {},
        body: refs.value.tmpl.body,
    }
    const request = {
        method: "POST",
        body: JSON.stringify(reqest)
    };

    const url = config.BASE_URL + "/v1/tmpl/manage/render";
    console.log("url", url)
    await fetchAPI(url, request)
        .then((result) => {
            console.log("result:", result)
            if (result.result.render !== null) {
                console.log("result: ", result.result.payload)
                refs.value.render = result.result.payload
            } else {
                refs.value.loadingText = "no data, check filters11"
            }
            if (refs.value.totalItems == 0) {
                refs.value.loadingText = "no data, check filters"
            }
        }).finally(() => {
            refs.value.loading = false;
        });
}

async function saveRequest() {
    const tmpl = refs.value.tmpl
    tmpl.created_at = null;
    tmpl.updated_at = null;
    const request = {
        method: "POST",
        body: JSON.stringify(tmpl)
    };
    console.log("save request:", request)

    const url = config.BASE_URL + "/v1/tmpl/manage/save";
    console.log("url", url)
    await fetchAPI(url, request)
        .then((result) => {
            console.log("saveResult:", result)
        }).finally(() => {
            refs.value.loading = false;
        });
}

</script>