import { useEffect, useMemo, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { Check, Edit2, Lock, LockOpen, Plus, Search, Trash2, X } from "lucide-react";
import * as Tabs from "@radix-ui/react-tabs";
import { PageContainer } from "../../components/layout/PageContainer";
import { Badge, Button, Input } from "../../components/ui/Form";
import { Select, SelectItem } from "../../components/ui/Select";
import { Modal } from "../../components/ui/Modal";
import { ConfirmDialog } from "../../components/ui/ConfirmDialog";
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
import { User } from "../../types/api";

const roleLabels: Record<User["role"], string> = {
  student: "Студент",
  teacher: "Преподаватель",
  dean: "Декан",
  admin: "Админ",
};

const statusLabels: Record<User["status"], string> = {
  active: "Активен",
  blocked: "Заблокирован",
  pending: "Ожидает",
};

const pageSize = 10;

export function UsersPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { success, error } = useNotificationStore();
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState(searchParams.get("tab") === "pending" ? "pending" : "all");
  const [roleFilter, setRoleFilter] = useState("all");
  const [statusFilter, setStatusFilter] = useState("all");
  const [searchQuery, setSearchQuery] = useState("");
  const [page, setPage] = useState(1);
  const [selectedPending, setSelectedPending] = useState<string[]>([]);
  const [createOpen, setCreateOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);
  const [actionUser, setActionUser] = useState<User | null>(null);
  const [blockDialogOpen, setBlockDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
  const [formData, setFormData] = useState({
    name: "",
    email: "",
    password: "",
    role: "student" as User["role"],
  });

  useEffect(() => {
    void fetchUsers();
  }, []);

  useEffect(() => {
    setPage(1);
    setSelectedPending([]);
  }, [activeTab, roleFilter, statusFilter, searchQuery]);

  const fetchUsers = async () => {
    try {
      setLoading(true);
      const response = await apiClient.get("/admin/users");
      setUsers(response?.users || []);
    } catch (err: any) {
      error(err.message || "Не удалось загрузить пользователей");
    } finally {
      setLoading(false);
    }
  };

  const filteredUsers = useMemo(() => {
    return users.filter((user) => {
      if (activeTab === "pending" && user.status !== "pending") return false;
      if (activeTab === "all" && user.status === "pending") return false;
      if (roleFilter !== "all" && user.role !== roleFilter) return false;
      if (statusFilter !== "all" && user.status !== statusFilter) return false;
      if (searchQuery.trim()) {
        const query = searchQuery.trim().toLowerCase();
        return user.name.toLowerCase().includes(query) || user.email.toLowerCase().includes(query);
      }
      return true;
    });
  }, [users, activeTab, roleFilter, statusFilter, searchQuery]);

  const pages = Math.max(1, Math.ceil(filteredUsers.length / pageSize));
  const pageUsers = filteredUsers.slice((page - 1) * pageSize, page * pageSize);

  const resetForm = () => {
    setFormData({ name: "", email: "", password: "", role: "student" });
    setFieldErrors({});
  };

  const validate = () => {
    const errors: Record<string, string> = {};
    if (!formData.name.trim()) errors.name = "Введите имя";
    if (!formData.email.trim()) errors.email = "Введите email";
    if (!editingUser && !formData.password) errors.password = "Введите пароль";
    setFieldErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSaveUser = async () => {
    if (!validate()) return;
    try {
      setSaving(true);
      if (editingUser) {
        await apiClient.put(`/admin/users/${editingUser.id}`, {
          name: formData.name,
          email: formData.email,
          role: formData.role,
        });
        success("Данные обновлены");
        setEditingUser(null);
      } else {
        await apiClient.post("/admin/users", formData);
        success("Пользователь создан");
        setCreateOpen(false);
      }
      resetForm();
      await fetchUsers();
    } catch (err: any) {
      if (err.status === 409 || /email/i.test(err.message || "")) {
        setFieldErrors((current) => ({ ...current, email: "Этот email уже занят" }));
      } else {
        error(err.message || "Не удалось сохранить пользователя");
      }
    } finally {
      setSaving(false);
    }
  };

  const openEditModal = (user: User) => {
    setEditingUser(user);
    setFormData({ name: user.name, email: user.email, password: "", role: user.role });
    setFieldErrors({});
  };

  const closeUserModal = () => {
    setCreateOpen(false);
    setEditingUser(null);
    resetForm();
  };

  const handleBlockUser = async () => {
    if (!actionUser) return;
    try {
      await apiClient.patch(`/admin/users/${actionUser.id}/block`, {});
      success(actionUser.status === "blocked" ? "Пользователь разблокирован" : "Пользователь заблокирован");
      setBlockDialogOpen(false);
      setActionUser(null);
      await fetchUsers();
    } catch (err: any) {
      error(err.message || "Не удалось изменить статус");
    }
  };

  const handleDeleteUser = async () => {
    if (!actionUser) return;
    try {
      await apiClient.delete(`/admin/users/${actionUser.id}`);
      success("Пользователь удалён");
      setDeleteDialogOpen(false);
      setActionUser(null);
      await fetchUsers();
    } catch (err: any) {
      error(err.message || "Не удалось удалить пользователя");
    }
  };

  const approveUsers = async (ids: string[]) => {
    try {
      await Promise.all(ids.map((id) => apiClient.patch(`/admin/users/${id}/approve`, {})));
      success(ids.length > 1 ? "Выбранные пользователи подтверждены" : "Пользователь подтверждён");
      setSelectedPending([]);
      await fetchUsers();
    } catch (err: any) {
      error(err.message || "Не удалось подтвердить пользователей");
    }
  };

  const rejectUsers = async (ids: string[]) => {
    try {
      await Promise.all(ids.map((id) => apiClient.delete(`/admin/users/${id}`)));
      success(ids.length > 1 ? "Выбранные пользователи отклонены" : "Пользователь отклонён");
      setSelectedPending([]);
      await fetchUsers();
    } catch (err: any) {
      error(err.message || "Не удалось отклонить пользователей");
    }
  };

  const resetPassword = async () => {
    if (!editingUser) return;
    try {
      await apiClient.patch(`/admin/users/${editingUser.id}/reset-password`, {});
      success("Пароль сброшен");
    } catch (err: any) {
      error(err.message || "Не удалось сбросить пароль");
    }
  };

  return (
    <PageContainer
      title="Пользователи"
      action={
        <Button onClick={() => setCreateOpen(true)} className="flex items-center gap-2">
          <Plus size={18} />
          Создать пользователя
        </Button>
      }
    >
      <div className="mb-6 grid grid-cols-1 gap-4 md:grid-cols-[1fr_180px_180px_auto]">
        <Input placeholder="Поиск по имени или email" value={searchQuery} onChange={(event) => setSearchQuery(event.target.value)} />
        <Select value={roleFilter} onValueChange={setRoleFilter}>
          <SelectItem value="all">Все роли</SelectItem>
          <SelectItem value="student">Студент</SelectItem>
          <SelectItem value="teacher">Преподаватель</SelectItem>
          <SelectItem value="dean">Декан</SelectItem>
          <SelectItem value="admin">Админ</SelectItem>
        </Select>
        <Select value={statusFilter} onValueChange={setStatusFilter}>
          <SelectItem value="all">Все статусы</SelectItem>
          <SelectItem value="active">Активен</SelectItem>
          <SelectItem value="blocked">Заблокирован</SelectItem>
        </Select>
        <Button variant="secondary" onClick={() => setPage(1)} className="flex items-center gap-2">
          <Search size={16} />
          Найти
        </Button>
      </div>

      <Tabs.Root value={activeTab} onValueChange={setActiveTab}>
        <Tabs.List className="mb-6 flex gap-2 border-b border-zinc-800">
          <Tabs.Trigger value="all" className="px-4 py-2 text-sm font-medium text-zinc-400 data-[state=active]:border-b-2 data-[state=active]:border-indigo-500 data-[state=active]:text-white">
            Все пользователи
          </Tabs.Trigger>
          <Tabs.Trigger value="pending" className="px-4 py-2 text-sm font-medium text-zinc-400 data-[state=active]:border-b-2 data-[state=active]:border-indigo-500 data-[state=active]:text-white">
            Ожидают подтверждения
          </Tabs.Trigger>
        </Tabs.List>

        {selectedPending.length > 0 && activeTab === "pending" && (
          <div className="mb-4 flex flex-wrap items-center gap-3 rounded-lg border border-zinc-800 bg-zinc-900 p-3">
            <span className="text-sm text-zinc-300">Выбрано: {selectedPending.length}</span>
            <Button size="sm" onClick={() => approveUsers(selectedPending)}>Подтвердить выбранных</Button>
            <Button size="sm" variant="danger" onClick={() => rejectUsers(selectedPending)}>Отклонить выбранных</Button>
          </div>
        )}

        <Tabs.Content value={activeTab}>
          <Table>
            <TableHead>
              <TableRow isHoverable={false}>
                {activeTab === "pending" && <TableHeaderCell className="w-10"> </TableHeaderCell>}
                <TableHeaderCell>Имя</TableHeaderCell>
                <TableHeaderCell>Email</TableHeaderCell>
                <TableHeaderCell>Роль</TableHeaderCell>
                <TableHeaderCell>Статус</TableHeaderCell>
                <TableHeaderCell>Дата регистрации</TableHeaderCell>
                <TableHeaderCell>Действия</TableHeaderCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow isHoverable={false}>
                  <TableCell colSpan={7} className="py-8 text-center text-zinc-400">Загрузка...</TableCell>
                </TableRow>
              ) : pageUsers.length === 0 ? (
                <TableEmpty message="Данные не найдены" />
              ) : (
                pageUsers.map((user) => (
                  <TableRow key={user.id} onClick={() => user.role === "teacher" && navigate(`/admin/users/${user.id}/teacher`)}>
                    {activeTab === "pending" && (
                      <TableCell>
                        <input
                          type="checkbox"
                          checked={selectedPending.includes(user.id)}
                          onClick={(event) => event.stopPropagation()}
                          onChange={(event) =>
                            setSelectedPending((current) =>
                              event.target.checked ? [...current, user.id] : current.filter((id) => id !== user.id)
                            )
                          }
                          className="h-4 w-4 rounded border-zinc-700 bg-zinc-800 accent-indigo-600"
                        />
                      </TableCell>
                    )}
                    <TableCell className="font-medium text-white">{user.name}</TableCell>
                    <TableCell>{user.email}</TableCell>
                    <TableCell>{roleLabels[user.role]}</TableCell>
                    <TableCell>
                      <Badge variant={user.status === "active" ? "success" : user.status === "blocked" ? "error" : "warning"}>
                        {statusLabels[user.status]}
                      </Badge>
                    </TableCell>
                    <TableCell>{new Date(user.created_at).toLocaleDateString("ru-RU")}</TableCell>
                    <TableCell>
                      {activeTab === "pending" ? (
                        <div className="flex gap-2" onClick={(event) => event.stopPropagation()}>
                          <Button size="sm" onClick={() => approveUsers([user.id])}><Check size={14} />Подтвердить</Button>
                          <Button size="sm" variant="danger" onClick={() => rejectUsers([user.id])}><X size={14} />Отклонить</Button>
                        </div>
                      ) : (
                        <div className="flex items-center gap-2" onClick={(event) => event.stopPropagation()}>
                          <button title="Редактировать" onClick={() => openEditModal(user)} className="rounded p-1.5 text-zinc-300 hover:bg-zinc-800"><Edit2 size={16} /></button>
                          <button title="Заблокировать/разблокировать" onClick={() => { setActionUser(user); setBlockDialogOpen(true); }} className="rounded p-1.5 text-zinc-300 hover:bg-zinc-800">
                            {user.status === "blocked" ? <LockOpen size={16} /> : <Lock size={16} />}
                          </button>
                          <button title="Удалить" onClick={() => { setActionUser(user); setDeleteDialogOpen(true); }} className="rounded p-1.5 text-red-400 hover:bg-red-950/40"><Trash2 size={16} /></button>
                        </div>
                      )}
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </Tabs.Content>
      </Tabs.Root>

      <div className="mt-4 flex items-center justify-between text-sm text-zinc-400">
        <span>Страница {page} из {pages}</span>
        <div className="flex gap-2">
          <Button size="sm" variant="secondary" disabled={page === 1} onClick={() => setPage((current) => current - 1)}>Назад</Button>
          <Button size="sm" variant="secondary" disabled={page === pages} onClick={() => setPage((current) => current + 1)}>Вперёд</Button>
        </div>
      </div>

      <Modal open={createOpen || !!editingUser} onOpenChange={(open) => !open && closeUserModal()} title={editingUser ? "Редактировать пользователя" : "Создать пользователя"} size="md">
        <div className="space-y-4">
          <Input label="Имя" value={formData.name} onChange={(event) => setFormData({ ...formData, name: event.target.value })} error={fieldErrors.name} />
          <Input label="Email" type="email" value={formData.email} onChange={(event) => setFormData({ ...formData, email: event.target.value })} error={fieldErrors.email} />
          {!editingUser && <Input label="Пароль" type="password" value={formData.password} onChange={(event) => setFormData({ ...formData, password: event.target.value })} error={fieldErrors.password} />}
          <Select label="Роль" value={formData.role} onValueChange={(value) => setFormData({ ...formData, role: value as User["role"] })}>
            <SelectItem value="student">Студент</SelectItem>
            <SelectItem value="teacher">Преподаватель</SelectItem>
            <SelectItem value="dean">Декан</SelectItem>
            <SelectItem value="admin">Админ</SelectItem>
          </Select>
          {editingUser && <Button variant="secondary" onClick={resetPassword}>Сбросить пароль</Button>}
          <div className="flex justify-end gap-3 pt-2">
            <Button variant="secondary" onClick={closeUserModal}>Отмена</Button>
            <Button isLoading={saving} onClick={handleSaveUser}>{editingUser ? "Сохранить" : "Создать"}</Button>
          </div>
        </div>
      </Modal>

      <ConfirmDialog
        open={blockDialogOpen}
        onOpenChange={setBlockDialogOpen}
        title={actionUser?.status === "blocked" ? "Разблокировать пользователя?" : "Заблокировать пользователя?"}
        description={`${actionUser?.status === "blocked" ? "Разблокировать" : "Заблокировать"} пользователя ${actionUser?.name}?`}
        confirmText="Подтвердить"
        cancelText="Отмена"
        isDangerous={actionUser?.status !== "blocked"}
        onConfirm={handleBlockUser}
      />

      <ConfirmDialog
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        title="Удалить пользователя?"
        description={`Удалить пользователя ${actionUser?.name}? Это действие необратимо.`}
        confirmText="Удалить"
        cancelText="Отмена"
        isDangerous
        onConfirm={handleDeleteUser}
      />
    </PageContainer>
  );
}
