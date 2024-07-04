import Context from "../../utils/context.ts";
import { useContext } from "react";

const CardBar = ({
  name,
  close,
  style,
}: {
  name: string;
  close?: boolean;
  style?: string;
}) => {
  const { setFollowing } = useContext(Context);
  return (
    <div className={"card-bar " + style}>
      <div className="card-title">{name.toUpperCase()}</div>

      {close && (
        <div
          className="button-close"
          onClick={() => {
            setFollowing("");
          }}
        >
          X
        </div>
      )}
    </div>
  );
};

export default CardBar;
