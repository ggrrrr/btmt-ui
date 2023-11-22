<template>
    <v-main>
        <v-card outlined justify="right">
            <v-card-text>
                <v-container fill-height no-gutters class="ma-0 mp-0">
                    <v-row>
                        <v-col no-gutters cols="4" sm="6" md="4">
                            <BtnLoadData @click="loadData" text="Load people"></BtnLoadData>
                        </v-col>
                        <v-col no-gutters cols="1" sm="6" md="4">
                            <BtnLoadData @click="store.showEdit = !store.showEdit" text="Add"></BtnLoadData>
                            <PersonDialog v-model="store.showEdit" />
                        </v-col>
                    </v-row>
                    <v-row justify="start">
                        <v-col no-gutters cols="4" class="justify-center">
                            <InputTextsList :list="searchTextFields.list" @click="loadData" label="Names"
                                hint="Email or names">
                            </InputTextsList>
                            <InputTextsList :list="searchPhonesFields.list" @click="loadData" label="Phones" hint="Phones">
                            </InputTextsList>
                            <InputTextsList :list="searchPINFields.list" @click="loadData" label="PIN" hint="ЕГН">
                            </InputTextsList>
                        </v-col>
                    </v-row>
                    <v-row justify="start" no-gutters>
                        <v-col class="text-left">
                            <ChipsList :list="searchTextFields" @click="delInputTexts"></ChipsList>
                            <ChipsList :list="searchPhonesFields" @click="delInputPhones"></ChipsList>
                            <ChipsList :list="searchPINFields" @click="delInputPIN"></ChipsList>
                        </v-col>
                    </v-row>
                </v-container>
            </v-card-text>
            <v-data-table-server :items-length="data.totalItems" :headers="data.headers" :items="list.list" multi-sort
                class="elevation-1">
                <template v-slot:top>
                </template>
                <template v-slot:no-data>
                    <BtnLoadData v-if="!data.loadingText" @click="loadData" text="Load people"></BtnLoadData>
                    <v-chip color="primary" veriant="text" v-else>{{ data.loadingText }}</v-chip>
                </template>
                <template v-slot:[`header.full_name`]="{ column }">
                    <div class="justify">
                        {{ column.title }}
                    </div>
                </template>
                <template v-slot:[`item.full_name`]="{ item }">
                    {{ formatFullName(item) }}
                </template>
                <template v-slot:[`item.created_at`]="{ item }">
                    <field-time-stamp :timeStamp="item.created_at"></field-time-stamp>
                </template>
                <template v-slot:[`item.dob`]="{ item }">
                    <FieldDOB :dob="item.dob"></FieldDOB>
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
import ChipsList from '@/components/ChipsList'
import BtnLoadData from '@/components/BtnLoadData';
import InputTextsList from '@/components/InputTextsList';
import FieldDOB from '@/components/FieldDOB';
import FieldTimeStamp from '@/components/FieldTimeStamp';
import FieldEmailMaps from '@/components/FieldEmailMaps';
import FieldPhonesMaps from '@/components/FieldPhonesMaps';
import FieldLabelsList from '@/components/FieldLabelsList';
import PersonDialog from '@/components/PersonDialog'

import { ref } from 'vue'
import { usePeopleStore } from "@/store/people";

const store = usePeopleStore()

const list = ref({ list: [] })

function delInputTexts(index) {
    console.log("delInputTexts", index)
    searchTextFields.value.list = pop(searchTextFields.value.list, index)
    loadData()
}

function delInputPhones(index) {
    console.log("delInputPhones", index)
    searchPhonesFields.value.list = pop(searchPhonesFields.value.list, index)
    loadData()
}
function delInputPIN(index) {
    console.log("delInputPhones", index)
    searchPINFields.value.list = pop(searchPINFields.value.list, index)
    loadData()
}

function pop(list, index) {
    let out = []
    list.forEach((v, i) => {
        if (i != index) {
            out.push(v)
        }
    })
    return out
}

function formatFullName(item) {
    if (item.name && item.full_name) {
        if (item.name.length > 0) {
            return `(${item.name}) ${item.full_name}`
        }
        return item.full_name
    }
}

const searchTextFields = ref({ list: [] })
const searchPhonesFields = ref({ list: [] })
const searchPINFields = ref({ list: [] })

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
        { title: 'Birthday', key: 'dob', align: 'end' },
        { title: 'Gender', key: 'gender', align: 'end' },
        { title: 'Age', key: 'age', align: 'end' },
        { title: 'Emails', key: 'emails', align: 'end' },
        { title: 'Names', key: 'full_name', align: 'end' },
        { title: 'Phones', key: 'phones', align: 'end' },
        { title: 'Labels', key: 'labels', align: 'end' },
        { title: 'Created', key: 'created_at', align: 'end' },
    ]
})

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

async function loadData() {
    let filter = {
        filters: {
        }
    }
    if (searchTextFields.value.list.length > 0) {
        filter.filters.texts = searchTextFields.value
    }
    if (searchPhonesFields.value.list.length > 0) {
        filter.filters.phones = searchPhonesFields.value
    }
    if (searchPINFields.value.list.length > 0) {
        filter.filters.pins = searchPINFields.value
    }
    console.log("filter", filter)
    const requestOptions = {
        // mode: "no-cors",
        method: "POST",
        body: JSON.stringify(filter),
    };
    data.value.loading = true;
    await sleep(2000)
    await fetchAPI("http://10.1.1.156:8000/rest/v1/people/list", requestOptions)
        .then((result) => {
            console.log("people.result", result)
            data.value.totalItems = 0;
            list.value.list = []
            result.result.forEach(
                (i) => {
                    console.log(i)
                    if (i.name === undefined) {
                        i.name = ""
                    }
                    if (i.full_name === undefined) {
                        i.full_name = ""
                    }
                    if (i.labels === undefined) {
                        i.labels = []
                    }
                    if (i.phones === undefined) {
                        i.phones = {}
                    }
                    if (i.emails === undefined) {
                        i.emails = {}
                    }
                    console.log("row", i)

                    list.value.list.push(i)
                    data.value.totalItems++;
                }
            )
            if (data.value.totalItems == 0) {
                data.value.loadingText = "no data, check filters"
            }
        }).finally(() => {
            data.value.loading = false;
        });
}

</script>