<script lang="ts">
	import Navbar from '$lib/components/navbar/Navbar.svelte';
	import Footer from '$lib/components/footer/Footer.svelte';

	let { children } = $props();
</script>

<svelte:head>
	<link href='https://fonts.googleapis.com/css?family=Inter' rel='stylesheet'>
</svelte:head>

<Navbar />
<main>
	{@render children()}
</main>
<Footer />

<style>
	:global(body) {
		margin: 0;
		padding: 0;
		font-family: 'Inter', sans-serif;
	}

	main {
		padding: 0 2rem;
	}
</style>
