"use client"

import { CreateExchangeForm } from "@/components/CreateExchangeForm/CreateExchangeForm";
import { DashboardPageContent } from "@/components/DashboardPageLayout/DashboardPageContent";
import { DashboardPageHeader } from "@/components/DashboardPageLayout/DashboardPageHeader";
import { DashboardPageLayout } from "@/components/DashboardPageLayout/DashboardPageLayout";
import { ExchangeList } from "@/components/ExchangeList/ExchangeList";
import { Button } from "@/components/ui/button";
import { Plus, X } from "lucide-react";
import { useState } from "react";

// List exchanges
// Create and update exchanges

export default function ExchangesPage() {
    const [showCreateNewForm, setShowCreateNewForm] = useState(false)

    return <DashboardPageLayout>
        <div className="flex items-center justify-between md:max-w-125 gap-5">
            <DashboardPageHeader title="Exchanges" description="This is a list of exchanges your current a part of!" />
            {!showCreateNewForm ? <Button onClick={() => setShowCreateNewForm(true)} size='sm'>
                <Plus />
                Create
            </Button> : <Button variant={'destructive'} onClick={() => setShowCreateNewForm(false)}><X /></Button>}
        </div>
        <DashboardPageContent>
            {showCreateNewForm ? <CreateExchangeForm /> :
                <ExchangePageList />}

        </DashboardPageContent>
    </DashboardPageLayout>
}


function ExchangePageList() {
    const exchanges = [{
        id: '123',
        name: 'Christmas 2026',
        exchangeDate: new Date('01/01/2026')
    }] as {
        id: string,
        name: string,
        exchangeDate: Date,
        exchangeDescription?: string,
        exchangeBudget?: number
    }[]

    if (exchanges.length === 0) return <EmptyExchangeList />

    return <ExchangeList exchanges={exchanges} />
}

function EmptyExchangeList() {
    return <div className="w-full text-center justify-center">
        <p className="text-sm">No exchanges to list yet</p>
    </div>
}