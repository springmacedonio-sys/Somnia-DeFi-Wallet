# 🌐 Eolia Frontend

**Eolia Frontend** is the user interface for the **Eolia Smart Wallet** — an Account Abstraction-based wallet built for the **OKX ETH CC Hackathon**.

This frontend allows users to:
  
- 🔑 Log in via OAuth (Google, GitHub, Apple)
- 💱 Swap tokens directly using **OKX DEX API**
- 📊 View balances, token lists, and transaction history
- 🛡️ Interact securely with the Eolia Signer & Bundler services

---

## ✨ Features

- **Next.js 15+ App Router** with TypeScript
- **OKX DEX API Integration** for on-wallet swaps
- **Account Abstraction (ERC‑4337)** support via Eolia Signer & Bundler
- **next-auth** for OAuth login
- **Viem** for blockchain RPC connections
- Clean, animated UI with **framer-motion**

---

## 📂 Project Structure

```plaintext
eolia-frontend
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
│   ├── lib/                    # API clients (OKX, balance, tx utils)
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
- Running instances of **Eolia Signer** and **Eolia Bundlr**

### Steps

1. **Clone the repository**
  
   ```bash
   git clone https://github.com/EoliaWallet/eolia-frontend.git
   cd eolia-frontend
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
   NEXT_OKX_API_KEY=...
   NEXT_OKX_SECRET_KEY=...
   NEXT_OKX_API_PASSPHRASE=...
   NEXT_OKX_PROJECT_ID=...
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
- **OKX DEX API**
- **next-auth** – OAuth
- **framer-motion** – Animations
- **Tailwind CSS** – Styling

---

## 📜 License

MIT
