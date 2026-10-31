import React, { lazy } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider } from './context/ThemeContext';
import { AuthProvider } from './context/AuthContext';
import { ToastProvider } from './context/ToastContext';
import { DashboardLayout } from './components/layout/DashboardLayout';
import { AuthLayout } from './components/layout/AuthLayout';
import { ProtectedRoute } from './components/auth/ProtectedRoute';
import { LoginPage } from './pages/LoginPage';
import { RegisterPage } from './pages/RegisterPage';
import { ForgotPasswordPage } from './pages/ForgotPasswordPage';
import { DashboardPage } from './pages/DashboardPage';
import { APP_ROUTES, getDefaultRouteForRole } from './routes/routes.config';
import { useAuth } from './hooks/useAuth';

const QuoteBuilderPage = lazy(() => import('./pages/staff/QuoteBuilderPage'));

const RootRedirect: React.FC = () => {
  const { role, isAuthenticated } = useAuth();
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }
  return <Navigate to={getDefaultRouteForRole(role)} replace />;
};

export const App: React.FC = () => {
  return (
    <BrowserRouter>
      <ThemeProvider>
        <AuthProvider>
          <ToastProvider>
            <Routes>
              {/* Root Landing Redirect */}
              <Route path="/" element={<RootRedirect />} />

              {/* Public Authentication Routes (Centered Stripe Surface Layout) */}
              <Route element={<AuthLayout />}>
                <Route path="/login" element={<LoginPage />} />
                <Route path="/register" element={<RegisterPage />} />
                <Route path="/forgot-password" element={<ForgotPasswordPage />} />
              </Route>

              {/* Central /dashboard route mapping */}
              <Route path="/dashboard" element={<DashboardPage />} />

              {/* Admin compatibility redirects */}
              <Route path="/admin/dashboard" element={<Navigate to="/admin/finance" replace />} />
              <Route path="/admin/products" element={<Navigate to="/admin/catalog" replace />} />

              {/* Portal compatibility redirects */}
              <Route path="/portal/dashboard" element={<Navigate to="/customer/dashboard" replace />} />
              <Route path="/portal/subscriptions" element={<Navigate to="/customer/subscriptions" replace />} />
              <Route path="/portal/invoices" element={<Navigate to="/customer/invoices" replace />} />

              {/* Main Application Layout with Protected Routes */}
              <Route element={<DashboardLayout />}>
                {APP_ROUTES.map((route) => {
                  const Element = route.element;
                  return (
                    <Route
                      key={route.id}
                      path={route.path}
                      element={
                        <ProtectedRoute allowedRoles={route.allowedRoles}>
                          <Element />
                        </ProtectedRoute>
                      }
                    />
                  );
                })}
                {/* Standalone sub-routes */}
                <Route
                  path="/staff/quotations/new"
                  element={
                    <ProtectedRoute allowedRoles={['staff', 'admin']}>
                      <QuoteBuilderPage />
                    </ProtectedRoute>
                  }
                />
              </Route>

              {/* Fallback Catch-all Route */}
              <Route path="*" element={<RootRedirect />} />
            </Routes>
          </ToastProvider>
        </AuthProvider>
      </ThemeProvider>
    </BrowserRouter>
  );
};

export default App;
