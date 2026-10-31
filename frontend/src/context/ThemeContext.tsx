import React, {
  createContext,
  useContext,
  useEffect,
  useState,
  useCallback,
} from 'react';

export type Theme = 'light' | 'dark' | 'system';

interface ThemeContextType {
  theme: Theme;
  resolvedTheme: 'light' | 'dark';
  toggleTheme: () => void;
  setTheme: (theme: Theme) => void;
  isDark: boolean;
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

const THEME_STORAGE_KEY = 'brightloop_theme';

function getSystemTheme(): 'light' | 'dark' {
  if (
    typeof window !== 'undefined' &&
    window.matchMedia &&
    window.matchMedia('(prefers-color-scheme: dark)').matches
  ) {
    return 'dark';
  }
  return 'light';
}

export const ThemeProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [theme, setThemeState] = useState<Theme>(() => {
    try {
      const saved = localStorage.getItem(THEME_STORAGE_KEY) as Theme | null;
      if (saved === 'light' || saved === 'dark' || saved === 'system') {
        return saved;
      }
    } catch {
      // Storage unavailable
    }
    return 'system';
  });

  const [systemTheme, setSystemTheme] = useState<'light' | 'dark'>(getSystemTheme);

  // Resolved theme ('light' or 'dark') based on user choice or OS preference
  const resolvedTheme: 'light' | 'dark' =
    theme === 'system' ? systemTheme : theme;
  const isDark = resolvedTheme === 'dark';

  // Listen to OS preference changes in real-time
  useEffect(() => {
    if (typeof window === 'undefined' || !window.matchMedia) return;
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    const handler = (e: MediaQueryListEvent) => {
      setSystemTheme(e.matches ? 'dark' : 'light');
    };
    mediaQuery.addEventListener('change', handler);
    return () => mediaQuery.removeEventListener('change', handler);
  }, []);

  // Sync DOM with resolved theme
  useEffect(() => {
    const root = document.documentElement;
    if (resolvedTheme === 'dark') {
      root.classList.add('dark');
      root.setAttribute('data-theme', 'dark');
      root.style.colorScheme = 'dark';
    } else {
      root.classList.remove('dark');
      root.setAttribute('data-theme', 'light');
      root.style.colorScheme = 'light';
    }
  }, [resolvedTheme]);

  // Persist user selection
  useEffect(() => {
    try {
      localStorage.setItem(THEME_STORAGE_KEY, theme);
    } catch {
      // Storage unavailable
    }
  }, [theme]);

  // Execute theme change with View Transitions API (or CSS transition fallback)
  const executeThemeChange = useCallback((updateFn: () => void) => {
    if (typeof window === 'undefined' || typeof document === 'undefined') {
      updateFn();
      return;
    }

    const doc = document as unknown as {
      startViewTransition?: (cb: () => void) => void;
      documentElement: HTMLElement;
    };

    if (typeof doc.startViewTransition === 'function') {
      doc.startViewTransition(() => {
        updateFn();
      });
    } else {
      // Fallback: inject global temporary transition class
      const root = doc.documentElement;
      root.classList.add('theme-transition');
      updateFn();
      window.setTimeout(() => {
        root.classList.remove('theme-transition');
      }, 300);
    }
  }, []);

  const setTheme = useCallback(
    (newTheme: Theme) => {
      executeThemeChange(() => {
        setThemeState(newTheme);
      });
    },
    [executeThemeChange]
  );

  const toggleTheme = useCallback(() => {
    executeThemeChange(() => {
      setThemeState((prev) => {
        const current = prev === 'system' ? systemTheme : prev;
        return current === 'dark' ? 'light' : 'dark';
      });
    });
  }, [executeThemeChange, systemTheme]);

  return (
    <ThemeContext.Provider
      value={{
        theme,
        resolvedTheme,
        toggleTheme,
        setTheme,
        isDark,
      }}
    >
      {children}
    </ThemeContext.Provider>
  );
};

export function useTheme(): ThemeContextType {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error('useTheme must be used within a ThemeProvider');
  }
  return context;
}

