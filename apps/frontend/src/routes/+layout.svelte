<script lang="ts">
    import '../app.css';
    import Navbar from '$lib/components/Navbar.svelte';
    import Footer from '$lib/components/Footer.svelte';
    import BackgroundFlower from '$lib/components/BackgroundFlower.svelte';
    import Toast from '$lib/components/Toast.svelte';
    import { onMount } from 'svelte';

    let { children } = $props();

    onMount(async () => {
        if ('serviceWorker' in navigator) {
            const { registerSW } = await import('virtual:pwa-register');
            registerSW({ immediate: true });
        }
    });
</script>

<Toast />
<BackgroundFlower />
<div class="relative z-10 flex flex-col min-h-screen">
    <Navbar />
    <main class="flex-1" style="min-height: calc(100dvh - 49px)">
        {@render children()}
    </main>
    <Footer />
</div>
