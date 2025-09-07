/**
 * verify-somnia-contracts.js
 * 
 * 专门用于在 Somnia Shannon Explorer 上验证合约的脚本
 * Script specifically for verifying contracts on Somnia Shannon Explorer
 * 
 * Usage:
 * node scripts/verify-somnia-contracts.js
 */

const { ethers } = require("ethers");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("🔍 开始验证 Somnia 测试网合约 / Starting Somnia testnet contract verification");
  
  // 读取部署信息 / Read deployment info
  const entryPointDeploy = JSON.parse(
    fs.readFileSync(path.join(__dirname, "../deployments/somnia/EntryPoint.json"), "utf8")
  );
  
  const factoryDeploy = JSON.parse(
    fs.readFileSync(path.join(__dirname, "../deployments/somnia/SimpleAccountFactory.json"), "utf8")
  );

  console.log("\n📋 部署地址信息 / Deployment Addresses:");
  console.log(`EntryPoint: ${entryPointDeploy.address}`);
  console.log(`SimpleAccountFactory: ${factoryDeploy.address}`);

  // 验证 EntryPoint 合约 / Verify EntryPoint contract
  console.log("\n🚀 验证 EntryPoint 合约 / Verifying EntryPoint contract...");
  try {
    const entryPointVerifyCmd = `npx hardhat verify --network somnia ${entryPointDeploy.address}`;
    console.log(`执行命令 / Executing command: ${entryPointVerifyCmd}`);
    
    // 由于 EntryPoint 构造函数没有参数，直接验证
    // Since EntryPoint constructor has no parameters, verify directly
    console.log("请手动执行以下命令 / Please manually execute the following command:");
    console.log(`cd eolia-contracts && ${entryPointVerifyCmd}`);
  } catch (error) {
    console.error("❌ EntryPoint 验证失败 / EntryPoint verification failed:", error.message);
  }

  // 验证 SimpleAccountFactory 合约 / Verify SimpleAccountFactory contract
  console.log("\n🚀 验证 SimpleAccountFactory 合约 / Verifying SimpleAccountFactory contract...");
  try {
    // SimpleAccountFactory 构造函数需要 EntryPoint 地址作为参数
    // SimpleAccountFactory constructor needs EntryPoint address as parameter
    const factoryVerifyCmd = `npx hardhat verify --network somnia ${factoryDeploy.address} "${entryPointDeploy.address}"`;
    console.log(`执行命令 / Executing command: ${factoryVerifyCmd}`);
    
    console.log("请手动执行以下命令 / Please manually execute the following command:");
    console.log(`cd eolia-contracts && ${factoryVerifyCmd}`);
  } catch (error) {
    console.error("❌ SimpleAccountFactory 验证失败 / SimpleAccountFactory verification failed:", error.message);
  }

  console.log("\n📝 验证说明 / Verification Notes:");
  console.log("1. 确保你的 .env 文件包含正确的私钥 / Ensure your .env file contains the correct private key");
  console.log("2. 如果验证失败，可能需要等待合约索引完成 / If verification fails, you may need to wait for contract indexing");
  console.log("3. 验证成功后，Shannon Explorer 应该能正常显示合约 / After successful verification, Shannon Explorer should display contracts normally");

  // 检查合约是否已部署 / Check if contracts are deployed
  console.log("\n🔍 检查合约部署状态 / Checking contract deployment status...");
  const provider = new ethers.JsonRpcProvider("https://dream-rpc.somnia.network");
  
  try {
    const entryPointCode = await provider.getCode(entryPointDeploy.address);
    const factoryCode = await provider.getCode(factoryDeploy.address);
    
    console.log(`EntryPoint 合约代码长度 / EntryPoint code length: ${entryPointCode === "0x" ? 0 : (entryPointCode.length - 2) / 2} bytes`);
    console.log(`Factory 合约代码长度 / Factory code length: ${factoryCode === "0x" ? 0 : (factoryCode.length - 2) / 2} bytes`);
    
    if (entryPointCode === "0x" || factoryCode === "0x") {
      console.error("❌ 合约未正确部署 / Contracts not properly deployed");
    } else {
      console.log("✅ 合约已成功部署 / Contracts successfully deployed");
    }
  } catch (error) {
    console.error("❌ 检查合约状态失败 / Failed to check contract status:", error.message);
  }

  console.log("\n🌐 浏览器链接 / Explorer Links:");
  console.log(`EntryPoint: https://shannon-explorer.somnia.network/address/${entryPointDeploy.address}`);
  console.log(`SimpleAccountFactory: https://shannon-explorer.somnia.network/address/${factoryDeploy.address}`);
}

main().catch((error) => {
  console.error("❌ 脚本执行失败 / Script execution failed:", error);
  process.exit(1);
});
