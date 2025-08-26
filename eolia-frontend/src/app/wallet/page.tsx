// components/Wallet.tsx
"use client";

import { ArrowRightLeft, House, Clock, WalletIcon, Settings } from 'lucide-react';
import { signIn, signOut, useSession } from "next-auth/react";
import { Button } from "@/components/ui/button";
import { useEffect, useState } from 'react';
import HomeSection from '../../components/sections/home-section';
import AssetsSection from '../../components/sections/assets-section';
import { AnimatePresence, motion } from 'motion/react';
import SwapSection from '../../components/sections/swap-section';
import HistorySection from '../../components/sections/history-section';
import Image from 'next/image';
import { useWeb3 } from "@/context/web3Context";
import SettingsSection from '@/components/sections/settings-section';
import { WalletProps } from '@/types/wallet';

const authProviders = [
  {
    name: "Email",
    icon: "/mail.svg",
  },
  {
    name: "Google",
    icon: "https://www.citypng.com/public/uploads/preview/google-logo-icon-gsuite-hd-701751694791470gzbayltphh.png",
  },
  {
    name: "Apple",
    icon: "https://upload.wikimedia.org/wikipedia/commons/f/fa/Apple_logo_black.svg",
  },
  {
    name: "Github",
    icon: "https://img.icons8.com/m_outlined/512/github.png",
  },
];

interface SessionData {
  user: {
    provider?: string;
    externalId?: string;
  };
}

