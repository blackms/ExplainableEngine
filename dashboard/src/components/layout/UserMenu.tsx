'use client';

import { useSession, signOut } from 'next-auth/react';

export function UserMenu() {
  const { data: session } = useSession();

  if (!session) return null;

  return (
    <div className="flex items-center gap-2">
      <span className="text-sm">{session.user?.name}</span>
      <button
        onClick={() => signOut()}
        className="text-sm text-muted-foreground hover:text-foreground"
      >
        Sign out
      </button>
    </div>
  );
}
