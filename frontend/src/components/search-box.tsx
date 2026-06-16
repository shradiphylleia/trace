import { FormEvent, useState } from "react";
import { Search } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

type SearchBoxProps = {
  initialValue?: string;
  placeholder?: string;
  onSearch: (query: string) => void;
};

export function SearchBox({ initialValue = "", placeholder = "Search by title, service, tag, or error text", onSearch }: SearchBoxProps) {
  const [query, setQuery] = useState(initialValue);

  function handleSubmit(event: FormEvent) {
    event.preventDefault();
    onSearch(query.trim());
  }

  return (
    <form className="flex flex-col gap-2 sm:flex-row" onSubmit={handleSubmit}>
      <div className="relative flex-1">
        <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
        <Input className="pl-9" value={query} placeholder={placeholder} onChange={(event) => setQuery(event.target.value)} />
      </div>
      <Button type="submit">
        <Search className="h-4 w-4" />
        Search
      </Button>
    </form>
  );
}
