/*
  verifyDeploy.js - Deployed contracts sanity checker (Somnia Testnet)

  中文说明（简要）:
  - 核对 EntryPoint 与 SimpleAccountFactory 的 senderCreator 是否一致
  - 解析 owner/salt，复算反事实地址并与期望 sender 对比
  - 构造 initCode 并用 getSenderAddress 验证推导 sender
  - 检查 accountImplementation 是否有代码、其 entryPoint() 是否指向正确地址
  - 打印 EntryPoint 存款与原生余额、以及 gas 估算（若提供）

  English (brief):
  - Verify EntryPoint and SimpleAccountFactory share the same senderCreator
  - Recompute counterfactual address from owner/salt and compare with expected sender
  - Build initCode and validate derived sender via getSenderAddress
  - Check accountImplementation has code and its entryPoint() returns the configured EP
  - Print EntryPoint deposit and native balance, plus gas estimate (if provided)

  Usage:
    node scripts/verifyDeploy.js --owner 0x... --salt 0 --sender 0x... [--rpc https://...] [--ep 0x...] [--factory 0x...]

  Defaults:
    - rpc:    deployments (Somnia) RPC https://dream-rpc.somnia.network
    - ep:     deployments/somnia/EntryPoint.json
    - factory:deployments/somnia/SimpleAccountFactory.json
*/

const fs = require("fs");
const path = require("path");
const { ethers } = require("ethers");

function getArg(name, def) {
  const i = process.argv.indexOf(`--${name}`);
  if (i !== -1 && i + 1 < process.argv.length) return process.argv[i + 1];
  return def;
}

async function main() {
  const owner = getArg("owner");
  const saltStr = getArg("salt", "0");
  const expectedSender = getArg("sender");
  const rpc = getArg("rpc", "https://dream-rpc.somnia.network");

  // Read deployed addresses
  const epDeploy = JSON.parse(
    fs.readFileSync(path.join(__dirname, "../deployments/somnia/EntryPoint.json"), "utf8")
  );
  const factoryDeploy = JSON.parse(
    fs.readFileSync(path.join(__dirname, "../deployments/somnia/SimpleAccountFactory.json"), "utf8")
  );

  const epAddr = getArg("ep", epDeploy.address);
  const factoryAddr = getArg("factory", factoryDeploy.address);

  if (!owner || !expectedSender) {
    console.error("[!] Missing required args --owner and --sender");
    process.exit(1);
  }

  const salt = ethers.toBigInt(saltStr);
  const provider = new ethers.JsonRpcProvider(rpc);

  console.log("Context / 上下文:");
  console.log("- RPC:", rpc);
  console.log("- EntryPoint:", epAddr);
  console.log("- Factory:", factoryAddr);
  console.log("- Owner:", owner);
  console.log("- Salt:", salt.toString());
  console.log("- Expected Sender:", expectedSender);

  // Minimal ABIs
  const factoryAbi = [
    "function senderCreator() view returns (address)",
    "function computeAccountAddress(address owner, uint256 salt) view returns (address)",
    "function accountImplementation() view returns (address)",
    "function createAccount(address owner, uint256 salt) returns (address)",
  ];
  const epAbi = [
    "function senderCreator() view returns (address)",
    "function getSenderAddress(bytes initCode) external",
    "function balanceOf(address) view returns (uint256)",
  ];
  const accountAbi = [
    "function entryPoint() view returns (address)",
    "function owner() view returns (address)",
  ];

  const factory = new ethers.Contract(factoryAddr, factoryAbi, provider);
  const ep = new ethers.Contract(epAddr, epAbi, provider);

  // 1) senderCreator alignment
  const fsc = await factory.senderCreator();
  const esc = await ep.senderCreator();
  console.log("senderCreator -> factory:", fsc, "ep:", esc);
  if (fsc.toLowerCase() !== esc.toLowerCase()) {
    console.error("[X] senderCreator mismatch. 工厂与 EP 的 senderCreator 不一致");
    process.exit(2);
  }

  // 2) recompute address via factory
  const recomputed = await factory.computeAccountAddress(owner, salt);
  console.log("computeAccountAddress(owner,salt) ->", recomputed);
  if (recomputed.toLowerCase() !== expectedSender.toLowerCase()) {
    console.error("[X] recomputed sender != expected sender. 复算地址与期望不一致");
  } else {
    console.log("[OK] recomputed matches expected sender. 复算地址一致");
  }

  // 3) initCode = factory || encode(createAccount(owner,salt))
  const iface = new ethers.Interface(factoryAbi);
  const encoded = iface.encodeFunctionData("createAccount", [owner, salt]);
  const initCode = factoryAddr + encoded.slice(2);
  console.log("initCode:", initCode);

  // getSenderAddress(initCode)
  try {
    await ep.getSenderAddress(initCode);
    console.warn("getSenderAddress: unexpectedly not reverted (should always revert)");
  } catch (e) {
    const data = e?.data || e?.error?.data || e?.info?.error?.data;
    const derived = data ? "0x" + data.slice(-40) : "";
    console.log("getSenderAddress(initCode) ->", derived || "<no revert data>");
    if (derived && derived.toLowerCase() !== expectedSender.toLowerCase()) {
      console.error("[X] derived sender != expected sender. 推导地址与期望不一致");
    } else if (derived) {
      console.log("[OK] derived matches expected sender. 推导地址一致");
    }
  }

  // 4) accountImplementation checks
  const impl = await factory.accountImplementation();
  const implCode = await provider.getCode(impl);
  console.log("accountImplementation:", impl, "code length:", implCode === "0x" ? 0 : (implCode.length - 2) / 2);

  // try reading entryPoint() from implementation
  try {
    const accountImpl = new ethers.Contract(impl, accountAbi, provider);
    const epFromImpl = await accountImpl.entryPoint();
    console.log("implementation.entryPoint() ->", epFromImpl);
    if (epFromImpl.toLowerCase() !== epAddr.toLowerCase()) {
      console.error("[X] implementation.entryPoint() != configured EP. 实现合约的 EntryPoint 不一致");
    } else {
      console.log("[OK] implementation.entryPoint matches EP");
    }
  } catch (e) {
    console.warn("[!] Failed to read implementation.entryPoint():", e.message);
  }

  // 5) deployment status of expected sender
  const senderCode = await provider.getCode(expectedSender);
  console.log("expected sender code length:", senderCode === "0x" ? 0 : (senderCode.length - 2) / 2);

  // 6) EP deposit & native balance
  try {
    const dep = await ep.balanceOf(expectedSender);
    const bal = await provider.getBalance(expectedSender);
    console.log("EntryPoint.deposit:", dep.toString(), "native balance:", bal.toString());
  } catch (e) {
    console.warn("[!] Failed to read EP deposit or native balance:", e.message);
  }

  console.log("\nDone.");
}

main().catch((e) => {
  console.error("Fatal:", e);
  process.exit(1);
});
