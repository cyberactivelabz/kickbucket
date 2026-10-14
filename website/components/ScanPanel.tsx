"use client";

import { useMemo, useState } from "react";

const words = ["dev", "prod", "staging", "backup", "backups", "assets", "logs", "data", "public", "private", "uploads", "media"];
const providerNames = ["s3", "gcs", "azure"] as const;

export default function ScanPanel() {
  const [keyword, setKeyword] = useState("acme");
  const [providers, setProviders] = useState<string[]>([...providerNames]);
  const candidates = useMemo(() => {
    const base = keyword.trim().toLowerCase();
    if (!base) return [];
    return Array.from(new Set([base, ...words.flatMap((word) => [`${base}-${word}`, `${word}-${base}`])]));
  }, [keyword]);

  function toggleProvider(provider: string) {
    setProviders((current) => current.includes(provider)
      ? current.filter((name) => name !== provider)
      : [...current, provider]);
  }

  return (
    <section aria-labelledby="scan-planner" className="border-y border-line bg-panel/60 py-8">
      <div className="grid gap-8 lg:grid-cols-[1fr_1.25fr] lg:items-start">
        <div>
          <p className="font-mono text-xs uppercase tracking-[0.18em] text-amber">Scan planner</p>
          <h2 id="scan-planner" className="mt-3 text-2xl font-semibold">Shape a scan before you run it.</h2>
          <p className="mt-3 max-w-md text-mute">Preview candidate names locally. The planner never sends a request or reads cloud storage.</p>
          <label className="mt-6 block font-mono text-sm" htmlFor="keyword">Keyword</label>
          <input
            id="keyword"
            value={keyword}
            onChange={(event) => setKeyword(event.target.value)}
            className="mt-2 w-full rounded-md border border-line bg-ink px-3 py-3 font-mono text-text outline-none focus:border-amber"
            placeholder="acme"
          />
          <fieldset className="mt-5">
            <legend className="font-mono text-sm">Providers</legend>
            <div className="mt-3 flex flex-wrap gap-2">
              {providerNames.map((provider) => {
                const active = providers.includes(provider);
                return (
                  <button
                    key={provider}
                    type="button"
                    aria-pressed={active}
                    onClick={() => toggleProvider(provider)}
                    className={`rounded-md border px-3 py-2 font-mono text-sm uppercase transition-colors ${active ? "border-amber bg-amber text-ink" : "border-line text-mute hover:border-text hover:text-text"}`}
                  >
                    {provider}
                  </button>
                );
              })}
            </div>
          </fieldset>
        </div>
        <div className="min-w-0 rounded-md border border-line bg-ink p-5">
          <div className="flex items-center justify-between gap-4 font-mono text-xs text-mute">
            <span>{providers.length} providers selected</span>
            <span>{candidates.length} candidates</span>
          </div>
          <div className="mt-4 grid max-h-64 gap-x-6 gap-y-2 overflow-y-auto sm:grid-cols-2">
            {candidates.map((candidate) => <code key={candidate} className="truncate text-sm text-mint">{candidate}</code>)}
            {candidates.length === 0 && <span className="font-mono text-sm text-mute">Enter a keyword to preview candidates.</span>}
          </div>
        </div>
      </div>
    </section>
  );
}