import { signIn } from '@/lib/auth';

export default function LoginPage() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-background">
      <div className="w-full max-w-sm space-y-6 p-8">
        <div className="text-center space-y-2">
          <h1 className="text-3xl font-bold">Explainable Engine</h1>
          <p className="text-sm text-muted-foreground">
            Sign in to access the dashboard
          </p>
        </div>
        <form
          action={async () => {
            'use server';
            await signIn('google', { redirectTo: '/' });
          }}
        >
          <button
            type="submit"
            className="w-full rounded-md bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground hover:bg-primary/90"
          >
            Continue with Google
          </button>
        </form>
      </div>
    </div>
  );
}
