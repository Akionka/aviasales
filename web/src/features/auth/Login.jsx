import { Controller, useForm } from "react-hook-form";
import { useDispatch } from "react-redux";
import { Navigate, useLocation, useNavigate } from "react-router-dom";
import { useLoginMutation } from "../../app/services/api";
import { useAuth } from "../../hooks/useAuth";
import { setCredentials } from "./authSlice";

import Avatar from "@mui/material/Avatar";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import Box from "@mui/material/Box";
import LockOutlinedIcon from "@mui/icons-material/LockOutlined";
import Typography from "@mui/material/Typography";
import Container from "@mui/material/Container";
import Alert from "@mui/material/Alert";

export const Login = () => {
  const dispatch = useDispatch();
  const natigate = useNavigate();
  const [login, { isLoading }] = useLoginMutation();
  const auth = useAuth();
  const location = useLocation();

  const fromPage = location.state?.from?.pathname || "/";

  const {
    handleSubmit,
    formState: { errors },
    setError,
    control,
  } = useForm();

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
    return <Navigate to={fromPage} />;
  }

  return (
    <Container component="main" maxWidth="xs">
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
            render={({ field: { onChange, onBlur, value, name, ref } }) => (
              <TextField
                margin="normal"
                required
                fullWidth
                id="login"
                label="Логин"
                name="login"
                autoComplete="login"
                autoFocus
                error={errors.password && errors.login}
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
            render={({ field: { onChange, onBlur, value, name, ref } }) => (
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
          <Button
            type="submit"
            fullWidth
            variant="contained"
            sx={{ mt: 3, mb: 2 }}
            disabled={isLoading}
          >
            Войти
          </Button>
          {errors.password && errors.login && (
            <Alert severity="error">{errors.password.message}</Alert>
          )}
        </Box>
      </Box>
    </Container>
  );
};
