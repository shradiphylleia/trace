import { Braces, GitPullRequest, ScrollText, TerminalSquare } from "lucide-react";
import type { Artifact } from "@/api/types";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";

type PreviewKind = "java" | "trace" | "git" | "misc";

type ArtifactPreviewProps = {
  artifact: Artifact;
};

export function ArtifactPreview({ artifact }: ArtifactPreviewProps) {
  if (!artifact.preview && artifact.artifact_type === "screenshot") {
    return <img className="max-h-[520px] rounded-md border object-contain" src={artifact.download_url} alt={artifact.title} />;
  }

  if (!artifact.preview) {
    return <p className="text-sm text-muted-foreground">No text preview is available for this file.</p>;
  }

  const kind = detectPreviewKind(artifact);
  const lines = artifact.preview.split(/\r?\n/).slice(0, 160);
  const Icon = iconForKind(kind);

  return (
    <div className="overflow-hidden rounded-md border bg-slate-950">
      <div className="flex items-center justify-between border-b border-slate-800 px-4 py-3">
        <div className="flex items-center gap-2 text-sm font-medium text-slate-100">
          <Icon className="h-4 w-4" />
          {labelForKind(kind)}
        </div>
        <Badge variant={badgeForKind(kind)}>{artifact.file_name}</Badge>
      </div>
      <pre className="max-h-[460px] overflow-auto p-4 font-mono text-xs leading-5 text-slate-100">
        {lines.map((line, index) => (
          <div key={`${index}-${line}`} className={cn("min-h-5 whitespace-pre-wrap", classForLine(kind, line))}>
            <span className="mr-4 inline-block w-8 select-none text-right text-slate-500">{index + 1}</span>
            <span>{line || " "}</span>
          </div>
        ))}
      </pre>
    </div>
  );
}

function detectPreviewKind(artifact: Artifact): PreviewKind {
  const text = `${artifact.file_name}\n${artifact.title}\n${artifact.description}\n${artifact.preview ?? ""}`.toLowerCase();
  if (text.includes(".java") || text.includes("exception in thread") || /\bat\s+[\w.$]+\(.*\.java:\d+\)/i.test(text)) {
    return "java";
  }
  if (text.includes("diff --git") || text.includes("fatal:") || text.includes("merge conflict") || text.includes("pull request")) {
    return "git";
  }
  if (artifact.artifact_type === "stack_trace" || text.includes("traceback") || text.includes("goroutine") || text.includes("panic:")) {
    return "trace";
  }
  return "misc";
}

function iconForKind(kind: PreviewKind) {
  if (kind === "java") return Braces;
  if (kind === "git") return GitPullRequest;
  if (kind === "trace") return TerminalSquare;
  return ScrollText;
}

function labelForKind(kind: PreviewKind) {
  if (kind === "java") return "Java exception preview";
  if (kind === "git") return "Git issue preview";
  if (kind === "trace") return "Trace preview";
  return "Artifact preview";
}

function badgeForKind(kind: PreviewKind) {
  if (kind === "java") return "warning" as const;
  if (kind === "git") return "info" as const;
  if (kind === "trace") return "danger" as const;
  return "default" as const;
}

function classForLine(kind: PreviewKind, line: string) {
  const lower = line.toLowerCase();
  if (kind === "java" && (line.includes("Exception") || line.includes("Caused by"))) {
    return "text-amber-200";
  }
  if (kind === "trace" && (lower.includes("panic:") || lower.includes("error") || lower.includes("failed"))) {
    return "text-red-200";
  }
  if (kind === "git" && (lower.startsWith("fatal:") || lower.includes("conflict"))) {
    return "text-red-200";
  }
  if (kind === "git" && (lower.startsWith("+") || lower.startsWith("-"))) {
    return lower.startsWith("+") ? "text-emerald-200" : "text-amber-200";
  }
  if (line.trim().startsWith("at ") || line.includes(".java:")) {
    return "text-sky-200";
  }
  if (lower.includes("warn")) {
    return "text-amber-200";
  }
  return "";
}
