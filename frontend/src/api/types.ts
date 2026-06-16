export type ArtifactType =
  | "stack_trace"
  | "log"
  | "api_payload"
  | "validation_report"
  | "screenshot";

export type ExpirationPolicy = "7d" | "14d" | "never";

export type Artifact = {
  id: string;
  short_code: string;
  title: string;
  description: string;
  artifact_type: ArtifactType;
  service_name: string;
  environment: string;
  tags: string[];
  creator: string;
  file_name: string;
  content_type: string;
  size_bytes: number;
  created_at: string;
  expires_at?: string;
  preview?: string;
  download_url: string;
  share_url: string;
};

export type UploadArtifactInput = {
  title: string;
  description: string;
  artifactType: ArtifactType;
  serviceName: string;
  environment: string;
  tags: string;
  creator: string;
  expiration: ExpirationPolicy;
  file: File;
};

export type UploadArtifactResponse = {
  id: string;
  short_code: string;
  short_url: string;
  artifact: Artifact;
};

export type SearchArtifactsParams = {
  q?: string;
  service?: string;
  tag?: string;
  limit?: number;
};
