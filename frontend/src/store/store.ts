import {writable} from "svelte/store";

export const activePage = writable<string>('home');