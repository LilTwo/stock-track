import React, { ReactElement, useState } from "react";
import { StockType } from "../../types/backend-types";
import {
  SelectDate,
  StockDateRangePicker,
} from "../date-range-picker/date-range-picker";
import { StockSelect } from "../stock-select/stock-select";
import "./trade-info-filter.scss";

export type FilterArgs = {
  stockType: StockType | null;
  startDate: string | null;
  endDate: string | null;
};

export function TradeInfoFilter(props: {
  filterArgs?: FilterArgs;
  onApply: (filterArgs: FilterArgs) => void;
}): ReactElement {
  // #TODO: need to sync filterArgs into filter components if props.filterArgs Change (useEffect)
  const [stockType, setStockType] = useState<StockType>();
  const [date, setDate] = useState<SelectDate>();

  const handleStockSelect = (newStockType: StockType) => {
    setStockType(newStockType);
  };
  const handleDateSelect = (newDate: SelectDate): void => {
    setDate(newDate);
  };
  const handleClick = () => {
    // only pass YYYY-MM-DD
    const startDate = date?.startDate?.toJSON()?.substring(0, 10) || null;
    const endDate = date?.endDate?.toJSON()?.substring(0, 10) || null;

    props.onApply({
      startDate,
      endDate,
      stockType: stockType || null,
    });
  };

  return (
    <div className={"outer-container"}>
      <div className={"filter-container"}>
        <StockSelect
          onSelect={handleStockSelect}
          className={"filter-stock-select"}
        />
        <StockDateRangePicker onSelect={handleDateSelect} />
      </div>
      <button onClick={handleClick} className={"filter-button"}>
        Apply
      </button>
    </div>
  );
}
