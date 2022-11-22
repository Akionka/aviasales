export const localDatetimeToUTC = (dt) =>
  new Date(
    Date.UTC(
      dt.getUTCFullYear(),
      dt.getUTCMonth(),
      dt.getUTCDate(),
      dt.getUTCHours(),
      dt.getUTCMinutes(),
      dt.getUTCSeconds(),
      dt.getUTCMilliseconds()
    ) +
      new Date().getTimezoneOffset() * 60000
  );