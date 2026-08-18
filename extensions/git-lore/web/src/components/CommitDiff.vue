<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  subject: string
  body?: string
  sha: string
  diff: string
}>()

interface DiffLine {
  text: string
  kind: 'add' | 'del' | 'hunk' | 'ctx'
}

const lines = computed(() =>
  props.diff.split('\n').map((text): DiffLine => {
    if (text.startsWith('+++') || text.startsWith('---')) return { text, kind: 'hunk' }
    if (text.startsWith('@@')) return { text, kind: 'hunk' }
    if (text.startsWith('+')) return { text, kind: 'add' }
    if (text.startsWith('-')) return { text, kind: 'del' }
    return { text, kind: 'ctx' }
  }),
)
</script>

<template>
  <div>
    <h3>{{ subject }}</h3>
    <p class="commit-meta">{{ sha }}</p>
    <pre v-if="body" class="markdown">{{ body }}</pre>
    <div class="diff">
      <div v-for="(line, i) in lines" :key="i" :class="line.kind">{{ line.text }}</div>
    </div>
  </div>
</template>
