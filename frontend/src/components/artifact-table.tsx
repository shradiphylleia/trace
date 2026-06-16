import type { Artifact } from "@/api/types";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { formatArtifactType } from "@/lib/artifact";
import { formatDate } from "@/lib/date";

type ArtifactTableProps = {
  artifacts: Artifact[];
  onOpenArtifact: (shortCode: string) => void;
};

export function ArtifactTable({ artifacts, onOpenArtifact }: ArtifactTableProps) {
  return (
    <Card>
      <CardContent className="overflow-x-auto p-0">
        <table className="w-full min-w-[760px] text-left text-sm">
          <thead className="border-b bg-muted/50 text-xs uppercase text-muted-foreground">
            <tr>
              <th className="px-5 py-3 font-semibold">Artifact</th>
              <th className="px-5 py-3 font-semibold">Type</th>
              <th className="px-5 py-3 font-semibold">Service</th>
              <th className="px-5 py-3 font-semibold">Created</th>
              <th className="px-5 py-3 font-semibold">Tags</th>
              <th className="px-5 py-3 text-right font-semibold">Action</th>
            </tr>
          </thead>
          <tbody>
            {artifacts.map((artifact) => (
              <tr key={artifact.id} className="border-b last:border-b-0">
                <td className="px-5 py-4">
                  <div className="font-medium text-foreground">{artifact.title}</div>
                  <div className="mt-1 max-w-md truncate text-xs text-muted-foreground">{artifact.description || artifact.file_name}</div>
                </td>
                <td className="px-5 py-4">
                  <Badge variant="outline">{formatArtifactType(artifact.artifact_type)}</Badge>
                </td>
                <td className="px-5 py-4">
                  <div>{artifact.service_name}</div>
                  <div className="text-xs text-muted-foreground">{artifact.environment}</div>
                </td>
                <td className="px-5 py-4 text-muted-foreground">{formatDate(artifact.created_at)}</td>
                <td className="px-5 py-4">
                  <div className="flex max-w-48 flex-wrap gap-1">
                    {artifact.tags.slice(0, 3).map((tag) => (
                      <Badge key={tag} variant="default">
                        {tag}
                      </Badge>
                    ))}
                  </div>
                </td>
                <td className="px-5 py-4 text-right">
                  <Button type="button" size="sm" variant="outline" onClick={() => onOpenArtifact(artifact.short_code)}>
                    Open
                  </Button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </CardContent>
    </Card>
  );
}
