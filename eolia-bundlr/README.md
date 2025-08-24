# 🚀 Eolia Bundlr

**Eolia Bundlr** is the custom **ERC‑4337 bundler** for the **Eolia Smart Wallet**.  
It validates incoming **UserOperations**, runs lightweight simulations, and relays valid ops to the **EntryPoint** contract on **XLayer**. It also exposes a small HTTP API for the frontend/backend to submit ops and query basic status.

This repository is one of **4 core components** in the Eolia stack:

1. **Frontend** (Next.js) – UI & OKX DEX integration  
2. **Signer** (Go) – Account creation, Turnkey signing, UserOp building  
3. **Bundlr** (this repo) – Validation + relay to EntryPoint  
4. **Smart Contracts** – EntryPoint / SmartAccount / Factory deployed on **XLayer**

---

## ✨ Features

- 🧠 **UserOp validation & simulation** before relay (gas sanity, nonce/initCode presence, basic checks)
- 📤 **Relay to EntryPoint** on **XLayer** via public RPC
- 🧵 **OpQueue** for basic queuing / backpressure control
- 🔌 **HTTP RPC** for submitting ops from the Signer / Frontend
- 📊 **Minimal tracking** to help the UI follow operation status
- 🛡️ **CORS** enabled for local development (localhost:3000, 127.0.0.1:8080)

> ⚠️ Note: This is a hackathon‑grade bundler focused on clarity. For production, add robust simulation, reputation, mempool logic, and rate limiting.

---

## 📂 Project Structure

```plaintext
EOLIA-BUNDLR
│
├── cmd/bundlr/        # Application entrypoint (main.go)
├── config/            # Config loader & config.yaml
├── internal/
│   ├── bundlr/        # Bundler core (loop, queue)
│   ├── rpc/           # HTTP router & handlers
│   ├── signer/        # (Helpers if bundler needs local signing)
│   └── validator/     # Validation logic + EntryPoint ABI
│       └── entrypoint/entrypoint.abi.json
├── types/             # Shared structs (UserOperation, etc.)
├── go.mod
├── go.sum
└── README.md
```

---

## ⚙️ Configuration

All runtime settings live in **`config/config.yaml`**.

```yaml
# config/config.yaml
chain_id: 196                       # XLayer Chain ID
rpc_url: "https://rpc.xlayer.tech"  # Public RPC URL for XLayer

# Deployed contract addresses on XLayer
entry_point: "0x379FF91b96c038ECb0dc6aCFb44366a39f0de566"   # EntryPoint
factory:     "0xC924da88e33fD1eD04f4A8a1f6BD14Ad030a3dC9"   # Account Factory

# Bundler sender (EOA) that pays for transactions
bundlr_address:    "YourAddressHere"     # Must have balance on XLayer
bundlr_private_key: "YourPrivateKeyHere" # For dev only — prefer env/VAULT in prod
```

> 🔒 **Security tip:** Avoid committing real private keys. Prefer environment variables or a KMS/Turnkey‑style signer in production.

---

## 🧰 Build & Run

### Local (Go)

```bash
# from repo root
go run ./cmd/bundlr
```

By default the server listens on **`:8181`** (see `cmd/bundlr/main.go`).

### Environment

- Go **1.21+**
- A working **XLayer** RPC endpoint
- Deployed **EntryPoint** & **Factory** addresses on XLayer
- A funded `bundlr_address` for paying tx fees

---

## 🌐 HTTP API (dev)

> The router is defined under `internal/rpc`. Endpoints may evolve; check the source if you change routes.

| Method | Path                    | Purpose                                    |
|-------:|-------------------------|--------------------------------------------|
| POST   | `/rpc/sendUserOp`       | Submit a **signed UserOperation**          |
| GET    | `/rpc/getUserOpReceipt` | Get basic status for a userOp/tx (if any)  |
| GET    | `/rpc/getChainId`       | Get the chain ID that bundlr's working on  |

### Example: submit a signed UserOperation

```bash
curl -X POST http://localhost:8181/bundle   -H "Content-Type: application/json"   -d '{
    "jsonrpc": "2.0",
    "method": "eth_sendUserOperation",
    "params": [
      {
        "sender": "0x...",
        "nonce": "0x...",
        "initCode": "0x...",
        "callData": "0x...",
        "accountGasLimits": "0x...",
        "preVerificationGas": "0x...",
        "gasFees": "0x...",
        "paymasterAndData": "0x",
        "signature": "0x..." 
      }
    ],
    "id" 1
  }'
```

---

## 🔄 Runtime Flow (High Level)

1. **Receive** signed UserOperation via HTTP (`/rpc/sendUserOp`)  
2. **Validate / simulate** the op (nonce/initCode presence, gas sanity)  
3. **Enqueue** into **OpQueue** for rate control  
4. **Relay** to **EntryPoint** on XLayer via `rpc_url`  
5. **Track** basic status so the UI can poll if needed

A simplified version of the loop is in `internal/bundlr/bundlr.go` and `internal/bundlr/opqueue.go`.

---

## 🧪 Dev Notes

- CORS is enabled for: `http://localhost:3000`, `http://127.0.0.1:8080`  
- Chain & contract addresses are exported to `types` at startup for easy access:
  - `types.ChainID = bundlr.ChainID`
  - `types.EntryPoint = bundlr.Validator.EntryPoint`

(See `cmd/bundlr/main.go` for the bootstrap sequence.)

---

## 🧩 Tech Stack

- **Go / Fiber** – HTTP server
- **XLayer** – EVM chain (Chain ID **196**)
- **ERC‑4337** – Account Abstraction
- **Custom Validator + OpQueue**

---

## 📜 License

MIT
