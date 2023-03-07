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
  } = useForm();

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
      console.log(err)
      if (err.status === 401) {
        setError("login", {
          type: "custom",
          message: "Неправильный логин или пароль!",
        });
        setError("password", {
          type: "custom",
          message: "Неправильный логин или пароль!",
        });
      } else {
        setError("login", {
          type: "custom",
          message: "Случилось что-то непредвиденное!",
        });
        setError("password", {
          type: "custom",
          message: `Случилось что-то непредвиденное!`,
        });
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
                    error={errors.password && errors.middle_name}
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
                    error={errors.password && errors.login}
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
                    error={errors.password && errors.login}
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
            Войти
          </Button>
          {errors.password && errors.login && (
            <Alert severity="error">{errors.password.message}</Alert>
          )}
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
