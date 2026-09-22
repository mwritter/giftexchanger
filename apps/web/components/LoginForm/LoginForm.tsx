"use client"

import { Button } from "@/components/ui/button"
import { Field, FieldDescription, FieldError } from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { requestMagicLink } from "@/lib/auth"
import { useForm } from "@tanstack/react-form"
import { useMutation } from "@tanstack/react-query"
import z from "zod"
import { useState } from "react"

const formSchema = z.object({
    email: z.email()
})

export function LoginForm() {
    const [sentTo, setSentTo] = useState<string | null>(null)

    const sendLink = useMutation({
        mutationFn: (email: string) => requestMagicLink(email),
        onSuccess: (_, email) => setSentTo(email)
    })

    const { handleSubmit, Field: FromField } = useForm({
        defaultValues: {
            email: ""
        },
        validators: {
            onSubmit: formSchema
        },
        onSubmit: ({ value }) => {
            sendLink.mutate(value.email)
        }
    })

    if (sentTo) {
        return <div className="text-center text-sm space-y-2">
            <p>Your login link has been sent! Check {sentTo} for login instructions.</p>
            {process.env.NODE_ENV === "development" && (
                <p className="text-muted-foreground">
                    Running locally? The link is printed in the API terminal.
                </p>
            )}
        </div>
    }

    return <form
        className="flex flex-col gap-4 w-full"
        id="login-form"
        onSubmit={(e) => {
            e.preventDefault()
            handleSubmit()
        }}
    >
        <FromField name="email">
            {(field) => {
                const isInvalid =
                    field.state.meta.isTouched && !field.state.meta.isValid
                return (
                    <Field data-invalid={isInvalid}>
                        <Label htmlFor={field.name}>
                            Email
                        </Label>
                        <Input
                            id={field.name}
                            name={field.name}
                            value={field.state.value}
                            onBlur={field.handleBlur}
                            onChange={e => field.handleChange(e.target.value)}
                            aria-invalid={isInvalid}
                            placeholder="Email"
                            autoComplete="off"
                            disabled={sendLink.isPending}
                        />
                        <FieldDescription>
                            We&apos;ll send you a login link.
                        </FieldDescription>
                        {isInvalid && (
                            <FieldError errors={field.state.meta.errors} />
                        )}
                    </Field>
                )
            }}
        </FromField>
        <Field>
            <Button type="submit" disabled={sendLink.isPending}>
                {sendLink.isPending ? "Sending link…" : "Login"}
            </Button>
            {sendLink.isError && (
                <FieldError>{sendLink.error.message}</FieldError>
            )}
        </Field>
    </form>
}
