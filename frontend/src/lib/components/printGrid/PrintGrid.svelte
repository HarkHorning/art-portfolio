<script lang="ts">
    import PrintTile from '$lib/components/printTile/PrintTile.svelte';
    import type { PrintTileInter } from '$lib/components/printTile/PrintTileInterface';

    let { prints, loading, error }: {
        prints: PrintTileInter[];
        loading: boolean;
        error: string | null;
    } = $props();
</script>

{#if loading}
    <p class="status">Loading...</p>
{:else if error}
    <p class="status error">{error}</p>
{:else if prints.length === 0}
    <p class="status">No prints available.</p>
{:else}
    <div class="grid-area">
        {#each prints as print (print.id)}
            <PrintTile
                id={print.id}
                title={print.title}
                url={print.url}
                portrait={print.portrait}
                sizes={print.sizes}
            />
        {/each}
    </div>
{/if}

<style>
    .grid-area {
        display: grid;
        grid-template-columns: repeat(4, 1fr);
        grid-auto-flow: dense;
        gap: 1rem;
    }

    .status {
        color: #666;
        font-style: italic;
    }

    .error {
        color: #c00;
    }

    @media screen and (max-width: 1100px) {
        .grid-area {
            grid-template-columns: repeat(3, 1fr);
        }
    }

    @media screen and (max-width: 750px) {
        .grid-area {
            grid-template-columns: repeat(2, 1fr);
        }
    }

    @media screen and (max-width: 450px) {
        .grid-area {
            grid-template-columns: repeat(1, 1fr);
        }
    }
</style>
