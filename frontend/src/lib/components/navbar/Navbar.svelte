<script lang="ts">
    import { page } from '$app/stores';
</script>

<nav>
    <a href="/" class="site-title">Hark Horning</a>
    <div class="links">
        <a href="/gallery" class:active={$page.url.pathname.startsWith('/gallery') || $page.url.pathname.startsWith('/art')}>Gallery</a>
        <a href="/shop" class:active={$page.url.pathname.startsWith('/shop')}>Shop</a>
        <a href="/about" class:active={$page.url.pathname === '/about'}>About</a>
    </div>
</nav>

<style>
    nav {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1.5rem 2rem;
        font-family: 'Inter', sans-serif;
    }

    .site-title {
        text-decoration: none;
        color: #000;
        font-size: 0.9rem;
        font-weight: 500;
        letter-spacing: 0.02em;
    }

    .links {
        display: flex;
        gap: 2rem;
    }

    a {
        text-decoration: none;
        color: #666;
        font-size: 0.9rem;
        letter-spacing: 0.05em;
        transition: color 0.2s;
    }

    a:hover { color: #000; }
    a.active { color: #000; }
</style>
