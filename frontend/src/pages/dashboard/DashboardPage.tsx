import {
  Users,
  GraduationCap,
  UserSquare2,
  ShieldAlert
} from "lucide-react";

import { PageContainer } from "@/components/layout/PageContainer";
import { StatCard } from "@/components/cards/StatCard";

import { users } from "@/mocks/users";

export function DashboardPage() {
  const pendingUsers = users.filter(
    (u) => u.status === "pending"
  ).length;

  const activeUsers = users.filter(
    (u) => u.status === "active"
  ).length;

  const blockedUsers = users.filter(
    (u) => u.status === "blocked"
  ).length;

  const students = users.filter(
    (u) => u.role === "student"
  ).length;

  return (
    <PageContainer
      title="Dashboard"
      description="Overview of university system activity"
    >
      <div className="grid gap-6 md:grid-cols-2 xl:grid-cols-4">
        <StatCard
          title="Pending approvals"
          value={pendingUsers}
          icon={ShieldAlert}
          description="Users waiting confirmation"
        />

        <StatCard
          title="Students"
          value={students}
          icon={GraduationCap}
          description="Registered students"
        />

        <StatCard
          title="Active users"
          value={activeUsers}
          icon={Users}
          description="Currently active accounts"
        />

        <StatCard
          title="Blocked users"
          value={blockedUsers}
          icon={UserSquare2}
          description="Suspended accounts"
        />
      </div>

      <div className="mt-8 grid gap-6 xl:grid-cols-[2fr,1fr]">
        <div className="rounded-2xl border border-zinc-800 bg-zinc-900 p-6">
          <div className="mb-6 flex items-center justify-between">
            <h2 className="text-xl font-semibold">
              Recent activity
            </h2>

            <button className="text-sm text-indigo-400">
              View all
            </button>
          </div>

          <div className="space-y-4">
            {[
              "Admin approved new teacher account",
              "Schedule updated for CS-101",
              "Grades changed in Mathematics",
              "Student transferred to another group"
            ].map((item) => (
              <div
                key={item}
                className="
                  flex items-center justify-between
                  rounded-xl border border-zinc-800
                  bg-zinc-950 px-4 py-4
                "
              >
                <div>
                  <p className="font-medium">
                    {item}
                  </p>

                  <p className="mt-1 text-sm text-zinc-500">
                    5 minutes ago
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="rounded-2xl border border-zinc-800 bg-zinc-900 p-6">
          <h2 className="mb-6 text-xl font-semibold">
            Quick actions
          </h2>

          <div className="space-y-3">
            {[
              "Create user",
              "Open schedule",
              "Manage ACL",
              "View audit logs"
            ].map((item) => (
              <button
                key={item}
                className="
                  flex h-14 w-full items-center rounded-xl
                  border border-zinc-800 bg-zinc-950 px-4
                  text-left font-medium transition
                  hover:border-indigo-500
                  hover:bg-indigo-500/10
                "
              >
                {item}
              </button>
            ))}
          </div>
        </div>
      </div>
    </PageContainer>
  );
}