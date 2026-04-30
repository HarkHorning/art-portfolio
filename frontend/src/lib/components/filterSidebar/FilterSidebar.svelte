<script lang="ts">
    import type { CategoryInter } from './CategoryInterface';

    let {
        categories,
        sizes,
        active,
        activeSize,
        activePriceRange,
        open,
        onSelect,
        onSizeSelect,
        onPriceSelect,
    }: {
        categories: CategoryInter[];
        sizes: string[] | null;
        active: string | null;
        activeSize: string | null;
        activePriceRange: number;
        open: boolean;
        onSelect: (slug: string | null) => void;
        onSizeSelect: (size: string | null) => void;
        onPriceSelect: (index: number) => void;
    } = $props();

    const priceRanges = [
        { label: 'All prices',   min: -1,    max: -1    },
        { label: 'Under $25',    min: -1,    max: 2499  },
        { label: '$25 – $50',    min: 2500,  max: 5000  },
        { label: '$50 – $100',   min: 5001,  max: 10000 },
        { label: '$100+',        min: 10001, max: -1    },
    ];
</script>

<aside class:open class:closed={!open}>
    <div class="content">
        <span class="label">Medium</span>
        <ul>
            <li>
                <button class:active={active === null} onclick={() => onSelect(null)}>All</button>
            </li>
            {#each categories as cat (cat.id)}
                <li>
                    <button class:active={active === cat.slug} onclick={() => onSelect(cat.slug)}>
                        {cat.name}
                    </button>
                </li>
            {/each}
        </ul>

        {#if sizes && sizes.length > 0}
            <span class="label" style="margin-top: 1.25rem;">Size</span>
            <ul>
                <li>
                    <button class:active={activeSize === null} onclick={() => onSizeSelect(null)}>All</button>
                </li>
                {#each sizes as size}
                    <li>
                        <button class:active={activeSize === size} onclick={() => onSizeSelect(size)}>
                            {size}"
                        </button>
                    </li>
                {/each}
            </ul>
        {/if}

        <!-- Price filter: re-enable by changing {#if false} to {#if true} -->
        {#if false}
        <span class="label" style="margin-top: 1.25rem;">Price</span>
        <ul>
            {#each priceRanges as range, i}
                <li>
                    <button class:active={activePriceRange === i} onclick={() => onPriceSelect(i)}>
                        {range.label}
                    </button>
                </li>
            {/each}
        </ul>
        {/if}
    </div>
</aside>

<style>
    aside {
        flex-shrink: 0;
        width: 140px;
        overflow: hidden;
        transition: width 0.2s ease;
        padding-top: 0.25rem;
    }

    aside.closed {
        width: 0;
        padding: 0;
    }

    .label {
        display: block;
        font-size: 0.7rem;
        letter-spacing: 0.1em;
        text-transform: uppercase;
        color: #999;
        margin-bottom: 0.75rem;
    }

    ul {
        list-style: none;
        margin: 0;
        padding: 0;
        display: flex;
        flex-direction: column;
        gap: 0.1rem;
    }

    li button {
        background: none;
        border: none;
        cursor: pointer;
        font-family: 'Inter', sans-serif;
        font-size: 0.85rem;
        color: #888;
        padding: 0.3rem 0;
        text-align: left;
        width: 100%;
        transition: color 0.15s;
        white-space: nowrap;
    }

    li button:hover { color: #000; }
    li button.active { color: #000; font-weight: 500; }
</style>
