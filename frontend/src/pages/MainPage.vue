<template>
  <div class="main-page">
    <SearchForm class="search-form" @submit="onFormSubmit" />
  </div>
</template>

<script setup>
import SearchForm from "@/components/SearchForm.vue";
import { useTitle } from "@vueuse/core";
import { format } from "date-fns";
import { useRouter } from "vue-router";

useTitle("Avia Calendar");

const router = useRouter();

function onFormSubmit(formData) {
  const dateFormat = "yyyy-MM-dd";

  router.push({
    path: "/results",
    query: {
      from: formData.from.toUpperCase(),
      to: formData.to.toUpperCase(),
      start_date: format(formData.startDate, dateFormat),
      end_date: format(formData.endDate, dateFormat),
      max_transfers_amount: formData.maxTransfersAmount,
      max_transfer_duration: formData.maxTransferDuration,
    },
  });
}
</script>

<style>
.search-form {
  padding: 16px;
}
</style>
