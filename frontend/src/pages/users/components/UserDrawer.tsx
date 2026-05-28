import { X } from "lucide-react";

import { User } from "@/types/user";
import { StatusBadge } from "@/components/status/StatusBadge";

interface Props {
  user: User | null;
  onClose: () => void;
}

export function UserDrawer({
  user,
  onClose
}: Props) {
  if (!user) return null;

  return (
    <div
      className="
        fixed right-0 top-0 z-50 h-screen w-[420px]
        border-l border-zinc-800 bg-zinc-950
        p-6 shadow-2xl
      "
    >
      <div className="mb-8 flex items-center justify-between">
        <h2 className="text-2xl font-bold">
          User details
        </h2>

        <button
          onClick={onClose}
          className="
            flex h-10 w-10 items-center justify-center
            rounded-xl border border-zinc-800
            hover:bg-zinc-900
          "
        >
          <X size={18} />
        </button>
      </div>

      <div className="mb-8 flex items-center gap-4">
        <div
          className="
            flex h-16 w-16 items-center justify-center
            rounded-2xl bg-indigo-500 text-xl font-bold
          "
        >
          {user.name[0]}
        </div>

        <div>
          <h3 className="text-xl font-semibold">
            {user.name}
          </h3>

          <p className="text-zinc-400">
            {user.email}
          </p>
        </div>
      </div>

      <div className="space-y-5">
        <div>
          <p className="mb-2 text-sm text-zinc-500">
            Role
          </p>

          <div className="capitalize">
            {user.role}
          </div>
        </div>

        <div>
          <p className="mb-2 text-sm text-zinc-500">
            Status
          </p>

          <StatusBadge status={user.status} />
        </div>

        <div>
          <p className="mb-2 text-sm text-zinc-500">
            Registered
          </p>

          <div>{user.createdAt}</div>
        </div>
      </div>

      <div className="mt-10 space-y-3">
        <button
          className="
            h-12 w-full rounded-xl bg-red-500/15
            font-medium text-red-400
            transition hover:bg-red-500/25
          "
        >
          Block user
        </button>

        <button
          className="
            h-12 w-full rounded-xl border border-zinc-800
            bg-zinc-900 font-medium
            transition hover:border-indigo-500
          "
        >
          Reset password
        </button>
      </div>
    </div>
  );
}