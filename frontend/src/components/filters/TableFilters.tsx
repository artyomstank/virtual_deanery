interface Props {
  search: string;
  setSearch: (value: string) => void;
}

export function TableFilters({
  search,
  setSearch
}: Props) {
  return (
    <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
      <input
        placeholder="Search users..."
        value={search}
        onChange={(e) =>
          setSearch(e.target.value)
        }
        className="
          h-12 w-full rounded-xl border border-zinc-800
          bg-zinc-900 px-4 outline-none
          transition focus:border-indigo-500
          lg:max-w-sm
        "
      />

      <div className="flex gap-3">
        <select
          className="
            h-12 rounded-xl border border-zinc-800
            bg-zinc-900 px-4 outline-none
          "
        >
          <option>All roles</option>
          <option>Admin</option>
          <option>Teacher</option>
          <option>Student</option>
        </select>

        <select
          className="
            h-12 rounded-xl border border-zinc-800
            bg-zinc-900 px-4 outline-none
          "
        >
          <option>All statuses</option>
          <option>Active</option>
          <option>Blocked</option>
          <option>Pending</option>
        </select>
      </div>
    </div>
  );
}