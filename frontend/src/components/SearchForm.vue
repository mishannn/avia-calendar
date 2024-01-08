<template>
  <v-form @submit.prevent="submit">
    <IATAInput v-model="formData.from" name="from" label="From (IATA)" />
    <IATAInput v-model="formData.to" name="to" label="To (IATA)" />
    <DatePicker
      v-model="formData.startDate"
      name="start_date"
      label="Start date"
    />
    <DatePicker v-model="formData.endDate" name="end_date" label="End date" />
    <v-text-field
      v-model="formData.maxTransfersAmount"
      label="Max transfers amount"
      name="max_transfers_amount"
      type="number"
    />
    <v-text-field
      v-model="formData.maxTransferDuration"
      label="Max transfer duration (hours)"
      name="max_transfer_duration"
      type="number"
    />

    <v-btn type="submit" color="primary" :disabled="!formValid" block>
      Search
    </v-btn>
  </v-form>
</template>

<script setup>
import { reactive, defineEmits } from "vue";
import DatePicker from "@/components/DatePicker.vue";
import { computed } from "vue";
import { watch } from "vue";
import { addDays } from "date-fns";
import IATAInput from "./IATAInput.vue";

const emit = defineEmits(["submit"]);

const formData = reactive({
  // from: "MOW",
  // to: "HKT",
  // startDate: new Date(),
  // endDate: addDays(new Date(), 30),
  from: "",
  to: "",
  startDate: undefined,
  endDate: undefined,
  maxTransfersAmount: 1,
  maxTransferDuration: 6,
});

watch(
  () => formData.startDate,
  (value) => {
    if (!formData.endDate) {
      formData.endDate = addDays(value, 14);
    }
  }
);

const formValid = computed(() => {
  return (
    !!formData.from &&
    !!formData.to &&
    !!formData.startDate &&
    !!formData.endDate
  );
});

function submit() {
  emit("submit", formData);
}
</script>
