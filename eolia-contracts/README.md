# Eolia Contracts

Smart Account & EntryPoint implementation for **Eolia Wallet**, deployed on Somnia Testnet for the hackathon project.

## 📂 Project Structure
  
```plaintext
eolia-contract/
├── contracts/               # Solidity smart contracts
│   ├── accounts/             # Account logic
│   ├── core/                 # Core interfaces
│   ├── interfaces/           # ERC-4337 + custom interfaces
│   └── utils/                 # Helper libraries
├── deploy/                   # hardhat-deploy scripts
│   ├── 0_deploy_entrypoint.ts
│   └── 1_deploy_factory.ts
├── deployments/somnia/       # Deployment artifacts
├── scripts/                  # Utility scripts (recover, create2, generate-userop)
├── test/                     # Contract tests
├── hardhat.config.ts         # Hardhat configuration
└── package.json              # Dependencies
```

## 🚀 Deployment

### 1. Install dependencies
  
```bash
pnpm install
```

### 2. Compile contracts
  
```bash
pnpm hardhat compile
```

### 3. Deploy EntryPoint & Factory
  
```bash
pnpm hardhat deploy --network somnia
```

### 4. Verify contracts
  
```bash
pnpm hardhat verify --network somnia <contractAddress> <constructorArgs...>
```

## 📜 Deployed Addresses (Somnia Testnet)
  
- **EntryPoint**: `0x...`
- **SimpleAccountFactory**: `0x...`

## 🛠 Helper Scripts
  
- `scripts/recoverAddress.ts` → Recover signer address from signature
- `scripts/create2account.ts` → Predict account address with CREATE2
- `scripts/generate-userop.ts` → Generate UserOperation struct for testing

## 🧩 Notes
  
- Compatible with ERC-4337
- Optimized for Somnia Testnet EVM (basefee opcode differences handled)
- Using `hardhat-deploy` for deterministic deployments

## 📄 License
  
MIT License
