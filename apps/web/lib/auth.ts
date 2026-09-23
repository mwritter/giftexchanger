export type User = {
  id: string
  email: string
  display_name: string
}

export const SESSION_COOKIE_NAME = "giftexchanger_session"

async function errorMessage(response: Response, fallback: string) {
  const data = (await response.json().catch(() => null)) as {
    error?: string
  } | null

  return data?.error ?? fallback
}

export async function requestMagicLink(email: string): Promise<void> {
  const response = await fetch("/api/auth/magic-link", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email }),
  })

  if (!response.ok) {
    throw new Error(
      await errorMessage(response, "We could not send your login link."),
    )
  }
}

export async function logout(): Promise<void> {
  const response = await fetch("/api/auth/logout", { method: "POST" })

  if (!response.ok) {
    throw new Error(await errorMessage(response, "We could not log you out."))
  }
}
