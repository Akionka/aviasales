// Файл web\src\features\auth\Login.jsx содержит код формы авторизации
import { Controller, useForm } from "react-hook-form";
import { useDispatch, useSelector } from "react-redux";
import { Navigate, useLocation, useNavigate } from "react-router-dom";
import { useLoginMutation, useWhoamiQuery } from "../../app/services/api";
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

export const Login = () => {
  const dispatch = useDispatch();
  const natigate = useNavigate();
  const [login, { isLoading }] = useLoginMutation();
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
    clearErrors
  } = useForm({reValidateMode: "onChange"});

  useEffect(() => {
    if (token) {
      refetch()
        .unwrap()
        .then((res) => dispatch(setCredentials({ user: res, token: token })));
    }
  }, [dispatch, refetch, token]);

  const onSubmit = async (data) => {
    try {
      const user = await login({
        login: data.login,
        password: data.password,
      }).unwrap();
      natigate(fromPage, { replace: true });
      dispatch(setCredentials(user));
      sessionStorage.setItem("auth_token", user.token);
    } catch (err) {
      if (err.status === 401) {
        setError("credentials", {
          type: "custom",
          message: "Неправильный логин или пароль!",
        });
      } else {
        setError("serverError", {
          type: "custom",
          message: "Случилось что-то непредвиденное!",
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
          Вход
        </Typography>
        <Box component="form" onSubmit={handleSubmit(onSubmit)} sx={{ mt: 1 }}>
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
                autoFocus
                error={errors.credentials}
                value={value}
                onChange={onChange}
                onBlur={onBlur}
                inputRef={ref}
              />
            )}
          />
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
                  error={errors.credentials}
                  value={value}
                  onChange={onChange}
                  onBlur={onBlur}
                  inputRef={ref} />
              )
            }
          />
          <Button
            type="submit"
            fullWidth
            variant="contained"
            sx={{ mt: 3, mb: 2 }}
            disabled={isLoading}
            onClick={() => clearErrors()}
          >
            Войти
          </Button>
          {errors.credentials?.message && <Alert severity="error">{errors.credentials.message}</Alert>}
          {errors.serverError?.message && <Alert severity="error">{errors.serverError.message}</Alert>}
          <Grid container justify-content="flex-end">
            <Grid item>
              <Link variant="body2" href="/signup">
                Регистрация
              </Link>
            </Grid>
          </Grid>
        </Box>
      </Box>
    </Container>
  );
};
