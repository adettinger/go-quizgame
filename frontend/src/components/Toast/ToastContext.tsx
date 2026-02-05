import { createContext, useContext, useState, type JSX, type ReactNode } from "react";
import * as RadixToast from '@radix-ui/react-toast';
import './ToastStyles.scss';

type ToastType = 'success' | 'error' | 'warning' | 'info';

interface ToastContextType {
    showToast: (type: ToastType, title: string, message: string | JSX.Element) => void;
}

const ToastContext = createContext<ToastContextType | undefined>(undefined);

export function ToastProvider({ children }: { children: ReactNode }) {
    const [open, setOpen] = useState(false);
    const [type, setType] = useState<ToastType>('success');
    const [title, setTitle] = useState('');
    const [message, setMessage] = useState<string | JSX.Element>('');

    const showToast = (toastType: ToastType, toastTitle: string, toastMessage: string | JSX.Element) => {
        setType(toastType);
        setTitle(toastTitle);
        setMessage(toastMessage);
        setOpen(true);
    };

    return (
        <ToastContext.Provider value={{ showToast }}>
            <RadixToast.Provider swipeDirection="right">
                {children}

                <RadixToast.Root
                    className={`ToastRoot ${type}`}
                    open={open}
                    onOpenChange={setOpen}
                    duration={5000}
                >
                    <RadixToast.Title className="ToastTitle">
                        {title}
                    </RadixToast.Title>
                    <RadixToast.Description className="ToastDescription">
                        {message}
                    </RadixToast.Description>
                    <RadixToast.Action className="ToastAction" asChild altText="Close">
                        <button className="ToastCloseButton">✕</button>
                    </RadixToast.Action>
                </RadixToast.Root>

                <RadixToast.Viewport className="ToastViewport" />
            </RadixToast.Provider>
        </ToastContext.Provider>
    );
}

export function useToast() {
    const context = useContext(ToastContext);
    if (context === undefined) {
        throw new Error('useToast must be used within a ToastProvider');
    }
    return context;
}