'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useState, useEffect } from 'react';
import {
  Home,
  Layers,
  Activity,
  FlaskConical,
  Settings,
  Menu,
  X,
  ShieldCheck,
  User,
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { useSession } from 'next-auth/react';

interface NavItem {
  name: string;
  href: string;
  icon: React.ComponentType<{ className?: string }>;
}

interface NavGroup {
  label: string;
  items: NavItem[];
}

const navigationGroups: NavGroup[] = [
  {
    label: 'Main',
    items: [
      { name: 'Dashboard', href: '/', icon: Home },
      { name: 'Explanations', href: '/audit', icon: Layers },
    ],
  },
  {
    label: 'Operations',
    items: [
      { name: 'Monitoring', href: '/monitor', icon: Activity },
    ],
  },
  {
    label: 'Developer',
    items: [
      { name: 'API Playground', href: '/playground', icon: FlaskConical },
    ],
  },
];

const bottomNav: NavItem[] = [
  { name: 'Settings', href: '/settings', icon: Settings },
];

export function Sidebar() {
  const pathname = usePathname();
  const { data: session } = useSession();
  const [mobileOpen, setMobileOpen] = useState(false);

  // Close mobile sidebar on route change
  useEffect(() => {
    setMobileOpen(false);
  }, [pathname]);

  // Lock body scroll when mobile sidebar is open
  useEffect(() => {
    if (mobileOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }
    return () => {
      document.body.style.overflow = '';
    };
  }, [mobileOpen]);

  function isActive(href: string) {
    if (href === '/') return pathname === '/';
    return pathname.startsWith(href);
  }

  return (
    <>
      {/* Mobile hamburger */}
      <button
        type="button"
        className="fixed top-3.5 left-3 z-50 rounded-md border border-border bg-background p-1.5 lg:hidden"
        onClick={() => setMobileOpen(!mobileOpen)}
        aria-label={mobileOpen ? 'Close menu' : 'Open menu'}
      >
        {mobileOpen ? <X className="size-5" /> : <Menu className="size-5" />}
      </button>

      {/* Backdrop */}
      {mobileOpen && (
        <div
          className="fixed inset-0 z-40 bg-black/50 backdrop-blur-sm lg:hidden"
          onClick={() => setMobileOpen(false)}
        />
      )}

      {/* Sidebar */}
      <aside
        className={cn(
          'fixed inset-y-0 left-0 z-40 flex w-60 flex-col border-r border-sidebar-border bg-sidebar transition-transform duration-250 ease-out lg:static lg:translate-x-0',
          mobileOpen ? 'translate-x-0' : '-translate-x-full'
        )}
      >
        {/* Logo */}
        <div className="flex h-14 items-center gap-2.5 border-b border-sidebar-border px-4 py-5">
          <div className="flex size-7 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
            <ShieldCheck className="size-4" />
          </div>
          <span className="text-sm font-semibold text-sidebar-foreground">
            Explainable Engine
          </span>
        </div>

        {/* Navigation groups */}
        <nav className="flex-1 overflow-y-auto px-3 py-4">
          {navigationGroups.map((group, groupIdx) => (
            <div key={group.label} className={cn(groupIdx > 0 && 'mt-6')}>
              <p className="mb-1 px-3 text-xs font-medium uppercase tracking-wider text-muted-foreground">
                {group.label}
              </p>
              <div className="space-y-1">
                {group.items.map((item) => {
                  const active = isActive(item.href);
                  return (
                    <Link
                      key={item.name}
                      href={item.href}
                      className={cn(
                        'flex items-center gap-3 rounded-md px-3 py-2 text-sm transition-colors',
                        active
                          ? 'border-l-2 border-primary bg-primary/10 font-medium text-primary'
                          : 'text-muted-foreground hover:bg-accent/50 hover:text-sidebar-foreground'
                      )}
                    >
                      <item.icon className="size-[18px] shrink-0" />
                      {item.name}
                    </Link>
                  );
                })}
              </div>
            </div>
          ))}
        </nav>

        {/* Bottom section */}
        <div className="border-t border-sidebar-border px-3 py-3">
          {/* Settings */}
          {bottomNav.map((item) => {
            const active = isActive(item.href);
            return (
              <Link
                key={item.name}
                href={item.href}
                className={cn(
                  'flex items-center gap-3 rounded-md px-3 py-2 text-sm transition-colors',
                  active
                    ? 'border-l-2 border-primary bg-primary/10 font-medium text-primary'
                    : 'text-muted-foreground hover:bg-accent/50 hover:text-sidebar-foreground'
                )}
              >
                <item.icon className="size-[18px] shrink-0" />
                {item.name}
              </Link>
            );
          })}

          {/* User info */}
          {session?.user && (
            <div className="mt-3 flex items-center gap-3 rounded-md px-3 py-2">
              <div className="flex size-8 items-center justify-center rounded-full bg-muted text-muted-foreground">
                <User className="size-4" />
              </div>
              <div className="min-w-0 flex-1">
                <p className="truncate text-sm font-medium text-sidebar-foreground">
                  {session.user.name}
                </p>
                <p className="truncate text-xs text-muted-foreground">
                  {session.user.email}
                </p>
              </div>
            </div>
          )}
        </div>
      </aside>
    </>
  );
}
