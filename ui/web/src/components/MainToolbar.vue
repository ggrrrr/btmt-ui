<template>
    <v-app-bar elevation="1" prominent>
        <template v-slot:prepend>
            <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
        </template>
        <template v-slot:append>
            <login-dialog />
            <v-chip v-if="loginStore.username" rounded color="primary">{{ loginStore.username }}</v-chip>
            <v-btn color="warning" title="exit" dense rounded v-show="loginStore.username" icon="mdi-logout"
                @click="loginStore.resetLogin" />
            <v-btn color="secondary" title="login" dense v-show="!loginStore.username" icon="mdi-login"
                @click="showLogin" />
            <v-btn color="secondary" title="refresh" dense v-show="loginStore.username" icon="mdi-refresh"
                @click="validateClick" />
        </template>
        <v-toolbar-title small>
            <!-- <v-spacer flat></v-spacer> -->
        </v-toolbar-title>
    </v-app-bar>
    <v-navigation-drawer v-model="drawer" temporary>
        <v-list-item color="primary" prepend-icon="mdi-home" title="Home" link to="/"></v-list-item>
        <v-divider></v-divider>
        <v-list-item color="primary" prepend-icon="mdi-account-group" link title="People" to="/people"></v-list-item>
        <v-list-item color="primary" prepend-icon="mdi-file-document-edit" link title="Templates"
            to="/templates/manage"></v-list-item>
        <v-list-item color="primary" prepend-icon="mdi-glass-mug-variant" link title="todo" to="/todo"></v-list-item>
        <v-list-item color="primary" prepend-icon="mdi-account-cog-outline" link title="Users"
            to="/users"></v-list-item>
    </v-navigation-drawer>
</template>

<script setup>
import { useLoginStore } from "@/store/auth";
import { ref } from 'vue'

import LoginDialog from './LoginDialog.vue'


let drawer = ref(false)
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
