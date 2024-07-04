import DataContainer from "./DataContainer.tsx";
import Separator from "./Separator.tsx";
import { useContext, useState } from "react";
import Context from "../../utils/context.ts";
import RocketLaunchIcon from "@mui/icons-material/RocketLaunch";
import GroupsIcon from "@mui/icons-material/Groups";
import PublicIcon from "@mui/icons-material/Public";
import { Diamond, PsychologyAlt } from "@mui/icons-material";

const action = {
  observation: "observing ",
  preparation: "preparing the launch to ",
  launch: "waiting rocket to arrived at ",
  perception: "perceiving ",
  alliance: "alliance with ",
  ended: "job done.",
};

function formatName2(name: string) {
  return name
    .split("_")
    .map((word) => word.toUpperCase())
    .join(" ");
}

const CountryData = ({ name, index }: { name: string; index: number }) => {
  const { data, setFollowing } = useContext(Context);
  const [expanded, setExpanded] = useState(false);

  return (
    <div
      className={"rocket " + (index % 2 === 0 ? "even" : "")}
      style={{
        display: "flex",
        flexDirection: "column",
        paddingLeft: 6,
        paddingRight: 12,
        paddingTop: 6,
        paddingBottom: 6,
        borderLeft: !data.countries[name].lastRocketStatus
          ? "gray solid 6px"
          : data.countries[name].lastRocketStatus === "landing"
            ? "green solid 6px"
            : "red solid 6px",
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
            <div
              className={"button-close"}
              style={{
                width: 14,
                height: 14,
                borderRadius: 6,
                background: "none",
              }}
              onClick={() => {
                setExpanded(!expanded);
              }}
            >
              {expanded ? "-" : "+"}
            </div>
          </div>
          <div
            style={{
              display: "flex",
              flexDirection: "row",
              alignItems: "flex-end",
              gap: 8,
            }}
          >
            <div style={{ margin: 0, fontSize: 14 }}>
              {action[data.countries[name].action as keyof typeof action]}
              {" " + data.countries[name].target !== "" &&
              data.countries[name].action !== "ended"
                ? data.countries[name].target.charAt(0).toUpperCase() +
                  data.countries[name].target.slice(1)
                : ""}
            </div>
            {data.countries[name].action !== "ended" && (
              <div className={"dot-flashing"} style={{ marginBottom: 3 }}></div>
            )}
          </div>
        </div>
      </div>
      <div className={"test-data " + (expanded ? "expanded" : "collapse")}>
        <div
          className="card-display"
          style={{
            fontSize: 12,
            marginBottom: 0,
            flex: 1,
            marginTop: 8,
            padding: 4,
          }}
        >
          <DataContainer
            icon={<Diamond fontSize={"small"} />}
            text={data.countries[name].resources.toString()}
          />
          <Separator />
          <DataContainer
            icon={<PsychologyAlt fontSize={"small"} />}
            text={data.countries[name].iq.toString()}
          />
        </div>
        <div
          style={{
            display: "flex",
            flexDirection: "row",
            alignItems: "center",
            marginTop: 8,
            gap: 8,
            color: "gray",
          }}
        >
          <GroupsIcon />
          <div
            className="card-display"
            style={{ fontSize: 12, marginBottom: 0, flex: 1 }}
          >
            {Object.values(data.countries[name].alliances).map(
              (value, index) => {
                return (
                  <>
                    {index !== 0 && <Separator />}
                    <DataContainer icon={""} text={value} />
                  </>
                );
              },
            )}
          </div>
        </div>

        <div
          style={{
            display: "flex",
            flexDirection: "row",
            alignItems: "center",
            marginTop: 8,
            gap: 8,
            color: "gray",
          }}
        >
          <PublicIcon />
          <div
            className="card-display"
            style={{ fontSize: 12, marginBottom: 0, flex: 1 }}
          >
            {Object.values(data.countries[name].planets).map((value, index) => {
              return (
                <>
                  {index !== 0 && <Separator />}
                  <DataContainer
                    icon={""}
                    text={value}
                    hover
                    func={() => {
                      setFollowing(value);
                    }}
                  />
                </>
              );
            })}
          </div>
        </div>

        {data.rockets[name] && (
          <div
            style={{
              display: "flex",
              flexDirection: "row",
              alignItems: "center",
              marginTop: 8,
              gap: 8,
              color: "gray",
            }}
          >
            <RocketLaunchIcon />
            <div
              className="card-display"
              style={{ fontSize: 12, marginBottom: 0, flex: 1 }}
            >
              <DataContainer
                icon={"x: "}
                text={data.rockets[name].x.toFixed(2)}
              />
              <Separator />
              <DataContainer
                icon={"y: "}
                text={data.rockets[name].y.toFixed(2)}
              />
              <Separator />
              <DataContainer
                icon={"z: "}
                text={data.rockets[name].z.toFixed(2)}
              />
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default CountryData;
