import { useWeb3 } from "@/context/web3Context";
import { WalletProps } from "@/types/wallet";
import { Minus, Plus, RefreshCw } from "lucide-react"
import { useState } from "react";

export default function HistorySection({ user }: { user: WalletProps | null }) {
  const { txHistory, fetchTxHistory } = useWeb3();
  const [loading, setLoading] = useState(false);
  
  const refreshHistory = async () => {
    try {
      setLoading(true);
      if (user?.accountAddress) {
        await fetchTxHistory(user.accountAddress as `0x${string}`);
      }
    } catch (error) {
      console.error("Error fetching TxHistory:", error);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="mt-8 w-full h-[600px] overflow-y-scroll scroll-smooth">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-xl font-bold text-gray-800">Transaction History</h2>
        <button
          onClick={refreshHistory}
          className={`absolute left-auto right-12 rounded-full font-medium transition cursor-pointer ${ loading ? 'text-[#965cef]' : 'text-[#b8b8b8] hover:text-[#965cef]'}`}
          disabled={loading}
        >
          <RefreshCw size={19} className={`${loading ? "animate-spin" : "block"} `} />
        </button>
      </div>

      <div className="flex flex-col gap-3 mt-4">
        {[...txHistory].reverse().map((tx, index) => (
          <div key={index} className="flex gap-1 bg-white shadow-md rounded-lg px-4 h-[90px] w-[340px]">
            <div className="flex items-center gap-2">
              <img
                className="h-8 w-8 rounded-full"
                src="https://altcoinsbox.com/wp-content/uploads/2023/03/okx-logo-black-and-white.jpg"
                alt="dex"
              />
              <div className="flex flex-col gap-[1px]">
                <span className="text-[15px] font-medium">{tx.dexName}</span>
                <span className="text-[13px] text-[#555]">
                  {new Date(tx.timestamp).toLocaleTimeString()} • {new Date(tx.timestamp).toLocaleDateString()}
                </span>
              </div>
            </div>

            {/* ORTA */}
            <div className="flex flex-col items-end gap-[6px] ml-auto mr-0 my-auto">
              {/* FROM */}
              <div className="text-sm font-semibold flex gap-[6px] items-center text-[14px] text-[#d86060]">
                <Minus className="inline-block" size={12} strokeWidth={2.5} />
                <p>{tx.fromAmount}</p>
                <img
                  className="rounded-full w-6 h-6"
                  src={`${(tx.fromToken)}.png`}
                  alt=""
                />
              </div>

              {/* TO */}
              <div className="text-sm font-semibold flex gap-[8px] items-center text-[14px] text-[#60d870]">
                <Plus className="inline-block" size={12} strokeWidth={2.5} />
                <p>{tx.toAmount}</p>
                <img
                  className="rounded-full w-6 h-6"
                  src={`${(tx.toToken)}.png`}
                  alt=""
                />
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}