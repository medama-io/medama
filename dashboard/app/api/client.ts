import { json } from '@remix-run/node';

import { EXPIRE_COOKIE, expireSession } from '@/utils/cookies';

import { type components, type paths } from './types';

const LOCALHOST = 'http://localhost:8080';

const DEFAULT_HEADERS = {
	'Content-Type': 'application/json',
};

export type ComponentSchema = keyof components['schemas'];

export interface DataResponse<
	T extends ComponentSchema | undefined = ComponentSchema
> {
	cookie?: string;
	data?: T extends ComponentSchema ? components['schemas'][T] : undefined;
	res: Response;
}

// We also need to consider that some endpoints return an array of objects instead of a single object
export interface DataResponseArray<
	T extends ComponentSchema | undefined = ComponentSchema
> {
	cookie?: string;
	data?: T extends ComponentSchema
		? Array<components['schemas'][T]>
		: undefined;
	res: Response;
}

export interface ClientOptions<
	Body extends ComponentSchema | undefined = ComponentSchema
> {
	cookie?: string | null;
	body?: Body extends ComponentSchema ? components['schemas'][Body] : undefined;
	method?: 'GET' | 'POST' | 'PATCH' | 'DELETE';
	noRedirect?: boolean;
}

const client = async (
	path: keyof paths,
	{ cookie, body, method, noRedirect }: ClientOptions
): Promise<Response> => {
	const res = await fetch(`${LOCALHOST}${path}`, {
		method: method ?? 'GET',
		headers: {
			...DEFAULT_HEADERS,
			...(cookie !== undefined && cookie !== null && { Cookie: cookie }),
		},
		...(body !== undefined && { body: JSON.stringify(body) }),
	});

	if (!res.ok) {
		// If the user is not logged in, redirect to the login page
		if (res.status === 401 && !noRedirect) {
			throw expireSession();
		}

		// If it is 401 and noRedirect is true, do not throw anything but expire the invalid session
		if (res.status === 401 && noRedirect) {
			res.headers.set('Set-Cookie', EXPIRE_COOKIE);
			return res;
		}

		throw json(res.statusText, {
			status: res.status,
		});
	}

	return res;
};

export { client, DEFAULT_HEADERS, LOCALHOST };
