<script lang="ts">
    // const title = "Artwork Name"
    // export const title = "Artwork Name";
    let { title, url, portrait } : { title: string, url: string, portrait: boolean } = $props()
</script>

<div class={portrait ? 'portrait' : 'landscape'}>
    <img src={url} alt="Art Piece" class="image" />
    <h2>{title}</h2>
</div>

<style>
    .portrait { border: 1px solid #ccc; padding: 1rem; border-radius: 8px; }
    .landscape { border: 1px solid #ccc; padding: 1rem; border-radius: 8px; grid-column: span 2; }
    .image { 
        width: 100%;
        height: 90%;
        object-fit: cover;
        /* object-fit: contain; */
        border-radius: 6px;
    }
</style>