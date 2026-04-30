<script lang="ts">
    import PrintGrid from '$lib/components/printGrid/PrintGrid.svelte';
    import PrintFilter from '$lib/components/printFilter/PrintFilter.svelte';
    import type { PrintTileInter } from '$lib/components/printTile/PrintTileInterface';

    const priceRanges = [
        { min: -1,    max: -1    },
        { min: -1,    max: 2499  },
        { min: 2500,  max: 5000  },
        { min: 5001,  max: 10000 },
        { min: 10001, max: -1    },
    ];

    let prints: PrintTileInter[] = $state([]);
    let sizes: string[] = $state([]);
    let loading = $state(true);
    let error: string | null = $state(null);
    let activeSize: string | null = $state(null);
    let activePriceRange = $state(0);
    let sidebarOpen = $state(true);

    $effect(() => {
        (async () => {
            try {
                const res = await fetch('/api/v1/print-sizes');
                if (res.ok) sizes = await res.json();
            } catch {}
        })();
    });

    $effect(() => {
        const size = activeSize;
        const range = priceRanges[activePriceRange];
        loading = true;
        error = null;

        (async () => {
            try {
                const params = new URLSearchParams();
                if (size) params.set('size', size);
                if (range.min >= 0) params.set('min_price', String(range.min));
                if (range.max >= 0) params.set('max_price', String(range.max));

                const url = `/api/v1/prints${params.size ? '?' + params.toString() : ''}`;
                const res = await fetch(url);
                if (!res.ok) throw new Error();
                prints = (await res.json()) ?? [];
            } catch {
                error = 'Unable to load prints. Please try again later.';
            } finally {
                loading = false;
            }
        })();
    });
</script>

<svelte:head>
    <title>Hark Horning — Shop</title>
</svelte:head>

<div class="shop-page">
    <div class="content" style="gap: {sidebarOpen ? '2rem' : '0'}">
        <PrintFilter
            {sizes}
            {activeSize}
            {activePriceRange}
            open={sidebarOpen}
            onSizeSelect={(size) => activeSize = size}
            onPriceSelect={(i) => activePriceRange = i}
        />
        <div class="grid-wrap">
            <button
                class="filter-toggle"
                onclick={() => sidebarOpen = !sidebarOpen}
            >
                {sidebarOpen ? '‹ Filters' : 'Filters ›'}
            </button>
            <PrintGrid {prints} {loading} {error} />
        </div>
    </div>
</div>

<style>
    .shop-page { width: 100%; }

    .content {
        display: flex;
        gap: 2rem;
        align-items: flex-start;
    }

    .grid-wrap {
        flex: 1;
        min-width: 0;
    }

    .filter-toggle {
        background: none;
        border: none;
        border-bottom: 1px solid #ccc;
        cursor: pointer;
        font-family: 'Inter', sans-serif;
        font-size: 0.8rem;
        letter-spacing: 0.04em;
        color: #555;
        padding: 0 0 0.2rem 0;
        margin-bottom: 1.25rem;
        display: block;
        transition: color 0.2s, border-color 0.2s;
    }

    .filter-toggle:hover {
        color: #000;
        border-color: #000;
    }
</style>
