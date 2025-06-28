import { createFileRoute } from '@tanstack/react-router';
import { LoginForm } from '@/components/LoginForm';

export const Route = createFileRoute('/login')({
  component: LoginPage,
});

function LoginPage() {
  return (
    <div className="h-screen w-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
      <LoginForm />
    </div>
  );
}