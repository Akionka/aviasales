// Файл web\src\features\report\reportPage.jsx содержит код для страницы отчёта
import { useEffect } from "react";
import { useParams } from "react-router-dom";
import { useGetReportByTicketIDQuery } from "../../app/services/api";
import { formatPhoneNumberIntl } from "react-phone-number-input";
import { localDatetimeToUTC } from "../../utils/dateConverter";

import styles from "./reportPage.module.css";

const GridItem = ({ value, label }) => {
  return (
    <div className={styles.ticket__grid_item}>
      <div className={styles.ticket__itemlabel}>{label}:</div>
      <div className={styles.ticket__itemvalue}>{value}</div>
    </div>
  );
};

const FlightItem = ({
  cityTo,
  cityFrom,
  depDate,
  depTime,
  arrDate,
  arrTime,
  line,
  seatNumber,
  seatClass,
}) => {
  return (
    <>
      <GridItem value={cityFrom} label="Из" />
      <div className={styles.ticket__subgrid}>
        <GridItem value={depDate} label="Дата" />
        <GridItem value={depTime} label="Время" />
      </div>
      <GridItem value={cityTo} label="В" />
      <div className={styles.ticket__subgrid}>
        <GridItem value={arrDate} label="Дата" />
        <GridItem value={arrTime} label="Время" />
      </div>
      <div className={styles.ticket__subgrid}>
        <GridItem value={line} label="Рейс" />
        <GridItem value={seatNumber} label="Место" />
      </div>
      <GridItem value={seatClass} label="Класс" />
    </>
  );
};

export const ReportPage = () => {
  const { ticketId } = useParams();
  const { data, error, isLoading } = useGetReportByTicketIDQuery({
    id: ticketId,
  });
  useEffect(() => {
    document.title = `Билет № ${ticketId}`;
  }, [ticketId]);
  if (isLoading) return;
  if (error)
    return (
      <div>
        Ошибка! {error.status} {error.data.error}
      </div>
    );
  return (
    <div className={styles.ticket}>
      <div className={styles.ticket__page}>
        <div className={styles.ticket__number}>Билет № {ticketId}</div>
        <div className={styles.ticket__grid}>
          <GridItem
            label="Имя пассажира"
            value={`${data.ticket.passenger_last_name} ${data.ticket.passenger_given_name}`}
          />
          <GridItem label="Перевозчик" value={`Журавли`} />
          <div className={styles.ticket__grid_item}>
            <div className={styles.ticket__subgrid}>
              <GridItem
                label="Дата рождения"
                value={new Date(
                  data.ticket.passenger_birth_date
                ).toLocaleDateString()}
              />
              <GridItem
                label="Пол"
                value={data.ticket.passenger_sex === 1 ? "Мужской" : "Женский"}
              />
            </div>
          </div>
          <GridItem
            label="Паспортные данные"
            value={data.ticket.passenger_passport_number}
          />
          {data.flights.map((f) => (
            <FlightItem
              arrDate={localDatetimeToUTC(
                new Date(f.arr_time_local)
              ).toLocaleDateString()}
              arrTime={localDatetimeToUTC(
                new Date(f.arr_time_local)
              ).toLocaleTimeString([], {
                hour: "2-digit",
                minute: "2-digit",
              })}
              cityFrom={f.dep_city}
              cityTo={f.arr_city}
              depDate={localDatetimeToUTC(
                new Date(f.dep_time_local)
              ).toLocaleDateString()}
              depTime={localDatetimeToUTC(
                new Date(f.dep_time_local)
              ).toLocaleTimeString([], {
                hour: "2-digit",
                minute: "2-digit",
              })}
              line={f.line_code}
              seatClass={
                f.class === "J"
                  ? "Бизнес"
                  : f.class === "W"
                  ? "Комфорт"
                  : "Эконом"
              }
              seatNumber={f.number}
            />
          ))}
          <div className={styles.ticket__grid_item} />
          <GridItem
            label="Общее время в пути"
            value={localDatetimeToUTC(
              new Date(data.total_time * 1000)
            ).toLocaleTimeString([], {
              hour: "2-digit",
              minute: "2-digit",
            })}
          />
          <GridItem label="Место продажи" value={data.booking_office.address} />
          <GridItem
            label="Телефон кассы"
            value={formatPhoneNumberIntl(
              "+" + data.booking_office.phone_number
            )}
          />
          <GridItem
            label="Кассир"
            value={`${data.cashier.last_name} ${data.cashier.first_name} ${data.cashier.middle_name}`}
          />
          <GridItem
            label="Дата и время продажи"
            value={new Date(data.purchase.date).toLocaleString()}
          />
          <GridItem
            label="Контактный телефон пассажира"
            value={formatPhoneNumberIntl("+" + data.purchase.contact_phone)}
          />
          <GridItem
            label="Контактный E-mail пассажира"
            value={data.purchase.contact_email}
          />
        </div>
      </div>
    </div>
  );
};
