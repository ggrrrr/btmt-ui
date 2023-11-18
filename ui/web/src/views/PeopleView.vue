<template>
    <v-main>
        <v-card outlined justify="right">
            <v-card-title>
                Poeple
                <v-spacer></v-spacer>
                <v-card-actions>
                    <v-btn class="mr-2" @click="loadData">refresh</v-btn>
                </v-card-actions>
            </v-card-title>
            <v-card-title>
                <v-container>
                    <v-row mx-auto>
                        <v-col class="mx-auto"> aaa</v-col>
                        <v-col cols="4">bbbb</v-col>
                    </v-row>
                </v-container>
            </v-card-title>
            <v-data-table-server :headers="data.headers" :items="list.list" :loading="data.loading" multi-sort
                class="elevation-1">
            </v-data-table-server>
        </v-card> </v-main>
</template>
<script setup>

import { fetchAPI } from "@/store/auth";
import { ref } from 'vue'

const list = ref({ list: [] })

const data = ref({
    itemsPerPage: 5,
    serverItems: [],
    loading: true,
    totalItems: 0,
    headers: [
        // { title: 'Id', key: 'id', align: 'end' },
        {
            title: 'Name',
            align: 'start',
            sortable: false,
            key: 'name',
        },
        { title: 'PIN', key: 'pin', align: 'end' },
        { title: 'Names', key: 'full_name', align: 'end' },
        { title: 'Phones', key: 'phones', align: 'end' },
        { title: 'Labels', key: 'labels', align: 'end' },
    ]
})

function loadData() {
    const requestOptions = {
        // mode: "no-cors",
        method: "POST",
        body: JSON.stringify({}),
    };
    data.value.loading = true;

    fetchAPI("http://10.1.1.156:8000/rest/v1/people/list", requestOptions)
        .then((result) => {
            console.log("asdasd", result)
            result.result.forEach(
                (i) => {
                    console.log(i)
                    list.value.list.push(i)
                }
            )
        }).finally(() => {
            data.value.loading = false;
        });
}

</script>