<template>
  <div class="results-page">
    <div v-if="loading" class="progress-container">
      <v-progress-circular indeterminate color="primary" />
    </div>
    <v-btn v-else color="primary" @click="openNewSearch">New search</v-btn>

    <div
      v-if="errors.length !== 0"
      style="display: flex; flex-direction: column; row-gap: 8px"
    >
      <v-alert v-for="error in errors" :text="error" type="warning" />
    </div>

    <div
      v-if="sortedResults.length !== 0"
      style="display: flex; flex-direction: column; row-gap: 16px"
    >
      <v-select
        v-model="sortBy"
        label="Sort by"
        :items="['Date', 'Price']"
        hide-details
      />

      <PriceCard
        v-for="result in sortedResults"
        :key="result.date"
        :tickets-event="result"
        :from="from"
        :to="to"
        :price-limit="priceLimit"
      />
    </div>
  </div>
</template>

<script setup>
import PriceCard from "@/components/PriceCard.vue";
import { useTitle } from "@vueuse/core";
import { useRouteQuery } from "@vueuse/router";
import { mean, median, quantileSeq } from "mathjs";
import { onUnmounted, ref, watch, computed } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();

const from = useRouteQuery("from");
const to = useRouteQuery("to");
const startDate = useRouteQuery("start_date");
const endDate = useRouteQuery("end_date");
const maxTransfersAmount = useRouteQuery("max_transfers_amount");
const maxTransferDuration = useRouteQuery("max_transfer_duration");

useTitle(`Prices from ${from.value} to ${to.value}`);

const allParamsFilled = computed(() => {
  return (
    !!from.value &&
    !!to.value &&
    !!startDate.value &&
    !!endDate.value &&
    !!maxTransfersAmount.value &&
    !!maxTransferDuration.value
  );
});

watch(
  allParamsFilled,
  (value) => {
    if (!value) {
      router.push("/");
    }
  },
  {
    immediate: true,
  }
);

const loading = ref(true);
const ticketsEvents = ref([]);
const errors = ref([]);

const searchParams = new URLSearchParams({
  from: from.value,
  to: to.value,
  start_date: startDate.value,
  end_date: endDate.value,
  max_transfers_amount: maxTransfersAmount.value,
  max_transfer_duration: maxTransferDuration.value,
});

const eventSource = new EventSource(
  `/api/v1/find-cheapest-tickets?${searchParams.toString()}`
);

eventSource.addEventListener("tickets", (event) => {
  let ticketsEvent = undefined;
  try {
    ticketsEvent = JSON.parse(event.data);
  } catch (err) {
    errors.value = [...errors.value, err.message];
    return;
  }

  ticketsEvents.value = [...ticketsEvents.value, ticketsEvent];
});

eventSource.addEventListener("marshalerror", (event) => {
  errors.value = [...errors.value, event.data];
});

eventSource.addEventListener("close", () => {
  eventSource.close();
  loading.value = false;
});

eventSource.addEventListener("error", () => {
  eventSource.close();
  errors.value = [...errors.value, "can't connect to server"];
  loading.value = false;
});

eventSource.addEventListener("requesterror", (event) => {
  eventSource.close();
  errors.value = [...errors.value, event.data];
  loading.value = false;
});

onUnmounted(() => {
  eventSource.close();
});

const priceLimit = computed(() => {
  const arr = ticketsEvents.value.reduce((prices, item) => {
    return [...prices, ...item.tickets.map((ticket) => ticket.price)];
  }, []);

  if (arr.length == 0) {
    return 0;
  }

  return mean(mean(arr));
});

const sortBy = ref("Date");
const sortedResults = computed(() => {
  return [...ticketsEvents.value].sort((a, b) => {
    if (sortBy.value == "Date") {
      return a.date.localeCompare(b.date);
    }

    if (sortBy.value == "Price") {
      if (a.tickets[0] && b.tickets[0]) {
        return a.tickets[0].price - b.tickets[0].price;
      } else if (!a.tickets[0]) {
        return 1
      } else {
        return -1
      }
    }

    return 0;
  });
});

function openNewSearch() {
  router.push("/");
}
</script>

<style lang="scss">
.results-page {
  padding: 16px;
  display: flex;
  flex-direction: column;
  row-gap: 16px;
}

.progress-container {
  display: flex;
  justify-content: center;
  padding: 8px;
}
</style>
