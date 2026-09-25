import { ExchangeListItems } from "./ExchangeListItem"

export type Exchange = {
    id: string,
    name: string,
    exchangeDate: Date,
    exchangeDescription?: string,
    exchangeBudgent?: string
}

type Props = {
    exchanges: Exchange[]
}


export function ExchangeList({ exchanges }: Props) {
    return <ul>
        {exchanges.map(e => <ExchangeListItems key={e.id} {...e} />)}
    </ul>
}