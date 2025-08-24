const Toggle = ({ isOn, onToggle }: { isOn: boolean; onToggle: () => void }) => {
  return (
    <div
      onClick={onToggle}
      className={`w-11 h-6 flex items-center rounded-full p-[2px] cursor-pointer transition-colors ${
        isOn ? "bg-purple-500" : "bg-gray-300"
      }`}
    >
      <div
        className={`w-5 h-5 bg-white rounded-full shadow-md transform duration-300 ${
          isOn ? "translate-x-[20px]" : "translate-x-0"
        }`}
      />
    </div>
  );
};

export default Toggle;
