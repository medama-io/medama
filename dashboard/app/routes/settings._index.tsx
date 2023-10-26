import {
	type LoaderFunctionArgs,
	type MetaFunction,
	redirect,
} from '@remix-run/node';

import { hasSession } from '@/utils/cookies';

export const meta: MetaFunction = () => {
	return [
		{ title: 'Settings | Medama' },
		{ name: 'description', content: 'Privacy focused web analytics.' },
	];
};

export const loader = ({ request }: LoaderFunctionArgs) => {
	// Check for session cookie and redirect to login if missing
	if (!hasSession(request)) {
		throw redirect('/login');
	}

	// Otherwise redirect to first settings page
	return redirect('/settings/account');
};
