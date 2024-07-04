import { useContext, useEffect, useRef, useState } from "react";
import { Text, useTexture } from "@react-three/drei";
import { useFrame, useThree } from "@react-three/fiber";
import * as THREE from "three";
import Context from "../utils/context.ts";

const Satellite = ({ name }: { name: string }) => {
  const satelliteRef = useRef<THREE.Mesh>(null!);
  const textRef = useRef<THREE.Mesh>(null!);
  const [satelliteTexture] = useTexture(["/assets/texture/" + name + ".jpg"]);
  const { camera } = useThree();
  const [zoomed, setZoomed] = useState(false);

  const { following, setFollowing, hovered, setHovered, data } =
    useContext(Context);
  useFrame(() => {
    // Handle orbit coordinates
    satelliteRef.current.position.set(
      data.satellites[name].x,
      0,
      data.satellites[name].z,
    );

    // Handle revolution
    satelliteRef.current.rotation.y += data.satellites[name].revolution;

    // Handle camera follow
    const name2 = data.satellites[name].parent;
    if (following === name) {
      const x =
        ((data.planets[name2].maxDistance - data.planets[name2].size * 4) /
          data.planets[name2].maxDistance) *
        data.planets[name2].x;
      const z =
        ((data.planets[name2].minDistance - data.planets[name2].size * 4) /
          data.planets[name2].minDistance) *
        data.planets[name2].z;

      if (!zoomed) {
        camera.position.set(
          x,
          data.planets[name2].y + data.planets[name2].size,
          z,
        );
        setZoomed(true);
      } else {
        camera.position.set(x, camera.position.y, z);
      }

      const newVector = new THREE.Vector3(
        data.planets[name2].x,
        data.planets[name2].y,
        data.planets[name2].z,
      );
      camera.lookAt(newVector);
    }

    // Handle text position
    textRef.current.position.set(
      satelliteRef.current.position.x + data.satellites[name].size + 10,
      satelliteRef.current.position.y + data.satellites[name].size + 10,
      satelliteRef.current.position.z,
    );
    textRef.current.lookAt(camera.position);
  });

  useEffect(() => {
    document.body.style.cursor = hovered !== "" ? "pointer" : "auto";
  }, [hovered]);

  return (
    <>
      <mesh
        ref={satelliteRef}
        position={[data.satellites[name].distance, 0, 0]}
        onPointerOver={() => setHovered(name)}
        onPointerOut={() => setHovered("")}
        onClick={() => {
          if (following === name) {
            setFollowing("");
          } else {
            setFollowing(name);
          }
        }}
      >
        <sphereGeometry args={[data.satellites[name].size, 32, 32]} />
        <meshPhongMaterial map={satelliteTexture} />
      </mesh>
      <mesh
        onPointerOver={() => setHovered(name)}
        onPointerOut={() => setHovered("")}
        onClick={() => {
          if (following === name) {
            setFollowing("");
          } else {
            setFollowing(name);
            setZoomed(false);
          }
        }}
      >
        <Text
          fillOpacity={
            following === data.satellites[name].parent || name !== following
              ? 0.7
              : 0
          }
          ref={textRef}
          position={[
            data.satellites[name].size,
            data.satellites[name].size / 2,
            0,
          ]}
          fontSize={4}
          color="white"
        >
          {name.charAt(0).toUpperCase() + name.slice(1)}
        </Text>
      </mesh>
    </>
  );
};

export default Satellite;
