<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  api,
  type BranchAssociation,
  type CommitDetail,
  type CommitInfo,
  type RemoteStatusResult,
  type RepoInfo,
  type SyncState,
  type WorkSummary,
} from '../api/client'
import RemoteBadge from '../components/RemoteBadge.vue'
import MarkdownView from '../components/MarkdownView.vue'
import CommitDiff from '../components/CommitDiff.vue'

const route = useRoute()
const router = useRouter()

const repo = ref<RepoInfo | null>(null)
const works = ref<WorkSummary[]>([])
const associations = ref<BranchAssociation[]>([])
const remoteStatus = ref<RemoteStatusResult | null>(null)
const selectedRemote = ref('origin')
const loading = ref(true)
const error = ref('')
const statusError = ref('')
const fetching = ref(false)

const workId = computed(() => (route.params.id as string) || '')
const filePath = computed(() => (route.params.path as string) || '')

const workFiles = ref<string[]>([])
const workSubject = ref('')
const workCommit = ref('')
const fileContent = ref('')
const commits = ref<CommitInfo[]>([])
const selectedCommit = ref<CommitDetail | null>(null)
const detailError = ref('')
const detailLoading = ref(false)

const currentWorkId = computed(() => {
  if (!repo.value) return ''
  const hit = associations.value.find((a) => a.branch === repo.value!.currentBranch)
  return hit?.workId || ''
})

const statusById = computed(() => {
  const map = new Map<string, SyncState>()
  for (const w of remoteStatus.value?.works || []) {
    map.set(w.id, w.state)
  }
  return map
})

const remoteOnlyWorks = computed(() =>
  (remoteStatus.value?.works || []).filter((w) => w.state === 'remote-only'),
)

async function loadShell() {
  loading.value = true
  error.value = ''
  try {
    const [r, w, a] = await Promise.all([
      api.repo(),
      api.works(),
      api.branchAssociations(),
    ])
    repo.value = r
    works.value = w
    associations.value = a
    if (r.remotes.length && !r.remotes.includes(selectedRemote.value)) {
      selectedRemote.value = r.remotes[0]
    }
    if (!workId.value && w.length) {
      const prefer = currentWorkId.value && w.some((x) => x.id === currentWorkId.value)
        ? currentWorkId.value
        : w[0].id
      await router.replace({ name: 'work', params: { id: prefer } })
    }
  } catch (e) {
    error.value = e instanceof Error ? e.message : String(e)
  } finally {
    loading.value = false
  }
}

async function refreshStatus() {
  statusError.value = ''
  try {
    remoteStatus.value = await api.remoteStatus(selectedRemote.value)
  } catch (e) {
    statusError.value = e instanceof Error ? e.message : String(e)
    remoteStatus.value = null
  }
}

async function fetchLore(id?: string) {
  fetching.value = true
  statusError.value = ''
  try {
    await api.fetchLore(selectedRemote.value, id)
    await loadShell()
    await refreshStatus()
  } catch (e) {
    statusError.value = e instanceof Error ? e.message : String(e)
  } finally {
    fetching.value = false
  }
}

async function loadWork(id: string) {
  if (!id) return
  detailLoading.value = true
  detailError.value = ''
  selectedCommit.value = null
  try {
    const [detail, log] = await Promise.all([api.work(id), api.log(id, 30)])
    workFiles.value = detail.files
    workSubject.value = detail.subject
    workCommit.value = detail.commit
    commits.value = log

    const path = filePath.value || detail.files[0] || ''
    if (path && !filePath.value && detail.files[0]) {
      await router.replace({
        name: 'file',
        params: { id, path: detail.files[0] },
      })
      return
    }
    if (path) {
      fileContent.value = await api.file(id, path)
    } else {
      fileContent.value = ''
    }
  } catch (e) {
    detailError.value = e instanceof Error ? e.message : String(e)
  } finally {
    detailLoading.value = false
  }
}

async function selectCommit(sha: string) {
  if (!workId.value) return
  detailError.value = ''
  try {
    selectedCommit.value = await api.commit(workId.value, sha)
  } catch (e) {
    detailError.value = e instanceof Error ? e.message : String(e)
  }
}

function shortSha(sha: string) {
  return sha.slice(0, 7)
}

onMounted(async () => {
  await loadShell()
  await refreshStatus()
  if (workId.value) await loadWork(workId.value)
})

watch(
  () => [route.params.id, route.params.path],
  async () => {
    if (workId.value) await loadWork(workId.value)
  },
)

