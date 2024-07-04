const DataContainer = ({
  icon,
  text,
  hover,
  func,
}: {
  icon: any;
  text: string;
  hover?: boolean;
  func?: any;
}) => {
  return (
    <div
      className={"icon-container " + (hover ? "icon-hover" : "")}
      onClick={() => {
        if (func) func();
      }}
    >
      {icon}
      {text}
    </div>
  );
};

export default DataContainer;
