import Slide from "@mui/material/Slide";
import { useContext } from "react";
import Context from "../utils/context.ts";
import { Diamond, Height, Speed, ThreeSixty } from "@mui/icons-material";
import CardBar from "./Card/CardBar.tsx";
import Separator from "./Card/Separator.tsx";
import DataContainer from "./Card/DataContainer.tsx";
import NightlightIcon from "@mui/icons-material/Nightlight";
import FlagIcon from "@mui/icons-material/Flag";

const SpaceCard = () => {
  const { following, data, setFollowing } = useContext(Context);

  if (following === "") return <></>;
  const name = data.planets[following]
    ? following
    : data.satellites[following] && data.satellites[following].parent;
  const childName = data.satellites[following] ? following : null;
  const parentData = data.planets[name];
  const childData = data.satellites[following]
    ? data.satellites[following]
    : null;

  return (
    <Slide direction="right" in={following !== ""} mountOnEnter unmountOnExit>
      {parentData && (
        <div className="card" style={{ fontSize: 12 }}>
          <div className="card-body">
            <CardBar name={name} close />
            <div className={"card-header " + name} />
            <div className="card-subtitle">
              {parentData.minDistance +
                " ~ " +
                parentData.maxDistance +
                " distance from The Sun"}
            </div>

            <div className="card-display">
              <DataContainer icon={"x: "} text={parentData.x.toFixed(2)} />
              <Separator />
              <DataContainer icon={"y: "} text={parentData.y.toFixed(2)} />
              <Separator />
              <DataContainer icon={"z: "} text={parentData.z.toFixed(2)} />
            </div>
            <div className="card-display" style={{ padding: 4 }}>
              <DataContainer
                icon={<ThreeSixty fontSize={"small"} />}
                text={parentData.revolution.toFixed(2)}
              />
              <Separator />
              <DataContainer
                icon={<Height fontSize={"small"} />}
                text={parentData.size.toFixed(2)}
              />
              <Separator />
              <DataContainer
                icon={<Speed fontSize={"small"} />}
                text={parentData.velocity.toFixed(2)}
              />
              <Separator />
              <DataContainer
                icon={<Diamond fontSize={"small"} />}
                text={parentData.resources.toString()}
              />
            </div>

            {parentData.satellites && (
              <div
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: 8,
                  color: "gray",
                }}
              >
                <NightlightIcon />
                <div className="card-display" style={{ margin: 0, flex: 1 }}>
                  {Object.values(parentData.satellites!).map((value, index) => {
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
            )}

            <div
              style={{
                display: "flex",
                gap: 8,
                marginTop: 8,
                alignItems: "center",
                color: "gray",
              }}
            >
              <FlagIcon />
              <div className="card-display" style={{ flex: 1, margin: 0 }}>
                {Object.values(
                  name !== "earth"
                    ? parentData.owners
                    : Object.keys(data.countries),
                ).map((value, index) => {
                  return (
                    <>
                      {index !== 0 && index !== 5 && <Separator />}
                      <DataContainer icon={""} text={value} />
                    </>
                  );
                })}
              </div>
            </div>
          </div>

          {childName && childData && (
            <Slide direction="right" in={true} mountOnEnter unmountOnExit>
              <div className={"card-body2"}>
                <CardBar name={childName.toUpperCase()} close />
                <div className={"card-header " + childName} />
                <div className="card-subtitle">
                  {childData.distance +
                    " distance from " +
                    name.charAt(0).toUpperCase() +
                    name.slice(1)}
                </div>
                <div className="card-display">
                  <DataContainer icon={"x: "} text={parentData.x.toFixed(2)} />
                  <Separator />
                  <DataContainer icon={"y: "} text={parentData.y.toFixed(2)} />
                  <Separator />
                  <DataContainer icon={"z: "} text={parentData.z.toFixed(2)} />
                </div>
                <div className="card-display" style={{ padding: 4 }}>
                  <DataContainer
                    icon={<ThreeSixty fontSize={"small"} />}
                    text={parentData.revolution.toFixed(2)}
                  />
                  <Separator />
                  <DataContainer
                    icon={<Height fontSize={"small"} />}
                    text={parentData.size.toFixed(2)}
                  />
                  <Separator />
                  <DataContainer
                    icon={<Speed fontSize={"small"} />}
                    text={parentData.velocity.toFixed(2)}
                  />
                  <Separator />
                  <DataContainer
                    icon={<Diamond fontSize={"small"} />}
                    text={childData.resources.toString()}
                  />
                </div>
                <div
                  style={{
                    display: "flex",
                    gap: 8,
                    marginTop: 8,
                    alignItems: "center",
                    color: "gray",
                  }}
                >
                  <FlagIcon />
                  <div className="card-display" style={{ flex: 1, margin: 0 }}>
                    {Object.values(childData.owners).map((value, index) => {
                      return (
                        <>
                          {index !== 0 && <Separator />}
                          <DataContainer icon={""} text={value} />
                        </>
                      );
                    })}
                  </div>
                </div>
              </div>
            </Slide>
          )}
        </div>
      )}
    </Slide>
  );
};

export default SpaceCard;
