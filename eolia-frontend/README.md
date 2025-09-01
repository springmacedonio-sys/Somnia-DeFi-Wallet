# 🌐 Somnia Frontend

**Somnia Frontend** is the user interface for the **Somnia Smart Wallet** — an Account Abstraction-based wallet built for the **Somnia DeFi Mini Hackathon**.

This frontend allows users to:
  
- 🔑 Log in via OAuth (Google, GitHub, Apple)
- 💱 Swap tokens directly using **Somnex**
- 📊 View balances, token lists, and transaction history
- 🛡️ Interact securely with the Somnia Signer & Bundler services

---

## ✨ Features

- **Next.js 15+ App Router** with TypeScript
- **Somnia DEX API Integration** for on-wallet swaps
- **Account Abstraction (ERC‑4337)** support via Somnia Signer & Bundler
- **next-auth** for OAuth login
- **Viem** for blockchain RPC connections
- Clean, animated UI with **framer-motion**

---

## 📂 Project Structure

```plaintext
somnia-frontend
│
├── public/                   # Static assets
├── src/
│   ├── app/                   # Next.js App Router pages
│   │   ├── api/auth/[…nextauth]  # OAuth backend routes
│   │   ├── fonts/              # Custom fonts
│   │   └── wallet/             # Wallet UI routes
│   ├── components/             # Reusable UI components
│   │   ├── sections/           # Page sections
│   │   └── ui/                 
│   ├── constants/              # Token lists, constants
│   ├── context/                # React context (e.g., Web3Context)
│   ├── lib/                    # API clients (Somnia, balance, tx utils)
│   ├── types/                  # Shared TS types
│   └── utils/                  # Formatting, decimals, helper functions
├── .env                        # Environment variables (not committed)
├── next.config.js              
├── eslint.config.mjs           
├── package.json                
└── README.md
```

---

## 🔧 Setup & Installation

### Requirements
  
- Node.js 18+
- pnpm / yarn / npm
- Running instances of **Somnia Signer** and **Somnia Bundlr**

### Steps

1. **Clone the repository**
  
   ```bash
   git clone https://github.com/springmacedonio-sys/Somnia-DeFi-Wallet.git
   cd somnia-frontend
   ```

2. **Install dependencies**
  
   ```bash
   pnpm install
   # or
   yarn install
   # or
   npm install
   ```

3. **Set environment variables** in `.env`:
  
   ```env
   NEXT_GOOGLE_CLIENT_ID=...
   NEXT_GOOGLE_CLIENT_SECRET=...
   NEXT_GITHUB_CLIENT_ID=...
   NEXT_GITHUB_CLIENT_SECRET=...
   NEXT_AUTH_SECRET=...
   NEXT_PUBLIC_ENTRYPOINT_ADDRESS=...
   ```

4. **Run the development server**
  
   ```bash
   pnpm dev
   ```

5. **Build for production**
  
   ```bash
   pnpm build
   pnpm start
   ```

---

## 🧩 Tech Stack

- **Next.js 15+ App Router**
- **TypeScript**
- **Viem** – Blockchain RPC
- **next-auth** – OAuth
- **framer-motion** – Animations
- **Tailwind CSS** – Styling

---

## 📜 License

MIT
