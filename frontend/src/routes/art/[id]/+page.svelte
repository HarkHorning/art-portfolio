<script lang="ts">
    import { page } from '$app/stores';
    import type { CategoryInter } from '$lib/components/filterSidebar/CategoryInterface';

    interface ArtDetail {
        id: number;
        title: string;
        description: string;
        portrait: boolean;
        url: string;
        categories: CategoryInter[];
    }

    let art: ArtDetail | null = $state(null);
    let loading = $state(true);
    let error: string | null = $state(null);

    $effect(() => {
        const id = $page.params.id;
        loading = true;
        error = null;

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
</script>

<div class="detail-page">
    <a href="/" class="back">← Back</a>

    {#if loading}
        <p class="status">Loading...</p>
    {:else if error}
        <p class="status error">{error}</p>
    {:else if art}
        <div class="detail" class:portrait={art.portrait} class:landscape={!art.portrait}>
            <div class="image-wrap">
                <img src={art.url} alt={art.title} />
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
        text-decoration: none;
        color: #888;
        font-size: 0.85rem;
        margin-bottom: 2rem;
        transition: color 0.2s;
    }

    .back:hover {
        color: #000;
    }

    /* Two-column layout for portrait pieces */
    .detail.portrait {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 3rem;
        align-items: start;
    }

    /* Single column for landscape pieces */
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
        margin: 0 0 1rem;
    }

    .categories {
        display: flex;
        flex-wrap: wrap;
        gap: 0.4rem;
        margin-bottom: 1.25rem;
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
        .detail.portrait {
            grid-template-columns: 1fr;
        }
    }
</style>
