"use client"

import Image from 'next/image';
import { MoveRight } from 'lucide-react';
import { Button } from "@/components/ui/button";
import { motion } from "framer-motion";

export default function Home() {

  return (
    <main className="relative flex flex-col items-center h-screen w-screen bg-black overflow-hidden font-mono">
      {/* 像素化景观背景 / Pixelated landscape background */}
      <div className="absolute inset-0">
        {/* 基础像素网格 / Base pixel grid */}
        <div className="w-full h-full opacity-40" style={{
          backgroundImage: `
            linear-gradient(rgba(0, 255, 255, 0.4) 1px, transparent 1px),
            linear-gradient(90deg, rgba(0, 255, 255, 0.4) 1px, transparent 1px)
          `,
          backgroundSize: '16px 16px'
        }}></div>
        
        {/* 像素化山脉效果 / Pixelated mountain effect */}
        <div className="absolute bottom-0 left-0 w-full h-1/3 opacity-60">
          <div className="w-full h-full" style={{
            backgroundImage: `
              linear-gradient(45deg, transparent 30%, rgba(0, 255, 255, 0.3) 30%, rgba(0, 255, 255, 0.3) 32%, transparent 32%),
              linear-gradient(-45deg, transparent 30%, rgba(0, 255, 255, 0.3) 30%, rgba(0, 255, 255, 0.3) 32%, transparent 32%),
              linear-gradient(90deg, transparent 30%, rgba(255, 255, 0, 0.2) 30%, rgba(255, 255, 0, 0.2) 32%, transparent 32%)
            `,
            backgroundSize: '64px 64px, 64px 64px, 32px 32px'
          }}></div>
        </div>

        {/* 像素化星空效果 / Pixelated starfield effect */}
        <div className="absolute inset-0 opacity-70">
          {[...Array(50)].map((_, i) => (
            <div
              key={i}
              className="absolute w-1 h-1 bg-cyan-400"
              style={{
                left: `${Math.random() * 100}%`,
                top: `${Math.random() * 100}%`,
                animation: `twinkle ${2 + Math.random() * 3}s infinite ${Math.random() * 2}s`
              }}
            />
          ))}
        </div>

        {/* 像素化科技线条 / Pixelated tech lines */}
        <div className="absolute inset-0 opacity-50">
          <div className="w-full h-full" style={{
            backgroundImage: `
              linear-gradient(90deg, transparent 49%, rgba(0, 255, 255, 0.6) 49%, rgba(0, 255, 255, 0.6) 51%, transparent 51%),
              linear-gradient(0deg, transparent 49%, rgba(255, 255, 0, 0.4) 49%, rgba(255, 255, 0, 0.4) 51%, transparent 51%)
            `,
            backgroundSize: '128px 128px, 128px 128px'
          }}></div>
        </div>

        {/* 像素化能量波纹 / Pixelated energy waves */}
        <div className="absolute inset-0 opacity-30">
          <div className="w-full h-full" style={{
            backgroundImage: `
              radial-gradient(circle at 20% 80%, rgba(0, 255, 255, 0.3) 0%, transparent 50%),
              radial-gradient(circle at 80% 20%, rgba(255, 255, 0, 0.2) 0%, transparent 50%),
              radial-gradient(circle at 40% 40%, rgba(255, 0, 255, 0.2) 0%, transparent 50%)
            `,
            backgroundSize: '200px 200px, 300px 300px, 250px 250px'
          }}></div>
        </div>
      </div>

      {/* 浮动像素粒子效果 - 增强版 / Enhanced floating pixel particles */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        {[...Array(25)].map((_, i) => (
          <motion.div
            key={i}
            className={`absolute ${i % 3 === 0 ? 'w-2 h-2 bg-cyan-400' : i % 3 === 1 ? 'w-3 h-3 bg-yellow-400' : 'w-1 h-1 bg-pink-400'} opacity-80`}
            style={{
              left: `${Math.random() * 100}%`,
              top: `${Math.random() * 100}%`,
            }}
            animate={{
              y: [0, -40, 0],
              x: [0, Math.random() * 20 - 10, 0],
              opacity: [0.8, 1, 0.8],
              scale: [1, 1.2, 1],
            }}
            transition={{
              duration: 5 + Math.random() * 4,
              repeat: Infinity,
              delay: Math.random() * 3,
            }}
          />
        ))}
      </div>

      {/* 像素化扫描线效果 / Pixelated scanline effect */}
      <div className="absolute inset-0 opacity-20 pointer-events-none">
        <div className="w-full h-full" style={{
          backgroundImage: `
            linear-gradient(0deg, transparent 0%, rgba(0, 255, 255, 0.1) 1%, transparent 2%)
          `,
          backgroundSize: '100% 4px',
          animation: 'scanline 8s linear infinite'
        }}></div>
      </div>

      <header className="w-full py-6 flex items-center justify-between max-w-[600px] mx-auto z-20">
        <div className="flex items-center gap-2">
          {/* 像素化 Logo / Pixelated Logo */}
          <div className="relative">
            <Image src="/logo.png" alt="Smart Logo" width={38} height={38} className="filter contrast-150 brightness-110" />
            <div className="absolute inset-0 border-2 border-cyan-400 opacity-70" style={{ clipPath: 'polygon(0 0, 100% 0, 90% 100%, 10% 100%)' }}></div>
          </div>
          <div className="flex flex-col">
            <p className="text-xl font-bold text-cyan-400 tracking-wider" style={{ textShadow: '0 0 8px rgba(0, 255, 255, 0.8)' }}>
              Smart
            </p>
            <div className="w-full h-1 bg-gradient-to-r from-cyan-400 to-transparent opacity-80"></div>
          </div>
        </div>
        
        {/* 保持原有按钮样式 / Keep original button style */}
        <Button className="bg-[#fafafa] text-[#1f1f1f] text-[12px] hover:bg-[#fff] border-[#c8c7c7] border w-[90px] cursor-pointer">
          <p>Contact Us</p>
        </Button>
      </header>

      <div key="landing" className="flex flex-col items-center justify-center mt-auto mb-0 relative z-10">
        {/* 主标题 - 像素游戏风格 / Main title - Pixel gaming style */}
        <motion.div className="text-center mb-8"
          initial={{ opacity: 0, y: 50 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.6, delay: 0.3 }}
        >
          <h1 className="text-center leading-tight font-bold text-[48px] md:text-[52px] mb-8 font-mono tracking-wider" style={{ textShadow: '0 0 15px rgba(0, 255, 255, 0.6)' }}>
            <span className="block leading-14 text-cyan-400 mb-2">DISCOVER THE</span>
            <span className="block leading-14 text-yellow-400">FREEDOM OF</span>
            <span className="block leading-14 text-pink-400">SELF-CUSTODY</span>
          </h1>
          
          {/* 像素化装饰边框 / Pixelated decorative border */}
          <div className="flex justify-center items-center gap-2 mb-4">
            <div className="w-4 h-4 bg-cyan-400"></div>
            <div className="w-8 h-2 bg-yellow-400"></div>
            <div className="w-4 h-4 bg-pink-400"></div>
            <div className="w-8 h-2 bg-cyan-400"></div>
            <div className="w-4 h-4 bg-yellow-400"></div>
          </div>
        </motion.div>

        {/* 副标题 - 复古游戏风格 / Subtitle - Retro gaming style */}
        <motion.p className="text-gray-300 font-mono text-[15px] text-center w-[500px] mb-6 leading-relaxed"
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.6, delay: 0.6 }}
        >
          <span className="text-cyan-400">[</span> One simple gateway to all your decentralized apps. 
          Own your keys, control your data — stay truly sovereign with Smart. <span className="text-cyan-400">]</span>
        </motion.p>

        {/* 行动按钮区域 - 像素游戏风格按钮 / Action buttons area - Pixel gaming style buttons */}
        <motion.div className="flex gap-6 mb-18 z-10 px-4"
          initial={{ opacity: 0, y: 40 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.6, delay: 0.9 }}
        >
          {/* Turnkey 安全按钮 - 像素科技风格 / Turnkey security button - Pixel tech style */}
          <motion.div
            className="relative group cursor-pointer"
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
          >
            <div className="relative bg-gradient-to-br from-gray-900 via-gray-800 to-black border-2 border-cyan-400 px-8 py-4 overflow-hidden w-[200px] h-[80px] flex items-center justify-center"
              style={{ 
                clipPath: 'polygon(0 0, 100% 0, 95% 100%, 5% 100%)',
                boxShadow: '0 0 20px rgba(0, 255, 255, 0.4), inset 0 0 20px rgba(0, 255, 255, 0.1)'
              }}
              onClick={() => window.open("https://turnkey.com", "_blank")}
            >
              {/* 动态呼吸光效 / Dynamic breathing glow effect */}
              <motion.div 
                className="absolute inset-0 bg-gradient-to-r from-cyan-400/20 to-transparent"
                animate={{
                  opacity: [0.2, 0.6, 0.2],
                }}
                transition={{
                  duration: 3,
                  repeat: Infinity,
                  ease: "easeInOut"
                }}
              />
              
              {/* 像素化边框装饰 / Pixelated border decoration */}
              <div className="absolute -top-1 -left-1 w-2 h-2 bg-cyan-400 opacity-80"></div>
              <div className="absolute -top-1 -right-1 w-2 h-2 bg-cyan-400 opacity-80"></div>
              <div className="absolute -bottom-1 -left-1 w-2 h-2 bg-cyan-400 opacity-80"></div>
              <div className="absolute -bottom-1 -right-1 w-2 h-2 bg-cyan-400 opacity-80"></div>
              
              {/* 按钮内容 / Button content */}
              <div className="flex flex-col items-center gap-2 relative z-10">
                <p className="text-sm font-bold tracking-wider text-cyan-400">SECURED BY</p>
                <img src="turnkey.svg" alt="turnkey logo" className='w-20 h-8 filter brightness-0 invert' />
              </div>
              
              {/* 扫描线效果 / Scanline effect */}
              <motion.div 
                className="absolute inset-0 bg-gradient-to-b from-transparent via-cyan-400/30 to-transparent"
                animate={{
                  y: [-100, 100],
                }}
                transition={{
                  duration: 2,
                  repeat: Infinity,
                  ease: "linear"
                }}
              />
            </div>
          </motion.div>

          {/* 开始使用按钮 - 主要 CTA 像素风格 / Start using button - Main CTA pixel style */}
          <motion.div
            className="relative group cursor-pointer"
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            animate={{
              y: [0, -8, 0],
            }}
            transition={{
              duration: 1.5,
              repeat: Infinity,
              ease: "easeInOut"
            }}
          >
            <div className="relative bg-gradient-to-br from-yellow-500 via-orange-500 to-red-500 border-2 border-yellow-400 px-8 py-4 overflow-hidden w-[200px] h-[80px] flex items-center justify-center"
              style={{ 
                clipPath: 'polygon(0 0, 100% 0, 95% 100%, 5% 100%)',
                boxShadow: '0 0 25px rgba(255, 255, 0, 0.6), inset 0 0 20px rgba(255, 255, 0, 0.2)'
              }}
              onClick={() => window.location.href = "/wallet"}
            >
              {/* 动态呼吸光效 / Dynamic breathing glow effect */}
              <motion.div 
                className="absolute inset-0 bg-gradient-to-r from-yellow-400/30 to-orange-400/30"
                animate={{
                  opacity: [0.3, 0.7, 0.3],
                }}
                transition={{
                  duration: 2.5,
                  repeat: Infinity,
                  ease: "easeInOut"
                }}
              />
              
              {/* 像素化边框装饰 / Pixelated border decoration */}
              <div className="absolute -top-1 -left-1 w-2 h-2 bg-yellow-400 opacity-80"></div>
              <div className="absolute -top-1 -right-1 w-2 h-2 bg-yellow-400 opacity-80"></div>
              <div className="absolute -bottom-1 -left-1 w-2 h-2 bg-yellow-400 opacity-80"></div>
              <div className="absolute -bottom-1 -right-1 w-2 h-2 bg-yellow-400 opacity-80"></div>
              
              {/* 按钮内容 / Button content */}
              <div className="flex flex-col items-center gap-2 relative z-10">
                <span className="text-sm font-bold tracking-wider text-black">START USING</span>
                <span className="text-sm font-bold tracking-wider text-black">EOLIA</span>
              </div>
              
              {/* 能量波纹效果 - 增强版 / Enhanced energy wave effect */}
              <motion.div 
                className="absolute inset-0 border-2 border-yellow-300"
                animate={{
                  scale: [1, 1.2, 1],
                  opacity: [0.8, 0.1, 0.8],
                }}
                transition={{
                  duration: 2,
                  repeat: Infinity,
                  ease: "easeInOut"
                }}
              />
              
              {/* 跃动光效 / Jumping glow effect */}
              <motion.div 
                className="absolute inset-0 bg-gradient-to-t from-yellow-400/20 to-transparent"
                animate={{
                  y: [0, -20, 0],
                  opacity: [0.2, 0.6, 0.2],
                }}
                transition={{
                  duration: 1.5,
                  repeat: Infinity,
                  ease: "easeInOut"
                }}
              />
              
              {/* 脉冲边框 / Pulsing border */}
              <motion.div 
                className="absolute inset-0 border-2 border-yellow-200"
                animate={{
                  scale: [1, 1.05, 1],
                  opacity: [0.5, 1, 0.5],
                }}
                transition={{
                  duration: 1.5,
                  repeat: Infinity,
                  ease: "easeInOut"
                }}
              />
            </div>
          </motion.div>
        </motion.div>

        {/* 手机模拟图 - 像素化处理 / Phone mockup - Pixelated treatment */}
        <motion.div className='relative z-10'
          initial={{ opacity: 0, y: 50 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: 300 }}
          transition={{ duration: 0.8, delay: 1.2 }}
        >
          <div className="relative">
            <Image
              src="/mockup-phone.png"
              alt="Hero Image"
              width={600}
              height={400}
              className="overflow-hidden z-10 filter contrast-125 brightness-110"
            />
            {/* 像素化边框效果 / Pixelated border effect */}
            <div className="absolute inset-0 border-4 border-cyan-400 opacity-60" style={{ 
              clipPath: 'polygon(0 0, 100% 0, 95% 100%, 5% 100%)',
              filter: 'drop-shadow(0 0 20px rgba(0, 255, 255, 0.6))'
            }}></div>
            
            {/* 像素化装饰元素 / Pixelated decorative elements */}
            <div className="absolute -top-4 -left-4 w-8 h-8 bg-yellow-400 opacity-80"></div>
            <div className="absolute -top-4 -right-4 w-8 h-8 bg-pink-400 opacity-80"></div>
            <div className="absolute -bottom-4 -left-4 w-8 h-8 bg-cyan-400 opacity-80"></div>
            <div className="absolute -bottom-4 -right-4 w-8 h-8 bg-yellow-400 opacity-80"></div>
          </div>
        </motion.div>

        {/* 底部像素化装饰条 / Bottom pixelated decorative bar */}
        <motion.div className="absolute bottom-8 left-1/2 transform -translate-x-1/2 w-full max-w-[600px] px-4"
          initial={{ opacity: 0, scaleX: 0 }}
          animate={{ opacity: 1, scaleX: 1 }}
          transition={{ duration: 1, delay: 1.5 }}
        >
          <div className="flex items-center justify-center gap-1">
            {[...Array(30)].map((_, i) => (
              <div key={i} className={`w-2 h-2 ${i % 4 === 0 ? 'bg-cyan-400' : i % 4 === 1 ? 'bg-yellow-400' : i % 4 === 2 ? 'bg-pink-400' : 'bg-cyan-400'} opacity-80`}></div>
            ))}
          </div>
        </motion.div>
      </div>

      {/* 自定义CSS动画 / Custom CSS animations */}
      <style jsx>{`
        @keyframes twinkle {
          0%, 100% { opacity: 0.3; }
          50% { opacity: 1; }
        }
        @keyframes scanline {
          0% { transform: translateY(-100%); }
          100% { transform: translateY(100vh); }
        }
      `}</style>
    </main>
  );
}