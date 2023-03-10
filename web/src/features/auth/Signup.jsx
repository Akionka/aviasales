// Файл web\src\features\auth\Login.jsx содержит код формы авторизации
import { Controller, useForm } from "react-hook-form";
import { useDispatch, useSelector } from "react-redux";
import { Navigate, useLocation, useNavigate } from "react-router-dom";
import { useLoginMutation, useSignupMutation, useWhoamiQuery } from "../../app/services/api";
import { useAuth } from "../../hooks/useAuth";
import { setCredentials } from "./authSlice";

import Avatar from "@mui/material/Avatar";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import Box from "@mui/material/Box";
import Grid from "@mui/material/Grid";
import LockOutlinedIcon from "@mui/icons-material/LockOutlined";
import Typography from "@mui/material/Typography";
import Container from "@mui/material/Container";
import Alert from "@mui/material/Alert";
import { useEffect } from "react";
import { Link } from "@mui/material";
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';

const makeNameString = (name, min, max) => {
  return z.string()
  .min(min, `Поле ${name} должно содержать минимум ${min} символа`)
  .max(max, `Поле ${name} может содержать максимум ${max} символа`)
}

const includesAny = (str, chars) => {
  console.log(chars.length)
  for (let i = 0; i < chars.length; i++) {
    if (str.includes(chars[i])) return true
  }
  return false
}

const signupSchema = z.object({
  first_name: makeNameString('Имя', 3, 64),
  last_name: makeNameString('Фамилия', 3, 64),
  middle_name: z.union([makeNameString('Отчество', 3, 64), z.literal('')]),
  login: makeNameString('Логин', 3, 32).regex(/^[a-zA-Z0-9]+$/, `Поле Логин может содержать только латинские заглавные и строчные буквы и цифры`),
  password: makeNameString('Пароль', 4, 16)
  .regex(/\d/, {message: "Пароль должен содержать хотя бы одну цифру"})
  .refine(str => str.toLowerCase() !== str, {message: "Пароль должен содержать хотя бы одну заглавную букву"})
  .refine(str => !includesAny(str, '*&{}|+'), {message: "Пароль не должен содержать запрещённые символы: *&{}|+"})
})

