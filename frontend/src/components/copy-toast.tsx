import { CheckCircle2 } from "lucide-react";

type CopyToastProps = {
  visible: boolean;
};

export function CopyToast({ visible }: CopyToastProps) {
  if (!visible) return null;

  return (
    <div className="fixed bottom-5 left-1/2 z-50 flex -translate-x-1/2 items-center gap-2 rounded-md border bg-card px-4 py-3 text-sm font-medium text-foreground shadow-lg">
      <CheckCircle2 className="h-4 w-4 text-emerald-600" />
      Copied to clipboard
    </div>
  );
}
