<template>
  <v-app>
    <main-toolbar />
    <error-dialog />
    <router-view />
    <main-footer />
  </v-app>
</template>

<script setup>
import ErrorDialog from "@/components/ErrorDialog";
import MainToolbar from "@/components/MainToolbar"
import MainFooter from '@/components/MainFooter'
</script>
