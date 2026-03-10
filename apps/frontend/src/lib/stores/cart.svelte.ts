export interface CartItem {
  id: string;
  title: string;
  price: number;
  size: string;
  image: string;
  quantity: number;
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
    return items.reduce((sum, item) => sum + item.quantity, 0);
  },

  get totalCents() {
    return items.reduce((sum, item) => sum + item.price * item.quantity, 0);
  },

  add(product: Omit<CartItem, 'quantity'>) {
    const existing = items.find((i) => i.id === product.id);
    if (existing) {
      existing.quantity += 1;
    } else {
      items.push({ ...product, quantity: 1 });
    }
    saveToStorage(items);
  },

  remove(id: string) {
    items = items.filter((i) => i.id !== id);
    saveToStorage(items);
  },

  updateQuantity(id: string, quantity: number) {
    if (quantity <= 0) {
      this.remove(id);
      return;
    }
    const item = items.find((i) => i.id === id);
    if (item) {
      item.quantity = quantity;
      saveToStorage(items);
    }
  },

  clear() {
    items = [];
    saveToStorage(items);
  },
};
