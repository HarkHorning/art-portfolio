<script lang="ts">
    interface Banner {
        id: number;
        art_tile_id: number;
        url: string;
        title: string;
        portrait: boolean;
    }

    interface Slide {
        banners: Banner[];
    }

    let { banners }: { banners: Banner[] } = $props();

    // Group banners into slides: portrait banners are paired, landscape get their own slide
    let slides = $derived.by(() => {
        const result: Slide[] = [];
        let i = 0;
        while (i < banners.length) {
            const b = banners[i];
            if (b.portrait && i + 1 < banners.length && banners[i + 1].portrait) {
                result.push({ banners: [b, banners[i + 1]] });
                i += 2;
            } else {
                result.push({ banners: [b] });
                i++;
            }
        }
        return result;
    });

    let current = $state(0);
    let paused = $state(false);

    function next() { current = (current + 1) % slides.length; }
    function prev() { current = (current - 1 + slides.length) % slides.length; }

    $effect(() => {
        if (slides.length <= 1) return;
        const id = setInterval(() => { if (!paused) next(); }, 5000);
        return () => clearInterval(id);
    });
</script>

{#if slides.length > 0}
<div
    class="hero"
    onmouseenter={() => paused = true}
    onmouseleave={() => paused = false}
    role="region"
    aria-label="Featured artwork"
>
    {#each slides as slide, i (i)}
        <div class="slide" class:active={i === current}>
            {#each slide.banners as banner (banner.id)}
                <a
                    href="/art/{banner.art_tile_id}"
                    class="panel"
                    tabindex={i === current ? 0 : -1}
                    aria-label={banner.title}
                >
                    <img src={banner.url} alt={banner.title} />
                </a>
            {/each}
        </div>
    {/each}

    {#if slides.length > 1}
        <button class="arrow left" onclick={prev} aria-label="Previous">‹</button>
        <button class="arrow right" onclick={next} aria-label="Next">›</button>

        <div class="dots">
            {#each slides as _, i}
                <button
                    class="dot"
                    class:active={i === current}
                    onclick={() => current = i}
                    aria-label="Go to slide {i + 1}"
                ></button>
            {/each}
        </div>
    {/if}
</div>
{/if}

<style>
    .hero {
        position: relative;
        width: 100%;
        height: 360px;
        border-radius: 8px;
        margin-bottom: 3rem;
        background: #fff;
    }

    .slide {
        position: absolute;
        inset: 0 160px;
        display: flex;
        gap: 0;
        opacity: 0;
        transition: opacity 0.7s ease;
        overflow: hidden;
    }

    .slide.active { opacity: 1; }

    .panel {
        flex: 1;
        display: block;
        overflow: hidden;
    }

    .panel + .panel {
        border-left: 2px solid #fff;
    }

    .panel img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        display: block;
        transition: transform 0.4s ease;
    }

    .panel:hover img { transform: scale(1.02); }

    .arrow {
        position: absolute;
        top: 50%;
        transform: translateY(-50%);
        background: rgba(255, 255, 255, 0.75);
        border: none;
        border-radius: 50%;
        width: 36px;
        height: 36px;
        font-size: 1.3rem;
        cursor: pointer;
        color: #111;
        transition: background 0.15s;
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 1;
    }

    .arrow:hover { background: rgba(255, 255, 255, 0.95); }
    .arrow.left  { left: 6px; }
    .arrow.right { right: 6px; }

    .dots {
        position: absolute;
        bottom: 0.75rem;
        left: 50%;
        transform: translateX(-50%);
        display: flex;
        gap: 0.4rem;
        z-index: 1;
    }

    .dot {
        width: 7px;
        height: 7px;
        border-radius: 50%;
        border: none;
        background: rgba(255, 255, 255, 0.5);
        cursor: pointer;
        padding: 0;
        transition: background 0.15s;
    }

    .dot.active { background: #fff; }

    @media (max-width: 1024px) {
        .slide { inset: 0 120px; }
    }

    @media (max-width: 768px) {
        .slide { inset: 0 80px; }
    }

    @media (max-width: 480px) {
        .hero { display: none; }
    }
</style>
