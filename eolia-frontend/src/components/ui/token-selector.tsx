// components/extra/TokenSelector.tsx
import { useState } from "react";
import { ChevronDown } from "lucide-react";
import { tokens } from "@/constants/tokens";

interface TokenSelectorProps {
  selected: typeof tokens[0];
  onChange: (token: typeof tokens[0]) => void;
  excludeSymbol?: string;
  disabled?: boolean;
}

const TokenSelector = ({ selected, onChange, excludeSymbol, disabled }: TokenSelectorProps) => {
  const [open, setOpen] = useState(false);

  const filteredTokens = tokens.filter(
    (t) => t.symbol !== excludeSymbol
  );

  return (
    <div className="relative">
      <button
        onClick={() => !disabled && setOpen(!open)}
        className={`flex items-center gap-2 border border-[#ccc] px-2 py-1 rounded-lg text-sm bg-white hover:shadow-sm transition-all ${disabled ? "opacity-40 cursor-not-allowed" : "cursor-pointer"}`}
      >
        <img src={selected.logo} alt={selected.symbol} className="w-5 h-5 rounded-full" />
        <span>{selected.symbol}</span>
        <ChevronDown size={14} />
      </button>

      {open && (
        <div className="absolute top-full mt-2 right-0 bg-white shadow-xl rounded-xl z-30 w-[180px] max-h-[250px] overflow-y-auto border border-[#ddd] animate-fadeIn">
          {filteredTokens.map((token) => (
            <div
              key={token.symbol}
              onClick={() => {
                onChange(token);
                setOpen(false);
              }}
              className="flex items-center gap-2 px-4 py-2 hover:bg-[#f3f3f3] cursor-pointer text-sm"
            >
              <img src={token.logo} alt={token.symbol} className="w-4 h-4 rounded-full" />
              <span className="font-medium">{token.symbol}</span>
              <span className="text-xs text-gray-500 ml-auto">{token.name}</span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default TokenSelector;