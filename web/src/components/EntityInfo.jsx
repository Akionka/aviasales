import { Button, Paper, Stack, Typography } from "@mui/material";

export const EntityInfo = ({ onDelete, items }) => {
  return (
    <Stack spacing={1}>
      <Paper elevation={0}>
        {items.map((item) => (
          <Typography>
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
