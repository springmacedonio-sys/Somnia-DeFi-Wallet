'use client';

import { tokens } from '@/constants/tokens';
import { sendUserOp } from '@/lib/tx';
import { parseAmount } from '@/utils/formatting';
import { AnimatePresence, motion } from 'framer-motion';
import { X } from 'lucide-react';
import { useState } from 'react';

export default function TransferModal({
  sender,
  isOpen,
  onClose,
}: {
  sender: `0x${string}`;
  isOpen: boolean;
  onClose: () => void;
}) {
  const [recipient, setRecipient] = useState<`0x${string}`>('0x');
  const [amount, setAmount] = useState('');
  const [selectedToken, setSelectedToken] = useState(tokens[0]);

  const handleSend = async () => {
    if (!recipient || !amount) return;
    const parsedAmount = parseAmount(amount, selectedToken.decimals);
    await sendUserOp({
      sender,
      userOp: {
        target: recipient,
        value: String(parsedAmount),
        data: '0x',
      }
    });
    onClose();
  };

  return (
    <AnimatePresence>
      {isOpen && (
        <motion.div
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
        >
          <motion.div
            className="bg-gradient-to-br from-[#261849] to-[#432C77] rounded-3xl shadow-2xl w-[90%] max-w-sm px-6 py-7 text-white relative"
            initial={{ scale: 0.9, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            exit={{ scale: 0.9, opacity: 0 }}
          >
            {/* Close Button */}
            <button
              onClick={onClose}
              className="absolute top-4 right-4 text-white/60 hover:text-white"
            >
              <X size={22} />
            </button>

            <h2 className="text-center text-lg font-semibold mb-6">Send Assets</h2>

            <div className="space-y-5">
              {/* Token Selector */}
              <div>
                <label className="text-sm text-white/80">Token</label>
                <select
                  className="w-full mt-1 px-4 py-2 bg-white/10 text-white border border-purple-500 rounded-xl appearance-none focus:outline-none focus:ring-2 focus:ring-purple-400 relative"
                  value={selectedToken.symbol}
                  onChange={(e) => {
                    const token = tokens.find((t) => t.symbol === e.target.value);
                    if (token) setSelectedToken(token);
                  }}
                >
                  {tokens.map((token) => (
                    <option key={token.symbol} value={token.symbol} className="text-black">
                      {token.symbol}
                    </option>
                  ))}
                </select>
              </div>

              {/* Recipient */}
              <div>
                <label className="text-sm text-white/80">Recipient Address</label>
                <input
                  type="text"
                  placeholder="0x..."
                  value={recipient}
                  onChange={(e) => setRecipient(e.target.value as `0x${string}`)}
                  className="w-full mt-1 px-4 py-2 bg-white/10 border border-white/20 rounded-xl text-sm placeholder-white/40 focus:outline-none focus:ring-2 focus:ring-purple-400"
                />
              </div>

              {/* Amount */}
              <div>
                <label className="text-sm text-white/80">Amount</label>
                <input
                  type="text"
                  placeholder="0.00"
                  value={amount}
                  onChange={(e) => setAmount(e.target.value)}
                  className="w-full mt-1 px-4 py-2 bg-white/10 border border-white/20 rounded-xl text-sm placeholder-white/40 focus:outline-none focus:ring-2 focus:ring-purple-400"
                />
              </div>

              {/* Send Button */}
              <button
                onClick={handleSend}
                className="w-full bg-white text-[#5C2BCB] font-semibold py-2.5 rounded-xl shadow-sm hover:bg-white/90 transition"
              >
                Send
              </button>
            </div>
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  );
}