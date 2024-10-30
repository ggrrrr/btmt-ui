<template>
    <v-main>
        <v-container>
            <v-row no-gutters>
                <v-col>
                    <v-textarea label="Row html" v-model="refs.row" variant="outlined"></v-textarea>
                </v-col>
                <v-col>
                    <span v-html="refs.render"></span>
                </v-col>
            </v-row>
            <v-row>
                <BtnLoadData @click="renderRequest" text="Render"></BtnLoadData>
            </v-row>
        </v-container>

    </v-main>
</template>

<script setup>
import BtnLoadData from '@/components/BtnLoadData.vue';

import { ref } from 'vue';
import { fetchAPI } from "@/store/auth";

import { useConfig } from "@/store/app";
const config = useConfig;

const refs = ref({
    row: `
<html><body>
<p>{{ .UserInfo.User }}</p>
<p>{{ .UserInfo.Device.DeviceInfo }}</p>
<p>{{ .UserInfo.ID }}</p>
{{ renderImg "http://localhost:8010/tmpl/image/IMG4945.JPG:1/resized" }}
</body></html>
`,
    render: ``,
})


async function renderRequest() {
    const reqest = {
        items: {},
        body: refs.value.row,
    }
    const request = {
        // mode: "no-cors",
        method: "POST",
        body: JSON.stringify(reqest)
    };

    const url = config.BASE_URL + "/tmpl/render";
    console.log("url", url)
    await fetchAPI(url, request)
        .then((result) => {
            console.log("result    :", result)
            // refs.value.totalItems = 0;
            // refs.value.data = []
            if (result.result.list !== null) {
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

</script>