'use client';

import { usePathname } from 'next/navigation';
import { ChevronRight } from 'lucide-react';
import { UserMenu } from './UserMenu';
import { ThemeToggle } from './ThemeToggle';

const routeLabels: Record<string, string> = {
  '/': 'Home',
  '/audit': 'Audit',
  '/monitor': 'Monitor',
  '/playground': 'Playground',
  '/settings': 'Settings',
  '/explain': 'Explanation',
};

function getBreadcrumbs(pathname: string) {
  if (pathname === '/') return [{ label: 'Home', href: '/' }];

  const segments = pathname.split('/').filter(Boolean);
  const crumbs: { label: string; href: string }[] = [];
  let path = '';

  for (const segment of segments) {
    path += `/${segment}`;
    const label = routeLabels[path] || segment;
    crumbs.push({ label, href: path });
  }

  return crumbs;
}

export function Header() {
  const pathname = usePathname();
  const breadcrumbs = getBreadcrumbs(pathname);

  return (
    <header className="flex h-14 shrink-0 items-center justify-between border-b border-border bg-background px-4 pl-14 md:pl-4">
      <nav className="flex items-center gap-1.5 text-sm">
        {breadcrumbs.map((crumb, i) => (
          <span key={crumb.href} className="flex items-center gap-1.5">
            {i > 0 && (
              <ChevronRight className="size-3.5 text-muted-foreground" />
            )}
            <span
              className={
                i === breadcrumbs.length - 1
                  ? 'font-medium text-foreground'
                  : 'text-muted-foreground'
              }
            >
              {crumb.label}
            </span>
          </span>
        ))}
      </nav>
      <div className="flex items-center gap-3">
        <ThemeToggle />
        <UserMenu />
      </div>
    </header>
  );
}
