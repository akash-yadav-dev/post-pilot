"use client";

import { motion } from "framer-motion";

export default function GlobalLoader() {
  return (
    <div className="fixed inset-0 z-[9999] flex items-center justify-center bg-white">
      <motion.div
        initial={{ scale: 0.9, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        exit={{ opacity: 0 }}
        transition={{ duration: 0.4 }}
        className="relative"
      >
        {/* Logo Container */}
        <motion.div
          animate={{ rotate: 360 }}
          transition={{ repeat: Infinity, duration: 3, ease: "linear" }}
          className="w-24 h-24 rounded-2xl border-[5px] border-gray-200 shadow-lg bg-white flex items-center justify-center"
        >
          {/* Glow */}
          <motion.div
            animate={{ opacity: [0.3, 0.8, 0.3] }}
            transition={{ repeat: Infinity, duration: 1.5 }}
            className="absolute inset-0 rounded-2xl bg-gradient-to-tr from-gray-200 via-white to-gray-300 blur-xl"
          />

          {/* Logo Letter */}
          <span className="text-3xl font-bold text-gray-800 relative z-10 [font-family:var(--font-leckerli-one)]">
            P
          </span>
        </motion.div>
      </motion.div>
    </div>
  );
}