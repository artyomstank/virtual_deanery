import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { ArrowLeft, Edit2, Plus, Save, X } from "lucide-react";
import * as Tabs from "@radix-ui/react-tabs";
import { PageContainer } from "../../components/layout/PageContainer";
import { Button, Input } from "../../components/ui/Form";
import { Modal } from "../../components/ui/Modal";
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
import { EmploymentRecord, Teacher } from "../../types/api";

export function TeacherProfilePage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { success, error } = useNotificationStore();
  const [teacher, setTeacher] = useState<Teacher | null>(null);
  const [employmentRecords, setEmploymentRecords] = useState<EmploymentRecord[]>([]);
  const [disciplines, setDisciplines] = useState<Array<{ id: string; name: string; group?: string; semester?: string }>>([]);
  const [loading, setLoading] = useState(true);
  const [editing, setEditing] = useState(false);
  const [addEmploymentOpen, setAddEmploymentOpen] = useState(false);
  const [formData, setFormData] = useState({
    name: "",
    email: "",
    department: "",
    academic_degree: "",
    employee_id: "",
    phone: "",
  });
  const [employmentData, setEmploymentData] = useState({
    position: "",
    employment_type: "",
    start_date: "",
    end_date: "",
    notes: "",
  });

  useEffect(() => {
    if (id) void loadTeacher();
  }, [id]);

  const loadTeacher = async () => {
    try {
      setLoading(true);
      const response = await apiClient.get(`/admin/users/${id}/teacher`);
      setTeacher(response);
      setFormData({
        name: response?.name || "",
        email: response?.email || "",
        department: response?.department || "",
        academic_degree: response?.academic_degree || "",
        employee_id: response?.employee_id || "",
        phone: response?.phone || "",
      });
      const [employmentRes, disciplinesRes] = await Promise.all([
        apiClient.get(`/admin/users/${id}/employment`).catch(() => ({ records: [] })),
        apiClient.get(`/admin/users/${id}/disciplines`).catch(() => ({ disciplines: [] })),
      ]);
      setEmploymentRecords(employmentRes?.records || []);
      setDisciplines(disciplinesRes?.disciplines || []);
    } catch (err: any) {
      error(err.message || "Не удалось загрузить карточку преподавателя");
    } finally {
      setLoading(false);
    }
  };

  const handleSaveProfile = async () => {
    if (!id) return;
    try {
      await apiClient.put(`/admin/users/${id}/teacher`, formData);
      success("Данные сохранены");
      setEditing(false);
      await loadTeacher();
    } catch (err: any) {
      error(err.message || "Не удалось сохранить данные");
    }
  };

  const handleAddEmployment = async () => {
    if (!id) return;
    if (!employmentData.position || !employmentData.employment_type || !employmentData.start_date) {
      error("Заполните должность, тип занятости и дату начала");
      return;
    }
    try {
      await apiClient.post(`/admin/users/${id}/employment`, employmentData);
      success("Запись добавлена");
      setAddEmploymentOpen(false);
      setEmploymentData({ position: "", employment_type: "", start_date: "", end_date: "", notes: "" });
      await loadTeacher();
    } catch (err: any) {
      error(err.message || "Не удалось добавить запись");
    }
  };

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-zinc-950">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-indigo-500 border-t-transparent" />
      </div>
    );
  }

  return (
    <PageContainer
      title={teacher?.name || "Карточка преподавателя"}
      action={
        <Button variant="secondary" onClick={() => navigate("/admin/users")} className="flex items-center gap-2">
          <ArrowLeft size={18} />
          Назад
        </Button>
      }
    >
      <div className="space-y-6">
        <section className="rounded-lg border border-zinc-800 bg-zinc-900 p-6">
          <div className="mb-6 flex items-center justify-between">
            <h2 className="text-lg font-semibold text-white">Личные данные</h2>
            {!editing ? (
              <Button variant="secondary" size="sm" onClick={() => setEditing(true)} className="flex items-center gap-2">
                <Edit2 size={16} />
                Редактировать
              </Button>
            ) : (
              <div className="flex gap-2">
                <Button size="sm" onClick={handleSaveProfile} className="flex items-center gap-2"><Save size={16} />Сохранить</Button>
                <Button variant="secondary" size="sm" onClick={() => { setEditing(false); void loadTeacher(); }} className="flex items-center gap-2"><X size={16} />Отмена</Button>
              </div>
            )}
          </div>

          <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
            {[
              ["name", "ФИО"],
              ["email", "Email"],
              ["department", "Кафедра"],
              ["academic_degree", "Учёная степень"],
              ["employee_id", "Табельный номер"],
              ["phone", "Контактный телефон"],
            ].map(([key, label]) =>
              editing ? (
                <Input key={key} label={label} value={(formData as any)[key]} onChange={(event) => setFormData({ ...formData, [key]: event.target.value })} />
              ) : (
                <div key={key}>
                  <p className="text-sm text-zinc-400">{label}</p>
                  <p className="mt-1 font-medium text-white">{(teacher as any)?.[key] || "-"}</p>
                </div>
              )
            )}
          </div>
        </section>

        <Tabs.Root defaultValue="employment">
          <Tabs.List className="mb-6 flex gap-2 border-b border-zinc-800">
            <Tabs.Trigger value="employment" className="px-4 py-2 text-sm font-medium text-zinc-400 data-[state=active]:border-b-2 data-[state=active]:border-indigo-500 data-[state=active]:text-white">
              История трудоустройства
            </Tabs.Trigger>
            <Tabs.Trigger value="disciplines" className="px-4 py-2 text-sm font-medium text-zinc-400 data-[state=active]:border-b-2 data-[state=active]:border-indigo-500 data-[state=active]:text-white">
              Назначенные дисциплины
            </Tabs.Trigger>
          </Tabs.List>

          <Tabs.Content value="employment">
            <div className="mb-4 flex justify-end">
              <Button size="sm" onClick={() => setAddEmploymentOpen(true)} className="flex items-center gap-2">
                <Plus size={16} />
                Добавить запись
              </Button>
            </div>
            <Table>
              <TableHead>
                <TableRow isHoverable={false}>
                  <TableHeaderCell>Должность</TableHeaderCell>
                  <TableHeaderCell>Тип занятости</TableHeaderCell>
                  <TableHeaderCell>Дата начала</TableHeaderCell>
                  <TableHeaderCell>Дата окончания</TableHeaderCell>
                  <TableHeaderCell>Примечания</TableHeaderCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {employmentRecords.length === 0 ? (
                  <TableEmpty message="Данные не найдены" />
                ) : (
                  employmentRecords.map((record) => (
                    <TableRow key={record.id}>
                      <TableCell>{record.position}</TableCell>
                      <TableCell>{record.employment_type}</TableCell>
                      <TableCell>{new Date(record.start_date).toLocaleDateString("ru-RU")}</TableCell>
                      <TableCell>{record.end_date ? new Date(record.end_date).toLocaleDateString("ru-RU") : "-"}</TableCell>
                      <TableCell>{record.notes || "-"}</TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </Tabs.Content>

          <Tabs.Content value="disciplines">
            <Table>
              <TableHead>
                <TableRow isHoverable={false}>
                  <TableHeaderCell>Дисциплина</TableHeaderCell>
                  <TableHeaderCell>Группа</TableHeaderCell>
                  <TableHeaderCell>Семестр</TableHeaderCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {disciplines.length === 0 ? (
                  <TableEmpty message="Данные не найдены" />
                ) : (
                  disciplines.map((item) => (
                    <TableRow key={item.id}>
                      <TableCell>{item.name}</TableCell>
                      <TableCell>{item.group || "-"}</TableCell>
                      <TableCell>{item.semester || "-"}</TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </Tabs.Content>
        </Tabs.Root>
      </div>

      <Modal open={addEmploymentOpen} onOpenChange={setAddEmploymentOpen} title="Добавить запись" size="md">
        <div className="space-y-4">
          <Input label="Должность" value={employmentData.position} onChange={(event) => setEmploymentData({ ...employmentData, position: event.target.value })} />
          <Input label="Тип занятости" value={employmentData.employment_type} onChange={(event) => setEmploymentData({ ...employmentData, employment_type: event.target.value })} />
          <Input label="Дата начала" type="date" value={employmentData.start_date} onChange={(event) => setEmploymentData({ ...employmentData, start_date: event.target.value })} />
          <Input label="Дата окончания" type="date" value={employmentData.end_date} onChange={(event) => setEmploymentData({ ...employmentData, end_date: event.target.value })} />
          <Input label="Примечания" value={employmentData.notes} onChange={(event) => setEmploymentData({ ...employmentData, notes: event.target.value })} />
          <div className="flex justify-end gap-3">
            <Button variant="secondary" onClick={() => setAddEmploymentOpen(false)}>Отмена</Button>
            <Button onClick={handleAddEmployment}>Добавить</Button>
          </div>
        </div>
      </Modal>
    </PageContainer>
  );
}
