'use client'

import { getPrice } from "@/lib/okx/client";
import { motion, AnimatePresence } from "framer-motion";
import { CheckCircle, Loader2 } from "lucide-react";
import { useEffect, useState } from "react";
import { formatAmountToUSD } from "../sections/swap-section";

interface UserOperationReceipt {
  actualGasCost: string;
  actualGasUsed: string;
  nonce: string;
  paymaster: string;
  sender: string;
  success: boolean;
  userOpHash: string;
  receipt: {
    blockHash: string;
    blockNumber: string;
    cumulativeGasUsed: string;
    effectiveGasPrice: string;
    gasUsed: string;
    logs: string[];
    logsBloom: string;
    transactionHash: string;
  };
}

const steps = [
  "Waiting in Bundlr",
  "Waiting for Block Confirmation"
];

function formatGas(hex: string) {
  try {
    return `${BigInt(hex).toLocaleString()} gas`;
  } catch {
    return hex;
  }
}

export default function ProcessSteps({
  userOpHash,
  onBack
}: {
  userOpHash: `0x${string}`;
  onBack: () => void;
}) {
  const [currentStep, setCurrentStep] = useState(0);
  const [txConfirmed, setTxConfirmed] = useState(false);
  const [receipt, setReceipt] = useState<UserOperationReceipt | null>(null);
  const [actualGasCostUSD, setActualGasCostUSD] = useState<string | null>(null);

  // Düz JSON-RPC polling
  useEffect(() => {
    if (!userOpHash || userOpHash === "0x") return;

    const interval = setInterval(async () => {
      try {
        const res = await fetch("http://127.0.0.1:8181/rpc/getUserOpReceipt", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            jsonrpc: "2.0",
            method: "eth_getUserOperationReceipt",
            params: [userOpHash],
            id: 1,
          }),
        });

        const json = await res.json();
        const result = json?.result;
        if (!result) return;

        const state = result.state;
        const rcp = result.receipt;

        if (state === "pending") {
          setCurrentStep(0);
        } else if (state === "bundled") {
          setCurrentStep(1);
        } else if (state === "sent") {
          setCurrentStep(1);

          let price = "0";

          try {
            price = await getPrice("0xe538905cf8410324e03a5a23c1c177a474d59b2b", 18);
          } catch (error) {
            console.error("Error fetching price:", error);
          }
          
          const actualGasCost = formatAmountToUSD(BigInt(rcp.actualGasCost), 18, Number(price));
          setActualGasCostUSD(actualGasCost);

          setReceipt(rcp);
          setTxConfirmed(true);
          clearInterval(interval);
        }
      } catch (err) {
        console.error("Polling failed:", err);
      }
    }, 500);

    return () => clearInterval(interval);
  }, [userOpHash]);

  return (
    <div className="w-full flex flex-col items-center text-center">
      <h2 className="text-xl font-semibold text-[#292929] mb-6">Processing Transaction</h2>
      <div className="flex flex-col gap-4 w-full mb-6">
        {steps.map((step, index) => {
          const isConfirmed = txConfirmed && index === 1;
          const isActive = index === currentStep;

          return (
            <motion.div
              key={step}
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.4, delay: index * 0.3 }}
              className={`flex flex-col px-4 py-3 rounded-xl shadow transition-all ${
                index < currentStep || isConfirmed
                  ? "bg-green-50 text-green-700"
                  : isActive
                  ? "bg-blue-50 text-blue-700"
                  : "bg-gray-100 text-gray-500"
              }`}
            >
              <div className="flex items-center gap-3">
                {index < currentStep || isConfirmed ? (
                  <CheckCircle size={20} className="text-green-500" />
                ) : isActive ? (
                  <Loader2 size={20} className="animate-spin text-blue-500" />
                ) : (
                  <div className="w-[20px] h-[20px] rounded-full border border-gray-400" />
                )}
                <span className="text-sm font-medium">{step}</span>
              </div>

              {index === 1 && receipt && (
                <AnimatePresence>
                  <motion.div
                    initial={{ opacity: 0, height: 0 }}
                    animate={{ opacity: 1, height: "auto" }}
                    exit={{ opacity: 0, height: 0 }}
                    transition={{ duration: 0.4 }}
                    className="mt-3 w-full text-left text-sm text-[#444] bg-white border border-[#ddd] rounded-xl p-4 overflow-hidden"
                  >
                    <div className="mb-3 p-2 rounded-md bg-green-50 border border-green-300 text-green-800 flex items-center gap-2">
                      <CheckCircle size={16} className="text-green-600" />
                      <span>Transaction Confirmed</span>
                    </div>

                    <div className="mb-2">
                      <span className="font-medium">Tx Hash:</span>{" "}
                      <a
                        href={`https://www.oklink.com/x-layer/tx/${receipt.receipt.transactionHash}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-[#9257f3] underline break-all"
                      >
                        {receipt.receipt.transactionHash}
                      </a>
                    </div>

                    <div className="mb-2">
                      <span className="font-medium">Success:</span>{" "}
                      <span className={receipt.success ? "text-green-600" : "text-red-600"}>
                        {receipt.success ? "true" : "false"}
                      </span>
                    </div>

                    <div className="mb-2">
                      <span className="font-medium">Gas Used:</span>{" "}
                      <span className="text-[#666]">
                        {formatGas(receipt.actualGasUsed)}
                      </span>
                    </div>

                    <div className="mb-2">
                      <span className="font-medium">Gas Cost:</span>{" "}
                      <span className="text-[#666]">
                        {actualGasCostUSD} $
                      </span>
                    </div>
                  </motion.div>
                </AnimatePresence>
              )}
            </motion.div>
          );
        })}
      </div>

      {txConfirmed && receipt && (
        <button
          onClick={onBack}
          className="mt-6 px-5 py-2.5 text-sm font-semibold rounded-lg bg-[#9257f3] text-white hover:bg-[#7b3df2] cursor-pointer shadow transition-all"
        >
          Back to Swap Page
        </button>
      )}
    </div>
  );
}