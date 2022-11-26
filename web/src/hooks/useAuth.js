// Файл web\src\hooks\useAuth.js содержит код хука, который предоставляет информацию об аутентификации
import { useMemo } from "react";
import { useSelector } from "react-redux";
import { selectCurrentUser } from "../features/auth/authSlice";

export const useAuth = () => {
  const user = useSelector(selectCurrentUser);
  return useMemo(() => ({ user }), [user]);
};
