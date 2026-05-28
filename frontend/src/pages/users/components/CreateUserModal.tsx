import * as Dialog from "@radix-ui/react-dialog";

import { X } from "lucide-react";

interface Props {
  open: boolean;
  onOpenChange: (value: boolean) => void;
}

export function CreateUserModal({
  open,
  onOpenChange
}: Props) {
  return (
    <Dialog.Root
      open={open}
      onOpenChange={onOpenChange}
    >
      <Dialog.Portal>
        <Dialog.Overlay
          className="
            fixed inset-0 z-40 bg-black/70
            backdrop-blur-sm
          "
        />

        <Dialog.Content
          className="
            fixed left-1/2 top-1/2 z-50
            w-full max-w-lg
            -translate-x-1/2 -translate-y-1/2
            rounded-2xl border border-zinc-800
            bg-zinc-950 p-6 shadow-2xl
          "
        >
          <div className="mb-8 flex items-center justify-between">
            <div>
              <Dialog.Title className="text-2xl font-bold">
                Create user
              </Dialog.Title>

              <Dialog.Description className="mt-2 text-sm text-zinc-400">
                New users created here become active immediately
              </Dialog.Description>
            </div>

            <button
              onClick={() =>
                onOpenChange(false)
              }
              className="
                flex h-10 w-10 items-center justify-center
                rounded-xl border border-zinc-800
                hover:bg-zinc-900
              "
            >
              <X size={18} />
            </button>
          </div>

          <div className="space-y-5">
            <div>
              <label className="mb-2 block text-sm text-zinc-400">
                Full name
              </label>

              <input
                placeholder="John Doe"
                className="
                  h-12 w-full rounded-xl border
                  border-zinc-800 bg-zinc-900
                  px-4 outline-none
                  transition focus:border-indigo-500
                "
              />
            </div>

            <div>
              <label className="mb-2 block text-sm text-zinc-400">
                Email
              </label>

              <input
                placeholder="john@uni.edu"
                className="
                  h-12 w-full rounded-xl border
                  border-zinc-800 bg-zinc-900
                  px-4 outline-none
                  transition focus:border-indigo-500
                "
              />
            </div>

            <div>
              <label className="mb-2 block text-sm text-zinc-400">
                Password
              </label>

              <input
                type="password"
                placeholder="••••••••"
                className="
                  h-12 w-full rounded-xl border
                  border-zinc-800 bg-zinc-900
                  px-4 outline-none
                  transition focus:border-indigo-500
                "
              />
            </div>

            <div>
              <label className="mb-2 block text-sm text-zinc-400">
                Role
              </label>

              <select
                className="
                  h-12 w-full rounded-xl border
                  border-zinc-800 bg-zinc-900
                  px-4 outline-none
                "
              >
                <option>student</option>
                <option>teacher</option>
                <option>admin</option>
              </select>
            </div>
          </div>

          <div className="mt-8 flex justify-end gap-3">
            <button
              onClick={() =>
                onOpenChange(false)
              }
              className="
                h-12 rounded-xl border border-zinc-800
                px-5 font-medium transition
                hover:bg-zinc-900
              "
            >
              Cancel
            </button>

            <button
              className="
                h-12 rounded-xl bg-indigo-500
                px-5 font-medium transition
                hover:bg-indigo-400
              "
            >
              Create user
            </button>
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  );
}