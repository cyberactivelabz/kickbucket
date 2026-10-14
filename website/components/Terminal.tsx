// Illustrative output. Replace with an asciinema recording once the CLI exists.
type Line = { text: string; tone?: "amber" | "coral" | "mint" | "mute" };

const lines: Line[] = [
  { text: "$ bucketmap scan acme", tone: "mute" },
  { text: "scanning 1,240 candidate names on s3, gcs, azure", tone: "mute" },
  { text: "s3     acme-backups    PUBLIC READ", tone: "coral" },
  { text: "azure  acme-logs       LISTABLE", tone: "amber" },
  { text: "gcs    acme-assets     private", tone: "mint" },
  { text: "3 buckets found, 2 exposed", tone: "mute" },
];

const toneClass = {
  amber: "text-amber",
  coral: "text-coral",
  mint: "text-mint",
  mute: "text-mute",
} as const;

export default function Terminal() {
  return (
    <pre
      aria-label="Example BucketMap scan output"
      className="overflow-x-auto rounded-md border border-line bg-panel p-5 font-mono text-sm leading-7"
    >
      {lines.map((l, i) => (
        <div
          key={l.text}
          className={`scan-line ${l.tone ? toneClass[l.tone] : ""}`}
          style={{ "--i": i } as React.CSSProperties}
        >
          {l.text}
        </div>
      ))}
    </pre>
  );
}
