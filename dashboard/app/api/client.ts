import { json } from '@remix-run/react';

import { expireSession } from '@/utils/cookies';

import type { components, paths } from './types';

// Determine if we are running on the server (during SSR or pre-rendering) or in the browser.
const isServer = typeof window === 'undefined';
const LOCAL_API = 'http://localhost:8080/api';
const API_BASE = isServer ? 'localhost' : window.location.hostname;

// If we are running locally (development), use 'http://localhost:8080'. Otherwise, use the current origin.
const API_URL =
	isServer || API_BASE === 'localhost'
		? LOCAL_API
		: `${window.location.origin}/api`;

const DEFAULT_HEADERS = {
	'Content-Type': 'application/json',
};

export type ComponentSchema = keyof components['schemas'];

export interface DataResponse<
	T extends ComponentSchema | undefined = ComponentSchema,
> {
	data?: T extends ComponentSchema ? components['schemas'][T] : undefined;
	res: Response;
}

// We also need to consider that some endpoints return an array of objects instead of a single object
export interface DataResponseArray<
	T extends ComponentSchema | undefined = ComponentSchema,
> {
	data?: T extends ComponentSchema
		? Array<components['schemas'][T]>
		: undefined;
	res: Response;
}

export type ClientOptions<
	Body extends ComponentSchema | undefined = ComponentSchema,
> = Partial<{
	body: Body extends ComponentSchema ? components['schemas'][Body] : undefined;
	query: Record<string, string | number | boolean | undefined>;
	method: 'GET' | 'POST' | 'PATCH' | 'DELETE';
	shouldRedirect: boolean;
	shouldThrow: boolean;
	pathKey: string;
}>;

const client = async (
	path: keyof paths,
	{
		body,
		method = 'GET',
		shouldRedirect = true,
		shouldThrow = true,
		pathKey,
		query,
	}: ClientOptions,
): Promise<Response> => {
	const url = new URL(
		// Replace any path closed in curly braces with the pathKey.
		// e.g. /website/{hostname}/mediums -> /website/example.com/mediums
		`${API_URL}${pathKey ? path.replace(/{.*}/, pathKey) : path}`,
	);

	// Add the query to the path.
	if (query !== undefined) {
		for (const [key, value] of Object.entries(query)) {
			if (value !== undefined) {
				// Handle empty filter values.
				url.searchParams.append(
					key,
					value === 'Direct/None' ? '' : String(value),
				);
			}
		}
	}

	const res = await fetch(url, {
		method,
		headers: {
			...DEFAULT_HEADERS,
		},
		credentials: 'include',
		...(body !== undefined && { body: JSON.stringify(body) }),
	});

	if (!res.ok && shouldThrow) {
		// If the user is not logged in, redirect to the login page
		if (res.status === 401 && shouldRedirect) {
			throw expireSession(shouldRedirect);
		}

		// If it is 401 and noRedirect is true, do not throw anything but expire the invalid session
		if (res.status === 401 && !shouldRedirect) {
			expireSession(!shouldRedirect);
			return res;
		}

		throw json(res.statusText, {
			status: res.status,
		});
	}

	return res;
};

export { API_BASE, API_URL, DEFAULT_HEADERS, client };
