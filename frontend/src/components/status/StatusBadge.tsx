import { UserStatus } from "@/types/user";

interface Props {
  status: UserStatus;
}

const styles = {
  active:
    "bg-emerald-500/15 text-emerald-400 border-emerald-500/20",

  blocked:
    "bg-red-500/15 text-red-400 border-red-500/20",

  pending:
    "bg-amber-500/15 text-amber-400 border-amber-500/20"
};

export function StatusBadge({
  status
}: Props) {
  return (
    <div
      className={`
        inline-flex items-center rounded-full border
        px-2.5 py-1 text-xs font-medium capitalize
        ${styles[status]}
      `}
    >
      {status}
    </div>
  );
}