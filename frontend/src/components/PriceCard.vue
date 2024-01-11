<template>
  <v-card>
    <template v-slot:title>
      <span :style="{ color: titleColor }">{{ title }}</span>
    </template>
    <v-card-text>
      <v-alert
        v-if="ticketsEvent.error"
        :text="ticketsEvent.error"
        type="warning"
      />
      <div v-else style="display: flex; flex-direction: column; row-gap: 16px;">
        <PriceCardItem
          v-for="ticket in ticketsEvent.tickets"
          :ticket="ticket"
          :price-limit="priceLimit"
        />
      </div>
    </v-card-text>
  </v-card>
</template>

<script setup>
import { format, isWeekend } from "date-fns";
import { computed } from "vue";
import { defineProps } from "vue";
import PriceCardItem from "@/components/PriceCardItem.vue";

const props = defineProps({
  from: {
    type: String,
    required: true,
  },
  to: {
    type: String,
    required: true,
  },
  ticketsEvent: {
    type: Object,
    required: true,
  },
  priceLimit: {
    type: Number,
    required: true,
  },
});

const date = computed(() => {
  return new Date(props.ticketsEvent.date);
});

const title = computed(() => {
  return format(date.value, "E, dd MMM yyyy");
});

const titleColor = computed(() => {
  return isWeekend(date.value) ? "#c50000" : undefined;
});
</script>
