export interface PrintSizeInter {
    id: number;
    print_id: number;
    size: string;
    price_cents: number;
    quantity_in_stock: number;
    sold: boolean;
}

export interface PrintTileInter {
    id: number;
    title: string;
    url: string;
    portrait: boolean;
    sizes: PrintSizeInter[];
}
