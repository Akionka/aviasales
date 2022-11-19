import { createSlice } from "@reduxjs/toolkit";

const initialState = {
  items: [],
  total_count: 0,
};

export const lineSlice = createSlice({
  name: "line",
  initialState,
  reducers: {
    setItems: (state, { payload: items }) => {
      state.items = items;
    },
    setTotalCount: (state, { payload: total_count }) => {
      state.total_count = total_count;
    },
  },
});

export const { setItems, setTotalCount } = lineSlice.actions;
export default lineSlice.reducer;
