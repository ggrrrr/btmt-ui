<template>
    <v-container no-gutters>
        <v-card outlined no-gutters justify="right">
            <v-card-title>
                <v-row>
                    <v-col no-gutters cols="4" sm="6" md="4">
                        <BtnLoadData @click="listRequest" text="List"></BtnLoadData>
                        <BtnLoadData to="/templates/manage/edit/new" text="Add"></BtnLoadData>
                    </v-col>

                </v-row>
            </v-card-title>
            <v-data-table-server :loading="refs.loading" :items-length="refs.totalItems" :headers="refs.headers"
                :items="refs.data" multi-sort class="elevation-1">
                <template v-slot:[`item.id`]="{ item }">
                    <BtnLoadData :to="item.editUrl" :text="item.id.substring(0, 10)"> {{ item.id.substring(0, 10) }}
                    </BtnLoadData>
                </template>
                <template v-slot:[`item.name`]="{ item }">
                    {{ item.name }}
                </template>
                <template v-slot:[`item.body`]="{ item }">
                    {{ item.body.substring(0, 20) }}...
                </template>
                <template v-slot:[`item.labels`]="{ item }">
                    <FieldLabelsList :labels="item.labels"></FieldLabelsList>
                </template>

                <template v-slot:[`item.src`]="{ item }">
                    <v-img :src="item.src" :alt="item.name" height="120" aspect-ratio="1">
                        <template v-slot:placeholder>
                            <div class="d-flex align-center justify-center fill-height">
                                <v-progress-circular color="grey-lighten-5" indeterminate></v-progress-circular>
                            </div>
                        </template>
                    </v-img>
                    <!-- <img src="http://localhost:8010/tmpl/image/glass-mug-variant2.png" :alt="item.name" /> -->
                </template>
                <template v-slot:[`item.created_at`]="{ item }">
                    <field-time-stamp :timeStamp="item.created_at"></field-time-stamp>
                </template>
                <template v-slot:[`item.actions`]="{ item }">
                    <v-btn title="Edit" :to="item.editUrl" rounded="xl" class="me-2">
                        <v-icon size="large">
                            mdi-pencil
                        </v-icon>
                    </v-btn>
                </template>

            </v-data-table-server>
        </v-card>
    </v-container>
</template>

<script setup>
import BtnLoadData from '@/components/BtnLoadData.vue';
import FieldTimeStamp from '@/components/FieldTimeStamp';
import FieldLabelsList from '@/components/FieldLabelsList.vue';

import { ref } from 'vue';
import { fetchAPI } from "@/store/auth";

import { useConfig } from "@/store/app";
import { Template } from '../svc/tmpl';
const config = useConfig;

const refs = ref({
    data: [],
    totalItems: 0,
    loading: false,
    headers: [
        { title: 'Id', key: 'id', align: '' },
        { title: 'Name', key: 'name', align: 'left', sortable: false },
        { title: 'Type', key: 'content_type', align: 'left', sortable: false },
        { title: 'Labels', key: 'labels', align: 'end', sortable: false },
        { title: 'Preview', key: 'body', align: '', sortable: false },
        { title: "Created", key: 'created_at', align: "right", sortable: false },
        { title: '', key: 'actions', sortable: false },
    ],
})

async function listRequest() {
    const reqest = {
        items: {},
        body: refs.value.row,
    }
    const request = {
        // mode: "no-cors",
        method: "POST",
        body: JSON.stringify(reqest)
    };

    const url = config.BASE_URL + "/v1/tmpl/manage/list";
    console.log("url", url)
    await fetchAPI(url, request)
        .then((result) => {
            // refs.value.totalItems = 0;
            // refs.value.data = []
            if (result.result.list !== null) {
                console.log("result: ", result.result)
                result.result.list.forEach(
                    (row) => {
                        console.log("result: ", row)
                        // const tmpl = {
                        //     id: row.Id,
                        //     name: row.Name,
                        //     body: row.Body,
                        //     labels: row.Labels,
                        //     created_at: row.CreatedAt,
                        //     editUrl: `/templates/manage/edit/${row.Id}`
                        // }
                        let tmpl = new Template(row);
                        refs.value.data.push(tmpl)

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