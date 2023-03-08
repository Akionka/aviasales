// Файл web\src\features\purchase\purchasesPage.jsx содержит код страницы с формой для таблицы "Покупки"
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
  useCreatePurchaseMutation,
  useDeletePurchaseByIDMutation,
  useGetPurchaseByIDQuery,
  useGetPurchasesQuery,
  useUpdatePurchaseByIDMutation,
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
    setRows((oldRows) => [...oldRows, { id: "", model_code: "", isNew: true }]);
    setRowModesModel((oldModel) => ({
      ...oldModel,
      "": { mode: GridRowModes.Edit, fieldToFocus: "id" },
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

export const PurchasesPage = () => {
  const auth = useAuth()
  const [paginationModel, setPaginationModel] = useState({
    pageSize: 25,
    page: 0,
  });

  const [rowModesModel, setRowModesModel] = useState({});
  const [rows, setRows] = useState([]);

  const [searchQuery, setSearchQuery] = useState("");

  const { data, isLoading } = useGetPurchasesQuery({
    page: paginationModel.page + 1,
    count: paginationModel.pageSize,
  });

  const {
    data: purchase,
    isLoading: isLoadingPurchase,
    error,
  } = useGetPurchaseByIDQuery({ id: searchQuery });

  useEffect(() => {
    setRows((prevRows) => {
      const newRows = prevRows?.filter((r) => r.isNew) || [];
      if (data?.items) {
        return [...data?.items, ...newRows] || [];
      }
      return prevRows;
    });
  }, [data?.items]);

  const [updatePurchase, { isLoadingUpdate }] = useUpdatePurchaseByIDMutation();
  const [deletePurchase, { isLoadingDelete }] = useDeletePurchaseByIDMutation();
  const [createPurchase, { isLoadingCreate }] = useCreatePurchaseMutation();

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
    deletePurchase({ id: row.id })
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
      headerName: "ID покупки",
      width: 150,
      editable: false,
      type: "number",
    },
    {
      field: "date",
      headerName: "Время покупки",
      width: 200,
      editable: true,
      type: "dateTime",
      valueFormatter: ({ value }) => value && new Date(value).toLocaleString(),
    },
    {
      field: "contact_phone",
      headerName: "Контактный телефон",
      width: 175,
      editable: true,
    },
    {
      field: "contact_email",
      headerName: "Контактный e-mail",
      width: 175,
      editable: true,
    },
    {
      field: "total_price",
      headerName: "Итоговая цена",
      width: 175,
      editable: true,
      type: "number",
    },
    {
      field: "booking_office_id",
      headerName: "ID кассы",
      width: 100,
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
            try {
              if (newRow.isNew) {
                const res = await createPurchase({
                  purchase: newRow,
                }).unwrap();
                setRows((prevRows) =>
                  prevRows.filter((row) => row.id !== oldRow.id)
                );
                return res;
              } else {
                const res = await updatePurchase({
                  id: oldRow.id,
                  purchase: newRow,
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
          label="ID покупки"
          type={"number"}
        />
      </Grid>
      <Grid item xs={5}>
        {isLoadingPurchase && <CircularProgress />}
        {!isLoadingPurchase && error && (
          <Typography>
            Ошибка! Код: {error?.status}. Сообщение: {error?.data?.error}
          </Typography>
        )}
        {!isLoadingPurchase && !error && purchase && searchQuery !== "" && (
          <EntityInfo
            items={columns
              .filter((col) => col.type !== "actions")
              .map((col) => {
                return {
                  label: col.headerName,
                  value: (() => {
                    switch (col.type) {
                      case "boolean":
                        return purchase[col.field] ? "Да" : "Нет";
                      case "date":
                        return purchase[col.field]
                          ? new Date(purchase[col.field]).toLocaleDateString()
                          : "";
                      case "dateTime":
                        return purchase[col.field]
                          ? new Date(purchase[col.field]).toLocaleString()
                          : "";
                      default:
                        return purchase[col.field];
                    }
                  })(),
                };
              })}
            onDelete={() =>
              deletePurchase({ id: purchase.id })
                .unwrap()
                .catch(({ data: { error } }) => alert(error))
            }
          />
        )}
      </Grid>
    </Grid>
  );
};
