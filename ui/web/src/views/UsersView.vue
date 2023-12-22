<template>
    <v-main>
        <v-card outlined justify="right">
            <v-card-text>
                <v-container fill-height no-gutters class="ma-0 mp-0">
                    <v-row>
                        <v-col no-gutters cols="4" sm="6" md="4">
                            <BtnLoadData :disabled="refs.loading" @click="loadData" text="Load"></BtnLoadData>
                        </v-col>
                        <v-col no-gutters cols="1" sm="6" md="4">
                            Add
                        </v-col>
                    </v-row>
                    <v-row>
                        <v-col no-gutters cols="4" class="">
                            <InputTextsList :list="searchTextFields.list" @click="loadData" label="Names"
                                hint="Email or names">
                            </InputTextsList>
                        </v-col>
                    </v-row>
                </v-container>
            </v-card-text>
            <v-data-table-server show-expand :loading="refs.loading" :items-length="refs.totalItems" :headers="refs.headers"
                :items="list.list" multi-sort class="elevation-1">
                <template v-slot:top>
                </template>
                <template v-slot:no-data>
                    <BtnLoadData v-if="!refs.loadingText" @click="loadData" text="Load"></BtnLoadData>
                    <v-chip color="primary" veriant="text" v-else>{{ refs.loadingText }}</v-chip>
                </template>
                <template v-slot:[`header.email`]="{ column }">
                    <div>
                        {{ column.title }}
                    </div>
                </template>
                <template v-slot:[`item.system_roles`]="{ item }">
                    <FieldLabelsList :labels="item.system_roles"></FieldLabelsList>
                </template>

                <template v-slot:[`item.actions`]="{ item }">
                    <v-btn title="Edit" @click="editItem(item)" rounded="xl" class="me-2">
                        <v-icon size="large">
                            mdi-pencil
                        </v-icon>
                    </v-btn>
                    <v-btn title="Delete" disabled color="warning" @click="editItem(item)" rounded="xl" class="me-2">
                        <v-icon size="large">
                            mdi-delete
                        </v-icon>
                    </v-btn>
                    <!-- <v-icon size="large" color="warning" @click="deleteItem(item)">
                    </v-icon> -->
                </template>
            </v-data-table-server>
        </v-card> </v-main>
</template>
<script setup>

import { fetchAPI } from "@/store/auth";
import { useConfig } from "@/store/app";

import BtnLoadData from '@/components/BtnLoadData';
import InputTextsList from '@/components/InputTextsList';
import FieldLabelsList from '@/components/FieldLabelsList';

import { ref } from 'vue'

// const store = usePeopleStore()
const config = useConfig;


const list = ref({ list: [] })

const searchTextFields = ref({ list: [] })

const refs = ref({
    filters: {
        texts: [],
    },
    itemsPerPage: 5,
    serverItems: [],
    loadingText: "",
    loading: false,
    totalItems: 0,
    headers: [
        { title: 'Email', key: 'email', align: 'end', sortable: false },
        { title: 'Status', key: 'status', align: 'end' },
        { title: 'Roles', key: 'system_roles', align: 'end' },

        { title: 'Created', key: 'created_at', align: 'end', sortable: false },
        { title: 'Actions', key: 'actions', sortable: false },


    ]
})

async function loadData() {
    const request = {
        // mode: "no-cors",
        method: "GET",
    };
    refs.value.loading = true;
    const url = config.BASE_URL + "/v1/auth/list";
    console.log("url", url)
    await fetchAPI(url, request)
        .then((result) => {
            console.log("authlist.result", result)
            refs.value.totalItems = 0;
            list.value.list = []
            if (result.result !== null) {
                result.result.forEach(
                    (j) => {
                        console.log(j)
                        list.value.list.push(j)
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