<script setup lang="ts">
import { computed } from 'vue'
import { marked } from 'marked'
import DOMPurify from 'dompurify'

const props = defineProps<{ source: string }>()

const html = computed(() => {
  const raw = marked.parse(props.source, { async: false }) as string
  return DOMPurify.sanitize(raw)
})
</script>

<template>
  <div class="markdown" v-html="html" />
</template>
