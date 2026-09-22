import { LoginForm } from "@/components/LoginForm/LoginForm"
import { getCurrentUser } from "@/lib/session"
import { redirect } from "next/navigation"

export default async function Login() {
    if (await getCurrentUser()) {
        redirect("/dashboard")
    }

    return <LoginForm />
}