watch(selectedRemote, () => {
  void refreshStatus()
})
</script>

<template>
  <div class="layout">
    <aside class="sidebar">
      <div class="sidebar-header">
        <h1>git-lore</h1>
        <div v-if="repo" class="meta">
          {{ repo.currentBranch || '(detached)' }} · {{ repo.root }}
        </div>
      </div>

      <div class="toolbar">
        <select v-model="selectedRemote" :disabled="!repo?.remotes.length">
          <option v-for="r in repo?.remotes || []" :key="r" :value="r">{{ r }}</option>
        </select>
        <button :disabled="fetching || !repo?.remotes.length" @click="refreshStatus">
          Status
        </button>
        <button
          class="primary"
          :disabled="fetching || !repo?.remotes.length"
          @click="fetchLore()"
        >
          Fetch
        </button>
      </div>

      <div v-if="remoteStatus && !remoteStatus.refspecConfigured" class="warn-banner">
        Lore fetch refspec not configured for {{ remoteStatus.remote }}. Fetch will add
        <code>refs/lore/*:refs/lore/*</code> (no force).
      </div>
      <div v-if="statusError" class="warn-banner">{{ statusError }}</div>

      <div v-if="loading" class="loading">Loading works…</div>
      <div v-else-if="error" class="error">{{ error }}</div>
      <template v-else>
        <ul v-if="works.length" class="work-list">
          <li v-for="w in works" :key="w.id">
            <router-link
              :to="{ name: 'work', params: { id: w.id } }"
              :class="{ active: w.id === workId }"
            >
              <div class="work-id">
                {{ w.id }}
                <RemoteBadge v-if="w.id === currentWorkId" label="branch" class="current" />
                <RemoteBadge v-if="statusById.get(w.id)" :state="statusById.get(w.id)" />
              </div>
              <div class="work-subject">{{ w.subject }}</div>
            </router-link>
          </li>
        </ul>
        <div v-else class="empty">
          No local Lore Works. Create one with the create-lore skill, or fetch from a remote.
          <pre>git for-each-ref refs/lore</pre>
        </div>

        <div v-if="remoteOnlyWorks.length" class="sidebar-header">
          <strong>Remote only</strong>
          <ul class="work-list">
            <li v-for="w in remoteOnlyWorks" :key="w.id">
              <div class="work-id">
                {{ w.id }}
                <RemoteBadge state="remote-only" />
                <button :disabled="fetching" @click="fetchLore(w.id)">Fetch</button>
              </div>
            </li>
          </ul>
        </div>
      </template>
    </aside>

    <main class="main">
      <div v-if="!workId" class="empty">Select a Lore Work.</div>
      <template v-else>
        <div class="detail-header">
          <h2>{{ workId }}</h2>
          <div class="sub">
            refs/lore/{{ workId }} · {{ shortSha(workCommit) }} · {{ workSubject }}
          </div>
        </div>

        <div v-if="detailLoading" class="loading">Loading…</div>
        <div v-else-if="detailError" class="error">{{ detailError }}</div>
        <div v-else class="panels">
          <nav class="file-list">
            <h3>Documents</h3>
            <ul>
              <li v-for="f in workFiles" :key="f">
                <router-link
                  :to="{ name: 'file', params: { id: workId, path: f } }"
                  :class="{ active: f === filePath }"
                  @click="selectedCommit = null"
                >
                  {{ f }}
                </router-link>
              </li>
            </ul>
          </nav>

          <section class="content">
            <CommitDiff
              v-if="selectedCommit"
              :subject="selectedCommit.subject"
              :body="selectedCommit.body"
              :sha="selectedCommit.sha"
              :diff="selectedCommit.diff"
            />
            <MarkdownView v-else-if="fileContent" :source="fileContent" />
            <div v-else class="empty">No document selected.</div>
          </section>

          <aside class="history">
            <h3>History</h3>
            <ul>
              <li v-for="c in commits" :key="c.sha">
                <button
                  class="commit"
                  :class="{ active: selectedCommit?.sha === c.sha }"
                  @click="selectCommit(c.sha)"
                >
                  <div class="commit-subject">{{ c.subject }}</div>
                  <div class="commit-meta">{{ shortSha(c.sha) }} · {{ c.date.slice(0, 10) }}</div>
                </button>
              </li>
            </ul>
          </aside>
        </div>
      </template>
    </main>
  </div>
</template>
