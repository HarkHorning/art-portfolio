<script lang="ts">
    import ArtGrid from "$lib/components/artGrid/ArtGrid.svelte";
    import FilterSidebar from "$lib/components/filterSidebar/FilterSidebar.svelte";
    import type { ArtTileInter } from "$lib/components/artTile/ArtTileInterface";
    import type { CategoryInter } from "$lib/components/filterSidebar/CategoryInterface";

    let tiles: ArtTileInter[] = $state([]);
    let categories: CategoryInter[] = $state([]);
    let loading = $state(true);
    let error: string | null = $state(null);
    let activeCategory: string | null = $state(null);
    let sidebarOpen = $state(true);

    $effect(() => {
        (async () => {
            try {
                const res = await fetch('/api/v1/categories');
                if (res.ok) categories = await res.json();
            } catch {}
        })();
    });

    $effect(() => {
        const category = activeCategory;
        loading = true;
        error = null;

        (async () => {
            try {
                const url = category ? `/api/v1/art?category=${category}` : '/api/v1/art';
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

<div class="art-page">
    <div class="header">
        <h2 class="page-header">My work:</h2>
    </div>

    <div class="content">
        <FilterSidebar
            {categories}
            active={activeCategory}
            open={sidebarOpen}
            onSelect={(slug) => activeCategory = slug}
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
    .art-page {
        width: 100%;
    }

    .page-header {
        font-weight: 400;
    }

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
