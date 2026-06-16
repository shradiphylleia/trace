import { FormEvent, useState } from "react";
import type { ReactNode } from "react";
import { ArrowLeft, UploadCloud } from "lucide-react";
import { uploadArtifact } from "@/api/client";
import type { Artifact, ArtifactType, ExpirationPolicy } from "@/api/types";
import { PageHeader } from "@/components/page-header";
import { ErrorState } from "@/components/state";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select } from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";

type UploadPageProps = {
  onUploaded: (artifact: Artifact) => void;
  onCancel: () => void;
};

export function UploadPage({ onUploaded, onCancel }: UploadPageProps) {
  const [artifactType, setArtifactType] = useState<ArtifactType>("stack_trace");
  const [expiration, setExpiration] = useState<ExpirationPolicy>("7d");
  const [file, setFile] = useState<File | null>(null);
  const [error, setError] = useState("");
  const [saving, setSaving] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    if (!file) {
      setError("Choose a file to upload.");
      return;
    }

    const form = new FormData(event.currentTarget);
    setSaving(true);
    try {
      const response = await uploadArtifact({
        title: String(form.get("title") ?? ""),
        description: String(form.get("description") ?? ""),
        artifactType,
        serviceName: String(form.get("serviceName") ?? ""),
        environment: String(form.get("environment") ?? ""),
        tags: String(form.get("tags") ?? ""),
        creator: String(form.get("creator") ?? ""),
        expiration,
        file,
      });
      onUploaded(response.artifact);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Upload failed");
    } finally {
      setSaving(false);
    }
  }

  return (
    <>
      <PageHeader
        title="Upload artifact"
        description="Add enough metadata that someone opening the link can start debugging without digging through a Slack thread."
        action={
          <Button type="button" variant="outline" onClick={onCancel}>
            <ArrowLeft className="h-4 w-4" />
            Back
          </Button>
        }
      />
      {error && <ErrorState message={error} />}
      <Card>
        <CardContent className="p-5">
          <form className="grid gap-5" onSubmit={handleSubmit}>
            <div className="grid gap-4 md:grid-cols-2">
              <Field label="Title">
                <Input name="title" required placeholder="Checkout API 500 in staging" />
              </Field>
              <Field label="Service name">
                <Input name="serviceName" required placeholder="payments" />
              </Field>
              <Field label="Artifact type">
                <Select value={artifactType} onChange={(event) => setArtifactType(event.target.value as ArtifactType)}>
                  <option value="stack_trace">Stack trace</option>
                  <option value="log">Log</option>
                  <option value="api_payload">API payload</option>
                  <option value="validation_report">Validation report</option>
                  <option value="screenshot">Screenshot</option>
                </Select>
              </Field>
              <Field label="Environment">
                <Input name="environment" required placeholder="staging" />
              </Field>
              <Field label="Creator">
                <Input name="creator" required placeholder="you@oracle.com" />
              </Field>
              <Field label="Expiration policy">
                <Select value={expiration} onChange={(event) => setExpiration(event.target.value as ExpirationPolicy)}>
                  <option value="7d">7 days</option>
                  <option value="14d">14 days</option>
                  <option value="never">Never</option>
                </Select>
              </Field>
            </div>
            <Field label="Description">
              <Textarea name="description" placeholder="Observed behavior, repro path, expected behavior, and any relevant deployment context." />
            </Field>
            <Field label="Tags">
              <Input name="tags" placeholder="checkout, sev2, qa" />
            </Field>
            <Field label="File upload">
              <label className="flex min-h-32 cursor-pointer flex-col items-center justify-center rounded-lg border border-dashed bg-muted/30 px-4 py-6 text-center hover:bg-muted/50">
                <UploadCloud className="mb-2 h-7 w-7 text-muted-foreground" />
                <span className="text-sm font-medium">{file ? file.name : "Choose artifact file"}</span>
                <span className="mt-1 text-xs text-muted-foreground">Logs, text, JSON, reports, and screenshots up to 25 MB</span>
                <input className="sr-only" type="file" onChange={(event) => setFile(event.target.files?.[0] ?? null)} />
              </label>
            </Field>
            <div className="flex justify-end gap-2">
              <Button type="button" variant="outline" onClick={onCancel}>
                Cancel
              </Button>
              <Button type="submit" disabled={saving}>
                <UploadCloud className="h-4 w-4" />
                {saving ? "Uploading" : "Create short link"}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </>
  );
}

function Field({ label, children }: { label: string; children: ReactNode }) {
  return (
    <div className="grid gap-2">
      <Label>{label}</Label>
      {children}
    </div>
  );
}
