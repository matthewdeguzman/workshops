import { writable } from 'svelte/store';
import type { QrCode } from './types';

export const qrCodes = writable<QrCode[]>([]);
export const sidebarOpen = writable(false);
