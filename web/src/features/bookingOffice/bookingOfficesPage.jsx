import SaveIcon from "@mui/icons-material/Save";
import AddIcon from "@mui/icons-material/Add";
import CancelIcon from "@mui/icons-material/Cancel";
import DeleteIcon from "@mui/icons-material/DeleteForeverOutlined";
import EditIcon from "@mui/icons-material/Edit";
import { Button, Skeleton } from "@mui/material";
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
  useGetOfficesQuery,
  useUpdateOfficeByIDMutation,
} from "../../app/services/api";

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
  const [page, setPage] = useState(0);
  const [rowCount, setRowCount] = useState(10);

  const [rowModesModel, setRowModesModel] = useState({});
  const [rows, setRows] = useState([]);

  const { data, isLoading } = useGetOfficesQuery({
    page: page + 1,
    count: rowCount,
  });

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
      width: 500,
      editable: true,
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
      <DataGrid
        autoHeight
        editMode="row"
        columns={columns}
        rows={rows}
        rowCount={data.total_count}
        rowsPerPageOptions={[5, 10, 15, 20, 25, 50, 100]}
        pageSize={rowCount}
        onPageSizeChange={(newRowCount) => setRowCount(newRowCount)}
        page={page}
        onPageChange={(newPage) => {
          setPage(newPage);
        }}
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
    </>
  );
};
