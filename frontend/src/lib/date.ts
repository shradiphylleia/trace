export function formatDate(value?: string | null) {
  if (!value) {
    return "Never";
  }
  return new Intl.DateTimeFormat(undefined, {
    year: "numeric",
    month: "short",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  }).format(new Date(value));
}

export function expirationLabel(value?: string | null) {
  if (!value) {
    return { text: "Never expires", tone: "success" as const };
  }
  const expiresAt = new Date(value).getTime();
  const now = Date.now();
  if (expiresAt <= now) {
    return { text: "Expired", tone: "danger" as const };
  }
  const days = Math.ceil((expiresAt - now) / (1000 * 60 * 60 * 24));
  if (days <= 2) {
    return { text: `Expires in ${days} day${days === 1 ? "" : "s"}`, tone: "warning" as const };
  }
  return { text: `Expires in ${days} day${days === 1 ? "" : "s"}`, tone: "info" as const };
}
