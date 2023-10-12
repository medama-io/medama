import {
	type LoaderFunctionArgs,
	type MetaFunction,
	redirect,
} from '@remix-run/node';

import { Login } from '@/components/login/Login';
import { getSession } from '@/utils/cookies';

export const meta: MetaFunction = () => {
	return [
		{ title: 'Login | Medama' },
		{ name: 'description', content: 'Login into Medama Analytics.' },
	];
};

export const loader = ({ request }: LoaderFunctionArgs) => {
	// If the user is already logged in, redirect them to the dashboard.
	if (getSession(request)) {
		throw redirect('/');
	}

	return { status: 200 };
};

export default function Index() {
	return <Login />;
}
