import type { ArtifactType } from "@/api/types";

export function formatArtifactType(type: ArtifactType) {
  return type.split("_").join(" ");
}
