import { Button, Stack } from '@mantine/core';

import classes from './Sidebar.module.css';

export const Sidebar = () => {
	return (
		<Stack className={classes.wrapper}>
			<Button>Account</Button>
			<Button>Advanced</Button>
		</Stack>
	);
};
