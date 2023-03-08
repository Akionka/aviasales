// Файл web\src\features\seat\seatsPage.jsx содержит код страницы с формой для таблицы "Места"
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
  useCreateSeatMutation,
  useDeleteSeatByIDMutation,
  useGetSeatByIDQuery,
  useGetSeatsQuery,
  useUpdateSeatByIDMutation,
} from "../../app/services/api";
import { useEffect, useState } from "react";
import {
  DataGrid,
  GridActionsCellItem,
  GridRowModes,
  GridToolbarContainer,
} from "@mui/x-data-grid";
import { EntityInfo } from "../../components/EntityInfo";

const EditToolbar = ({ setRows, setRowModesModel }) => {
  const handleClick = () => {
    setRows((oldRows) => [...oldRows, { id: 0, class: "Y", isNew: true }]);
    setRowModesModel((oldModel) => ({
      ...oldModel,
      0: { mode: GridRowModes.Edit, fieldToFocus: "id" },
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

export const SeatsPage = () => {
  const [paginationModel, setPaginationModel] = useState({
    pageSize: 25,
    page: 0,
  });

  const [rowModesModel, setRowModesModel] = useState({});
  const [rows, setRows] = useState([]);

  const [searchQuery, setSearchQuery] = useState("");

  const { data, isLoading } = useGetSeatsQuery({
    page: paginationModel.page + 1,
    count: paginationModel.pageSize,
  });

  const {
    data: seat,
    isLoading: isLoadingSeat,
    error,
  } = useGetSeatByIDQuery({ id: searchQuery });

  useEffect(() => {
    setRows((prevRows) => {
      const newRows = prevRows?.filter((r) => r.isNew) || [];
      if (data?.items) {
        return [...data?.items, ...newRows] || [];
      }
      return prevRows;
    });
  }, [data?.items]);

  const [updateSeat, { isLoadingUpdate }] = useUpdateSeatByIDMutation();
  const [deleteSeat, { isLoadingDelete }] = useDeleteSeatByIDMutation();
  const [createSeat, { isLoadingCreate }] = useCreateSeatMutation();

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
    deleteSeat({ id: row.id })
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
      headerName: "ID места",
      width: 150,
      editable: true,
      type: "number",
    },
    {
      field: "model_code",
      headerName: "Код модели самолёта",
      width: 175,
      editable: true,
    },
    {
      field: "number",
      headerName: "Номер места",
      width: 175,
      editable: true,
    },
    {
      field: "class",
      headerName: "Класс места",
      width: 175,
      editable: true,
      type: "singleSelect",
      valueOptions: [
        { value: "J", label: "Бизнес" },
        { value: "W", label: "Комфорт" },
        { value: "Y", label: "Эконом" },
      ],
    },
    {
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
    },
  ];

  if (isLoading)
    return <Skeleton variant="rectangular" width={512} height={512} />;

  return (
    <>
      {" "}
      <Grid rowSpacing={3} columnSpacing={3} container>
        <Grid item xs={12}>
          <DataGrid
            getRowHeight={() => "auto"}
            autoHeight
            editMode="row"
            getRowId={(row) => row.id}
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
            experimentalFeatures={{ newEditingApi: true }}
            processRowUpdate={async (newRow, oldRow) => {
              try {
                if (newRow.isNew) {
                  const res = await createSeat({
                    seat: newRow,
                  }).unwrap();
                  setRows((prevRows) =>
                    prevRows.filter((row) => row.id !== oldRow.id)
                  );
                  return res;
                } else {
                  const res = await updateSeat({
                    id: oldRow.id,
                    seat: newRow,
                  }).unwrap();
                  return res;
                }
              } catch (error) {
                throw new Error(error.data.error);
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
            label="ID места"
            type={"number"}
          />
        </Grid>
        <Grid item xs={5}>
          {isLoadingSeat && <CircularProgress />}
          {!isLoadingSeat && error && (
            <Typography>
              Ошибка! Код: {error?.status}. Сообщение: {error?.data?.error}
            </Typography>
          )}
          {!isLoadingSeat && !error && seat && searchQuery !== "" && (
            <EntityInfo
              items={columns
                .filter((col) => col.type !== "actions")
                .map((col) => {
                  return { label: col.headerName, value: seat[col.field] };
                })}
              onDelete={() =>
                deleteSeat({ id: seat.id })
                  .unwrap()
                  .catch(({ data: { error } }) => alert(error))
              }
            />
          )}
        </Grid>
      </Grid>
    </>
  );
};
