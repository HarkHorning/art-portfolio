<script lang="ts">
    import { default as ArtTile } from "$lib/components/artTile/ArtTile.svelte"
    import type { ArtTileInter } from "../artTile/ArtTileInterface";

    let tiles: ArtTileInter[] = $state([]);

    $effect(() => {
        (async () => {
            const res = await fetch(`/api/v1/art`);
            tiles = await res.json();
        })();
    });
</script>

<div class="grid-area">
    {#each tiles as tile (tile.id)}
        <ArtTile title={tile.title} url={tile.url} portrait={tile.portrait} />
    {/each}
</div>

<style>
    .grid-area {
        display: grid;
        grid-template-columns: repeat(3, 1fr);
    }

    @media screen and (max-width: 850px) {
        .grid-area {
            display: grid;
            grid-template-columns: repeat(2, 1fr);
        }
    }

    @media screen and (max-width: 450px) {
        .grid-area {
            display: grid;
            grid-template-columns: repeat(1, 1fr);
        }
    }
</style>
