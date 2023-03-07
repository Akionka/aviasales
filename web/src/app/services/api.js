// Файл web\src\app\services\api.js содержит код для взаимодействия с сервером
import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";

export const api = createApi({
  baseQuery: fetchBaseQuery({
    baseUrl: "http://127.0.0.1:8080/api",
    tagTypes: [
      "Airport",
      "BookingOffice",
      "Cashier",
      "Flight",
      "FlightInTicket",
      "Line",
      "Liner",
      "LinerModel",
      "Purchase",
      "Seat",
      "Ticket",
      "Session",
    ],
    prepareHeaders: (headers, { getState }) => {
      const token = getState().auth.token;
      if (token) {
        headers.set("Authorization", `Bearer ${token}`);
      }
      return headers;
    },
  }),
  endpoints: (builder) => ({
    login: builder.mutation({
      query: (credentials) => ({
        url: "session",
        method: "post",
        body: credentials,
      }),
    }),
    signup: builder.mutation({
      query: (credentials) => ({
        url: "user",
        method: "post",
        body: credentials,
      }),
      invalidatesTags: (result, error, {id}) => [
        {type: "Cashier", id: id},
        {type: "Cashier", id: "PARTIAL-LIST"},
      ]
    }),
    whoami: builder.query({
      query: () => "whoami",
      providesTags: ["Session"],
    }),

    getTimezones: builder.query({
      query: () => "timezones",
    }),

    getAirports: builder.query({
      query: ({ page, count }) => `airports?page=${page}&count=${count}`,
      providesTags: (result, error, params) =>
        result
          ? [
              ...result.items.map(({ iata_code }) => ({
                type: "Airport",
                id: iata_code,
              })),
              { type: "Airport", id: "PARTIAL-LIST" },
            ]
          : [{ type: "Airport", id: "PARTIAL-LIST" }],
    }),
    createAirport: builder.mutation({
      query: ({ airport }) => ({
        url: `airports`,
        method: "post",
        body: airport,
      }),
      invalidatesTags: (result, error, { airport: { iata_code } }) => [
        { type: "Airport", id: iata_code },
        { type: "Airport", id: "PARTIAL-LIST" },
      ],
    }),

    getOffices: builder.query({
      query: ({ page, count }) => `booking_offices?page=${page}&count=${count}`,
      providesTags: (result, error, params) =>
        result
          ? [
              ...result.items.map(({ id }) => ({ type: "BookingOffice", id })),
              { type: "BookingOffice", id: "PARTIAL-LIST" },
            ]
          : [{ type: "BookingOffice", id: "PARTIAL-LIST" }],
    }),
    createOffice: builder.mutation({
      query: ({ office }) => ({
        url: `booking_offices`,
        method: "post",
        body: office,
      }),
      invalidatesTags: (result, error, { office: { id } }) => [
        { type: "BookingOffice", id: id },
        { type: "BookingOffice", id: "PARTIAL-LIST" },
      ],
    }),

    getCashiers: builder.query({
      query: ({ page, count }) => `cashiers?page=${page}&count=${count}`,
      providesTags: (result, error, params) =>
        result
          ? [
              ...result.items.map(({ id }) => ({ type: "Cashier", id: id })),
              { type: "Cashier", id: "PARTIAL-LIST" },
            ]
          : [{ type: "Cashier", id: "PARTIAL-LIST" }],
    }),
    createCashier: builder.mutation({
      query: ({ cashier }) => ({
        url: `cashiers`,
        method: "post",
        body: cashier,
      }),
      invalidatesTags: (result, error, { cashier: { id } }) => [
        { type: "Cashier", id: id },
        { type: "Cashier", id: "PARTIAL-LIST" },
      ],
    }),

    getFlightInTickets: builder.query({
      query: ({ page, count }) =>
        `flight_in_tickets?page=${page}&count=${count}`,
      providesTags: (result, error, params) =>
        result
          ? [
              ...result.items.map(({ id }) => ({ type: "FlightInTicket", id })),
              { type: "FlightInTicket", id: "PARTIAL-LIST" },
            ]
          : [{ type: "FlightInTicket", id: "PARTIAL-LIST" }],
    }),
    createFlightInTicket: builder.mutation({
      query: ({ flight_in_ticket }) => ({
        url: `flight_in_tickets`,
        method: "post",
        body: flight_in_ticket,
      }),
      invalidatesTags: (result, error, { flight_in_ticket: { id } }) => [
        { type: "FlightInTicket", id: id },
        { type: "FlightInTicket", id: "PARTIAL-LIST" },
      ],
    }),

    getFlights: builder.query({
      query: ({ page, count }) => `flights?page=${page}&count=${count}`,
      providesTags: (result, error, params) =>
        result
          ? [
              ...result.items.map(({ id }) => ({ type: "Flight", id: id })),
              { type: "Flight", id: "PARTIAL-LIST" },
            ]
          : [{ type: "Flight", id: "PARTIAL-LIST" }],
    }),
    createFlight: builder.mutation({
      query: ({ flight }) => ({
        url: `flights`,
        method: "post",
        body: flight,
      }),
      invalidatesTags: (result, error, { flight: { id } }) => [
        { type: "Flight", id: id },
        { type: "Flight", id: "PARTIAL-LIST" },
      ],
    }),

    getLines: builder.query({
      query: ({ page, count }) => `lines?page=${page}&count=${count}`,
      providesTags: (result, error, params) =>
        result
          ? [
              ...result.items.map(({ line_code }) => ({
                type: "Line",
                id: line_code,
              })),
              { type: "Line", id: "PARTIAL-LIST" },
            ]
          : [{ type: "Line", id: "PARTIAL-LIST" }],
    }),
    createLine: builder.mutation({
      query: ({ line }) => ({
        url: `lines`,
        method: "post",
        body: line,
      }),
      invalidatesTags: (result, error, { line: { line_code } }) => [
        { type: "Line", id: line_code },
        { type: "Line", id: "PARTIAL-LIST" },
      ],
    }),

    getLiners: builder.query({
      query: ({ page, count }) => `liners?page=${page}&count=${count}`,
      providesTags: (result, error, params) =>
        result
          ? [
              ...result.items.map(({ iata_code }) => ({
                type: "Liner",
                id: iata_code,
              })),
              { type: "Liner", id: "PARTIAL-LIST" },
            ]
          : [{ type: "Liner", id: "PARTIAL-LIST" }],
    }),
    createLiner: builder.mutation({
      query: ({ liner }) => ({
        url: `liners`,
        method: "post",
        body: liner,
      }),
      invalidatesTags: (result, error, { liner: { iata_code } }) => [
        { type: "Liner", id: iata_code },
        { type: "Liner", id: "PARTIAL-LIST" },
      ],
    }),

    getLinerModels: builder.query({
      query: ({ page, count }) => `liner_models?page=${page}&count=${count}`,
      providesTags: (result, error, params) =>
        result
          ? [
              ...result.items.map(({ iata_type_code }) => ({
                type: "LinerModel",
                id: iata_type_code,
              })),
              { type: "LinerModel", id: "PARTIAL-LIST" },
            ]
          : [{ type: "LinerModel", id: "PARTIAL-LIST" }],
    }),
    createLinerModel: builder.mutation({
      query: ({ linermodel }) => ({
        url: `liner_models`,
        method: "post",
        body: linermodel,
      }),
      invalidatesTags: (result, error, { linermodel: { iata_type_code } }) => [
        { type: "LinerModel", id: iata_type_code },
        { type: "LinerModel", id: "PARTIAL-LIST" },
      ],
    }),

    getPurchases: builder.query({
      query: ({ page, count }) => `purchases?page=${page}&count=${count}`,
      providesTags: (result, error, params) =>
        result
          ? [
              ...result.items.map(({ id }) => ({ type: "Purchase", id: id })),
              { type: "Purchase", id: "PARTIAL-LIST" },
            ]
          : [{ type: "Purchase", id: "PARTIAL-LIST" }],
    }),
    createPurchase: builder.mutation({
      query: ({ purchase }) => ({
        url: `purchases`,
        method: "post",
        body: purchase,
      }),
      invalidatesTags: (result, error, { purchase: { id } }) => [
        { type: "Purchase", id: id },
        { type: "Purchase", id: "PARTIAL-LIST" },
      ],
    }),

    getSeats: builder.query({
      query: ({ page, count }) => `seats?page=${page}&count=${count}`,
      providesTags: (result, error, params) =>
        result
          ? [
              ...result.items.map(({ id }) => ({ type: "Seat", id: id })),
              { type: "Seat", id: "PARTIAL-LIST" },
            ]
          : [{ type: "Seat", id: "PARTIAL-LIST" }],
    }),
    createSeat: builder.mutation({
      query: ({ seat }) => ({
        url: `seats`,
        method: "post",
        body: seat,
      }),
      invalidatesTags: (result, error, { seat: { id } }) => [
        { type: "Seat", id: id },
        { type: "Seat", id: "PARTIAL-LIST" },
      ],
    }),

    getTickets: builder.query({
      query: ({ page, count }) => `tickets?page=${page}&count=${count}`,
      providesTags: (result, error, params) =>
        result
          ? [
              ...result.items.map(({ id }) => ({ type: "Ticket", id: id })),
              { type: "Ticket", id: "PARTIAL-LIST" },
            ]
          : [{ type: "Ticket", id: "PARTIAL-LIST" }],
    }),
    createTicket: builder.mutation({
      query: ({ ticket }) => ({
        url: `tickets`,
        method: "post",
        body: ticket,
      }),
      invalidatesTags: (result, error, { ticket: { id } }) => [
        { type: "Ticket", id: id },
        { type: "Ticket", id: "PARTIAL-LIST" },
      ],
    }),

    getAirportByCode: builder.query({
      query: ({ code }) => `airports/${code}`,
      providesTags: (result, error, params) =>
        result ? [{ type: "Airport", id: result.iata_code }] : [],
    }),
    deleteAirportByCode: builder.mutation({
      query: ({ code }) => ({
        url: `airports/${code}`,
        method: "delete",
      }),
      invalidatesTags: (result, error, { code }) => [
        { type: "Airport", id: code },
        { type: "Airport", id: "PARTIAL-LIST" },
      ],
    }),
    updateAirportByCode: builder.mutation({
      query: ({ code, airport }) => ({
        url: `airports/${code}`,
        method: "put",
        body: airport,
      }),
      invalidatesTags: (result, error, { code }) => [
        { type: "Airport", id: code },
        { type: "Airport", id: "PARTIAL-LIST" },
        { type: "Line", id: "PARTIAL-LIST" },
      ],
    }),

    getOfficeByID: builder.query({
      query: ({ id }) => `booking_offices/${id}`,
      providesTags: (result, error, params) =>
        result ? [{ type: "BookingOffice", id: result.id }] : [],
    }),
    deleteOfficeByID: builder.mutation({
      query: ({ id }) => ({
        url: `booking_offices/${id}`,
        method: "delete",
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: "BookingOffice", id: id },
        { type: "BookingOffice", id: "PARTIAL-LIST" },
      ],
    }),
    updateOfficeByID: builder.mutation({
      query: ({ id, office }) => ({
        url: `booking_offices/${id}`,
        method: "put",
        body: office,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: "BookingOffice", id: id },
        { type: "BookingOffice", id: "PARTIAL-LIST" },
      ],
    }),

    getCashierByID: builder.query({
      query: ({ id }) => `cashiers/${id}`,
      providesTags: (result, error, params) =>
        result ? [{ type: "Cashier", id: result.id }] : [],
    }),
    deleteCashierByID: builder.mutation({
      query: ({ id }) => ({
        url: `cashiers/${id}`,
        method: "delete",
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: "Cashier", id: id },
        { type: "Cashier", id: "PARTIAL-LIST" },
      ],
    }),
    updateCashierByID: builder.mutation({
      query: ({ id, cashier }) => ({
        url: `cashiers/${id}`,
        method: "put",
        body: cashier,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: "Cashier", id: id },
        { type: "Cashier", id: "PARTIAL-LIST" },
      ],
    }),
    updateCashierPassword: builder.mutation({
      query: ({ id, new_password }) => ({
        url: `cashiers/${id}/password`,
        method: "put",
        body: {
          password: new_password,
        },
      }),
    }),

    getFlightInTicketByID: builder.query({
      query: ({ id }) => `flight_in_tickets/${id}`,
      providesTags: (result, error, params) =>
        result ? [{ type: "FlightInTicket", id: result.id }] : [],
    }),
    deleteFlightInTicketByID: builder.mutation({
      query: ({ id }) => ({
        url: `flight_in_tickets/${id}`,
        method: "delete",
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: "FlightInTicket", id: id },
        { type: "FlightInTicket", id: "PARTIAL-LIST" },
      ],
    }),
    updateFlightInTicketByID: builder.mutation({
      query: ({ id, flight_in_ticket }) => ({
        url: `flight_in_tickets/${id}`,
        method: "put",
        body: flight_in_ticket,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: "FlightInTicket", id: id },
        { type: "FlightInTicket", id: "PARTIAL-LIST" },
      ],
    }),

    getFlightByID: builder.query({
      query: ({ id }) => `flights/${id}`,
      providesTags: (result, error, params) =>
        result ? [{ type: "Flight", id: result.id }] : [],
    }),
    deleteFlightByID: builder.mutation({
      query: ({ id }) => ({
        url: `flights/${id}`,
        method: "delete",
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: "Flight", id: id },
        { type: "Flight", id: "PARTIAL-LIST" },
      ],
    }),
    updateFlightByID: builder.mutation({
      query: ({ id, flight }) => ({
        url: `flights/${id}`,
        method: "put",
        body: flight,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: "Flight", id: id },
        { type: "Flight", id: "PARTIAL-LIST" },
      ],
    }),

    getLineByCode: builder.query({
      query: ({ code }) => `lines/${code}`,
      providesTags: (result, error, params) =>
        result ? [{ type: "Line", id: result.line_code }] : [],
    }),
    deleteLineByCode: builder.mutation({
      query: ({ code }) => ({
        url: `lines/${code}`,
        method: "delete",
      }),
      invalidatesTags: (result, error, { code }) => [
        { type: "Line", id: code },
        { type: "Line", id: "PARTIAL-LIST" },
      ],
    }),
    updateLineByCode: builder.mutation({
      query: ({ code, line }) => ({
        url: `lines/${code}`,
        method: "put",
        body: line,
      }),
      invalidatesTags: (result, error, { code }) => [
        { type: "Line", id: code },
        { type: "Line", id: "PARTIAL-LIST" },
        { type: "Flight", id: "PARTIAL-LIST" },
      ],
    }),

    getLinerByCode: builder.query({
      query: ({ code }) => `liners/${code}`,
    }),
    deleteLinerByCode: builder.mutation({
      query: ({ code }) => ({
        url: `liners/${code}`,
        method: "delete",
      }),
      invalidatesTags: (result, error, { code }) => [
        { type: "Liner", id: code },
        { type: "Liner", id: "PARTIAL-LIST" },
      ],
    }),
    updateLinerByCode: builder.mutation({
      query: ({ code, liner }) => ({
        url: `liners/${code}`,
        method: "put",
        body: liner,
      }),
      invalidatesTags: (result, error, { code }) => [
        { type: "Liner", id: code },
        { type: "Liner", id: "PARTIAL-LIST" },
        { type: "Flight", id: "PARTIAL-LIST" },
      ],
    }),

    getLinerModelByCode: builder.query({
      query: ({ code }) => `liner_models/${code}`,
      providesTags: (result, error, params) =>
        result ? [{ type: "Line", id: result.iata_code }] : [],
    }),
    deleteLinerModelByCode: builder.mutation({
      query: ({ code }) => ({
        url: `liner_models/${code}`,
        method: "delete",
      }),
      invalidatesTags: (result, error, { code }) => [
        { type: "LinerModel", id: code },
        { type: "LinerModel", id: "PARTIAL-LIST" },
      ],
    }),
    updateLinerModelByCode: builder.mutation({
      query: ({ code, liner_model }) => ({
        url: `liner_models/${code}`,
        method: "put",
        body: liner_model,
      }),
      invalidatesTags: (result, error, { code }) => [
        { type: "LinerModel", id: code },
        { type: "LinerModel", id: "PARTIAL-LIST" },
        { type: "Liner", id: "PARTIAL-LIST" },
        { type: "Seat", id: "PARTIAL-LIST" },
      ],
    }),

    getPurchaseByID: builder.query({
      query: ({ id }) => `purchases/${id}`,
      providesTags: (result, error, params) =>
        result ? [{ type: "Purchase", id: result.id }] : [],
    }),
    deletePurchaseByID: builder.mutation({
      query: ({ id }) => ({
        url: `purchases/${id}`,
        method: "delete",
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: "Purchase", id: id },
        { type: "Purchase", id: "PARTIAL-LIST" },
      ],
    }),
    updatePurchaseByID: builder.mutation({
      query: ({ id, purchase }) => ({
        url: `purchases/${id}`,
        method: "put",
        body: purchase,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: "Purchase", id: id },
        { type: "Purchase", id: "PARTIAL-LIST" },
      ],
    }),

    getSeatByID: builder.query({
      query: ({ id }) => `seats/${id}`,
      providesTags: (result, error, params) =>
        result ? [{ type: "Seat", id: result.id }] : [],
    }),
    deleteSeatByID: builder.mutation({
      query: ({ id }) => ({
        url: `seats/${id}`,
        method: "delete",
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: "Seat", id: id },
        { type: "Seat", id: "PARTIAL-LIST" },
      ],
    }),
    updateSeatByID: builder.mutation({
      query: ({ id, seat }) => ({
        url: `seats/${id}`,
        method: "put",
        body: seat,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: "Seat", id: id },
        { type: "Seat", id: "PARTIAL-LIST" },
      ],
    }),

    getTicketByID: builder.query({
      query: ({ id }) => `tickets/${id}`,
      providesTags: (result, error, params) =>
        result ? [{ type: "Ticket", id: result.id }] : [],
    }),
    deleteTicketByID: builder.mutation({
      query: ({ id }) => ({
        url: `tickets/${id}`,
        method: "delete",
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: "Ticket", id: id },
        { type: "Ticket", id: "PARTIAL-LIST" },
      ],
    }),
    updateTicketByID: builder.mutation({
      query: ({ id, ticket }) => ({
        url: `tickets/${id}`,
        method: "put",
        body: ticket,
      }),
      invalidatesTags: (result, error, { id }) => [
        { type: "Ticket", id: id },
        { type: "Ticket", id: "PARTIAL-LIST" },
      ],
    }),

    getReportByTicketID: builder.query({
      query: ({ id }) => `tickets/${id}/report`,
      providesTags: (result, error, { id }) => [{ type: "Report", id: id }],
    }),
  }),
});

