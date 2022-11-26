// Файл web\src\App.js содержит код приложения на верхнем уровне
import React, { useEffect } from "react";
import { useDispatch } from "react-redux";
import { Routes, Route, Navigate } from "react-router-dom";
import { setToken } from "./features/auth/authSlice";

import { Login } from "./features/auth/Login";
import { MainMenu } from "./features/mainMenu/MainMenu";
import { ProtectedRoute } from "./hoc/ProtectedRoute";
import { AirportsPage } from "./features/airport/airportsPage";
import { BookingOfficesPage } from "./features/bookingOffice/bookingOfficesPage";
import { CashiersPage } from "./features/cashier/cashiersPage";
import { FlightInTicketsPage } from "./features/flightInTicket/flightInTicketsPage";
import { FlightsPage } from "./features/flight/flightsPage";
import { LinesPage } from "./features/line/linesPage";
import { LinerModelsPage } from "./features/linerModel/linerModelsPage";
import { LinersPage } from "./features/liner/linersPage";
import { PurchasesPage } from "./features/purchase/purchasesPage";
import { SeatsPage } from "./features/seat/seatsPage";
import { TicketsPage } from "./features/ticket/ticketsPage";
import { ReportPage } from "./features/report/reportPage";

function App() {
  const dispatch = useDispatch();

  useEffect(() => {
    const token = sessionStorage.getItem("auth_token");
    if (token) {
      dispatch(setToken(token));
    }
  }, [dispatch]);

  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route
        path="/"
        element={
          <ProtectedRoute>
            <MainMenu />
          </ProtectedRoute>
        }
      >
        <Route path="airport" element={<AirportsPage />} />
        <Route path="booking_office" element={<BookingOfficesPage />} />
        <Route path="cashier" element={<CashiersPage />} />
        <Route path="flight_in_ticket" element={<FlightInTicketsPage />} />
        <Route path="flight" element={<FlightsPage />} />
        <Route path="line" element={<LinesPage />} />
        <Route path="liner_model" element={<LinerModelsPage />} />
        <Route path="liner" element={<LinersPage />} />
        <Route path="purchase" element={<PurchasesPage />} />
        <Route path="seat" element={<SeatsPage />} />
        <Route path="ticket" element={<TicketsPage />} />
        <Route path="*" element={<Navigate to="/" />} />
      </Route>
      <Route
        path="ticket/:ticketId/report/"
        element={
          <ProtectedRoute>
            <ReportPage />
          </ProtectedRoute>
        }
      />
    </Routes>
  );
}

export default App;
