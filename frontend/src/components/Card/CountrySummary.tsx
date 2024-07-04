import { useContext } from "react";
import Context from "../../utils/context.ts";

function formatName2(name: string) {
  return name
    .split("_")
    .map((word) => word.toUpperCase())
    .join(" ");
}

const CountrySummary = ({ name, index }: { name: string; index: number }) => {
  const { data } = useContext(Context);
  return (
    <div
      className={"rocket " + (index % 2 === 0 ? "even" : "")}
      style={{
        display: "flex",
        flexDirection: "column",
        paddingLeft: 12,
        paddingRight: 12,
        paddingTop: 6,
        paddingBottom: 6,
      }}
    >
      <div
        style={{
          display: "flex",
          flex: 1,
          alignItems: "center",
          gap: 8,
        }}
      >
        <div className={name + " flag"} style={{ height: 26, width: 40 }} />
        <div
          style={{
            display: "flex",
            flexDirection: "column",
            flex: 1,
          }}
        >
          <div
            style={{
              display: "flex",
              alignItems: "center",
              justifyContent: "space-between",
              flex: 1,
            }}
          >
            <div style={{ margin: 0, fontSize: 16 }} className={"country-name"}>
              {formatName2(name)}
            </div>
          </div>
        </div>
      </div>

      <div style={{ fontSize: 16, marginTop: 8 }}>
        planets:{" " + data.countries[name].presence}
      </div>
      <div style={{ fontSize: 16 }}>
        territory:{" " + data.countries[name].territory}/15
        {data.countries[name].territory
          ? " (" +
            ((data.countries[name].territory! * 100) / 15).toFixed(2) +
            "%)"
          : ""}
      </div>
      <div style={{ fontSize: 16 }}>
        alliances:{" " + data.countries[name].alliances.length}
      </div>
      <div style={{ fontSize: 16, marginTop: 8, textDecoration: "underline" }}>
        strategy:
      </div>
      <div style={{ fontSize: 16 }}>
        observation:{" " + data.countries[name].sObservation}
      </div>
      <div style={{ fontSize: 16 }}>
        alliance:{" " + data.countries[name].sAlliance}
      </div>

      <div style={{ fontSize: 16 }}>
        skipping:{" " + data.countries[name].sSkipping}
      </div>
      <div style={{ fontSize: 16 }}>
        launch:{" " + data.countries[name].sLaunch}
      </div>
    </div>
  );
};

export default CountrySummary;
