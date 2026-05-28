import { useNavigate } from "react-router-dom";

export function LoginPage() {
  const navigate = useNavigate();

  return (
    <div className="flex min-h-screen items-center justify-center p-6">
      <div className="w-full max-w-md rounded-2xl border border-zinc-800 bg-zinc-900 p-8 shadow-soft">
        <div className="mb-8">
          <h1 className="text-3xl font-bold">
            Welcome back
          </h1>

          <p className="mt-2 text-zinc-400">
            Sign in to continue
          </p>
        </div>

        <div className="space-y-4">
          <input
            placeholder="Email"
            className="h-12 w-full rounded-xl border border-zinc-800 bg-zinc-950 px-4 outline-none focus:border-indigo-500"
          />

          <input
            type="password"
            placeholder="Password"
            className="h-12 w-full rounded-xl border border-zinc-800 bg-zinc-950 px-4 outline-none focus:border-indigo-500"
          />

          <button
            onClick={() => navigate("/dashboard")}
            className="h-12 w-full rounded-xl bg-indigo-500 font-medium transition hover:bg-indigo-400"
          >
            Sign In
          </button>
        </div>
      </div>
    </div>
  );
}