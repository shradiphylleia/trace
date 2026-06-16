import { Plus } from "lucide-react";
import type { Artifact } from "@/api/types";
import { ArtifactTable } from "@/components/artifact-table";
import { PageHeader } from "@/components/page-header";
import { SearchBox } from "@/components/search-box";
import { EmptyState, LoadingState } from "@/components/state";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

type DashboardPageProps = {
  artifacts: Artifact[];
  loading: boolean;
  onSearch: (query: string) => void;
  onOpenArtifact: (shortCode: string) => void;
  onUpload: () => void;
};

export function DashboardPage({ artifacts, loading, onSearch, onOpenArtifact, onUpload }: DashboardPageProps) {
  return (
    <>
      <PageHeader
        title="TraceShare"
        description="Collect logs, stack traces, payloads, validation reports, and screenshots into one short debugging link."
        action={
          <Button type="button" onClick={onUpload}>
            <Plus className="h-4 w-4" />
            Upload artifact
          </Button>
        }
      />

      <Card>
        <CardHeader>
          <CardTitle>Search artifacts</CardTitle>
        </CardHeader>
        <CardContent>
          <SearchBox onSearch={onSearch} />
        </CardContent>
      </Card>

      <section className="grid gap-3">
        <div className="flex items-center justify-between">
          <h3 className="text-base font-semibold">Recent artifacts</h3>
          <span className="text-sm text-muted-foreground">{artifacts.length} shown</span>
        </div>
        {loading && <LoadingState label="Loading recent artifacts" />}
        {!loading && artifacts.length === 0 && (
          <EmptyState title="No artifacts yet" description="Upload the first trace, log, payload, report, or screenshot to create a shareable link." />
        )}
        {!loading && artifacts.length > 0 && <ArtifactTable artifacts={artifacts} onOpenArtifact={onOpenArtifact} />}
      </section>
    </>
  );
}
