import { useFrame, useLoader, useThree } from "@react-three/fiber";
import { useContext, useEffect, useMemo, useRef, useState } from "react";
import context from "../utils/context.ts";
import * as THREE from "three";
// @ts-expect-error exist but typescript cause the error
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader";

const Rocket = ({ id }: { id: string }) => {
  // @ts-expect-error work but IDK why error here
  const obj = useMemo(() => {
    // eslint-disable-next-line react-hooks/rules-of-hooks
    return useLoader(GLTFLoader, "/assets/model/tardis/scene.gltf");
  });

  const clonedObj = obj.scene.clone();

  const rocketRef = useRef<THREE.Mesh>(null!);

  const { data } = useContext(context);
  const { hovered, setHovered } = useContext(context);
  const { camera } = useThree();
  const [fontSize, setFontSize] = useState(30);

  useEffect(() => {
    document.body.style.cursor = hovered !== "" ? "pointer" : "auto";
  }, [hovered]);

  useFrame(() => {
    // Handle orbit coordinates
    if (data.rockets[id]) {
      rocketRef.current.position.set(
        data.rockets[id].x,
        data.rockets[id].y,
        data.rockets[id].z,
      );

      if (data.planets[data.rockets[id].destination]) {
        const destCoord = new THREE.Vector3(
          data.planets[data.rockets[id].destination].x,
          data.planets[data.rockets[id].destination].y,
          data.planets[data.rockets[id].destination].z,
        );
        data.rockets[id].distanceTo =
          rocketRef.current.position.distanceTo(destCoord);
      }
    }

    const distance = rocketRef.current.position.distanceTo(camera.position);
    rocketRef.current.rotation.y += 0.01;
    // Handle text size
    const newFontSize = Math.max(30, distance / 150);

    // angle form by sun and planet saturn
    setFontSize(newFontSize);
  });
  return (
    <>
      <mesh
        ref={rocketRef}
        scale={fontSize}
        rotation-x={Math.PI / 12}
        onPointerOver={() => setHovered("rocket")}
        onPointerOut={() => setHovered("")}
      >
        <primitive object={clonedObj} />
      </mesh>
    </>
  );
};

export default Rocket;
