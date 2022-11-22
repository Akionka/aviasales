export const formatPassportNumber = (passport_number) =>
  passport_number &&
  `${passport_number.slice(0, 4)} ${passport_number.slice(4, 10)}`;
