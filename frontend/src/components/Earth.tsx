import { Text, useTexture } from "@react-three/drei";
import { useFrame, useThree } from "@react-three/fiber";
import { useContext, useEffect, useRef, useState } from "react";
import * as THREE from "three";
// @ts-expect-error works but typescript cause the error
import fragmentShader from "../shaders/earth/earthFragment.glsl";
// @ts-expect-error works but typescript cause the error
import vertexShader from "../shaders/earth/earthVertex.glsl";
// @ts-expect-error works but typescript cause the error
import atmosphereFragment from "../shaders/earth/atmosphereFragment.glsl";
// @ts-expect-error works but typescript cause the error
import atmosphereVertex from "../shaders/earth/atmosphereVertex.glsl";
import { IPlanetProps } from "../utils/model.tsx";
import Satellite from "./Satellite.tsx";
import Context from "../utils/context.ts";

const Earth = ({ name }: IPlanetProps) => {
  const earthRef = useRef<THREE.Group>(null!);
  const { camera } = useThree();
  const textRef = useRef<THREE.Mesh>(null!);
  const [earthTexture] = useTexture(["/assets/texture/" + name + ".jpg"]);

  const [fontSize, setFontSize] = useState(14);
  const [opacity, setOpacity] = useState(1);
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

  useEffect(() => {
    document.body.style.cursor = hovered !== "" ? "pointer" : "auto";
  }, [hovered]);
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
  useFrame(() => {
    //Calculate the Earth's position based on its angle from the Sun
    earthRef.current.position.set(
      data.planets["earth"].x,
      0,
      data.planets["earth"].z,
    );

    earthRef.current.rotation.y += data.planets[name].revolution;

    // Handle camera follow
    const distance = earthRef.current.position.distanceTo(camera.position);

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
          earthRef.current.position.y + data.planets[name].size,
          z,
        );
        setZoomed(true);
      } else {
        camera.position.set(x, camera.position.y, z);
      }
      camera.lookAt(earthRef.current.position);
    }

    // Handle opacity
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
      data.planets["earth"].x + data.planets[name].size + 10,
      data.planets[name].size + 10,
      data.planets["earth"].z,
    );
    textRef.current.lookAt(camera.position);
  });

  const sphere = {
    globeTexture: {
      value: earthTexture,
    },
  };

  const hoveredMaterial = new THREE.LineBasicMaterial({
    color: data.planets[name].color,
    side: THREE.DoubleSide,
  });
  const material = new THREE.LineBasicMaterial({
    side: THREE.DoubleSide,
  });

  return (
    <>
      <group ref={earthRef} position={[data.planets[name].size, 0, 0]}>
        <mesh
          castShadow={true}
          receiveShadow={true}
          onClick={() => {
            setFollowing(following !== name ? name : "");
            setZoomed(false);
          }}
          onPointerOver={() => setHovered(name)}
          onPointerOut={() => setHovered("")}
        >
          <sphereGeometry args={[data.planets[name].size, 32, 32]} />
          <shaderMaterial
            vertexShader={vertexShader}
            fragmentShader={fragmentShader}
            uniforms={sphere}
          />
        </mesh>
        <mesh scale={[1.1, 1.1, 1.1]}>
          <sphereGeometry args={[data.planets[name].size, 32, 32]} />
          <shaderMaterial
            vertexShader={atmosphereVertex}
            fragmentShader={atmosphereFragment}
            side={THREE.BackSide}
            blending={THREE.AdditiveBlending}
          />
        </mesh>
        <Satellite name={"moon"} />
        <line
          geometry={getGeometry(
            "moon",
            data.satellites["moon"].distance + data.planets[name].size,
            data.satellites["moon"].distance + data.planets[name].size,
          )}
          material={
            hovered === "moon" || following === "moon"
              ? hoveredMaterial
              : material
          }
          rotation-x={Math.PI / 2}
        />
      </group>
      <mesh
        onPointerOver={() => setHovered(name)}
        onPointerOut={() => setHovered("")}
        onClick={() => {
          setFollowing(following !== name ? name : "");
          setZoomed(false);
        }}
      >
        <Text
          fillOpacity={opacity}
          ref={textRef}
          position={[7.8341 + 10, 7.8341, 0]}
          fontSize={fontSize}
          color="white"
        >
          EARTH
        </Text>
      </mesh>
    </>
  );
};

export default Earth;
