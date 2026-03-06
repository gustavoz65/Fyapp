"use client";

import { createContext, useContext, useState, ReactNode } from "react";

interface HealthModalContextType {
  isOpen: boolean;
  openModal: () => void;
  closeModal: () => void;
}

const HealthModalContext = createContext<HealthModalContextType | null>(null);

export function HealthModalProvider({ children }: { children: ReactNode }) {
  const [isOpen, setIsOpen] = useState(false);

  const openModal = () => setIsOpen(true);
  const closeModal = () => setIsOpen(false);

  return (
    <HealthModalContext.Provider value={{ isOpen, openModal, closeModal }}>
      {children}
    </HealthModalContext.Provider>
  );
}

export function useHealthModal() {
  const context = useContext(HealthModalContext);
  if (!context) {
    throw new Error("useHealthModal must be used within HealthModalProvider");
  }
  return context;
}
