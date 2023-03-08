// Файл web\src\features\bookingOffice\bookingOfficesPage.jsx содержит код страницы с формой для таблицы "Касса"
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
import { useEffect, useState } from "react";
import {
  DataGrid,
  GridActionsCellItem,
  GridRowModes,
  GridToolbarContainer,
} from "@mui/x-data-grid";
import {
  useCreateOfficeMutation,
  useDeleteOfficeByIDMutation,
  useGetOfficeByIDQuery,
  useGetOfficesQuery,
  useUpdateOfficeByIDMutation,
} from "../../app/services/api";
import { EntityInfo } from "../../components/EntityInfo";
import { formatPhoneNumberIntl } from "react-phone-number-input";
import { useAuth } from "../../hooks/useAuth";

const EditToolbar = ({ setRows, setRowModesModel }) => {
  const handleClick = () => {
    setRows((oldRows) => [
      ...oldRows,
      { id: 0, address: "", phone_number: "", isNew: true },
    ]);
    setRowModesModel((oldModel) => ({
      ...oldModel,
      0: { mode: GridRowModes.Edit, fieldToFocus: "address" },
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

export const BookingOfficesPage = () => {
  const auth = useAuth()
  const [paginationModel, setPaginationModel] = useState({
    pageSize: 25,
    page: 0,
  });

  const [rowModesModel, setRowModesModel] = useState({});
  const [rows, setRows] = useState([]);

  const [searchQuery, setSearchQuery] = useState("");

  const { data, isLoading } = useGetOfficesQuery({
    page: paginationModel.page + 1,
    count: paginationModel.pageSize,
  });

  const {
    data: office,
    isLoading: isLoadingOffice,
    error,
  } = useGetOfficeByIDQuery({ id: searchQuery });

  useEffect(() => {
    setRows((prevRows) => {
      const newRows = prevRows?.filter((r) => r.isNew) || [];
      if (data?.items) {
        return [...data?.items, ...newRows] || [];
      }
      return prevRows;
    });
  }, [data?.items]);

  const [updateOffice, { isLoadingUpdate }] = useUpdateOfficeByIDMutation();
  const [deleteOffice, { isLoadingDelete }] = useDeleteOfficeByIDMutation();
  const [createOffice, { isLoadingCreate }] = useCreateOfficeMutation();

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
    deleteOffice({ id: row.id })
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
    { field: "address", headerName: "Адрес", width: 500, editable: true },
    {
      field: "phone_number",
      headerName: "Номер телефона",
      width: 150,
      editable: true,
      valueFormatter: ({value}) => value && formatPhoneNumberIntl('+'+value)
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
                  const res = await createOffice({
                    office: { ...newRow, id: Number(newRow.id) },
                  }).unwrap();
                  setRows((prevRows) =>
                    prevRows.filter((row) => row.id !== oldRow.id)
                  );
                  return res;
                } else {
                  const res = await updateOffice({
                    id: oldRow.id,
                    office: newRow,
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
            label="ID кассы"
            type={"number"}
          />
        </Grid>
        <Grid item xs={5}>
          {isLoadingOffice && <CircularProgress />}
          {!isLoadingOffice && error && (
            <Typography>
              Ошибка! Код: {error?.status}. Сообщение: {error?.data?.error}
            </Typography>
          )}
          {!isLoadingOffice && !error && office && searchQuery !== "" && (
            <EntityInfo
              items={columns
                .filter((col) => col.type !== "actions")
                .map((col) => {
                  return { label: col.headerName, value: office[col.field] };
                })}
              onDelete={() =>
                deleteOffice({ id: office.id })
                  .unwrap()
                  .catch(({ data: { error } }) => alert(error))
              }
            />
          )}
        </Grid>
      </Grid>
  );
};
