import { Link, useLocation } from '@remix-run/react';
import classes from './Sidebar.module.css';

const SETTINGS_MAP = [
	{ label: 'Account', path: 'account' },
	{ label: 'Websites', path: 'websites' },
	{ label: 'Spam', path: 'spam' },
	{ label: 'Tracker', path: 'tracker' },
	{ label: 'Usage', path: 'usage' },
] as const;

export const Sidebar = () => {
	const { pathname } = useLocation();

	const options = SETTINGS_MAP.map((setting) => {
		const active = pathname.startsWith(`/settings/${setting.path}`);

		return (
			<Link
				key={setting.path}
				prefetch="intent"
				to={`/settings/${setting.path}`}
				data-active={active}
			>
				{setting.label}
			</Link>
		);
	});

	return <div className={classes.wrapper}>{options}</div>;
};
