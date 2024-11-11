<template>
  <v-layout class="rounded rounded-md">
    <main-toolbar />
    <error-dialog />
    <v-main>
      <router-view />
    </v-main>
  </v-layout>
  <main-footer />
</template>

<script setup>
import ErrorDialog from "@/components/ErrorDialog";
import MainToolbar from "@/components/MainToolbar"
import MainFooter from '@/components/MainFooter'
</script>
