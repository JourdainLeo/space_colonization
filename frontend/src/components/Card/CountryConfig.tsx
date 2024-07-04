import React from "react";

function formatName2(name: string) {
  return name
    .split("_")
    .map((word) => word.toUpperCase())
    .join(" ");
}

const CountryConfig = ({
  name,
  index,
  init,
  setInit,
}: {
  name: string;
  index: number;
  init: { [key: string]: number };
  setInit: React.Dispatch<React.SetStateAction<{ [key: string]: number }>>;
}) => {
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
            <input
              onChange={(e) => {
                setInit({ ...init, [name]: Number(e.target.value) });
              }}
              value={init[name]}
              style={{ borderRadius: 8, border: "none", padding: 6 }}
            />
          </div>
        </div>
      </div>
    </div>
  );
};

export default CountryConfig;
