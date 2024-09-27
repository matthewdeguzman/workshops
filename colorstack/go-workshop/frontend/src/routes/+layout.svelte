<script lang="ts">
	import '../app.css';
	import { FilePenLine, Trash, PanelLeft, House, Heart } from 'lucide-svelte';
	import type { QrCode } from '$lib/types';
	import type { LayoutData } from './$types';
	import { sidebarOpen, qrCodes } from '$lib/stores';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	$: id = $page.params.id;

	export let data: LayoutData;
	$qrCodes = data.codes ?? [];

	async function handleDeletion(e: Event, code: QrCode) {
		e.stopPropagation();
		e.preventDefault();
		try {
			await fetch(`${import.meta.env.VITE_BASE_URL}/api/delete/${code.id}`, {
				method: 'DELETE',
				headers: {
					Authorization: 'Bearer 1234567890'
				}
			});
			if (id === code.id) goto('/');
			$qrCodes = $qrCodes.filter((c: QrCode) => c.id !== code.id);
		} catch (err) {
			console.error(err);
		}
	}
</script>

<svelte:head>
	<link rel="preconnect" href={import.meta.env.VITE_BASE_URL} />
</svelte:head>

<div class="absolute top-4 left-4 z-20 flex gap-1.5">
	<button on:click={() => ($sidebarOpen = !$sidebarOpen)} class="icon ml-[-6px]">
		<PanelLeft />
	</button>
	<a href="/" on:click={() => ($sidebarOpen = false)} class="icon ml-[-6px]">
		<House />
	</a>
</div>
<aside
	class="absolute top-0 left-[-16rem] h-screen w-64 bg-white shadow-lg z-10 transform transition-transform duration-300 pt-14 px-2"
	class:open={$sidebarOpen}
>
	<div class="flex flex-col">
		{#each $qrCodes as code (code.id)}
			<a
				on:click={() => {
					$sidebarOpen = false;
				}}
				href={`/qr/${code.id}`}
				class="group rounded-lg p-2 hover:bg-neutral-300 transform transition-all flex items-center relative"
			>
				<h1 class="text-left w-[calc(100%-56px)] truncate">
					{code.title}
				</h1>
				<div class="absolute right-1 opacity-100 flex gap-1">
					<a
						href={`/qr/${code.id}?edit=true`}
						class="transition-all hover:bg-sky-200 p-1 rounded-md"
					>
						<FilePenLine size={18} />
					</a>
					<button
						class="transition-all hover:bg-rose-200 p-1 rounded-md"
						on:click={(e) => handleDeletion(e, code)}
					>
						<Trash size={18} />
					</button>
				</div>
			</a>
		{/each}
	</div>
</aside>
<slot />
<section class="absolute bottom-1 w-full flex flex-col items-center">
	<h1>
		Made with <Heart size={16} fill="red" class="inline-block translate-y-[-2px] text-[red]" /> by Matthew
		DeGuzman
	</h1>
</section>

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

	:global(.icon) {
		@apply p-1.5 rounded-lg hover:bg-neutral-300 transition-all;
	}

	.open {
		transform: translateX(16rem);
	}
</style>
