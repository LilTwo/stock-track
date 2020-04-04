import React, {ReactElement, useEffect, useState} from "react";
import moment from "moment";
import "react-dates/initialize";
import "react-dates/lib/css/_datepicker.css";
import { DateRangePicker } from "react-dates";
import { StockType } from "../../types/backend-types";

export type SelectDate = {
  startDate: moment.Moment | null;
  endDate: moment.Moment | null;
};

export function StockDateRangePicker(props: {
  startDate?: moment.Moment;
  endDate?: moment.Moment;
  onSelect: (dates: SelectDate) => void;
}): ReactElement {
  const [startDate, setStartDate] = useState(props.startDate || null);
  const [endDate, setEndDate] = useState(props.endDate || null);
  const [focusedInput, setFocusedInput] = useState<string | null>(null);

  useEffect(
    (startDate?: moment.Moment, endDate?: moment.Moment) => {
      setStartDate(startDate || null);
      setEndDate(endDate || null);
    },
    [props.startDate, props.endDate]
  );

  const handleFocusChange = (focusedInput: string): void => {
    setFocusedInput(focusedInput);
  };

  const handleDateChange = (dates: SelectDate): void => {
    const { startDate, endDate } = dates;
    setStartDate(startDate);
    setEndDate(endDate);
    props.onSelect(dates);
  };
  return (
    <DateRangePicker
      onDatesChange={handleDateChange}
      onFocusChange={handleFocusChange}
      startDate={startDate}
      endDate={endDate}
      startDateId={"start_date_id"}
      endDateId={"end_date_id"}
      focusedInput={focusedInput}
      isOutsideRange={(): boolean => false}
    />
  );
}
