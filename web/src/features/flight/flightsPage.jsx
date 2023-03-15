// Файл web\src\features\flight\flightsPage.jsx содержит код страницы с формой для таблицы "Полёт"
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
  useCreateFlightMutation,
  useDeleteFlightByIDMutation,
  useGetFlightByIDQuery,
  useGetFlightsQuery,
  useUpdateFlightByIDMutation,
} from "../../app/services/api";
import { useEffect, useState } from "react";
import {
  DataGrid,
  GridActionsCellItem,
  GridRowModes,
  GridToolbarContainer,
} from "@mui/x-data-grid";
import { EntityInfo } from "../../components/EntityInfo";
import { localDatetimeToUTC } from "../../utils/dateConverter";
import { useAuth } from "../../hooks/useAuth";

const EditToolbar = ({ setRows, setRowModesModel }) => {
  const handleClick = () => {
    setRows((oldRows) => [
      ...oldRows,
      {
        id: 0,
        dep_date: new Date(),
        is_hot: false,
        line_code: "JR",
        liner_code: "RA",
        isNew: true,
      },
    ]);
    setRowModesModel((oldModel) => ({
      ...oldModel,
      0: { mode: GridRowModes.Edit, fieldToFocus: "dep_date" },
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

export const FlightsPage = () => {
  const auth = useAuth()
  const [paginationModel, setPaginationModel] = useState({
    pageSize: 25,
    page: 0,
  });

  const [rowModesModel, setRowModesModel] = useState({});
  const [rows, setRows] = useState([]);

  const [searchQuery, setSearchQuery] = useState("");

  const { data, isLoading } = useGetFlightsQuery({
    page: paginationModel.page + 1,
    count: paginationModel.pageSize,
  });

  const {
    data: flight,
    isLoading: isLoadingFlight,
    error,
  } = useGetFlightByIDQuery({ id: searchQuery });

  useEffect(() => {
    setRows((prevRows) => {
      const newRows = prevRows?.filter((r) => r.isNew) || [];
      if (data?.items) {
        return [...data?.items, ...newRows] || [];
      }
      return prevRows;
    });
  }, [data?.items]);

  const [updateFlight, { isLoadingUpdate }] = useUpdateFlightByIDMutation();
  const [deleteFlight, { isLoadingDelete }] = useDeleteFlightByIDMutation();
  const [createFlight, { isLoadingCreate }] = useCreateFlightMutation();

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
    deleteFlight({ id: row.id })
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
      headerName: "ID",
      width: 100,
      editable: false,
      type: "number",
    },
    {
      field: "dep_date",
      headerName: "Дата отправления",
      width: 200,
      editable: true,
      type: "date",
      valueGetter: ({ value }) => value && new Date(value),
      valueFormatter: ({ value }) =>
        value && localDatetimeToUTC(new Date(value)).toLocaleDateString(),
    },
    {
      field: "is_hot",
      headerName: "Горячий рейс?",
      type: "boolean",
      width: 125,
      editable: true,
    },
    {
      field: "line_code",
      headerName: "Код рейса",
      width: 200,
      editable: true,
    },
    {
      field: "liner_code",
      headerName: "Код самолёта",
      width: 200,
      editable: true,
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
                newRow.dep_date.getFullYear(),
                newRow.dep_date.getMonth(),
                newRow.dep_date.getDate(),
                0,
                0,
                0,
                0
              )
            );
            try {
              if (newRow.isNew) {
                const res = await createFlight({
                  flight: { ...newRow, dep_date: dateInUTC.toISOString() },
                }).unwrap();
                setRows((prevRows) =>
                  prevRows.filter((row) => row.id !== oldRow.id)
                );
                return res;
              } else {
                const res = await updateFlight({
                  id: oldRow.id,
                  flight: { ...newRow, dep_date: dateInUTC.toISOString() },
                }).unwrap();
                return res;
              }
            } catch (error) {
              const err_count = Object.keys(error.data.error).length
              throw new Error(`${err_count > 1 ? 'Поля' : 'Поле'} ${Object.keys(error.data.error).join(', ')} ${err_count > 1 ? 'заполнены' : 'заполнено'} неправильно`);
            }
          }}
          onProcessRowUpdateError={(error) => {
            alert(error);
          }}
        />
      </Grid>
      <Grid item xs={12}>
        <TextField
          variant="outlined"
          value={searchQuery}
          onChange={handleSearchQueryChange}
          size="small"
          label="ID полёта"
          type={"number"}
        />
      </Grid>
      <Grid item xs={5}>
        {isLoadingFlight && <CircularProgress />}
        {!isLoadingFlight && error && (
          <Typography>
            Ошибка! Код: {error?.status}. Сообщение: {error?.data?.error}
          </Typography>
        )}
        {!isLoadingFlight && !error && flight && searchQuery !== "" && (
          <EntityInfo
            items={columns
              .filter((col) => col.type !== "actions")
              .map((col) => {
                return {
                  label: col.headerName,
                  value: (() => {
                    switch (col.type) {
                      case "boolean":
                        return flight[col.field] ? "Да" : "Нет";
                      case "date":
                        return flight[col.field]
                          ? new Date(flight[col.field]).toLocaleDateString()
                          : "";
                      default:
                        return flight[col.field];
                    }
                  })(),
                };
              })}
            onDelete={() =>
              deleteFlight({ id: flight.id })
                .unwrap()
                .catch(({ data: { error } }) => alert(error))
            }
          />
        )}
      </Grid>
    </Grid>
  );
};
