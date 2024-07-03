import { Stack, UnstyledButton } from '@mantine/core';

import classes from './Sidebar.module.css';
import { Link, useLocation } from '@remix-run/react';

const SETTINGS_MAP = [
	{ label: 'Account', path: 'account' },
	{ label: 'Usage', path: 'usage' },
] as const;

export const Sidebar = () => {
	const { pathname } = useLocation();

	const options = SETTINGS_MAP.map((setting) => {
		const active = pathname.startsWith(`/settings/${setting.path}`);

		return (
			<UnstyledButton
				key={setting.path}
				component={Link}
				to={`/settings/${setting.path}`}
				data-active={active}
			>
				{setting.label}
			</UnstyledButton>
		);
	});

	return <Stack className={classes.wrapper}>{options}</Stack>;
};
