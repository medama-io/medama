const sortBy =
	// biome-ignore lint/suspicious/noExplicitAny: Generic function.
		<T extends Record<string, any>>(key: keyof T) =>
		(a: T, b: T) =>
			a[key] > b[key] ? 1 : b[key] > a[key] ? -1 : 0;

export { sortBy };
