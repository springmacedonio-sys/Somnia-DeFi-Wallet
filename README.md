# 🚀 Somnia DeFi Wallet

**Somnia DeFi Wallet** is a smart wallet project based on Account Abstraction, built specifically for the **Somnia DeFi Mini Hackathon**.

This is a complete DeFi wallet ecosystem with 4 core microservice components, each focusing on specific functional areas, working together to provide users with secure and convenient decentralized financial services.

---

## 🏗️ Project Architecture

Somnia DeFi Wallet adopts a microservices architecture design with the following 4 core components:

### 1. 🌐 [Frontend Service](./eolia-frontend/) - User Interface
- **Tech Stack**: Next.js 15+, TypeScript, Tailwind CSS
- **Features**: OAuth login, wallet management, token swapping, transaction history
- **Highlights**: Modern UI design, responsive layout, smooth user experience

### 2. 🔐 [Signer Service](./eolia-signer/) - Secure Signing Service
- **Tech Stack**: Go, Fiber, PostgreSQL, **Turnkey API Integration**
- **Features**: Account creation, key management, UserOperation signing, transaction tracking
- **Highlights**: Enterprise-grade security powered by **Turnkey**, support for multiple OAuth providers

### 3. 📦 [Bundlr Service](./eolia-bundlr/) - ERC-4337 Bundler
- **Tech Stack**: Go, Fiber, Ethereum
- **Features**: UserOperation validation, simulation execution, batch bundling, on-chain relay
- **Highlights**: High-performance validation engine, optimized for Somnia testnet

### 4. 📜 [Smart Contracts](./eolia-contracts/) - Smart Contracts
- **Tech Stack**: Solidity, Hardhat, ERC-4337
- **Features**: Smart accounts, EntryPoint, factory contracts
- **Highlights**: Fully compatible with ERC-4337 standard, deployed on Somnia testnet

---

## 🚀 Quick Start

### Prerequisites
- Node.js 18+
- Go 1.21+
- PostgreSQL
- Docker (optional)

### 1. Clone the Project
```bash
git clone https://github.com/springmacedonio-sys/Somnia-DeFi-Wallet.git
cd Somnia-DeFi-Wallet
```

### 2. Start Services
Each microservice has its own README and startup instructions. It's recommended to start them in the following order:

1. **Deploy Smart Contracts** - See [Smart Contracts README](./eolia-contracts/)
2. **Start Signer Service** - See [Signer Service README](./eolia-signer/)
3. **Start Bundlr Service** - See [Bundlr Service README](./eolia-bundlr/)
4. **Start Frontend Application** - See [Frontend Service README](./eolia-frontend/)

### 3. Environment Configuration
Each service requires corresponding environment variable configuration. For detailed instructions, please refer to each service's README file.

---

## 🔧 Development Guide

### Local Development Environment
- All services support local development mode
- Complete Docker configuration included
- Hot reload and debugging support

### API Documentation
- **Signer Service**: RESTful API with JWT authentication, **Turnkey integration**
- **Bundlr Service**: HTTP RPC interface, compatible with ERC-4337 standard
- Detailed API documentation available in each service's README

### Testing
- Smart contracts include complete test suites
- Frontend includes component testing
- Backend services include unit tests and integration tests

---

## 🌟 Core Features

- **🔐 Account Abstraction**: Fully compatible with ERC-4337 standard
- **🛡️ Enterprise Security**: **Turnkey-powered key management system**
- **⚡ High Performance**: Go backend, high concurrency support
- **🎨 Modern UI**: Next.js frontend, responsive design
- **🔗 Multi-chain Support**: Optimized for Somnia network

---

## 🔑 Turnkey Integration

**Somnia DeFi Wallet** leverages **Turnkey** for enterprise-grade security and key management:

- **Secure Key Generation**: All private keys are generated and stored securely through Turnkey
- **Hardware Security**: Integration with hardware security modules (HSMs)
- **Audit Trail**: Complete audit trail for all cryptographic operations
- **Compliance**: Enterprise-grade compliance and security standards
- **Scalability**: Designed to handle high-volume operations securely


---

## 📚 Technical Documentation

- [ERC-4337 Account Abstraction Standard](https://eips.ethereum.org/EIPS/eip-4337)
- [Somnia Network Documentation](https://docs.somnia.network/)
- [Next.js Documentation](https://nextjs.org/docs)
- [Go Fiber Documentation](https://docs.gofiber.io/)
- [Turnkey API Documentation](https://docs.turnkey.com/)

---

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

---

## 🙏 Acknowledgments

Special thanks to the Somnia network team for technical support and infrastructure to this project.

---

*Built with ❤️ for the Somnia DeFi community*
