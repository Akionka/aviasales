// Файл web\src\components\EntityInfo.jsx содержит код компонента, отвечающего за отображение информации о сущности
import { Button, Paper, Stack, Typography } from "@mui/material";

export const EntityInfo = ({ onDelete, items }) => {
  return (
    <Stack spacing={1}>
      <Paper elevation={0}>
        {items.map((item) => (
          <Typography key={item.label}>
            {item.label}: {item.value}
          </Typography>
        ))}
      </Paper>
      <Button variant="contained" onClick={onDelete}>
        Удалить
      </Button>
    </Stack>
  );
};
