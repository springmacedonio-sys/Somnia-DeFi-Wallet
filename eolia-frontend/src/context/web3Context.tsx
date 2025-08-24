"use client"

import { createContext, useContext, useState } from "react";
import { getTokenBalance } from "@/lib/getBalance";
import { tokens } from "@/constants/tokens";
import { getPrice } from "@/lib/okx/client";
import { trimToDecimals } from "@/utils/decimals";
import { delay } from "@/utils/utils";

type TokenBalance = {
  symbol: string;
  name: string;
  logo: string;
  balance: number;
  usdValue: number;
};

type TxRecord = {
	dexName:    string;
	fromToken:  string;
	fromAmount: string;
	toToken:    string;
	toAmount:   string;
	timestamp:  string;
}

type Web3ContextType = {
  totalBalance: number;
  balances: TokenBalance[];
  txHistory: TxRecord[];
  fetchBalances: (address: `0x${string}`) => Promise<void>;
  fetchTxHistory: (address: `0x${string}`) => Promise<void>;
};

function formatAmount(amount: bigint, decimals: number): string {
  const s = amount.toString().padStart(decimals + 1, "0");
  const int = s.slice(0, -decimals);
  const frac = s.slice(-decimals).replace(/0+$/, "");
  return frac ? `${int}.${frac}` : int;
}

const Web3Context = createContext<Web3ContextType | undefined>(undefined);

export const Web3Provider = ({ children }: { children: React.ReactNode }) => {
  const [balances, setBalances] = useState<TokenBalance[]>([]);
  const [totalBalance, setTotalBalance] = useState<number>(0);
  const [txHistory, setTxHistory] = useState<TxRecord[]>([]);

  const fetchBalances = async (address: `0x${string}`) => {
    if (!address) return;

    let total: number = 0;

    const results: TokenBalance[] = [];

    for (const token of tokens) {
      try {
        const balance = await getTokenBalance(address, token.address as `0x${string}`, token.decimals);
        const price = await getPrice(token.address as `0x${string}`, token.decimals);

        total += balance * Number(price);
        results.push({
          symbol: token.symbol,
          name: token.name,
          logo: token.logo,
          balance,
          usdValue: balance * Number(price),
        });

        await delay(200);
      } catch (e) {
        console.error(`Error loading ${token.symbol}`, e);
        results.push({
          symbol: token.symbol,
          name: token.name,
          logo: token.logo,
          balance: 0,
          usdValue: 0,
        });
      }
    }

    setBalances(results);
    setTotalBalance(total);
  };

  const fetchTxHistory = async () => {
    try {;
      const resTx = await fetch('http://localhost:8080/txHistory', {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
      })

      const data = await resTx.json();

      const _txHistory = (data.txHistory as TxRecord[]).map((tx: any) => {
        const _fromAmount = tokens.find(token => token.symbol === tx.fromToken) ? trimToDecimals(Number(formatAmount(BigInt(tx.fromAmount), tokens.find(token => token.symbol === tx.fromToken)!.decimals)), 6) : tx.fromAmount;
        const _toAmount = tokens.find(token => token.symbol === tx.toToken) ? trimToDecimals(Number(formatAmount(BigInt(tx.toAmount), tokens.find(token => token.symbol === tx.toToken)!.decimals)), 6) : tx.toAmount;

        return {
          dexName: tx.dexName,
          fromToken: tx.fromToken,
          fromAmount: _fromAmount,
          toToken: tx.toToken,
          toAmount: _toAmount,
          timestamp: tx.timestamp,
        };
      });

      setTxHistory(_txHistory || []);
    } catch (err) {
      console.error("Error fetching tx history:", err);
      setTxHistory([]);
    }
  };

  return (
    <Web3Context.Provider value={{ balances, totalBalance, txHistory, fetchBalances, fetchTxHistory }}>
      {children}
    </Web3Context.Provider>
  );
};

export const useWeb3 = (): Web3ContextType => {
  const context = useContext(Web3Context);
  if (!context) {
    throw new Error("useWeb3 must be used within a Web3Provider");
  }
  return context;
};