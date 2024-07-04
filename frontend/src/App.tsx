import "./App.css";
import { useState } from "react";
import SpaceCard from "./components/SpaceCard.tsx";
import { MsgType, SpaceData } from "./utils/model.tsx";
import { solar_system } from "./utils/base_data.tsx";
import useWebSocket from "react-use-websocket";
import Context from "./utils/context.ts";
import {
  receiveCountry,
  receiveEnd,
  receiveRocket,
  receiveStar,
} from "./utils/utils.ts";
import CountryCard from "./components/CountryCard.tsx";
import { Canvas } from "@react-three/fiber";
import SolarSystem from "./components/SolarSystem.tsx";

const WS_URL = "ws://127.0.0.1:8080";

function App() {
  const [following, setFollowing] = useState("");
  const [data, setData] = useState<SpaceData>(solar_system);
  const [hovered, setHovered] = useState("");
  const [planetGeometries, setPlanetGeometries] = useState<any>({});
  const [started, setStarted] = useState(false);
  const [init, setInit] = useState({});
  const [end, setEnd] = useState(false);
  const { sendJsonMessage } = useWebSocket(WS_URL, {
    shouldReconnect: () => true,
    onOpen: () => {
      console.log("WebSocket connection established.");
    },

    onMessage: (event) => {
      const msg = JSON.parse(event.data);
      const newData = { ...data };
      let res;

      setStarted(true);
      switch (msg.type) {
        case "star":
          setData(receiveStar(msg, newData));
          break;
        case "end":
          setEnd(true);
          setData(receiveEnd(msg, newData));
          break;
        case "rocket":
          setData(receiveRocket(msg, newData));
          break;
        case "country":
          res = receiveCountry(msg, newData);
          if (res) {
            setData(res);
          }
          break;
      }
    },
  });
  const sendMessage = (message: MsgType) => {
    sendJsonMessage(message);
  };

  return (
    <Context.Provider
      value={{
        following,
        setFollowing,
        data,
        setData,
        hovered,
        setHovered,
        planetGeometries,
        setPlanetGeometries,
        started,
        init,
        setInit,
        sendMessage,
        end,
      }}
    >
      <SpaceCard />
      <CountryCard />
      <Canvas
        style={{ flex: 16 }}
        camera={{
          fov: 45,
          near: 0.1,
          far: 1000000,
          position: [0, 1000, 800],
          aspect: window.innerWidth / window.innerHeight,
        }}
        frameloop={"always"}
      >
        <SolarSystem />
      </Canvas>
    </Context.Provider>
  );
}

export default App;
