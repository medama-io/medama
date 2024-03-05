import { Button, Grid, Group, Stack, Text, TextInput } from '@mantine/core';

import { type components } from '@/api/types';

import classes from './Settings.module.css';
import { Sidebar } from './Sidebar';

interface SettingsProps {
	user: components['schemas']['UserGet'];
}

export const Settings = ({ user }: SettingsProps) => {
	return (
		<Stack className={classes.wrapper}>
			<Text fw={700} fz={32} p={24}>
				Settings
			</Text>
			<Grid>
				<Grid.Col span={4}>
					<Sidebar />
				</Grid.Col>
				<Grid.Col span={8} py="xl">
					<Group>
						<Stack justify="flex-start" align="flex-start">
							<Text size="xl" fw={700}>
								Account details
							</Text>
							<Text size="sm">Edit your username and password.</Text>
						</Stack>
						<Stack>
							<TextInput
								label="Username"
								placeholder="Username"
								value={user.username}
							/>
							<TextInput
								label="Password"
								placeholder="Password"
								type="password"
							/>
						</Stack>
					</Group>
					<Button variant="light" color="blue">
						Save
					</Button>
				</Grid.Col>
			</Grid>
		</Stack>
	);
};
