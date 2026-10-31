import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import { ToastProvider } from './context/ToastContext';
import { DashboardLayout } from './components/layout/DashboardLayout';
import { ProtectedRoute } from './components/auth/ProtectedRoute';
import { APP_ROUTES, getDefaultRouteForRole } from './routes/routes.config';
import { useAuth } from './hooks/useAuth';

const RootRedirect: React.FC = () => {
  const { role, isAuthenticated } = useAuth();
  if (!isAuthenticated) {
    return <Navigate to="/admin/dashboard" replace />;
  }
  return <Navigate to={getDefaultRouteForRole(role)} replace />;
};

export const App: React.FC = () => {
  return (
    <BrowserRouter>
      <AuthProvider>
        <ToastProvider>
          <Routes>
            {/* Root Landing Redirect */}
            <Route path="/" element={<RootRedirect />} />

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
            </Route>

            {/* Fallback Catch-all Route */}
            <Route path="*" element={<RootRedirect />} />
          </Routes>
        </ToastProvider>
      </AuthProvider>
    </BrowserRouter>
  );
};

export default App;
