import { useEffect, useState } from "react";
import type { Artifact } from "@/api/types";
import { searchArtifacts } from "@/api/client";
import { PageHeader } from "@/components/page-header";
import { SearchBox } from "@/components/search-box";
import { EmptyState, ErrorState, LoadingState } from "@/components/state";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { formatArtifactType } from "@/lib/artifact";
import { formatDate } from "@/lib/date";

type SearchResultsPageProps = {
  initialQuery: string;
  onOpenArtifact: (shortCode: string) => void;
};

export function SearchResultsPage({ initialQuery, onOpenArtifact }: SearchResultsPageProps) {
  const [query, setQuery] = useState(initialQuery);
  const [artifacts, setArtifacts] = useState<Artifact[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    runSearch(initialQuery);
  }, [initialQuery]);

  async function runSearch(nextQuery: string) {
    setQuery(nextQuery);
    setLoading(true);
    setError("");
    try {
      const items = await searchArtifacts({ q: nextQuery, limit: 50 });
      setArtifacts(items);
    } catch (err) {
      setArtifacts([]);
      setError(err instanceof Error ? err.message : "Search failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <>
      <PageHeader title="Search results" description="Find debugging context by title, service, tags, or error text." />
      <Card>
        <CardContent className="p-5">
          <SearchBox initialValue={query} onSearch={runSearch} />
        </CardContent>
      </Card>
      {loading && <LoadingState label="Searching artifacts" />}
      {error && <ErrorState message={error} />}
      {!loading && !error && artifacts.length === 0 && (
        <EmptyState title="No matching artifacts" description="Try a service name, a shorter error phrase, or one of the tags used by the team." />
      )}
      {!loading && !error && artifacts.length > 0 && (
        <div className="grid gap-3">
          {artifacts.map((artifact) => (
            <Card key={artifact.id}>
              <CardContent className="grid gap-3 p-5 md:grid-cols-[minmax(0,1fr)_auto] md:items-center">
                <div className="min-w-0">
                  <div className="flex flex-wrap items-center gap-2">
                    <h3 className="font-semibold">{artifact.title}</h3>
                    <Badge variant="outline">{formatArtifactType(artifact.artifact_type)}</Badge>
                  </div>
                  <div className="mt-1 text-sm text-muted-foreground">
                    {artifact.service_name} / {formatDate(artifact.created_at)}
                  </div>
                  <p className="mt-2 line-clamp-2 text-sm text-muted-foreground">
                    {artifact.preview || artifact.description || "No preview available."}
                  </p>
                </div>
                <Button type="button" variant="outline" onClick={() => onOpenArtifact(artifact.short_code)}>
                  Open
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </>
  );
}
