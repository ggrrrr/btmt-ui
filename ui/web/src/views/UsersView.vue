<template>
    <v-main>
        <v-card outlined justify="right">
            <v-card-text>
                <v-container fill-height no-gutters class="ma-0 mp-0">
                    <v-row>
                        <v-col no-gutters cols="4" sm="6" md="4">
                            <BtnLoadData :disabled="refs.loading" @click="loadData" text="Load Users"></BtnLoadData>
                        </v-col>
                        <v-col no-gutters cols="1" sm="6" md="4">
                            <BtnLoadData @click="showNewUser" text="Add"></BtnLoadData>
                            <v-dialog v-model="refs.edit.show" max-width="900">
                                <UserEditForm v-model="refs.edit" @add="addUser" @save="saveUser" />
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
                        </v-col>
                    </v-row>
                </v-container>
            </v-card-text>
            <v-data-table-server :loading="refs.loading" :items-length="refs.totalItems" :headers="refs.headers"
                :items="list.list" multi-sort class="elevation-1">
                <template v-slot:top>
                </template>
                <template v-slot:no-data>
                    <BtnLoadData v-if="!refs.loadingText" @click="loadData" text="Load Users"></BtnLoadData>
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
                <template v-slot:[`item.created_at`]="{ item }">
                    <field-time-stamp :timeStamp="item.created_at"></field-time-stamp>
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
import { User, EditUser } from "@/store/users";

import BtnLoadData from '@/components/BtnLoadData';
import InputTextsList from '@/components/InputTextsList';
import FieldLabelsList from '@/components/FieldLabelsList';
import UserEditForm from '@/components/UserEditForm'
import FieldTimeStamp from '@/components/FieldTimeStamp';

import { ref } from 'vue'

// const store = useUsersStore
const config = useConfig;


const list = ref({ list: [] })

const searchTextFields = ref({ list: [] })

const refs = ref({
    edit: new EditUser(),
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

function editItem(user) {
    console.log("editItem", user)
    refs.value.edit.user = user
    refs.value.edit.show = true
    refs.value.edit.isNew = false
}

function showNewUser() {
    refs.value.edit.user = new User()
    refs.value.edit.show = true
    refs.value.edit.isNew = true
}

function addUser() {
}

function saveUser() {
}


async function loadData() {
    const request = {
        // mode: "no-cors",
        method: "GET",
    };
    refs.value.loading = true;
    const url = config.BASE_URL + "/v1/auth/user/list";
    console.log("url", url)
    await fetchAPI(url, request)
        .then((result) => {
            console.log("authlist.result", result)
            refs.value.totalItems = 0;
            list.value.list = []
            if (result.result !== null) {
                result.result.forEach(
                    (j) => {
                        let i = Object.assign(new User(), j);
                        console.log(i)
                        list.value.list.push(i)
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