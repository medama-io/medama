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

export const clientLoader = async (_: ClientLoaderFunctionArgs) => {
	return redirect('/settings/account');
};

export const ErrorBoundary = () => {
	return redirect('/settings/account');
};
