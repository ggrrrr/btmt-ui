<template>
    <v-app-bar elevation="1" prominent>
        <template v-slot:prepend>
            <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
        </template>
        <template v-slot:append>
            <v-label>{{ loginStore.email }}</v-label>
            <v-btn dense v-show="loginStore.email" icon="mdi-logout" @click="loginStore.resetLogin"></v-btn>
            <v-btn dense v-show="!loginStore.email" icon="mdi-login" @click="showLogin"></v-btn>
            <v-btn dense v-show="loginStore.email" icon="mdi-refresh" @click="validateClick"></v-btn>
        </template>
        <v-toolbar-title small><router-link to="/"><v-icon>mdi-home</v-icon></router-link>
            <!-- <v-spacer flat></v-spacer> -->
            <login-dialog />
        </v-toolbar-title>
    </v-app-bar>
    <v-navigation-drawer v-model="drawer" temporary>
        <v-list-item title="My Application" subtitle="Vuetify"></v-list-item>
        <v-divider></v-divider>
        <v-list-item link title="people" to="/people"></v-list-item>
        <v-list-item link title="List Item 2"></v-list-item>
        <v-list-item link title="List Item 3"></v-list-item>
    </v-navigation-drawer>
</template>

<script setup>
import { useLoginStore } from "@/store/auth";
import LoginDialog from './LoginDialog.vue'
import { ref } from 'vue'


let drawer = ref(true)
let loginStore = useLoginStore()

function showLogin() {
    console.log("showLogin")
    loginStore.showLogin = true
}

function validateClick() {
    console.log("validateClick")
    loginStore.validateRequest()
}


</script>
