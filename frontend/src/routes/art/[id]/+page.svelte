<script lang="ts">
    import { page } from '$app/stores';
    import type { CategoryInter } from '$lib/components/filterSidebar/CategoryInterface';

    interface ImageInter {
        id: number;
        variant: string;
        url: string;
        sort_order: number;
    }

    interface ArtDetail {
        id: number;
        title: string;
        description: string;
        portrait: boolean;
        categories: CategoryInter[];
        images: ImageInter[];
        size?: string;
        price_cents?: number;
        sold: boolean;
    }

    function formatPrice(cents: number): string {
        return `$${(cents / 100).toFixed(0)}`;
    }

    let art: ArtDetail | null = $state(null);
    let loading = $state(true);
    let error: string | null = $state(null);
    let selectedIndex = $state(0);

    let displayImages = $derived(
        art ? art.images.filter(img => img.variant === 'low').sort((a, b) => a.sort_order - b.sort_order) : []
    );

    let currentImage = $derived(displayImages[selectedIndex] ?? null);

    $effect(() => {
        const id = $page.params.id;
        loading = true;
        error = null;
        selectedIndex = 0;

        (async () => {
            try {
                const res = await fetch(`/api/v1/art/${id}`);
                if (!res.ok) throw new Error();
                art = await res.json();
            } catch {
                error = 'Could not load this piece.';
            } finally {
                loading = false;
            }
        })();
    });

    function prev() {
        selectedIndex = (selectedIndex - 1 + displayImages.length) % displayImages.length;
    }

    function next() {
        selectedIndex = (selectedIndex + 1) % displayImages.length;
    }
</script>

<svelte:head>
    <title>{art ? `${art.title} — Hark Horning` : 'Hark Horning'}</title>
</svelte:head>

<div class="detail-page">
    <button onclick={() => history.back()} class="back">← Back</button>

    {#if loading}
        <p class="status">Loading...</p>
    {:else if error}
        <p class="status error">{error}</p>
    {:else if art}
        <div class="detail">
            <div class="image-wrap">
                {#if currentImage}
                    <img src={currentImage.url} alt={art.title} />
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
                <h1>{art.title}</h1>
                {#if art.categories.length > 0}
                    <div class="categories">
                        {#each art.categories as cat (cat.id)}
                            <span class="tag">{cat.name}</span>
                        {/each}
                    </div>
                {/if}
                {#if art.size || art.price_cents != null}
                    <div class="art-meta">
                        {#if art.size}<span class="size">{art.size}"</span>{/if}
                        {#if art.sold}
                            <span class="sold">Sold</span>
                        {:else if art.price_cents != null}
                            <span class="price">{formatPrice(art.price_cents)}</span>
                        {/if}
                    </div>
                {/if}
                {#if art.description}
                    <p class="description">{art.description}</p>
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

    .thumb.active {
        border-color: #111;
    }

    .info {
        padding-top: 0.5rem;
    }

    h1 {
        font-size: 1.5rem;
        font-weight: 400;
        margin: 0 0 1rem;
    }

    .categories {
        display: flex;
        flex-wrap: wrap;
        gap: 0.4rem;
        margin-bottom: 1.25rem;
    }

    .art-meta {
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

    .tag {
        font-size: 0.75rem;
        letter-spacing: 0.04em;
        color: #666;
        border: 1px solid #ddd;
        border-radius: 3px;
        padding: 0.2rem 0.6rem;
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
        .detail {
            grid-template-columns: 1fr;
        }
    }
</style>
