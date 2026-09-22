import { HomeHero } from "@/components/HomeHero/HomeHero";
import { getCurrentUser } from "@/lib/session";
import { redirect } from "next/navigation";

// This home page is for unauthenticated users
// If authenticated redirect to user dashboard

async function Home() {
  if (await getCurrentUser()) {
    redirect("/dashboard");
  }

  return <HomeHero />;
}

export default Home;
