<script lang="ts">
  interface Props {
    id: string;
    title: string;
    price: number;
    size: string;
    image: string;
  }

  let { id, title, price, size, image }: Props = $props();
</script>

<a href="/shop/{id}" class="group block">
  <div class="aspect-[3/4] overflow-hidden rounded-lg bg-sand-dark">
    <img
      src={image}
      alt={title}
      class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
    />
  </div>

  <div class="mt-3 space-y-1">
    <h3 class="text-sm font-medium truncate">{title}</h3>
    <div class="flex items-center justify-between text-sm">
      <span class="text-ink-muted">{size}</span>
      <span class="font-display font-semibold">${(price / 100).toFixed(2)}</span>
    </div>
  </div>
</a>
