// Файл web\src\features\cashier\cashiersPage.jsx содержит код страницы с формой для таблицы "Кассир"
import SaveIcon from "@mui/icons-material/Save";
import AddIcon from "@mui/icons-material/Add";
import CancelIcon from "@mui/icons-material/Cancel";
import DeleteIcon from "@mui/icons-material/DeleteForeverOutlined";
import EditIcon from "@mui/icons-material/Edit";
import PasswordIcon from "@mui/icons-material/Password";
import {
  Button,
  CircularProgress,
  Grid,
  Skeleton,
  Stack,
  TextField,
  Typography,
} from "@mui/material";
import { useEffect, useState } from "react";
import {
  DataGrid,
  GridActionsCellItem,
  GridRowModes,
  GridToolbarContainer,
  useGridApiContext,
} from "@mui/x-data-grid";
import {
  useCreateCashierMutation,
  useDeleteCashierByIDMutation,
  useGetCashierByIDQuery,
  useGetCashiersQuery,
  useUpdateCashierByIDMutation,
  useUpdateCashierPasswordMutation,
} from "../../app/services/api";
import { EntityInfo } from "../../components/EntityInfo";

const EditToolbar = ({ setRows, setRowModesModel }) => {
  const handleClick = () => {
    setRows((oldRows) => [
      ...oldRows,
      {
        id: 0,
        login: "",
        last_name: "",
        first_name: "",
        middle_name: "",
        isNew: true,
      },
    ]);
    setRowModesModel((oldModel) => ({
      ...oldModel,
      0: { mode: GridRowModes.Edit, fieldToFocus: "login" },
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

const PasswordCell = (props) => {
  const { id, value, field } = props;
  const apiRef = useGridApiContext();

  const handleValueChange = (event) => {
    const newValue = event.target.value;
    apiRef.current.setEditCellValue({ id, field, value: newValue });
  };

  return (
    <>
      <Stack direction="row" spacing={2}>
        <TextField value={value} onChange={handleValueChange} />
        <Button
          endIcon={<PasswordIcon />}
          title="Изменить пароль"
          variant="contained"
          size="small"
          onClick={() => props.onClick(id, value)}
        >
          Изменить пароль
        </Button>
      </Stack>
    </>
  );
};

export const CashiersPage = () => {
  const [paginationModel, setPaginationModel] = useState({
    pageSize: 25,
    page: 0,
  });

  const [rowModesModel, setRowModesModel] = useState({});
  const [rows, setRows] = useState([]);

  const [searchQuery, setSearchQuery] = useState("");

  const { data, isLoading } = useGetCashiersQuery({
    page: paginationModel.page + 1,
    count: paginationModel.pageSize,
  });

  const {
    data: cashier,
    isLoading: isLoadingCashier,
    error,
  } = useGetCashierByIDQuery({ id: searchQuery });

  useEffect(() => {
    setRows((prevRows) => {
      const newRows = prevRows?.filter((r) => r.isNew) || [];
      if (data?.items) {
        return (
          [
            ...data?.items.map((item) => {
              return { ...item, password: item.password || "" };
            }),
            ...newRows,
          ] || []
        );
      }
      return prevRows;
    });
  }, [data?.items]);

  const [updateCashier, { isLoadingUpdate }] = useUpdateCashierByIDMutation();
  const [updateCashierPassword] = useUpdateCashierPasswordMutation();
  const [deleteCashier, { isLoadingDelete }] = useDeleteCashierByIDMutation();
  const [createCashier, { isLoadingCreate }] = useCreateCashierMutation();

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
    deleteCashier({ id: row.id })
      .unwrap()
      .catch((error) => alert(error.data.error));
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

  const handlePasswordChangeClick = (id, new_password) => {
    updateCashierPassword({ id, new_password })
      .unwrap()
      .then(() => alert("Успешно изменен пароль"))
      .catch((error) => alert(error.data.error));
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
    { field: "login", headerName: "Логин", width: 200, editable: true },
    {
      field: "last_name",
      headerName: "Фамилия",
      width: 200,
      editable: true,
    },
    {
      field: "first_name",
      headerName: "Имя",
      width: 200,
      editable: true,
    },
    {
      field: "middle_name",
      headerName: "Отчество",
      width: 200,
      editable: true,
    },
    {
      field: "password",
      headerName: "Пароль",
      width: 500,
      editable: true,
      renderEditCell: (props) => (
        <PasswordCell {...props} onClick={handlePasswordChangeClick} />
      ),
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
                  const res = await createCashier({
                    cashier: { ...newRow, id: Number(newRow.id) },
                  }).unwrap();
                  setRows((prevRows) =>
                    prevRows.filter((row) => row.id !== oldRow.id)
                  );
                  return res;
                } else {
                  const res = await updateCashier({
                    id: oldRow.id,
                    cashier: newRow,
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
            label="ID кассира"
            type={"number"}
          />
        </Grid>
        <Grid item xs={5}>
          {isLoadingCashier && <CircularProgress />}
          {!isLoadingCashier && error && (
            <Typography>
              Ошибка! Код: {error?.status}. Сообщение: {error?.data?.error}
            </Typography>
          )}
          {!isLoadingCashier && !error && cashier && searchQuery !== "" && (
            <EntityInfo
              items={columns
                .filter(
                  (col) => col.type !== "actions" && col.field !== "password"
                )
                .map((col) => {
                  return { label: col.headerName, value: cashier[col.field] };
                })}
              onDelete={() =>
                deleteCashier({ id: cashier.id })
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
