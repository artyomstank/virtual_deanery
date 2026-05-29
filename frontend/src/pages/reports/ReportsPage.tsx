import { useEffect, useState } from "react";
import { Download } from "lucide-react";
import { PageContainer } from "../../components/layout/PageContainer";
import { Button, Input } from "../../components/ui/Form";
import { Select, SelectItem } from "../../components/ui/Select";
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
import { AttendanceReport, AuditLog, GradeReport } from "../../types/api";

export function ReportsPage() {
  const { success, error } = useNotificationStore();
  const [grades, setGrades] = useState<GradeReport[]>([]);
  const [attendance, setAttendance] = useState<AttendanceReport[]>([]);
  const [audit, setAudit] = useState<AuditLog[]>([]);
  const [loading, setLoading] = useState(false);
  const [filters, setFilters] = useState({ group: "all", discipline: "all", period: "" });

  useEffect(() => {
    void loadReports();
  }, [filters.group, filters.discipline, filters.period]);

  const query = new URLSearchParams(filters).toString();

  const loadReports = async () => {
    try {
      setLoading(true);
      const [gradesRes, attendanceRes, auditRes] = await Promise.all([
        apiClient.get(`/reports/grades?${query}`),
        apiClient.get(`/reports/attendance?${query}`),
        apiClient.get(`/admin/audit?${query}`),
      ]);
      setGrades(gradesRes?.reports || []);
      setAttendance(attendanceRes?.reports || []);
      setAudit(auditRes?.logs || []);
    } catch (err: any) {
      error(err.message || "Не удалось загрузить отчёты");
    } finally {
      setLoading(false);
    }
  };

  const exportReport = async (type: string) => {
    try {
      const blob = await apiClient.download(`/reports/export?type=${type}&${query}`);
      const url = URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = url;
      link.download = `${type}-${Date.now()}.xlsx`;
      document.body.appendChild(link);
      link.click();
      link.remove();
      URL.revokeObjectURL(url);
      success("Файл Excel скачан");
    } catch (err: any) {
      error(err.message || "Не удалось экспортировать отчёт");
    }
  };

  return (
    <PageContainer title="Отчётность">
      <div className="grid grid-cols-1 gap-8">
        <ReportCard title="История ведомостей" type="audit" onExport={exportReport}>
          <Filters filters={filters} setFilters={setFilters} />
          <Table>
            <TableHead><TableRow isHoverable={false}><TableHeaderCell>Пользователь</TableHeaderCell><TableHeaderCell>Действие</TableHeaderCell><TableHeaderCell>Дата/время</TableHeaderCell></TableRow></TableHead>
            <TableBody>{loading ? loadingRow(3) : audit.length === 0 ? <TableEmpty message="Данные не найдены" /> : audit.map((log) => <TableRow key={log.id}><TableCell>{log.user_name}</TableCell><TableCell>{log.action}</TableCell><TableCell>{new Date(log.timestamp).toLocaleString("ru-RU")}</TableCell></TableRow>)}</TableBody>
          </Table>
        </ReportCard>

        <ReportCard title="Статистика успеваемости" type="grades" onExport={exportReport}>
          <Filters filters={filters} setFilters={setFilters} />
          <Table>
            <TableHead><TableRow isHoverable={false}><TableHeaderCell>Студент</TableHeaderCell><TableHeaderCell>Дисциплина</TableHeaderCell><TableHeaderCell>Оценка</TableHeaderCell><TableHeaderCell>Дата</TableHeaderCell></TableRow></TableHead>
            <TableBody>{loading ? loadingRow(4) : grades.length === 0 ? <TableEmpty message="Данные не найдены" /> : grades.map((row) => <TableRow key={row.id}><TableCell>{row.student_name}</TableCell><TableCell>{row.discipline}</TableCell><TableCell>{row.grade}</TableCell><TableCell>{new Date(row.date).toLocaleDateString("ru-RU")}</TableCell></TableRow>)}</TableBody>
          </Table>
        </ReportCard>

        <ReportCard title="Процент посещаемости" type="attendance" onExport={exportReport}>
          <Filters filters={filters} setFilters={setFilters} />
          <Table>
            <TableHead><TableRow isHoverable={false}><TableHeaderCell>Студент</TableHeaderCell><TableHeaderCell>Группа</TableHeaderCell><TableHeaderCell>Посещаемость</TableHeaderCell></TableRow></TableHead>
            <TableBody>{loading ? loadingRow(3) : attendance.length === 0 ? <TableEmpty message="Данные не найдены" /> : attendance.map((row) => <TableRow key={row.id}><TableCell>{row.student_name}</TableCell><TableCell>{row.group}</TableCell><TableCell>{row.attendance_percent}%</TableCell></TableRow>)}</TableBody>
          </Table>
        </ReportCard>
      </div>
    </PageContainer>
  );
}

function ReportCard({ title, type, onExport, children }: { title: string; type: string; onExport: (type: string) => void; children: React.ReactNode }) {
  return (
    <section className="rounded-lg border border-zinc-800 bg-zinc-900 p-6">
      <div className="mb-6 flex items-center justify-between gap-4">
        <h2 className="text-lg font-semibold text-white">{title}</h2>
        <Button variant="secondary" size="sm" onClick={() => onExport(type)} className="flex items-center gap-2">
          <Download size={16} />
          Экспорт в Excel
        </Button>
      </div>
      {children}
    </section>
  );
}

function Filters({ filters, setFilters }: { filters: { group: string; discipline: string; period: string }; setFilters: (filters: { group: string; discipline: string; period: string }) => void }) {
  return (
    <div className="mb-6 grid grid-cols-1 gap-4 md:grid-cols-3">
      <Select value={filters.group} onValueChange={(group) => setFilters({ ...filters, group })}>
        <SelectItem value="all">Все группы</SelectItem>
        <SelectItem value="CS-101">CS-101</SelectItem>
        <SelectItem value="CS-102">CS-102</SelectItem>
      </Select>
      <Select value={filters.discipline} onValueChange={(discipline) => setFilters({ ...filters, discipline })}>
        <SelectItem value="all">Все дисциплины</SelectItem>
        <SelectItem value="mathematics">Математика</SelectItem>
        <SelectItem value="physics">Физика</SelectItem>
      </Select>
      <Input type="month" value={filters.period} onChange={(event) => setFilters({ ...filters, period: event.target.value })} />
    </div>
  );
}

function loadingRow(colSpan: number) {
  return <TableRow isHoverable={false}><TableCell colSpan={colSpan} className="py-8 text-center">Загрузка...</TableCell></TableRow>;
}
