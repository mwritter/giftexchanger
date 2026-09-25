"use client"

import { Button } from "@/components/ui/button";
import { Card, CardAction, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Calendar, Edit, Trash, User2, Wallet, X } from "lucide-react";
import { Exchange } from "./ExchangeList";

export function ExchangeListItems({ name, exchangeDate, exchangeBudgent, exchangeDescription }: Exchange) {
    return <li>
        <Card>
            <CardHeader>
                <CardTitle>{name}</CardTitle>
                <CardDescription>{exchangeDescription}</CardDescription>
                <CardAction>
                    <Button variant={'ghost'}>
                        <Edit size={12} />
                    </Button>
                    <Button variant={'destructive'}>
                        <Trash size={12} />
                    </Button>
                </CardAction>
            </CardHeader>
            <CardContent>
                <div className="flex items-center gap-2">
                    <User2 size={15} /> 2 participants
                </div>
                <div className="flex items-center gap-2">
                    <Calendar size={15} /> {exchangeDate.toDateString()}
                </div>
                <div className="flex items-center gap-2">
                    <Wallet size={15} /> ${exchangeBudgent} budget
                </div>
            </CardContent>
            <CardFooter>
                <p className="text-sm">Created by you</p>
            </CardFooter>
        </Card>
    </li>
}