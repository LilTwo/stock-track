import React, { useState } from "react";
import ReactDOM from "react-dom";
import Plot from "react-plotlyjs-ts";
import { ajax } from "rxjs/ajax";
import {
  FilterArgs,
  TradeInfoFilter,
} from "./components/trade-info-filter/trade-info-filter";
import {filter} from "rxjs/operators";

type PlotData = {
  x: number | string[];
  y: number[];
  mode?: string;
  type?: string;
  name: string;
  hoverInfo: string;
};

function composeUrl(filterArgs: FilterArgs): string {
  const params = new URLSearchParams();
  const { stockType, startDate, endDate } = filterArgs;
  params.set("code", stockType?.Code.toString() || "");
  params.set("startDate", startDate || "");
  params.set("endDate", endDate || "");
  return `/api/trade-info?${params.toString()}`;
}

function Root() {
  const numLines = 2;
  const [plotData, setPlotData] = useState<PlotData[]>(
    Array(numLines).fill({} as PlotData)
  );

  const handleApply = (filterArgs: FilterArgs, index: number) => {
    console.log("handle", filterArgs, index);
    ajax.get(composeUrl(filterArgs)).subscribe((resp) => {
      const data: { Date: string; Data: number }[] | null = resp.response;
      if (data === null) {
        // #TODO: figure out a way to indicate this state
        console.log("no data avaliable");
        return;
      }
      const lineName = `${filterArgs.stockType?.Code} ${filterArgs.stockType?.Name}`;

      const line: PlotData = {
        x: data.map((d) => d.Date),
        y: data.map((d) => d.Data),
        type: "scatter",
        name: lineName,
        hoverInfo: lineName,
      };
      plotData[index] = line;
      setPlotData([...plotData]);
    });
  };

  const filters = Array(numLines).fill(0).map((_, index) => {
      const localHandleApply = (filterArgs: FilterArgs) => {
        handleApply(filterArgs, index);
      };
      return <TradeInfoFilter onApply={localHandleApply} key={index} />;
    });

  return (
    <div>
      <Plot data={plotData} />
      {filters}
    </div>
  );
}

ReactDOM.render(<Root />, document.getElementById("root"));
