import { AlertCircle, Inbox, Loader2 } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";

export function LoadingState({ label = "Loading artifacts" }: { label?: string }) {
  return (
    <Card>
      <CardContent className="flex items-center gap-3 p-5 text-sm text-muted-foreground">
        <Loader2 className="h-4 w-4 animate-spin" />
        {label}
      </CardContent>
    </Card>
  );
}

export function EmptyState({ title, description }: { title: string; description: string }) {
  return (
    <Card>
      <CardContent className="flex flex-col items-center justify-center gap-2 p-8 text-center">
        <Inbox className="h-8 w-8 text-muted-foreground" />
        <div className="font-medium">{title}</div>
        <p className="max-w-md text-sm text-muted-foreground">{description}</p>
      </CardContent>
    </Card>
  );
}

export function ErrorState({ message }: { message: string }) {
  return (
    <Card className="border-destructive/30">
      <CardContent className="flex items-center gap-3 p-5 text-sm text-destructive">
        <AlertCircle className="h-4 w-4" />
        {message}
      </CardContent>
    </Card>
  );
}
