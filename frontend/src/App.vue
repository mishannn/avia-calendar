<template>
  <v-app>
    <v-app-bar color="primary">
      <v-app-bar-title>{{ title }}</v-app-bar-title>
    </v-app-bar>

    <v-main class="main-container">
      <div class="router-container">
        <router-view />
      </div>
    </v-main>
  </v-app>
</template>

<script setup>
import { useTitle } from "@vueuse/core";

const title = useTitle(undefined, { observe: true });
</script>

<style>
.main-container {
  display: flex;
  justify-content: center;
}

.router-container {
  flex: 1 0 auto;
  max-width: min(100%, 600px);
}
</style>
