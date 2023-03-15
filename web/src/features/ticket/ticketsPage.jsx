// Файл web\src\features\ticket\ticketsPage.jsx содержит код страницы с формой для таблицы "Билеты"
import SaveIcon from "@mui/icons-material/Save";
import AddIcon from "@mui/icons-material/Add";
import CancelIcon from "@mui/icons-material/Cancel";
import DeleteIcon from "@mui/icons-material/DeleteForeverOutlined";
import EditIcon from "@mui/icons-material/Edit";
import {
  Button,
  CircularProgress,
  Grid,
  Skeleton,
  TextField,
  Typography,
} from "@mui/material";
import {
  useCreateTicketMutation,
  useDeleteTicketByIDMutation,
  useGetTicketByIDQuery,
  useGetTicketsQuery,
  useUpdateTicketByIDMutation,
} from "../../app/services/api";
import { useEffect, useState } from "react";
import {
  DataGrid,
  GridActionsCellItem,
  GridRowModes,
  GridToolbarContainer,
} from "@mui/x-data-grid";
import { EntityInfo } from "../../components/EntityInfo";
import { useNavigate } from "react-router-dom";
import { formatPassportNumber } from "../../utils/formatters/passport";
import { localDatetimeToUTC } from "../../utils/dateConverter";
import { formatGender } from "../../utils/formatters/gender";
import { useAuth } from "../../hooks/useAuth";

const EditToolbar = ({ setRows, setRowModesModel }) => {
  const handleClick = () => {
    setRows((oldRows) => [
      ...oldRows,
      { id: 0, passenger_sex: 0, isNew: true },
    ]);
    setRowModesModel((oldModel) => ({
      ...oldModel,
      0: { mode: GridRowModes.Edit, fieldToFocus: "passenger_last_name" },
    }));
  };

  return (
    <GridToolbarContainer>
      <Button startIcon={<AddIcon />} onClick={handleClick}>
        Добавить запись
      </Button>
    </GridToolbarContainer>
  );
};

