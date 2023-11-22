<template>
    <v-dialog v-model="show" max-width="900">
        <div class="d-flex flex-wrap">
            <div class="flex-1-0 ma-1 pa-1">
                <v-card max-width="800">
                    <v-card-title>
                        <v-row no-gutters>
                            <v-col class="d-flex justify-start"><v-chip color="primary">Person</v-chip></v-col>
                            <v-col class="d-flex justify-end"><v-btn rounded>Add</v-btn></v-col>
                        </v-row>
                    </v-card-title>
                    <v-container>
                        <v-row no-gutters>
                            <v-col class="mr-0 pr-1" cols="3">
                                <v-text-field label="PIN" hint="ЕГН"></v-text-field>
                            </v-col>
                            <v-col class="mr-0 pr-1" cols="3">
                                <v-text-field label="Name" hint="Joro"></v-text-field>
                            </v-col>
                            <v-col class="ml-0 pl-1" cols="">
                                <v-text-field label="Full names" hint="Varban Krushev"></v-text-field>
                            </v-col>
                        </v-row>
                        <v-row no-gutters>
                            <InputTypeValue :rules="phoneRules" v-model="staticData.phones" label="Phone"
                                :typeItems="staticData.phoneTypes">
                            </InputTypeValue>
                            <v-col>
                                <v-chip size="x-small" rounded v-for="(val, index) in staticData.phones" :key="index">
                                    {{ index }}
                                    {{ val }}
                                </v-chip>
                            </v-col>
                        </v-row>
                        <v-row class="pt-2" no-gutters>
                            <InputTypeValue :rules="emailRules" v-model="staticData.emails" label="Email"
                                :typeItems="staticData.emailTypes">
                            </InputTypeValue>
                            <v-col>
                                <v-chip size="x-small" rounded v-for="(val, index) in staticData.emails" :key="index">
                                    {{ index }}
                                    {{ val }}
                                </v-chip>
                            </v-col>
                        </v-row>
                    </v-container>
                </v-card>
            </div>
            <div class="ma-1 pa-1">
                <v-card min-width="200">
                    <v-card-title>Quick search</v-card-title>
                    <v-card-text>
                        asdasd
                        {{ staticData.phones }}
                    </v-card-text>
                </v-card>
            </div>
        </div>
    </v-dialog>
</template>


<script setup>
import { ref, onMounted } from 'vue';
import { useDisplay } from 'vuetify'
import InputTypeValue from './InputTypeValue.vue';

onMounted(() => {
    console.log("onMounted", useDisplay())
})

const staticData = ref({
    phones: {},
    phoneTypes: [
        "main",
        "home",
    ],
    emails: {},
    emailTypes: [
        "main",
        "home",
    ]
})

const phoneRules = [
    value => !!value || 'Required.',
]


const emailRules = [
    value => !!value || 'Required.',
    value => {
        const pattern = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
        return pattern.test(value) || 'Invalid e-mail.'
    },
]

const show = ref(true)

</script>