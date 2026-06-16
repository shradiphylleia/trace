import type { ReactNode } from "react";
import type { LucideIcon } from "lucide-react";
import { Bug, Server } from "lucide-react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

type NavRoute = {
  name: string;
  query?: string;
  shortCode?: string;
};

type NavItem = {
  label: string;
  icon: LucideIcon;
  route: NavRoute;
};

type AppShellProps = {
  activePage: string;
  navItems: NavItem[];
  onNavigate: (route: NavRoute) => void;
  children: ReactNode;
};

export function AppShell({ activePage, navItems, onNavigate, children }: AppShellProps) {
  return (
    <div className="page-shell">
      <div className="border-b bg-card">
        <div className="mx-auto flex max-w-7xl flex-col gap-4 px-4 py-4 sm:px-6 lg:flex-row lg:items-center lg:justify-between lg:px-8">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-md bg-primary text-primary-foreground">
              <Bug className="h-5 w-5" />
            </div>
            <div>
              <h1 className="text-lg font-semibold tracking-normal">TraceShare</h1>
              <p className="text-sm text-muted-foreground">Debugging artifacts for QA and engineering handoffs</p>
            </div>
          </div>
          <div className="flex flex-wrap items-center gap-2">
            {navItems.map((item) => {
              const Icon = item.icon;
              const page = typeof item.route.name === "string" ? item.route.name : "";
              return (
                <Button
                  key={item.label}
                  type="button"
                  variant={activePage === page ? "secondary" : "ghost"}
                  size="sm"
                  onClick={() => onNavigate(item.route)}
                >
                  <Icon className="h-4 w-4" />
                  {item.label}
                </Button>
              );
            })}
          </div>
        </div>
      </div>
      <main className="content-shell">
        <div className="flex items-center gap-2 text-xs font-medium uppercase tracking-wide text-muted-foreground">
          <Server className="h-3.5 w-3.5" />
          Internal engineering tool
        </div>
        <div className={cn("grid gap-6")}>{children}</div>
      </main>
      <footer className="border-t bg-card/60">
        <div className="mx-auto flex max-w-7xl items-center justify-center px-4 py-4 text-xs text-muted-foreground sm:px-6 lg:px-8">
          Made with love by Shraddha
        </div>
      </footer>
    </div>
  );
}
