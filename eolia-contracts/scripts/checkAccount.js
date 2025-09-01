import { ethers } from "ethers";

// RPC 地址（Somnia Testnet 或你当前网络）
const RPC_URL = "https://dream-rpc.somnia.network"; 
const provider = new ethers.JsonRpcProvider(RPC_URL);

// 需要检查的账户地址
const accountAddress = "0xAcDa50f052a14B1F845275E107592c1fb2864676";

async function checkAccountDeployment() {
  try {
    // 查询链上 Bytecode
    const code = await provider.getCode(accountAddress);
    
    if (code === "0x") {
      console.log("账户尚未部署，Bytecode length = 0");
    } else {
      console.log("账户已部署，Bytecode length =", code.length / 2 - 1);
    }
  } catch (err) {
    console.error("查询失败:", err);
  }
}

checkAccountDeployment();
