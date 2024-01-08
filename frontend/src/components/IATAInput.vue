<template>
  <v-menu v-model="isOpen">
    <template v-slot:activator="{ props }">
      <v-text-field
        :label="label"
        :name="name"
        :model-value="modelValue"
        autocomplete="off"
        @update:model-value="onChange"
        v-bind="props"
      />
    </template>
    <v-list lines="one">
      <div v-if="loading" class="progress-container">
        <v-progress-circular indeterminate color="primary" />
      </div>
      <template v-else>
        <v-list-item
          v-if="!debouncedModelValue"
          title="Start enter text to show hints"
          disabled
        />
        <v-list-item
          v-else-if="items.length === 0"
          title="Can't find IATA codes"
          disabled
        />
        <template v-else>
          <v-list-item
            v-for="item in items"
            :key="item.code"
            :title="item.code"
            :subtitle="`${item.name}, ${item.country_name}`"
            @click="onSelectHint(item)"
          />
        </template>
      </template>
    </v-list>
  </v-menu>
</template>

<script setup>
import { ref } from "vue";
import { refDebounced } from "@vueuse/core";
import { watch } from "vue";
import { computed } from "vue";

const props = defineProps({
  modelValue: {
    type: String,
  },
  label: {
    type: String,
  },
  name: {
    type: String,
  },
});

const emit = defineEmits(["update:model-value"]);

const isOpen = ref(false);

const loading = ref(false);
const items = ref([]);

function onChange(value) {
  emit("update:model-value", value);
}

function onSelectHint(item) {
  emit("update:model-value", item.code);
}

const debouncedModelValue = refDebounced(
  computed(() => props.modelValue),
  500
);

watch(debouncedModelValue, () => {
  if (!isOpen.value) {
    items.value = [];
    return;
  }

  searchIATA();
});

watch(isOpen, () => {
  if (isOpen.value) {
    searchIATA();
  }
});

async function searchIATA() {
  if (!debouncedModelValue.value) {
    items.value = [];
    return;
  }

  loading.value = true;
  try {
    const response = await fetch(
      `/api/v1/find-iata?filter=${encodeURIComponent(
        debouncedModelValue.value
      )}`
    );

    if (response.status !== 200) {
      items.value = [];
      return;
    }

    items.value = await response.json();
  } catch (err) {
    items.value = [];
  } finally {
    loading.value = false;
  }
}
</script>

<style lang="scss">
.progress-container {
  display: flex;
  justify-content: center;
  padding: 8px;
}
</style>
