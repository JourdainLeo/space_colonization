import { createContext } from "react";
import {MsgType, SpaceData} from "./model.tsx";
import { solar_system } from "./base_data.tsx";

// Any for the setters here are actually type React.Dispatch<React.SetStateAction<TYPE>>
// But I don't know ho to init in the Context these value when it's not the place where the states are (they are in App.tsx)
// And the goal is not to have TYPE | null because I will have to check every time and that's no the goal
interface IContext {
  init: { [key: string]: number };
  setInit: any;
  started: boolean;
  hovered: string;
  setHovered: any;
  following: string;
  setFollowing: any;
  data: SpaceData;
  setData: any;
  planetGeometries: {
    [key: string]: THREE.BufferGeometry<THREE.NormalBufferAttributes>;
  };
  setPlanetGeometries: any;
  sendMessage: (msg: MsgType) => void;
  end: boolean;
}

const Context = createContext<IContext>({
  init: {},
  setInit: false,
  started: false,
  hovered: "",
  setHovered: false,
  following: "",
  setFollowing: false,
  data: solar_system,
  setData: false,
  planetGeometries: {},
  setPlanetGeometries: false,
  sendMessage: () => {},
  end: false,
});
export default Context;
