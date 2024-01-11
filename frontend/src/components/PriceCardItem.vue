<template>
  <div
    class="price-card-item"
    style="
      display: grid;
      grid-template-columns: 1fr 1fr 1fr 1fr;
      column-gap: 11px;
    "
  >
    <LabelledValue label="Price">
      <span :style="{ color: priceColor }">{{ ticket.price }} RUB</span>
    </LabelledValue>
    <LabelledValue label="Transfers">
      {{ ticket.transfers_amount }}
    </LabelledValue>
    <LabelledValue label="Flights">
      <v-dialog width="500">
        <template v-slot:activator="{ props }">
          <a style="text-decoration: none; color: #1867c0" v-bind="props">
            Open
          </a>
        </template>

        <template v-slot:default="{ isActive }">
          <v-card title="Flights">
            <v-card-text>
              <v-timeline side="end" truncate-line="both">
                <v-timeline-item v-for="flight in ticket.flights" dot-color="#9f9f9f" icon="mdi-airplane">
                  <div>Depart from {{ flight.from }} at {{ flight.departure_time }}</div>
                  <div>Arrive to {{ flight.to }} at {{ flight.arrival_time }}</div>
                  <div>With {{ flight.airline_name }}</div>
                </v-timeline-item>
              </v-timeline>
            </v-card-text>

            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn text="Close" @click="isActive.value = false"></v-btn>
            </v-card-actions>
          </v-card>
        </template>
      </v-dialog>
    </LabelledValue>
    <LabelledValue label="Aviasales">
      <a
        :href="ticket.search_link"
        target="_blank"
        style="text-decoration: none; color: #1867c0"
      >
        Open
      </a>
    </LabelledValue>
  </div>
</template>

<script setup>
import { computed, defineProps } from "vue";
import LabelledValue from "./LabelledValue.vue";

const props = defineProps({
  ticket: {
    type: Object,
    required: true,
  },
  priceLimit: {
    type: Number,
    required: true,
  },
});

const priceColor = computed(() => {
  return props.ticket.price <= props.priceLimit ? "#089c08" : undefined;
});

function openDetails() {
  console.log("openDetails");
}
</script>
