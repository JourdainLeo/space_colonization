// SolarSystem component
import "../App.css";
import Earth from "./Earth.tsx";
import { OrbitControls } from "@react-three/drei";
import Sun from "./Sun.tsx";
import Planet from "./Planet.tsx";
import { solar_system } from "../utils/base_data.tsx";
import Rocket from "./Rocket.tsx";
import Context from "../utils/context.ts";
import { useContext } from "react";
import { SkyBox } from "./Skybox.tsx";

const SolarSystem = () => {
  const { data } = useContext(Context);

  return (
    <>
      {/*<Perf/>*/}
      <color attach="background" args={["black"]} />
      {<OrbitControls />}
      <ambientLight intensity={1.5} />
      <Sun />
      <SkyBox />
      <Earth name={"earth"} />

      {Object.entries(solar_system.planets).map(([key, value]) => {
        if (key !== "earth") {
          return <Planet key={key} name={key} ring={value.ring} />;
        }
      })}
      {data.rockets &&
        Object.entries(data.rockets).map(([key]) => {
          return <Rocket key={key} id={key} />;
        })}
    </>
  );
};

export default SolarSystem;
