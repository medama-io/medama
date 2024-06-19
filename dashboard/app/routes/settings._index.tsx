import { userLoggedIn } from '@/api/user';
import { hasSession } from '@/utils/cookies';
import { type MetaFunction, redirect } from '@remix-run/react';

export const meta: MetaFunction = () => {
	return [
		{ title: 'Settings | Medama' },
		{ name: 'description', content: 'Privacy focused web analytics.' },
	];
};

export const clientLoader = async () => {
	await userLoggedIn();

	// Otherwise redirect to first settings page
	return redirect('/settings/account');
};
