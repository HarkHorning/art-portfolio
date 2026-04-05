<script lang="ts">
    import { default as ArtGrid } from "$lib/components/artGrid/ArtGrid.svelte"

</script>

<div class="art-page">
    <h1>I'm Hark,</h1>
    <p>welcome to my portfolio</p>
    <div>
        <ArtGrid />
    </div>
</div>

<style>
    .art-page {
        font-family: 'CustomFont', sans-serif;
        width: full;
    }
</style>
