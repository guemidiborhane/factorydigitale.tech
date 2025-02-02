import { type ReactNode } from "preact/compat";
import Navbar from "ui/Navbar";
import { useDirection } from "~/hooks";
import Toaster from "ui/Toaster";

export default function RootLayout({ children }: { children?: ReactNode }) {
    return (
        <div dir={useDirection()}>
            <Navbar />

            {children}

            <Toaster />
        </div>
    )
}
