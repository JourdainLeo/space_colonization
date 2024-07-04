import { SpaceData } from "./model.tsx";

export const receiveEnd = (msg: any, data: SpaceData) => {
  const status = msg.status;

  Object.values(status).map((country: any) => {
    data.countries[country.country].resources = country.budget;
    data.countries[country.country].presence = country.presence;
    data.countries[country.country].territory = country.planets;
    data.countries[country.country].sObservation =
      country.strategy[0].observation;
    data.countries[country.country].sAlliance = country.strategy[0].alliance;
    data.countries[country.country].sLaunch = country.strategy[0].launch;
    data.countries[country.country].sSkipping = country.strategy[0].skipping;
  });

  return data;
};

export const receiveStar = (msg: any, data: SpaceData) => {
  // msg is an object with planets and satellites

  const d: any = msg.data;

  Object.values(d).map((star: any) => {
    if (star.type === "planet") {
      data.planets[star.name].x = star.x;
      data.planets[star.name].y = star.y;
      data.planets[star.name].z = star.z;
      data.planets[star.name].resources = star.resources;
      data.planets[star.name].owners = star.owners;
    } else if (star.type === "satellite") {
      data.satellites[star.name].x = star.x;
      data.satellites[star.name].y = star.y;
      data.satellites[star.name].z = star.z;
      data.satellites[star.name].resources = star.resources;
      data.satellites[star.name].owners = star.owners;
    }
  });
  return data;
};
export const receivePlanet = (msg: any, data: SpaceData) => {
  data.planets[msg.name].x = msg.x;
  data.planets[msg.name].y = msg.y;
  data.planets[msg.name].z = msg.z;
  data.planets[msg.name].resources = msg.resources;
  data.planets[msg.name].owners = msg.owners;
  return data;
};

export const receiveSatellite = (msg: any, data: SpaceData) => {
  data.satellites[msg.name].x = msg.x;
  data.satellites[msg.name].y = msg.y;
  data.satellites[msg.name].z = msg.z;
  data.satellites[msg.name].resources = msg.resources;
  data.satellites[msg.name].owners = msg.owners;
  return data;
};

export const receiveCountry = (msg: any, data: SpaceData) => {
  if (!data.countries[msg.country]) {
    data.countries[msg.country] = {
      action: "",
      target: "",
      resources: 0,
      planets: [],
      alliances: [],
      iq: 0,
    };
  }

  data.countries[msg.country].action = msg.action;
  if (msg.action === "alliance")
    data.countries[msg.country].alliances.push(msg.target);
  if (msg.qi) data.countries[msg.country].iq = msg.qi;
  data.countries[msg.country].planets = msg.planets;
  if (msg.target) data.countries[msg.country].target = msg.target;
  data.countries[msg.country].resources = msg.budget;

  return data;
};

export const receiveRocket = (msg: any, data: SpaceData) => {
  if (!data.rockets[msg.sender]) {
    data.rockets[msg.sender] = {
      x: 0,
      y: 0,
      z: 0,
      destination: "",
      sender: "",
      status: "",
      distanceTo: 0,
    };
  }
  data.rockets[msg.sender].x = msg.x;
  data.rockets[msg.sender].y = msg.y;
  data.rockets[msg.sender].z = msg.z;
  data.rockets[msg.sender].destination = msg.destination;
  data.rockets[msg.sender].sender = msg.sender;
  data.rockets[msg.sender].status = msg.status;

  if (data.rockets[msg.sender] && msg.status !== "alive") {
    data.countries[msg.sender].lastRocketStatus = msg.status;
    delete data.rockets[msg.sender];
  }
  return data;
};
