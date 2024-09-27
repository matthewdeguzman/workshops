<script lang="ts">
	import { page } from '$app/stores';
	import { quadOut } from 'svelte/easing';
	import { goto } from '$app/navigation';
	import { fly, fade } from 'svelte/transition';
	import type { PageServerData } from './$types';
	import { FilePenLine } from 'lucide-svelte';

	export let data: PageServerData;
	$: code = data.code;
	$: id = $page.params.id;
	$: editing = $page.url.searchParams.get('edit') === 'true';
	let formElement: HTMLFormElement;

	async function handleUpdate(e: SubmitEvent) {
		e.preventDefault();
		const formData = new FormData(formElement);
		try {
			await fetch(`http://localhost:8080/api/update/`, {
				method: 'PUT',
				headers: {
					Authorization: 'Bearer 1234567890'
				},
				body: formData
			});
		} catch (error) {
			console.error(error);
		}
		await goto(`/qr/${id}?edit=false`);
		editing = false;
	}
</script>

<section class="h-screen grid grid-cols-1 items-center justify-items-center bg-white">
	<button
		class="icon absolute top-4 right-4"
		on:click={() => {
			goto(`/qr/${id}?edit=true`);
			editing = true;
		}}
	>
		<FilePenLine />
	</button>
	{#if data}
		{#if !editing}
			<div
				class="flex flex-col items-center gap-2 absolute"
				in:fly={{ delay: 300, duration: 250, y: 100, easing: quadOut }}
				out:fly={{ delay: 250, duration: 250, y: -100, easing: quadOut }}
			>
				<h1 class="text-center text-6xl font-bold bg-transparent w-full">{code.title}</h1>
				<p class="text-center text-2xl bg-transparent w-full">{code.description}</p>
				<img
					class="rounded-lg"
					src={`http://localhost:8080/qr/${id}`}
					alt="QR code"
					width="512"
					height="512"
				/>
			</div>
		{:else}
			<form
				class="flex flex-col items-center gap-2 absolute max-w-[600px] w-full"
				bind:this={formElement}
				in:fly={{ delay: 500, duration: 250, y: 100, easing: quadOut }}
				out:fade={{ duration: 300, easing: quadOut }}
				on:submit={handleUpdate}
			>
				<div class="w-full">
					<label for="title" class="font-bold">Event Title </label>
					<input
						class="event-input"
						type="text"
						name="title"
						placeholder="Event Title"
						bind:value={code.title}
					/>
				</div>
				<div class="w-full">
					<label for="description" class="font-bold">Event Description</label>
					<input
						class="event-input"
						type="text"
						name="description"
						placeholder="Event Description"
						bind:value={code.description}
					/>
				</div>
				<input class="event-input" type="hidden" name="id" value={id} />
				<div class="flex gap-4">
					<button
						class="btn !bg-transparent hover:!bg-neutral-300 border border-solid border-black !text-black"
						on:click={async () => {
							await goto(`/qr/${id}?edit=false`);
							editing = false;
						}}>Cancel</button
					>
					<button class="btn" type="submit">Save Event</button>
				</div>
			</form>
		{/if}
	{:else}
		<h1>Loading...</h1>
	{/if}
</section>
