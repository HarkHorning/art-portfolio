<script lang="ts">
    import { page } from '$app/stores';

    interface PrintDetail {
        id: number;
        title: string;
        description: string;
        portrait: boolean;
        url: string;
        price_cents: number;
        size: string;
        sold: boolean;
        quantity_in_stock: number;
    }

    let print: PrintDetail | null = $state(null);
    let loading = $state(true);
    let error: string | null = $state(null);

    $effect(() => {
        const id = $page.params.id;
        loading = true;
        error = null;

        (async () => {
            try {
                const res = await fetch(`/api/v1/prints/${id}`);
                if (!res.ok) throw new Error();
                print = await res.json();
            } catch {
                error = 'Could not load this print.';
            } finally {
                loading = false;
            }
        })();
    });

    function formatPrice(cents: number): string {
        return `$${(cents / 100).toFixed(0)}`;
    }
</script>

<svelte:head>
    <title>{print ? `${print.title} — Hark Horning` : 'Hark Horning'}</title>
</svelte:head>

<div class="detail-page">
    <button onclick={() => history.back()} class="back">← Back</button>

    {#if loading}
        <p class="status">Loading...</p>
    {:else if error}
        <p class="status error">{error}</p>
    {:else if print}
        <div class="detail" class:portrait={print.portrait} class:landscape={!print.portrait}>
            <div class="image-wrap">
                <img src={print.url} alt={print.title} />
            </div>
            <div class="info">
                <h1>{print.title}</h1>
                <div class="print-meta">
                    <span class="size">{print.size}"</span>
                    {#if print.sold}
                        <span class="sold">Sold</span>
                    {:else if print.quantity_in_stock === 0}
                        <span class="sold">Out of stock</span>
                    {:else}
                        <span class="price">{formatPrice(print.price_cents)}</span>
                        <span class="stock">({print.quantity_in_stock} in stock)</span>
                    {/if}
                </div>
                {#if print.description}
                    <p class="description">{print.description}</p>
                {/if}
            </div>
        </div>
    {/if}
</div>

<style>
    .detail-page {
        width: 100%;
        max-width: 1100px;
    }

    .back {
        display: inline-block;
        background: none;
        border: none;
        padding: 0;
        cursor: pointer;
        color: #888;
        font-size: 0.85rem;
        font-family: inherit;
        margin-bottom: 2rem;
        transition: color 0.2s;
    }

    .back:hover {
        color: #000;
    }

    .detail.portrait {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 3rem;
        align-items: start;
    }

    .detail.landscape {
        display: grid;
        grid-template-columns: 1fr;
        gap: 1.5rem;
        max-width: 800px;
    }

    .image-wrap img {
        width: 100%;
        height: auto;
        border-radius: 8px;
        display: block;
    }

    .info {
        padding-top: 0.5rem;
    }

    h1 {
        font-size: 1.5rem;
        font-weight: 400;
        margin: 0 0 0.75rem;
    }

    .print-meta {
        display: flex;
        gap: 1rem;
        align-items: baseline;
        margin-bottom: 1.25rem;
    }

    .size {
        font-size: 0.85rem;
        color: #999;
    }

    .price {
        font-size: 1.1rem;
        font-weight: 500;
        color: #111;
    }

    .sold {
        font-size: 0.85rem;
        color: #999;
        font-style: italic;
    }

    .stock {
        font-size: 0.8rem;
        color: #aaa;
    }

    .description {
        color: #555;
        font-size: 0.9rem;
        line-height: 1.7;
        margin: 0;
    }

    .status {
        color: #666;
        font-style: italic;
    }

    .error {
        color: #c00;
    }

    @media (max-width: 650px) {
        .detail.portrait {
            grid-template-columns: 1fr;
        }
    }
</style>
