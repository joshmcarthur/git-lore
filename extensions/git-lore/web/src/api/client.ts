export type SyncState =
  | 'synced'
  | 'local-only'
  | 'remote-only'
  | 'ahead'
  | 'behind'
  | 'diverged'

export interface RepoInfo {
  root: string
  currentBranch: string
  remotes: string[]
}

export interface WorkSummary {
  id: string
  ref: string
  commit: string
  subject: string
}

export interface WorkDetail extends WorkSummary {
  files: string[]
}

export interface BranchAssociation {
  branch: string
  workId: string
}

export interface CommitInfo {
  sha: string
  subject: string
  author: string
  date: string
}

export interface CommitDetail extends CommitInfo {
  parents: string[]
  body: string
  diff: string
}

export interface WorkRemoteStatus {
  id: string
  localSha?: string
  remoteSha?: string
  state: SyncState
}

export interface RemoteStatusResult {
  remote: string
  refspecConfigured: boolean
  works: WorkRemoteStatus[]
}

export interface FetchResult {
  remote: string
  workId?: string
  refspecAdded: boolean
  output: string
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(path, init)
  if (!res.ok) {
    let msg = res.statusText
    try {
      const body = await res.json()
      if (body?.error) msg = body.error
    } catch {
      /* ignore */
    }
    throw new Error(msg)
  }
  const ct = res.headers.get('content-type') || ''
  if (ct.includes('application/json')) {
    return res.json() as Promise<T>
  }
  return res.text() as unknown as T
}

export const api = {
  repo: () => request<RepoInfo>('/api/repo'),
  branchAssociations: () => request<BranchAssociation[]>('/api/branch-associations'),
  works: () => request<WorkSummary[]>('/api/works'),
  work: (id: string) => request<WorkDetail>(`/api/works/${encodeURIComponent(id)}`),
  file: async (id: string, path: string) => {
    const res = await fetch(
      `/api/works/${encodeURIComponent(id)}/files/${path.split('/').map(encodeURIComponent).join('/')}`,
    )
    if (!res.ok) {
      let msg = res.statusText
      try {
        const body = await res.json()
        if (body?.error) msg = body.error
      } catch {
        /* ignore */
      }
      throw new Error(msg)
    }
    return res.text()
  },
  log: (id: string, limit = 20) =>
    request<CommitInfo[]>(`/api/works/${encodeURIComponent(id)}/log?limit=${limit}`),
  commit: (id: string, sha: string) =>
    request<CommitDetail>(
      `/api/works/${encodeURIComponent(id)}/commits/${encodeURIComponent(sha)}`,
    ),
  remoteStatus: (remote = 'origin') =>
    request<RemoteStatusResult>(`/api/remote/status?remote=${encodeURIComponent(remote)}`),
  fetchLore: (remote = 'origin', workId?: string) =>
    request<FetchResult>('/api/remote/fetch', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ remote, workId: workId || undefined }),
    }),
}
