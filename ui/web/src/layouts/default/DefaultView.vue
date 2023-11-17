<template>
  <v-main>
    <main-toolbar />
    <error-dialog />
    <login-form />
    <router-view />
    <main-footer />
  </v-main>
</template>

<script setup>
import ErrorDialog from "@/components/ErrorDialog";
import LoginForm from '@/components/LoginForm'
import MainToolbar from "@/components/MainToolbar"
import MainFooter from '@/components/MainFooter'
</script>
