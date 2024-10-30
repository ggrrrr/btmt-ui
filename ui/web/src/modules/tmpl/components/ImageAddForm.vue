<template>
    <v-card max-width="400">
        <v-card-title>
            <v-row no-gutters>
                <v-col class="d-flex justify-start">
                    <!-- <v-chip color="primary">Upload image </v-chip> -->
                    <v-label>{{ refs.file.fileName }}</v-label>

                </v-col>
                <v-col class="d-flex justify-end">
                    <BtnLoadData :disabled="!refs.file.ready" @click="inputFileUpload" text="Upload"></BtnLoadData>
                </v-col>
            </v-row>
        </v-card-title>
        <v-container>
            <v-form>
                <v-row no-gutters>
                    <v-col class="mr-0 pr-1">
                        <v-file-input density="compact" color="primary" show-size @change="formChange"
                            v-model="refs.file.fileForm" type='file' name="image" @update="inputFileUpload">
                        </v-file-input>
                    </v-col>
                </v-row>
            </v-form>
        </v-container>
    </v-card>
</template>

<script setup>
import { ref } from 'vue'
import BtnLoadData from '@/components/BtnLoadData.vue';

import { fetchAPI } from "@/store/auth";
import { useConfig } from "@/store/app";

const emits = defineEmits(['save'])


const config = useConfig;

const refs = ref({
    file: {
        fileName: null,
        fileForm: null,
        ready: false,
    }
})

function formChange() {
    console.log("formChange", refs.value.file.fileForm);
    if (!refs.value.file.fileForm) {
        console.log("fileForm empty")
        refs.value.file.ready = false
        return
    }

    const regexpSize = /[^0-9a-z\-.]/gi;


    if (refs.value.file.fileForm.name) {
        refs.value.file.ready = true
        const match = refs.value.file.fileForm.name.replace(regexpSize, '');
        refs.value.file.fileName = match
    }

}

async function inputFileUpload() {
    console.log('fileForm', refs.value.file.fileForm)

    if (refs.value.file.fileForm == null) {
        console.log("no files")
        return
    }

    let formData = new FormData();
    formData.append("file", refs.value.file.fileFormта);
    // formData.append("file", refs.value.file.fileName);
    let request = {
        method: 'POST',
        body: formData,
    };
    const url = config.BASE_URL + "/tmpl/image";
    await fetchAPI(url, request)
        .then((result) => {
            console.log("file.result", result)
            emits('save')
        }).finally(() => {
            console.log("file.result", "finally")
        });

}

</script>