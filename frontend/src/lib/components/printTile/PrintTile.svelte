<script lang="ts">
    import type { PrintSizeInter } from './PrintTileInterface';

    let { id, title, url, portrait, sizes } : {
        id: number;
        title: string;
        url: string;
        portrait: boolean;
        sizes: PrintSizeInter[];
    } = $props();

    const availableSizes = sizes.filter(s => !s.sold && s.quantity_in_stock > 0);
</script>

<a href="/shop/{id}" class={portrait ? 'portrait' : 'landscape'} class:unavailable={availableSizes.length === 0}>
    <img src={url} alt={title} class="image" />
    <div class="meta">
        <h2>{title}</h2>
        <div class="details">
            {#if availableSizes.length === 0}
                <span class="sold-label">Out of stock</span>
            {/if}
            <!-- Price: re-enable by wrapping the block below in {#if true} -->
            {#if false}
                {#if fromPrice !== null}
                    <span class="price">from {formatPrice(fromPrice)}</span>
                    <span class="size-count">{availableSizes.length} {availableSizes.length === 1 ? 'size' : 'sizes'}</span>
                {/if}
            {/if}
        </div>
    </div>
</a>

<style>
    a {
        display: block;
        text-decoration: none;
        color: inherit;
        border: 1px solid #ccc;
        padding: 1rem;
        border-radius: 8px;
        transition: border-color 0.2s, box-shadow 0.2s;
    }

    a:hover {
        border-color: #999;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
    }

    a.unavailable {
        opacity: 0.6;
    }

    .landscape {
        grid-column: span 2;
    }

    .image {
        width: 100%;
        height: 90%;
        object-fit: cover;
        border-radius: 6px;
        display: block;
    }

    .meta {
        margin-top: 0.6rem;
        display: flex;
        justify-content: space-between;
        align-items: baseline;
        gap: 0.5rem;
    }

    h2 {
        margin: 0;
        font-size: 0.85rem;
        font-weight: 400;
        color: #333;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .details {
        display: flex;
        gap: 0.5rem;
        flex-shrink: 0;
        font-size: 0.8rem;
        align-items: baseline;
    }

    .price {
        color: #333;
        font-weight: 500;
    }

    .size-count {
        color: #aaa;
        font-size: 0.75rem;
    }

    .sold-label {
        color: #999;
        font-style: italic;
    }
</style>
