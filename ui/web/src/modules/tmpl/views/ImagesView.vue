<template>
    <v-main>
        <v-card outlined justify="right">
            <v-card-text>
                <v-container fill-height no-gutters class="ma-0 mp-0">
                    <v-row>
                        <v-col no-gutters cols="4" sm="6" md="4">
                            <BtnLoadData @click="loadImages" text="List"></BtnLoadData>
                            <BtnLoadData @click="showUploadImage" text="add"></BtnLoadData>
                            <v-dialog v-model="refs.addImage" max-width="600">
                                <!-- <PersonEditForm v-model="refs.edit" @add="addPerson" @save="savePerson" /> -->
                                <ImageAddForm @save="showUploadImage"></ImageAddForm>
                            </v-dialog>
                            <!-- <BtnLoadData :disabled="refs.loading" @click="loadData" text="Load people"></BtnLoadData> -->
                        </v-col>
                    </v-row>
                </v-container>
            </v-card-text>
            <v-data-table-server :loading="refs.loading" :items-length="refs.totalItems" :headers="refs.headers"
                :items="refs.data" multi-sort class="elevation-1">
                <template v-slot:[`item.src`]="{ item }">
                    {{ item.src }}
                    <img :src="item.src" :alt="item.name" />
                    <!-- <img src="http://localhost:8010/tmpl/image/glass-mug-variant2.png" :alt="item.name" /> -->
                </template>
            </v-data-table-server>
        </v-card>
    </v-main>


</template>

<script setup>
import { ref } from 'vue';
import ImageAddForm from '../components/ImageAddForm.vue';
import BtnLoadData from '@/components/BtnLoadData.vue';
import { fetchAPI } from "@/store/auth";
import { useConfig } from "@/store/app";

const config = useConfig;

const refs = ref({
    data: [],
    totalItems: 0,
    loading: false,
    addImage: false,
    headers: [
        { title: 'Id', key: 'id', align: '' },
        { title: 'Preview', key: 'src', align: '', sortable: false },
        { title: 'Name', key: 'name', align: 'end', sortable: false }
    ],
    file: {
        fileName: null,
        fileForm: null,
        ready: false,
    }
})
function showUploadImage() {
    refs.value.addImage = !refs.value.addImage
    console.log("refs.value.addImage", refs.value.addImage)
}

async function loadImages() {
    const request = {
        // mode: "no-cors",
        method: "GET",
    };
    refs.value.loading = true;
    const url = config.BASE_URL + "/tmpl/images";
    console.log("url", url)
    await fetchAPI(url, request)
        .then((result) => {
            console.log("result    :", result)
            // refs.value.totalItems = 0;
            // refs.value.data = []
            if (result.result.list !== null) {
                console.log("result.list: ", result.result.list)
                result.result.list.forEach(
                    (row) => {
                        console.log("row: ", row)
                        console.log("result j: ", row["Id"])
                        let id = row['Id']
                        let fileInfo = {
                            id: id,
                            name: row['Name'],
                            src: `http://localhost:8010/tmpl/image/${id}/resized`
                        }
                        refs.value.data.push(fileInfo)
                        refs.value.totalItems++
                    }
                )
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