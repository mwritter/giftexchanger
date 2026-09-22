"use client"

import { Button } from "@/components/ui/button"
import { logout } from "@/lib/auth"
import { useMutation } from "@tanstack/react-query"
import { useRouter } from "next/navigation"

export function LogoutButton() {
    const router = useRouter()

    const signOut = useMutation({
        mutationFn: logout,
        onSuccess: () => {
            router.replace("/")
        }
    })

    return <div className="flex flex-col items-start gap-2">
        <Button
            variant="outline"
            disabled={signOut.isPending}
            onClick={() => signOut.mutate()}
        >
            {signOut.isPending ? "Logging out…" : "Log out"}
        </Button>
        {signOut.isError && (
            <p className="text-sm text-destructive">{signOut.error.message}</p>
        )}
    </div>
}
