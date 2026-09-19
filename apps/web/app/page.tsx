import { Button } from "@/components/ui/button";

function Home() {
  return (
    <main className="flex flex-1 flex-col items-center justify-center gap-6 p-8">
      <div className="space-y-2 text-center">
        <h1 className="text-3xl font-semibold tracking-tight">GiftSwap</h1>
        <p className="text-muted-foreground">
          Self-hostable Secret Santa for friends and families.
        </p>
      </div>
      <div className="flex flex-wrap items-center justify-center gap-2">
        <Button>Create exchange</Button>
        <Button variant="outline">Join with invite</Button>
      </div>
    </main>
  );
}

export default Home;
