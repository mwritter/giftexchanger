"use client"

import { LoginForm } from "@/components/LoginForm/LoginForm";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from "@/components/ui/dialog";
import { useRouter } from "next/navigation";

export default function Page() {
    const { back } = useRouter()
    return <Dialog open onOpenChange={back}>
        <DialogContent>
            <DialogHeader>
                <DialogTitle>Login with your email</DialogTitle>
                <DialogDescription>
                    We use magic links to login.  We send you a email that will give you instructions on login in.
                </DialogDescription>
            </DialogHeader>
            <LoginForm />
        </DialogContent>
    </Dialog>
}