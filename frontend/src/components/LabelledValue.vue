<template>
  <div>
    <div class="text-grey">{{ label }}</div>
    <div>
      <slot />
    </div>
  </div>
</template>

<script setup>
import { defineProps } from "vue";

const props = defineProps({
  label: {
    type: String,
    required: true,
  },
  value: {
    type: String,
    default: "-",
  },
});
</script>
