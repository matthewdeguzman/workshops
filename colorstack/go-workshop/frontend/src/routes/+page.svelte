<script lang="ts">
	import { goto } from '$app/navigation';
	import type { QrCode } from '$lib/types';
	import { qrCodes, sidebarOpen } from '$lib/stores';

	const token = '1234567890';

	let form: HTMLFormElement;

	async function handleGenerate(e: SubmitEvent) {
		e.preventDefault();
		const formData = new FormData(form);
		try {
			const resp = await fetch(`${import.meta.env.VITE_BASE_URL}/api/generate`, {
				method: 'POST',
				headers: {
					Authorization: `Bearer ${token}`
				},
				body: formData
			});

			if (!resp.ok) {
				console.error(resp.statusText);
				return;
			}

			const json = (await resp.json()) as QrCode;
			qrCodes.update((codes) => [...codes, json]);

			$sidebarOpen = false;
			await goto(`/qr/${json.id}`);
		} catch (error) {
			console.error(error);
		}
	}
</script>

<section class="h-screen grid grid-cols-1 items-center justify-items-center relative">
	<div>
		<form
			class="flex flex-col gap-2 relative"
			action={`${import.meta.env.VITE_BASE_URL}/api/generate`}
			method="POST"
			on:submit={handleGenerate}
			bind:this={form}
		>
			<h1 class="title">Awesome QR Code Generator</h1>
			<input class="event-input" type="text" name="url" placeholder="URL" required />
			<input class="event-input" type="text" name="title" placeholder="Event Title" />
			<input class="event-input" type="text" name="description" placeholder="Event Description" />
			<button type="submit" class="btn">Generate QR Code</button>
		</form>
	</div>
</section>
