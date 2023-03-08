// Файл web\src\features\mainMenu\MainMenu.jsx содержит код компонента "Главное меню"
import * as React from "react";
import Box from "@mui/material/Box";
import Drawer from "@mui/material/Drawer";
import AppBar from "@mui/material/AppBar";
import Toolbar from "@mui/material/Toolbar";
import List from "@mui/material/List";
import Typography from "@mui/material/Typography";
import Divider from "@mui/material/Divider";
import ListItem from "@mui/material/ListItem";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemText from "@mui/material/ListItemText";
import Button  from "@mui/material/Button";
import { Outlet, useNavigate } from "react-router-dom";
import { useAuth } from "../../hooks/useAuth";
import { useDispatch } from "react-redux";
import { setCredentials } from "../auth/authSlice";

const drawerWidth = 240;

const routes = [
  "airport",
  "booking_office",
  "cashier",
  "flight",
  "flight_in_ticket",
  "line",
  "liner_model",
  "liner",
  "purchase",
  "seat",
  "ticket",
];

export const MainMenu = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch()
  const auth = useAuth();

  const handleLogout = (e) => {
    dispatch(setCredentials({user: null, token: null}))
  }

  return (
    <Box sx={{ display: "flex", flexGrow: 1}}>
      <AppBar
        position="fixed"
        sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}
      >
        <Toolbar>
          <Typography variant="h6" noWrap component="div" sx={{flexGrow: 1}}>
            Журавли
          </Typography>
          <Box sx={{display: "flex", gap: 1}}>
            <Typography variant="h6">Здравствуйте, {auth.user.last_name} {auth.user.first_name} {auth.user.middle_name}</Typography>
            <Button color="inherit" onClick={handleLogout}>Выйти</Button>
          </Box>
        </Toolbar>
      </AppBar>
      <Drawer
        variant="permanent"
        sx={{
          width: drawerWidth,
          flexShrink: 0,
          [`& .MuiDrawer-paper`]: {
            width: drawerWidth,
            boxSizing: "border-box",
          },
        }}
      >
        <Toolbar />
        <Box sx={{ overflow: "auto" }}>
          <List>
            {[
              "Аэропорты",
              "Кассы",
              "Кассиры",
              "Полёты",
              "Полёты в билете",
              "Рейсы",
              "Модели самолётов",
              "Самолёты",
              "Покупки",
              "Места",
              "Билеты",
            ].map((text, index) => (
              <ListItem key={text} disablePadding>
                <ListItemButton>
                  <ListItemText
                    primary={text}
                    onClick={() => navigate(routes[index])}
                  />
                </ListItemButton>
              </ListItem>
            ))}
          </List>
          <Divider />
          <ListItem key={"Генерация билета"} disablePadding>
            <ListItemButton>
              <ListItemText
                primary={"Генерация билета"}
                onClick={() => navigate("/ticket/1/report")}
              />
            </ListItemButton>
          </ListItem>
        </Box>
      </Drawer>
      <Box component="main" sx={{ flexGrow: 1, p: 3 }}>
        <Toolbar />
        <Outlet />
      </Box>
    </Box>
  );
};
