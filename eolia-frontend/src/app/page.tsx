"use client"

import Image from 'next/image';
import { MoveRight } from 'lucide-react';
import { Button } from "@/components/ui/button";
import { GridPatternLinearGradient } from "@/components/ui/grid-pattern-linear-gradient";
import { motion } from "framer-motion";

export default function Home() {

  return (
    <main className="relative flex flex-col items-center h-screen w-screen bg-[#f4f4fd] overflow-hidden font-sfpro">
      <header className="w-full py-6 flex items-center justify-between max-w-[600px] mx-auto z-20">
        <div className="flex items-center gap-2">
          <Image src="/logo.png" alt="Smart Logo" width={38} height={38} />
          <p className="text-xl font-semibold text-[#0e0e0e]">Smart</p>
        </div>
        <Button className="bg-[#fafafa] text-[#1f1f1f] text-[12px] hover:bg-[#fff] border-[#c8c7c7] border w-[90px] cursor-pointer">
          <p>Contact Us</p>
        </Button>
      </header>

      <div key="landing" className="flex flex-col items-center justify-center mt-auto mb-0">
        <motion.h1 className="text-center leading-tight font-bold text-[48px] md:text-[52px] mb-8"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.4, delay: 0.6 }}
        >
          <span className="block leading-14">Discover the</span>
          <span className="block leading-14">freedom of self-custody</span>
        </motion.h1>
        <motion.p className="text-[#383838] font-medium text-[15px] text-center w-[500px] mb-6"
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.4, delay: 0.8 }}
        >
          One simple gateway to all your decentralized apps. Own your keys, control your data — stay truly sovereign with Smart.
        </motion.p>
        <motion.div className="flex gap-5 mb-18 z-10"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.4, delay: 1 }}
        >
          <Button
            style={{ boxShadow: 'rgba(0, 0, 0, 0.24) 0px 3px 8px' }}
            className="bg-[#fafafa] text-[#1f1f1f] text-[12px] hover:bg-[#fff] border-[#c8c7c7] border w-[200px] cursor-pointer"
            onClick={() => window.open("https://turnkey.com", "_blank")}
          >
            <p className='text-black'>Secured By</p>
            <img src="turnkey.svg" alt="turnkey logo" className='w-[80px]' />
          </Button>
          <Button
            style={{ boxShadow: 'rgba(0, 0, 0, 0.24) 0px 3px 8px' }}
            className="bg-[#292929] text-[#e0e0e0] text-[12px] hover:bg-[#313131] border-[#757577] border w-[170px] cursor-pointer"
            onClick={() => window.location.href = "/wallet"}
          >
            <p>Start using Eolia</p>
            <MoveRight />
          </Button>
        </motion.div>
        <GridPatternLinearGradient />
        <motion.div className='z-10'
          initial={{ opacity: 0, y: 0 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: 300 }}
          transition={{ duration: 0.4, delay: 0.6 }}
        >
          <Image
            src="/mockup-phone.png"
            alt="Hero Image"
            width={600}
            height={400}
            className="overflow-hidden z-10"
          />
        </motion.div>
      </div>
    </main>
  );
}