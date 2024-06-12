import {
	type ClientLoaderFunctionArgs,
	useLoaderData,
	json,
	type MetaFunction,
	redirect,
} from '@remix-run/react';

import type { components } from '@/api/types';
import { userGet } from '@/api/user';
import { Settings } from '@/components/settings/Settings';
import { hasSession } from '@/utils/cookies';

interface LoaderData {
	user: components['schemas']['UserGet'];
}

export const meta: MetaFunction = () => {
	return [
		{ title: 'Settings | Medama' },
		{ name: 'description', content: 'Privacy focused web analytics.' },
	];
};

const ACCEPTED_SETTINGS = new Set(['account', 'advanced']);

export const clientLoader = async ({ request }: ClientLoaderFunctionArgs) => {
	// If pathname does not match accepted settings pages, 404
	const url = new URL(request.url);
	const pathname = url.pathname.replace('/settings', '');
	if (!ACCEPTED_SETTINGS.has(pathname.replace('/', ''))) {
		throw json('Not found', {
			status: 404,
		});
	}

	// Check for session cookie and redirect to login if missing
	if (!hasSession()) {
		throw redirect('/login');
	}

	const { data } = await userGet();

	if (!data) {
		throw json('Failed to get user.', {
			status: 500,
		});
	}

	return json<LoaderData>({
		user: data,
	});
};

export default function Index() {
	const { user } = useLoaderData<LoaderData>();

	if (!user) {
		return;
	}

	return (
		<main>
			<Settings user={user} />
		</main>
	);
}
