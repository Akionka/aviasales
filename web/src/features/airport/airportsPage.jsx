// Файл web\src\features\airport\airportsPage.jsx содержит код страницы с формой для таблицы "Аэропорт"
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
  useCreateAirportMutation,
  useDeleteAirportByCodeMutation,
  useGetAirportByCodeQuery,
  useGetAirportsQuery,
  useGetTimezonesQuery,
  useUpdateAirportByCodeMutation,
} from "../../app/services/api";
import { useEffect, useMemo, useState } from "react";
import {
  DataGrid,
  GridActionsCellItem,
  GridRowModes,
  GridToolbarContainer,
} from "@mui/x-data-grid";
import { EntityInfo } from "../../components/EntityInfo";
import {useAuth} from "../../hooks/useAuth"

const EditToolbar = ({ setRows, setRowModesModel }) => {
  const handleClick = () => {
    setRows((oldRows) => [
      ...oldRows,
      { iata_code: "New", city: "", timezone: "", isNew: true },
    ]);
    setRowModesModel((oldModel) => ({
      ...oldModel,
      New: { mode: GridRowModes.Edit, fieldToFocus: "iata_code" },
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

export const AirportsPage = () => {
  const auth = useAuth()
  const [paginationModel, setPaginationModel] = useState({
    pageSize: 25,
    page: 0,
  });

  const [rowModesModel, setRowModesModel] = useState({});
  const [rows, setRows] = useState([]);

  const [searchQuery, setSearchQuery] = useState("");

  const { data, isLoading } = useGetAirportsQuery({
    page: paginationModel.page + 1,
    count: paginationModel.pageSize,
  });

  const {
    data: airport,
    isLoading: isLoadingAirport,
    error,
  } = useGetAirportByCodeQuery({ code: searchQuery });

  const { data: timezones, isLoadingTimezones } = useGetTimezonesQuery();

  useEffect(() => {
    setRows((prevRows) => {
      const newRows = prevRows?.filter((r) => r.isNew) || [];
      if (data?.items) {
        return [...data?.items, ...newRows] || [];
      }
      return prevRows;
    });
  }, [data?.items]);

  const [updateAirport, { isLoadingUpdate }] = useUpdateAirportByCodeMutation();
  const [deleteAirport, { isLoadingDelete }] = useDeleteAirportByCodeMutation();
  const [createAirport, { isLoadingCreate }] = useCreateAirportMutation();

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
    deleteAirport({ code: row.id })
      .unwrap()
      .catch(({ data: { error } }) => alert(error));
  };

  const handleCancelClick = (row) => () => {
    setRowModesModel({
      ...rowModesModel,
      [row.id]: { mode: GridRowModes.View, ignoreModifications: true },
    });

    const editedRow = rows.find((r) => r.iata_code === row.id);
    if (editedRow.isNew) {
      setRows(rows.filter((r) => r.iata_code !== row.id));
    }
  };

  const handleSearchQueryChange = (e) => {
    setSearchQuery(e.target.value);
  };
  const columns = [
    { field: "iata_code", headerName: "Код IATA", width: 100, editable: true },
    { field: "city", headerName: "Город", width: 200, editable: true },
    {
      field: "timezone",
      headerName: "Часовой пояс",
      type: "singleSelect",
      width: 200,
      editable: true,
      valueOptions: useMemo(() => timezones || [], [timezones]),
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

  if (isLoading || isLoadingTimezones)
    return <Skeleton variant="rectangular" width={512} height={512} />;

  return (
    <Grid rowSpacing={3} columnSpacing={3} container>
      <Grid item xs={12}>
        <DataGrid
          autoHeight
          editMode="row"
          getRowId={(row) => row.iata_code}
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
                const res = await createAirport({ airport: newRow }).unwrap();
                setRows((prevRows) =>
                  prevRows.filter((row) => row.iata_code !== oldRow.iata_code)
                );
                return res;
              } else {
                const res = await updateAirport({
                  code: oldRow.iata_code,
                  airport: newRow,
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
          label="IATA код аэропорта"
        />
      </Grid>
      <Grid item xs={5}>
        {isLoadingAirport && <CircularProgress />}
        {!isLoadingAirport && error && (
          <Typography>
            Ошибка! Код: {error?.status}. Сообщение: {error?.data?.error}
          </Typography>
        )}
        {!isLoadingAirport && !error && airport && searchQuery !== "" && (
          <EntityInfo
            items={columns
              .filter((col) => col.type !== "actions")
              .map((col) => {
                return { label: col.headerName, value: airport[col.field] };
              })}
            onDelete={() =>
              deleteAirport({ code: airport.iata_code })
                .unwrap()
                .catch(({ data: { error } }) => alert(error))
            }
          />
        )}
      </Grid>
    </Grid>
  );
};
