// Файл web\src\features\line\linerModelsPage.jsx содержит код страницы с формой для таблицы "Модели самолётов"
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
  useCreateLinerModelMutation,
  useDeleteLinerModelByCodeMutation,
  useGetLinerModelByCodeQuery,
  useGetLinerModelsQuery,
  useUpdateLinerModelByCodeMutation,
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
      { iata_type_code: "", name: "", isNew: true },
    ]);
    setRowModesModel((oldModel) => ({
      ...oldModel,
      "": { mode: GridRowModes.Edit, fieldToFocus: "iata_type_code" },
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

export const LinerModelsPage = () => {
  const auth = useAuth()
  const [paginationModel, setPaginationModel] = useState({
    pageSize: 25,
    page: 0,
  });

  const [rowModesModel, setRowModesModel] = useState({});
  const [rows, setRows] = useState([]);

  const [searchQuery, setSearchQuery] = useState("");

  const { data, isLoading } = useGetLinerModelsQuery({
    page: paginationModel.page + 1,
    count: paginationModel.pageSize,
  });

  const {
    data: linerModel,
    isLoading: isLoadingLinerModel,
    error,
  } = useGetLinerModelByCodeQuery({ code: searchQuery });

  useEffect(() => {
    setRows((prevRows) => {
      const newRows = prevRows?.filter((r) => r.isNew) || [];
      if (data?.items) {
        return [...data?.items, ...newRows] || [];
      }
      return prevRows;
    });
  }, [data?.items]);

  const [updateLinerModel, { isLoadingUpdate }] =
    useUpdateLinerModelByCodeMutation();
  const [deleteLinerModel, { isLoadingDelete }] =
    useDeleteLinerModelByCodeMutation();
  const [createLinerModel, { isLoadingCreate }] = useCreateLinerModelMutation();

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
    deleteLinerModel({ code: row.id })
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
      field: "iata_type_code",
      headerName: "Код модели самолёта",
      width: 172,
      editable: true,
    },
    {
      field: "name",
      headerName: "Название",
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
    return <Skeleton variant="rectangular" width={'auto'} height={256} />;

  return (
      <Grid rowSpacing={3} columnSpacing={3} container>
        <Grid item xs={12}>
          <DataGrid
            autoHeight
            editMode="row"
            getRowId={(row) => row.iata_type_code}
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
                  const res = await createLinerModel({
                    liner_model: newRow,
                  }).unwrap();
                  setRows((prevRows) =>
                    prevRows.filter(
                      (row) => row.iata_type_code !== oldRow.iata_type_code
                    )
                  );
                  return res;
                } else {
                  const res = await updateLinerModel({
                    code: oldRow.iata_type_code,
                    liner_model: newRow,
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
            label="Код модели самолёта"
          />
        </Grid>
        <Grid item xs={5}>
          {isLoadingLinerModel && <CircularProgress />}
          {!isLoadingLinerModel && error && (
            <Typography>
              Ошибка! Код: {error?.status}. Сообщение: {error?.data?.error}
            </Typography>
          )}
          {!isLoadingLinerModel && !error && linerModel && searchQuery !== "" && (
            <EntityInfo
              items={columns
                .filter((col) => col.type !== "actions")
                .map((col) => {
                  return {
                    label: col.headerName,
                    value: linerModel[col.field],
                  };
                })}
              onDelete={() =>
                deleteLinerModel({ code: linerModel.iata_type_code })
                  .unwrap()
                  .catch(({ data: { error } }) => alert(error))
              }
            />
          )}
        </Grid>
      </Grid>
  );
};
