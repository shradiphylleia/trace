import { FormEvent, useEffect, useState } from "react";
import { Copy, Download, Search } from "lucide-react";
import type { Artifact } from "@/api/types";
import { ArtifactPreview } from "@/components/artifact-preview";
import { CopyToast } from "@/components/copy-toast";
import { PageHeader } from "@/components/page-header";
import { EmptyState, ErrorState, LoadingState } from "@/components/state";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { formatArtifactType } from "@/lib/artifact";
import { formatDate, expirationLabel } from "@/lib/date";

type ArtifactDetailsPageProps = {
  shortCode: string;
  loadArtifact: (shortCode: string) => Promise<Artifact>;
  onLookup: (shortCode: string) => void;
  onSearch: (query: string) => void;
};

export function ArtifactDetailsPage({ shortCode, loadArtifact, onLookup, onSearch }: ArtifactDetailsPageProps) {
  const [lookupCode, setLookupCode] = useState(shortCode);
  const [artifact, setArtifact] = useState<Artifact | null>(null);
  const [loading, setLoading] = useState(Boolean(shortCode));
  const [error, setError] = useState("");
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    setLookupCode(shortCode);
    if (!shortCode) {
      setArtifact(null);
      setLoading(false);
      setError("");
      return;
    }

    let ignore = false;
    setLoading(true);
    setError("");
    loadArtifact(shortCode)
      .then((item) => {
        if (!ignore) setArtifact(item);
      })
      .catch((err) => {
        if (!ignore) {
          setArtifact(null);
          setError(err instanceof Error ? err.message : "Artifact not found");
        }
      })
      .finally(() => {
        if (!ignore) setLoading(false);
      });
    return () => {
      ignore = true;
    };
  }, [shortCode, loadArtifact]);

  function handleLookup(event: FormEvent) {
    event.preventDefault();
    if (lookupCode.trim()) onLookup(lookupCode.trim());
  }

  async function copyLink() {
    if (!artifact) return;
    await navigator.clipboard.writeText(artifact.share_url);
    setCopied(true);
    window.setTimeout(() => setCopied(false), 1500);
  }

  const expiration = expirationLabel(artifact?.expires_at);

  return (
    <>
      <PageHeader
        title="Artifact details"
        description="Open a short code, copy the debugging link, or download the original artifact."
      />
      <Card>
        <CardContent className="p-5">
          <form className="flex flex-col gap-2 sm:flex-row" onSubmit={handleLookup}>
            <Input value={lookupCode} placeholder="abc123" onChange={(event) => setLookupCode(event.target.value)} />
            <Button type="submit">
              <Search className="h-4 w-4" />
              Lookup
            </Button>
          </form>
        </CardContent>
      </Card>
      {loading && <LoadingState label="Loading artifact" />}
      {error && <ErrorState message={error} />}
      {!loading && !error && !artifact && <EmptyState title="No artifact selected" description="Paste a short code to open a TraceShare artifact." />}
      {artifact && (
        <div className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_360px]">
          <Card>
            <CardHeader>
              <CardTitle>{artifact.title}</CardTitle>
            </CardHeader>
            <CardContent className="grid gap-4">
              <p className="text-sm text-muted-foreground">{artifact.description || "No description provided."}</p>
              <div className="flex flex-wrap gap-2">
                <Badge variant="outline">{formatArtifactType(artifact.artifact_type)}</Badge>
                <Badge variant={expiration.tone}>{expiration.text}</Badge>
                {artifact.tags.map((tag) => (
                  <Badge key={tag}>{tag}</Badge>
                ))}
              </div>
              <div className="rounded-md border bg-muted/20 p-4">
                <div className="mb-2 text-sm font-medium">Artifact preview</div>
                <ArtifactPreview artifact={artifact} />
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardHeader>
              <CardTitle>Metadata</CardTitle>
            </CardHeader>
            <CardContent className="grid gap-4 text-sm">
              <Metadata label="Short code" value={artifact.short_code} />
              <Metadata label="Service" value={artifact.service_name} />
              <Metadata label="Environment" value={artifact.environment} />
              <Metadata label="Creator" value={artifact.creator} />
              <Metadata label="Created" value={formatDate(artifact.created_at)} />
              <Metadata label="Expires" value={formatDate(artifact.expires_at)} />
              <Metadata label="File" value={artifact.file_name} />
              <Metadata label="Share URL" value={artifact.share_url} />
              <div className="grid gap-2 pt-2">
                <Button type="button" onClick={copyLink}>
                  <Copy className="h-4 w-4" />
                  Copy link
                </Button>
                <Button type="button" variant="outline" onClick={() => window.open(artifact.download_url, "_blank")}>
                  <Download className="h-4 w-4" />
                  Download
                </Button>
                <Button type="button" variant="ghost" onClick={() => onSearch(artifact.service_name)}>
                  Search service
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      )}
      <CopyToast visible={copied} />
    </>
  );
}

function Metadata({ label, value }: { label: string; value: string }) {
  return (
    <div className="grid gap-1">
      <div className="text-xs font-semibold uppercase text-muted-foreground">{label}</div>
      <div className="break-words text-foreground">{value}</div>
    </div>
  );
}
