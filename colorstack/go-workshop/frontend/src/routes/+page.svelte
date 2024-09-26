<script lang="ts">
	import { goto } from '$app/navigation';
	import type { QrCode } from '$lib/types';

	const token = '1234567890';

	let form: HTMLFormElement;

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		const formData = new FormData(form);
		console.log(Object.fromEntries(formData));
		try {
			const resp = await fetch('http://localhost:8080/api/generate', {
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
			await goto(`/qr/${json.id}`);
		} catch (error) {
			console.error(error);
		}
	}
</script>

<section class="h-screen grid grid-cols-1 items-center justify-items-center">
	<form
		class="flex flex-col gap-2"
		action="http://localhost:8080/api/generate"
		method="POST"
		on:submit={handleSubmit}
		bind:this={form}
	>
		<h1 class="title">Awesome QR Code Generator</h1>
		<input class="event-input" type="text" name="url" placeholder="URL" required />
		<input class="event-input" type="text" name="title" placeholder="Event Title" />
		<input class="event-input" type="text" name="description" placeholder="Event Description" />
		<button type="submit" class="btn">Generate QR Code</button>
	</form>
</section>
