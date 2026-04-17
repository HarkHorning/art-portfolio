<script lang="ts">
    import { page } from '$app/stores';
    import type { PrintSizeInter } from '$lib/components/printTile/PrintTileInterface';

    interface ImageInter {
        id: number;
        variant: string;
        url: string;
        sort_order: number;
    }

    interface PrintDetail {
        id: number;
        title: string;
        description: string;
        portrait: boolean;
        url: string;
        sizes: PrintSizeInter[];
        images: ImageInter[];
    }

    let print: PrintDetail | null = $state(null);
    let loading = $state(true);
    let error: string | null = $state(null);
    let selectedSize: PrintSizeInter | null = $state(null);
    let selectedIndex = $state(0);

    let displayImages = $derived(
        print ? print.images.filter(img => img.variant === 'low').sort((a, b) => a.sort_order - b.sort_order) : []
    );

    let currentImage = $derived(displayImages[selectedIndex] ?? null);

    $effect(() => {
        const id = $page.params.id;
        loading = true;
        error = null;
        selectedIndex = 0;

        (async () => {
            try {
                const res = await fetch(`/api/v1/prints/${id}`);
                if (!res.ok) throw new Error();
                print = await res.json();
                if (print) {
                    selectedSize = print.sizes.find(s => !s.sold && s.quantity_in_stock > 0) ?? print.sizes[0] ?? null;
                }
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

    function isAvailable(s: PrintSizeInter): boolean {
        return !s.sold && s.quantity_in_stock > 0;
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
        <div class="detail">
            <div class="image-wrap">
                {#if currentImage}
                    <img src={currentImage.url} alt={print.title} />
                {:else}
                    <img src={print.url} alt={print.title} />
                {/if}
                {#if displayImages.length > 1}
                    <div class="thumbnails">
                        {#each displayImages as img, i (img.id)}
                            <button
                                class="thumb"
                                class:active={i === selectedIndex}
                                onclick={() => selectedIndex = i}
                                aria-label="View image {i + 1}"
                            >
                                <img src={img.url} alt="" />
                            </button>
                        {/each}
                    </div>
                {/if}
            </div>
            <div class="info">
                <h1>{print.title}</h1>

                {#if print.sizes.length > 0}
                    <div class="size-selector">
                        {#each print.sizes as size}
                            <button
                                class="size-btn"
                                class:selected={selectedSize?.id === size.id}
                                class:unavailable={!isAvailable(size)}
                                onclick={() => selectedSize = size}
                            >
                                {size.size}"
                            </button>
                        {/each}
                    </div>

                    {#if selectedSize}
                        <div class="price-row">
                            {#if !isAvailable(selectedSize)}
                                <span class="sold">
                                    {selectedSize.sold ? 'Sold' : 'Out of stock'}
                                </span>
                            {:else}
                                <span class="price">{formatPrice(selectedSize.price_cents)}</span>
                                <span class="stock">({selectedSize.quantity_in_stock} in stock)</span>
                            {/if}
                        </div>
                    {/if}
                {/if}

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

    .back:hover { color: #000; }

    .detail {
        display: grid;
        grid-template-columns: 1fr 320px;
        gap: 3rem;
        align-items: start;
    }

    .image-wrap img {
        width: 100%;
        height: auto;
        border-radius: 8px;
        display: block;
    }

    .thumbnails {
        display: flex;
        gap: 0.5rem;
        margin-top: 0.75rem;
        flex-wrap: wrap;
    }

    .thumb {
        width: 60px;
        height: 60px;
        padding: 0;
        border: 2px solid transparent;
        border-radius: 4px;
        cursor: pointer;
        overflow: hidden;
        background: none;
        transition: border-color 0.15s;
    }

    .thumb img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        border-radius: 2px;
    }

    .thumb.active { border-color: #111; }

    .info { padding-top: 0.5rem; }

    h1 {
        font-size: 1.5rem;
        font-weight: 400;
        margin: 0 0 1.25rem;
    }

    .size-selector {
        display: flex;
        flex-wrap: wrap;
        gap: 0.5rem;
        margin-bottom: 1rem;
    }

    .size-btn {
        padding: 0.35rem 0.75rem;
        border: 1px solid #ccc;
        border-radius: 4px;
        background: #fff;
        cursor: pointer;
        font-family: inherit;
        font-size: 0.85rem;
        color: #333;
        transition: border-color 0.15s, background 0.15s;
    }

    .size-btn:hover { border-color: #999; }
    .size-btn.selected { border-color: #111; background: #111; color: #fff; }
    .size-btn.unavailable { opacity: 0.4; cursor: default; }

    .price-row {
        display: flex;
        align-items: baseline;
        gap: 0.75rem;
        margin-bottom: 1.25rem;
    }

    .price { font-size: 1.1rem; font-weight: 500; color: #111; }
    .stock { font-size: 0.8rem; color: #aaa; }
    .sold { font-size: 0.85rem; color: #999; font-style: italic; }

    .description {
        color: #555;
        font-size: 0.9rem;
        line-height: 1.7;
        margin: 0;
    }

    .status { color: #666; font-style: italic; }
    .error { color: #c00; }

    @media (max-width: 650px) {
        .detail { grid-template-columns: 1fr; }
    }
</style>
