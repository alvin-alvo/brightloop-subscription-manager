import React from 'react';
import { Sun, Moon } from 'lucide-react';
import { useTheme } from '../context/ThemeContext';

interface ThemeToggleProps {
  className?: string;
}

export const ThemeToggle: React.FC<ThemeToggleProps> = ({ className = '' }) => {
  const { isDark, toggleTheme } = useTheme();

  return (
    <button
      onClick={toggleTheme}
      type="button"
      className={`relative inline-flex items-center justify-center w-9 h-9 rounded-full text-slate-500 hover:text-stripe-ink dark:text-slate-400 dark:hover:text-white hover:bg-slate-100 dark:hover:bg-slate-800 border border-transparent hover:border-slate-200 dark:hover:border-slate-700 transition-colors focus:outline-none focus:ring-2 focus:ring-stripe-primary/20 ${className}`}
      aria-label={isDark ? 'Switch to light mode' : 'Switch to dark mode'}
      title={isDark ? 'Switch to light mode' : 'Switch to dark mode'}
    >
      <div className="relative w-5 h-5 flex items-center justify-center overflow-hidden">
        {/* Sun Icon (Active in Light Mode) */}
        <Sun
          className={`w-4 h-4 text-amber-500 absolute transform transition-all duration-500 ease-[cubic-bezier(0.22,1,0.36,1)] ${
            isDark
              ? 'rotate-90 scale-0 opacity-0 pointer-events-none'
              : 'rotate-0 scale-100 opacity-100'
          }`}
        />

        {/* Moon Icon (Active in Dark Mode) */}
        <Moon
          className={`w-4 h-4 text-indigo-400 absolute transform transition-all duration-500 ease-[cubic-bezier(0.22,1,0.36,1)] ${
            isDark
              ? 'rotate-0 scale-100 opacity-100'
              : '-rotate-90 scale-0 opacity-0 pointer-events-none'
          }`}
        />
      </div>
    </button>
  );
};

export default ThemeToggle;

