import { createBrowserRouter, Navigate } from "react-router-dom";
import { useAuthStore } from "../store/auth.store";
import { AuthLayout } from "./layouts/AuthLayout";
import { AdminLayout } from "./layouts/AdminLayout";
import { LoginPage } from "../pages/auth/LoginPage";
import { DashboardPage } from "../pages/dashboard/DashboardPage";
import { UsersPage } from "../pages/users/UsersPage";
import { DataManagementPage } from "../pages/data/DataManagementPage";
import { ReportsPage } from "../pages/reports/ReportsPage";
import { ACLPage } from "../pages/acl/ACLPage";
import { TeacherProfilePage } from "../pages/teachers/TeacherProfilePage";

// Protected route wrapper
function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const token = useAuthStore((state) => state.token);
  
  if (!token) {
    return <Navigate to="/login" replace />;
  }
  
  return children;
}

export const router = createBrowserRouter([
  {
    path: "/login",
    element: <AuthLayout />,
    children: [
      {
        index: true,
        element: <LoginPage />,
      },
    ],
  },
  {
    path: "/admin",
    element: (
      <ProtectedRoute>
        <AdminLayout />
      </ProtectedRoute>
    ),
    children: [
      {
        index: true,
        element: <DashboardPage />,
      },
      {
        path: "users",
        element: <UsersPage />,
      },
      {
        path: "users/:id/teacher",
        element: <TeacherProfilePage />,
      },
      {
        path: "data",
        element: <DataManagementPage />,
      },
      {
        path: "reports",
        element: <ReportsPage />,
      },
      {
        path: "acl",
        element: <ACLPage />,
      },
    ],
  },
  {
    path: "/",
    element: <Navigate to="/admin" replace />,
  },
  {
    path: "*",
    element: <Navigate to="/admin" replace />,
  },
]);
