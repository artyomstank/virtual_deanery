import { useEffect, useMemo, useState } from "react";
import * as Tabs from "@radix-ui/react-tabs";
import { Edit2, Plus, Trash2, UserMinus, Users } from "lucide-react";
import { PageContainer } from "../../components/layout/PageContainer";
import { Button, Input } from "../../components/ui/Form";
import { Select, SelectItem } from "../../components/ui/Select";
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
import { Department, ExamSchedule, Group, Student } from "../../types/api";

type Entity = "groups" | "exams" | "departments";

const emptyForms = {
  groups: { name: "", course: "", curator_name: "", student_count: "" },
  exams: { discipline: "", date: "", time: "", room: "", teacher_name: "", group: "" },
  departments: { name: "", office: "", head_name: "" },
};

export function DataManagementPage() {
  const { success, error } = useNotificationStore();
  const [activeTab, setActiveTab] = useState("students");
  const [loading, setLoading] = useState(false);
  const [students, setStudents] = useState<Student[]>([]);
  const [groups, setGroups] = useState<Group[]>([]);
  const [exams, setExams] = useState<ExamSchedule[]>([]);
  const [departments, setDepartments] = useState<Department[]>([]);
  const [groupFilter, setGroupFilter] = useState("all");
  const [courseFilter, setCourseFilter] = useState("all");
  const [studentStatusFilter, setStudentStatusFilter] = useState("all");
  const [examGroupFilter, setExamGroupFilter] = useState("all");
  const [examDateFilter, setExamDateFilter] = useState("");
  const [modalEntity, setModalEntity] = useState<Entity | null>(null);
  const [editingItem, setEditingItem] = useState<any>(null);
  const [formData, setFormData] = useState<Record<string, string>>({});
  const [formWarning, setFormWarning] = useState("");

  useEffect(() => {
    void loadData();
  }, [activeTab]);

  const loadData = async () => {
    try {
      setLoading(true);
      if (activeTab === "students") {
        const [studentsRes, groupsRes] = await Promise.all([
          apiClient.get("/admin/students"),
          apiClient.get("/admin/groups").catch(() => ({ groups: [] })),
        ]);
        setStudents(studentsRes?.students || []);
        setGroups(groupsRes?.groups || []);
      }
      if (activeTab === "groups") {
        const response = await apiClient.get("/admin/groups");
        setGroups(response?.groups || []);
      }
      if (activeTab === "exams") {
        const response = await apiClient.get("/admin/exam-schedule");
        setExams(response?.exams || []);
      }
      if (activeTab === "departments") {
        const response = await apiClient.get("/admin/departments");
        setDepartments(response?.departments || []);
      }
    } catch (err: any) {
      error(err.message || "Не удалось загрузить учебные данные");
    } finally {
      setLoading(false);
    }
  };

  const filteredStudents = useMemo(() => students.filter((student) => {
    if (groupFilter !== "all" && student.group !== groupFilter) return false;
    if (courseFilter !== "all" && String(student.course) !== courseFilter) return false;
    if (studentStatusFilter !== "all" && student.status !== studentStatusFilter) return false;
    return true;
  }), [students, groupFilter, courseFilter, studentStatusFilter]);

  const filteredExams = useMemo(() => exams.filter((exam: any) => {
    if (examGroupFilter !== "all" && exam.group !== examGroupFilter) return false;
    if (examDateFilter && exam.date !== examDateFilter) return false;
    return true;
  }), [exams, examGroupFilter, examDateFilter]);

  const openModal = (entity: Entity, item?: any) => {
    setModalEntity(entity);
    setEditingItem(item || null);
    setFormWarning("");
    setFormData(item ? Object.fromEntries(Object.entries(item).map(([key, value]) => [key, String(value ?? "")])) : emptyForms[entity]);
  };

  const closeModal = () => {
    setModalEntity(null);
    setEditingItem(null);
    setFormData({});
    setFormWarning("");
  };

  const saveEntity = async () => {
    if (!modalEntity) return;
    const path = modalEntity === "exams" ? "/admin/exam-schedule" : `/admin/${modalEntity}`;
    try {
      const payload = {
        ...formData,
        course: formData.course ? Number(formData.course) : undefined,
        student_count: formData.student_count ? Number(formData.student_count) : undefined,
      };
      if (modalEntity === "exams") {
        const conflict = exams.some((exam) => exam.id !== editingItem?.id && exam.room === formData.room && exam.date === formData.date && exam.time === formData.time);
        if (conflict) {
          setFormWarning("Конфликт аудитории и времени");
          return;
        }
      }
      if (editingItem) {
        await apiClient.put(`${path}/${editingItem.id}`, payload);
        success("Данные обновлены");
      } else {
        await apiClient.post(path, payload);
        success("Запись создана");
      }
      closeModal();
      await loadData();
    } catch (err: any) {
      error(err.message || "Не удалось сохранить запись");
    }
  };

  const deleteEntity = async (entity: Entity, id: string) => {
    const path = entity === "exams" ? "/admin/exam-schedule" : `/admin/${entity}`;
    try {
      await apiClient.delete(`${path}/${id}`);
      success("Запись удалена");
      await loadData();
    } catch (err: any) {
      error(err.message || "Не удалось удалить запись");
    }
  };

  const updateStudentStatus = async (student: Student, patch: Record<string, string>) => {
    try {
      await apiClient.patch(`/admin/students/${student.id}`, patch);
      success(patch.status === "expelled" ? "Студент отчислен" : "Данные студента обновлены");
      await loadData();
    } catch (err: any) {
      error(err.message || "Не удалось изменить студента");
    }
  };

  return (
    <PageContainer title="Учебные данные">
      <Tabs.Root value={activeTab} onValueChange={setActiveTab}>
        <Tabs.List className="mb-6 flex flex-wrap gap-2 border-b border-zinc-800">
          {[
            ["students", "Студенты"],
            ["groups", "Группы"],
            ["exams", "Расписание экзаменов"],
            ["departments", "Кафедры"],
          ].map(([value, label]) => (
            <Tabs.Trigger key={value} value={value} className="px-4 py-2 text-sm font-medium text-zinc-400 data-[state=active]:border-b-2 data-[state=active]:border-indigo-500 data-[state=active]:text-white">
              {label}
            </Tabs.Trigger>
          ))}
        </Tabs.List>

        <Tabs.Content value="students">
          <div className="mb-6 grid grid-cols-1 gap-4 md:grid-cols-3">
            <Select value={groupFilter} onValueChange={setGroupFilter}>
              <SelectItem value="all">Все группы</SelectItem>
              {groups.map((group) => <SelectItem key={group.id} value={group.name}>{group.name}</SelectItem>)}
            </Select>
            <Select value={courseFilter} onValueChange={setCourseFilter}>
              <SelectItem value="all">Все курсы</SelectItem>
              {[1, 2, 3, 4, 5].map((course) => <SelectItem key={course} value={String(course)}>{course} курс</SelectItem>)}
            </Select>
            <Select value={studentStatusFilter} onValueChange={setStudentStatusFilter}>
              <SelectItem value="all">Все статусы</SelectItem>
              <SelectItem value="active">Активен</SelectItem>
              <SelectItem value="expelled">Отчислен</SelectItem>
            </Select>
          </div>
          <Table>
            <TableHead>
              <TableRow isHoverable={false}>
                <TableHeaderCell>ФИО</TableHeaderCell>
                <TableHeaderCell>Группа</TableHeaderCell>
                <TableHeaderCell>Курс</TableHeaderCell>
                <TableHeaderCell>Статус</TableHeaderCell>
                <TableHeaderCell>Действия</TableHeaderCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? <TableRow isHoverable={false}><TableCell colSpan={5} className="py-8 text-center">Загрузка...</TableCell></TableRow> :
                filteredStudents.length === 0 ? <TableEmpty message="Данные не найдены" /> :
                filteredStudents.map((student) => (
                  <TableRow key={student.id} className={student.status === "expelled" ? "bg-zinc-900/60 opacity-60" : ""}>
                    <TableCell>{student.name}</TableCell>
                    <TableCell>{student.group}</TableCell>
                    <TableCell>{student.course}</TableCell>
                    <TableCell>{student.status === "active" ? "Активен" : "Отчислен"}</TableCell>
                    <TableCell>
                      <div className="flex flex-wrap gap-2">
                        <Button size="sm" variant="secondary"><Edit2 size={14} />Редактировать</Button>
                        <Button size="sm" variant="secondary" onClick={() => updateStudentStatus(student, { group: groups[0]?.name || student.group })}><Users size={14} />Перевести</Button>
                        <Button size="sm" variant="danger" onClick={() => updateStudentStatus(student, { status: "expelled" })}><UserMinus size={14} />Отчислить</Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
            </TableBody>
          </Table>
        </Tabs.Content>

        <Tabs.Content value="groups">
          <CrudHeader label="Добавить группу" onClick={() => openModal("groups")} />
          <Table>
            <TableHead><TableRow isHoverable={false}><TableHeaderCell>Название</TableHeaderCell><TableHeaderCell>Курс</TableHeaderCell><TableHeaderCell>Куратор</TableHeaderCell><TableHeaderCell>Количество студентов</TableHeaderCell><TableHeaderCell>Действия</TableHeaderCell></TableRow></TableHead>
            <TableBody>{renderRows(groups, "groups", loading, openModal, deleteEntity)}</TableBody>
          </Table>
        </Tabs.Content>

        <Tabs.Content value="exams">
          <div className="mb-6 grid grid-cols-1 gap-4 md:grid-cols-[1fr_180px_auto]">
            <Select value={examGroupFilter} onValueChange={setExamGroupFilter}><SelectItem value="all">Все группы</SelectItem>{groups.map((group) => <SelectItem key={group.id} value={group.name}>{group.name}</SelectItem>)}</Select>
            <Input type="date" value={examDateFilter} onChange={(event) => setExamDateFilter(event.target.value)} />
            <Button onClick={() => openModal("exams")}><Plus size={18} />Добавить экзамен</Button>
          </div>
          <Table>
            <TableHead><TableRow isHoverable={false}><TableHeaderCell>Дисциплина</TableHeaderCell><TableHeaderCell>Дата</TableHeaderCell><TableHeaderCell>Время</TableHeaderCell><TableHeaderCell>Аудитория</TableHeaderCell><TableHeaderCell>Преподаватель</TableHeaderCell><TableHeaderCell>Действия</TableHeaderCell></TableRow></TableHead>
            <TableBody>{renderRows(filteredExams, "exams", loading, openModal, deleteEntity)}</TableBody>
          </Table>
        </Tabs.Content>

        <Tabs.Content value="departments">
          <CrudHeader label="Добавить кафедру" onClick={() => openModal("departments")} />
          <Table>
            <TableHead><TableRow isHoverable={false}><TableHeaderCell>Название</TableHeaderCell><TableHeaderCell>Кабинет</TableHeaderCell><TableHeaderCell>Заведующий</TableHeaderCell><TableHeaderCell>Действия</TableHeaderCell></TableRow></TableHead>
            <TableBody>{renderRows(departments, "departments", loading, openModal, deleteEntity)}</TableBody>
          </Table>
        </Tabs.Content>
      </Tabs.Root>

      <Modal open={!!modalEntity} onOpenChange={(open) => !open && closeModal()} title={editingItem ? "Редактировать запись" : "Создать запись"} size="md">
        {modalEntity && (
          <div className="space-y-4">
            {Object.keys(emptyForms[modalEntity]).map((key) => (
              <Input key={key} label={fieldLabel(key)} type={key === "date" ? "date" : key === "time" ? "time" : "text"} value={formData[key] || ""} onChange={(event) => setFormData({ ...formData, [key]: event.target.value })} />
            ))}
            {formWarning && <p className="rounded-md border border-red-800 bg-red-950/40 px-3 py-2 text-sm text-red-300">{formWarning}</p>}
            <div className="flex justify-end gap-3">
              <Button variant="secondary" onClick={closeModal}>Отмена</Button>
              <Button onClick={saveEntity}>Сохранить</Button>
            </div>
          </div>
        )}
      </Modal>
    </PageContainer>
  );
}

function CrudHeader({ label, onClick }: { label: string; onClick: () => void }) {
  return <div className="mb-6 flex justify-end"><Button onClick={onClick}><Plus size={18} />{label}</Button></div>;
}

function fieldLabel(key: string) {
  const labels: Record<string, string> = {
    name: "Название",
    course: "Курс",
    curator_name: "Куратор",
    student_count: "Количество студентов",
    discipline: "Дисциплина",
    date: "Дата",
    time: "Время",
    room: "Аудитория",
    teacher_name: "Преподаватель",
    group: "Группа",
    office: "Кабинет",
    head_name: "Заведующий",
  };
  return labels[key] || key;
}

function renderRows(items: any[], entity: Entity, loading: boolean, openModal: (entity: Entity, item?: any) => void, deleteEntity: (entity: Entity, id: string) => void) {
  if (loading) return <TableRow isHoverable={false}><TableCell colSpan={6} className="py-8 text-center">Загрузка...</TableCell></TableRow>;
  if (items.length === 0) return <TableEmpty message="Данные не найдены" />;
  return items.map((item) => (
    <TableRow key={item.id}>
      {entity === "groups" && <><TableCell>{item.name}</TableCell><TableCell>{item.course}</TableCell><TableCell>{item.curator_name}</TableCell><TableCell>{item.student_count}</TableCell></>}
      {entity === "exams" && <><TableCell>{item.discipline}</TableCell><TableCell>{new Date(item.date).toLocaleDateString("ru-RU")}</TableCell><TableCell>{item.time}</TableCell><TableCell>{item.room}</TableCell><TableCell>{item.teacher_name}</TableCell></>}
      {entity === "departments" && <><TableCell>{item.name}</TableCell><TableCell>{item.office}</TableCell><TableCell>{item.head_name}</TableCell></>}
      <TableCell>
        <div className="flex gap-2">
          <Button size="sm" variant="secondary" onClick={() => openModal(entity, item)}><Edit2 size={14} />Редактировать</Button>
          <Button size="sm" variant="danger" onClick={() => deleteEntity(entity, item.id)}><Trash2 size={14} />Удалить</Button>
        </div>
      </TableCell>
    </TableRow>
  ));
}
