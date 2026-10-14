type Line = { text: string; tone?: "amber" | "coral" | "mint" | "mute" };

const lines: Line[] = [
  { text: "$ kickbucket scan acme", tone: "mute" },
  { text: "Only scan cloud storage you own or are authorized to test.", tone: "mute" },
  { text: "s3    acme-backups                   LISTABLE", tone: "amber" },
  { text: "azure acme-logs                      LISTABLE", tone: "amber" },
  { text: "gcs   acme-assets                    PRIVATE", tone: "mint" },
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
      aria-label="Example Kickbucket scan output"
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
