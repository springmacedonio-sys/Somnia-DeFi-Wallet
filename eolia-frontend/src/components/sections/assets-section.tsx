"use client"

import { useState } from "react"
import { useWeb3 } from "@/context/web3Context";
import { tokens } from "@/constants/tokens";
import { LoaderCircle, RefreshCw } from "lucide-react";
import { WalletProps } from "@/types/wallet";
import { trimToDecimals } from "@/utils/decimals";

const dummyCollectibles = [
  { name: "Doodle Devils", image: "https://lh3.googleusercontent.com/TjiVAvF_PH5FdVYO7bxmEtH4bElf4RkhEz10xBo6MP4r_9sGpmMpGVDUQln4CVGhEEpEX40kQP-u5ddlxNmmy3yeDIbs3gGXxoE" },
  { name: "Pudgy Penguins", image: "https://img-cdn.magiceden.dev/autoquality:size:512000:20:80/rs:fill:250:0:0/plain/https%3A%2F%2Fbafybeia7qfyvgwfxkhwjoorxqe6nezjk7jqv3bblo4cbs4hpjcecnky4vy.ipfs.w3s.link%2FIMG_6833.png" },
  { name: "Bored Ape Yacht Club", image: "https://img-cdn.magiceden.dev/autoquality:size:512000:20:80/rs:fill:250:0:0/plain/https%3A%2F%2Fbafybeicj5occh26bnn6eyvysll7sta2idw2t6b6xb6gbhgxuqcyf4dsbg4.ipfs.w3s.link%2FIMG_8872.png" },
]

const AssetsSection = ({ user }: { user: WalletProps | null }) => {
  const [activeTab, setActiveTab] = useState<"tokens" | "collectibles">("tokens")
  const { balances, fetchBalances } = useWeb3();

  const [loading, setLoading] = useState(false);

  const refreshBalances = async () => {
    try {
      setLoading(true);
      if (user?.accountAddress) {
        await fetchBalances(user.accountAddress as `0x${string}`);
      }
    } catch (error) {
      console.error("Error fetching balances:", error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="w-[400px] h-[640px] relative flex flex-col items-center rounded-t-xl bg-[#f7f7f7] px-[20px] py-5 overflow-y-scroll scroll-smooth">
      
      {/* Tab Selector */}
      <div className="flex justify-center gap-2 bg-white rounded-full px-2 py-1 mb-5 shadow-md">
        <button
          onClick={() => setActiveTab("tokens")}
          className={`text-sm px-4 py-1 rounded-full font-medium transition cursor-pointer ${
            activeTab === "tokens"
              ? "bg-[#965cef] text-white"
              : "text-[#555] hover:text-black"
          }`}
        >
          Tokens
        </button>
        <button
          onClick={() => setActiveTab("collectibles")}
          className={`text-sm px-4 py-1 rounded-full font-medium transition cursor-pointer ${
            activeTab === "collectibles"
              ? "bg-[#965cef] text-white"
              : "text-[#555] hover:text-black"
          }`}
        >
          Collectibles
        </button>
        <button
          onClick={refreshBalances}
          className={`absolute left-auto right-12 py-1 rounded-full font-medium transition cursor-pointer ${ loading ? 'text-[#965cef]' : 'text-[#b8b8b8] hover:text-[#965cef]'}`}
          disabled={loading}
        >
          <RefreshCw size={20} className={`${loading ? "animate-spin" : "block"} `} />
        </button>
      </div>

      {/* Token List */}
      {activeTab === "tokens" && (
        <div className="w-full flex flex-col gap-2">
          {tokens.map((token, idx) => (
            <div key={idx} className="bg-white rounded-lg px-4 py-3 flex items-center justify-between shadow-sm">
              <div className="flex items-center gap-3">
                <img src={token.logo} className="w-7 h-7 rounded-full" alt={token.name} />
                <div className="flex flex-col">
                  <span className="text-sm font-semibold">{token.name}</span>
                  <span className="text-xs text-[#666]">{token.symbol}</span>
                </div>
              </div>
              <div className="flex flex-col">
                <span className="text-sm font-semibold flex ml-auto mr-0">{balances[idx] ? trimToDecimals(balances[idx].balance, 5) : <LoaderCircle size={18} className="animate-spin" />} {token.symbol}</span>
                <span className="text-xs text-[#666] flex w-max ml-auto mr-0">{balances[idx] ? trimToDecimals(balances[idx].usdValue, 2) : <LoaderCircle size={18} className="animate-spin" />} $</span>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Collectibles List */}
      {activeTab === "collectibles" && (
        <div className="w-full grid grid-cols-3 gap-3 mt-1">
          {dummyCollectibles.map((nft, idx) => (
            <div key={idx} className="bg-white rounded-lg shadow-sm p-1 flex flex-col items-center">
              <img src={nft.image} alt={nft.name} className="w-full h-[80px] object-cover rounded-md" />
              <p className="text-[11px] mt-1 text-center">{nft.name}</p>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default AssetsSection