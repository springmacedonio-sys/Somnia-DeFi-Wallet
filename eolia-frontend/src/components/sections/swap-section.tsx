// SwapSection.tsx
'use client'

import { useEffect, useRef, useState } from "react";
import { ArrowDownUp, RefreshCcw } from "lucide-react";
import TokenSelector from "../ui/token-selector";
import { calculateGas, getPrice, getSwapQuote, Swap } from "@/lib/okx/client";
import ConfirmModal from "../ui/confirm-modal";
import ProcessSteps from "../ui/process-steps";
import { tokens } from "@/constants/tokens";
import { WalletProps } from "@/types/wallet";
import { trimToDecimals } from "@/utils/decimals";
import { formatAmount, formatAmountToUSD, parseAmount } from "@/utils/formatting";

const SwapSection = ({ user }: { user: WalletProps | null }) => {
  const [fromToken, setFromToken] = useState(tokens[0]);
  const [toToken, setToToken] = useState(tokens[1]);
  const [fromAmount, setFromAmount] = useState("");
  const [toAmount, setToAmount] = useState("");
  const [quote, setQuote] = useState<string | null>(null);
  const [quoteTimer, setQuoteTimer] = useState(15);
  const [isFetchingQuote, setIsFetchingQuote] = useState(false);
  const [rawGasFee, setRawGasFee] = useState<string | null>(null);
  const [gasFee, setGasFee] = useState<string | null>(null);
  const [stage, setStage] = useState<"idle" | "gas-calculated">("idle");
  const [showModal, setShowModal] = useState(false);
  const [showProcess, setShowProcess] = useState(false);
  const [userHash, setUserHash] = useState<`0x${string}`>("0x");
  const [events, setEvents] = useState<{ label: string; line: string }[]>([]);

  const quoteIntervalRef = useRef<NodeJS.Timeout | null>(null);
  const countdownIntervalRef = useRef<NodeJS.Timeout | null>(null);
  const lastGasTimestamp = useRef<number | null>(null);

  const disabled =
    fromToken.symbol === toToken.symbol ||
    !quote ||
    quote === "Insufficient liquidity" ||
    isFetchingQuote

  const buttonDisabled =
    fromToken.symbol === toToken.symbol ||
    !quote ||
    quote === "Insufficient liquidity" ||
    isFetchingQuote ||
    !fromAmount || 
    Number(fromAmount) <= 0;

  const handleSwapTokens = () => {
    if (fromToken.symbol === toToken.symbol) return;
    const temp = fromToken;
    setFromToken(toToken);
    setToToken(temp);
  };

  const fetchQuote = async () => {
    setIsFetchingQuote(true);
    console.log("Fetching quote...");
    try {
      const amountToSend = parseAmount("1", fromToken.decimals);
      const quoteResult = await getSwapQuote(
        fromToken.address,
        toToken.address,
        amountToSend.toString(),
        "0.005"
      );
      const humanReadable = formatAmount(quoteResult.toTokenAmount, toToken.decimals);
      setQuote(trimToDecimals(Number(humanReadable), 6));
    } catch {
      setQuote(null);
    } finally {
      setIsFetchingQuote(false);
      setQuoteTimer(15);
    }
  };

  const calculateTo = (_fromAmount: string) => {
    if (Number(_fromAmount) > 0 && quote && !isNaN(Number(quote))) {
      const result = Number(_fromAmount) * Number(quote);
      setToAmount(result.toFixed(6));
    } else {
      setToAmount("");
    }
  };

  const calculateGasFees = async () => {
    setIsFetchingQuote(true);
    setStage("idle");

    try {
      const { totalGas, events } = await calculateGas({
        sender: user?.accountAddress as `0x${string}`,
        fromToken: fromToken.address as `0x${string}`,
        toToken: toToken.address as `0x${string}`,
        amount: parseAmount(fromAmount, fromToken.decimals).toString(),
        slippage: "0.005"
      })
      const OKB_PRICE = await getPrice("0xe538905cf8410324e03a5a23c1c177a474d59b2b", 18);
      if (!totalGas || !OKB_PRICE || OKB_PRICE === "") {
        console.error("Failed to fetch gas fee or OKB price");
        return;
      }

      const formattedGasFee = formatAmountToUSD(totalGas, 18, Number(OKB_PRICE));
      setRawGasFee(formatAmount(totalGas, 18));
      setGasFee(formattedGasFee);
      setEvents(events);
      lastGasTimestamp.current = Date.now();

      setStage("gas-calculated");
    } catch (error) {
      console.error("Error calculating gas fees:", error);
      setGasFee(null);
      setEvents([]);
    } finally {
      setIsFetchingQuote(false);
    }
  };

  const handleConfirmSwap = async () => {
    setUserHash("0x");
    setShowModal(false);
    setShowProcess(true);

    const userOpHash = await Swap({
      sender: user?.accountAddress as `0x${string}`,
      fromToken: fromToken.address as `0x${string}`,
      toToken: toToken.address as `0x${string}`,
      amount: parseAmount(fromAmount, fromToken.decimals).toString(),
      slippage: "0.005"
    })
    if (!userOpHash) {
      console.error("Failed to get user operation hash");
      return;
    }
    setUserHash(userOpHash);
  };

  useEffect(() => {
    clearInterval(quoteIntervalRef.current!);
    clearInterval(countdownIntervalRef.current!);

    setQuoteTimer(15);
    fetchQuote();
    setFromAmount("");
    setToAmount("");
    setGasFee(null);
    setStage("idle");

    quoteIntervalRef.current = setInterval(fetchQuote, 15000);
    countdownIntervalRef.current = setInterval(() => {
      setQuoteTimer((prev) => {
        const next = prev === 0 ? 15 : prev - 1;

        if (lastGasTimestamp.current && Date.now() - lastGasTimestamp.current >= 15000) {
          setGasFee(null);
          setStage("idle");
          lastGasTimestamp.current = null;
        }

        return next;
      });
    }, 1000);

    return () => {
      clearInterval(quoteIntervalRef.current!);
      clearInterval(countdownIntervalRef.current!);
    };
  }, [fromToken, toToken]);

  useEffect(() => {
    // Clear gas fee if input or token changes
    setGasFee(null);
    setStage("idle");
    lastGasTimestamp.current = null;
  }, [fromAmount, fromToken, toToken]);

  return (
    <div className="w-[400px] h-[640px] relative flex flex-col items-center rounded-t-xl bg-[#f7f7f7] px-[40px] overflow-y-scroll pt-8">
      <ConfirmModal
        isOpen={showModal}
        onClose={() => setShowModal(false)}
        onConfirm={handleConfirmSwap}
        data={{
          title: "Transaction",
          description: "You're just one step away from completing your transaction.",
          actions: events,
          confirmLabel: "Confirm Swap"
        }}
      />
      {showProcess ? (
        <ProcessSteps userOpHash={userHash} onBack={() => setShowProcess(false)} />
      ) : (
        <>
          <div className="text-center mb-6">
            <h2 className="text-2xl font-semibold text-[#292929]">Swap Tokens</h2>
            <p className="text-sm text-[#7d7d7d] mt-1">Exchange your tokens instantly</p>
          </div>

          {/* Input Section */}
          <div className="w-full space-y-4 relative">
            {/* From */}
            <div className="bg-white px-3 py-3 rounded-xl shadow">
              <p className="text-xs text-[#7d7d7d] mb-1">From</p>
              <div className="flex items-center gap-2">
                <input
                  type="number"
                  value={fromAmount}
                  onChange={(e) => {
                    setFromAmount(e.target.value);
                    calculateTo(e.target.value);
                  }}
                  placeholder="0.0"
                  className="text-xl font-semibold text-[#292929] outline-none bg-transparent w-full border-[1px] border-[#f5f5f5] rounded-lg px-3 py-2"
                  readOnly={disabled}
                />
                <TokenSelector
                  selected={fromToken}
                  onChange={setFromToken}
                  excludeSymbol={toToken.symbol}
                  disabled={disabled}
                />
              </div>
            </div>

            {/* Flip Button - Positioned */}
            <div
              onClick={handleSwapTokens}
              className="bg-[#9257f3] py-[6px] flex w-[100px] mx-auto items-center justify-center rounded-full text-white cursor-pointer hover:scale-105 transition-all shadow-lg"
            >
              <ArrowDownUp size={20} />
            </div>

            {/* To */}
            <div className="bg-white px-4 py-3 rounded-xl shadow">
              <p className="text-xs text-[#7d7d7d] mb-1">To</p>
              <div className="flex items-center gap-2">
                <input
                  type="number"
                  value={toAmount}
                  onChange={() => {}}
                  placeholder="0.0"
                  className="text-xl font-semibold text-[#292929] outline-none bg-transparent w-full"
                  readOnly
                />
                <TokenSelector
                  selected={toToken}
                  onChange={setToToken}
                  excludeSymbol={fromToken.symbol}
                  disabled={disabled}
                />
              </div>
            </div>
          </div>

          {/* Quote Info */}
          <div className="text-sm text-[#7d7d7d] mt-1 mb-5 flex items-center gap-2">
            <RefreshCcw
              size={14}
              onClick={isFetchingQuote ? undefined : fetchQuote}
              className={`cursor-pointer transition ${isFetchingQuote ? "opacity-50" : "hover:text-[#965cef] active:rotate-90"}`}
            />
            {quote ? (
              <span>1 {fromToken.symbol} ≈ {quote} {toToken.symbol}</span>
            ) : (
              <span className="text-red-500">No quote available</span>
            )}
            <span className="text-xs text-gray-500 ml-1">({quoteTimer}s)</span>
          </div>

          {/* Gas Fee Display */}
          {stage === "gas-calculated" && gasFee && (
            <div className="bg-white p-4 rounded-xl w-full shadow text-sm text-[#292929] mb-4 text-center">
              <p className="font-medium mb-1">Estimated Gas Fee</p>
              <p className="text-base font-semibold">{gasFee} $ <span className="text-xs text-gray-500">({rawGasFee} OKB)</span></p>
              <p className="text-xs text-gray-500 mt-1">This is required to process your transaction on-chain.</p>
            </div>
          )}

          {/* Main Action Button */}
          <button
            disabled={buttonDisabled}
            onClick={stage === "idle" ? calculateGasFees : () => setShowModal(true)}
            className={`w-full ${buttonDisabled ? "bg-gray-300 cursor-not-allowed" : "bg-[#965cef] hover:bg-[#8a52fc] cursor-pointer"} text-white text-sm font-semibold py-3 rounded-full transition`}
          >
            {stage === "idle" ? "Calculate Gas Fees" : isFetchingQuote ? "Processing..." : "Swap"}
          </button>
        </>
      )}
    </div>
  );
};

export default SwapSection;
export { formatAmountToUSD };