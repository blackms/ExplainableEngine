'use client';

import { usePathname } from 'next/navigation';
import { ChevronRight, Search } from 'lucide-react';
import { ThemeToggle } from './ThemeToggle';
import { UserMenu } from './UserMenu';

const routeLabels: Record<string, string> = {
  '/': 'Dashboard',
  '/audit': 'Explanations',
  '/monitor': 'Monitoring',
  '/playground': 'API Playground',
  '/settings': 'Settings',
  '/explain': 'Explanation',
};

function getBreadcrumbs(pathname: string) {
  if (pathname === '/') return [{ label: 'Dashboard', href: '/' }];

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
    <header className="flex h-14 shrink-0 items-center justify-between border-b border-border bg-background px-4 pl-14 lg:pl-4">
      {/* Left: Breadcrumbs */}
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

      {/* Right: Search hint + theme + user */}
      <div className="flex items-center gap-3">
        <button
          type="button"
          className="hidden items-center gap-2 rounded-md border border-border px-3 py-1.5 text-sm text-muted-foreground transition-colors hover:bg-accent sm:flex"
          aria-label="Search"
        >
          <Search className="size-3.5" />
          <span>Search</span>
          <kbd className="ml-1 rounded border border-border bg-muted px-1.5 py-0.5 text-[10px] font-medium text-muted-foreground">
            &#8984;K
          </kbd>
        </button>
        <ThemeToggle />
        <UserMenu />
      </div>
    </header>
  );
}
