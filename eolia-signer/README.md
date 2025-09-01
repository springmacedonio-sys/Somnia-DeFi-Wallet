# 🚀 Somnia Signer

**Somnia Signer** is the secure signing backend service for **Somnia Smart Wallet** — an Account Abstraction-based wallet built for the **Somnia DeFi Mini Hackathon**.

This service handles:
  
- Secure account creation & key management via **Turnkey**
- OAuth-based login & authentication
- Preparing & signing **UserOperations**
- Transaction history tracking for the frontend

Somnia Signer is **1 of the 4 core components** of the Somnia ecosystem:

1. **Frontend** (Next.js) – User interface  
2. **Somnia Signer** (this service) – Secure signing & account logic  
3. **Somnia Bundler** – UserOperation validation & relay to EntryPoint  
4. **Smart Contracts** – EntryPoint, SmartAccount, Factory deployed on Somnia Testnet

---

## ✨ Features

- 🔑 **Secure Key Management** – All private keys are generated & stored with **Turnkey**, never exposed to backend or frontend.
- 🌐 **OAuth Login** – Google, GitHub, Apple login via NextAuth on frontend.
- 📦 **UserOperation Builder** – Nonce, `initCode`, gas values are set in the signer before signing.
- 📝 **Transaction Tracking** – Bundler-signed UserOps are tracked and stored in DB for frontend.
- 🔒 **JWT Authentication** – All sensitive endpoints require JWT tokens.

---

## 📂 Project Structure

```plaintext
SOMNIA-SIGNER
│
├── db/               # Database connection & errors
├── ethclient/        # Somnia Testnet RPC client & ABIs
├── handler/          # HTTP route handlers
├── middleware/       # JWT middleware
├── models/           # Data models
├── signer/           # SmartSigner core logic
├── turnkey/          # Turnkey API client
├── types/            # Shared structs & types
├── utils/            # Utilities & Docker configs
├── main.go           # Application entrypoint
└── README.md
```

---

## 🛠 Setup & Installation

### Requirements
  
- Go 1.21+
- PostgreSQL
- Turnkey API credentials
- Somnia Testnet RPC endpoint

### Steps

1. **Clone the repository**
  
   ```bash
   git clone https://github.com/springmacedonio-sys/Somnia-DeFi-Wallet.git
   cd eolia-signer
   ```

2. **Set environment variables** in `.env`
  
   ```env
   DB_CONN_STR=postgres://user:password@localhost:5432/somnia
   JWT_SECRET=supersecret
   ```

3. **Install dependencies**
  
   ```bash
   go mod tidy
   ```

4. **Run with Docker (optional)**
  
   ```bash
   docker-compose up --build
   ```

5. **Run locally**
  
   ```bash
   go run main.go
   ```

---

## 🌐 API Endpoints

| Method | Endpoint              | Description                          | Auth |
|--------|-----------------------|---------------------------------------|------|
| POST   | `/auth`               | Login                                 | ❌   |
| POST   | `/auth/register`      | Register new account                  | ❌   |
| POST   | `/auth/logout`        | Logout                                | ✅   |
| GET    | `/auth/:wallet_name`  | Check if wallet name is available     | ❌   |
| GET    | `/me`                 | Get current user info                 | ✅   |
| POST   | `/sign`               | Sign a UserOperation                  | ✅   |
| POST   | `/addTx`              | Save tx hash to history               | ✅   |
| GET    | `/txHistory`          | Get transaction history               | ✅   |

✅ = Requires JWT authentication

---

## 🧩 Tech Stack

- **Go Fiber** – HTTP framework
- **PostgreSQL** – Database
- **Turnkey API** – Secure signing
- **Somnia Testnet** – EVM-compatible blockchain
- **ERC-4337** – Account Abstraction standard

---

## 📜 License

MIT
