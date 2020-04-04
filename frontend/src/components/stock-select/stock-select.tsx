import React, { ReactElement, useEffect, useState } from "react";
import { Subject } from "rxjs";
import { startWith, switchMap, throttleTime } from "rxjs/operators";
import { ajax } from "rxjs/ajax";
import Select from "react-select";
import { StockType } from "../../types/backend-types";

type Option = {
  value: StockType;
  label: ReactElement;
};

function toLabel(stock: StockType): ReactElement {
  return (
    <div>
      <b>{stock.Code}</b> {stock.Name}
      <div style={{ fontSize: "11px" }}>{stock.ISN}</div>
    </div>
  );
}

function stockTypeToOption(stockType?: StockType): Option | null {
  return stockType
    ? {
        label: toLabel(stockType),
        value: stockType,
      }
    : null;
}

export function StockSelect(props: {
  onSelect: (stockType: StockType) => void;
  stockType?: StockType;
  className?: string;
}): ReactElement {
  const [options, setOptions] = useState<Option[]>([]);
  const search$ = useState(new Subject<string>())[0];
  const [value, setValue] = useState<Option | null>(null);

  useEffect(
    (stockType?: StockType) => {
      setValue(stockTypeToOption(stockType));
    },
    [props.stockType]
  );

  useEffect(() => {
    search$
      .pipe(
        startWith(""),
        throttleTime(500, undefined, { leading: true, trailing: true }),
        switchMap((search) => ajax.get(`/api/auto-complete?search=${search}`))
      )
      .subscribe((resp) => {
        const data: StockType[] = resp.response;
        const result = data.map((stockType) => ({
          label: toLabel(stockType),
          value: stockType,
        }));
        setOptions(result);
      });
  }, []);

  const handleInputChange = (val: string): void => {
    search$.next(val);
  };

  const handleOnChange = (val: any): void => {
    const stockType = val.value as StockType;
    const valueEle = (
      <span>
        <b>{stockType.Code}</b> {stockType.Name}
      </span>
    );
    setValue({ label: valueEle, value: val.value });
    props.onSelect(stockType);
  };

  return (
    <Select
      className={props.className}
      options={options}
      onInputChange={handleInputChange}
      filterOption={() => true}
      onChange={handleOnChange}
      value={value}
    />
  );
}
