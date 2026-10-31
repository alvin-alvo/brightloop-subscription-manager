import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { getDefaultRouteForRole } from '../routes/routes.config';

export const DashboardPage: React.FC = () => {
  const { user, role, isAuthenticated } = useAuth();

  if (!isAuthenticated || !user) {
    return <Navigate to="/login" replace />;
  }

  const destination = getDefaultRouteForRole(role);
  return <Navigate to={destination} replace />;
};

export default DashboardPage;

