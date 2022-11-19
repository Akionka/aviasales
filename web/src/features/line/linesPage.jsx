import SaveIcon from "@mui/icons-material/Save";
import AddIcon from "@mui/icons-material/Add";
import CancelIcon from "@mui/icons-material/Cancel";
import DeleteIcon from "@mui/icons-material/DeleteForeverOutlined";
import EditIcon from "@mui/icons-material/Edit";
import { Button, Skeleton, TextField } from "@mui/material";
import {
  useCreateLineMutation,
  useDeleteLineByCodeMutation,
  useGetLinesQuery,
  useUpdateLineByCodeMutation,
} from "../../app/services/api";
import { useEffect, useState } from "react";
import {
  DataGrid,
  GridActionsCellItem,
  GridRowModes,
  GridToolbarContainer,
  useGridApiContext,
} from "@mui/x-data-grid";
import { TimePicker } from "@mui/x-date-pickers/TimePicker";

const EditToolbar = ({ setRows, setRowModesModel }) => {
  const handleClick = () => {
    setRows((oldRows) => [...oldRows, { line_code: 'JR', dep_time: '10:10', arr_time: '12:12', isNew: true }]);
    setRowModesModel((oldModel) => ({
      ...oldModel,
      'JR': { mode: GridRowModes.Edit, fieldToFocus: "dep_date" },
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

const TimeEditCell = (props) => {
  const { id, value, field } = props;
  const apiRef = useGridApiContext();

  const handleValueChange = (newValue) => {
    apiRef.current.setEditCellValue({ id, field, value: newValue });
  };

  return (
    <TimePicker
      label=""
      ampm={false}
      value={value}
      onChange={handleValueChange}
      renderInput={(params) => <TextField {...params} />}
    />
  );
};

const timeType = {
  valueFormatter: ({ value }) =>
    value &&
    `${String(value.getHours()).padStart(2, "0")}:${String(
      value.getMinutes()
    ).padStart(2, "0")}`,
  valueGetter: ({ value }) => {
    if (value) {
      const [hours, minutes, seconds] = value.split(":");
      const date = new Date();
      date.setHours(hours);
      date.setMinutes(minutes);
      date.setSeconds(seconds);
      return date;
    }
  },
  renderEditCell: (props) => <TimeEditCell {...props} />,
};

export const LinesPage = () => {
  const [page, setPage] = useState(0);
  const [rowCount, setRowCount] = useState(10);

  const [rowModesModel, setRowModesModel] = useState({});
  const [rows, setRows] = useState([]);

  const { data, isLoading } = useGetLinesQuery({
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

  const [updateLine, { isLoadingUpdate }] = useUpdateLineByCodeMutation();
  const [deleteLine, { isLoadingDelete }] = useDeleteLineByCodeMutation();
  const [createLine, { isLoadingCreate }] = useCreateLineMutation();

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
    deleteLine({ code: row.id })
      .unwrap()
      .catch(({ data: { error } }) => alert(error));
  };

  const handleCancelClick = (row) => () => {
    setRowModesModel({
      ...rowModesModel,
      [row.id]: { mode: GridRowModes.View, ignoreModifications: true },
    });
    const editedRow = rows.find((r) => r.line_code === row.id);
    if (editedRow.isNew) {
      setRows(rows.filter((r) => r.line_code !== row.id));
    }
  };

  const columns = [
    {
      field: "line_code",
      headerName: "Код рейса",
      width: 100,
      editable: true,
    },
    {
      field: "dep_time",
      headerName: "Время вылета",
      width: 200,
      editable: true,
      ...timeType,
    },
    {
      field: "arr_time",
      headerName: "Время прибытия",
      width: 200,
      editable: true,
      ...timeType,
    },
    {
      field: "base_price",
      headerName: "Базовая цена",
      width: 100,
      editable: true,
      type: "number",
    },
    {
      field: "dep_airport",
      headerName: "Код эропорта вылета",
      width: 200,
      editable: true,
    },
    {
      field: "arr_airport",
      headerName: "Код аэропорта прибытия",
      width: 200,
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
        getRowId={(row) => row.line_code}
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
              const res = await createLine({
                line: { ...newRow, dep_time: newRow.dep_time.format("HH:mm"), arr_time: newRow.arr_time.format("HH:mm") },
              }).unwrap();
              setRows((prevRows) =>
                prevRows.filter((row) => row.line_code !== oldRow.line_code)
              );
              return res;
            } else {
              const res = await updateLine({
                code: oldRow.line_code,
                line: { ...newRow, dep_time: newRow.dep_time.format("HH:mm"), arr_time: newRow.arr_time.format("HH:mm") },
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
