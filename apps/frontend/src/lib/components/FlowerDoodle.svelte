<script lang="ts">
  interface Props {
    size?: number;
    class?: string;
  }

  let { size = 32, class: className = '' }: Props = $props();
</script>

<svg
  width={size}
  height={size}
  viewBox="0 0 64 64"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  stroke-linecap="round"
  class={className}
>
  <path d="M32 20c-4-8-14-10-16-4s6 10 16 16c-10 2-20 0-18 6s12 2 18-6c2 10 0 20 6 18s2-12-6-18c10-2 20 0 18-6s-12-2-18 6c4-10 6-18 0-18s-6 12 0 18" />
  <circle cx="32" cy="32" r="4" fill="currentColor" />
</svg>
