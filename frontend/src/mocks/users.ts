import { User } from "@/types/user";

export const users: User[] = [
  {
    id: "1",
    name: "Alexander Petrov",
    email: "alex.petrov@uni.edu",
    role: "admin",
    status: "active",
    createdAt: "2026-01-10"
  },

  {
    id: "2",
    name: "Maria Ivanova",
    email: "maria.ivanova@uni.edu",
    role: "teacher",
    status: "active",
    createdAt: "2026-01-12"
  },

  {
    id: "3",
    name: "John Smith",
    email: "john.smith@student.edu",
    role: "student",
    status: "blocked",
    createdAt: "2026-01-18"
  },

  {
    id: "4",
    name: "Emma Wilson",
    email: "emma.wilson@student.edu",
    role: "student",
    status: "pending",
    createdAt: "2026-01-20"
  },

  {
    id: "5",
    name: "Daniel Brown",
    email: "daniel.brown@uni.edu",
    role: "teacher",
    status: "pending",
    createdAt: "2026-01-22"
  }
];