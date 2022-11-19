import { createSlice } from "@reduxjs/toolkit";

const initialState = {
  items: [],
  total_count: 0,
};

export const airportSlice = createSlice({
  name: "airport",
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

export const { setItems, setTotalCount } = airportSlice.actions;
export default airportSlice.reducer;
