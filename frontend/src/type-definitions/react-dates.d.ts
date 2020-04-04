declare module "react-dates" {
  import moment from "moment";
  export type DateRangePickerProps = {
    startDateId: string;
    startDate: moment.Moment | null;
    endDateId: string;
    endDate: moment.Moment | null;
    focusedInput: string | null;
    onFocusChange(focusedInput: string): void;
    onDatesChange(dates: {
      startDate: moment.Moment | null;
      endDate: moment.Moment | null;
    }): void;
    isOutsideRange(date: moment.Moment): boolean;
  };

  declare class DateRangePicker extends React.ComponentElement<
    DateRangePickerProps,
    any
  > {
    render(): JSX.Element;
    context: any;
    setState: any;
    forceUpdate: any;
    state: any;
    props: DateRangePickerProps;
    refs: any;
  }
  export { DateRangePicker };
}
