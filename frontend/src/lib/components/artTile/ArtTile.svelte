<script lang="ts">
    let { id, title, url, portrait } : { id: number, title: string, url: string, portrait: boolean } = $props()
</script>

<a href="/art/{id}" class={portrait ? 'portrait' : 'landscape'}>
    <img src={url} alt={title} class="image" />
    <h2>{title}</h2>
</a>

<style>
    a {
        display: block;
        text-decoration: none;
        color: inherit;
        border: 1px solid #ccc;
        padding: 1rem;
        border-radius: 8px;
        transition: border-color 0.2s, box-shadow 0.2s;
    }

    a:hover {
        border-color: #999;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
    }

    .landscape {
        grid-column: span 2;
    }

    .image {
        width: 100%;
        height: 90%;
        object-fit: cover;
        border-radius: 6px;
        display: block;
    }

    h2 {
        margin: 0.6rem 0 0;
        font-size: 0.85rem;
        font-weight: 400;
        color: #333;
    }
</style>
