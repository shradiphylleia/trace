import { useEffect, useMemo, useState } from "react";
import { Activity, FileSearch, Plus, Search } from "lucide-react";
import { getArtifact, searchArtifacts } from "@/api/client";
import type { Artifact } from "@/api/types";
import { AppShell } from "@/components/app-shell";
import { DashboardPage } from "@/pages/dashboard";
import { UploadPage } from "@/pages/upload";
import { ArtifactDetailsPage } from "@/pages/artifact-details";
import { SearchResultsPage } from "@/pages/search-results";

type Route =
  | { name: "dashboard" }
  | { name: "upload" }
  | { name: "artifact"; shortCode: string }
  | { name: "search"; query: string };

const navItems = [
  { label: "Dashboard", icon: Activity, route: { name: "dashboard" } as Route },
  { label: "Upload", icon: Plus, route: { name: "upload" } as Route },
  { label: "Search", icon: Search, route: { name: "search", query: "" } as Route },
  { label: "Lookup", icon: FileSearch, route: { name: "artifact", shortCode: "" } as Route },
];

function routeFromLocation(): Route {
  const url = new URL(window.location.href);
  const page = url.searchParams.get("page");
  if (page === "upload") return { name: "upload" };
  if (page === "artifact") return { name: "artifact", shortCode: url.searchParams.get("code") ?? "" };
  if (page === "search") return { name: "search", query: url.searchParams.get("q") ?? "" };
  return { name: "dashboard" };
}

function writeRoute(route: Route) {
  const url = new URL(window.location.href);
  url.search = "";
  if (route.name !== "dashboard") {
    url.searchParams.set("page", route.name);
  }
  if (route.name === "artifact" && route.shortCode) {
    url.searchParams.set("code", route.shortCode);
  }
  if (route.name === "search" && route.query) {
    url.searchParams.set("q", route.query);
  }
  window.history.pushState({}, "", url);
}

export default function App() {
  const [route, setRouteState] = useState<Route>(() => routeFromLocation());
  const [recentArtifacts, setRecentArtifacts] = useState<Artifact[]>([]);
  const [recentLoading, setRecentLoading] = useState(true);

  useEffect(() => {
    function onPopState() {
      setRouteState(routeFromLocation());
    }
    window.addEventListener("popstate", onPopState);
    return () => window.removeEventListener("popstate", onPopState);
  }, []);

  useEffect(() => {
    let ignore = false;
    setRecentLoading(true);
    searchArtifacts({ limit: 8 })
      .then((items) => {
        if (!ignore) setRecentArtifacts(items??[]);
      })
      .catch(() => {
        if (!ignore) setRecentArtifacts([]);
      })
      .finally(() => {
        if (!ignore) setRecentLoading(false);
      });
    return () => {
      ignore = true;
    };
  }, []);

  const activePage = useMemo(() => route.name, [route.name]);

  function setRoute(nextRoute: Route) {
    writeRoute(nextRoute);
    setRouteState(nextRoute);
  }

  async function handleArtifactUploaded(artifact: Artifact) {
    setRecentArtifacts((items) => [artifact, ...items.filter((item) => item.id !== artifact.id)].slice(0, 8));
    setRoute({ name: "artifact", shortCode: artifact.short_code });
  }

  return (
    <AppShell activePage={activePage} navItems={navItems} onNavigate={(nextRoute) => setRoute(nextRoute as Route)}>
      {route.name === "dashboard" && (
        <DashboardPage
          artifacts={recentArtifacts}
          loading={recentLoading}
          onSearch={(query) => setRoute({ name: "search", query })}
          onOpenArtifact={(shortCode) => setRoute({ name: "artifact", shortCode })}
          onUpload={() => setRoute({ name: "upload" })}
        />
      )}
      {route.name === "upload" && <UploadPage onUploaded={handleArtifactUploaded} onCancel={() => setRoute({ name: "dashboard" })} />}
      {route.name === "artifact" && (
        <ArtifactDetailsPage
          shortCode={route.shortCode}
          loadArtifact={getArtifact}
          onLookup={(shortCode) => setRoute({ name: "artifact", shortCode })}
          onSearch={(query) => setRoute({ name: "search", query })}
        />
      )}
      {route.name === "search" && (
        <SearchResultsPage
          initialQuery={route.query}
          onOpenArtifact={(shortCode) => setRoute({ name: "artifact", shortCode })}
        />
      )}
    </AppShell>
  );
}