export const Signup = () => {
  const dispatch = useDispatch();
  const natigate = useNavigate();
  const [signup, { isLoading: isLoadingSignup }] = useSignupMutation();
  const [login, { isLoading: isLoadingLogin }] = useLoginMutation();
  const auth = useAuth();
  const token = useSelector((state) => state.auth.token);
  const location = useLocation();
  const { refetch } = useWhoamiQuery();

  const fromPage = location.state?.from?.pathname || "/";

  const {
    handleSubmit,
    formState: { errors },
    setError,
    control,
  } = useForm({resolver: zodResolver(signupSchema)});

  useEffect(() => {
    if (token) {
      refetch()
        .unwrap()
        .then((res) => dispatch(setCredentials({ user: res, token: token })));
    }
  }, [dispatch, refetch, token]);

  const onSubmit = async (data) => {
    try {
      await signup({
        first_name: data.first_name,
        last_name: data.last_name,
        middle_name: data.middle_name,
        login: data.login,
        password: data.password,
      }).unwrap();
      const user = await login({
        login: data.login,
        password: data.password,
      }).unwrap();
      natigate('/', { replace: true });
      dispatch(setCredentials(user));
      sessionStorage.setItem("auth_token", user.token);
    } catch (err) {
      switch (err.status) {
        case 400:
          for (const key in err.data.error) {
            if (Object.hasOwnProperty.call(err.data.error, key)) {
              const element = err.data.error[key];
              setError(key, {
                type: "custom",
                message: element
              })
            }
          }
          break
        default:
          setError("serverError", {
            type: "custom",
            message: "Произошла серверная ошибка. Обратитесь к администратору",
          });
          break;
      }
    }
  };

  if (auth.user) {
    return <Navigate to={fromPage} replace={true} />;
  }
  return (
    <Container component="main" maxWidth="sm">
      <Box
        sx={{
          marginTop: 8,
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
        }}
      >
        <Avatar sx={{ m: 1, bgcolor: "primary.main" }}>
          <LockOutlinedIcon />
        </Avatar>
        <Typography component="h1" variant="h5">
          Регистрация
        </Typography>
        <Box component="form" onSubmit={handleSubmit(onSubmit)} sx={{ mt: 1 }}>
          <Grid container spacing={2}>
            <Grid item xs={12} sm={4}>
              <Controller
                name="first_name"
                control={control}
                defaultValue={""}
                render={({ field: { onChange, onBlur, value, ref } }) => (
                  <TextField
                    margin="normal"
                    required
                    fullWidth
                    id="first_name"
                    label="Имя"
                    name="first_name"
                    autoComplete="first_name"
                    autoFocus
                    error={errors.first_name}
                    value={value}
                    onChange={onChange}
                    onBlur={onBlur}
                    inputRef={ref}
                  />
                )}
              />
            </Grid>
            <Grid item xs={12} sm={4}>
              <Controller
                name="last_name"
                control={control}
                defaultValue={""}
                render={({ field: { onChange, onBlur, value, ref } }) => (
                  <TextField
                    margin="normal"
                    required
                    fullWidth
                    id="last_name"
                    label="Фамилия"
                    name="last_name"
                    autoComplete="last_name"
                    error={errors.last_name}
                    value={value}
                    onChange={onChange}
                    onBlur={onBlur}
                    inputRef={ref}
                  />
                )}
              />
            </Grid>
            <Grid item xs={12} sm={4}>
              <Controller
                name="middle_name"
                control={control}
                defaultValue={""}
                render={({ field: { onChange, onBlur, value, ref } }) => (
                  <TextField
                    margin="normal"
                    fullWidth
                    id="middle_name"
                    label="Отчество"
                    name="middle_name"
                    autoComplete="middle_name"
                    error={errors.middle_name}
                    value={value}
                    onChange={onChange}
                    onBlur={onBlur}
                    inputRef={ref}
                  />
                )}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <Controller
                name="login"
                control={control}
                defaultValue={""}
                render={({ field: { onChange, onBlur, value, ref } }) => (
                  <TextField
                    margin="normal"
                    required
                    fullWidth
                    id="login"
                    label="Логин"
                    name="login"
                    autoComplete="login"
                    error={errors.login}
                    value={value}
                    onChange={onChange}
                    onBlur={onBlur}
                    inputRef={ref}
                  />
                )}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <Controller
                name="password"
                control={control}
                defaultValue={""}
                render={({ field: { onChange, onBlur, value, ref } }) => (
                  <TextField
                    margin="normal"
                    required
                    fullWidth
                    id="password"
                    label="Пароль"
                    name="password"
                    autoComplete="password"
                    type="password"
                    error={errors.password}
                    value={value}
                    onChange={onChange}
                    onBlur={onBlur}
                    inputRef={ref}
                  />
                )}
              />
            </Grid>
          </Grid>
          <Button
            type="submit"
            fullWidth
            variant="contained"
            sx={{ mt: 3, mb: 2 }}
            disabled={isLoadingSignup || isLoadingLogin}
          >
            Зарегистрироваться
          </Button>
          {errors.first_name?.message && <Alert severity="error">{errors.first_name.message}</Alert>}
          {errors.last_name?.message && <Alert severity="error">{errors.last_name.message}</Alert>}
          {errors.middle_name?.message && <Alert severity="error">{errors.middle_name.message}</Alert>}
          {errors.login?.message && <Alert severity="error">{errors.login.message}</Alert>}
          {errors.password?.message && <Alert severity="error">{errors.password.message}</Alert>}
          {errors.credentials?.message && <Alert severity="error">{errors.credentials.message}</Alert>}
          {errors.serverError?.message && <Alert severity="error">{errors.serverError.message}</Alert>}
          <Grid container justify-content="flex-end">
            <Grid item>
              <Link variant="body2" href="/signup">
                Авторизация
              </Link>
            </Grid>
          </Grid>
        </Box>
      </Box>
    </Container>
  );
};
