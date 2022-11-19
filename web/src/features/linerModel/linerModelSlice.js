import { createSlice } from "@reduxjs/toolkit";

const initialState = {
  items: [],
  total_count: 0,
};

export const linerModelSlice = createSlice({
  name: "liner_model",
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

export const { setItems, setTotalCount } = linerModelSlice.actions;
export default linerModelSlice.reducer;
