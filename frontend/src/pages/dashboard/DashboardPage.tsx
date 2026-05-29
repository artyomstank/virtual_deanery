import { useEffect, useState } from "react";
import { Clock, GraduationCap, UserCheck, Users } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { PageContainer } from "../../components/layout/PageContainer";
import { Button } from "../../components/ui/Form";
import {
  Table,
  TableBody,
  TableCell,
  TableEmpty,
  TableHead,
  TableHeaderCell,
  TableRow,
} from "../../components/ui/Table";
import { apiClient } from "../../lib/api";
import { useNotificationStore } from "../../store/notification.store";
import { AdminStats, AuditLog } from "../../types/api";

function StatCard({
  title,
  value,
  icon,
  accent = "text-indigo-400",
  warning = false,
  onClick,
}: {
  title: string;
  value: number;
  icon: React.ReactNode;
  accent?: string;
  warning?: boolean;
  onClick?: () => void;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={`rounded-lg border p-6 text-left transition ${
        warning
          ? "border-amber-700 bg-amber-950/30 hover:bg-amber-950/45"
          : "border-zinc-800 bg-zinc-900 hover:bg-zinc-800/80"
      } ${onClick ? "cursor-pointer" : "cursor-default"}`}
    >
      <div className="flex items-center justify-between gap-4">
        <div>
          <p className="text-sm text-zinc-400">{title}</p>
          <p className={`mt-2 text-3xl font-bold ${warning ? "text-amber-300" : "text-white"}`}>
            {value}
          </p>
        </div>
        <div className={accent}>{icon}</div>
      </div>
    </button>
  );
}

export function DashboardPage() {
  const navigate = useNavigate();
  const { error: showError } = useNotificationStore();
  const [stats, setStats] = useState<AdminStats | null>(null);
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    void fetchData();
  }, []);

  const fetchData = async () => {
    try {
      setLoading(true);
      const [statsRes, auditRes] = await Promise.all([
        apiClient.get("/admin/stats"),
        apiClient.get("/admin/audit?limit=10"),
      ]);
      setStats(statsRes);
      setLogs(auditRes?.logs || []);
    } catch (err: any) {
      showError(err.message || "Не удалось загрузить дашборд");
    } finally {
      setLoading(false);
    }
  };

  return (
    <PageContainer title="Дашборд">
      <div className="mb-8 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard title="Всего пользователей" value={stats?.total_users || 0} icon={<Users size={32} />} />
        <StatCard title="Студентов" value={stats?.students || 0} icon={<GraduationCap size={32} />} accent="text-blue-400" />
        <StatCard title="Преподавателей" value={stats?.teachers || 0} icon={<UserCheck size={32} />} accent="text-emerald-400" />
        <StatCard
          title="Ожидают подтверждения"
          value={stats?.pending_approval || 0}
          icon={<Clock size={32} />}
          accent={(stats?.pending_approval || 0) > 0 ? "text-amber-400" : "text-zinc-500"}
          warning={(stats?.pending_approval || 0) > 0}
          onClick={() => navigate("/admin/users?tab=pending")}
        />
      </div>

      <div>
        <div className="mb-6 flex items-center justify-between gap-4">
          <h2 className="text-lg font-semibold text-white">Последние действия</h2>
          <Button variant="secondary" size="sm" onClick={() => navigate("/admin/reports")}>
            Смотреть все
          </Button>
        </div>

        <Table>
          <TableHead>
            <TableRow isHoverable={false}>
              <TableHeaderCell>Пользователь</TableHeaderCell>
              <TableHeaderCell>Действие</TableHeaderCell>
              <TableHeaderCell>Дата/время</TableHeaderCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {loading ? (
              <TableRow isHoverable={false}>
                <TableCell colSpan={3} className="py-8 text-center">
                  <div className="flex items-center justify-center gap-2 text-zinc-400">
                    <div className="h-4 w-4 animate-spin rounded-full border-2 border-indigo-500 border-t-transparent" />
                    Загрузка...
                  </div>
                </TableCell>
              </TableRow>
            ) : logs.length === 0 ? (
              <TableEmpty message="Данные не найдены" />
            ) : (
              logs.slice(0, 10).map((log) => (
                <TableRow key={log.id}>
                  <TableCell>{log.user_name}</TableCell>
                  <TableCell>{log.action}</TableCell>
                  <TableCell>{new Date(log.timestamp).toLocaleString("ru-RU")}</TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
    </PageContainer>
  );
}
