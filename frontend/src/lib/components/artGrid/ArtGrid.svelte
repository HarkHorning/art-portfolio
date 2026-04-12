<script lang="ts">
    import ArtTile from "$lib/components/artTile/ArtTile.svelte";
    import type { ArtTileInter } from "../artTile/ArtTileInterface";

    let { tiles, loading, error }: {
        tiles: ArtTileInter[];
        loading: boolean;
        error: string | null;
    } = $props();
</script>

{#if loading}
    <p class="status">Loading...</p>
{:else if error}
    <p class="status error">{error}</p>
{:else if tiles.length === 0}
    <p class="status">No artwork to display.</p>
{:else}
    <div class="grid-area">
        {#each tiles as tile (tile.id)}
            <ArtTile id={tile.id} title={tile.title} url={tile.url} portrait={tile.portrait} />
        {/each}
    </div>
{/if}

<style>
    .grid-area {
        display: grid;
        grid-template-columns: repeat(3, 1fr);
        gap: 1rem;
    }

    .status {
        color: #666;
        font-style: italic;
    }

    .error {
        color: #c00;
    }

    @media screen and (max-width: 850px) {
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
