<template>
  <v-menu v-model="showDatePicker" :close-on-content-click="false">
    <template v-slot:activator="{ props }">
      <v-text-field
        :label="label"
        :name="name"
        :model-value="formattedDate"
        readonly
        v-bind="props"
      />
    </template>
    <v-date-picker :model-value="modelValue" @update:model-value="onChange" />
  </v-menu>
</template>

<script setup>
import { format } from "date-fns";
import { computed } from "vue";
import { ref, defineProps } from "vue";

const props = defineProps({
  modelValue: {
    type: Date,
  },
  label: {
    type: String,
  },
  name: {
    type: String,
  },
});

const emit = defineEmits(["update:model-value"]);

const formattedDate = computed(() => {
  if (!props.modelValue) {
    return "";
  }

  return format(props.modelValue, "E, dd MMM yyyy");
});

const showDatePicker = ref(false);

function onChange(value) {
  showDatePicker.value = false;
  emit("update:model-value", value);
}
</script>
