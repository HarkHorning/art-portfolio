<script lang="ts">
    import ArtGrid from "$lib/components/artGrid/ArtGrid.svelte";
    import FilterSidebar from "$lib/components/filterSidebar/FilterSidebar.svelte";
    import type { ArtTileInter } from "$lib/components/artTile/ArtTileInterface";
    import type { CategoryInter } from "$lib/components/filterSidebar/CategoryInterface";

    const priceRanges = [
        { min: -1,    max: -1    },
        { min: -1,    max: 2499  },
        { min: 2500,  max: 5000  },
        { min: 5001,  max: 10000 },
        { min: 10001, max: -1    },
    ];

    let tiles: ArtTileInter[] = $state([]);
    let categories: CategoryInter[] = $state([]);
    let sizes: string[] = $state([]);
    let loading = $state(true);
    let error: string | null = $state(null);
    let activeCategory: string | null = $state(null);
    let activeSize: string | null = $state(null);
    let activePriceRange = $state(0);
    let sidebarOpen = $state(true);

    $effect(() => {
        (async () => {
            try {
                const [catRes, sizeRes] = await Promise.all([
                    fetch('/api/v1/categories'),
                    fetch('/api/v1/art-sizes'),
                ]);
                if (catRes.ok) categories = await catRes.json();
                if (sizeRes.ok) sizes = await sizeRes.json();
            } catch {}
        })();
    });

    $effect(() => {
        const category = activeCategory;
        const size = activeSize;
        const range = priceRanges[activePriceRange];
        loading = true;
        error = null;

        (async () => {
            try {
                const params = new URLSearchParams();
                if (category) params.set('category', category);
                if (size) params.set('size', size);
                if (range.min >= 0) params.set('min_price', String(range.min));
                if (range.max >= 0) params.set('max_price', String(range.max));

                const url = `/api/v1/art${params.size ? '?' + params.toString() : ''}`;
                const res = await fetch(url);
                if (!res.ok) throw new Error();
                tiles = (await res.json()) ?? [];
            } catch {
                error = "Unable to load artwork. Please try again later.";
            } finally {
                loading = false;
            }
        })();
    });
</script>

<svelte:head>
    <title>Hark Horning — Gallery</title>
</svelte:head>

<div class="art-page">
    <div class="content" style="gap: {sidebarOpen ? '2rem' : '0'}">
        <FilterSidebar
            {categories}
            {sizes}
            active={activeCategory}
            {activeSize}
            {activePriceRange}
            open={sidebarOpen}
            onSelect={(slug) => activeCategory = slug}
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
            <ArtGrid {tiles} {loading} {error} />
        </div>
    </div>
</div>

<style>
    .art-page { width: 100%; }

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
