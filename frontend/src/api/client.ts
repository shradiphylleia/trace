import type { Artifact, SearchArtifactsParams, UploadArtifactInput, UploadArtifactResponse } from "./types";

const apiBaseUrl = import.meta.env.VITE_API_BASE_URL ?? "";

type SearchResponse = {
  items: Artifact[];
};

export class ApiError extends Error {
  status: number;

  constructor(message: string, status: number) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${apiBaseUrl}${path}`, init);
  const contentType = response.headers.get("content-type") ?? "";

  if (!response.ok) {
    if (contentType.includes("application/json")) {
      const body = (await response.json()) as { error?: string };
      throw new ApiError(body.error ?? "Request failed", response.status);
    }
    throw new ApiError("Request failed", response.status);
  }

  return (await response.json()) as T;
}

export async function searchArtifacts(params: SearchArtifactsParams = {}) {
  const query = new URLSearchParams();
  if (params.q) query.set("q", params.q);
  if (params.service) query.set("service", params.service);
  if (params.tag) query.set("tag", params.tag);
  if (params.limit) query.set("limit", String(params.limit));

  const response = await request<SearchResponse>(`/api/search?${query.toString()}`);
  return response.items;
}

export async function getArtifact(shortCode: string) {
  return request<Artifact>(`/api/artifacts/${encodeURIComponent(shortCode)}`);
}

export async function uploadArtifact(input: UploadArtifactInput) {
  const form = new FormData();
  form.set("title", input.title);
  form.set("description", input.description);
  form.set("artifact_type", input.artifactType);
  form.set("service_name", input.serviceName);
  form.set("environment", input.environment);
  form.set("tags", input.tags);
  form.set("creator", input.creator);
  form.set("expiration", input.expiration);
  form.set("file", input.file);

  return request<UploadArtifactResponse>("/api/artifacts", {
    method: "POST",
    body: form,
  });
}
