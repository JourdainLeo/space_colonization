import { useThree } from "@react-three/fiber";
import { CubeTexture, CubeTextureLoader } from "three";
import { useState } from "react";

export function SkyBox() {
  const { scene } = useThree();
  const loader = new CubeTextureLoader();
  const [texture, setTexture] = useState<CubeTexture>();
  // The CubeTextureLoader load method takes an array of urls representing all 6 sides of the cube.
  if (!texture) {
    setTexture(
      loader.load([
        "/assets/skybox/sky6ft.png",
        "/assets/skybox/sky6bk.png",
        "/assets/skybox/sky6up.png",
        "/assets/skybox/sky6dn.png",
        "/assets/skybox/sky6rt.png",
        "/assets/skybox/sky6lf.png",
      ]),
    );
  }

  // Set the scene background property to the resulting texture.
  if (texture) scene.background = texture;
  return null;
}
