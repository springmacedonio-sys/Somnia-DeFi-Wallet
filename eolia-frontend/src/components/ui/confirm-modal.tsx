'use client';
import { motion, AnimatePresence } from "framer-motion";
import { X } from "lucide-react";

const backdrop = {
  hidden: { opacity: 0 },
  visible: { opacity: 1 }
};

const modal = {
  hidden: { y: "-10%", opacity: 0 },
  visible: { y: "0%", opacity: 1 }
};

export default function ConfirmModal({
  isOpen,
  onClose,
  onConfirm,
  data
}: {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
  data: {
    title: string;
    description?: string;
    actions: {
      label: string;
      line: string;
    }[];
    confirmLabel: string;
  };
}) {
  return (
    <AnimatePresence>
      {isOpen && (
        <motion.div
          className="fixed top-0 left-0 w-full h-full bg-[#0000001f] backdrop-blur-sm z-50 flex items-center justify-center shadow-xl"
          variants={backdrop}
          initial="hidden"
          animate="visible"
          exit="hidden"
        >
          <motion.div
            className="bg-white p-6 rounded-2xl shadow-2xl w-[400px] text-center relative"
            variants={modal}
          >
            <button
              onClick={onClose}
              className="absolute cursor-pointer top-3 right-3 text-gray-400 hover:text-gray-600"
            >
              <X size={20} />
            </button>

            <p className="text-lg font-semibold mb-3 text-[#070417]">
              {data.title}
            </p>
            {data.description && (
              <p className="text-sm mb-10 text-[#8684a0]">{data.description}</p>
            )}

            <div className="text-left space-y-3 mb-6">
              {data.actions.map((action, i) => (
                <div
                  key={i}
                  className="w-full border-[2px] border-[#f0eff6] rounded-lg px-4 py-3 bg-transparent flex justify-between items-center"
                >
                  <div className="flex items-center gap-2">
                    <span className="text-sm font-semibold text-[#1a172a]">
                      {action.label}
                    </span>
                    <span className="text-sm font-light text-[#1a172a]">
                      {action.line}
                    </span>
                  </div>
                </div>
              ))}
            </div>

            <button
              onClick={onConfirm}
              className="w-full py-2 rounded-xl bg-[#965cef] text-white font-medium cursor-pointer transition"
            >
              {data.confirmLabel}
            </button>
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  );
}