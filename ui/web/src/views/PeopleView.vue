<template>
    <v-main>
        <v-card outlined justify="right">
            <v-card-title>
                <v-col cols="12" sm="6" md="4">
                    <BtnLoadData @click="loadData" text="Load people"></BtnLoadData>
                </v-col>
            </v-card-title>
            <v-data-table-server :no-data-text="data.loadingText" :items-length="data.totalItems" :headers="data.headers"
                :items="list.list" :loading="data.loading" multi-sort class="elevation-1">
                <template v-slot:no-data>
                    <BtnLoadData v-if="!data.loadingText" @click="loadData" text="Load people"></BtnLoadData>
                    <div v-else>{{ data.loadingText }}</div>
                </template>
                <template v-slot:[`header.full_name`]="{ column }">
                    <div class="justify">

                        <v-text-field density="compact" v-model="searchTextField" @keydown.enter="addSearchField"
                            @click:append-inner="addSearchField" append-inner-icon="mdi-plus" :label="column.title"
                            type="text" hint="Email or name"></v-text-field>
                        <v-chip size="small" append-icon="mdi-close" v-for="(name, index) in searchTextFields.list"
                            :key="index" @click="delSearchField(index)"> {{ name }}</v-chip>
                    </div>

                </template>
                <template v-slot:[`item.full_name`]="{ item }">
                    {{ formatFullName(item) }}
                </template>
                <template v-slot:[`item.created_at`]="{ item }">
                    <field-time-stamp :timeStamp="item.created_at"></field-time-stamp>
                </template>
                <template v-slot:[`item.emails`]="{ item }">
                    <field-email-maps :emails="item.emails"></field-email-maps>
                </template>
                <template v-slot:[`item.phones`]="{ item }">
                    <FieldPhonesMaps :phones="item.phones"></FieldPhonesMaps>
                </template>
                <template v-slot:[`item.labels`]="{ item }">
                    <FieldLabelsList :labels="item.labels"></FieldLabelsList>
                </template>
            </v-data-table-server>
        </v-card> </v-main>
</template>
<script setup>

import { fetchAPI } from "@/store/auth";
import BtnLoadData from '@/components/BtnLoadData';
import FieldTimeStamp from '@/components/FieldTimeStamp';
import FieldEmailMaps from '@/components/FieldEmailMaps';
import FieldPhonesMaps from '@/components/FieldPhonesMaps';
import FieldLabelsList from '@/components/FieldLabelsList';

import { ref } from 'vue'

const list = ref({ list: [] })

function formatFullName(item) {
    if (item.name.length > 0) {
        return `(${item.name}) ${item.full_name}`
    }
    return item.full_name
}

const searchTextField = ref("")

const searchTextFields = ref({ list: [] })

function delSearchField(index) {
    let out = []
    searchTextFields.value.list.forEach((v, i) => {
        if (i != index) {
            out.push(v)
        }
    })
    searchTextFields.value.list = out
    console.log("searchTextFields", index, searchTextFields.value)
    // searchTextFields.value.pop(index)
    loadData()
}

function addSearchField() {
    console.log("addSearchField", searchTextField.value, searchTextField.value.length)
    if (searchTextField.value.length > 1) {
        let ok = true
        searchTextFields.value.list.forEach((v, index) => {
            console.log("addSearchField", v, index, searchTextField.value)
            if (v == searchTextField.value) {
                ok = false
            }
        })
        if (ok) {
            searchTextFields.value.list.push(searchTextField.value)
        }
        // searchTextFields.value.list.push(searchTextField.value)
        console.log("searchTextFields", searchTextFields.value)
        searchTextField.value = ""
        loadData()
    }
}
const data = ref({
    filters: {
        phones: [],
        texts: [],
    },
    itemsPerPage: 5,
    serverItems: [],
    loadingText: "",
    loading: false,
    totalItems: 0,
    headers: [
        // { title: 'Id', key: 'id', align: 'end' },
        {
            title: 'Name',
            align: ' d-none',
            key: 'name',
        },
        { title: 'PIN', key: 'pin', align: 'end' },
        { title: 'Emails', key: 'emails', align: 'end' },
        { title: 'Names', key: 'full_name', align: 'end' },
        { title: 'Phones', key: 'phones', align: 'end' },
        { title: 'Labels', key: 'labels', align: 'end' },
        { title: 'Created', key: 'created_at', align: 'end' },
    ]
})

function loadData() {
    let filter = {
        filters: {
        }
    }
    if (searchTextFields.value.list.length > 0) {
        filter.filters.texts = searchTextFields.value
    }
    console.log("filter", filter)
    const requestOptions = {
        // mode: "no-cors",
        method: "POST",
        body: JSON.stringify(filter),
    };
    data.value.loading = true;

    fetchAPI("http://10.1.1.156:8000/rest/v1/people/list", requestOptions)
        .then((result) => {
            console.log("people.result", result)
            data.value.totalItems = 0;
            list.value.list = []
            result.result.forEach(
                (i) => {
                    console.log(i)
                    list.value.list.push(i)
                    data.value.totalItems++;
                }
            )
            if (data.value.totalItems == 0) {
                data.value.loadingText = "no data, please reduce filter"
            }
        }).finally(() => {
            data.value.loading = false;
        });
}

</script>