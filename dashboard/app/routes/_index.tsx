import {
	type LoaderFunctionArgs,
	type MetaFunction,
	redirect,
} from '@remix-run/node';

import { getSession } from '@/utils/cookies';

export const meta: MetaFunction = () => {
	return [
		{ title: 'Medama | Privacy Focused Web Analytics.' },
		{ name: 'description', content: 'Privacy focused web analytics.' },
	];
};

export const loader = ({ request }: LoaderFunctionArgs) => {
	// Check for session cookie and redirect to login if missing
	if (!getSession(request)) {
		throw redirect('/login');
	}

	return { status: 200 };
};

export default function Index() {
	return (
		<div>
			<h1>Homepage</h1>
		</div>
	);
}
