<script>
	import '../app.css';
</script>

<slot />

<style lang="postcss">
	:global(.title) {
		@apply text-4xl font-bold;
	}

	:global(.btn) {
		@apply bg-black text-white p-2 rounded-lg hover:bg-neutral-600 transition-all duration-300;
	}

	:global(.event-input) {
		@apply border border-gray-300 p-2 rounded-lg w-full;
	}
</style>
