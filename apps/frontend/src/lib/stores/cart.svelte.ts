export interface CartItem {
  productId: string;
  variantId: string;
  title: string;
  priceCents: number;
  size: string;
  color: string;
  image: string;
}

const STORAGE_KEY = 'retrosnack-cart';

function loadFromStorage(): CartItem[] {
  if (typeof window === 'undefined') return [];
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    return raw ? JSON.parse(raw) : [];
  } catch {
    return [];
  }
}

function saveToStorage(items: CartItem[]) {
  if (typeof window === 'undefined') return;
  localStorage.setItem(STORAGE_KEY, JSON.stringify(items));
}

let items = $state<CartItem[]>(loadFromStorage());

export const cart = {
  get items() {
    return items;
  },

  get count() {
    return items.length;
  },

  get totalCents() {
    return items.reduce((sum, item) => sum + item.priceCents, 0);
  },

  has(productId: string) {
    return items.some((i) => i.productId === productId);
  },

  add(item: CartItem) {
    if (items.some((i) => i.variantId === item.variantId)) return;
    items.push(item);
    saveToStorage(items);
  },

  remove(variantId: string) {
    items = items.filter((i) => i.variantId !== variantId);
    saveToStorage(items);
  },

  clear() {
    items = [];
    saveToStorage(items);
  },
};
