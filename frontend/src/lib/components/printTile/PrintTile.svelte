<script lang="ts">
    let { id, title, url, portrait, price_cents, size, sold, quantity_in_stock } : {
        id: number;
        title: string;
        url: string;
        portrait: boolean;
        price_cents: number;
        size: string;
        sold: boolean;
        quantity_in_stock: number;
    } = $props();

    function formatPrice(cents: number): string {
        return `$${(cents / 100).toFixed(0)}`;
    }
</script>

<a href="/prints/{id}" class={portrait ? 'portrait' : 'landscape'} class:sold>
    <img src={url} alt={title} class="image" />
    <div class="meta">
        <h2>{title}</h2>
        <div class="details">
            <span class="size">{size}"</span>
            {#if sold}
                <span class="price sold-label">Sold</span>
            {:else if quantity_in_stock === 0}
                <span class="price sold-label">Out of stock</span>
            {:else}
                <span class="price">{formatPrice(price_cents)}</span>
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

    a.sold {
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
    }

    .size {
        color: #999;
    }

    .price {
        color: #333;
        font-weight: 500;
    }

    .sold-label {
        color: #999;
        font-style: italic;
        font-weight: 400;
    }
</style>
