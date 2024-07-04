//http://hyperphysics.phy-astr.gsu.edu/hbase/Solar/soldata3.html#c1
// distance : UA * 100
// velocity : km/s * 10^-2
// revolution : (years) -> multiply by 10^2 and then divide by 2

// Components props

export interface IPlanetProps {
  name: string;
  ring?: boolean;
}

export interface Star {
  size: number;
  velocity: number;
  revolution: number;
  x: number;
  y: number;
  z: number;
  owners: string[];
  resources: number;
}

// Data type
export interface SatelliteType extends Star {
  parent: string;
  distance: number;
}

export type SatellitesData = { [key: string]: SatelliteType };

interface PlanetType extends Star {
  minDistance: number;
  maxDistance: number;
  color: string;
  ring?: boolean;
  ringRotation?: number;
  ringTextureRotation?: number;
  ringDistance?: number;
  ringSize?: number;
  satellites?: string[];
}

interface RocketType {
  x: number;
  y: number;
  z: number;
  destination: string;
  sender: string;
  status: string;
  distanceTo: number;
}

export type RocketsData = { [key: string]: RocketType };

export type PlanetsData = { [key: string]: PlanetType };

interface CountryType {
  action: string;
  target: string;
  resources: number;
  alliances: string[];
  planets: string[];
  iq: number;
  presence?: number;
  territory?: number;
  sObservation?: string;
  sLaunch?: string;
  sSkipping?: string;
  sAlliance?: string;
  lastRocketStatus?: string;
}

export type CountriesData = { [key: string]: CountryType };

export interface MsgType {
  countries: { [key: string]: number };
}

export interface SpaceData {
  planets: PlanetsData;
  satellites: SatellitesData;
  rockets: RocketsData;
  countries: CountriesData;
}
