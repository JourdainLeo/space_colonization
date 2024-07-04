import { Text, useTexture } from "@react-three/drei";
import { useFrame, useThree } from "@react-three/fiber";
import { useContext, useEffect, useRef, useState } from "react";
import * as THREE from "three";
// @ts-expect-error works but typescript cause the error
import fragmentShader from "../shaders/sun/sunFragment.glsl";
// @ts-expect-error works but typescript cause the error
import vertexShader from "../shaders/sun/sunVertex.glsl";
// @ts-expect-error works but typescript cause the error
import atmosphereFragment from "../shaders/sun/solarFragment.glsl";
// @ts-expect-error works but typescript cause the error
import atmosphereVertex from "../shaders/sun/solarVertex.glsl";

import { solar_system } from "../utils/base_data.tsx";
import Context from "../utils/context.ts";

const material = new THREE.LineBasicMaterial({
  side: THREE.DoubleSide,
});

const Sun = () => {
  const sunRef = useRef<THREE.Mesh>(null!);
  const [sunTexture] = useTexture(["/assets/texture/sun.jpg"]);
  const textRef = useRef<THREE.Mesh>(null!);
  const { planetGeometries, setPlanetGeometries } = useContext(Context);
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

  const { following, setFollowing, hovered, setHovered } = useContext(Context);

  useFrame(() => {
    sunRef.current.rotation.y -= 0.003;
    const distance = sunRef.current.position.distanceTo(camera.position);
    const newFontSize = Math.max(14, distance / 100);
    setFontSize(newFontSize);
    textRef.current.lookAt(camera.position);
  });

  const sphere = {
    globeTexture: {
      value: sunTexture,
    },
  };

  useEffect(() => {
    document.body.style.cursor = hovered !== "" ? "pointer" : "auto";
  }, [hovered]);

  const [fontSize, setFontSize] = useState(14);

  const { camera } = useThree();
  return (
    <>
      <mesh
        ref={sunRef}
        position={[0, 0, 0]}
        onClick={() => {
          camera.lookAt(0, 0, 0);
          setFollowing("");
        }}
        onPointerOver={() => setHovered("sun")}
        onPointerOut={() => setHovered("")}
      >
        <sphereGeometry args={[100, 32, 32]} />
        <pointLight castShadow={true} />
        <shaderMaterial
          vertexShader={vertexShader}
          fragmentShader={fragmentShader}
          uniforms={sphere}
        />
        <Text
          ref={textRef}
          position={[10, 85, 0]}
          fontSize={fontSize}
          color="white"
        >
          SUN
        </Text>
      </mesh>
      <mesh scale={[1.1, 1.1, 1.1]}>
        <sphereGeometry args={[100, 32, 32]} />
        <shaderMaterial
          vertexShader={atmosphereVertex}
          fragmentShader={atmosphereFragment}
          side={THREE.BackSide}
          blending={THREE.AdditiveBlending}
        />
      </mesh>

      {Object.entries(solar_system.planets).map(([key, value], index) => {
        const hoveredMaterial = new THREE.LineBasicMaterial({
          color: solar_system.planets[key].color,
          side: THREE.DoubleSide,
        });
        return (
          <line
            key={key + index}
            geometry={getGeometry(
              key,
              value.maxDistance + 100,
              value.minDistance + 100,
            )}
            material={
              hovered === key || following === key ? hoveredMaterial : material
            }
            rotation-x={Math.PI / 2}
          />
        );
      })}
    </>
  );
};

export default Sun;
