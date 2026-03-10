export interface MockProduct {
  id: string;
  title: string;
  price: number;
  size: string;
  image: string;
}

export const mockProducts: MockProduct[] = [
  {
    id: '1',
    title: 'vintage denim jacket',
    price: 4500,
    size: 'M',
    image: 'https://placehold.co/600x800/CFC0B2/1A1A1A?text=denim+jacket',
  },
  {
    id: '2',
    title: 'floral midi skirt',
    price: 2800,
    size: 'S',
    image: 'https://placehold.co/600x800/CFC0B2/1A1A1A?text=midi+skirt',
  },
  {
    id: '3',
    title: 'corduroy overalls',
    price: 5200,
    size: 'L',
    image: 'https://placehold.co/600x800/CFC0B2/1A1A1A?text=overalls',
  },
  {
    id: '4',
    title: 'knit cardigan',
    price: 3400,
    size: 'M',
    image: 'https://placehold.co/600x800/CFC0B2/1A1A1A?text=cardigan',
  },
  {
    id: '5',
    title: 'retro band tee',
    price: 2200,
    size: 'S',
    image: 'https://placehold.co/600x800/CFC0B2/1A1A1A?text=band+tee',
  },
  {
    id: '6',
    title: 'pleated trousers',
    price: 3800,
    size: 'M',
    image: 'https://placehold.co/600x800/CFC0B2/1A1A1A?text=trousers',
  },
  {
    id: '7',
    title: 'leather crossbody bag',
    price: 4000,
    size: 'OS',
    image: 'https://placehold.co/600x800/CFC0B2/1A1A1A?text=crossbody',
  },
  {
    id: '8',
    title: 'wool blend coat',
    price: 6800,
    size: 'L',
    image: 'https://placehold.co/600x800/CFC0B2/1A1A1A?text=wool+coat',
  },
];
