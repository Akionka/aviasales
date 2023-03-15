// Файл web\src\features\flightInTicket\flightInTicketsPage.jsx содержит код страницы с формой для таблицы "Полёт в билете"
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
  useCreateFlightInTicketMutation,
  useDeleteFlightInTicketByIDMutation,
  useGetFlightInTicketByIDQuery,
  useGetFlightInTicketsQuery,
  useGetFlightsQuery,
  useGetSeatsQuery,
  useUpdateFlightInTicketByIDMutation,
} from "../../app/services/api";
import { useEffect, useState } from "react";
import {
  DataGrid,
  GridActionsCellItem,
  GridRowModes,
  GridToolbarContainer,
} from "@mui/x-data-grid";
import { EntityInfo } from "../../components/EntityInfo";
import { useAuth } from "../../hooks/useAuth";

const EditToolbar = ({ setRows, setRowModesModel }) => {
  const handleClick = () => {
    setRows((oldRows) => [
      ...oldRows,
      { id: 0, flight_id: 0, seat_id: 0, ticket_id: 0, isNew: true },
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

export const FlightInTicketsPage = () => {
  const auth = useAuth()
  const [paginationModel, setPaginationModel] = useState({
    pageSize: 25,
    page: 0,
  });

  const [rowModesModel, setRowModesModel] = useState({});
  const [rows, setRows] = useState([]);

  const [searchQuery, setSearchQuery] = useState("");

  const { data, isLoading } = useGetFlightInTicketsQuery({
    page: paginationModel.page + 1,
    count: paginationModel.pageSize,
  });

  const {data: flights, isLoading: isLoadingFlights} = useGetFlightsQuery({page: 0, count: -1})
  const {data: seats, isLoading: isLoadingSeats} = useGetSeatsQuery({page: 0, count: -1})

  const {
    data: flightInTicket,
    isLoading: isLoadingFlightInTicket,
    error,
  } = useGetFlightInTicketByIDQuery({ id: searchQuery });

  useEffect(() => {
    setRows((prevRows) => {
      const newRows = prevRows?.filter((r) => r.isNew) || [];
      if (data?.items) {
        return [...data?.items, ...newRows] || [];
      }
      return prevRows;
    });
  }, [data?.items]);

  const [updateFlightInTicket, { isLoadingUpdate }] =
    useUpdateFlightInTicketByIDMutation();
  const [deleteFlightInTicket, { isLoadingDelete }] =
    useDeleteFlightInTicketByIDMutation();
  const [createFlightInTicket, { isLoadingCreate }] =
    useCreateFlightInTicketMutation();

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
    deleteFlightInTicket({ id: row.id })
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
console.log(flights)
  const columns = [
    {
      field: "id",
      headerName: "ID",
      width: 100,
      editable: false,
      type: "number",
    },
    {
      field: "flight_id",
      headerName: "Полёт",
      width: 200,
      editable: true,
      type: "singleSelect",
      valueOptions: flights?.items?.map(f => {return {value: f.id, label: `${f.line_code}[${new Date(f.dep_date).toLocaleDateString()}]`}})
    },
    {
      field: "ticket_id",
      headerName: "Номер билета",
      width: 125,
      editable: true,
      type: "number",
    },
    {
      field: "seat_id",
      headerName: "Место",
      width: 200,
      editable: true,
      type: "singleSelect",
      valueOptions: seats?.items?.map(s => {return {value: s.id, label: `${s.model_code} - ${s.number} - ${s.class}`}})
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
            try {
              if (newRow.isNew) {
                const res = await createFlightInTicket({
                  flight_in_ticket: newRow,
                }).unwrap();
                setRows((prevRows) =>
                  prevRows.filter((row) => row.id !== oldRow.id)
                );
                return res;
              } else {
                const res = await updateFlightInTicket({
                  id: oldRow.id,
                  flight_in_ticket: newRow,
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
          label="ID полёта в билете"
          type={"number"}
        />
      </Grid>
      <Grid item xs={5}>
        {isLoadingFlightInTicket && <CircularProgress />}
        {!isLoadingFlightInTicket && error && (
          <Typography>
            Ошибка! Код: {error?.status}. Сообщение: {error?.data?.error}
          </Typography>
        )}
        {!isLoadingFlightInTicket &&
          !error &&
          flightInTicket &&
          searchQuery !== "" && (
            <EntityInfo
              items={columns
                .filter((col) => col.type !== "actions")
                .map((col) => {
                  return {
                    label: col.headerName,
                    value: flightInTicket[col.field],
                  };
                })}
              onDelete={() =>
                deleteFlightInTicket({ id: flightInTicket.id })
                  .unwrap()
                  .catch(({ data: { error } }) => alert(error))
              }
            />
          )}
      </Grid>
    </Grid>
  );
};
