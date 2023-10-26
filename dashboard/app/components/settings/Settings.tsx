import { Box, Grid, Group, Stack } from '@mantine/core';

import { type components } from '@/api/types';

import classes from './Settings.module.css';
import { Sidebar } from './Sidebar';

interface SettingsProps {
	user: components['schemas']['UserGet'];
}

export const Settings = ({ user }: SettingsProps) => {
	return (
		<Stack className={classes.wrapper}>
			<h1>Settings</h1>
			<Grid>
				<Grid.Col span={4}>
					<Sidebar />
				</Grid.Col>
				<Grid.Col span={8}>
					<Group>
						<Box>Account</Box>
						<Box>Site</Box>
						{JSON.stringify(user)}
					</Group>
				</Grid.Col>
			</Grid>
		</Stack>
	);
};
