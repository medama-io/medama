import { Group, Stack } from '@mantine/core';

import { Sidebar } from './Sidebar';

import classes from './Layout.module.css';

interface SettingsLayoutProps {
	children: React.ReactNode;
}

export const SettingsLayout = ({ children }: SettingsLayoutProps) => {
	return (
		<Group className={classes.wrapper}>
			<Sidebar />
			<Stack flex={1}>{children}</Stack>
		</Group>
	);
};
