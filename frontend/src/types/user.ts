export type UserRole =
  | "admin"
  | "teacher"
  | "student";

export type UserStatus =
  | "active"
  | "blocked"
  | "pending";

export interface User {
  id: string;
  name: string;
  email: string;
  role: UserRole;
  status: UserStatus;
  createdAt: string;
}