export const TicketsPage = () => {
  const navigate = useNavigate();
  const auth = useAuth()
  const [paginationModel, setPaginationModel] = useState({
    pageSize: 25,
    page: 0,
  });

  const [rowModesModel, setRowModesModel] = useState({});
  const [rows, setRows] = useState([]);

  const [searchQuery, setSearchQuery] = useState("");

  const { data, isLoading } = useGetTicketsQuery({
    page: paginationModel.page + 1,
    count: paginationModel.pageSize,
  });

  const {
    data: ticket,
    isLoading: isLoadingTicket,
    error,
  } = useGetTicketByIDQuery({ id: searchQuery });

  useEffect(() => {
    setRows((prevRows) => {
      const newRows = prevRows?.filter((r) => r.isNew) || [];
      if (data?.items) {
        return [...data?.items, ...newRows] || [];
      }
      return prevRows;
    });
  }, [data?.items]);

  const [updateTicket, { isLoadingUpdate }] = useUpdateTicketByIDMutation();
  const [deleteTicket, { isLoadingDelete }] = useDeleteTicketByIDMutation();
  const [createTicket, { isLoadingCreate }] = useCreateTicketMutation();

  const handleRowEditStart = (params, event) => {
    event.defaultMuiPrevented = true;
  };

  const handleRowEditStop = (params, event) => {
    event.defaultMuiPrevented = true;
  };

  const handleEditClick = (row) => () => {
    setRowModesModel({
      ...rowModesModel,
      [row.id]: { mode: GridRowModes.Edit },
    });
  };

  const handleSaveClick = (row) => () => {
    setRowModesModel({
      ...rowModesModel,
      [row.id]: { mode: GridRowModes.View },
    });
  };

  const handleDeleteClick = (row) => () => {
    deleteTicket({ id: row.id })
      .unwrap()
      .catch(({ data: { error } }) => alert(error));
  };

  const handleCancelClick = (row) => () => {
    setRowModesModel({
      ...rowModesModel,
      [row.id]: { mode: GridRowModes.View, ignoreModifications: true },
    });

    const editedRow = rows.find((r) => r.id === row.id);
    if (editedRow.isNew) {
      setRows(rows.filter((r) => r.id !== row.id));
    }
  };

  const handleSearchQueryChange = (e) => {
    setSearchQuery(e.target.value);
  };

  const columns = [
    {
      field: "id",
      headerName: "ID билета",
      width: 150,
      editable: false,
    },
    {
      field: "passenger_last_name",
      headerName: "Фамилия пассажира",
      width: 175,
      editable: true,
    },
    {
      field: "passenger_given_name",
      headerName: "Имя пассажира",
      width: 175,
      editable: true,
    },
    {
      field: "passenger_birth_date",
      headerName: "Дата рождения пассажира",
      width: 200,
      editable: true,
      type: "date",
      valueFormatter: ({ value }) =>
        value && localDatetimeToUTC(new Date(value)).toLocaleDateString(),
    },
    {
      field: "passenger_passport_number",
      headerName: "Номер паспорта пассажира",
      width: 210,
      editable: true,
      valueFormatter: ({ value }) => formatPassportNumber(value),
    },
    {
      field: "passenger_sex",
      headerName: "Пол пассажира",
      width: 150,
      editable: true,
      type: "singleSelect",
      valueOptions: [
        { value: 1, label: "Мужской" },
        { value: 2, label: "Женский" },
      ],
      valueFormatter: ({ value }) => formatGender(value),
    },
    {
      field: "purchase_id",
      headerName: "Номер покупки",
      width: 125,
      editable: true,
      type: "number",
    }]
    if (auth.user.role_id === 2) {
      columns.push({
        field: "actions",
        headerName: "Действия",
        type: "actions",
        width: 140,
        getActions: (row) => {
          const isInEditMode = rowModesModel[row.id]?.mode === GridRowModes.Edit;
          if (isInEditMode) {
            return [
              <GridActionsCellItem
                icon={<SaveIcon />}
                label="Save"
                onClick={handleSaveClick(row)}
              />,
              <GridActionsCellItem
                icon={<CancelIcon />}
                label="Cancel"
                className="textPrimary"
                onClick={handleCancelClick(row)}
                color="inherit"
              />,
            ];
          }
          return [
            <GridActionsCellItem
              icon={<EditIcon />}
              label="Edit"
              className="textPrimary"
              onClick={handleEditClick(row)}
              color="inherit"
            />,
            <GridActionsCellItem
              icon={<DeleteIcon />}
              label="Delete"
              onClick={handleDeleteClick(row)}
              color="inherit"
            />,
          ];
        },

      })
    }

  if (isLoading)
    return <Skeleton variant="rectangular" width={512} height={512} />;

  return (
    <Grid rowSpacing={3} columnSpacing={3} container>
      <Grid item xs={12}>
        <DataGrid
          autoHeight
          editMode="row"
          columns={columns}
          rows={rows}
          rowCount={data.total_count}
          pageSizeOptions={[5, 10, 15, 20, 25, 50, 100]}
          paginationModel={paginationModel}
          onPaginationModelChange={setPaginationModel}
          rowModesModel={rowModesModel}
          onRowModesModelChange={(newModel) => setRowModesModel(newModel)}
          onRowEditStart={handleRowEditStart}
          onRowEditStop={handleRowEditStop}
          components={{
            Toolbar: EditToolbar,
          }}
          componentsProps={{
            toolbar: { setRows, setRowModesModel },
          }}
          paginationMode="server"
          loading={
            isLoading || isLoadingUpdate || isLoadingDelete || isLoadingCreate
          }
          processRowUpdate={async (newRow, oldRow) => {
            const dateInUTC = new Date(
              Date.UTC(
                newRow.passenger_birth_date.getFullYear(),
                newRow.passenger_birth_date.getMonth(),
                newRow.passenger_birth_date.getDate(),
                0,
                0,
                0,
                0
              )
            );
            try {
              if (newRow.isNew) {
                const res = await createTicket({
                  ticket: { ...newRow, passenger_birth_date: dateInUTC },
                }).unwrap();
                setRows((prevRows) =>
                  prevRows.filter((row) => row.id !== oldRow.id)
                );
                return res;
              } else {
                const res = await updateTicket({
                  id: oldRow.id,
                  ticket: { ...newRow, passenger_birth_date: dateInUTC },
                }).unwrap();
                return res;
              }
            } catch (error) {
              const err_count = Object.keys(error.data.error).length
              throw new Error(`${err_count > 1 ? 'Поля' : 'Поле'} ${Object.keys(error.data.error).join(', ')} ${err_count > 1 ? 'заполнены' : 'заполнено'} неправильно`);
            }
          }}
          onProcessRowUpdateError={(error) => {
            alert(`{error.length > 1 ? 'Поля' : 'Поле'} {Object.keys(error).join(', '))} {error.length > 1 ? 'заполнены' : 'заполнено'} неправильно`)
          }}
        />
      </Grid>
      <Grid item xs={2}>
        <TextField
          variant="outlined"
          value={searchQuery}
          onChange={handleSearchQueryChange}
          size="small"
          label="ID билета"
          type={"number"}
        />
      </Grid>
      <Grid item xs={10}>
        {!isLoadingTicket && !error && ticket && searchQuery !== "" && (
          <Button
            variant="contained"
            onClick={() => navigate(`/ticket/${searchQuery}/report`)}
          >
            Просмотреть отчёт
          </Button>
        )}
      </Grid>
      <Grid item xs={5}>
        {isLoadingTicket && <CircularProgress />}
        {!isLoadingTicket && error && (
          <Typography>
            Ошибка! Код: {error?.status}. Сообщение: {error?.data?.error}
          </Typography>
        )}
        {!isLoadingTicket && !error && ticket && searchQuery !== "" && (
          <EntityInfo
            items={columns
              .filter((col) => col.type !== "actions")
              .map((col) => {
                return {
                  label: col.headerName,
                  value: (() => {
                    switch (col.field) {
                      // case "passenger_birth_date":
                      // return ticket[col.field] ? "Да" : "Нет";
                      case "passenger_birth_date":
                        return ticket[col.field]
                          ? new Date(ticket[col.field]).toLocaleDateString()
                          : "";
                      case "passenger_sex":
                        return formatGender(ticket[col.field]);
                      case "passenger_passport_number":
                        return formatPassportNumber(ticket[col.field]);
                      default:
                        return ticket[col.field];
                    }
                  })(),
                };
              })}
            onDelete={() =>
              deleteTicket({ id: ticket.id })
                .unwrap()
                .catch(({ data: { error } }) => alert(error))
            }
          />
        )}
      </Grid>
    </Grid>
  );
};
