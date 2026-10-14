import Terminal from "@/components/Terminal";
import ScanPanel from "@/components/ScanPanel";

const REPO = "cyberactivelabz/kickbucket";
const GITHUB = `https://github.com/${REPO}`;

type Release = { assets?: { download_count: number }[] };

async function getStats() {
  try {
    const opts = { next: { revalidate: 3600 } };
    const [repo, releases] = await Promise.all([
      fetch(`https://api.github.com/repos/${REPO}`, opts).then((r) => r.json()),
      fetch(`https://api.github.com/repos/${REPO}/releases`, opts).then((r) => r.json()),
    ]);
    const downloads = Array.isArray(releases)
      ? (releases as Release[])
          .flatMap((r) => r.assets ?? [])
          .reduce((sum, a) => sum + a.download_count, 0)
      : 0;
    return { stars: (repo.stargazers_count as number) ?? 0, downloads };
  } catch {
    return { stars: 0, downloads: 0 };
  }
}

const features = [
  { title: "Fast", body: "A concurrent Go engine checks thousands of candidate names in seconds." },
  { title: "Multi-cloud", body: "S3, GCS and Azure today. New providers are a single package to add." },
  { title: "CLI, web and desktop", body: "Same engine everywhere, plus a Go library for your own tools." },
];

export default async function Home() {
  const { stars, downloads } = await getStats();

  return (
    <main className="mx-auto max-w-5xl px-6">
      <header className="flex items-center justify-between py-6 font-mono text-sm">
        <span className="font-medium">kickbucket</span>
        <nav className="flex gap-6 text-mute">
          <a className="hover:text-text" href="https://docs.kickbucket.dev">Docs</a>
          <a className="hover:text-text" href={GITHUB}>GitHub</a>
          <a className="hover:text-text" href={`https://github.com/sponsors/${REPO.split("/")[0]}`}>Sponsor</a>
        </nav>
      </header>

      <section className="py-16">
        <h1 className="max-w-3xl text-4xl font-semibold leading-tight sm:text-6xl">
          The Nmap for cloud storage.
        </h1>
        <p className="mt-5 max-w-xl text-lg text-mute">
          Find exposed S3, GCS and Azure buckets before someone else does. Open source, written in Go.
        </p>
        <code className="mt-8 block w-fit max-w-full select-all overflow-x-auto rounded-md border border-line bg-panel px-4 py-3 font-mono text-sm text-amber">
          go install github.com/{REPO}/cmd/kickbucket@latest
        </code>
        <div className="mt-10"><Terminal /></div>
      </section>

      <ScanPanel />

      <section className="grid gap-10 border-t border-line py-14 sm:grid-cols-3">
        {features.map((f) => (
          <div key={f.title} className="border-l-2 border-amber pl-4">
            <h2 className="font-semibold">{f.title}</h2>
            <p className="mt-2 text-mute">{f.body}</p>
          </div>
        ))}
      </section>

      <section className="flex flex-wrap gap-x-12 gap-y-6 border-t border-line py-14 font-mono">
        <div><div className="text-3xl">{stars.toLocaleString()}</div><div className="text-sm text-mute">GitHub stars</div></div>
        <div><div className="text-3xl">{downloads.toLocaleString()}</div><div className="text-sm text-mute">release downloads</div></div>
      </section>

      <section className="flex flex-wrap gap-4 border-t border-line py-14">
        <a href={GITHUB} className="rounded-md bg-amber px-5 py-3 font-semibold text-ink">View on GitHub</a>
        <a href={`${GITHUB}/releases`} className="rounded-md border border-line px-5 py-3">Download</a>
        <a href="https://docs.kickbucket.dev" className="rounded-md border border-line px-5 py-3">Read the docs</a>
      </section>

      <footer className="border-t border-line py-8 text-sm text-mute">
        MIT licensed. Use only on assets you are authorized to test.
      </footer>
    </main>
  );
}
