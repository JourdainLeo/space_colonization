import { Text, useTexture } from "@react-three/drei";
import { useFrame, useThree } from "@react-three/fiber";
import { useContext, useEffect, useRef, useState } from "react";
import * as THREE from "three";
import Satellite from "./Satellite.tsx";
import { IPlanetProps } from "../utils/model.tsx";
import Context from "../utils/context.ts";

const Planet = ({ name, ring }: IPlanetProps) => {
  const planetRef = useRef<any>(null!);
  const textRef = useRef<THREE.Mesh>(null!);
  const { camera } = useThree();
  const input = ["/assets/texture/" + name + ".jpg"];
  if (ring) {
    input.push("/assets/texture/" + name + "_ring.png");
  }
  const [planetTexture, ringTexture] = useTexture(input);
  const [fontSize, setFontSize] = useState(2);
  const [opacity, setOpacity] = useState(1);
  const ringRef = useRef<THREE.Mesh>(null!);
  const [zoomed, setZoomed] = useState(false);
  const {
    following,
    setFollowing,
    hovered,
    setHovered,
    data,
    planetGeometries,
    setPlanetGeometries,
  } = useContext(Context);

  if (data.planets[name].ringTextureRotation) {
    ringTexture.rotation = data.planets[name].ringTextureRotation!;
  }

  useEffect(() => {
    document.body.style.cursor = hovered !== "" ? "pointer" : "auto";
  }, [hovered]);

  useFrame(() => {
    // Handle orbit coordinates
    planetRef.current.position.set(
      data.planets[name].x,
      0,
      data.planets[name].z,
    );

    // Handle ring position
    if (ring) {
      ringRef.current.position.set(
        data.planets[name].x,
        data.planets[name].y,
        data.planets[name].z,
      );
    }

    // Handle revolution
    planetRef.current.rotation.y += data.planets[name].revolution;

    // Handle camera follow
    if (following === name) {
      const x =
        ((data.planets[name].maxDistance - data.planets[name].size * 4) /
          data.planets[name].maxDistance) *
        data.planets[name].x;
      const z =
        ((data.planets[name].minDistance - data.planets[name].size * 4) /
          data.planets[name].minDistance) *
        data.planets[name].z;

      if (!zoomed) {
        camera.position.set(
          x,
          planetRef.current.position.y + data.planets[name].size,
          z,
        );
        setZoomed(true);
      } else {
        camera.position.set(x, camera.position.y, z);
      }

      camera.lookAt(planetRef.current.position);
    }

    const distance = planetRef.current.position.distanceTo(camera.position);

    // Handle text opacity
    let opacity;
    const maxD = 800;
    const minD = 300;
    if (distance > maxD) {
      opacity = 1;
    } else if (distance > minD) {
      opacity = (distance - minD) / (maxD - minD);
    } else {
      opacity = 0;
    }
    setOpacity(opacity);

    // Handle text size
    const newFontSize = Math.max(2, distance / 100);
    setFontSize(newFontSize);

    // Handle text position
    textRef.current.position.set(
      data.planets[name].x + data.planets[name].size + 10,
      data.planets[name].y + data.planets[name].size + 10,
      data.planets[name].z,
    );
    textRef.current.lookAt(camera.position);
  });

  const getGeometry = (key: string, x: number, z: number) => {
    if (planetGeometries[key]) {
      return planetGeometries[key];
    }

    const points = new THREE.EllipseCurve(
      0,
      0,
      x,
      z,
      0,
      2 * Math.PI,
      false,
    ).getPoints(1000);
    const geometry = new THREE.BufferGeometry().setFromPoints(points);

    setPlanetGeometries(
      (prevGeometries: THREE.BufferGeometry<THREE.NormalBufferAttributes>) => ({
        ...prevGeometries,
        [key]: geometry,
      }),
    );

    return geometry;
  };

  const text = (
    <mesh
      onPointerOver={() => setHovered(name)}
      onPointerOut={() => setHovered("")}
      onClick={() => {
        setFollowing(following === name ? "" : name);
        setZoomed(false);
      }}
    >
      <Text
        fillOpacity={opacity}
        ref={textRef}
        position={[data.planets[name].size + 10, data.planets[name].size, 0]}
        fontSize={fontSize}
        color="white"
      >
        {name.toUpperCase()}
      </Text>
    </mesh>
  );

  let planet;
  const material = new THREE.LineBasicMaterial({
    side: THREE.DoubleSide,
  });
  if (data.planets[name].satellites) {
    planet = (
      <group ref={planetRef}>
        <mesh
          onClick={() => {
            setFollowing(following === name ? "" : name);
            setZoomed(false);
          }}
          onPointerOver={() => setHovered(name)}
          onPointerOut={() => setHovered("")}
        >
          <sphereGeometry args={[data.planets[name].size, 32, 32]} />
          <meshPhongMaterial map={planetTexture} />
        </mesh>
        {Object.values(data.planets[name].satellites!).map((value, index) => {
          return <Satellite key={value + index} name={value} />;
        })}
        {Object.values(data.planets[name].satellites!).map((value, index) => {
          const hoveredMaterial = new THREE.LineBasicMaterial({
            color: data.planets[name].color,
            side: THREE.DoubleSide,
          });
          return (
            <line
              key={name + index}
              geometry={getGeometry(
                value,
                data.satellites[value].distance + data.planets[name].size,
                data.satellites[value].distance + data.planets[name].size,
              )}
              material={
                hovered === value || following === value
                  ? hoveredMaterial
                  : material
              }
              rotation-x={Math.PI / 2}
            />
          );
        })}
      </group>
    );
  } else {
    planet = (
      <mesh
        ref={planetRef}
        onClick={() => {
          setFollowing(following === name ? "" : name);
          setZoomed(false);
        }}
        onPointerOver={() => setHovered(name)}
        onPointerOut={() => setHovered("")}
      >
        <sphereGeometry args={[data.planets[name].size, 32, 32]} />
        <meshPhongMaterial map={planetTexture} />
      </mesh>
    );
  }

  const ringMesh = ring && (
    <mesh ref={ringRef} rotation-x={data.planets[name].ringRotation}>
      <torusGeometry
        args={[
          data.planets[name].ringDistance,
          data.planets[name].ringSize,
          2.5,
          100,
        ]}
      />
      <meshPhongMaterial map={ringTexture} />
    </mesh>
  );

  return (
    <>
      {planet}
      {text}
      {ringMesh}
    </>
  );
};

export default Planet;
