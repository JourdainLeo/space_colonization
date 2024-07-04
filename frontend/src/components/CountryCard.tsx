import CardBar from "./Card/CardBar.tsx";
import Context from "../utils/context.ts";
import { useContext, useState } from "react";
import CountryData from "./Card/CountryData.tsx";
import CountryConfig from "./Card/CountryConfig.tsx";
import CountrySummary from "./Card/CountrySummary.tsx";

const CountryCard = () => {
  const { data, started, sendMessage, end } = useContext(Context);
  const country = [
    "france",
    "north_korea",
    "south_korea",
    "united_states",
    "united_kingdom",
    "india",
    "china",
    "russia",
    "central_africa",
    "australia",
    "italy",
    "germany",
    "mexico",
    "canada",
    "japan",
  ];
  const [init, setInit] = useState<{ [key: string]: number }>({
    france: 0,
    north_korea: 0,
    south_korea: 0,
    united_states: 0,
    united_kingdom: 0,
    india: 0,
    china: 0,
    russia: 0,
    central_africa: 0,
    australia: 0,
    italy: 0,
    germany: 0,
    mexico: 0,
    canada: 0,
    japan: 0,
  });

  const shuffle = () => {
    const temp = { ...init };
    // for each country put a random number between 10000 and infinite
    Object.keys(init).map((c) => {
      temp[c] = Math.floor(Math.random() * 100000);
    });
    setInit(temp);
  };

  return (
    <div
      className="card "
      style={{
        right: 0,
        width: 375,
      }}
    >
      <div
        className="card-body"
        style={{
          height: "97%",
          padding: 0,
          display: "flex",
          flexDirection: "column",
        }}
      >
        <CardBar
          name={end ? "Summary" : started ? "Countries" : "Config"}
          style={"p12"}
        />
        {end ? (
          <div style={{ overflowY: "auto" }}>
            {Object.entries(data.countries).map(([key], index) => {
              return <CountrySummary key={key} name={key} index={index} />;
            })}
          </div>
        ) : started ? (
          <div style={{ overflowY: "auto" }}>
            {Object.entries(data.countries).map(([key], index) => {
              return <CountryData key={key} name={key} index={index} />;
            })}
          </div>
        ) : (
          <div>
            <div
              style={{
                display: "flex",
                flexDirection: "row",
                justifyContent: "center",
                marginBottom: 12,
                padding: 12,
                gap: 12,
              }}
            >
              <button
                className={"button"}
                style={{
                  color: "black",
                  border: "black solid 1px",
                  fontSize: 18,
                }}
                onClick={() => {
                  shuffle();
                }}
              >
                SHUFFLE
              </button>
              <button
                className={"button"}
                style={{
                  color: "black",
                  border: "black solid 1px",
                  fontSize: 18,
                }}
                onClick={() => {
                  sendMessage({ countries: init });
                }}
              >
                START
              </button>
            </div>
            {Object.values(country).map((key, index) => {
              return (
                <CountryConfig
                  key={index + key}
                  name={key}
                  index={index}
                  init={init}
                  setInit={setInit}
                />
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
};

export default CountryCard;
