// Файл web\src\hoc\ProtectedRoute.js содержит код для компоненты высшего порядка, которая позволяет запретить неавторизованный доступ
import { Navigate, useLocation } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";

export const ProtectedRoute = ({ children }) => {
  const auth = useAuth();
  const location = useLocation();

  if (!auth.user) {
    return <Navigate to="/login" state={{ from: location }} />;
  }

  return children;
};
