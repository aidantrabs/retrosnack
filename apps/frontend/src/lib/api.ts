const BASE_URL = import.meta.env.PUBLIC_API_URL ?? 'http://localhost:8080';

async function request<T>(path: string, init?: RequestInit): Promise<T> {
    const res = await fetch(`${BASE_URL}${path}`, {
        headers: { 'Content-Type': 'application/json', ...init?.headers },
        ...init,
    });
    if (!res.ok) {
        const body = await res.text();
        throw new Error(`HTTP ${res.status}: ${body}`);
    }
    return res.json() as Promise<T>;
}

export const api = {
    products: {
        list: (limit = 20, offset = 0) =>
            request<Product[]>(`/api/products?limit=${limit}&offset=${offset}`),
        get: (id: string) => request<Product>(`/api/products/${id}`),
        variants: (id: string) => request<Variant[]>(`/api/products/${id}/variants`),
    },
    categories: {
        list: () => request<Category[]>('/api/categories'),
    },
    inventory: {
        get: (variantId: string) => request<StockItem>(`/api/inventory/${variantId}`),
    },
    orders: {
        create: (items: OrderItemInput[]) =>
            request<Order>('/api/orders', {
                method: 'POST',
                body: JSON.stringify({ items }),
            }),
        get: (id: string) => request<Order>(`/api/orders/${id}`),
    },
    checkout: {
        create: (orderId: string) =>
            request<CheckoutSession>('/api/checkout', {
                method: 'POST',
                body: JSON.stringify({ order_id: orderId }),
            }),
    },
    payments: {
        config: () => request<PaymentConfig>('/api/payments/config'),
        process: (orderId: string, sourceId: string) =>
            request<PaymentResult>('/api/payments/process', {
                method: 'POST',
                body: JSON.stringify({ order_id: orderId, source_id: sourceId }),
            }),
    },
    drops: {
        list: () => request<Drop[]>('/api/drops'),
        get: (slug: string) => request<Drop>(`/api/drops/${slug}`),
        products: (slug: string) => request<Product[]>(`/api/drops/${slug}/products`),
    },
};

// types — kept in sync with go models

export interface Product {
    id: string;
    title: string;
    description: string;
    brand: string;
    condition: 'new' | 'excellent' | 'good' | 'fair';
    price_cents: number;
    category_id: string;
    seller_id: string | null;
    instagram_post_url: string;
    drop_id: string | null;
    notes: string;
    images: ProductImage[];
    drop?: Drop;
    created_at: string;
    updated_at: string;
}

export interface Drop {
    id: string;
    name: string;
    slug: string;
    description: string;
    instagram_url: string;
    released_at: string | null;
    created_at: string;
}

export interface ProductImage {
    id: string;
    product_id: string;
    url: string;
    position: number;
}

export interface Variant {
    id: string;
    product_id: string;
    size: string;
    color: string;
    sku: string;
    created_at: string;
}

export interface Category {
    id: string;
    name: string;
    slug: string;
    parent_id: string | null;
}

export interface StockItem {
    variant_id: string;
    quantity: number;
    reserved: number;
    available: number;
}

export interface Order {
    id: string;
    user_id: string | null;
    status: 'pending' | 'paid' | 'shipped' | 'delivered' | 'cancelled';
    total_cents: number;
    checkout_session_id: string | null;
    items: OrderItem[];
    created_at: string;
    updated_at: string;
}

export interface OrderItem {
    id: string;
    order_id: string;
    variant_id: string;
    quantity: number;
    price_cents: number;
}

export interface OrderItemInput {
    variant_id: string;
    quantity: number;
    price_cents: number;
}

export interface CheckoutSession {
    id: string;
    order_id: string;
    url: string;
}

export interface PaymentConfig {
    application_id: string;
    location_id: string;
    environment: string;
}

export interface PaymentResult {
    order_id: string;
    payment_id: string;
    status: string;
}
