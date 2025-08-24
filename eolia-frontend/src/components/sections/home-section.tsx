import { ArrowRightLeft, Copy, DollarSign, Minus, Plus, QrCode, ScanLine, Send, LogOut, LoaderCircle } from "lucide-react"
import TransferModal from '../TransferModal';
import { useState } from 'react';
import { useWeb3 } from "@/context/web3Context";
import { WalletProps } from "@/types/wallet";

const HomeSection = ({ user, logOut }: { user: WalletProps | null, logOut: () => Promise<void> }) => {
  const [modalOpen, setModalOpen] = useState(false);
  const { totalBalance, txHistory } = useWeb3();
  return (
    <div className="w-[400px] h-[640px] relative flex flex-col items-center rounded-t-xl bg-[#f7f7f7] px-[40px] overflow-y-scroll scroll-smooth">
      <TransferModal sender={user?.accountAddress as `0x${string}`} isOpen={modalOpen} onClose={() => setModalOpen(false)} />

      <img className="absolute w-[400px] h-[350px] bg-transparent rounded-t-xl" src="/wallpaper.png" alt="background" />

      <div className='flex z-10 text-[#eee9ed] items-center mt-6 w-full'>
        <img className='w-[27px] h-[27px] rounded-full' src={user?.profileImageUrl ? user.profileImageUrl : "https://assets.raribleuserdata.com/prod/v1/image/t_image_big/aHR0cHM6Ly9pcGZzLnJhcmlibGV1c2VyZGF0YS5jb20vaXBmcy9RbVNjWXRHdmhRQXpYNVNjaldVZEo1QlY3eE1KSkc4Ulo5REoyMkRHNkFhanVrLzMzMDgucG5n"} alt="" />
        <div className='flex items-center cursor-pointer group select-none'>
          <Copy onClick={() => { navigator.clipboard.writeText(user?.accountAddress || "") }} className="inline-block ml-2 group-active:scale-120 transition" size={14} />
          <p className='text-[17px] ml-1'>{user?.walletName}</p>
        </div>
        <div onClick={() => {logOut()}} className='flex items-center justify-center p-3 h-[28px] bg-[#a479fa79] hover:bg-[#8a52fc] transition text-white hover:scale-105 rounded-full ml-auto mr-0 cursor-pointer'>
          <LogOut color='#c4c3c3' size={14}/>
        </div>
      </div>

      {/* Balance */}
      <div className="flex flex-col items-center text-center z-10 mt-14">
        <p className="text-sm text-[#b29bdf] leading-[12px]">Current Balance</p>
        <h2 className="text-[42px] mt-[2px] font-bold text-[#eee9ed] leading-[46px]">{totalBalance ? `$${totalBalance.toFixed(2)}` : <LoaderCircle className="animate-spin mt-4" />}</h2>
      </div>

      {/* Actions */}
      <div className="flex justify-between mt-15 z-10 gap-10">
        {["Receive", "Transfer", "Swap", "Buy"].map((label, i) => (
          <div key={i} className="flex flex-col items-center gap-1">
            <div className="w-10 h-10 rounded-full bg-[#a479fa79] flex items-center justify-center text-xl text-[#ffffff] cursor-pointer hover:bg-[#8a52fc] transition hover:scale-105">
              {label === "Receive" && <QrCode className="inline-block" size={18} />}
              {label === "Transfer" && <Send onClick={() => setModalOpen(true)} className="inline-block" size={18} />}
              {label === "Swap" && <ArrowRightLeft className="inline-block" size={18} />}
              {label === "Buy" && <DollarSign className="inline-block" size={18} />}
            </div>
            <span className="text-xs mt-1 text-[#eee9ed]">{label}</span>
          </div>
        ))}
      </div>

      {/* Receive Assets */}
      <div className="mt-6 bg-gray-100 rounded-xl px-4 pt-4 pb-5 z-10 w-full shadow-2xl">
        <div className='flex items-center justify-between'>
          <div className='flex flex-col'>
            <p className="text-lg font-semibold">Receive Assets</p>
            <p className="text-xs text-[#292929] w-[160px]">Copy your unique address to receive money.</p>
          </div>
          <div className='flex items-center justify-center bg-[#e9e9e9] rounded-full h-[48px] w-[48px]'>
            <div className='flex rounded-sm bg-[#956fe0] p-1'>
              <ScanLine size={15} strokeWidth={2.5} color="#fff" />
            </div>
          </div>
        </div>
        <button onClick={() => { navigator.clipboard.writeText(user?.accountAddress || "") }} className="w-full flex items-center active:scale-95 justify-center gap-2 bg-[#965cef] hover:bg-[#8a52fc] text-white text-sm font-medium py-2 cursor-pointer rounded-full transition mt-5">
          <Copy className="inline-block ml-2 group-active:scale-120 transition" size={16} strokeWidth={2.5} />
          <p>Copy Address</p>
        </button>
      </div>

      {/* Transaction History */}
      <div className="mt-6 mb-4 w-full">
        <div className="flex justify-between items-center mb-2">
          <p className="text-lg font-semibold">Transaction History</p>
          <button className="text-sm text-[#303030] cursor-pointer">See All</button>
        </div>
        <div className="flex flex-col gap-1 mt-4">
          {txHistory.length === 0 ? (
            <div className="w-full h-[100px] flex items-center justify-center text-sm text-gray-500 bg-white shadow-md rounded-lg px-4">
              No transactions yet. Make your first swap to get started 🚀
            </div>
          ) : ( txHistory.slice(-3).reverse().map((tx, index) => (
            <div key={index} className="flex gap-1 bg-white shadow-md rounded-lg px-4 h-[90px]">
              <div className="flex items-center gap-2">
                <img
                  className="h-7 w-7 rounded-full"
                  src="https://altcoinsbox.com/wp-content/uploads/2023/03/okx-logo-black-and-white.jpg"
                  alt="dex"
                />
                <div className="flex flex-col">
                  <span className="text-[12px] font-medium">{tx.dexName}</span>
                  <span className="text-[11px] text-[#555]">
                    {new Date(tx.timestamp).toLocaleTimeString()} • {new Date(tx.timestamp).toLocaleDateString()}
                  </span>
                </div>
              </div>

              {/* ORTA */}
              <div className="flex flex-col items-end gap-[6px] ml-auto mr-0 my-auto">
                {/* FROM */}
                <div className="text-sm font-semibold flex gap-[6px] items-center text-[#d86060]">
                  <Minus className="inline-block" size={12} strokeWidth={2.5} />
                  <p>{tx.fromAmount}</p>
                  <img
                    className="rounded-full w-5 h-5"
                    src={`${(tx.fromToken)}.png`}
                    alt=""
                  />
                </div>

                {/* TO */}
                <div className="text-sm font-semibold flex gap-[6px] items-center text-[#60d870]">
                  <Plus className="inline-block" size={12} strokeWidth={2.5} />
                  <p>{tx.toAmount}</p>
                  <img
                    className="rounded-full w-5 h-5"
                    src={`${(tx.toToken)}.png`}
                    alt=""
                  />
                </div>
              </div>
            </div>
          ))
          )}
        </div>
      </div>
    </div>
  )
}

export default HomeSection