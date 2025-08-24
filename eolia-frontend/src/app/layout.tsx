import type { Metadata } from "next";
import "./globals.css";
import { Inter } from "next/font/google";
import { sfPro } from "./fonts/fonts";
import { ClientLayout } from "@/components/client-layout";
import { Web3Provider } from "@/context/web3Context";

const inter = Inter({
  subsets: ['latin'],
  variable: '--font-inter',
  display: 'swap',
})

export const metadata: Metadata = {
  title: "Eolia Wallet",
  description: "Eolia Wallet is a self-custodial wallet that gives you full control over your digital assets and data.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <ClientLayout>
        <body
          className={`${inter.variable} ${sfPro.variable} antialiased`}
        >
          <Web3Provider>
            {children}
          </Web3Provider>
        </body>
      </ClientLayout>
    </html>
  );
}
