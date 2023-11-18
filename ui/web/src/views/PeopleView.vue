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
            <v-data-table :headers="headers.headers" :items="list.list" multi-sort class="elevation-1">
            </v-data-table>
        </v-card> </v-main>
</template>
<script setup>

import { fetchAPI } from "@/store/auth";
import { ref } from 'vue'

const list = ref({ list: [] })

const headers = ref({
    headers: [
        {
            text: "Names",
            align: "start",
            sortable: true,
            value: "Name",
        },
        { text: "Email", value: "Email", sortable: true },
        { text: "Labels", value: "Labels" },
        { text: "Phones", value: "Phones" },
        { text: "PIN", value: "Pin" },
        { text: "Date of birth", value: "Dob" },
        { text: "Age", value: "Age" },
        { text: "Actions", value: "actions", sortable: false },
    ]
})


function loadData() {
    const requestOptions = {
        // mode: "no-cors",
        method: "POST",
        body: JSON.stringify({}),
    };

    fetchAPI("http://10.1.1.156:8000/rest/v1/people/list", requestOptions)
        .then((result) => {
            console.log("asdasd", result)
            result.result.forEach(
                (i) => {
                    console.log(i)
                    list.value.list.push(i)
                }
            )
        });
}

</script>