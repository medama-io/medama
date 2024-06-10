import { hasSession } from '@/utils/cookies';
import {
	type ClientLoaderFunctionArgs,
	type MetaFunction,
	redirect,
} from '@remix-run/react';

export const meta: MetaFunction = () => {
	return [
		{ title: 'Settings | Medama' },
		{ name: 'description', content: 'Privacy focused web analytics.' },
	];
};

export const clientLoader = ({ request }: ClientLoaderFunctionArgs) => {
	// Check for session cookie and redirect to login if missing
	if (!hasSession(request)) {
		throw redirect('/login');
	}

	// Otherwise redirect to first settings page
	return redirect('/settings/account');
};
