<script lang="ts">
    import { page } from '$app/stores';
</script>

<div class="error-page">
    <h1 class="code">{$page.status}</h1>

    {#if $page.status === 404}
        <p class="message">This page doesn't exist. Perhaps it was sold.</p>
    {:else}
        <p class="message">Something went wrong on our end.</p>
    {/if}

    <a href="/" class="home-link">← Back to the gallery</a>
</div>

<style>
    .error-page {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        min-height: 60vh;
        text-align: center;
        gap: 1rem;
    }

    .code {
        font-size: 6rem;
        font-weight: 700;
        margin: 0;
        color: #111;
        line-height: 1;
    }

    .message {
        font-size: 1.1rem;
        color: #555;
        margin: 0;
    }

    .home-link {
        margin-top: 1rem;
        color: #111;
        text-decoration: none;
        font-size: 0.95rem;
        border-bottom: 1px solid #111;
        padding-bottom: 1px;
    }

    .home-link:hover {
        color: #555;
        border-color: #555;
    }
</style>
