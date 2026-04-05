<script lang="ts">

	let { children } = $props();
</script>

<svelte:head>
	<link href='https://fonts.googleapis.com/css?family=Inter' rel='stylesheet'>
</svelte:head>

{@render children()}
