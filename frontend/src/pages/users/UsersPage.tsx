import { useMemo, useState } from "react";

import {
  Check,
  Plus,
  X
} from "lucide-react";

import { users } from "@/mocks/users";

import { User } from "@/types/user";

import { PageContainer } from "@/components/layout/PageContainer";
import { TableFilters } from "@/components/filters/TableFilters";
import { StatusBadge } from "@/components/status/StatusBadge";


import { UserDrawer } from "./components/UserDrawer";
import { CreateUserModal } from "./components/CreateUserModal";

export function UsersPage() {
  const [search, setSearch] = useState("");

  const [selectedUser, setSelectedUser] =
    useState<User | null>(null);

  const filteredUsers = useMemo(() => {
    return users.filter((user) => {
      return (
        user.name
          .toLowerCase()
          .includes(search.toLowerCase()) ||
        user.email
          .toLowerCase()
          .includes(search.toLowerCase())
      );
    });
  }, [search]);

  return (
    <>
      <PageContainer
        title="Users"
        description="Manage users, permissions and account statuses"
        action={
          <button
            className="
              flex h-12 items-center gap-2 rounded-xl
              bg-indigo-500 px-5 font-medium
              transition hover:bg-indigo-400
            "
          >
            <Plus size={18} />

            Create user
          </button>
        }
      >
        <div className="space-y-6">
          <TableFilters
            search={search}
            setSearch={setSearch}
          />

          <div
            className="
              overflow-hidden rounded-2xl border
              border-zinc-800 bg-zinc-900
            "
          >
            <table className="w-full">
              <thead className="bg-zinc-950">
                <tr className="text-left text-sm text-zinc-400">
                  <th className="px-6 py-4">
                    Name
                  </th>

                  <th className="px-6 py-4">
                    Email
                  </th>

                  <th className="px-6 py-4">
                    Role
                  </th>

                  <th className="px-6 py-4">
                    Status
                  </th>

                  <th className="px-6 py-4">
                    Registered
                  </th>
                </tr>
              </thead>

              <tbody>
                {filteredUsers.map((user) => (
                  <tr
                    key={user.id}
                    onClick={() =>
                      setSelectedUser(user)
                    }
                    className="
                      cursor-pointer border-t border-zinc-800
                      transition hover:bg-zinc-800/50
                    "
                  >
                    <td className="px-6 py-5 font-medium">
                      {user.name}
                    </td>

                    <td className="px-6 py-5 text-zinc-400">
                      {user.email}
                    </td>

                    <td className="px-6 py-5 capitalize">
                      {user.role}
                    </td>

                    <td className="px-6 py-5">
                      <StatusBadge
                        status={user.status}
                      />
                    </td>

                    <td className="px-6 py-5 text-zinc-400">
                      {user.createdAt}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </PageContainer>

      <UserDrawer
        user={selectedUser}
        onClose={() =>
          setSelectedUser(null)
        }
      />
    </>
  );
}