export const {
  useSignupMutation,
  useLoginMutation,
  useWhoamiQuery,
  useGetTimezonesQuery,
  useGetAirportsQuery,
  useCreateAirportMutation,
  useGetOfficesQuery,
  useCreateOfficeMutation,
  useGetCashiersQuery,
  useCreateCashierMutation,
  useGetFlightInTicketsQuery,
  useCreateFlightInTicketMutation,
  useGetFlightsQuery,
  useCreateFlightMutation,
  useGetLinesQuery,
  useCreateLineMutation,
  useGetLinersQuery,
  useCreateLinerMutation,
  useGetLinerModelsQuery,
  useCreateLinerModelMutation,
  useGetPurchasesQuery,
  useCreatePurchaseMutation,
  useGetSeatsQuery,
  useCreateSeatMutation,
  useGetTicketsQuery,
  useCreateTicketMutation,
  useGetAirportByCodeQuery,
  useDeleteAirportByCodeMutation,
  useUpdateAirportByCodeMutation,
  useGetOfficeByIDQuery,
  useDeleteOfficeByIDMutation,
  useUpdateOfficeByIDMutation,
  useGetCashierByIDQuery,
  useDeleteCashierByIDMutation,
  useUpdateCashierByIDMutation,
  useUpdateCashierPasswordMutation,
  useGetFlightInTicketByIDQuery,
  useDeleteFlightInTicketByIDMutation,
  useUpdateFlightInTicketByIDMutation,
  useGetFlightByIDQuery,
  useDeleteFlightByIDMutation,
  useUpdateFlightByIDMutation,
  useGetLineByCodeQuery,
  useDeleteLineByCodeMutation,
  useUpdateLineByCodeMutation,
  useGetLinerByCodeQuery,
  useDeleteLinerByCodeMutation,
  useUpdateLinerByCodeMutation,
  useGetLinerModelByCodeQuery,
  useDeleteLinerModelByCodeMutation,
  useUpdateLinerModelByCodeMutation,
  useGetPurchaseByIDQuery,
  useDeletePurchaseByIDMutation,
  useUpdatePurchaseByIDMutation,
  useGetSeatByIDQuery,
  useDeleteSeatByIDMutation,
  useUpdateSeatByIDMutation,
  useGetTicketByIDQuery,
  useDeleteTicketByIDMutation,
  useUpdateTicketByIDMutation,
  useGetReportByTicketIDQuery,
} = api;
