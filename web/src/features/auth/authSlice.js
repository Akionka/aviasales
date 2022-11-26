// Файл web\src\features\auth\authSlice.js содержит код, создающий срез для работы с аутентификацией
import { createSlice } from "@reduxjs/toolkit";

const initialState = {
  token: null,
  user: null,
};

export const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    setCredentials: (state, { payload: { user, token } }) => {
      state.user = user;
      state.token = token;
    },
    setToken: (state, { payload: token }) => {
      state.token = token;
    },
  },
});

export const { setCredentials, setToken } = authSlice.actions;
export default authSlice.reducer;
export const selectCurrentUser = (state) => state.auth.user;
