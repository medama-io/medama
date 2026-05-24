import { redirect } from 'react-router';
import type { Route } from './+types/settings._index';

export const meta: Route.MetaFunction = () => {
	return [
		{ title: 'Settings | Medama' },
		{ name: 'description', content: 'Privacy focused web analytics.' },
	];
};

export const clientLoader = async () => {
	return redirect('/settings/account');
};

export const ErrorBoundary = () => {
	return redirect('/settings/account');
};
