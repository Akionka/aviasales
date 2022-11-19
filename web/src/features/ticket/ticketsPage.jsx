import SaveIcon from "@mui/icons-material/Save";
import AddIcon from "@mui/icons-material/Add";
import CancelIcon from "@mui/icons-material/Cancel";
import DeleteIcon from "@mui/icons-material/DeleteForeverOutlined";
import EditIcon from "@mui/icons-material/Edit";
import { Button, Skeleton } from "@mui/material";
import {
  useCreateTicketMutation,
  useDeleteTicketByIDMutation,
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
  const [page, setPage] = useState(0);
  const [rowCount, setRowCount] = useState(10);

  const [rowModesModel, setRowModesModel] = useState({});
  const [rows, setRows] = useState([]);

  const { data, isLoading } = useGetTicketsQuery({
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

  const [updateTicket, { isLoadingUpdate }] =    useUpdateTicketByIDMutation();
  const [deleteTicket, { isLoadingDelete }] =    useDeleteTicketByIDMutation();
  const [createTicket, { isLoadingCreate }] =    useCreateTicketMutation();

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
      valueFormatter: ({ value }) => value && new Date(value).toLocaleDateString(),
    },
    {
      field: "passenger_passport_number",
      headerName: "Номер паспорта пассажира",
      width: 210,
      editable: true,
      valueFormatter: ({ value }) => value && `${value.slice(0, 4)} ${value.slice(4, 10)}`,
    },
    {
      field: "passenger_sex",
      headerName: "Пол пассажира",
      width: 150,
      editable: true,
      type: "singleSelect",
      valueOptions: [
        {value: 1, label: "Мужской"},
        {value: 2, label: "Женский"},
      ],
      valueFormatter: ({ value }) => value && value === 1 ? 'Мужской' : 'Женский',
    },
    {
      field: "purchase_id",
      headerName: "ID покупки",
      width: 100,
      editable: true,
      type: "number"
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
                ticket: {...newRow, passenger_birth_date: dateInUTC},
              }).unwrap();
              setRows((prevRows) =>
                prevRows.filter((row) => row.id !== oldRow.id)
              );
              return res;
            } else {
              const res = await updateTicket({
                id: oldRow.id,
                ticket: {...newRow, passenger_birth_date: dateInUTC},
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
