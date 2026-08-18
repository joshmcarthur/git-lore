<script setup lang="ts">
import { computed } from 'vue'
import type { SyncState } from '../api/client'

const props = defineProps<{
  state?: SyncState
  label?: string
}>()

const text = computed(() => props.label || props.state || '')
const cls = computed(() => (props.state ? `badge ${props.state}` : 'badge'))
</script>

<template>
  <span v-if="text" :class="cls">{{ text }}</span>
</template>
