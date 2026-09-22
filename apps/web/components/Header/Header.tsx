import { Gift } from "lucide-react";

export function Header() {
    return <div className="flex gap-2 items-center">
        <Gift />
        <p className="text-2xl font-bold">GiftSwap</p>
    </div>
}