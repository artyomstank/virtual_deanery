import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button, Input } from "../../components/ui/Form";
import { useAuthStore } from "../../store/auth.store";

export function LoginPage() {
  const navigate = useNavigate();
  const login = useAuthStore((state) => state.login);
  const isLoading = useAuthStore((state) => state.isLoading);
  const clearError = useAuthStore((state) => state.clearError);

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loginError, setLoginError] = useState("");
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    clearError();
    setLoginError("");

    const nextErrors: Record<string, string> = {};
    if (!email.trim()) nextErrors.email = "Введите email";
    if (!password) nextErrors.password = "Введите пароль";

    if (Object.keys(nextErrors).length > 0) {
      setFieldErrors(nextErrors);
      return;
    }

    setFieldErrors({});

    try {
      await login(email, password);
      navigate("/admin", { replace: true });
    } catch {
      setLoginError("Неверный email или пароль");
      setFieldErrors({ email: " ", password: " " });
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-950 p-6">
      <div className="w-full max-w-md rounded-lg border border-zinc-800 bg-zinc-900 p-8 shadow-xl">
        <div className="mb-8 text-center">
          <h1 className="text-3xl font-bold text-white">Виртуальный деканат</h1>
          <p className="mt-2 text-zinc-400">Вход в систему</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            label="Email"
            type="email"
            placeholder="user@example.com"
            value={email}
            onChange={(event) => {
              setEmail(event.target.value);
              setLoginError("");
            }}
            error={fieldErrors.email}
            autoComplete="email"
          />

          <Input
            label="Пароль"
            type="password"
            placeholder="Введите пароль"
            value={password}
            onChange={(event) => {
              setPassword(event.target.value);
              setLoginError("");
            }}
            error={fieldErrors.password}
            autoComplete="current-password"
          />

          {loginError && (
            <p className="rounded-md border border-red-800 bg-red-950/40 px-3 py-2 text-sm text-red-300">
              {loginError}
            </p>
          )}

          <Button type="submit" isLoading={isLoading} className="w-full">
            Войти
          </Button>
        </form>

        <div className="mt-4 text-center">
          <a href="#" className="text-sm text-indigo-400 hover:text-indigo-300">
            Забыли пароль?
          </a>
        </div>

        <div className="mt-6 rounded-lg border border-zinc-800 bg-zinc-950 p-3 text-sm text-zinc-400">
          <p className="font-medium text-zinc-200">Демо-вход</p>
          <p>Email: admin@dekanat.local</p>
          <p>Пароль: admin123</p>
        </div>
      </div>
    </div>
  );
}
