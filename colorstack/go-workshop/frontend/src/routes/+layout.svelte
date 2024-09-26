<script lang="ts">
	import '../app.css';
	import { PanelLeft, House } from 'lucide-svelte';
	import type { QrCode } from '$lib/types';
	import type { LayoutData } from './$types';

	let sidebarOpen = true;
	export let data: LayoutData;
	const codes: QrCode[] = data.codes ?? [];
</script>

<div class="absolute top-4 left-4 z-20 flex gap-1.5">
	<button
		on:click={() => (sidebarOpen = !sidebarOpen)}
		class="p-1.5 rounded-lg hover:bg-neutral-300 transition-all ml-[-6px]"
	>
		<PanelLeft />
	</button>
	<a
		href="/"
		on:click={() => (sidebarOpen = false)}
		class="p-1.5 rounded-lg hover:bg-neutral-300 transition-all ml-[-6px]"
	>
		<House />
	</a>
</div>
<aside
	class="absolute top-0 left-[-16rem] h-screen w-64 bg-white shadow-lg z-10 transform transition-transform duration-300 pt-14 px-2"
	class:open={sidebarOpen}
>
	<div class="flex flex-col">
		{#each codes as code}
			<a
				on:click={() => {
					sidebarOpen = false;
				}}
				href={`/qr/${code.id}`}
				class="rounded-lg p-2 hover:bg-neutral-300 transform transition-all"
			>
				<h1 class="truncate text-left">{code.title}</h1>
			</a>
		{/each}
	</div>
</aside>
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

	.open {
		transform: translateX(16rem);
	}
</style>
