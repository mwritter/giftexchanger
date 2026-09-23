import { LoginForm } from "@/components/LoginForm/LoginForm"
import { getCurrentUser } from "@/lib/session"
import { Gift } from "lucide-react"
import Link from "next/link"
import { redirect } from "next/navigation"

export default async function Login() {
    if (await getCurrentUser()) {
        redirect("/dashboard")
    }

    return <div className="flex flex-1 p-10">
        <Link href={"/"} className="fixed flex items-center gap-2 top-10 left-10 font-bold text-2xl max-sm:hidden">
            <Gift /> GiftExchanger
        </Link>
        <div className="flex flex-col justify-center items-center flex-1">
            <div className="flex flex-col items-center w-full max-w-120">
                <Link href={"/"} className="flex gap-2 items-center font-bold text-xl sm:hidden max-sm:mb-10">
                    <Gift /> GiftExchanger
                </Link>
                <LoginForm />
            </div>
        </div>
    </div>
}
