<template>
    <v-main>
        <v-card outlined justify="right">
            <v-card-text>
                <v-container fill-height no-gutters class="ma-0 mp-0">
                    <v-row>
                        <v-col no-gutters cols="4" sm="6" md="4">
                            <BtnLoadData :disabled="refs.loading" @click="loadData" text="Load people"></BtnLoadData>
                        </v-col>
                        <v-col no-gutters cols="1" sm="6" md="4">
                            <BtnLoadData @click="showNewPerson" text="Add"></BtnLoadData>
                            <v-dialog v-model="refs.edit.show" max-width="900">
                                <PersonEditForm v-model="refs.edit" @add="addPerson" @save="savePerson" />
                            </v-dialog>
                            <!-- <div class="d-flex justify-center ma-4">
                                <div class="d-flex justify-space-between" style="width: 60%">
                                    <v-progress-circular></v-progress-circular>
                                </div>
                            </div>
 -->
                        </v-col>
                    </v-row>
                    <v-row>
                        <v-col no-gutters cols="4" class="">
                            <InputTextsList :list="searchTextFields.list" @click="loadData" label="Names"
                                hint="Email or names">
                            </InputTextsList>
                            <InputTextsList :list="searchPhonesFields.list" @click="loadData" label="Phones" hint="Phones">
                            </InputTextsList>
                            <InputTextsList :list="searchPINFields.list" @click="loadData" label="PIN" hint="ЕГН">
                            </InputTextsList>
                        </v-col>
                    </v-row>
                    <v-row no-gutters>
                        <v-col class="text-left">
                            <ChipsList :list="searchTextFields" @click="delInputTexts"></ChipsList>
                            <ChipsList :list="searchPhonesFields" @click="delInputPhones"></ChipsList>
                            <ChipsList :list="searchPINFields" @click="delInputPIN"></ChipsList>
                        </v-col>
                    </v-row>
                </v-container>
            </v-card-text>
            <v-data-table-server show-expand :loading="refs.loading" :items-length="refs.totalItems" :headers="refs.headers"
                :items="list.list" multi-sort class="elevation-1">
                <template v-slot:top>
                </template>
                <template v-slot:no-data>
                    <BtnLoadData v-if="!refs.loadingText" @click="loadData" text="Load people"></BtnLoadData>
                    <v-chip color="primary" veriant="text" v-else>{{ refs.loadingText }}</v-chip>
                </template>
                <template v-slot:[`header.full_name`]="{ column }">
                    <div>
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
                    <FieldBirthday :dob="item.dob"></FieldBirthday>
                </template>
                <template v-slot:[`item.id_numbers`]="{ item }">
                    <field-email-maps :emails="item.id_numbers"></field-email-maps>
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
                <template v-slot:expanded-row="{ columns, item }">
                    <tr>
                        <td :colspan="columns.length">
                            <v-card class="bg-primary-lighten-1">
                                <v-row>
                                    <v-col>
                                        <FieldBirthday :dob="item.dob" />
                                    </v-col>
                                    <v-col>
                                        <field-email-maps :emails="item.id_numbers"></field-email-maps>
                                    </v-col>
                                    <v-col>
                                        {{ item.age }}
                                    </v-col>
                                    <v-col>
                                        {{ item.gender }}
                                    </v-col>
                                </v-row>
                            </v-card>
                        </td>
                    </tr>
                </template>
            </v-data-table-server>
        </v-card> </v-main>
</template>
<script setup>

import { fetchAPI } from "@/store/auth";
import { useConfig } from "@/store/app";

import ChipsList from '@/components/ChipsList'
import BtnLoadData from '@/components/BtnLoadData';
import InputTextsList from '@/components/InputTextsList';
import FieldBirthday from '@/components/FieldBirthday';
import FieldTimeStamp from '@/components/FieldTimeStamp';
import FieldEmailMaps from '@/components/FieldEmailMaps';
import FieldPhonesMaps from '@/components/FieldPhonesMaps';
import FieldLabelsList from '@/components/FieldLabelsList';
import PersonEditForm from '@/components/PersonEditForm'

import { ref } from 'vue'
import { Person, Dob, EditPerson } from "@/store/people";

// const store = usePeopleStore()
const config = useConfig;


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
    }
    if (item.full_name.length == 0) {
        return "----"

    }
    return item.full_name
}

const searchTextFields = ref({ list: [] })
const searchPhonesFields = ref({ list: [] })
const searchPINFields = ref({ list: [] })

const refs = ref({
    edit: new EditPerson(),
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
        { title: 'Id', key: 'id', align: ' d-none' },
        { title: 'Emails', key: 'emails', align: 'end', sortable: false },
        { title: 'Names', key: 'full_name', align: 'end' },
        { title: 'Phones', key: 'phones', align: 'end' },
        { title: 'Labels', key: 'labels', align: 'end' },
        { title: 'Attributes', key: 'attrs', align: 'end' },

        { title: 'Gender', key: 'gender', align: ' d-none' },
        { title: 'IDs', key: 'id_numbers', align: ' d-none' },
        { title: 'Age', key: 'age', align: ' d-none' },
        { title: 'Birthday', key: 'dob', align: ' d-none' },

        { title: 'Created', key: 'created_at', align: 'end', sortable: false },
        { title: 'Actions', key: 'actions', sortable: false },

        { title: 'Name', key: 'name', align: ' d-none' },

    ]
})

function addPerson(person) {
    console.log("addPerson", person, refs.value.edit)
    refs.value.edit.show = false
}

function savePerson(person) {
    console.log("savePerson", person, refs.value.edit)
    refs.value.edit.show = false
}

function editItem(person) {
    console.log("editItem", person)
    refs.value.edit.person = person
    refs.value.edit.show = true
}

function showNewPerson() {
    refs.value.edit.person = new Person()
    refs.value.edit.show = true
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
    const request = {
        // mode: "no-cors",
        method: "POST",
        body: JSON.stringify(filter),
    };
    refs.value.loading = true;
    const url = config.BASE_URL + "/v1/people/list";
    console.log("url", url)
    await fetchAPI(url, request)
        .then((result) => {
            console.log("people.result", result)
            refs.value.totalItems = 0;
            list.value.list = []
            if (result.result !== null) {
                result.result.forEach(
                    (j) => {
                        let i = Object.assign(new Person(), j);
                        if (i.dob !== undefined) {
                            let dob = Object.assign(new Dob(), i.dob)
                            i.dob = dob
                        }
                        console.log(i)
                        list.value.list.push(i)
                        refs.value.totalItems++;
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