<template>
  <v-main>
    <the-toolbar />
    <request-error />
    <login-form />
    <router-view />
    <the-footer />
  </v-main>
</template>

<script setup>
import RequestError from "@/components/RequestError.vue";
import TheFooter from '@/components/TheFooter'
import LoginForm from '@/components/LoginForm'
import TheToolbar from "@/components/TheToolbar"
</script>
