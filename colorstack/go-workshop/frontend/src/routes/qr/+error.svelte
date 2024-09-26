<script lang="ts">
	import { page } from '$app/stores';
</script>

<section class="h-screen grid grid-cols-1 items-center justify-items-center">
	<div class="flex flex-col gap-2">
		<h1 class="title">{$page.error?.message || 'Uh oh! Something went wrong'}</h1>
		<a class="btn w-fit" href="/">Go back home</a>
	</div>
</section>
