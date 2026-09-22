import { Header } from "@/components/Header/Header"
import { LogoutButton } from "@/components/LogoutButton/LogoutButton"
import { requireUser } from "@/lib/session"

export default async function Dashboard() {
    const user = await requireUser()

    return <main className="flex flex-1 flex-col gap-6 p-8">
        <Header />
        <div className="flex flex-col gap-4 w-full max-w-120">
            <div className="space-y-1">
                <h1 className="text-2xl font-bold tracking-tight">
                    {user.display_name || "Your account"}
                </h1>
                <p className="text-sm text-muted-foreground">{user.email}</p>
            </div>
            <LogoutButton />
        </div>
    </main>
}
