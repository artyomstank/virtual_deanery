import { LucideIcon } from "lucide-react";

interface Props {
  title: string;
  value: number | string;
  icon: LucideIcon;
  description?: string;
}

export function StatCard({
  title,
  value,
  icon: Icon,
  description
}: Props) {
  return (
    <div
      className="
        group cursor-pointer rounded-2xl border border-zinc-800
        bg-zinc-900 p-6 transition-all duration-300
        hover:-translate-y-1
        hover:border-indigo-500/50
        hover:bg-zinc-900/80
      "
    >
      <div className="mb-5 flex items-start justify-between">
        <div>
          <p className="text-sm text-zinc-400">
            {title}
          </p>

          <h2 className="mt-2 text-4xl font-bold">
            {value}
          </h2>
        </div>

        <div
          className="
            flex h-14 w-14 items-center justify-center
            rounded-2xl bg-indigo-500/15 text-indigo-400
          "
        >
          <Icon size={24} />
        </div>
      </div>

      {description && (
        <p className="text-sm text-zinc-500">
          {description}
        </p>
      )}
    </div>
  );
}