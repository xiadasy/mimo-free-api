const ENV_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? ''

function getBaseUrl(): string {
  return localStorage.getItem('api_base_url') || ENV_BASE_URL
}

export async function apiFetch(path: string, init?: RequestInit): Promise<Response> {
  return fetch(`${getBaseUrl()}${path}`, init)
}
