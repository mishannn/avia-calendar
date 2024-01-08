<template>
  <v-card>
    <template v-slot:title>
      <span :style="{ color: titleColor }">{{ title }}</span>
    </template>
    <v-card-text>
      <v-alert v-if="price.error" :text="price.error" type="warning" />
      <div
        v-else
        style="
          display: grid;
          grid-template-columns: 1fr 1fr 1fr;
          column-gap: 11px;
        "
      >
        <LabelledValue label="Price">
          <span :style="{ color: priceColor }">{{ price.price }} RUB</span>
        </LabelledValue>
        <LabelledValue label="Transfers">
          {{ price.transfers_amount }}
        </LabelledValue>
        <LabelledValue label="Aviasales">
          <a
            :href="price.search_link"
            target="_blank"
            style="text-decoration: none; color: #1867c0"
          >
            Open
          </a>
        </LabelledValue>
      </div>
    </v-card-text>
  </v-card>
</template>

<script setup>
import { format, isWeekend } from "date-fns";
import { computed } from "vue";
import { defineProps } from "vue";
import LabelledValue from "@/components/LabelledValue.vue";

const props = defineProps({
  from: {
    type: String,
    required: true,
  },
  to: {
    type: String,
    required: true,
  },
  price: {
    type: Object,
    required: true,
  },
  priceLimit: {
    type: Number,
    required: true,
  },
});

const date = computed(() => {
  return new Date(props.price.date);
});

const title = computed(() => {
  return format(date.value, "E, dd MMM yyyy");
});

const titleColor = computed(() => {
  return isWeekend(date.value) ? "#c50000" : undefined;
});

const priceColor = computed(() => {
  return props.price.price <= props.priceLimit ? "#089c08" : undefined;
});
</script>
