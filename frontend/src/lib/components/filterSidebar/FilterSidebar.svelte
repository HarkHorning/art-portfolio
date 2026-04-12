<script lang="ts">
    import type { CategoryInter } from './CategoryInterface';

    let {
        categories,
        active,
        open,
        onSelect
    }: {
        categories: CategoryInter[];
        active: string | null;
        open: boolean;
        onSelect: (slug: string | null) => void;
    } = $props();
</script>

<aside class:open class:closed={!open}>
    <div class="content">
        <span class="label">Mediums</span>
        <ul>
            <li>
                <button
                    class:active={active === null}
                    onclick={() => onSelect(null)}
                >
                    All
                </button>
            </li>
            {#each categories as cat (cat.id)}
                <li>
                    <button
                        class:active={active === cat.slug}
                        onclick={() => onSelect(cat.slug)}
                    >
                        {cat.name}
                    </button>
                </li>
            {/each}
        </ul>
    </div>
</aside>

<style>
    aside {
        flex-shrink: 0;
        width: 140px;
        overflow: hidden;
        transition: width 0.2s ease;
        padding-top: 0.25rem;
    }

    aside.closed {
        width: 0;
        padding: 0;
    }

    .label {
        display: block;
        font-size: 0.7rem;
        letter-spacing: 0.1em;
        text-transform: uppercase;
        color: #999;
        margin-bottom: 0.75rem;
    }

    ul {
        list-style: none;
        margin: 0;
        padding: 0;
        display: flex;
        flex-direction: column;
        gap: 0.1rem;
    }

    li button {
        background: none;
        border: none;
        cursor: pointer;
        font-family: 'Inter', sans-serif;
        font-size: 0.85rem;
        color: #888;
        padding: 0.3rem 0;
        text-align: left;
        width: 100%;
        transition: color 0.15s;
        white-space: nowrap;
    }

    li button:hover {
        color: #000;
    }

    li button.active {
        color: #000;
        font-weight: 500;
    }
</style>
