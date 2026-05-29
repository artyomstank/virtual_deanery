import { useEffect, useState } from "react";
import { Check } from "lucide-react";
import { PageContainer } from "../../components/layout/PageContainer";
import { apiClient } from "../../lib/api";
import { useNotificationStore } from "../../store/notification.store";
import { ACLPermission } from "../../types/api";

const resources = ["users", "students", "grades", "attendance", "exam_schedule", "debts", "reports", "profile"];
const roles = ["student", "teacher", "dean", "admin"];

const resourceLabels: Record<string, string> = {
  users: "users",
  students: "students",
  grades: "grades",
  attendance: "attendance",
  exam_schedule: "exam_schedule",
  debts: "debts",
  reports: "reports",
  profile: "profile",
};

const roleLabels: Record<string, string> = {
  student: "студент",
  teacher: "преподаватель",
  dean: "декан",
  admin: "админ",
};

export function ACLPage() {
  const { error } = useNotificationStore();
  const [permissions, setPermissions] = useState<ACLPermission[]>([]);
  const [loading, setLoading] = useState(true);
  const [savedKey, setSavedKey] = useState("");

  useEffect(() => {
    void loadPermissions();
  }, []);

  const loadPermissions = async () => {
    try {
      setLoading(true);
      const response = await apiClient.get("/admin/acl");
      setPermissions(response?.permissions || []);
    } catch (err: any) {
      error(err.message || "Не удалось загрузить права доступа");
    } finally {
      setLoading(false);
    }
  };

  const getPermission = (resource: string, role: string) =>
    permissions.find((permission) => permission.resource === resource && permission.role === role) || {
      resource,
      role,
      read: false,
      write: false,
      delete: false,
    };

  const patchPermission = async (resource: string, role: string, field: "read" | "write" | "delete") => {
    const current = getPermission(resource, role);
    const nextValue = !current[field];
    const key = `${resource}-${role}-${field}`;
    setPermissions((items) => {
      const exists = items.some((item) => item.resource === resource && item.role === role);
      if (!exists) return [...items, { ...current, [field]: nextValue }];
      return items.map((item) => item.resource === resource && item.role === role ? { ...item, [field]: nextValue } : item);
    });

    try {
      await apiClient.patch("/admin/acl", { resource, role, permission: field, value: nextValue });
      setSavedKey(key);
      window.setTimeout(() => setSavedKey((value) => value === key ? "" : value), 2000);
    } catch (err: any) {
      setPermissions((items) => items.map((item) => item.resource === resource && item.role === role ? { ...item, [field]: current[field] } : item));
      error(err.message || "Не удалось обновить право");
    }
  };

  return (
    <PageContainer title="Матрица прав доступа">
      <div className="grid grid-cols-1 gap-6 xl:grid-cols-[1fr_280px]">
        <div className="overflow-x-auto rounded-lg border border-zinc-800">
          <table className="w-full min-w-[920px]">
            <thead>
              <tr className="border-b border-zinc-800 bg-zinc-900">
                <th className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wider text-zinc-400">Роль / ресурс</th>
                {resources.map((resource) => (
                  <th key={resource} className="px-4 py-3 text-center text-xs font-semibold uppercase tracking-wider text-zinc-400">
                    {resourceLabels[resource]}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-800">
              {loading ? (
                <tr><td colSpan={resources.length + 1} className="px-4 py-8 text-center text-zinc-400">Загрузка...</td></tr>
              ) : (
                roles.map((role) => (
                  <tr key={role} className="bg-zinc-950 hover:bg-zinc-900">
                    <td className="px-4 py-3 text-sm font-medium text-white">{roleLabels[role]}</td>
                    {resources.map((resource) => {
                      const permission = getPermission(resource, role);
                      return (
                        <td key={`${role}-${resource}`} className="px-4 py-3">
                          <div className="flex items-center justify-center gap-2">
                            {(["read", "write", "delete"] as const).map((field) => {
                              const label = field === "read" ? "R" : field === "write" ? "W" : "D";
                              const key = `${resource}-${role}-${field}`;
                              return (
                                <label key={field} className="flex items-center gap-1 text-xs text-zinc-400">
                                  <input
                                    type="checkbox"
                                    checked={permission[field]}
                                    onChange={() => patchPermission(resource, role, field)}
                                    className="h-4 w-4 rounded border-zinc-700 bg-zinc-800 accent-indigo-600"
                                  />
                                  {label}
                                  {savedKey === key && <Check size={14} className="text-emerald-400" />}
                                </label>
                              );
                            })}
                          </div>
                        </td>
                      );
                    })}
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        <aside className="rounded-lg border border-zinc-800 bg-zinc-900 p-4">
          <h2 className="mb-3 text-sm font-semibold text-white">Легенда</h2>
          <div className="space-y-3 text-sm text-zinc-400">
            <p><span className="font-semibold text-white">R</span> — чтение и просмотр ресурса.</p>
            <p><span className="font-semibold text-white">W</span> — создание и редактирование.</p>
            <p><span className="font-semibold text-white">D</span> — удаление записей.</p>
            <p>Ресурсы соответствуют API-модулям: пользователи, студенты, оценки, посещаемость, расписание экзаменов, долги, отчёты и профиль.</p>
          </div>
        </aside>
      </div>
    </PageContainer>
  );
}
