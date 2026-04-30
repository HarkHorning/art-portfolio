<script lang="ts">
    let content = $state('');
    let photoUrl = $state('');

    $effect(() => {
        (async () => {
            try {
                const [contentRes, photoRes] = await Promise.all([
                    fetch('/api/v1/content/about_content'),
                    fetch('/api/v1/content/artist_photo_url'),
                ]);
                if (contentRes.ok) {
                    const data = await contentRes.json();
                    content = data.value ?? '';
                }
                if (photoRes.ok) {
                    const data = await photoRes.json();
                    photoUrl = data.value ?? '';
                }
            } catch {}
        })();
    });

    let paragraphs = $derived(
        content.split('\n').map(p => p.trim()).filter(p => p.length > 0)
    );
</script>

<svelte:head>
    <title>Hark Horning — About</title>
</svelte:head>

<div class="about-page">
    <div class="about-layout">
        {#if photoUrl}
            <div class="photo-col">
                <img src={photoUrl} alt="Hark Horning" class="artist-photo" />
            </div>
        {/if}
        <div class="text-col" class:no-photo={!photoUrl}>
            <h1>About Me</h1>
            {#each paragraphs as p}
                <p>{p}</p>
            {/each}
        </div>
    </div>
</div>

<style>
    .about-page {
        width: 100%;
        max-width: 1000px;
        padding-top: 2.5rem;
    }

    .about-layout {
        display: flex;
        gap: 5rem;
        align-items: flex-start;
    }

    .photo-col {
        flex-shrink: 0;
    }

    .artist-photo {
        width: 280px;
        height: 360px;
        object-fit: cover;
        border-radius: 6px;
        display: block;
    }

    .text-col {
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 1.4rem;
        padding-top: 0.5rem;
    }

    h1 {
        font-size: 1.6rem;
        font-weight: 400;
        letter-spacing: 0.02em;
        color: #1a1a1a;
        margin: 0 0 0.75rem;
    }

    p {
        color: #4a4a4a;
        line-height: 1.9;
        margin: 0;
        font-size: 1rem;
    }

    @media (max-width: 650px) {
        .about-layout {
            flex-direction: column;
            gap: 1.75rem;
        }

        .artist-photo {
            width: 100%;
            height: 300px;
        }
    }
</style>
