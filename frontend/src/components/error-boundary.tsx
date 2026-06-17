import { Component, type ErrorInfo, type ReactNode } from "react";
import { AlertCircle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";

type ErrorBoundaryProps = {
  children: ReactNode;
};

type ErrorBoundaryState = {
  hasError: boolean;
};

export class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  state: ErrorBoundaryState = { hasError: false };

  static getDerivedStateFromError() {
    return { hasError: true };
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    console.error("TraceShare UI error", error, info);
  }

  render() {
    if (!this.state.hasError) {
      return this.props.children;
    }

    return (
      <div className="min-h-screen bg-background p-6">
        <Card className="mx-auto mt-16 max-w-xl border-destructive/30">
          <CardContent className="grid gap-4 p-6 text-center">
            <AlertCircle className="mx-auto h-8 w-8 text-destructive" />
            <div>
              <h1 className="text-lg font-semibold">TraceShare hit a UI error</h1>
              <p className="mt-1 text-sm text-muted-foreground">Refresh the page or check whether the API response shape changed.</p>
            </div>
            <Button type="button" onClick={() => window.location.reload()}>
              Refresh
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }
}