export default function Wallet() {
  const { fetchBalances, fetchTxHistory } = useWeb3();

  const [ page, setPage ] = useState("home");
  
  const { data: session, status } = useSession();

  const [phase, setPhase] = useState<"checking" | "needs-register" | "needs-auth" | "ready">("checking");
  
  const [walletName, setWalletName] = useState("");
  const [walletNameValid, setWalletNameValid] = useState(false);
  const [errorMessage, setErrorMessage] = useState("");
  const [profileImageUrl, setProfileImageUrl] = useState("");

  const [user, setUser] = useState<WalletProps | null>(null);

  useEffect(() => {
    (async () => {
      if (!session && status === "loading") {
        setPhase("checking");
        return;
      }
      try {
        const resMe = await fetch("http://localhost:8080/me", { credentials: "include" });
        if (resMe.ok) {
          const data = await resMe.json();
          setUser({
            walletName: data.private_user_info.wallet_name,
            accountAddress: data.private_user_info.account_address,
            profileImageUrl: data.private_user_info.profile_image_url,
            authProvider: data.private_user_info.auth_provider,
            authExternalId: data.private_user_info.auth_external_id,
            lastLogin: data.private_user_info.last_login,
            createdAt: data.private_user_info.created_at,
          });
          setPhase("ready");
          return;
        } else {
          const errorData = await resMe.json();
          console.warn("Error auto login:", errorData);
        }

        if ( session && (session as SessionData).user?.externalId && (session as SessionData).user?.provider) {
          const resAuth = await fetch("http://localhost:8080/auth", {
            method: "POST",
            credentials: "include",
            body: JSON.stringify({
              auth_provider: (session as SessionData).user.provider,
              auth_external_id: (session as SessionData).user.externalId,
            }),
            headers: { "Content-Type": "application/json" },
          });

          if (resAuth.ok) {
            const data = await resAuth.json();
            setUser({
              walletName: data.private_user_info.wallet_name,
              accountAddress: data.private_user_info.account_address,
              profileImageUrl: data.private_user_info.profile_image_url,
              authProvider: data.private_user_info.auth_provider,
              authExternalId: data.private_user_info.auth_external_id,
              lastLogin: data.private_user_info.last_login,
              createdAt: data.private_user_info.created_at,
            });
            setPhase("ready");
          } else {
            const errorData = await resAuth.json();
            console.warn("Error logging in:", errorData);
            setPhase("needs-register");
          }
        } else {
          setPhase("needs-auth");
        }
      } catch (e) {
        console.error("auth error", e);
        if (session) {
          setPhase("needs-register");
        } else {
          setPhase("needs-auth");
        }
      }
    })();
  }, [status]);

  useEffect(() => {
    if (user?.accountAddress) {
      fetchBalances(user.accountAddress as `0x${string}`);
      fetchTxHistory(user.accountAddress as `0x${string}`);
    }
  }, [user?.accountAddress]);

  const handleCreateWallet = async () => {
    try {
      if (!walletName ) return;
      setPhase("checking");
      const resRegister = await fetch("http://localhost:8080/auth/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({
          auth_provider: (session as SessionData).user.provider,
          auth_external_id: (session as SessionData).user.externalId,
          wallet_name: walletName,
          profile_image_url: profileImageUrl,
        }),
      })
      if (resRegister.ok) {
        const data = await resRegister.json();
        setUser({
          walletName: data.private_user_info.wallet_name,
          accountAddress: data.private_user_info.account_address,
          profileImageUrl: data.private_user_info.profile_image_url,
          authProvider: data.private_user_info.auth_provider,
          authExternalId: data.private_user_info.auth_external_id,
          lastLogin: data.private_user_info.last_login,
          createdAt: data.private_user_info.created_at,
        });
        setPhase("ready");
      } else {
        const errorData = await resRegister.json();
        console.error("Error creating wallet:", errorData);
      }
    } catch (error) {
      console.error("Error creating wallet:", error);
    }
  }

  const checkIfUserNameValid = async (name: string) => {
    setWalletName(name);
    if (!name) return false;

    if (name.length < 3 || name.length > 20 || !/^[a-zA-Z0-9_.-]+$/.test(name)) {
      setWalletNameValid(false);
      setErrorMessage("Wallet name must be 3-20 characters long and can only contain letters, numbers, underscores, hyphens, and dots.");
      return false;
    }

    const res = await fetch(`http://localhost:8080/auth/${encodeURIComponent(name)}`);
    if (!res.ok) {
      console.error("Error checking username:", res.statusText);
      return false;
    }
    const data = await res.json();

    if (data.occupied) {
      setErrorMessage("This wallet name is already taken. Please choose another.");
      setWalletNameValid(false);
    }

    setErrorMessage("");
    setWalletNameValid(true);
  }

  const handleLogOut = async () => {
    try {
      await fetch("http://localhost:8080/auth/logout", {
        method: "POST",
        credentials: "include",
      });
    } finally {
      setUser(null);
      setPhase("checking")
      signOut();
    }
  }

  return (
    <main style={{ backgroundImage: `url(bg-a.png)` }} className="bg-no-repeat bg-cover bg-center relative flex flex-col items-center h-screen w-screen gap-10 bg-[#f4f4fd] overflow-hidden font-sfpro">
      <header className="w-full py-6 flex items-center justify-between max-w-[600px] mx-auto z-20">
        <div className="flex items-center gap-2">
          <Image src="/logo.png" alt="Eolia Logo" width={38} height={38} />
          <p className="text-xl font-semibold text-[#1f1f1f]">Eolia</p>
        </div>
        <Button className="bg-[#fafafa] text-[#1f1f1f] text-[12px] hover:bg-[#fff] border-[#c8c7c7] border w-[90px] cursor-pointer">
          <p>Contact Us</p>
        </Button>
      </header>

      <AnimatePresence mode='wait'>
        { phase === "checking" ? (
          <motion.div
            key="loading"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.4 }}
            className="absolute top-0 w-full h-full flex items-center justify-center bg-[#f7f7f7]"
          >
            <div className="bg-white shadow-xl p-6 rounded-xl flex items-center justify-center">
              <Image src="/ring.svg" alt="Loading" width={50} height={50} />
            </div>
          </motion.div>
        ) : phase === "ready" ? (
          <motion.div 
            key="ready"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.4 }}
            className='w-[400px] h-[700px] relative flex flex-col items-center shadow-xl rounded-xl bg-[#f7f7f7] px-[40px]'
          >
            <AnimatePresence mode="wait">
              <motion.div
                key={page}
                initial={{ opacity: 0, y: 0 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0 }}
                transition={{ duration: 0.4 }}
              >
                {page === "home" && <HomeSection logOut={handleLogOut} user={user} />}
                {page === "assets" && <AssetsSection user={user} />}
                {page === "swap" && <SwapSection user={user} />}
                {page === "history" && <HistorySection user={user} />}
                {page === "settings" && <SettingsSection />}
              </motion.div>
            </AnimatePresence>

            <div className='absolute top-auto bottom-0 rounded-b-xl w-full h-[60px] flex items-center justify-between px-12 z-10 bg-[#f7f7f7]'>
              <div className='absolute top-0 left-0 bg-[#e6e6e6] w-full h-[1px]' />
              <House onClick={() => setPage("home")} size={24} color={page == 'home' ? '#9257f3' : '#121212'} className='cursor-pointer hover:scale-105 transition' />
              <WalletIcon onClick={() => setPage("assets")} size={24} color={page == 'assets' ? '#9257f3' : '#121212'} className='cursor-pointer hover:scale-105 transition' />
              <div onClick={() => setPage("swap")} className='bg-[#9257f3] w-[40px] h-[40px] rounded-full flex items-center justify-center hover:rotate-180 cursor-pointer transition'>
                <ArrowRightLeft size={20} color='#f7f7f7' />
              </div>
              <Clock onClick={() => setPage("history")} size={24} color={page == 'history' ? '#9257f3' : '#121212'} className='cursor-pointer hover:scale-105 transition' />
              <Settings onClick={() => setPage("settings")} size={24} color={page == 'settings' ? '#9257f3' : '#121212'} className='cursor-pointer hover:scale-105 transition' />
            </div>
          </motion.div>
        ) : phase === "needs-register" ? (
          <motion.div
            key="needs-register"
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: 30 }}
            transition={{ duration: 0.4 }}
            className="w-[400px] h-[700px] flex items-center justify-center bg-white rounded-3xl border border-[#e4e7ec] shadow-[0_8px_28px_rgba(0,0,0,0.06)] px-4"
          >
            <div className="w-full h-full p-8 flex flex-col gap-6">
              {/* Logo Section */}
              <div className="flex items-center justify-center mb-2">
                <Image src="/logo.png" alt="Smart Logo" width={48} height={48} />
              </div>

              {/* Title & Description */}
              <div className="text-center">
                <h2 className="text-[22px] font-bold text-[#111] mb-1">Create Your Smart Wallet</h2>
                <p className="text-sm text-[#555] max-w-[320px] mx-auto leading-snug">
                  Start by choosing a wallet name and optional profile image.
                </p>
              </div>

              {/* Input Fields */}
              <div className="flex flex-col gap-5 w-full mt-2">
                <div className="flex flex-col gap-1">
                  <label className="text-sm font-semibold text-[#333]">Wallet Name</label>
                  <input
                    type="text"
                    value={walletName}
                    onChange={(e) => checkIfUserNameValid(e.target.value)}
                    placeholder="e.g. arcyx.eth"
                    className="px-4 py-2.5 border border-[#d1d5db] rounded-xl bg-[#fafafa] text-sm focus:outline-none focus:ring-2 focus:ring-[#111] focus:bg-white transition duration-150"
                  />
                  {errorMessage && <p className="text-sm text-red-500 mt-2">{errorMessage}</p>}
                </div>

                <div className="flex flex-col gap-1">
                  <label className="text-sm font-semibold text-[#333]">Profile Image URL (optional)</label>
                  <input
                    type="text"
                    value={profileImageUrl}
                    onChange={(e) => setProfileImageUrl(e.target.value)}
                    placeholder="https://your-image-url.com"
                    className="px-4 py-2.5 border border-[#d1d5db] rounded-xl bg-[#fafafa] text-sm focus:outline-none focus:ring-2 focus:ring-[#111] focus:bg-white transition duration-150"
                  />
                </div>
              </div>

              {/* Preview (Optional) */}
              {profileImageUrl && (
                <div className="flex justify-center mt-2">
                  <img
                    src={profileImageUrl}
                    alt="preview"
                    className="h-16 w-16 rounded-full border border-[#ddd] shadow"
                  />
                </div>
              )}

              {/* Submit Button */}
              <button
                onClick={handleCreateWallet}
                className={`mt-4 w-full py-2.5 rounded-xl ${walletNameValid ? "bg-[#111] text-white hover:scale-[1.015] active:scale-[0.98] hover:bg-[#000] cursor-pointer" : "bg-[#f0f0f0] text-[#aaa]"} text-sm font-semibold transition shadow-md`}
                disabled={!walletNameValid}
              >
                Create Wallet
              </button>

              {/* Footer Note */}
              <p className="text-[11px] text-center text-[#999] mt-auto">
                Powered by <span className="font-semibold">Eolia</span>
              </p>
            </div>
          </motion.div>
        ) : phase === "needs-auth" ? (
          <motion.div
            key="needs-auth"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.4 }}
            className="w-[400px] h-[700px] flex items-center justify-center bg-white rounded-2xl shadow-xl border border-[#eee] px-4"
          >
            <div className="w-full h-full bg-white rounded-2xl p-8 flex flex-col gap-6 items-center">
              {/* Logo Section */}
              <div className="flex items-center justify-center mb-2">
                <Image src="/logo.png" alt="Smart Logo" width={48} height={48} />
              </div>
              
              <div className="text-center">
                <h2 className="text-[20px] font-bold text-[#111] mb-3">Sign In To Your Smart Wallet</h2>
                <p className="text-sm text-[#555] max-w-[320px] mx-auto leading-snug">
                  Choose a method to continue.
                </p>
              </div>

              <div className="flex flex-col gap-4 w-full mt-5">
                {authProviders.map((provider) => (
                  <Button
                    key={provider.name}
                    className="flex items-center cursor-pointer justify-start gap-4 bg-[#fafafa] text-[#1f1f1f] hover:bg-white transition border active:scale-95 border-[#ddd] px-4 py-2 rounded-xl w-full h-[48px]"
                    onClick={() => signIn(provider.name.toLowerCase(), { callbackUrl: "/wallet" })}
                  >
                    <img src={provider.icon} alt={`${provider.name} icon`} width={20} height={20} />
                    <span className="text-sm font-medium">{provider.name}</span>
                  </Button>
                ))}
              </div>

              {/* Footer Note */}
              <p className="text-[11px] text-center text-[#999] mt-auto">
                Powered by <span className="font-semibold">Eolia</span>
              </p>
            </div>
          </motion.div>
        ) : null}
      </AnimatePresence>
    </main>
  );
}